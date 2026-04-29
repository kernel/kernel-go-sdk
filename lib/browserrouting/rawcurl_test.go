package browserrouting

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRawCURLRoundTripper(t *testing.T) {
	var sawPath, sawRaw string
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawPath = r.URL.Path
		sawRaw = r.URL.RawQuery
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte("ok"))
	}))
	defer up.Close()

	rt := newRawCURLRoundTripper(up.URL+"/browser/kernel", "jwt1", http.DefaultTransport)

	client := &http.Client{Transport: rt}
	req, err := http.NewRequest(http.MethodGet, "https://example.org/foo?bar=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusTeapot {
		t.Fatalf("status: %d", res.StatusCode)
	}
	b, _ := io.ReadAll(res.Body)
	if string(b) != "ok" {
		t.Fatalf("body: %q", b)
	}
	if sawPath != "/browser/kernel/curl/raw" {
		t.Fatalf("path: %s", sawPath)
	}
	if sawRaw == "" {
		t.Fatal("expected query")
	}
}

func TestRawCURLRoundTripperSuppressesGoTransportHeaders(t *testing.T) {
	var sawUserAgent, sawAcceptEncoding string
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawUserAgent = r.Header.Get("User-Agent")
		sawAcceptEncoding = r.Header.Get("Accept-Encoding")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer up.Close()

	client := NewHTTPClient(up.URL+"/browser/kernel", "jwt1", http.DefaultClient)
	req, err := http.NewRequest(http.MethodGet, "https://example.org/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if sawUserAgent != "" {
		t.Fatalf("unexpected user-agent: %q", sawUserAgent)
	}
	if sawAcceptEncoding != "" {
		t.Fatalf("unexpected accept-encoding: %q", sawAcceptEncoding)
	}
}

func TestRawCURLRoundTripperPreservesExplicitTransportHeaders(t *testing.T) {
	var sawUserAgent, sawAcceptEncoding string
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawUserAgent = r.Header.Get("User-Agent")
		sawAcceptEncoding = r.Header.Get("Accept-Encoding")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer up.Close()

	client := NewHTTPClient(up.URL+"/browser/kernel", "jwt1", http.DefaultClient)
	req, err := http.NewRequest(http.MethodGet, "https://example.org/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("User-Agent", "custom-agent")
	req.Header.Set("Accept-Encoding", "identity")
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if sawUserAgent != "custom-agent" {
		t.Fatalf("user-agent: %q", sawUserAgent)
	}
	if sawAcceptEncoding != "identity" {
		t.Fatalf("accept-encoding: %q", sawAcceptEncoding)
	}
}

func TestRawCURLRoundTripperRelativeURL(t *testing.T) {
	rt := newRawCURLRoundTripper("https://x/browser/kernel", "j", http.DefaultTransport)
	req, _ := http.NewRequest(http.MethodGet, "/relative", nil)
	_, err := rt.RoundTrip(req)
	if err == nil {
		t.Fatal("expected error for relative url")
	}
}
