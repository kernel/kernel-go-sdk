package credentials

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kernel/kernel-go-sdk/internal"
)

const (
	// DefaultAPIBaseURL is the production Kernel API.
	DefaultAPIBaseURL = "https://api.onkernel.com"
	// DefaultAuthBaseURL is the production Kernel authorization server. It is
	// used only together with [DefaultAPIBaseURL].
	DefaultAuthBaseURL = "https://auth.onkernel.com"

	grantTypeJWTBearer = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	tokenEndpointPath  = "/token"

	// maxAssertionSize is the largest identity token the token endpoint
	// accepts. Real OIDC tokens are far smaller.
	maxAssertionSize = 16 << 10
	// maxTokenResponseSize bounds how much of a token endpoint response is
	// read.
	maxTokenResponseSize = 1 << 20

	// maxExpiresIn rejects an expires_in too large to be a number of seconds.
	maxExpiresIn = 1e9

	defaultExchangeTimeout = 30 * time.Second
)

// Token is a Kernel access token and what the authorization server said
// about it. A federation exchange fills in every field.
type Token struct {
	// AccessToken is sent to the Kernel API as "Authorization: Bearer ...".
	AccessToken string
	// ExpiresAt is when the access token stops being valid. A federation
	// exchange always sets it; the zero value, from another provider, means
	// the expiry is unknown. There is no refresh token: exchange a fresh
	// identity token to get a new access token.
	ExpiresAt time.Time
	// Scope lists the granted permissions, separated by spaces.
	Scope string
	// OrgID is the organization the token belongs to.
	OrgID string
	// AccessScope is "project" for a token bound to one project.
	AccessScope string
	// ProjectID is the project of the federation rule. A request that names
	// another project is refused.
	ProjectID string
	// Principal is the service account the token acts as.
	Principal Principal
	// FederationRuleID is the rule that matched the identity token.
	FederationRuleID string
}

// Principal identifies who a [Token] acts as.
type Principal struct {
	// Type is "service_account".
	Type string
	// ID is the service account ID.
	ID string
}

// FederationExchangeParams are the inputs of [ExchangeFederationAssertion].
type FederationExchangeParams struct {
	// Assertion is the identity token (a JWT) to exchange. Required. Most
	// providers issue single-use tokens, so an assertion that was already
	// sent, even in a request that failed, must not be sent again.
	Assertion string
	// FederationRuleID is the federation rule ("wifr_...") to check the
	// assertion against. Required.
	FederationRuleID string
	// Resource is the API the token is for, such as "https://api.onkernel.com"
	// with no trailing slash. When empty the server chooses.
	Resource string
	// AuthBaseURL is the Kernel authorization server. It defaults to
	// KERNEL_AUTH_BASE_URL, then to [DefaultAuthBaseURL].
	AuthBaseURL string
	// HTTPClient sends the request. It defaults to a client with a 30 second
	// timeout. An [*http.Client] is used with redirects turned off; any other
	// implementation must not redirect or retry the request itself.
	HTTPClient HTTPDoer
	// UserAgent replaces the SDK's User-Agent header.
	UserAgent string
}

// ExchangeFederationAssertion trades an identity token for a Kernel access
// token in a single request to the authorization server. It needs no Kernel
// credential.
//
// A refused exchange is reported as a [*WorkloadIdentityError].
func ExchangeFederationAssertion(ctx context.Context, params FederationExchangeParams) (*Token, error) {
	if params.Assertion == "" {
		return nil, configError("ExchangeFederationAssertion: Assertion is required")
	}
	if params.FederationRuleID == "" {
		return nil, configError("ExchangeFederationAssertion: FederationRuleID is required")
	}
	authBaseURL := params.AuthBaseURL
	if authBaseURL == "" {
		authBaseURL = os.Getenv(EnvAuthBaseURL)
	}
	if authBaseURL == "" {
		authBaseURL = DefaultAuthBaseURL
	}
	_, endpoint, err := tokenEndpoint(authBaseURL)
	if err != nil {
		return nil, err
	}
	client := params.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: defaultExchangeTimeout}
	}
	return exchange(ctx, exchangeRequest{
		client:           withoutRedirects(client),
		endpoint:         endpoint,
		assertion:        params.Assertion,
		federationRuleID: params.FederationRuleID,
		resource:         params.Resource,
		userAgent:        params.UserAgent,
	})
}

// FederationConfig configures [NewFederationProvider].
type FederationConfig struct {
	// IdentityTokenProvider returns a fresh identity token for every
	// exchange. Required.
	IdentityTokenProvider IdentityTokenFunc
	// FederationRuleID is the federation rule ("wifr_...") to check the
	// identity token against. Required.
	FederationRuleID string
	// APIBaseURL is the Kernel API the access token will be used with. It
	// defaults to [DefaultAPIBaseURL].
	APIBaseURL string
	// Resource is the API the token is requested for. It defaults to the
	// origin of APIBaseURL, such as "https://api.onkernel.com".
	Resource string
	// AuthBaseURL is the Kernel authorization server. It defaults to
	// KERNEL_AUTH_BASE_URL, then to [DefaultAuthBaseURL] when APIBaseURL is
	// exactly the production API. For any other API it must be set: a custom
	// API is never paired with the production authorization server
	// implicitly.
	AuthBaseURL string
	// HTTPClient sends the token requests. It defaults to a client with a 30
	// second timeout. An [*http.Client] is used with redirects turned off;
	// any other implementation must not redirect or retry a request itself.
	HTTPClient HTTPDoer
	// UserAgent replaces the SDK's User-Agent header on token requests.
	UserAgent string
}

// NewFederationProvider returns an [AccessTokenProvider] that exchanges an
// identity token for a Kernel access token. Every call asks
// cfg.IdentityTokenProvider for a fresh token and makes a new exchange; wrap
// the provider in a [TokenCache] to reuse the result.
//
// The identity-token provider can read the audience its token should carry
// with [ResolveAudience].
func NewFederationProvider(cfg FederationConfig) AccessTokenProvider {
	endpoints, cfgErr := resolveEndpoints(cfg)
	if cfgErr == nil && cfg.IdentityTokenProvider == nil {
		cfgErr = configError("an identity token provider is required")
	}
	if cfgErr == nil && cfg.FederationRuleID == "" {
		cfgErr = configError("a federation rule ID is required")
	}
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: defaultExchangeTimeout}
	}
	client = withoutRedirects(client)

	// A consumed single-use source cannot say why the exchange that used it
	// failed, so that failure is kept to report alongside. The digest of the
	// last assertion tells a source that repeats itself from one whose token
	// was spent elsewhere.
	var (
		mu            sync.Mutex
		lastFailure   error
		lastAssertion [sha256.Size]byte
		sentBefore    bool
	)

	return func(ctx context.Context, _ bool) (*Token, error) {
		if cfgErr != nil {
			return nil, cfgErr
		}
		var source sourceInfo
		assertion, err := cfg.IdentityTokenProvider(withSourceInfo(withAudience(ctx, endpoints.audience), &source))
		if err != nil {
			var transient transientError
			wrapped := &WorkloadIdentityError{
				Message:   "could not get an identity token",
				Err:       err,
				retryable: errors.As(err, &transient),
			}
			if errors.Is(err, ErrIdentityTokenConsumed) {
				mu.Lock()
				if lastFailure != nil {
					wrapped.Hint = "The exchange that used it failed: " + lastFailure.Error()
				}
				mu.Unlock()
			}
			return nil, wrapped
		}
		assertion = strings.TrimSpace(assertion)
		if assertion == "" {
			return nil, configError("the identity token provider returned an empty token")
		}

		digest := sha256.Sum256([]byte(assertion))
		mu.Lock()
		repeated := sentBefore && digest == lastAssertion
		lastAssertion, sentBefore = digest, true
		mu.Unlock()

		requestClient := client
		if override := httpClientFrom(ctx); override != nil {
			requestClient = withoutRedirects(override)
		}
		token, err := exchange(ctx, exchangeRequest{
			client:           requestClient,
			endpoint:         endpoints.token,
			assertion:        assertion,
			federationRuleID: cfg.FederationRuleID,
			resource:         endpoints.resource,
			userAgent:        cfg.UserAgent,
		})
		var wiErr *WorkloadIdentityError
		if errors.As(err, &wiErr) && wiErr.replay {
			wiErr.Hint = replayHint(source, repeated)
		}
		mu.Lock()
		lastFailure = err
		mu.Unlock()
		return token, err
	}
}

type httpClientKey struct{}

// withHTTPClient makes a federation provider called with ctx send its token
// request with client instead of the one it was configured with. A client
// option uses it to exchange through the HTTP client of whichever request
// needs the token.
func withHTTPClient(ctx context.Context, client HTTPDoer) context.Context {
	return context.WithValue(ctx, httpClientKey{}, client)
}

func httpClientFrom(ctx context.Context) HTTPDoer {
	client, _ := ctx.Value(httpClientKey{}).(HTTPDoer)
	return client
}

type federationEndpoints struct {
	token    string
	resource string
	// audience is what the built-in identity-token sources request by
	// default: the authorization server's URL.
	audience string
}

func resolveEndpoints(cfg FederationConfig) (federationEndpoints, error) {
	apiBaseURL := cfg.APIBaseURL
	if apiBaseURL == "" {
		apiBaseURL = DefaultAPIBaseURL
	}
	apiURL, err := url.Parse(apiBaseURL)
	if err != nil || apiURL.Scheme == "" || apiURL.Host == "" {
		return federationEndpoints{}, configError("invalid API base URL %q", apiBaseURL)
	}
	apiOrigin := apiURL.Scheme + "://" + strings.ToLower(apiURL.Host)
	// Only the production API itself implies the production authorization
	// server; a path on the same host could be a gateway to anything.
	production := apiOrigin == DefaultAPIBaseURL && (apiURL.EscapedPath() == "" || apiURL.EscapedPath() == "/") &&
		apiURL.RawQuery == "" && apiURL.User == nil

	authBaseURL := cfg.AuthBaseURL
	if authBaseURL == "" {
		authBaseURL = os.Getenv(EnvAuthBaseURL)
	}
	if authBaseURL == "" {
		if !production {
			// Only the origin is printed: the rest of a gateway URL can hold
			// credentials.
			return federationEndpoints{}, configError(
				"the API base URL (on %s) is not %s, so its authorization server is unknown: set the auth base URL (FederationOptions.AuthBaseURL or %s)",
				apiOrigin, DefaultAPIBaseURL, EnvAuthBaseURL)
		}
		authBaseURL = DefaultAuthBaseURL
	}
	authBaseURL, endpoint, err := tokenEndpoint(authBaseURL)
	if err != nil {
		return federationEndpoints{}, err
	}

	resource := cfg.Resource
	if resource == "" {
		resource = apiOrigin
	}
	return federationEndpoints{token: endpoint, resource: resource, audience: authBaseURL}, nil
}

// tokenEndpoint returns the authorization server's URL without a trailing
// slash, and its token endpoint. It refuses a URL that would send the
// identity token in cleartext; plain http is allowed only for a loopback
// host. The messages leave the URL out where it may hold credentials.
func tokenEndpoint(authBaseURL string) (base, endpoint string, err error) {
	u, parseErr := url.Parse(authBaseURL)
	if parseErr != nil || u.Host == "" {
		return "", "", configError("invalid auth base URL: want a URL such as %s", DefaultAuthBaseURL)
	}
	// net/http would turn credentials into an Authorization header, which
	// the token endpoint rejects, and a query or fragment would end up
	// after "/token".
	if u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", "", configError("invalid auth base URL: it must not contain credentials, a query or a fragment")
	}
	base = u.Scheme + "://" + u.Host + strings.TrimRight(u.EscapedPath(), "/")
	switch host := strings.ToLower(u.Hostname()); {
	case u.Scheme == "https":
	case u.Scheme == "http" && (host == "localhost" || host == "127.0.0.1" || host == "::1"):
	default:
		return "", "", configError("refusing to send an identity token to the non-https auth base URL %q", base)
	}
	return base, base + tokenEndpointPath, nil
}

// withoutRedirects returns client with redirects turned off when it is an
// [*http.Client]: a 307 or 308 would replay the request body, which carries
// the identity token, to a host that was never checked. client is not
// modified. Any other [HTTPDoer] keeps its own policy; the token request is
// built so that net/http cannot replay its body either way.
func withoutRedirects(client HTTPDoer) HTTPDoer {
	c, ok := client.(*http.Client)
	if !ok {
		return client
	}
	clone := *c
	clone.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &clone
}

type exchangeRequest struct {
	client           HTTPDoer
	endpoint         string
	assertion        string
	federationRuleID string
	resource         string
	userAgent        string
}

// exchange posts the jwt-bearer grant. The server rejects any parameter
// other than these four, and any Authorization header.
func exchange(ctx context.Context, r exchangeRequest) (*Token, error) {
	if len(r.assertion) > maxAssertionSize {
		return nil, configError("the identity token is %d bytes, over the %d byte limit", len(r.assertion), maxAssertionSize)
	}
	form := url.Values{
		"grant_type":         {grantTypeJWTBearer},
		"assertion":          {r.assertion},
		"federation_rule_id": {r.federationRuleID},
	}
	if r.resource != "" {
		form.Set("resource", r.resource)
	}
	// The body is deliberately not rewindable (no GetBody): an identity token
	// is single-use, so nothing below this function may send it twice, be it
	// a redirect or a transport retry.
	body := form.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.endpoint, io.NopCloser(strings.NewReader(body)))
	if err != nil {
		return nil, &WorkloadIdentityError{Message: "could not build the token request", Err: err}
	}
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	userAgent := r.userAgent
	if userAgent == "" {
		userAgent = fmt.Sprintf("Kernel/Go %s", internal.PackageVersion)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, &WorkloadIdentityError{
			Message:   "token request failed",
			Err:       redactError(err, r.assertion),
			retryable: true,
		}
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxTokenResponseSize))
	if err != nil {
		return nil, &WorkloadIdentityError{
			StatusCode: resp.StatusCode,
			RequestID:  requestID(resp, r.assertion),
			Message:    "could not read the token endpoint's response",
			Err:        redactError(err, r.assertion),
			retryable:  true,
		}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, exchangeError(resp, respBody, r.assertion)
	}

	// A success body holds the access token, so nothing from it goes into an
	// error.
	invalid := func(message string) error {
		return &WorkloadIdentityError{StatusCode: resp.StatusCode, RequestID: requestID(resp, r.assertion), Message: message}
	}
	var parsed struct {
		AccessToken string   `json:"access_token"`
		TokenType   string   `json:"token_type"`
		ExpiresIn   *float64 `json:"expires_in"`
		Scope       string   `json:"scope"`
		OrgID       string   `json:"org_id"`
		AccessScope string   `json:"access_scope"`
		ProjectID   string   `json:"project_id"`
		Principal   struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"principal"`
		FederationRuleID string `json:"federation_rule_id"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, invalid("the token endpoint's response is not valid JSON")
	}
	if parsed.AccessToken == "" {
		return nil, invalid("the token endpoint's response has no access_token")
	}
	if parsed.TokenType != "" && !strings.EqualFold(parsed.TokenType, "Bearer") {
		return nil, invalid("the token endpoint returned an unsupported token_type (want Bearer)")
	}
	// The token lives for minutes, so a response that does not say for how
	// long cannot be cached safely.
	if parsed.ExpiresIn == nil || *parsed.ExpiresIn <= 0 || *parsed.ExpiresIn > maxExpiresIn {
		return nil, invalid("the token endpoint returned no valid expires_in")
	}
	return &Token{
		AccessToken:      parsed.AccessToken,
		ExpiresAt:        time.Now().Add(time.Duration(*parsed.ExpiresIn * float64(time.Second))),
		Scope:            parsed.Scope,
		OrgID:            parsed.OrgID,
		AccessScope:      parsed.AccessScope,
		ProjectID:        parsed.ProjectID,
		Principal:        Principal{Type: parsed.Principal.Type, ID: parsed.Principal.ID},
		FederationRuleID: parsed.FederationRuleID,
	}, nil
}
