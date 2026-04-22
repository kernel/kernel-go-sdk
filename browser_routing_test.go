package kernel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kernel/kernel-go-sdk/option"
)

func TestBrowserRoutingWarmsCacheAndRoutesAllowlistedSubresources(t *testing.T) {
	var calls []struct {
		Path string
		Auth string
	}
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, struct {
			Path string
			Auth string
		}{Path: r.URL.Path + "?" + r.URL.RawQuery, Auth: r.Header.Get("Authorization")})

		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/browsers":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"session_id": "sess-1",
				"base_url":   srv.URL + "/browser/kernel",
				"cdp_ws_url": "wss://browser-session.test/browser/cdp?jwt=token-abc",
			})
		default:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"duration_ms": 1,
				"exit_code":   0,
				"stderr_b64":  "",
				"stdout_b64":  "",
			})
		}
	}))
	defer srv.Close()

	client := NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("sk_test"),
		option.WithHTTPClient(srv.Client()),
		WithBrowserRouting(BrowserRoutingConfig{Enabled: true, DirectToVMSubresources: []string{"process"}}),
	)

	if _, err := client.Browsers.New(context.Background(), BrowserNewParams{}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Browsers.Process.Exec(context.Background(), "sess-1", BrowserProcessExecParams{Command: "echo"}); err != nil {
		t.Fatal(err)
	}

	if route, ok := client.BrowserRouteCache.Load("sess-1"); !ok || route.JWT != "token-abc" {
		t.Fatalf("expected warmed browser route cache, got %#v ok=%v", route, ok)
	}
	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}
	if calls[1].Path != "/browser/kernel/process/exec?jwt=token-abc" {
		t.Fatalf("expected direct VM path, got %q", calls[1].Path)
	}
	if calls[1].Auth != "" {
		t.Fatalf("expected authorization header removed, got %q", calls[1].Auth)
	}
}

func TestBrowserRoutingSkipsSubresourcesOutsideConfiguredAllowlist(t *testing.T) {
	var paths []string
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/browsers":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"session_id": "sess-1",
				"base_url":   srv.URL + "/browser/kernel",
				"cdp_ws_url": "wss://browser-session.test/browser/cdp?jwt=token-abc",
			})
		default:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"duration_ms": 1,
				"exit_code":   0,
				"stderr_b64":  "",
				"stdout_b64":  "",
			})
		}
	}))
	defer srv.Close()

	client := NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("sk_test"),
		option.WithHTTPClient(srv.Client()),
		WithBrowserRouting(BrowserRoutingConfig{Enabled: true, DirectToVMSubresources: []string{"computer"}}),
	)

	if _, err := client.Browsers.New(context.Background(), BrowserNewParams{}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Browsers.Process.Exec(context.Background(), "sess-1", BrowserProcessExecParams{Command: "echo"}); err != nil {
		t.Fatal(err)
	}

	if got := paths[len(paths)-1]; got != "/browsers/sess-1/process/exec" {
		t.Fatalf("expected control-plane path, got %q", got)
	}
}
