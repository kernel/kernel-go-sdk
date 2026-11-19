package credentials_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kernel/kernel-go-sdk/lib/credentials"
)

func TestExchangeFederationAssertion(t *testing.T) {
	cleanEnv(t)
	auth := newAuthServer(t)
	assertion := providerToken(t, providerGitHub, auth.URL)

	before := time.Now()
	token, err := credentials.ExchangeFederationAssertion(context.Background(), credentials.FederationExchangeParams{
		Assertion:        assertion,
		FederationRuleID: testRuleID,
		Resource:         "https://api.dev.onkernel.com",
		AuthBaseURL:      auth.URL + "/",
	})
	if err != nil {
		t.Fatalf("ExchangeFederationAssertion: %v", err)
	}

	want := credentials.Token{
		AccessToken:      "kernel-access-token-1",
		Scope:            "browsers:create browsers:use browsers:delete",
		OrgID:            "org_1",
		AccessScope:      "project",
		ProjectID:        "proj_1",
		Principal:        credentials.Principal{Type: "service_account", ID: "svc_1"},
		FederationRuleID: testRuleID,
	}
	got := *token
	got.ExpiresAt = time.Time{}
	if got != want {
		t.Errorf("token = %+v, want %+v", got, want)
	}
	if lo, hi := before.Add(298*time.Second), time.Now().Add(298*time.Second); token.ExpiresAt.Before(lo) || token.ExpiresAt.After(hi) {
		t.Errorf("ExpiresAt = %v, want now + expires_in (298s)", token.ExpiresAt)
	}

	requests := auth.exchanges()
	if len(requests) != 1 {
		t.Fatalf("token endpoint got %d requests, want 1", len(requests))
	}
	req := requests[0]
	wantBody := "assertion=" + assertion +
		"&federation_rule_id=" + testRuleID +
		"&grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Ajwt-bearer" +
		"&resource=https%3A%2F%2Fapi.dev.onkernel.com"
	if req.Body != wantBody {
		t.Errorf("request body = %q\nwant %q", req.Body, wantBody)
	}
	if got := req.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
		t.Errorf("Content-Type = %q", got)
	}
	if got := req.Header.Get("Authorization"); got != "" {
		t.Errorf("Authorization = %q, want none", got)
	}
	if got := req.Header.Get("User-Agent"); !strings.HasPrefix(got, "Kernel/Go ") {
		t.Errorf("User-Agent = %q, want the SDK's", got)
	}
}

func TestExchangeFederationAssertionOmitsEmptyResource(t *testing.T) {
	cleanEnv(t)
	auth := newAuthServer(t)
	_, err := credentials.ExchangeFederationAssertion(context.Background(), credentials.FederationExchangeParams{
		Assertion:        providerToken(t, providerAWS, auth.URL),
		FederationRuleID: testRuleID,
		AuthBaseURL:      auth.URL,
		UserAgent:        "kernel-cli/1.2.3",
	})
	if err != nil {
		t.Fatalf("ExchangeFederationAssertion: %v", err)
	}
	req := auth.exchanges()[0]
	if _, ok := req.Form["resource"]; ok {
		t.Errorf("resource was sent: %q", req.Body)
	}
	if got := req.Header.Get("User-Agent"); got != "kernel-cli/1.2.3" {
		t.Errorf("User-Agent = %q, want the override", got)
	}
}

func TestExchangeFederationAssertionAuthBaseURLFromEnv(t *testing.T) {
	cleanEnv(t)
	auth := newAuthServer(t)
	t.Setenv(credentials.EnvAuthBaseURL, auth.URL)
	_, err := credentials.ExchangeFederationAssertion(context.Background(), credentials.FederationExchangeParams{
		Assertion:        providerToken(t, providerKubernetes, auth.URL),
		FederationRuleID: testRuleID,
	})
	if err != nil {
		t.Fatalf("ExchangeFederationAssertion: %v", err)
	}
	if n := len(auth.exchanges()); n != 1 {
		t.Fatalf("token endpoint got %d requests, want 1", n)
	}
}

func TestExchangeFederationAssertionValidatesInput(t *testing.T) {
	cleanEnv(t)
	tests := []struct {
		name   string
		params credentials.FederationExchangeParams
		want   string
	}{
		{"no assertion", credentials.FederationExchangeParams{FederationRuleID: testRuleID}, "Assertion is required"},
		{"no rule", credentials.FederationExchangeParams{Assertion: "a.b.c"}, "FederationRuleID is required"},
		{"cleartext", credentials.FederationExchangeParams{Assertion: "a.b.c", FederationRuleID: testRuleID, AuthBaseURL: "http://auth.example.com"}, "non-https"},
		{"not a URL", credentials.FederationExchangeParams{Assertion: "a.b.c", FederationRuleID: testRuleID, AuthBaseURL: "auth.example.com"}, "invalid auth base URL"},
		{"credentials in the URL", credentials.FederationExchangeParams{Assertion: "a.b.c", FederationRuleID: testRuleID, AuthBaseURL: "https://user:hunter2@auth.example.com"}, "must not contain credentials"},
		{"query in the URL", credentials.FederationExchangeParams{Assertion: "a.b.c", FederationRuleID: testRuleID, AuthBaseURL: "https://auth.example.com?tenant=1"}, "must not contain credentials, a query"},
		{"fragment in the URL", credentials.FederationExchangeParams{Assertion: "a.b.c", FederationRuleID: testRuleID, AuthBaseURL: "https://auth.example.com#x"}, "must not contain credentials, a query or a fragment"},
		{"oversized assertion", credentials.FederationExchangeParams{Assertion: strings.Repeat("a", 16*1024+1), FederationRuleID: testRuleID, AuthBaseURL: "https://auth.example.com"}, "over the 16384 byte limit"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.params.HTTPClient = doerFunc(func(*http.Request) (*http.Response, error) {
				t.Error("a request was sent")
				return nil, errors.New("stop")
			})
			_, err := credentials.ExchangeFederationAssertion(context.Background(), tt.params)
			var wiErr *credentials.WorkloadIdentityError
			if !errors.As(err, &wiErr) || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want a WorkloadIdentityError containing %q", err, tt.want)
			}
			if strings.Contains(err.Error(), "hunter2") {
				t.Fatalf("error %q contains the URL's password", err)
			}
		})
	}
}

func TestTokenEndpointMustBeHTTPSExceptLoopback(t *testing.T) {
	cleanEnv(t)
	identity := func(context.Context) (string, error) { return "a.b.c", nil }
	tests := []struct {
		authBaseURL string
		allowed     bool
	}{
		{"https://auth.onkernel.com", true},
		{"http://localhost:8080", true},
		{"http://127.0.0.1:8080", true},
		{"http://[::1]:8080", true},
		{"http://auth.onkernel.com", false},
		{"http://127.0.0.1.example.com", false},
		{"http://10.0.0.5", false},
		{"ftp://auth.onkernel.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.authBaseURL, func(t *testing.T) {
			var sent bool
			provider := credentials.NewFederationProvider(credentials.FederationConfig{
				IdentityTokenProvider: identity,
				FederationRuleID:      testRuleID,
				AuthBaseURL:           tt.authBaseURL,
				HTTPClient: doerFunc(func(*http.Request) (*http.Response, error) {
					sent = true
					return nil, errors.New("stop")
				}),
			})
			_, err := provider(context.Background(), false)
			if sent != tt.allowed {
				t.Fatalf("request sent = %v, want %v (error: %v)", sent, tt.allowed, err)
			}
			if !tt.allowed && !strings.Contains(err.Error(), "non-https") {
				t.Fatalf("error = %v, want a refusal to use a non-https endpoint", err)
			}
		})
	}
}

type doerFunc func(*http.Request) (*http.Response, error)

func (f doerFunc) Do(req *http.Request) (*http.Response, error) { return f(req) }

func TestFederationProviderUsesFreshIdentityTokenForEveryExchange(t *testing.T) {
	cleanEnv(t)
	auth := newAuthServer(t)
	source := &freshTokens{t: t, provider: providerGitHub}
	provider := credentials.NewFederationProvider(credentials.FederationConfig{
		IdentityTokenProvider: source.get,
		FederationRuleID:      testRuleID,
		AuthBaseURL:           auth.URL,
	})

	seen := map[string]bool{}
	for i := 1; i <= 3; i++ {
		token, err := provider(context.Background(), false)
		if err != nil {
			t.Fatalf("exchange %d: %v", i, err)
		}
		if token.ExpiresAt.IsZero() {
			t.Fatalf("exchange %d: no expiry", i)
		}
		seen[auth.exchanges()[i-1].Form.Get("assertion")] = true
	}
	if source.count() != 3 || len(seen) != 3 {
		t.Fatalf("identity source called %d times, %d distinct assertions sent; want 3 and 3", source.count(), len(seen))
	}
	// The default resource is the production API, and the audience handed to
	// the source is the authorization server.
	req := auth.exchanges()[0]
	if got := req.Form.Get("resource"); got != "https://api.onkernel.com" {
		t.Errorf("resource = %q, want https://api.onkernel.com", got)
	}
	if got := tokenAudiences(req.Form.Get("assertion")); len(got) != 1 || got[0] != auth.URL {
		t.Errorf("identity token audience = %v, want [%s]", got, auth.URL)
	}
}

func TestTokenCacheKeepsExchangeDetails(t *testing.T) {
	cleanEnv(t)
	auth := newAuthServer(t)
	source := &freshTokens{t: t, provider: providerKubernetes}
	cache := credentials.NewTokenCache(credentials.NewFederationProvider(credentials.FederationConfig{
		IdentityTokenProvider: source.get,
		FederationRuleID:      testRuleID,
		AuthBaseURL:           auth.URL,
	}))

	token, err := cache.Token(context.Background())
	if err != nil {
		t.Fatalf("Token: %v", err)
	}
	if token.Principal.ID != "svc_1" || token.ProjectID != "proj_1" || token.FederationRuleID != testRuleID {
		t.Fatalf("token = %+v, want the principal, project and rule from the exchange", token)
	}
	// The caller gets its own copy.
	token.AccessToken = "changed"
	again, err := cache.Token(context.Background())
	if err != nil || again.AccessToken != "kernel-access-token-1" {
		t.Fatalf("second Token = %+v, %v", again, err)
	}
	if n := len(auth.exchanges()); n != 1 {
		t.Fatalf("%d exchanges, want 1", n)
	}
}

// TestReplayHintNamesTheSource checks that an "already used" refusal is
// explained in terms of where the identity token came from.
func TestReplayHintNamesTheSource(t *testing.T) {
	cleanEnv(t)
	tokenFile := filepath.Join(t.TempDir(), "token")
	spentFile := filepath.Join(t.TempDir(), "spent")
	spent := uniqueToken("spent.static")
	same := uniqueToken("callback")
	for path, token := range map[string]string{tokenFile: uniqueToken("file"), spentFile: uniqueToken("spent.file")} {
		if err := os.WriteFile(path, []byte(token+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name   string
		source credentials.IdentityTokenFunc
		// exchanges is how many exchanges the provider makes; the last one
		// is the one refused. With 1, the token was spent by someone else.
		exchanges int
		want      string
		wantNot   string
	}{
		{
			name:      "file that was not rotated",
			source:    credentials.File(tokenFile),
			exchanges: 2,
			want: "The identity token file " + tokenFile + " held the same token as in the previous exchange, " +
				"so whatever writes the file must rotate it before each refresh. " +
				"For a Kubernetes projected token, the issuer can allow reuse instead.",
			wantNot: "KERNEL_IDENTITY_TOKEN ",
		},
		{
			name:      "file whose token was spent elsewhere",
			source:    credentials.File(spentFile),
			exchanges: 1,
			want:      "The token in the identity token file " + spentFile + " was already exchanged, so whatever writes the file must rotate it",
			wantNot:   "KERNEL_IDENTITY_TOKEN ",
		},
		{
			name:      "static value",
			source:    credentials.Static(spent),
			exchanges: 1,
			want:      "a static KERNEL_IDENTITY_TOKEN can be exchanged only once.",
		},
		{
			name:      "callback that repeats itself",
			source:    func(context.Context) (string, error) { return same, nil },
			exchanges: 2,
			want:      "The identity token provider returned the same token as for the previous exchange; it must return a fresh one every time it is called.",
			wantNot:   "KERNEL_IDENTITY_TOKEN",
		},
		{
			name:      "callback whose token was spent elsewhere",
			source:    func(context.Context) (string, error) { return spent, nil },
			exchanges: 1,
			want:      "Every exchange needs a fresh identity token; an identity token can be exchanged only once.",
			wantNot:   "KERNEL_IDENTITY_TOKEN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := newAuthServer(t)
			// Tokens "spent elsewhere" have been exchanged by another process.
			for _, token := range []string{spent, readToken(t, spentFile)} {
				if _, err := credentials.ExchangeFederationAssertion(context.Background(), credentials.FederationExchangeParams{
					Assertion: token, FederationRuleID: testRuleID, AuthBaseURL: auth.URL,
				}); err != nil {
					t.Fatal(err)
				}
			}
			provider := credentials.NewFederationProvider(credentials.FederationConfig{
				IdentityTokenProvider: tt.source,
				FederationRuleID:      testRuleID,
				AuthBaseURL:           auth.URL,
			})
			var err error
			for i := 1; i <= tt.exchanges; i++ {
				_, err = provider(context.Background(), false)
				if i < tt.exchanges && err != nil {
					t.Fatalf("exchange %d: %v", i, err)
				}
			}
			var wiErr *credentials.WorkloadIdentityError
			if !errors.As(err, &wiErr) || wiErr.Code != "invalid_grant" {
				t.Fatalf("error = %v, want the replay refusal", err)
			}
			if !strings.Contains(wiErr.Hint, tt.want) {
				t.Errorf("hint = %q\nwant it to contain %q", wiErr.Hint, tt.want)
			}
			if tt.wantNot != "" && strings.Contains(wiErr.Hint, tt.wantNot) {
				t.Errorf("hint %q mentions %q", wiErr.Hint, tt.wantNot)
			}
		})
	}
}

func readToken(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(data))
}

func TestFederationProviderEndpointDefaults(t *testing.T) {
	identity := func(ctx context.Context) (string, error) { return credentials.ResolveAudience(ctx, "") }
	// capture reports the token endpoint, resource and default audience a
	// provider resolves to, without a network.
	capture := func(cfg credentials.FederationConfig) (endpoint, resource, audience string, err error) {
		cfg.IdentityTokenProvider = identity
		cfg.FederationRuleID = testRuleID
		cfg.HTTPClient = doerFunc(func(req *http.Request) (*http.Response, error) {
			endpoint = req.URL.String()
			_ = req.ParseForm()
			resource, audience = req.PostForm.Get("resource"), req.PostForm.Get("assertion")
			return nil, errors.New("stop")
		})
		_, err = credentials.NewFederationProvider(cfg)(context.Background(), false)
		return
	}

	t.Run("production", func(t *testing.T) {
		cleanEnv(t)
		endpoint, resource, audience, _ := capture(credentials.FederationConfig{})
		if endpoint != "https://auth.onkernel.com/token" || resource != "https://api.onkernel.com" || audience != "https://auth.onkernel.com" {
			t.Fatalf("got endpoint %q, resource %q, audience %q", endpoint, resource, audience)
		}
	})
	t.Run("production with a trailing slash and path", func(t *testing.T) {
		cleanEnv(t)
		endpoint, resource, _, _ := capture(credentials.FederationConfig{APIBaseURL: "https://api.onkernel.com/"})
		if endpoint != "https://auth.onkernel.com/token" || resource != "https://api.onkernel.com" {
			t.Fatalf("got endpoint %q, resource %q", endpoint, resource)
		}
	})
	for _, apiBaseURL := range []string{
		"https://api.dev.onkernel.com/",
		"https://api.onkernel.com/gateway/",
		"https://api.onkernel.com/%2F/",
		"https://api.onkernel.com//",
		"https://user:secret-in-url@api.onkernel.com/?token=secret-in-url",
		"https://api.onkernel.com:8443/",
		"http://api.onkernel.com/",
		"https://api.onkernel.com.example.com/",
	} {
		t.Run("no production auth for "+apiBaseURL, func(t *testing.T) {
			cleanEnv(t)
			endpoint, _, _, err := capture(credentials.FederationConfig{APIBaseURL: apiBaseURL})
			if endpoint != "" {
				t.Fatalf("a request was sent to %q", endpoint)
			}
			if err == nil || !strings.Contains(err.Error(), "KERNEL_AUTH_BASE_URL") {
				t.Fatalf("error = %v, want one naming KERNEL_AUTH_BASE_URL", err)
			}
			if strings.Contains(err.Error(), "secret-in-url") {
				t.Fatalf("error %q contains credentials from the URL", err)
			}
		})
	}
	t.Run("auth base URL with a path", func(t *testing.T) {
		cleanEnv(t)
		endpoint, _, audience, _ := capture(credentials.FederationConfig{AuthBaseURL: "https://sso.example.com/kernel/"})
		if endpoint != "https://sso.example.com/kernel/token" || audience != "https://sso.example.com/kernel" {
			t.Fatalf("got endpoint %q, audience %q", endpoint, audience)
		}
	})
	t.Run("auth base URL with an escaped path", func(t *testing.T) {
		cleanEnv(t)
		endpoint, _, audience, _ := capture(credentials.FederationConfig{AuthBaseURL: "https://sso.example.com/tenant%3Fblue/"})
		if endpoint != "https://sso.example.com/tenant%3Fblue/token" || audience != "https://sso.example.com/tenant%3Fblue" {
			t.Fatalf("got endpoint %q, audience %q", endpoint, audience)
		}
	})
	t.Run("custom API with KERNEL_AUTH_BASE_URL", func(t *testing.T) {
		cleanEnv(t)
		t.Setenv(credentials.EnvAuthBaseURL, "https://auth.dev.onkernel.com/")
		endpoint, resource, audience, _ := capture(credentials.FederationConfig{APIBaseURL: "https://api.dev.onkernel.com/"})
		if endpoint != "https://auth.dev.onkernel.com/token" || resource != "https://api.dev.onkernel.com" || audience != "https://auth.dev.onkernel.com" {
			t.Fatalf("got endpoint %q, resource %q, audience %q", endpoint, resource, audience)
		}
	})
	t.Run("explicit values win", func(t *testing.T) {
		cleanEnv(t)
		t.Setenv(credentials.EnvAuthBaseURL, "https://auth.dev.onkernel.com")
		t.Setenv(credentials.EnvIdentityTokenAudience, "custom-audience")
		endpoint, resource, audience, _ := capture(credentials.FederationConfig{
			APIBaseURL:  "https://api.dev.onkernel.com",
			AuthBaseURL: "https://auth.example.com",
			Resource:    "https://mcp.example.com",
		})
		if endpoint != "https://auth.example.com/token" || resource != "https://mcp.example.com" || audience != "custom-audience" {
			t.Fatalf("got endpoint %q, resource %q, audience %q", endpoint, resource, audience)
		}
	})
}

func TestFederationProviderRejectsBadIdentityTokens(t *testing.T) {
	cleanEnv(t)
	auth := newAuthServer(t)
	huge := strings.Repeat("a", 16*1024+1)
	tests := []struct {
		name     string
		identity credentials.IdentityTokenFunc
		want     string
	}{
		{"oversized", func(context.Context) (string, error) { return huge, nil }, "over the 16384 byte limit"},
		{"empty", func(context.Context) (string, error) { return " \n", nil }, "empty token"},
		{"source error", func(context.Context) (string, error) { return "", errors.New("no metadata server") }, "no metadata server"},
		{"missing", nil, "identity token provider is required"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := credentials.NewFederationProvider(credentials.FederationConfig{
				IdentityTokenProvider: tt.identity,
				FederationRuleID:      testRuleID,
				AuthBaseURL:           auth.URL,
			})
			_, err := provider(context.Background(), false)
			var wiErr *credentials.WorkloadIdentityError
			if !errors.As(err, &wiErr) || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want a WorkloadIdentityError containing %q", err, tt.want)
			}
			if strings.Contains(err.Error(), huge[:64]) {
				t.Fatal("the error contains the identity token")
			}
		})
	}
	if n := len(auth.exchanges()); n != 0 {
		t.Fatalf("token endpoint got %d requests, want 0", n)
	}
}

func TestFederationProviderAcceptsLargestAssertion(t *testing.T) {
	cleanEnv(t)
	auth := newAuthServer(t)
	provider := credentials.NewFederationProvider(credentials.FederationConfig{
		IdentityTokenProvider: func(context.Context) (string, error) { return strings.Repeat("a", 16*1024), nil },
		FederationRuleID:      testRuleID,
		AuthBaseURL:           auth.URL,
	})
	if _, err := provider(context.Background(), false); err != nil {
		t.Fatalf("a 16 KiB assertion was refused: %v", err)
	}
}

func TestExchangeErrors(t *testing.T) {
	cleanEnv(t)
	const assertion = "header.payload-that-must-not-leak.signature"
	tests := []struct {
		name        string
		status      int
		body        string
		header      http.Header
		wantCode    string
		wantDesc    string
		wantMessage []string
	}{
		{
			name:        "invalid grant",
			status:      400,
			body:        `{"error":"invalid_grant","error_description":"Credential is invalid, expired, or revoked"}`,
			wantCode:    "invalid_grant",
			wantDesc:    "Credential is invalid, expired, or revoked",
			wantMessage: []string{"status 400", "invalid_grant: Credential is invalid, expired, or revoked"},
		},
		{
			name:        "replayed assertion gets a hint",
			status:      400,
			body:        `{"error":"invalid_grant","error_description":"Credential was already used"}`,
			header:      http.Header{"X-Request-Id": {"req_abc"}},
			wantCode:    "invalid_grant",
			wantDesc:    "Credential was already used",
			wantMessage: []string{"request id req_abc", "Every exchange needs a fresh identity token; an identity token can be exchanged only once."},
		},
		{
			name:        "invalid target",
			status:      400,
			body:        `{"error":"invalid_target","error_description":"Unknown resource","error_uri":"https://kernel.sh/docs/wif"}`,
			wantCode:    "invalid_target",
			wantDesc:    "Unknown resource",
			wantMessage: []string{"invalid_target: Unknown resource", "https://kernel.sh/docs/wif", "Set Resource to the Kernel API origin with no trailing slash."},
		},
		{
			name:        "temporarily unavailable",
			status:      503,
			body:        `{"error":"temporarily_unavailable","error_description":"Try again"}`,
			wantCode:    "temporarily_unavailable",
			wantDesc:    "Try again",
			wantMessage: []string{"status 503"},
		},
		{
			name:   "other fields are dropped",
			status: 400,
			body: `{"error":"invalid_grant","error_description":"bad token ` + assertion + `",` +
				`"assertion":"` + assertion + `","access_token":"secret-access-token","debug":{"jwt":"` + assertion + `"}}`,
			wantCode:    "invalid_grant",
			wantDesc:    "bad token [redacted]",
			wantMessage: []string{"invalid_grant: bad token [redacted]"},
		},
		{
			name:        "not JSON",
			status:      502,
			body:        "<html>" + assertion + " secret-access-token</html>",
			wantMessage: []string{"status 502", "unexpected response"},
		},
		{
			name:        "redirect is not followed",
			status:      307,
			header:      http.Header{"Location": {"https://elsewhere.example.com/token"}},
			wantMessage: []string{"status 307"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requests int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				for key, values := range tt.header {
					w.Header()[key] = values
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			_, err := credentials.ExchangeFederationAssertion(context.Background(), credentials.FederationExchangeParams{
				Assertion:        assertion,
				FederationRuleID: testRuleID,
				AuthBaseURL:      server.URL,
			})
			var wiErr *credentials.WorkloadIdentityError
			if !errors.As(err, &wiErr) {
				t.Fatalf("error = %v (%T), want a *WorkloadIdentityError", err, err)
			}
			if wiErr.StatusCode != tt.status || wiErr.Code != tt.wantCode || wiErr.Description != tt.wantDesc {
				t.Errorf("got status %d, code %q, description %q; want %d, %q, %q",
					wiErr.StatusCode, wiErr.Code, wiErr.Description, tt.status, tt.wantCode, tt.wantDesc)
			}
			if want := tt.header.Get("X-Request-Id"); wiErr.RequestID != want {
				t.Errorf("RequestID = %q, want %q", wiErr.RequestID, want)
			}
			message := err.Error()
			for _, want := range tt.wantMessage {
				if !strings.Contains(message, want) {
					t.Errorf("error %q does not contain %q", message, want)
				}
			}
			for _, secret := range []string{assertion, "secret-access-token"} {
				if strings.Contains(message, secret) {
					t.Errorf("error %q leaks %q", message, secret)
				}
			}
			if requests != 1 {
				t.Errorf("token endpoint got %d requests, want 1", requests)
			}
		})
	}
}

func TestExchangeRejectsBadSuccessResponses(t *testing.T) {
	cleanEnv(t)
	const accessToken = "secret-access-token"
	tests := []struct {
		name string
		body string
		want string
	}{
		{"not JSON", "ok " + accessToken, "not valid JSON"},
		{"no access token", `{"token_type":"Bearer","expires_in":300}`, "no access_token"},
		{"wrong token type", `{"token_type":"mac","access_token":"` + accessToken + `","expires_in":300}`, "unsupported token_type"},
		{"expired", `{"token_type":"Bearer","access_token":"` + accessToken + `","expires_in":0}`, "no valid expires_in"},
		{"no lifetime", `{"token_type":"Bearer","access_token":"` + accessToken + `"}`, "no valid expires_in"},
		{"expires_in is not a number", `{"token_type":"Bearer","access_token":"` + accessToken + `","expires_in":"soon"}`, "not valid JSON"},
		{"larger than 1 MiB", `{"access_token":"` + accessToken + `","padding":"` + strings.Repeat("x", 1<<20) + `"}`, "not valid JSON"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()
			_, err := credentials.ExchangeFederationAssertion(context.Background(), credentials.FederationExchangeParams{
				Assertion:        "a.b.c",
				FederationRuleID: testRuleID,
				AuthBaseURL:      server.URL,
			})
			var wiErr *credentials.WorkloadIdentityError
			if !errors.As(err, &wiErr) || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want a WorkloadIdentityError containing %q", err, tt.want)
			}
			if strings.Contains(err.Error(), accessToken) {
				t.Fatalf("error %q leaks the access token", err)
			}
		})
	}
}

func TestExchangeAcceptsLenientSuccessResponses(t *testing.T) {
	cleanEnv(t)
	tests := []struct {
		name string
		body string
	}{
		{"lower-case token type", `{"token_type":"bearer","access_token":"t","expires_in":1}`},
		{"no token type", `{"access_token":"t","expires_in":600}`},
		{"fractional lifetime", `{"access_token":"t","expires_in":1.5}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()
			token, err := credentials.ExchangeFederationAssertion(context.Background(), credentials.FederationExchangeParams{
				Assertion:        "a.b.c",
				FederationRuleID: testRuleID,
				AuthBaseURL:      server.URL,
			})
			if err != nil {
				t.Fatalf("ExchangeFederationAssertion: %v", err)
			}
			if !token.ExpiresAt.After(time.Now()) {
				t.Fatalf("ExpiresAt = %v, want a time in the future", token.ExpiresAt)
			}
		})
	}
}

// TestExchangeNeverResendsTheAssertion covers HTTP clients the SDK cannot
// reconfigure: a wrapper around an http.Client that follows redirects.
func TestExchangeNeverResendsTheAssertion(t *testing.T) {
	cleanEnv(t)
	var elsewhere int
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		elsewhere++
	}))
	defer target.Close()
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Redirect(w, r, target.URL+"/token", http.StatusTemporaryRedirect)
	}))
	defer server.Close()

	following := &http.Client{}
	_, err := credentials.ExchangeFederationAssertion(context.Background(), credentials.FederationExchangeParams{
		Assertion:        "a.b.c",
		FederationRuleID: testRuleID,
		AuthBaseURL:      server.URL,
		HTTPClient:       doerFunc(following.Do),
	})
	var wiErr *credentials.WorkloadIdentityError
	if !errors.As(err, &wiErr) || wiErr.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("error = %v, want the 307", err)
	}
	if requests != 1 || elsewhere != 0 {
		t.Fatalf("%d requests to the token endpoint and %d to the redirect target, want 1 and 0", requests, elsewhere)
	}
}

func TestExchangeErrorsFromHTTPClientAreRedacted(t *testing.T) {
	cleanEnv(t)
	const assertion = "header.payload-that-must-not-leak.signature"
	cause := errors.New("proxy refused")
	t.Run("transport error", func(t *testing.T) {
		_, err := credentials.ExchangeFederationAssertion(context.Background(), credentials.FederationExchangeParams{
			Assertion:        assertion,
			FederationRuleID: testRuleID,
			AuthBaseURL:      "https://auth.example.com",
			HTTPClient: doerFunc(func(req *http.Request) (*http.Response, error) {
				body, _ := io.ReadAll(req.Body)
				return nil, fmt.Errorf("%w for body %s", cause, body)
			}),
		})
		if !errors.Is(err, cause) || !strings.Contains(err.Error(), "proxy refused") {
			t.Fatalf("error = %v, want it to wrap the client's error", err)
		}
		for _, format := range []string{"%v", "%+v", "%#v"} {
			if text := fmt.Sprintf(format, err); strings.Contains(text, assertion) {
				t.Fatalf("the error formatted with %s leaks the assertion: %s", format, text)
			}
		}
	})
	t.Run("request ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request-Id", "echo-"+assertion)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"error":"invalid_grant"}`)
		}))
		defer server.Close()
		_, err := credentials.ExchangeFederationAssertion(context.Background(), credentials.FederationExchangeParams{
			Assertion:        assertion,
			FederationRuleID: testRuleID,
			AuthBaseURL:      server.URL,
		})
		if err == nil || strings.Contains(err.Error(), assertion) || !strings.Contains(err.Error(), "request id echo-[redacted]") {
			t.Fatalf("error = %v, want the request ID without the assertion", err)
		}
	})
}

func TestExchangeRespectsContext(t *testing.T) {
	cleanEnv(t)
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	defer server.Close()
	defer close(release)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := credentials.ExchangeFederationAssertion(ctx, credentials.FederationExchangeParams{
		Assertion:        "a.b.c",
		FederationRuleID: testRuleID,
		AuthBaseURL:      server.URL,
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context.DeadlineExceeded", err)
	}
}
