package credentials_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kernel/kernel-go-sdk/lib/credentials"
)

const testRuleID = "wifr_test123"

// federationEnv is every environment variable the credential chain reads.
var federationEnv = []string{
	"KERNEL_API_KEY",
	"KERNEL_BASE_URL",
	"KERNEL_CUSTOM_HEADERS",
	credentials.EnvFederationRuleID,
	credentials.EnvIdentityTokenSource,
	credentials.EnvIdentityTokenFile,
	credentials.EnvIdentityToken,
	credentials.EnvIdentityTokenAudience,
	credentials.EnvAuthBaseURL,
	"ACTIONS_ID_TOKEN_REQUEST_URL",
	"ACTIONS_ID_TOKEN_REQUEST_TOKEN",
	"GCE_METADATA_HOST",
}

// cleanEnv unsets the credential chain's variables for the test, which may
// itself be running in a GitHub Actions job.
func cleanEnv(t *testing.T) {
	t.Helper()
	for _, key := range federationEnv {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}
}

// Identity-token fixtures. The header and claim shapes are the ones the real
// providers issue; the signature is a placeholder because nothing in the SDK
// verifies it.

type provider string

const (
	providerGitHub     provider = "github_actions"
	providerAWS        provider = "aws"
	providerGoogle     provider = "google"
	providerKubernetes provider = "kubernetes"
)

var jwtSerial struct {
	sync.Mutex
	n int
}

func providerToken(t *testing.T, p provider, audience string) string {
	t.Helper()
	jwtSerial.Lock()
	jwtSerial.n++
	jti := fmt.Sprintf("jti-%d", jwtSerial.n)
	jwtSerial.Unlock()

	iat := time.Now().Unix()
	var header, claims map[string]any
	switch p {
	case providerGitHub:
		header = map[string]any{"alg": "RS256", "kid": "38826b17-6a30-5f9b-b169-8beb8202f723", "typ": "JWT", "x5t": "ykNaY4qM_ta4k2TgZOCEYLkcYlA"}
		claims = map[string]any{
			"iss":                 "https://token.actions.githubusercontent.com",
			"aud":                 audience,
			"sub":                 "repo:kernel@123/app@456:ref:refs/heads/main",
			"repository":          "kernel/app",
			"repository_id":       "456",
			"repository_owner_id": "123",
			"jti":                 jti,
			"iat":                 iat,
			"nbf":                 iat - 300,
			"exp":                 iat + 300,
		}
	case providerAWS:
		header = map[string]any{"alg": "RS256", "kid": "RSA_0", "typ": "JWT"}
		claims = map[string]any{
			"iss": "https://abcd1234-5678-90ab-cdef-1234567890ab.tokens.sts.global.api.aws",
			"aud": audience,
			"sub": "arn:aws:iam::123456789012:role/service/team/app-role",
			"https://sts.amazonaws.com/": map[string]any{
				"aws_account":   "123456789012",
				"source_region": "us-east-1",
			},
			"jti": jti,
			"iat": iat,
			"exp": iat + 300,
		}
	case providerGoogle:
		header = map[string]any{"alg": "RS256", "kid": "c5d6e7f8a9b0c1d2e3f4", "typ": "JWT"}
		claims = map[string]any{
			"iss":            "https://accounts.google.com",
			"aud":            audience,
			"azp":            "117000000000000000000",
			"sub":            "117000000000000000000",
			"email":          "app@project.iam.gserviceaccount.com",
			"email_verified": true,
			"iat":            iat,
			"exp":            iat + 3600,
		}
	case providerKubernetes:
		header = map[string]any{"alg": "RS256", "kid": "kXw2yT0aV3"}
		claims = map[string]any{
			"iss": "https://kubernetes.default.svc.cluster.local",
			"aud": []string{audience},
			"sub": "system:serviceaccount:default:app",
			"kubernetes.io": map[string]any{
				"namespace":      "default",
				"serviceaccount": map[string]any{"name": "app", "uid": "5b0f1c2e"},
			},
			"jti": jti,
			"iat": iat,
			"nbf": iat,
			"exp": iat + 3600,
		}
	default:
		t.Fatalf("unknown provider %q", p)
	}
	return jwtSegment(t, header) + "." + jwtSegment(t, claims) + ".c2lnbmF0dXJl"
}

func jwtSegment(t *testing.T, v any) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return base64.RawURLEncoding.EncodeToString(data)
}

// tokenAudiences returns the aud claim of a fixture token as a list.
func tokenAudiences(token string) []string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil
	}
	var claims struct {
		Aud any `json:"aud"`
	}
	if json.Unmarshal(payload, &claims) != nil {
		return nil
	}
	switch aud := claims.Aud.(type) {
	case string:
		return []string{aud}
	case []any:
		var out []string
		for _, a := range aud {
			if s, ok := a.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// authServer imitates the Kernel authorization server's /token endpoint,
// including the checks that shape the client: exactly four parameters, no
// Authorization header, and single-use assertions.
type authServer struct {
	*httptest.Server
	t *testing.T

	mu sync.Mutex
	// expiresIn is the lifetime of issued tokens, in seconds.
	expiresIn float64
	// allowReuse turns replay protection off, as for Google tokens.
	allowReuse bool
	// respond, when set, can answer exchange n (1-based) itself by returning
	// a status other than 0.
	respond func(n int) (status int, body string)
	// resource is the only resource the server accepts, when set.
	resource string

	requests []tokenRequest
	seen     map[string]bool
	issued   map[string]bool
}

type tokenRequest struct {
	Header http.Header
	Form   url.Values
	Body   string
}

func newAuthServer(t *testing.T) *authServer {
	t.Helper()
	s := &authServer{t: t, expiresIn: 298, seen: map[string]bool{}, issued: map[string]bool{}}
	s.Server = httptest.NewServer(http.HandlerFunc(s.handle))
	t.Cleanup(s.Close)
	return s
}

func oauthError(w http.ResponseWriter, status int, code, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": code, "error_description": description})
}

func (s *authServer) handle(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	form, _ := url.ParseQuery(string(body))

	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = append(s.requests, tokenRequest{Header: r.Header.Clone(), Form: form, Body: string(body)})
	n := len(s.requests)

	if r.Method != http.MethodPost || r.URL.Path != "/token" {
		http.NotFound(w, r)
		return
	}
	if r.Header.Get("Authorization") != "" {
		oauthError(w, 400, "invalid_request", "Unsupported client authentication")
		return
	}
	if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
		oauthError(w, 400, "invalid_request", "Unsupported content type")
		return
	}
	for key, values := range form {
		switch key {
		case "grant_type", "assertion", "federation_rule_id", "resource":
		default:
			oauthError(w, 400, "invalid_request", "Unsupported parameter: "+key)
			return
		}
		if len(values) != 1 {
			oauthError(w, 400, "invalid_request", "Repeated parameter: "+key)
			return
		}
	}
	assertion := form.Get("assertion")
	if form.Get("grant_type") != "urn:ietf:params:oauth:grant-type:jwt-bearer" || assertion == "" || form.Get("federation_rule_id") == "" {
		oauthError(w, 400, "invalid_request", "Missing required parameter")
		return
	}
	if resource, ok := form["resource"]; ok && s.resource != "" && resource[0] != s.resource {
		oauthError(w, 400, "invalid_target", "Unknown resource")
		return
	}
	if s.respond != nil {
		if status, body := s.respond(n); status != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = io.WriteString(w, body)
			return
		}
	}
	if s.seen[assertion] && !s.allowReuse {
		oauthError(w, 400, "invalid_grant", "Credential was already used")
		return
	}
	s.seen[assertion] = true

	accessToken := fmt.Sprintf("kernel-access-token-%d", n)
	s.issued[accessToken] = true
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"token_type":         "Bearer",
		"access_token":       accessToken,
		"expires_in":         s.expiresIn,
		"scope":              "browsers:create browsers:use browsers:delete",
		"org_id":             "org_1",
		"access_scope":       "project",
		"project_id":         "proj_1",
		"principal":          map[string]string{"type": "service_account", "id": "svc_1"},
		"federation_rule_id": form.Get("federation_rule_id"),
	})
}

func (s *authServer) exchanges() []tokenRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]tokenRequest(nil), s.requests...)
}

// revokeAll makes every token issued so far invalid at the API.
func (s *authServer) revokeAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.issued = map[string]bool{}
}

func (s *authServer) valid(accessToken string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.issued[accessToken]
}

// freshTokens is an identity-token source that mints a new fixture token,
// for the audience the SDK asks for, on every call.
type freshTokens struct {
	t        *testing.T
	provider provider

	mu     sync.Mutex
	issued []string
}

func (f *freshTokens) get(ctx context.Context) (string, error) {
	audience, err := credentials.ResolveAudience(ctx, "")
	if err != nil {
		return "", err
	}
	token := providerToken(f.t, f.provider, audience)
	f.mu.Lock()
	defer f.mu.Unlock()
	f.issued = append(f.issued, token)
	return token, nil
}

func (f *freshTokens) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.issued)
}
