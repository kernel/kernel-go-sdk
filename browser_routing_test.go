package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
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

	handlerDone := make(chan struct{})
	requestCanceled := make(chan struct{})
	streamPath := make(chan string, 1)
	var requestCanceledOnce sync.Once
	var streamPathOnce sync.Once
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

		streamPathOnce.Do(func() { streamPath <- r.URL.RequestURI() })
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		select {
		case <-r.Context().Done():
			requestCanceledOnce.Do(func() { close(requestCanceled) })
		case <-handlerDone:
		}
	}))
	t.Cleanup(func() {
		close(handlerDone)
		srv.Close()
	})

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

func TestBrowserRoutingSubresourcesFromEnvDefaults(t *testing.T) {
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
	if got := browserRoutingSubresourcesFromEnv(); len(got) != 5 || got[0] != "curl" || got[1] != "telemetry/stream" || got[2] != "computer" || got[3] != "playwright" || got[4] != "process" {
		t.Fatalf("expected default subresources [curl telemetry/stream computer playwright process], got %#v", got)
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

func TestBrowserRoutingDefaultsRouteToVM(t *testing.T) {
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
		case "/browser/kernel/computer/screenshot":
			w.Header().Set("Content-Type", "image/png")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte{0x89, 0x50, 0x4e, 0x47})
		case "/browser/kernel/process/exec":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"duration_ms": 1,
				"exit_code":   0,
				"stderr_b64":  "",
				"stdout_b64":  "",
			})
		default:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
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
	screenshot, err := client.Browsers.Computer.CaptureScreenshot(context.Background(), "sess-1", BrowserComputerCaptureScreenshotParams{})
	if err != nil {
		t.Fatal(err)
	}
	if screenshot != nil {
		_ = screenshot.Body.Close()
	}
	if _, err := client.Browsers.Playwright.Execute(context.Background(), "sess-1", BrowserPlaywrightExecuteParams{Code: "return 1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Browsers.Process.Exec(context.Background(), "sess-1", BrowserProcessExecParams{Command: "echo"}); err != nil {
		t.Fatal(err)
	}

	if len(calls) != 4 {
		t.Fatalf("expected 4 calls, got %d %#v", len(calls), calls)
	}
	if calls[1].Path != "/browser/kernel/computer/screenshot?jwt=token-abc" {
		t.Fatalf("expected direct VM screenshot path, got %q", calls[1].Path)
	}
	if calls[1].Auth != "" {
		t.Fatalf("expected authorization header removed, got %q", calls[1].Auth)
	}
	if calls[2].Path != "/browser/kernel/playwright/execute?jwt=token-abc" {
		t.Fatalf("expected direct VM playwright path, got %q", calls[2].Path)
	}
	if calls[2].Auth != "" {
		t.Fatalf("expected authorization header removed, got %q", calls[2].Auth)
	}
	if calls[3].Path != "/browser/kernel/process/exec?jwt=token-abc" {
		t.Fatalf("expected direct VM process path, got %q", calls[3].Path)
	}
	if calls[3].Auth != "" {
		t.Fatalf("expected authorization header removed, got %q", calls[3].Auth)
	}
}

func TestBrowserRoutingDefaultsKeepFsAndTelemetryEventsOnAPI(t *testing.T) {
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

	var paths []string
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		switch {
		case r.URL.Path == "/browsers":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"session_id": "sess-1",
				"base_url":   srv.URL + "/browser/kernel",
				"cdp_ws_url": "wss://browser-session.test/browser/cdp?jwt=token-abc",
			})
		case r.URL.Path == "/browsers/sess-1/fs/read_file":
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("x"))
		case r.URL.Path == "/browsers/sess-1/telemetry/events":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]any{})
		default:
			w.Header().Set("Content-Type", "application/json")
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
	fsResp, err := client.Browsers.Fs.ReadFile(context.Background(), "sess-1", BrowserFReadFileParams{Path: "/tmp/x"})
	if err != nil {
		t.Fatal(err)
	}
	if fsResp != nil {
		_ = fsResp.Body.Close()
	}
	if _, err := client.Browsers.Telemetry.Events(context.Background(), "sess-1", BrowserTelemetryEventsParams{}); err != nil {
		t.Fatal(err)
	}

	want := []string{
		"/browsers",
		"/browsers/sess-1/fs/read_file",
		"/browsers/sess-1/telemetry/events",
	}
	if len(paths) != len(want) {
		t.Fatalf("expected paths %v, got %v", want, paths)
	}
	for i := range want {
		if paths[i] != want[i] {
			t.Fatalf("expected paths %v, got %v", want, paths)
		}
	}
}
