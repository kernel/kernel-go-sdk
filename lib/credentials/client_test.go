package credentials_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	kernel "github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/lib/credentials"
	"github.com/kernel/kernel-go-sdk/option"
)

const testAPIKey = "sk_test_api_key"

// apiServer imitates the Kernel API: it accepts the API key and the access
// tokens its authorization server has issued and not revoked, and answers
// 401 to anything else.
type apiServer struct {
	*httptest.Server
	auth *authServer

	mu sync.Mutex
	// respond, when set, can answer request n (1-based) itself by returning a
	// status other than 0.
	respond  func(n int) int
	requests []apiRequest
}

type apiRequest struct {
	Method        string
	Path          string
	Authorization string
	Body          string
	Status        int
}

func newAPIServer(t *testing.T, auth *authServer) *apiServer {
	t.Helper()
	s := &apiServer{auth: auth}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		authorization := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authorization, "Bearer ")
		status := http.StatusOK
		s.mu.Lock()
		if token != testAPIKey && !s.auth.valid(token) {
			status = http.StatusUnauthorized
		}
		if s.respond != nil {
			if forced := s.respond(len(s.requests) + 1); forced != 0 {
				status = forced
			}
		}
		s.requests = append(s.requests, apiRequest{r.Method, r.URL.RequestURI(), authorization, string(body), status})
		s.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if status != http.StatusOK {
			w.WriteHeader(status)
			_, _ = io.WriteString(w, `{"code":"unauthorized","message":"invalid credentials"}`)
			return
		}
		if r.URL.Path == "/profiles" {
			// Three pages of one profile each.
			next := map[string]string{"": "1", "1": "2", "2": "0"}[r.URL.Query().Get("offset")]
			w.Header().Set("X-Next-Offset", next)
			_, _ = io.WriteString(w, `[{"id":"profile"}]`)
			return
		}
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *apiServer) seen() []apiRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]apiRequest(nil), s.requests...)
}

// authorizations lists the Authorization header and outcome of each request,
// such as "Bearer kernel-access-token-1 200".
func (s *apiServer) authorizations() []string {
	var out []string
	for _, r := range s.seen() {
		out = append(out, fmt.Sprintf("%s %d", r.Authorization, r.Status))
	}
	return out
}

type fixture struct {
	auth *authServer
	api  *apiServer
}

// newFixture starts an authorization server and an API, with a clean
// environment.
func newFixture(t *testing.T) *fixture {
	t.Helper()
	cleanEnv(t)
	auth := newAuthServer(t)
	return &fixture{auth: auth, api: newAPIServer(t, auth)}
}

// federatedClient returns a client that authenticates with source against the
// fixture's servers.
func (f *fixture) federatedClient(source credentials.IdentityTokenFunc, opts ...option.RequestOption) kernel.Client {
	opts = append([]option.RequestOption{
		option.WithBaseURL(f.api.URL),
		option.WithFederationTokenProvider(source, option.FederationOptions{
			FederationRuleID: testRuleID,
			AuthBaseURL:      f.auth.URL,
		}),
	}, opts...)
	return kernel.NewClient(opts...)
}

func get(client kernel.Client, opts ...option.RequestOption) error {
	var res map[string]any
	return client.Get(context.Background(), "auth/context", nil, &res, opts...)
}

func mustGet(t *testing.T, client kernel.Client, opts ...option.RequestOption) {
	t.Helper()
	if err := get(client, opts...); err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func equal(t *testing.T, what string, got, want []string) {
	t.Helper()
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("%s:\n  got  %q\n  want %q", what, got, want)
	}
}

func assertions(auth *authServer) []string {
	var out []string
	for _, r := range auth.exchanges() {
		out = append(out, r.Form.Get("assertion"))
	}
	return out
}

func distinct(values []string) int {
	set := map[string]bool{}
	for _, v := range values {
		set[v] = true
	}
	return len(set)
}

func TestClientGitHubActionsWithoutCode(t *testing.T) {
	f := newFixture(t)
	runner, runnerRequests := githubTokenService(t)
	t.Setenv("KERNEL_BASE_URL", f.api.URL)
	t.Setenv(credentials.EnvAuthBaseURL, f.auth.URL)
	t.Setenv(credentials.EnvFederationRuleID, testRuleID)
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", runner.URL+"/token?api-version=2.0")
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "runtime-token")

	client := kernel.NewClient()
	for i := 0; i < 3; i++ {
		mustGet(t, client)
	}

	equal(t, "API requests", f.api.authorizations(), []string{
		"Bearer kernel-access-token-1 200",
		"Bearer kernel-access-token-1 200",
		"Bearer kernel-access-token-1 200",
	})
	exchanges := f.auth.exchanges()
	if len(exchanges) != 1 || len(*runnerRequests) != 1 {
		t.Fatalf("%d exchanges and %d runner requests, want 1 and 1", len(exchanges), len(*runnerRequests))
	}
	form := exchanges[0].Form
	if got := form.Get("federation_rule_id"); got != testRuleID {
		t.Errorf("federation_rule_id = %q", got)
	}
	if got := form.Get("resource"); got != f.api.URL {
		t.Errorf("resource = %q, want the API origin %q", got, f.api.URL)
	}
	if got := tokenAudiences(form.Get("assertion")); len(got) != 1 || got[0] != f.auth.URL {
		t.Errorf("identity token audience = %v, want the auth base URL %q", got, f.auth.URL)
	}
}

func init() {
	// A source registered the way the awssts package registers aws_sts.
	credentials.RegisterIdentityTokenSource("test_registered", func(audience string) credentials.IdentityTokenFunc {
		return func(ctx context.Context) (string, error) {
			audience, err := credentials.ResolveAudience(ctx, audience)
			if err != nil {
				return "", err
			}
			return "registered.source.token-for-" + audience, nil
		}
	})
}

func TestClientEnvironmentSources(t *testing.T) {
	tests := []struct {
		name string
		// setup prepares the source and returns the assertion the token
		// endpoint should receive ("" when it cannot be known in advance).
		setup func(t *testing.T, f *fixture) (wantAssertion string)
	}{
		{
			name: "Kubernetes projected token file",
			setup: func(t *testing.T, f *fixture) string {
				token := providerToken(t, providerKubernetes, f.auth.URL)
				path := filepath.Join(t.TempDir(), "token")
				if err := os.WriteFile(path, []byte(token+"\n"), 0o600); err != nil {
					t.Fatal(err)
				}
				t.Setenv(credentials.EnvIdentityTokenFile, path)
				return token
			},
		},
		{
			name: "literal token",
			setup: func(t *testing.T, f *fixture) string {
				token := providerToken(t, providerAWS, f.auth.URL)
				t.Setenv(credentials.EnvIdentityToken, token)
				return token
			},
		},
		{
			name: "Google Cloud metadata server",
			setup: func(t *testing.T, f *fixture) string {
				var token string
				metadata := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					token = providerToken(t, providerGoogle, r.URL.Query().Get("audience"))
					_, _ = io.WriteString(w, token)
				}))
				t.Cleanup(metadata.Close)
				t.Setenv("GCE_METADATA_HOST", strings.TrimPrefix(metadata.URL, "http://"))
				t.Setenv(credentials.EnvIdentityTokenSource, "gcp_metadata")
				return ""
			},
		},
		{
			name: "registered source",
			setup: func(t *testing.T, f *fixture) string {
				t.Setenv(credentials.EnvIdentityTokenSource, "test_registered")
				return "registered.source.token-for-" + f.auth.URL
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			t.Setenv("KERNEL_BASE_URL", f.api.URL)
			t.Setenv(credentials.EnvAuthBaseURL, f.auth.URL)
			t.Setenv(credentials.EnvFederationRuleID, testRuleID)
			want := tt.setup(t, f)

			mustGet(t, kernel.NewClient())

			equal(t, "API requests", f.api.authorizations(), []string{"Bearer kernel-access-token-1 200"})
			got := assertions(f.auth)
			if len(got) != 1 {
				t.Fatalf("%d exchanges, want 1", len(got))
			}
			if want != "" && got[0] != want {
				t.Errorf("assertion = %q, want the source's token %q unchanged", got[0], want)
			}
			if aud := tokenAudiences(got[0]); want == "" && (len(aud) != 1 || aud[0] != f.auth.URL) {
				t.Errorf("identity token audience = %v, want the auth base URL", aud)
			}
		})
	}
}

func TestClientIncompleteEnvironmentFailsRequests(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want string
	}{
		{"no identity token source", nil, "KERNEL_FEDERATION_RULE_ID is set but no identity token source could be configured"},
		{"aws_sts without the import", map[string]string{credentials.EnvIdentityTokenSource: "aws_sts"}, `_ "github.com/kernel/kernel-go-sdk/lib/credentials/awssts"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			t.Setenv("KERNEL_BASE_URL", f.api.URL)
			t.Setenv(credentials.EnvAuthBaseURL, f.auth.URL)
			t.Setenv(credentials.EnvFederationRuleID, testRuleID)
			for key, value := range tt.env {
				t.Setenv(key, value)
			}
			client := kernel.NewClient()

			err := get(client)
			var wiErr *credentials.WorkloadIdentityError
			if !errors.As(err, &wiErr) || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want a WorkloadIdentityError containing %q", err, tt.want)
			}
			if n := len(f.api.seen()); n != 0 {
				t.Fatalf("%d unauthenticated requests reached the API", n)
			}
			// An explicit credential still works.
			mustGet(t, client, option.WithAPIKey(testAPIKey))
		})
	}
}

func TestClientAPIKeyInEnvironmentWinsAndWarnsOnce(t *testing.T) {
	f := newFixture(t)
	t.Setenv("KERNEL_BASE_URL", f.api.URL)
	t.Setenv("KERNEL_API_KEY", testAPIKey)
	t.Setenv(credentials.EnvAuthBaseURL, f.auth.URL)
	t.Setenv(credentials.EnvFederationRuleID, testRuleID)
	t.Setenv(credentials.EnvIdentityToken, "a.b.c")

	var logged bytes.Buffer
	log.SetOutput(&logged)
	defer log.SetOutput(os.Stderr)
	credentials.ResetEnvWarningForTest()

	for i := 0; i < 3; i++ {
		mustGet(t, kernel.NewClient())
	}

	equal(t, "API requests", f.api.authorizations(), []string{
		"Bearer " + testAPIKey + " 200", "Bearer " + testAPIKey + " 200", "Bearer " + testAPIKey + " 200",
	})
	if n := len(f.auth.exchanges()); n != 0 {
		t.Fatalf("%d exchanges, want 0", n)
	}
	if n := strings.Count(logged.String(), "KERNEL_API_KEY and KERNEL_FEDERATION_RULE_ID are both set; KERNEL_API_KEY takes precedence"); n != 1 {
		t.Fatalf("warning logged %d times, want once; log:\n%s", n, logged.String())
	}
}

func TestClientAPIKeyStillWorksAlone(t *testing.T) {
	f := newFixture(t)
	t.Setenv("KERNEL_BASE_URL", f.api.URL)
	t.Setenv("KERNEL_API_KEY", testAPIKey)

	var logged bytes.Buffer
	log.SetOutput(&logged)
	defer log.SetOutput(os.Stderr)
	credentials.ResetEnvWarningForTest()

	mustGet(t, kernel.NewClient())
	equal(t, "API requests", f.api.authorizations(), []string{"Bearer " + testAPIKey + " 200"})
	if logged.Len() != 0 {
		t.Fatalf("unexpected log output: %s", logged.String())
	}
}

func TestClientExplicitOptionsWinOverEnvironment(t *testing.T) {
	t.Run("WithAPIKey over federation variables", func(t *testing.T) {
		f := newFixture(t)
		t.Setenv(credentials.EnvAuthBaseURL, f.auth.URL)
		t.Setenv(credentials.EnvFederationRuleID, testRuleID)
		t.Setenv(credentials.EnvIdentityToken, "a.b.c")

		mustGet(t, kernel.NewClient(option.WithBaseURL(f.api.URL), option.WithAPIKey(testAPIKey)))
		equal(t, "API requests", f.api.authorizations(), []string{"Bearer " + testAPIKey + " 200"})
		if n := len(f.auth.exchanges()); n != 0 {
			t.Fatalf("%d exchanges, want 0", n)
		}
	})
	t.Run("WithFederationTokenProvider over KERNEL_API_KEY", func(t *testing.T) {
		f := newFixture(t)
		t.Setenv("KERNEL_API_KEY", testAPIKey)
		source := &freshTokens{t: t, provider: providerGitHub}

		mustGet(t, f.federatedClient(source.get))
		equal(t, "API requests", f.api.authorizations(), []string{"Bearer kernel-access-token-1 200"})
	})
	t.Run("WithFederationTokenProvider over federation variables", func(t *testing.T) {
		f := newFixture(t)
		t.Setenv(credentials.EnvFederationRuleID, "wifr_from_env")
		t.Setenv(credentials.EnvIdentityTokenFile, filepath.Join(t.TempDir(), "missing"))
		source := &freshTokens{t: t, provider: providerGitHub}

		mustGet(t, f.federatedClient(source.get))
		exchanges := f.auth.exchanges()
		if len(exchanges) != 1 || exchanges[0].Form.Get("federation_rule_id") != testRuleID {
			t.Fatalf("exchanges = %+v, want one for the explicit rule", exchanges)
		}
	})
}

func TestClientLastCredentialOptionWins(t *testing.T) {
	f := newFixture(t)
	source := &freshTokens{t: t, provider: providerGitHub}
	federation := option.WithFederationTokenProvider(source.get, option.FederationOptions{
		FederationRuleID: testRuleID,
		AuthBaseURL:      f.auth.URL,
	})
	base := option.WithBaseURL(f.api.URL)

	mustGet(t, kernel.NewClient(base, option.WithAPIKey("sk_replaced"), federation))
	mustGet(t, kernel.NewClient(base, federation, option.WithAPIKey(testAPIKey)))
	// A request option replaces the client's credential in both directions.
	mustGet(t, kernel.NewClient(base, federation), option.WithAPIKey(testAPIKey))
	mustGet(t, kernel.NewClient(base, option.WithAPIKey("sk_replaced")), federation)
	mustGet(t, kernel.NewClient(base, federation), option.WithHeader("Authorization", "Bearer "+testAPIKey))

	equal(t, "API requests", f.api.authorizations(), []string{
		"Bearer kernel-access-token-1 200",
		"Bearer " + testAPIKey + " 200",
		"Bearer " + testAPIKey + " 200",
		"Bearer kernel-access-token-1 200",
		"Bearer " + testAPIKey + " 200",
	})
	if n := len(f.auth.exchanges()); n != 1 {
		t.Fatalf("%d exchanges, want 1: the option's token cache is shared by the clients using it", n)
	}
}

func TestClientRejectsBadFederationOptions(t *testing.T) {
	f := newFixture(t)
	source := &freshTokens{t: t, provider: providerGitHub}
	tests := []struct {
		name    string
		baseURL string
		opt     option.RequestOption
		want    string
	}{
		{name: "nil provider", opt: option.WithFederationTokenProvider(nil, option.FederationOptions{FederationRuleID: testRuleID}), want: "provider is nil"},
		{name: "no rule", opt: option.WithFederationTokenProvider(source.get, option.FederationOptions{}), want: "FederationRuleID is required"},
		{
			name: "custom base URL without an auth base URL",
			opt:  option.WithFederationTokenProvider(source.get, option.FederationOptions{FederationRuleID: testRuleID}),
			want: "KERNEL_AUTH_BASE_URL",
		},
		{
			name:    "path on the production host without an auth base URL",
			baseURL: "https://api.onkernel.com/gateway/",
			opt:     option.WithFederationTokenProvider(source.get, option.FederationOptions{FederationRuleID: testRuleID}),
			want:    "KERNEL_AUTH_BASE_URL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.baseURL == "" {
				tt.baseURL = f.api.URL
			}
			started := time.Now()
			err := get(kernel.NewClient(option.WithBaseURL(tt.baseURL), tt.opt))
			var wiErr *credentials.WorkloadIdentityError
			if !errors.As(err, &wiErr) || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want a WorkloadIdentityError containing %q", err, tt.want)
			}
			if elapsed := time.Since(started); elapsed > 300*time.Millisecond {
				t.Fatalf("the request took %v: a configuration error must not be retried", elapsed)
			}
		})
	}
	if source.count() != 0 || len(f.auth.exchanges()) != 0 || len(f.api.seen()) != 0 {
		t.Fatalf("a misconfigured client reached out: %d identity tokens, %d exchanges, %d API requests",
			source.count(), len(f.auth.exchanges()), len(f.api.seen()))
	}
}

func TestClientRetriesOnceAfter401WithFreshAssertion(t *testing.T) {
	f := newFixture(t)
	source := &freshTokens{t: t, provider: providerGitHub}
	client := f.federatedClient(source.get)

	mustGet(t, client)
	f.auth.revokeAll()

	body := map[string]any{"name": "replayed"}
	var res map[string]any
	if err := client.Post(context.Background(), "browsers", body, &res); err != nil {
		t.Fatalf("request after revocation failed: %v", err)
	}

	equal(t, "API requests", f.api.authorizations(), []string{
		"Bearer kernel-access-token-1 200",
		"Bearer kernel-access-token-1 401",
		"Bearer kernel-access-token-2 200",
	})
	requests := f.api.seen()
	if requests[1].Body == "" || requests[1].Body != requests[2].Body {
		t.Fatalf("the retry sent body %q, the first attempt %q", requests[2].Body, requests[1].Body)
	}
	if got := assertions(f.auth); len(got) != 2 || distinct(got) != 2 {
		t.Fatalf("assertions = %q, want two different ones", got)
	}
}

func TestClientSecond401Surfaces(t *testing.T) {
	f := newFixture(t)
	source := &freshTokens{t: t, provider: providerGitHub}
	client := f.federatedClient(source.get)
	// Nothing the authorization server issues is accepted by this API.
	f.api.mu.Lock()
	f.api.auth = newAuthServer(t)
	f.api.mu.Unlock()

	err := get(client)
	var apiErr *kernel.Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("error = %v, want the API's 401", err)
	}
	equal(t, "API requests", f.api.authorizations(), []string{
		"Bearer kernel-access-token-1 401",
		"Bearer kernel-access-token-2 401",
	})
	if got := assertions(f.auth); len(got) != 2 || distinct(got) != 2 {
		t.Fatalf("assertions = %q, want two different ones", got)
	}
}

// TestClientReplaysOnceWhenRetriesInterleave pins the limit of one replay per
// request: the replay itself fails with a retryable status, and the client's
// own retry is rejected again.
func TestClientReplaysOnceWhenRetriesInterleave(t *testing.T) {
	f := newFixture(t)
	// Nothing the authorization server issues is accepted by this API, and
	// its second answer is a 503.
	f.api.mu.Lock()
	f.api.auth = newAuthServer(t)
	f.api.respond = func(n int) int {
		if n == 2 {
			return http.StatusServiceUnavailable
		}
		return 0
	}
	f.api.mu.Unlock()
	source := &freshTokens{t: t, provider: providerGitHub}

	err := get(f.federatedClient(source.get))
	var apiErr *kernel.Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("error = %v, want the API's 401", err)
	}
	equal(t, "API requests", f.api.authorizations(), []string{
		"Bearer kernel-access-token-1 401",
		"Bearer kernel-access-token-2 503",
		"Bearer kernel-access-token-2 401",
	})
	if n := len(f.auth.exchanges()); n != 2 {
		t.Fatalf("%d exchanges, want 2", n)
	}
}

func TestClientDoesNotRetry401WhenBodyCannotBeReplayed(t *testing.T) {
	f := newFixture(t)
	source := &freshTokens{t: t, provider: providerGitHub}
	client := f.federatedClient(source.get)
	mustGet(t, client)
	f.auth.revokeAll()

	// A plain reader has no GetBody, so the request cannot be sent twice.
	stream := io.MultiReader(strings.NewReader("streamed"))
	err := client.Post(context.Background(), "browsers", nil, nil, option.WithRequestBody("application/octet-stream", stream))
	var apiErr *kernel.Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("error = %v, want the API's 401", err)
	}
	// The rejected token is still dropped, so the next request recovers.
	mustGet(t, client)
	equal(t, "API requests", f.api.authorizations(), []string{
		"Bearer kernel-access-token-1 200",
		"Bearer kernel-access-token-1 401",
		"Bearer kernel-access-token-2 200",
	})
}

func TestClientConcurrentRequestsShareOneExchange(t *testing.T) {
	f := newFixture(t)
	source := &freshTokens{t: t, provider: providerGitHub}
	client := f.federatedClient(source.get)

	run := func() {
		var wg sync.WaitGroup
		for i := 0; i < 30; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := get(client); err != nil {
					t.Errorf("request failed: %v", err)
				}
			}()
		}
		wg.Wait()
	}

	run()
	if n := len(f.auth.exchanges()); n != 1 {
		t.Fatalf("%d exchanges for 30 concurrent requests, want 1", n)
	}

	// All 30 are rejected with the revoked token; one refresh serves them all.
	f.auth.revokeAll()
	run()
	if n := len(f.auth.exchanges()); n != 2 {
		t.Fatalf("%d exchanges after a revocation seen by 30 concurrent requests, want 2", n)
	}
	if n := source.count(); n != 2 {
		t.Fatalf("%d identity tokens fetched, want 2", n)
	}
}

func TestClientRefreshesBeforeExpiry(t *testing.T) {
	f := newFixture(t)
	f.auth.expiresIn = 0.4
	source := &freshTokens{t: t, provider: providerGitHub}
	client := f.federatedClient(source.get)

	mustGet(t, client)
	mustGet(t, client)
	time.Sleep(450 * time.Millisecond)
	mustGet(t, client)

	// The expired token is never sent: the client goes by expires_in.
	equal(t, "API requests", f.api.authorizations(), []string{
		"Bearer kernel-access-token-1 200",
		"Bearer kernel-access-token-1 200",
		"Bearer kernel-access-token-2 200",
	})
	if got := assertions(f.auth); distinct(got) != 2 {
		t.Fatalf("assertions = %q, want two different ones", got)
	}
}

func TestClientExchangeRefusalIsNotRetried(t *testing.T) {
	f := newFixture(t)
	f.auth.respond = func(int) (int, string) {
		return 400, `{"error":"invalid_grant","error_description":"Credential is invalid, expired, or revoked"}`
	}
	source := &freshTokens{t: t, provider: providerGitHub}

	err := get(f.federatedClient(source.get))
	var wiErr *credentials.WorkloadIdentityError
	if !errors.As(err, &wiErr) || wiErr.StatusCode != 400 || wiErr.Code != "invalid_grant" {
		t.Fatalf("error = %v, want the invalid_grant refusal", err)
	}
	if exchanges, tokens := len(f.auth.exchanges()), source.count(); exchanges != 1 || tokens != 1 {
		t.Fatalf("%d exchanges with %d identity tokens, want 1 and 1", exchanges, tokens)
	}
	if n := len(f.api.seen()); n != 0 {
		t.Fatalf("%d API requests, want 0", n)
	}
}

func TestClientRetriesUnavailableTokenEndpointWithFreshAssertions(t *testing.T) {
	f := newFixture(t)
	f.auth.respond = func(n int) (int, string) {
		if n == 1 {
			return 503, `{"error":"temporarily_unavailable","error_description":"Try again"}`
		}
		return 0, ""
	}
	source := &freshTokens{t: t, provider: providerGitHub}

	mustGet(t, f.federatedClient(source.get))

	if got := assertions(f.auth); len(got) != 2 || distinct(got) != 2 {
		t.Fatalf("assertions = %q, want two different ones", got)
	}
	equal(t, "API requests", f.api.authorizations(), []string{"Bearer kernel-access-token-2 200"})
}

func TestClientRetriesUnavailableIdentityProvider(t *testing.T) {
	f := newFixture(t)
	var runnerRequests int
	runner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		runnerRequests++
		if runnerRequests == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = fmt.Fprintf(w, `{"value":%q}`, providerToken(t, providerGitHub, r.URL.Query().Get("audience")))
	}))
	defer runner.Close()
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", runner.URL+"/token?api-version=2.0")
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "runtime-token")

	mustGet(t, f.federatedClient(credentials.GitHubActions("")))

	if runnerRequests != 2 || len(f.auth.exchanges()) != 1 {
		t.Fatalf("%d runner requests and %d exchanges, want 2 and 1", runnerRequests, len(f.auth.exchanges()))
	}

	// A refusal, unlike an outage, is not retried.
	forbidden := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		runnerRequests++
		w.WriteHeader(http.StatusForbidden)
	}))
	defer forbidden.Close()
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", forbidden.URL+"/token?api-version=2.0")
	runnerRequests = 0
	err := get(f.federatedClient(credentials.GitHubActions("")))
	if err == nil || !strings.Contains(err.Error(), "status 403") || runnerRequests != 1 {
		t.Fatalf("error = %v after %d runner requests, want the 403 after 1", err, runnerRequests)
	}
}

func TestClientStaticTokenIsNeverReplayed(t *testing.T) {
	t.Run("after a revocation", func(t *testing.T) {
		f := newFixture(t)
		t.Setenv("KERNEL_BASE_URL", f.api.URL)
		t.Setenv(credentials.EnvAuthBaseURL, f.auth.URL)
		t.Setenv(credentials.EnvFederationRuleID, testRuleID)
		t.Setenv(credentials.EnvIdentityToken, providerToken(t, providerGitHub, f.auth.URL))
		client := kernel.NewClient()

		mustGet(t, client)
		f.auth.revokeAll()
		err := get(client)

		var wiErr *credentials.WorkloadIdentityError
		if !errors.As(err, &wiErr) || !errors.Is(err, credentials.ErrIdentityTokenConsumed) {
			t.Fatalf("error = %v, want a WorkloadIdentityError wrapping ErrIdentityTokenConsumed", err)
		}
		if !strings.Contains(err.Error(), "exchanged only once") {
			t.Fatalf("error %q does not explain why", err)
		}
		if n := len(f.auth.exchanges()); n != 1 {
			t.Fatalf("%d exchanges, want 1: the token must not be sent twice", n)
		}
	})
	t.Run("by a second client", func(t *testing.T) {
		f := newFixture(t)
		t.Setenv("KERNEL_BASE_URL", f.api.URL)
		t.Setenv(credentials.EnvAuthBaseURL, f.auth.URL)
		t.Setenv(credentials.EnvFederationRuleID, testRuleID)
		t.Setenv(credentials.EnvIdentityToken, providerToken(t, providerGitHub, f.auth.URL))

		mustGet(t, kernel.NewClient())
		err := get(kernel.NewClient())

		if !errors.Is(err, credentials.ErrIdentityTokenConsumed) {
			t.Fatalf("second client's error = %v, want ErrIdentityTokenConsumed", err)
		}
		if n := len(f.auth.exchanges()); n != 1 {
			t.Fatalf("%d exchanges, want 1: the token must not be sent twice", n)
		}
	})
	t.Run("after a failed exchange", func(t *testing.T) {
		f := newFixture(t)
		f.auth.respond = func(int) (int, string) {
			return 503, `{"error":"temporarily_unavailable","error_description":"Try again"}`
		}
		client := f.federatedClient(credentials.Static(providerToken(t, providerGitHub, f.auth.URL)))

		err := get(client)
		if !errors.Is(err, credentials.ErrIdentityTokenConsumed) {
			t.Fatalf("error = %v, want ErrIdentityTokenConsumed", err)
		}
		for _, want := range []string{"The exchange that used it failed", "status 503", "temporarily_unavailable"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not contain %q", err, want)
			}
		}
		if n := len(f.auth.exchanges()); n != 1 {
			t.Fatalf("%d exchanges, want 1: the token must not be sent twice", n)
		}
	})
}

func TestClientKeepsTokenOnItsOwnOrigin(t *testing.T) {
	f := newFixture(t)
	var gotAuthorization []string
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuthorization = append(gotAuthorization, r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer other.Close()
	source := &freshTokens{t: t, provider: providerGitHub}
	client := f.federatedClient(source.get)
	mustGet(t, client)

	var res map[string]any
	if err := client.Get(context.Background(), other.URL+"/elsewhere", nil, &res); err != nil {
		t.Fatalf("request to another origin failed: %v", err)
	}
	equal(t, "Authorization sent to another origin", gotAuthorization, []string{""})
}

func TestClientSeparatesTokensByAPI(t *testing.T) {
	f := newFixture(t)
	secondAPI := newAPIServer(t, f.auth)
	source := &freshTokens{t: t, provider: providerGitHub}
	federation := option.WithFederationTokenProvider(source.get, option.FederationOptions{
		FederationRuleID: testRuleID,
		AuthBaseURL:      f.auth.URL,
	})

	mustGet(t, kernel.NewClient(option.WithBaseURL(f.api.URL), federation))
	mustGet(t, kernel.NewClient(option.WithBaseURL(secondAPI.URL), federation))

	exchanges := f.auth.exchanges()
	if len(exchanges) != 2 {
		t.Fatalf("%d exchanges, want one per API", len(exchanges))
	}
	if first, second := exchanges[0].Form.Get("resource"), exchanges[1].Form.Get("resource"); first != f.api.URL || second != secondAPI.URL {
		t.Fatalf("resources = %q, %q; want each client's own API", first, second)
	}
	equal(t, "second API requests", secondAPI.authorizations(), []string{"Bearer kernel-access-token-2 200"})
}

func TestClientPaginationKeepsCredential(t *testing.T) {
	f := newFixture(t)
	source := &freshTokens{t: t, provider: providerGitHub}
	client := f.federatedClient(source.get)

	pager := client.Profiles.ListAutoPaging(context.Background(), kernel.ProfileListParams{})
	var profiles int
	for pager.Next() {
		profiles++
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("pagination failed: %v", err)
	}
	if profiles != 3 {
		t.Fatalf("%d profiles, want 3", profiles)
	}
	equal(t, "API requests", f.api.authorizations(), []string{
		"Bearer kernel-access-token-1 200",
		"Bearer kernel-access-token-1 200",
		"Bearer kernel-access-token-1 200",
	})
}

// TestClientPaginationReplaysOncePerPage checks that the pages of one call,
// which share the middleware, each get their own replay.
func TestClientPaginationReplaysOncePerPage(t *testing.T) {
	tests := []struct {
		name string
		opts []option.RequestOption
	}{
		{name: "default headers"},
		{name: "without the retry count header", opts: []option.RequestOption{option.WithHeaderDel("X-Stainless-Retry-Count")}},
		{
			name: "with a later middleware that strips the idempotency key",
			opts: []option.RequestOption{option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				req.Header.Del("Idempotency-Key")
				return next(req)
			})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			source := &freshTokens{t: t, provider: providerGitHub}
			client := f.federatedClient(source.get, tt.opts...)

			page, err := client.Profiles.List(context.Background(), kernel.ProfileListParams{})
			if err != nil {
				t.Fatalf("first page: %v", err)
			}
			for i := 2; i <= 3; i++ {
				f.auth.revokeAll()
				if page, err = page.GetNextPage(); err != nil {
					t.Fatalf("page %d after a revocation: %v", i, err)
				}
			}
			equal(t, "API requests", f.api.authorizations(), []string{
				"Bearer kernel-access-token-1 200",
				"Bearer kernel-access-token-1 401",
				"Bearer kernel-access-token-2 200",
				"Bearer kernel-access-token-2 401",
				"Bearer kernel-access-token-3 200",
			})
		})
	}
}

func TestClientRequestContextBoundsTokenExchange(t *testing.T) {
	cleanEnv(t)
	release := make(chan struct{})
	auth := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	defer auth.Close()
	defer close(release)
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("an unauthenticated request reached the API")
	}))
	defer api.Close()

	source := &freshTokens{t: t, provider: providerGitHub}
	client := kernel.NewClient(
		option.WithBaseURL(api.URL),
		option.WithFederationTokenProvider(source.get, option.FederationOptions{FederationRuleID: testRuleID, AuthBaseURL: auth.URL}),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	started := time.Now()
	var res map[string]any
	err := client.Get(ctx, "auth/context", nil, &res)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(started); elapsed > 2*time.Second {
		t.Fatalf("the request returned after %v, long after its deadline", elapsed)
	}
}

func TestClientDebugLogDoesNotSeeTokens(t *testing.T) {
	f := newFixture(t)
	source := &freshTokens{t: t, provider: providerGitHub}
	var logged bytes.Buffer
	client := f.federatedClient(source.get, option.WithDebugLog(log.New(&logged, "", 0)))

	mustGet(t, client)

	if logged.Len() == 0 {
		t.Fatal("the debug logger saw nothing")
	}
	for _, secret := range append(assertions(f.auth), "kernel-access-token-1") {
		if strings.Contains(logged.String(), secret) {
			t.Fatalf("the debug log contains a token:\n%s", logged.String())
		}
	}
}
