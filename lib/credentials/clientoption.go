package credentials

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
)

// FederationOptions configures workload identity federation for a client.
type FederationOptions struct {
	// FederationRuleID is the federation rule ("wifr_...") to check the
	// identity token against. Required.
	FederationRuleID string
	// Resource is the API the token is requested for. It defaults to the
	// origin of the client's base URL, such as "https://api.onkernel.com".
	Resource string
	// AuthBaseURL is the Kernel authorization server. It defaults to
	// KERNEL_AUTH_BASE_URL, then to https://auth.onkernel.com when the client
	// talks to the production API. A client with any other base URL must set
	// it.
	AuthBaseURL string
}

// ClientOption returns the request option behind
// [github.com/kernel/kernel-go-sdk/option.WithFederationTokenProvider]; see
// there for how it behaves.
func ClientOption(provider IdentityTokenFunc, opts FederationOptions) requestconfig.RequestOption {
	o := &clientOption{identityTokenProvider: provider, opts: opts}
	switch {
	case provider == nil:
		o.err = configError("WithFederationTokenProvider: the identity token provider is nil")
	case opts.FederationRuleID == "":
		o.err = configError("WithFederationTokenProvider: FederationRuleID is required")
	}
	return o
}

// EnvClientOption returns the request option for the federation settings in
// the environment (see [FederationFromEnv]). It does nothing when
// KERNEL_API_KEY is set, which takes precedence, or when
// KERNEL_FEDERATION_RULE_ID is not. A federation setup that is incomplete
// fails the requests that would have used it.
func EnvClientOption() requestconfig.RequestOption {
	if apiKeyWinsInEnv() {
		return requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error { return nil })
	}
	provider, opts, err := FederationFromEnv()
	if provider == nil && err == nil {
		return requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error { return nil })
	}
	return &clientOption{identityTokenProvider: provider, opts: opts, err: err}
}

// clientOption authenticates a client's requests with federation. One value
// serves every request of the clients it was given to, so the token caches
// live here: one per API base URL, because a token is never reused across
// APIs.
type clientOption struct {
	identityTokenProvider IdentityTokenFunc
	opts                  FederationOptions
	err                   error

	mu     sync.Mutex
	caches map[string]*TokenCache
}

type activeCredentialKey struct{}

// Apply makes o the credential of the request. Credentials follow the usual
// rule for options: the last one wins. An API key or Authorization header
// applied before o is cleared here, and one applied after o makes the
// middleware stand down, as does a later federation option.
func (o *clientOption) Apply(r *requestconfig.RequestConfig) error {
	r.APIKey = ""
	r.Request.Header.Del("Authorization")
	r.Context = context.WithValue(r.Context, activeCredentialKey{}, o)
	r.Middlewares = append(r.Middlewares, o.middleware(r))
	return nil
}

func (o *clientOption) middleware(r *requestconfig.RequestConfig) func(*http.Request, func(*http.Request) (*http.Response, error)) (*http.Response, error) {
	// replayed limits the 401 replay to one per request, however many times
	// the client's retry loop comes back here. The middleware is shared by
	// every page of a paginated call, so requests are told apart by their
	// Idempotency-Key: it stays the same across the retries of a request, and
	// requestconfig gives each page a new one.
	var (
		replayMu sync.Mutex
		replayed map[string]bool
	)
	firstReplay := func(key string) bool {
		replayMu.Lock()
		defer replayMu.Unlock()
		if replayed[key] {
			return false
		}
		if replayed == nil {
			replayed = map[string]bool{}
		}
		replayed[key] = true
		return true
	}
	return func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
		if active, _ := r.Context.Value(activeCredentialKey{}).(*clientOption); active != o {
			return next(req)
		}
		if req.Header.Get("Authorization") != "" {
			return next(req)
		}
		// The token is for the client's API only. A request can name another
		// host with an absolute URL.
		if r.BaseURL == nil || !sameOrigin(req.URL, r.BaseURL) {
			return next(req)
		}

		cache, err := o.cache(r.BaseURL.String())
		if err != nil {
			return failRequest(req, err)
		}
		// Token requests use the HTTP client of the request that needs the
		// token, so its proxy and transport settings apply, but not the
		// client's middleware: that is where request logging lives, and the
		// body of a token request is the identity token.
		var client HTTPDoer = r.HTTPClient
		if r.CustomHTTPDoer != nil {
			client = r.CustomHTTPDoer
		}
		ctx := withHTTPClient(req.Context(), client)
		// Read now: a later middleware may change the headers in place.
		idempotencyKey := req.Header.Get("Idempotency-Key")
		token, err := cache.Token(ctx)
		if err != nil {
			return failRequest(req, err)
		}
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		resp, err := next(req)
		// A later middleware can send the request elsewhere, as direct browser
		// routing does. A 401 from there says nothing about this token.
		if err != nil || resp.StatusCode != http.StatusUnauthorized || !sameOrigin(req.URL, r.BaseURL) {
			return resp, err
		}

		// The token was rejected, most likely revoked. Get a new one, which
		// takes a fresh identity token, and try once more.
		cache.invalidateToken(token.AccessToken)
		if req.Body != nil && req.GetBody == nil {
			return resp, nil
		}
		if !firstReplay(idempotencyKey) {
			return resp, nil
		}
		var body io.ReadCloser
		if req.GetBody != nil {
			if body, err = req.GetBody(); err != nil {
				return resp, nil
			}
		}
		// The rejected response is released before the exchange, which may
		// need the same connection.
		_ = resp.Body.Close()
		token, err = cache.Token(ctx)
		if err != nil {
			if body != nil {
				_ = body.Close()
			}
			return failRequest(req, err)
		}
		if body != nil {
			req.Body = body
		}
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		return next(req)
	}
}

func (o *clientOption) cache(apiBaseURL string) (*TokenCache, error) {
	if o.err != nil {
		return nil, o.err
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if cache, ok := o.caches[apiBaseURL]; ok {
		return cache, nil
	}
	cache := NewTokenCache(NewFederationProvider(FederationConfig{
		IdentityTokenProvider: o.identityTokenProvider,
		FederationRuleID:      o.opts.FederationRuleID,
		APIBaseURL:            apiBaseURL,
		Resource:              o.opts.Resource,
		AuthBaseURL:           o.opts.AuthBaseURL,
	}))
	if o.caches == nil {
		o.caches = map[string]*TokenCache{}
	}
	o.caches[apiBaseURL] = cache
	return cache, nil
}

// failRequest returns err for a request that could not be authenticated. The
// client retries a request that produced no response, which suits a token
// endpoint that was unreachable or overloaded: the retry makes a new exchange
// with a fresh identity token. Any other failure would repeat itself, so it
// comes with a response that tells the client not to retry.
func failRequest(req *http.Request, err error) (*http.Response, error) {
	var wiErr *WorkloadIdentityError
	if errors.As(err, &wiErr) && wiErr.retryable {
		return nil, err
	}
	return &http.Response{
		Status:  "credentials unavailable",
		Header:  http.Header{"X-Should-Retry": {"false"}},
		Body:    http.NoBody,
		Request: req,
	}, err
}

func sameOrigin(a, b *url.URL) bool {
	return strings.EqualFold(a.Scheme, b.Scheme) && strings.EqualFold(a.Host, b.Host)
}
