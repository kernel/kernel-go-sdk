package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kernel/kernel-go-sdk/option"
)

func TestBrowserRoutingWarmsCacheAndRoutesAllowlistedSubresources(t *testing.T) {
	t.Setenv(browserRoutingSubresourcesEnv, "process")

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

func TestBrowserRoutingRewritesTelemetryStreamToVM(t *testing.T) {
	t.Setenv(browserRoutingSubresourcesEnv, "telemetry/stream")

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

		switch r.URL.Path {
		case "/browsers":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"session_id": "sess-1",
				"base_url":   srv.URL + "/browser/kernel",
				"cdp_ws_url": "wss://browser-session.test/browser/cdp?jwt=token-abc",
			})
		default:
			// Telemetry stream is an SSE response; emit a single frame and close.
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("id: 1\ndata: {\"category\":\"api\"}\n\n"))
		}
	}))
	defer srv.Close()

	client := NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("sk_test"),
		option.WithHTTPClient(srv.Client()),
	)

	if _, err := client.Browsers.New(context.Background(), BrowserNewParams{}); err != nil {
		t.Fatal(err)
	}

	stream := client.Browsers.Telemetry.StreamStreaming(context.Background(), "sess-1", BrowserTelemetryStreamParams{})
	for stream.Next() {
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}

	if route, ok := client.BrowserRouteCache.Load("sess-1"); !ok || route.JWT != "token-abc" {
		t.Fatalf("expected warmed browser route cache, got %#v ok=%v", route, ok)
	}
	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}
	if calls[1].Path != "/browser/kernel/telemetry/stream?jwt=token-abc" {
		t.Fatalf("expected direct VM telemetry stream path, got %q", calls[1].Path)
	}
	if calls[1].Auth != "" {
		t.Fatalf("expected authorization header removed, got %q", calls[1].Auth)
	}
}

func TestBrowserRoutingTelemetryStreamPreservesContextCancellation(t *testing.T) {
	t.Setenv(browserRoutingSubresourcesEnv, "telemetry/stream")

	requestCanceled := make(chan struct{})
	streamPath := make(chan string, 1)
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/browsers" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"session_id": "sess-1",
				"base_url":   srv.URL + "/browser/kernel",
				"cdp_ws_url": "wss://browser-session.test/browser/cdp?jwt=token-abc",
			})
			return
		}

		streamPath <- r.URL.RequestURI()
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
		close(requestCanceled)
	}))
	defer srv.Close()

	client := NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("sk_test"),
		option.WithHTTPClient(srv.Client()),
	)
	if _, err := client.Browsers.New(context.Background(), BrowserNewParams{}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	stream := client.Browsers.Telemetry.StreamStreaming(ctx, "sess-1", BrowserTelemetryStreamParams{})
	if got := <-streamPath; got != "/browser/kernel/telemetry/stream?jwt=token-abc" {
		t.Fatalf("expected direct VM telemetry path, got %q", got)
	}
	next := make(chan bool, 1)
	go func() {
		next <- stream.Next()
	}()

	cancel()
	select {
	case got := <-next:
		if got {
			t.Fatal("expected canceled stream to stop")
		}
		if !errors.Is(stream.Err(), context.Canceled) {
			t.Fatalf("expected context cancellation, got %v", stream.Err())
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for canceled stream to stop")
	}

	select {
	case <-requestCanceled:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for direct VM request cancellation")
	}
}

func TestBrowserRoutingSkipsSubresourcesOutsideConfiguredAllowlist(t *testing.T) {
	t.Setenv(browserRoutingSubresourcesEnv, "computer")

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

func TestBrowserRoutingSubresourcesFromEnvDefaultsToCurl(t *testing.T) {
	original, ok := os.LookupEnv(browserRoutingSubresourcesEnv)
	if err := os.Unsetenv(browserRoutingSubresourcesEnv); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if !ok {
			_ = os.Unsetenv(browserRoutingSubresourcesEnv)
			return
		}
		_ = os.Setenv(browserRoutingSubresourcesEnv, original)
	})
	if got := browserRoutingSubresourcesFromEnv(); len(got) != 2 || got[0] != "curl" || got[1] != "telemetry/stream" {
		t.Fatalf("expected default subresources [curl telemetry/stream], got %#v", got)
	}

	t.Setenv(browserRoutingSubresourcesEnv, "")
	if got := browserRoutingSubresourcesFromEnv(); len(got) != 0 {
		t.Fatalf("expected empty env to disable routing, got %#v", got)
	}

	t.Setenv(browserRoutingSubresourcesEnv, "curl, process")
	got := browserRoutingSubresourcesFromEnv()
	want := []string{"curl", "process"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %#v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %#v", want, got)
		}
	}
}
