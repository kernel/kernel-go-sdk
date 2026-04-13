package browserscope

import (
	"net/http"
	"strings"

	"github.com/kernel/kernel-go-sdk/option"
)

// BrowserSessionMiddleware prepares requests for a browser session base_url:
// it strips the control-plane browsers/{session_id} path prefix, attaches jwt,
// and removes Authorization so the Kernel API key is not forwarded to the browser.
func BrowserSessionMiddleware(sessionID, jwt string) option.Middleware {
	marker := "/browsers/" + sessionID + "/"
	return func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		req.Header.Del("Authorization")

		if jwt != "" {
			q := req.URL.Query()
			if q.Get("jwt") == "" {
				q.Set("jwt", jwt)
				req.URL.RawQuery = q.Encode()
			}
		}

		if sessionID != "" {
			if idx := strings.Index(req.URL.Path, marker); idx >= 0 {
				prefix := strings.TrimRight(req.URL.Path[:idx], "/")
				rest := strings.TrimPrefix(req.URL.Path[idx+len(marker):], "/")
				if rest == "" {
					rest = "/"
				} else {
					rest = "/" + rest
				}
				req.URL.Path = prefix + rest
			}
		}

		return next(req)
	}
}
