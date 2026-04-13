package kernel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kernel/kernel-go-sdk/lib/browserscope"
	"github.com/kernel/kernel-go-sdk/option"
)

func TestForBrowserRewritesToKernelPaths(t *testing.T) {
	var gotPath string
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		if r.URL.Query().Get("jwt") == "" {
			http.Error(w, "jwt", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"duration_ms": 1,
			"exit_code":   0,
			"stderr_b64":  "",
			"stdout_b64":  "",
		})
	}))
	defer srv.Close()

	c := NewClient(
		option.WithBaseURL("https://api.example/"),
		option.WithAPIKey("sk_test"),
		option.WithHTTPClient(srv.Client()),
	)

	b := &BrowserGetResponse{
		SessionID: "sid-1",
		BaseURL:   srv.URL + "/browser/kernel",
		CdpWsURL:  "wss://x/browser/cdp?jwt=session-jwt",
	}

	sess, err := c.ForBrowser(b)
	if err != nil {
		t.Fatal(err)
	}
	if sess.SessionID() != "sid-1" {
		t.Fatalf("session id: %s", sess.SessionID())
	}

	_, err = sess.Process.Exec(context.Background(), BrowserProcessExecParams{Command: "true"})
	if err != nil {
		t.Fatal(err)
	}

	if gotPath != "/browser/kernel/process/exec" {
		t.Fatalf("path: got %q", gotPath)
	}
	if gotAuth != "" {
		t.Fatalf("authorization should be empty, got %q", gotAuth)
	}
}

func TestForBrowserRefNormalize(t *testing.T) {
	ref := browserscope.Ref{
		SessionID: "s",
		BaseURL:   "https://x/browser/kernel",
		JWT:       "direct",
	}
	norm, err := ref.Normalize()
	if err != nil {
		t.Fatal(err)
	}
	if norm.JWT != "direct" {
		t.Fatalf("jwt: %q", norm.JWT)
	}
}
