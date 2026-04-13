package browserscope

import (
	"net/http"
	"strings"

	"github.com/kernel/kernel-go-sdk/option"
)

const kernelHTTPPrefix = "/browser/kernel"

// MetroKernelMiddleware prepares requests for metro-api's /browser/kernel proxy:
// it strips the control-plane browsers/{session_id} path prefix, attaches jwt,
// and removes Authorization so the Kernel API key is not forwarded to the VM.
func MetroKernelMiddleware(sessionID, jwt string) option.Middleware {
	prefix := kernelHTTPPrefix + "/browsers/" + sessionID + "/"
	return func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		req.Header.Del("Authorization")

		if jwt != "" {
			q := req.URL.Query()
			if q.Get("jwt") == "" {
				q.Set("jwt", jwt)
				req.URL.RawQuery = q.Encode()
			}
		}

		if sessionID != "" && strings.HasPrefix(req.URL.Path, prefix) {
			rest := strings.TrimPrefix(req.URL.Path, prefix)
			if rest == "" {
				rest = "/"
			}
			if !strings.HasPrefix(rest, "/") {
				rest = "/" + rest
			}
			req.URL.Path = kernelHTTPPrefix + rest
		}

		return next(req)
	}
}
