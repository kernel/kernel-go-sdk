package kernel

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/lib/browserscope"
	"github.com/kernel/kernel-go-sdk/option"
)

// BrowserSessionClient is a browser-scoped view of a browser session. Subresources
// use the session base_url and do not repeat the browser id in method
// signatures. SessionID is exposed for future routing extensions.
type BrowserSessionClient struct {
	sessionID  string
	opts       []option.RequestOption
	kernelBase string
	jwt        string

	Replays    BrowserSessionReplayService
	Fs         BrowserSessionFService
	Process    BrowserSessionProcessService
	Logs       BrowserSessionLogService
	Computer   BrowserSessionComputerService
	Playwright BrowserSessionPlaywrightService
}

// SessionID returns the control-plane browser session id.
func (b *BrowserSessionClient) SessionID() string { return b.sessionID }

// HTTPClient returns an [http.Client] that performs egress HTTP through the
// browser's Chrome network stack via the browser session /curl/raw proxy. Each
// request must use an absolute http(s) URL; it is not rewritten to expose
// /curl/raw in the public API.
func (b *BrowserSessionClient) HTTPClient() *http.Client {
	cfg, err := requestconfig.PreRequestOptions(b.opts...)
	if err != nil {
		return &http.Client{
			Transport: browserscope.NewRawCURLRoundTripper(b.kernelBase, b.jwt, nil),
		}
	}
	underlying := cfg.HTTPClient
	if underlying == nil {
		underlying = http.DefaultClient
	}
	rt := underlying.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &http.Client{
		Transport: browserscope.NewRawCURLRoundTripper(b.kernelBase, b.jwt, rt),
		Timeout:   underlying.Timeout,
	}
}

// ForBrowser returns a [BrowserSessionClient] for the given browser value.
// Supported types are [browserscope.Ref], [*BrowserNewResponse], [*BrowserGetResponse],
// [*BrowserListResponse], and [*BrowserUpdateResponse].
func (c *Client) ForBrowser(v any, opts ...option.RequestOption) (*BrowserSessionClient, error) {
	ref, err := browserSessionRefFrom(v)
	if err != nil {
		return nil, err
	}
	norm, err := ref.Normalize()
	if err != nil {
		return nil, err
	}

	scoped := slices.Concat(
		c.Options,
		[]option.RequestOption{
			option.WithBaseURL(norm.BaseURL),
			option.WithMiddleware(browserscope.BrowserSessionMiddleware(norm.SessionID, norm.JWT)),
		},
		opts,
	)

	innerFs := NewBrowserFService(scoped...)
	return &BrowserSessionClient{
		sessionID:  norm.SessionID,
		opts:       scoped,
		kernelBase: norm.BaseURL,
		jwt:        norm.JWT,
		Replays:    BrowserSessionReplayService{inner: NewBrowserReplayService(scoped...), id: norm.SessionID},
		Fs: BrowserSessionFService{
			inner: innerFs,
			id:    norm.SessionID,
			Watch: BrowserSessionFWatchService{inner: innerFs.Watch, id: norm.SessionID},
		},
		Process:    BrowserSessionProcessService{inner: NewBrowserProcessService(scoped...), id: norm.SessionID},
		Logs:       BrowserSessionLogService{inner: NewBrowserLogService(scoped...), id: norm.SessionID},
		Computer:   BrowserSessionComputerService{inner: NewBrowserComputerService(scoped...), id: norm.SessionID},
		Playwright: BrowserSessionPlaywrightService{inner: NewBrowserPlaywrightService(scoped...), id: norm.SessionID},
	}, nil
}

func browserSessionRefFrom(v any) (browserscope.Ref, error) {
	switch t := v.(type) {
	case browserscope.Ref:
		return t, nil
	case *BrowserNewResponse:
		if t == nil {
			return browserscope.Ref{}, fmt.Errorf("kernel: ForBrowser: nil *BrowserNewResponse")
		}
		return browserscope.Ref{SessionID: t.SessionID, BaseURL: t.BaseURL, CdpWsURL: t.CdpWsURL}, nil
	case *BrowserGetResponse:
		if t == nil {
			return browserscope.Ref{}, fmt.Errorf("kernel: ForBrowser: nil *BrowserGetResponse")
		}
		return browserscope.Ref{SessionID: t.SessionID, BaseURL: t.BaseURL, CdpWsURL: t.CdpWsURL}, nil
	case *BrowserListResponse:
		if t == nil {
			return browserscope.Ref{}, fmt.Errorf("kernel: ForBrowser: nil *BrowserListResponse")
		}
		return browserscope.Ref{SessionID: t.SessionID, BaseURL: t.BaseURL, CdpWsURL: t.CdpWsURL}, nil
	case *BrowserUpdateResponse:
		if t == nil {
			return browserscope.Ref{}, fmt.Errorf("kernel: ForBrowser: nil *BrowserUpdateResponse")
		}
		return browserscope.Ref{SessionID: t.SessionID, BaseURL: t.BaseURL, CdpWsURL: t.CdpWsURL}, nil
	default:
		return browserscope.Ref{}, fmt.Errorf("kernel: ForBrowser: unsupported type %T", v)
	}
}
