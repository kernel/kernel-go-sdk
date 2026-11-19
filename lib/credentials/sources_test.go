package credentials_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kernel/kernel-go-sdk/lib/credentials"
)

func TestFileRereadsOnEveryCall(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	source := credentials.File(path)

	if _, err := source(context.Background()); err == nil {
		t.Fatal("a missing file returned a token")
	}
	for _, content := range []string{"first.jwt.token\n", "  second.jwt.token  \n\n"} {
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		got, err := source(context.Background())
		if err != nil {
			t.Fatalf("File: %v", err)
		}
		if want := strings.TrimSpace(content); got != want {
			t.Fatalf("token = %q, want %q", got, want)
		}
	}
	if err := os.WriteFile(path, []byte(" \n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := source(context.Background()); err == nil || !strings.Contains(err.Error(), "is empty") {
		t.Fatalf("error for an empty file = %v", err)
	}
	if _, err := credentials.File("")(context.Background()); err == nil {
		t.Fatal("an empty path returned a token")
	}
}

// uniqueToken returns a token no earlier test, or earlier run of this test,
// has used: a static token is remembered for the life of the process.
func uniqueToken(label string) string {
	return fmt.Sprintf("%s.%d.jwt", label, time.Now().UnixNano())
}

func TestStaticIsGoodForOneUse(t *testing.T) {
	token := uniqueToken("only.once")
	source := credentials.Static(token + "\n")
	got, err := source(context.Background())
	if err != nil || got != token {
		t.Fatalf("first call = %q, %v", got, err)
	}
	// A second source for the same token is what a second client built from
	// the same KERNEL_IDENTITY_TOKEN gets.
	for i, source := range []credentials.IdentityTokenFunc{source, credentials.Static(token)} {
		_, err := source(context.Background())
		if !errors.Is(err, credentials.ErrIdentityTokenConsumed) {
			t.Fatalf("call %d error = %v, want ErrIdentityTokenConsumed", i+2, err)
		}
		if !strings.Contains(err.Error(), "exchanged only once") || !strings.Contains(err.Error(), "KERNEL_IDENTITY_TOKEN_FILE") {
			t.Fatalf("error %q does not explain the limit or the alternative", err)
		}
	}
	if _, err := credentials.Static(uniqueToken("another"))(context.Background()); err != nil {
		t.Fatalf("a different token was refused: %v", err)
	}
	for range 2 {
		if _, err := credentials.Static(" ")(context.Background()); err == nil || errors.Is(err, credentials.ErrIdentityTokenConsumed) {
			t.Fatalf("error for an empty value = %v", err)
		}
	}
}

func TestResolveAudience(t *testing.T) {
	cleanEnv(t)
	ctx := context.Background()
	if _, err := credentials.ResolveAudience(ctx, ""); err == nil {
		t.Fatal("no audience anywhere, but none was reported missing")
	}
	if got, _ := credentials.ResolveAudience(ctx, "explicit"); got != "explicit" {
		t.Fatalf("audience = %q, want explicit", got)
	}
	t.Setenv(credentials.EnvIdentityTokenAudience, "from-env")
	if got, _ := credentials.ResolveAudience(ctx, ""); got != "from-env" {
		t.Fatalf("audience = %q, want from-env", got)
	}
	if got, _ := credentials.ResolveAudience(ctx, "explicit"); got != "explicit" {
		t.Fatalf("audience = %q, want explicit over the environment", got)
	}
}

// githubTokenService imitates the GitHub Actions runner's OIDC endpoint. It
// records the requests it gets and answers with a fixture token for the
// requested audience.
func githubTokenService(t *testing.T) (*httptest.Server, *[]*http.Request) {
	t.Helper()
	var requests []*http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r)
		if r.Header.Get("Authorization") != "Bearer runtime-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"count": 1,
			"value": providerToken(t, providerGitHub, r.URL.Query().Get("audience")),
		})
	}))
	t.Cleanup(server.Close)
	return server, &requests
}

func TestGitHubActions(t *testing.T) {
	cleanEnv(t)
	server, requests := githubTokenService(t)
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", server.URL+"/_apis/oidc/token?api-version=2.0")
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "runtime-token")

	const audience = "https://auth.onkernel.com/?a=b&c"
	source := credentials.GitHubActions(audience)
	first, err := source(context.Background())
	if err != nil {
		t.Fatalf("GitHubActions: %v", err)
	}
	second, err := source(context.Background())
	if err != nil {
		t.Fatalf("GitHubActions: %v", err)
	}
	if first == second {
		t.Fatal("two calls returned the same token; every call must fetch a new one")
	}
	if got := tokenAudiences(first); len(got) != 1 || got[0] != audience {
		t.Fatalf("token audience = %v, want [%s]", got, audience)
	}

	req := (*requests)[0]
	if req.Method != http.MethodGet || req.URL.Path != "/_apis/oidc/token" {
		t.Errorf("request = %s %s", req.Method, req.URL.Path)
	}
	if want := "api-version=2.0&audience=https%3A%2F%2Fauth.onkernel.com%2F%3Fa%3Db%26c"; req.URL.RawQuery != want {
		t.Errorf("query = %q, want %q", req.URL.RawQuery, want)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer runtime-token" {
		t.Errorf("Authorization = %q", got)
	}
}

func TestGitHubActionsOutsideAWorkflow(t *testing.T) {
	cleanEnv(t)
	_, err := credentials.GitHubActions("aud")(context.Background())
	if err == nil || !strings.Contains(err.Error(), "id-token: write") {
		t.Fatalf("error = %v, want one that names the id-token: write permission", err)
	}
}

func TestGitHubActionsFailureDoesNotLeakRuntimeToken(t *testing.T) {
	cleanEnv(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("denied for runtime-token"))
	}))
	defer server.Close()
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", server.URL+"/token?api-version=2.0")
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "runtime-token")

	_, err := credentials.GitHubActions("aud")(context.Background())
	if err == nil || !strings.Contains(err.Error(), "status 403") {
		t.Fatalf("error = %v, want the status", err)
	}
	if strings.Contains(err.Error(), "runtime-token") {
		t.Fatalf("error %q leaks the runtime token", err)
	}
}

func TestGCPMetadata(t *testing.T) {
	cleanEnv(t)
	var requests []*http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r)
		if r.Header.Get("Metadata-Flavor") != "Google" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, _ = w.Write([]byte(providerToken(t, providerGoogle, r.URL.Query().Get("audience")) + "\n"))
	}))
	defer server.Close()
	t.Setenv("GCE_METADATA_HOST", strings.TrimPrefix(server.URL, "http://"))
	t.Setenv(credentials.EnvIdentityTokenAudience, "https://auth.dev.onkernel.com")

	token, err := credentials.GCPMetadata("")(context.Background())
	if err != nil {
		t.Fatalf("GCPMetadata: %v", err)
	}
	if strings.TrimSpace(token) != token {
		t.Errorf("token %q was not trimmed", token)
	}
	if got := tokenAudiences(token); len(got) != 1 || got[0] != "https://auth.dev.onkernel.com" {
		t.Errorf("token audience = %v, want the one from KERNEL_IDENTITY_TOKEN_AUDIENCE", got)
	}
	req := requests[0]
	if want := "/computeMetadata/v1/instance/service-accounts/default/identity"; req.URL.Path != want {
		t.Errorf("path = %q, want %q", req.URL.Path, want)
	}
	if want := "audience=https%3A%2F%2Fauth.dev.onkernel.com&format=full"; req.URL.RawQuery != want {
		t.Errorf("query = %q, want %q", req.URL.RawQuery, want)
	}
}

func TestGCPMetadataFailure(t *testing.T) {
	cleanEnv(t)
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	t.Setenv("GCE_METADATA_HOST", strings.TrimPrefix(server.URL, "http://"))
	_, err := credentials.GCPMetadata("aud")(context.Background())
	if err == nil || !strings.Contains(err.Error(), "status 404") {
		t.Fatalf("error = %v, want the status", err)
	}
}

func TestRegisterIdentityTokenSource(t *testing.T) {
	expectPanic := func(name string, register func()) {
		t.Helper()
		defer func() {
			if recover() == nil {
				t.Errorf("%s: no panic", name)
			}
		}()
		register()
	}
	source := func(string) credentials.IdentityTokenFunc { return credentials.Static("x") }
	expectPanic("nil source", func() { credentials.RegisterIdentityTokenSource("test_nil", nil) })
	expectPanic("built-in name", func() { credentials.RegisterIdentityTokenSource(credentials.SourceGitHubActions, source) })
	expectPanic("reserved name", func() { credentials.RegisterIdentityTokenSource(credentials.SourceFile, source) })

	// Unique per run, because registrations last for the whole process.
	name := fmt.Sprintf("test_twice_%d", time.Now().UnixNano())
	credentials.RegisterIdentityTokenSource(name, source)
	expectPanic("registered twice", func() { credentials.RegisterIdentityTokenSource(name, source) })
}
