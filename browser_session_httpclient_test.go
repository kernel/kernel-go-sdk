package kernel

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/lib/browserrouting"
	"github.com/kernel/kernel-go-sdk/option"
)

func TestBrowserSessionHTTPClientRawCurl(t *testing.T) {
	var sawRaw string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/browser/kernel/curl/raw" {
			http.NotFound(w, r)
			return
		}
		sawRaw = r.URL.RawQuery
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("proxied"))
	}))
	defer srv.Close()

	c := NewClient(
		option.WithBaseURL("https://api.example/"),
		option.WithAPIKey("sk"),
		option.WithHTTPClient(srv.Client()),
	)

	storeBrowserRouteCache(c.Options, browserrouting.Ref{
		SessionID: "sid",
		BaseURL:   srv.URL + "/browser/kernel",
		CdpWsURL:  "wss://x/browser/cdp?jwt=j1",
	})

	hc, err := c.Browsers.HTTPClient("sid")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodGet, "https://httpbin.org/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := hc.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if string(body) != "proxied" {
		t.Fatalf("body %q", body)
	}
	if sawRaw == "" {
		t.Fatal("expected raw query on curl/raw")
	}
}

func TestBrowserSessionHTTPClientRequiresCachedRoute(t *testing.T) {
	c := NewClient(
		option.WithBaseURL("https://api.example/"),
		option.WithAPIKey("sk"),
	)

	storeBrowserRouteCache(c.Options, browserrouting.Ref{
		SessionID: "sid",
		BaseURL:   "https://browser-session.test/browser/kernel",
		CdpWsURL:  "wss://x/browser/cdp?jwt=j1",
	})
	c.BrowserRouteCache.Delete("sid")

	_, err := c.Browsers.HTTPClient("sid")
	if err == nil {
		t.Fatal("expected cached route lookup failure")
	}
}

func TestBrowserHTTPClientPropagatesRequestConfigError(t *testing.T) {
	c := NewClient(
		option.WithBaseURL("https://api.example/"),
		option.WithAPIKey("sk"),
	)

	storeBrowserRouteCache(c.Options, browserrouting.Ref{
		SessionID: "sid",
		BaseURL:   "https://browser-session.test/browser/kernel",
		CdpWsURL:  "wss://x/browser/cdp?jwt=j1",
	})

	wantErr := errors.New("request config failed")
	hc, err := c.Browsers.HTTPClient("sid", requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error {
		return wantErr
	}))
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %q, got %v", wantErr, err)
	}
	if hc != nil {
		t.Fatal("expected nil client when request config fails")
	}
}
