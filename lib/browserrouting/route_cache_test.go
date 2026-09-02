package browserrouting

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestDirectVMRoutingMiddlewareClearsStaleRawPath(t *testing.T) {
	cache := NewRouteCache()
	cache.Store(Route{
		SessionID: "sess-1",
		BaseURL:   "https://browser.example/browser/kernel",
		JWT:       "jwt-123",
	})

	middleware := DirectVMRoutingMiddleware(cache, []string{"process"})
	reqURL, err := url.Parse("https://api.example/browsers/sess-1/process/exec")
	if err != nil {
		t.Fatal(err)
	}
	reqURL.RawPath = "/browsers/sess-1/process/%65xec"

	req := &http.Request{
		URL:    reqURL,
		Header: http.Header{"Authorization": []string{"Bearer sk_test"}},
	}

	var got *http.Request
	_, err = middleware(req, func(next *http.Request) (*http.Response, error) {
		got = next
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("expected middleware to invoke next handler")
	}
	if got.URL.Path != "/browser/kernel/process/exec" {
		t.Fatalf("expected rewritten path, got %q", got.URL.Path)
	}
	if got.URL.RawPath != "" {
		t.Fatalf("expected stale raw path to be cleared, got %q", got.URL.RawPath)
	}
	if got.URL.Query().Get("jwt") != "jwt-123" {
		t.Fatalf("expected jwt query param, got %q", got.URL.Query().Get("jwt"))
	}
	if got.Header.Get("Authorization") != "" {
		t.Fatalf("expected authorization to be stripped, got %q", got.Header.Get("Authorization"))
	}
}

func TestDirectVMRoutingMiddlewareAllowlistMatching(t *testing.T) {
	// Pins the segment-boundary allowlist: telemetry/stream (live SSE) routes to
	// the VM, telemetry/events (historical, served by the control plane from S2)
	// does NOT, and a stream-prefixed-but-different path is not matched.
	cases := []struct {
		name       string
		path       string
		routedToVM bool
	}{
		{"telemetry stream -> VM", "/browsers/sess-1/telemetry/stream", true},
		{"telemetry events -> control plane", "/browsers/sess-1/telemetry/events", false},
		{"telemetry streaming-config not a stream prefix", "/browsers/sess-1/telemetry/streaming-config", false},
		{"bare telemetry -> control plane", "/browsers/sess-1/telemetry", false},
		{"curl proxy -> VM", "/browsers/sess-1/curl/raw", true},
		{"computer screenshot -> VM", "/browsers/sess-1/computer/screenshot", true},
		{"playwright execute -> VM", "/browsers/sess-1/playwright/execute", true},
		{"process exec -> VM", "/browsers/sess-1/process/exec", true},
		{"process stdout stream -> VM", "/browsers/sess-1/process/proc-1/stdout/stream", true},
		{"fs read_file -> VM", "/browsers/sess-1/fs/read_file", true},
		{"fs watch events -> VM", "/browsers/sess-1/fs/watch/watch-1/events", true},
		{"fs-prefixed segment not matched", "/browsers/sess-1/fsx/read_file", false},
		{"logs stream -> VM", "/browsers/sess-1/logs/stream", true},
		{"logs-prefixed segment not matched", "/browsers/sess-1/logstream", false},
		{"extensions -> control plane", "/browsers/sess-1/extensions", false},
		{"replays -> control plane", "/browsers/sess-1/replays/rec-1", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cache := NewRouteCache()
			cache.Store(Route{
				SessionID: "sess-1",
				BaseURL:   "https://browser.example/browser/kernel",
				JWT:       "jwt-123",
			})
			middleware := DirectVMRoutingMiddleware(cache, []string{"curl", "telemetry/stream", "computer", "playwright", "process", "fs", "logs"})

			reqURL, err := url.Parse("https://api.example" + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			req := &http.Request{URL: reqURL, Header: http.Header{}}

			var got *http.Request
			if _, err := middleware(req, func(next *http.Request) (*http.Response, error) {
				got = next
				return nil, nil
			}); err != nil {
				t.Fatal(err)
			}
			if got == nil {
				t.Fatal("expected middleware to invoke next handler")
			}
			routedToVM := got.URL.Host == "browser.example"
			if routedToVM != tc.routedToVM {
				t.Fatalf("path %q routedToVM=%v, want %v (host=%q path=%q)", tc.path, routedToVM, tc.routedToVM, got.URL.Host, got.URL.Path)
			}
		})
	}
}

func TestParseCacheLifecycleRejectsBrowserSubresourcePaths(t *testing.T) {
	cases := []string{
		"/browsers/sess-1/process/exec",
		"/browsers/sess-1/browsers",
	}

	for _, path := range cases {
		reqURL, err := url.Parse("https://api.example" + path)
		if err != nil {
			t.Fatal(err)
		}
		lifecycle, err := parseCacheLifecycle(&http.Request{
			Method: http.MethodGet,
			URL:    reqURL,
			Header: http.Header{},
		})
		if err != nil {
			t.Fatal(err)
		}
		if lifecycle.sniffResponse || lifecycle.evictSessionID != "" {
			t.Fatalf("expected subresource path %q to be ignored, got %#v", path, lifecycle)
		}
	}
}

func TestDirectVMRoutingMiddlewarePopulatesCacheFromJSONResponse(t *testing.T) {
	cache := NewRouteCache()
	middleware := DirectVMRoutingMiddleware(cache, nil)
	body := `{"session_id":"sess-1","base_url":"https://browser.example/browser/kernel","cdp_ws_url":"wss://browser.example/browser/cdp?jwt=jwt-123"}`

	reqURL, err := url.Parse("https://api.example/browsers")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: http.Header{},
	}

	res, err := middleware(req, func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	route, ok := cache.Load("sess-1")
	if !ok {
		t.Fatal("expected browser route cache to be warmed")
	}
	if route.BaseURL != "https://browser.example/browser/kernel" || route.JWT != "jwt-123" {
		t.Fatalf("unexpected cached route %#v", route)
	}

	gotBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotBody) != body {
		t.Fatalf("expected response body to be preserved, got %q", gotBody)
	}
}

func TestDirectVMRoutingMiddlewarePopulatesCacheFromNestedJSONResponse(t *testing.T) {
	cache := NewRouteCache()
	middleware := DirectVMRoutingMiddleware(cache, nil)
	body := `{"items":[{"session_id":"sess-2","base_url":"https://browser.example/browser/kernel","cdp_ws_url":"wss://browser.example/browser/cdp?jwt=jwt-234"}]}`

	reqURL, err := url.Parse("https://api.example/v1/browsers")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    reqURL,
		Header: http.Header{},
	}

	_, err = middleware(req, func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	route, ok := cache.Load("sess-2")
	if !ok {
		t.Fatal("expected nested browser metadata to warm the cache")
	}
	if route.JWT != "jwt-234" {
		t.Fatalf("expected nested browser metadata jwt to be cached, got %q", route.JWT)
	}
}

func TestDirectVMRoutingMiddlewarePopulatesCacheFromBrowserPoolAcquireResponse(t *testing.T) {
	cache := NewRouteCache()
	middleware := DirectVMRoutingMiddleware(cache, nil)
	body := `{"session_id":"sess-3","base_url":"https://browser.example/browser/kernel","cdp_ws_url":"wss://browser.example/browser/cdp?jwt=jwt-345"}`

	reqURL, err := url.Parse("https://api.example/browser_pools/pool-1/acquire")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: http.Header{},
	}

	_, err = middleware(req, func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	route, ok := cache.Load("sess-3")
	if !ok {
		t.Fatal("expected browser pool acquire response to warm the cache")
	}
	if route.JWT != "jwt-345" {
		t.Fatalf("expected browser pool acquire jwt to be cached, got %q", route.JWT)
	}
}

func TestDirectVMRoutingMiddlewareSkipsCacheSniffingForNonBrowserMetadataPaths(t *testing.T) {
	cache := NewRouteCache()
	middleware := DirectVMRoutingMiddleware(cache, nil)
	body := `{"session_id":"sess-1","base_url":"https://browser.example/browser/kernel","cdp_ws_url":"wss://browser.example/browser/cdp?jwt=jwt-123"}`

	reqURL, err := url.Parse("https://api.example/projects")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    reqURL,
		Header: http.Header{},
	}

	_, err = middleware(req, func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := cache.Load("sess-1"); ok {
		t.Fatal("expected non-browser metadata response to skip cache warm-up")
	}
}

func TestDirectVMRoutingMiddlewareEvictsCacheOnSuccessfulBrowserDelete(t *testing.T) {
	cache := NewRouteCache()
	cache.Store(Route{
		SessionID: "sess-1",
		BaseURL:   "https://browser.example/browser/kernel",
		JWT:       "jwt-123",
	})
	middleware := DirectVMRoutingMiddleware(cache, nil)

	reqURL, err := url.Parse("https://api.example/browsers/sess-1")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodDelete,
		URL:    reqURL,
		Header: http.Header{},
	}

	_, err = middleware(req, func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := cache.Load("sess-1"); ok {
		t.Fatal("expected successful browser delete to evict cached route")
	}
}

func TestDirectVMRoutingMiddlewareEvictsCacheOnSuccessfulBrowserPoolRelease(t *testing.T) {
	cache := NewRouteCache()
	cache.Store(Route{
		SessionID: "sess-1",
		BaseURL:   "https://browser.example/browser/kernel",
		JWT:       "jwt-123",
	})
	middleware := DirectVMRoutingMiddleware(cache, nil)
	releaseBody := `{"session_id":"sess-1","reuse":false}`

	reqURL, err := url.Parse("https://api.example/browser_pools/pool-1/release")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method:        http.MethodPost,
		URL:           reqURL,
		Header:        http.Header{"Content-Type": []string{"application/json"}},
		Body:          io.NopCloser(strings.NewReader(releaseBody)),
		ContentLength: int64(len(releaseBody)),
	}

	var gotRequestBody string
	_, err = middleware(req, func(next *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(next.Body)
		if err != nil {
			return nil, err
		}
		gotRequestBody = string(body)
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if gotRequestBody != releaseBody {
		t.Fatalf("expected release request body to be preserved, got %q", gotRequestBody)
	}
	if _, ok := cache.Load("sess-1"); ok {
		t.Fatal("expected successful browser pool release to evict cached route")
	}
}

func TestDirectVMRoutingMiddlewareDeleteWinsOverJSONCacheSniff(t *testing.T) {
	cache := NewRouteCache()
	cache.Store(Route{
		SessionID: "sess-1",
		BaseURL:   "https://browser.example/browser/kernel",
		JWT:       "jwt-123",
	})
	middleware := DirectVMRoutingMiddleware(cache, nil)

	reqURL, err := url.Parse("https://api.example/browsers/sess-1")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodDelete,
		URL:    reqURL,
		Header: http.Header{},
	}

	_, err = middleware(req, func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(strings.NewReader(
				`{"session_id":"sess-1","base_url":"https://browser.example/browser/kernel","cdp_ws_url":"wss://browser.example/browser/cdp?jwt=jwt-123"}`,
			)),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := cache.Load("sess-1"); ok {
		t.Fatal("expected delete response to leave cached route evicted")
	}
}

func TestDirectVMRoutingMiddlewareFallsBackOnStaleJWT(t *testing.T) {
	cache := NewRouteCache()
	cache.Store(Route{
		SessionID: "sess-1",
		BaseURL:   "https://browser.example/browser/kernel",
		JWT:       "jwt-123",
	})

	middleware := DirectVMRoutingMiddleware(cache, []string{"computer"})
	reqURL, err := url.Parse("https://api.example/browsers/sess-1/computer/screenshot")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: http.Header{"Authorization": []string{"Bearer sk_test"}},
		Host:   "api.example",
	}

	var calls []string
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		calls = append(calls, next.URL.String())
		if next.URL.Host == "browser.example" {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader("Invalid JWT")),
			}, nil
		}
		if next.Header.Get("Authorization") != "Bearer sk_test" {
			t.Fatalf("expected restored authorization, got %q", next.Header.Get("Authorization"))
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("png")),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 after fallback, got %d", res.StatusCode)
	}
	if len(calls) != 2 {
		t.Fatalf("expected vm then control-plane call, got %v", calls)
	}
	if !strings.Contains(calls[0], "browser.example") || !strings.Contains(calls[0], "jwt=jwt-123") {
		t.Fatalf("expected first call on VM with jwt, got %q", calls[0])
	}
	if !strings.Contains(calls[1], "api.example/browsers/sess-1/computer/screenshot") {
		t.Fatalf("expected second call on control plane, got %q", calls[1])
	}
	if _, ok := cache.Load("sess-1"); ok {
		t.Fatal("expected stale jwt to evict cached route")
	}
}

func TestDirectVMRoutingMiddlewareKeepsRefreshedRouteAfterStaleJWT(t *testing.T) {
	cache := NewRouteCache()
	cache.Store(Route{
		SessionID: "sess-1",
		BaseURL:   "https://browser.example/browser/kernel",
		JWT:       "jwt-123",
	})

	middleware := DirectVMRoutingMiddleware(cache, []string{"computer"})
	reqURL, err := url.Parse("https://api.example/browsers/sess-1/computer/screenshot")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: http.Header{"Authorization": []string{"Bearer sk_test"}},
		Host:   "api.example",
	}

	_, err = middleware(req, func(next *http.Request) (*http.Response, error) {
		if next.URL.Host == "browser.example" {
			cache.Store(Route{
				SessionID: "sess-1",
				BaseURL:   "https://browser.example/browser/kernel",
				JWT:       "jwt-FRESH",
			})
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader("Invalid JWT")),
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("png")),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	route, ok := cache.Load("sess-1")
	if !ok {
		t.Fatal("expected refreshed route to survive stale jwt fallback")
	}
	if route.JWT != "jwt-FRESH" {
		t.Fatalf("expected jwt-FRESH, got %q", route.JWT)
	}
}

func TestDirectVMRoutingMiddlewareRewindsBodyOnStaleJWTFallback(t *testing.T) {
	cache := NewRouteCache()
	cache.Store(Route{
		SessionID: "sess-1",
		BaseURL:   "https://browser.example/browser/kernel",
		JWT:       "jwt-123",
	})

	body := []byte(`{"code":"return 1"}`)
	middleware := DirectVMRoutingMiddleware(cache, []string{"playwright"})
	reqURL, err := url.Parse("https://api.example/browsers/sess-1/playwright/execute")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: http.Header{"Authorization": []string{"Bearer sk_test"}},
		Host:   "api.example",
		Body:   io.NopCloser(strings.NewReader(string(body))),
		GetBody: func() (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(string(body))), nil
		},
		ContentLength: int64(len(body)),
	}

	var gotBodies []string
	_, err = middleware(req, func(next *http.Request) (*http.Response, error) {
		b, readErr := io.ReadAll(next.Body)
		if readErr != nil {
			return nil, readErr
		}
		gotBodies = append(gotBodies, string(b))
		if next.URL.Host == "browser.example" {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader("Invalid JWT")),
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"success":true}`)),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(gotBodies) != 2 {
		t.Fatalf("expected two bodies, got %v", gotBodies)
	}
	if gotBodies[0] != string(body) || gotBodies[1] != string(body) {
		t.Fatalf("expected rewound body on fallback, got %v", gotBodies)
	}
}

func TestDirectVMRoutingMiddlewareKeepsAuthResponseWhenBodyCannotRewind(t *testing.T) {
	cache := NewRouteCache()
	cache.Store(Route{
		SessionID: "sess-1",
		BaseURL:   "https://browser.example/browser/kernel",
		JWT:       "jwt-123",
	})

	middleware := DirectVMRoutingMiddleware(cache, []string{"playwright"})
	reqURL, err := url.Parse("https://api.example/browsers/sess-1/playwright/execute")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: http.Header{"Authorization": []string{"Bearer sk_test"}},
		Host:   "api.example",
		Body:   io.NopCloser(strings.NewReader(`{"code":"return 1"}`)),
	}

	var calls int
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		calls++
		_, _ = io.ReadAll(next.Body)
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(strings.NewReader("Invalid JWT")),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("expected no control-plane retry without GetBody, got %d calls", calls)
	}
	if _, ok := cache.Load("sess-1"); ok {
		t.Fatal("expected the stale route to be evicted even without a retry")
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected original 401, got %d", res.StatusCode)
	}
	got, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected readable 401 body, got %v", err)
	}
	if string(got) != "Invalid JWT" {
		t.Fatalf("expected Invalid JWT, got %q", got)
	}
}

func TestDirectVMRoutingMiddlewareKeepsAuthResponseWhenGetBodyFails(t *testing.T) {
	cache := NewRouteCache()
	cache.Store(Route{
		SessionID: "sess-1",
		BaseURL:   "https://browser.example/browser/kernel",
		JWT:       "jwt-123",
	})

	middleware := DirectVMRoutingMiddleware(cache, []string{"playwright"})
	reqURL, err := url.Parse("https://api.example/browsers/sess-1/playwright/execute")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: http.Header{"Authorization": []string{"Bearer sk_test"}},
		Host:   "api.example",
		Body:   io.NopCloser(strings.NewReader(`{"code":"return 1"}`)),
		GetBody: func() (io.ReadCloser, error) {
			return nil, io.ErrUnexpectedEOF
		},
	}

	var calls int
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		calls++
		_, _ = io.ReadAll(next.Body)
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(strings.NewReader("Invalid JWT")),
		}, nil
	})
	if err != nil {
		t.Fatalf("expected original auth response, got err %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected no control-plane retry when GetBody fails, got %d calls", calls)
	}
	if _, ok := cache.Load("sess-1"); ok {
		t.Fatal("expected the stale route to be evicted even without a retry")
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected original 401, got %d", res.StatusCode)
	}
	got, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected readable 401 body, got %v", err)
	}
	if string(got) != "Invalid JWT" {
		t.Fatalf("expected Invalid JWT, got %q", got)
	}
}
