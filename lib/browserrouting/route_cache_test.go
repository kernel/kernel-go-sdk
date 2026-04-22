package browserrouting

import (
	"net/http"
	"net/url"
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
