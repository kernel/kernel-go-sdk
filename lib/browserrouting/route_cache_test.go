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

func TestParseBrowserMetadataPathRejectsSubresourcePaths(t *testing.T) {
	if sessionID, ok := parseBrowserMetadataPath("/browsers/sess-1/process/exec"); ok || sessionID != "" {
		t.Fatalf("expected subresource path to be rejected, got sessionID=%q ok=%v", sessionID, ok)
	}
	if sessionID, ok := parseBrowserMetadataPath("/browsers/sess-1/browsers"); ok || sessionID != "" {
		t.Fatalf("expected nested browsers subresource path to be rejected, got sessionID=%q ok=%v", sessionID, ok)
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
