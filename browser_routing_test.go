package kernel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
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
	want := []string{"curl", "telemetry/stream", "computer", "playwright", "process", "fs", "logs/stream"}
	if got := browserRoutingSubresourcesFromEnv(); !slices.Equal(got, want) {
		t.Fatalf("expected default subresources %v, got %#v", want, got)
	}

	t.Setenv(browserRoutingSubresourcesEnv, "")
	if got := browserRoutingSubresourcesFromEnv(); len(got) != 0 {
		t.Fatalf("expected empty env to disable routing, got %#v", got)
	}

	t.Setenv(browserRoutingSubresourcesEnv, "curl, process")
	wantConfigured := []string{"curl", "process"}
	if got := browserRoutingSubresourcesFromEnv(); !slices.Equal(got, wantConfigured) {
		t.Fatalf("expected %v, got %#v", wantConfigured, got)
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

func TestBrowserRoutingDefaultsKeepControlPlaneSubresourcesOnAPI(t *testing.T) {
	unsetBrowserRoutingEnv(t)

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
			_ = json.NewEncoder(w).Encode([]any{})
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
	if _, err := client.Browsers.Telemetry.Events(context.Background(), "sess-1", BrowserTelemetryEventsParams{}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Browsers.Replays.List(context.Background(), "sess-1"); err != nil {
		t.Fatal(err)
	}

	want := []string{
		"/browsers",
		"/browsers/sess-1/telemetry/events",
		"/browsers/sess-1/replays",
	}
	if !slices.Equal(paths, want) {
		t.Fatalf("expected paths %v, got %v", want, paths)
	}
}

func TestBrowserRoutingDefaultsRouteFsAndLogsToVM(t *testing.T) {
	unsetBrowserRoutingEnv(t)

	type call struct {
		Path string
		Auth string
		Body string
	}
	var calls []call
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		calls = append(calls, call{Path: r.URL.RequestURI(), Auth: r.Header.Get("Authorization"), Body: string(body)})

		switch r.URL.Path {
		case "/browsers":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"session_id": "sess-1",
				"base_url":   srv.URL + "/browser/kernel",
				"cdp_ws_url": "wss://browser-session.test/browser/cdp?jwt=token-abc",
			})
		case "/browser/kernel/fs/read_file":
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte{0x00, 0x01})
		case "/browser/kernel/fs/list_files":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]any{})
		case "/browser/kernel/logs/stream", "/browser/kernel/fs/watch/watch-1/events":
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("data: {\"event\":\"log\",\"message\":\"hello\",\"timestamp\":\"2020-01-01T00:00:00Z\"}\n\n"))
		default:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{})
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
	if _, err := client.Browsers.Fs.ListFiles(context.Background(), "sess-1", BrowserFListFilesParams{Path: "/tmp"}); err != nil {
		t.Fatal(err)
	}
	readRes, err := client.Browsers.Fs.ReadFile(context.Background(), "sess-1", BrowserFReadFileParams{Path: "/tmp/x"})
	if err != nil {
		t.Fatal(err)
	}
	contents, err := io.ReadAll(readRes.Body)
	if err != nil {
		t.Fatal(err)
	}
	_ = readRes.Body.Close()
	if !bytes.Equal(contents, []byte{0x00, 0x01}) {
		t.Fatalf("expected binary body from VM, got %v", contents)
	}
	if err := client.Browsers.Fs.WriteFile(context.Background(), "sess-1", strings.NewReader("payload"), BrowserFWriteFileParams{
		Path: "/tmp/x",
		Mode: String("600"),
	}); err != nil {
		t.Fatal(err)
	}
	if err := client.Browsers.Fs.Upload(context.Background(), "sess-1", BrowserFUploadParams{
		Files: []BrowserFUploadParamsFile{
			{DestPath: "/tmp/one", File: strings.NewReader("one")},
			{DestPath: "/tmp/two", File: strings.NewReader("two")},
		},
	}); err != nil {
		t.Fatal(err)
	}
	watchStream := client.Browsers.Fs.Watch.EventsStreaming(context.Background(), "watch-1", BrowserFWatchEventsParams{IDOrName: "sess-1"})
	for watchStream.Next() {
	}
	if err := watchStream.Err(); err != nil {
		t.Fatal(err)
	}
	logStream := client.Browsers.Logs.StreamStreaming(context.Background(), "sess-1", BrowserLogStreamParams{
		Source: BrowserLogStreamParamsSourcePath,
		Path:   String("/var/log/x"),
		Follow: Bool(true),
	})
	for logStream.Next() {
	}
	if err := logStream.Err(); err != nil {
		t.Fatal(err)
	}

	wantPaths := []string{
		"/browsers",
		"/browser/kernel/fs/list_files?jwt=token-abc&path=%2Ftmp",
		"/browser/kernel/fs/read_file?jwt=token-abc&path=%2Ftmp%2Fx",
		"/browser/kernel/fs/write_file?jwt=token-abc&mode=600&path=%2Ftmp%2Fx",
		"/browser/kernel/fs/upload?jwt=token-abc",
		"/browser/kernel/fs/watch/watch-1/events?jwt=token-abc",
		"/browser/kernel/logs/stream?follow=true&jwt=token-abc&path=%2Fvar%2Flog%2Fx&source=path",
	}
	gotPaths := make([]string, 0, len(calls))
	for _, c := range calls {
		gotPaths = append(gotPaths, c.Path)
	}
	if !slices.Equal(gotPaths, wantPaths) {
		t.Fatalf("expected paths %v, got %v", wantPaths, gotPaths)
	}
	for _, c := range calls[1:] {
		if c.Auth != "" {
			t.Fatalf("expected authorization header removed for %q, got %q", c.Path, c.Auth)
		}
	}
	if calls[3].Body != "payload" {
		t.Fatalf("expected write_file body forwarded to the VM, got %q", calls[3].Body)
	}
	uploadBody := calls[4].Body
	// The VM accepts dot-indexed multipart array names (files.<index>.<field>).
	for _, want := range []string{`name="files.0.dest_path"`, `name="files.0.file"`, `name="files.1.dest_path"`} {
		if !strings.Contains(uploadBody, want) {
			t.Fatalf("expected upload body to contain %s, got %q", want, uploadBody)
		}
	}
}

func TestBrowserRoutingLogsStreamPreservesContextCancellation(t *testing.T) {
	unsetBrowserRoutingEnv(t)

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
	stream := client.Browsers.Logs.StreamStreaming(ctx, "sess-1", BrowserLogStreamParams{
		Source:            BrowserLogStreamParamsSourceSupervisor,
		SupervisorProcess: String("chromium"),
	})
	if got := <-streamPath; got != "/browser/kernel/logs/stream?jwt=token-abc&source=supervisor&supervisor_process=chromium" {
		t.Fatalf("expected direct VM logs path, got %q", got)
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

func TestBrowserRoutingFsFallsBackToControlPlaneOnStaleJWT(t *testing.T) {
	unsetBrowserRoutingEnv(t)

	type call struct {
		Path string
		Auth string
		Body string
	}
	var calls []call
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		calls = append(calls, call{Path: r.URL.RequestURI(), Auth: r.Header.Get("Authorization"), Body: string(body)})

		switch {
		case r.URL.Path == "/browsers":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"session_id": "sess-1",
				"base_url":   srv.URL + "/browser/kernel",
				"cdp_ws_url": "wss://browser-session.test/browser/cdp?jwt=token-abc",
			})
		case strings.HasPrefix(r.URL.Path, "/browser/kernel/"):
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Invalid JWT"))
		default:
			w.WriteHeader(http.StatusNoContent)
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
	if err := client.Browsers.Fs.Upload(context.Background(), "sess-1", BrowserFUploadParams{
		Files: []BrowserFUploadParamsFile{{DestPath: "/tmp/one", File: strings.NewReader("one")}},
	}); err != nil {
		t.Fatal(err)
	}

	if len(calls) != 3 {
		t.Fatalf("expected browser create plus VM and control-plane upload, got %#v", calls)
	}
	if calls[1].Path != "/browser/kernel/fs/upload?jwt=token-abc" {
		t.Fatalf("expected direct VM upload first, got %q", calls[1].Path)
	}
	if calls[2].Path != "/browsers/sess-1/fs/upload" {
		t.Fatalf("expected control-plane fallback, got %q", calls[2].Path)
	}
	if calls[2].Auth != "Bearer sk_test" {
		t.Fatalf("expected API key on the fallback request, got %q", calls[2].Auth)
	}
	if calls[2].Body != calls[1].Body {
		t.Fatalf("expected the fallback to replay the upload body")
	}
	if _, ok := client.BrowserRouteCache.Load("sess-1"); ok {
		t.Fatal("expected the stale route to be evicted")
	}
}

func TestBrowserRoutingEnvOverrideKeepsFsAndLogsOnControlPlane(t *testing.T) {
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
			_ = json.NewEncoder(w).Encode([]any{})
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
	if _, err := client.Browsers.Fs.ListFiles(context.Background(), "sess-1", BrowserFListFilesParams{Path: "/tmp"}); err != nil {
		t.Fatal(err)
	}

	want := []string{"/browsers", "/browsers/sess-1/fs/list_files"}
	if !slices.Equal(paths, want) {
		t.Fatalf("expected paths %v, got %v", want, paths)
	}
}

func unsetBrowserRoutingEnv(t *testing.T) {
	t.Helper()
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
}
