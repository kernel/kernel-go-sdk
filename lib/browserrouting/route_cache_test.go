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

// browserGoneResponse builds a routed-VM 404 with the browser_gone body code.
func browserGoneResponse() *http.Response {
	body := `{"code":"browser_gone","message":"browser not found"}`
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// telemetryEventsRequest builds a GET /browsers/{id}/telemetry/events request
// carrying a control-plane Authorization header.
func telemetryEventsRequest(t *testing.T) *http.Request {
	t.Helper()
	reqURL, err := url.Parse("https://api.example/browsers/sess-1/telemetry/events")
	if err != nil {
		t.Fatal(err)
	}
	return &http.Request{
		Method: http.MethodGet,
		URL:    reqURL,
		Host:   "api.example",
		Header: http.Header{"Authorization": []string{"Bearer sk_test"}},
	}
}

func warmTelemetryCache(t *testing.T) *RouteCache {
	t.Helper()
	cache := NewRouteCache()
	cache.Store(Route{
		SessionID: "sess-1",
		BaseURL:   "https://browser.example/browser/kernel",
		JWT:       "jwt-123",
	})
	return cache
}

func TestDirectVMRoutingMiddlewareFallsBackToControlPlaneOnBrowserGone(t *testing.T) {
	cache := warmTelemetryCache(t)
	// "telemetry" is enabled explicitly here; the default subresource list is
	// intentionally not modified by this PR.
	middleware := DirectVMRoutingMiddleware(cache, []string{"telemetry"})
	req := telemetryEventsRequest(t)

	var calls []*http.Request
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		// Snapshot the request state observed on each call.
		clone := *next
		cloneURL := *next.URL
		clone.URL = &cloneURL
		clone.Header = next.Header.Clone()
		calls = append(calls, &clone)
		if len(calls) == 1 {
			return browserGoneResponse(), nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"events":[]}`)),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(calls) != 2 {
		t.Fatalf("expected exactly one control-plane re-issue (2 calls), got %d", len(calls))
	}

	// First call routed to the VM: Authorization stripped, jwt injected.
	if calls[0].URL.Host != "browser.example" {
		t.Fatalf("expected first call routed to VM host, got %q", calls[0].URL.Host)
	}
	if calls[0].Header.Get("Authorization") != "" {
		t.Fatalf("expected VM call to drop Authorization, got %q", calls[0].Header.Get("Authorization"))
	}
	if calls[0].URL.Query().Get("jwt") != "jwt-123" {
		t.Fatalf("expected VM call to carry jwt, got %q", calls[0].URL.Query().Get("jwt"))
	}

	// Second call replays the original control-plane request.
	if calls[1].URL.Host != "api.example" {
		t.Fatalf("expected fallback to original CP host, got %q", calls[1].URL.Host)
	}
	if calls[1].URL.Path != "/browsers/sess-1/telemetry/events" {
		t.Fatalf("expected fallback to original CP path, got %q", calls[1].URL.Path)
	}
	if calls[1].Header.Get("Authorization") != "Bearer sk_test" {
		t.Fatalf("expected fallback to restore Authorization, got %q", calls[1].Header.Get("Authorization"))
	}
	if calls[1].URL.Query().Get("jwt") != "" {
		t.Fatalf("expected fallback to drop jwt query param, got %q", calls[1].URL.Query().Get("jwt"))
	}

	if _, ok := cache.Load("sess-1"); ok {
		t.Fatal("expected fallback to evict the cached route")
	}

	gotBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotBody) != `{"events":[]}` {
		t.Fatalf("expected control-plane response body, got %q", gotBody)
	}
}

func TestDirectVMRoutingMiddlewareReturnsControlPlaneErrorWithoutLooping(t *testing.T) {
	cache := warmTelemetryCache(t)
	middleware := DirectVMRoutingMiddleware(cache, []string{"telemetry"})
	req := telemetryEventsRequest(t)

	callCount := 0
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		callCount++
		if callCount == 1 {
			return browserGoneResponse(), nil
		}
		// Control plane also reports the browser gone; must NOT loop.
		return browserGoneResponse(), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 2 {
		t.Fatalf("expected exactly 2 calls (VM + one CP re-issue), got %d", callCount)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected control-plane response returned as-is, got %d", res.StatusCode)
	}
}

func TestDirectVMRoutingMiddlewareNoFallbackForIneligiblePath(t *testing.T) {
	cache := warmTelemetryCache(t)
	// Route a different subresource that is NOT in the fallback registry.
	middleware := DirectVMRoutingMiddleware(cache, []string{"process"})
	reqURL, err := url.Parse("https://api.example/browsers/sess-1/process/exec")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    reqURL,
		Host:   "api.example",
		Header: http.Header{"Authorization": []string{"Bearer sk_test"}},
	}

	callCount := 0
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		callCount++
		return browserGoneResponse(), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 1 {
		t.Fatalf("expected no fallback for ineligible path, got %d calls", callCount)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected VM 404 returned, got %d", res.StatusCode)
	}
	if _, ok := cache.Load("sess-1"); !ok {
		t.Fatal("expected route to remain cached for ineligible path")
	}
	gotBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(gotBody), "browser_gone") {
		t.Fatalf("expected VM 404 body returned intact, got %q", gotBody)
	}
}

func TestDirectVMRoutingMiddlewareNoFallbackOnTransientServerError(t *testing.T) {
	for _, status := range []int{http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout} {
		cache := warmTelemetryCache(t)
		middleware := DirectVMRoutingMiddleware(cache, []string{"telemetry"})
		req := telemetryEventsRequest(t)

		callCount := 0
		res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
			callCount++
			return &http.Response{
				StatusCode: status,
				Header:     http.Header{},
				Body:       io.NopCloser(strings.NewReader("upstream error")),
			}, nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if callCount != 1 {
			t.Fatalf("status %d: expected transient error to propagate without fallback, got %d calls", status, callCount)
		}
		if res.StatusCode != status {
			t.Fatalf("status %d: expected propagated unchanged, got %d", status, res.StatusCode)
		}
		if _, ok := cache.Load("sess-1"); !ok {
			t.Fatalf("status %d: expected route to remain cached", status)
		}
	}
}

func TestDirectVMRoutingMiddlewareNoFallbackOnConnectionError(t *testing.T) {
	cache := warmTelemetryCache(t)
	middleware := DirectVMRoutingMiddleware(cache, []string{"telemetry"})
	req := telemetryEventsRequest(t)

	callCount := 0
	_, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		callCount++
		return nil, io.ErrUnexpectedEOF
	})
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected connection error to propagate, got %v", err)
	}
	if callCount != 1 {
		t.Fatalf("expected no fallback on connection error, got %d calls", callCount)
	}
	if _, ok := cache.Load("sess-1"); !ok {
		t.Fatal("expected route to remain cached after connection error")
	}
}

func TestDirectVMRoutingMiddlewareNoFallbackOnSuccess(t *testing.T) {
	cache := warmTelemetryCache(t)
	middleware := DirectVMRoutingMiddleware(cache, []string{"telemetry"})
	req := telemetryEventsRequest(t)

	callCount := 0
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		callCount++
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"events":[]}`)),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 1 {
		t.Fatalf("expected no fallback on success, got %d calls", callCount)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 returned, got %d", res.StatusCode)
	}
	if _, ok := cache.Load("sess-1"); !ok {
		t.Fatal("expected route to remain cached on success")
	}
}

func TestDirectVMRoutingMiddlewareNoFallbackOnNon404BrowserGone(t *testing.T) {
	// A live VM's own 404 with a different code must NOT trigger fallback.
	cache := warmTelemetryCache(t)
	middleware := DirectVMRoutingMiddleware(cache, []string{"telemetry"})
	req := telemetryEventsRequest(t)

	callCount := 0
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		callCount++
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"code":"not_found","message":"no such event"}`)),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 1 {
		t.Fatalf("expected no fallback for non-browser_gone 404, got %d calls", callCount)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected VM 404 returned, got %d", res.StatusCode)
	}
	gotBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(gotBody), "not_found") {
		t.Fatalf("expected VM 404 body returned intact, got %q", gotBody)
	}
}

func TestDirectVMRoutingMiddlewareNoFallbackForEligibleSubresourceWrongSuffix(t *testing.T) {
	// The "telemetry" subresource is routed and the body is browser_gone, but the
	// suffix (/foo) is not the registered "/events" — registry keying is exact on
	// {subresource, suffix}, so this must NOT fall back.
	cache := warmTelemetryCache(t)
	middleware := DirectVMRoutingMiddleware(cache, []string{"telemetry"})
	reqURL, err := url.Parse("https://api.example/browsers/sess-1/telemetry/foo")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    reqURL,
		Host:   "api.example",
		Header: http.Header{"Authorization": []string{"Bearer sk_test"}},
	}

	callCount := 0
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		callCount++
		return browserGoneResponse(), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 1 {
		t.Fatalf("expected no fallback for eligible subresource with wrong suffix, got %d calls", callCount)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected VM 404 returned, got %d", res.StatusCode)
	}
	if _, ok := cache.Load("sess-1"); !ok {
		t.Fatal("expected route to remain cached for wrong-suffix path")
	}
}

func TestDirectVMRoutingMiddlewareNoFallbackOnOther4xx(t *testing.T) {
	// A 403 on an eligible GET must NOT fall back — only a 404 with the
	// browser_gone body code is authoritative for fallback.
	cache := warmTelemetryCache(t)
	middleware := DirectVMRoutingMiddleware(cache, []string{"telemetry"})
	req := telemetryEventsRequest(t)

	callCount := 0
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		callCount++
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"code":"forbidden","message":"nope"}`)),
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 1 {
		t.Fatalf("expected no fallback on 403, got %d calls", callCount)
	}
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 returned unchanged, got %d", res.StatusCode)
	}
	if _, ok := cache.Load("sess-1"); !ok {
		t.Fatal("expected route to remain cached on 403")
	}
	gotBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(gotBody), "forbidden") {
		t.Fatalf("expected 403 body returned intact, got %q", gotBody)
	}
}

func TestDirectVMRoutingMiddlewareNoFallbackForNonGetMethod(t *testing.T) {
	cache := warmTelemetryCache(t)
	middleware := DirectVMRoutingMiddleware(cache, []string{"telemetry"})
	reqURL, err := url.Parse("https://api.example/browsers/sess-1/telemetry/events")
	if err != nil {
		t.Fatal(err)
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Host:   "api.example",
		Header: http.Header{"Authorization": []string{"Bearer sk_test"}},
	}

	callCount := 0
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		callCount++
		return browserGoneResponse(), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 1 {
		t.Fatalf("expected no fallback for POST, got %d calls", callCount)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected VM 404 returned for POST, got %d", res.StatusCode)
	}
}

func TestDirectVMRoutingMiddlewareNonRoutedRequestUntouched(t *testing.T) {
	// No cached route -> request is never routed, so fallback never applies even
	// for an eligible path + browser_gone body.
	cache := NewRouteCache()
	middleware := DirectVMRoutingMiddleware(cache, []string{"telemetry"})
	req := telemetryEventsRequest(t)

	var got *http.Request
	callCount := 0
	res, err := middleware(req, func(next *http.Request) (*http.Response, error) {
		callCount++
		got = next
		return browserGoneResponse(), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 1 {
		t.Fatalf("expected non-routed request to make exactly one call, got %d", callCount)
	}
	if got.URL.Host != "api.example" {
		t.Fatalf("expected non-routed request untouched, got host %q", got.URL.Host)
	}
	if got.Header.Get("Authorization") != "Bearer sk_test" {
		t.Fatalf("expected non-routed Authorization untouched, got %q", got.Header.Get("Authorization"))
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected response returned unchanged, got %d", res.StatusCode)
	}
}
