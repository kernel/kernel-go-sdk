package kernel

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

	sess, err := c.ForBrowser(&BrowserGetResponse{
		SessionID: "sid",
		BaseURL:   srv.URL + "/browser/kernel",
		CdpWsURL:  "wss://x/browser/cdp?jwt=j1",
	})
	if err != nil {
		t.Fatal(err)
	}

	hc := sess.HTTPClient()
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
