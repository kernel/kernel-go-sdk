package browserscope

import (
	"fmt"
	"strings"
)

// Ref identifies a browser session for browser-scoped HTTP calls. SessionID is
// reserved for future client-side routing; browser-scoped requests rewrite
// the /browsers/{SessionID}/ path segment against the returned base_url.
type Ref struct {
	SessionID string
	BaseURL   string
	JWT       string
	CdpWsURL  string
}

// Normalize validates fields and fills JWT from CdpWsURL when JWT is empty.
func (r Ref) Normalize() (Ref, error) {
	if strings.TrimSpace(r.BaseURL) == "" {
		return Ref{}, fmt.Errorf("browserscope: base_url is required")
	}
	if strings.TrimSpace(r.SessionID) == "" {
		return Ref{}, fmt.Errorf("browserscope: session_id is required")
	}
	out := r
	out.BaseURL = strings.TrimSpace(r.BaseURL)
	out.SessionID = strings.TrimSpace(r.SessionID)
	out.JWT = strings.TrimSpace(r.JWT)
	out.CdpWsURL = strings.TrimSpace(r.CdpWsURL)
	if out.JWT == "" {
		src := out.CdpWsURL
		if src == "" {
			return Ref{}, fmt.Errorf("browserscope: jwt or cdp_ws_url is required")
		}
		jwt, err := jwtFromWebSocketURL(src)
		if err != nil {
			return Ref{}, err
		}
		out.JWT = jwt
	}
	return out, nil
}
