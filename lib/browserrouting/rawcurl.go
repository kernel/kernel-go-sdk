package browserrouting

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// RawCURLRoundTripper implements browser-egress HTTP by tunneling through the
// browser session base_url /curl/raw endpoint.
type RawCURLRoundTripper struct {
	browserBaseURL string
	jwt            string
	underlying     http.RoundTripper
}

var _ http.RoundTripper = (*RawCURLRoundTripper)(nil)

// NewHTTPClient returns an [http.Client] that performs browser egress HTTP via
// the browser session base_url and internal /curl/raw path.
func NewHTTPClient(browserBaseURL, jwt string, underlying *http.Client) *http.Client {
	if underlying == nil {
		underlying = http.DefaultClient
	}
	rt := underlying.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &http.Client{
		Transport: newRawCURLRoundTripper(browserBaseURL, jwt, rt),
		Timeout:   underlying.Timeout,
	}
}

// newRawCURLRoundTripper returns an [http.RoundTripper] that maps each request
// to {base_url}/curl/raw?jwt=...&url=<absolute-target>, preserving method,
// headers, and body. The caller's request URL must be an absolute http(s) URL.
func newRawCURLRoundTripper(browserBaseURL, jwt string, underlying http.RoundTripper) *RawCURLRoundTripper {
	if underlying == nil {
		underlying = http.DefaultTransport
	}
	return &RawCURLRoundTripper{
		browserBaseURL: strings.TrimRight(strings.TrimSpace(browserBaseURL), "/"),
		jwt:            strings.TrimSpace(jwt),
		underlying:     underlying,
	}
}

// RoundTrip implements [http.RoundTripper].
func (t *RawCURLRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL == nil || !req.URL.IsAbs() {
		return nil, fmt.Errorf("browserrouting: raw curl requires an absolute request URL (got %q)", req.URL)
	}
	if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
		return nil, fmt.Errorf("browserrouting: raw curl requires http or https scheme")
	}
	if t.browserBaseURL == "" {
		return nil, fmt.Errorf("browserrouting: browser base_url is required")
	}
	if t.jwt == "" {
		return nil, fmt.Errorf("browserrouting: jwt is required for raw curl")
	}

	target := req.URL.String()
	proxyURL, err := url.Parse(t.browserBaseURL + "/curl/raw")
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("jwt", t.jwt)
	q.Set("url", target)
	proxyURL.RawQuery = q.Encode()

	out := req.Clone(req.Context())
	out.URL = proxyURL
	out.Host = proxyURL.Host
	out.RequestURI = ""

	return t.underlying.RoundTrip(out)
}
