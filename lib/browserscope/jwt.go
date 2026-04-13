package browserscope

import (
	"fmt"
	"net/url"
	"strings"
)

// jwtFromWebSocketURL extracts the session jwt query parameter from a browser
// websocket URL (for example cdp_ws_url or webdriver_ws_url).
func jwtFromWebSocketURL(wsURL string) (string, error) {
	wsURL = strings.TrimSpace(wsURL)
	if wsURL == "" {
		return "", fmt.Errorf("browserscope: empty websocket url")
	}
	u, err := url.Parse(wsURL)
	if err != nil {
		return "", fmt.Errorf("browserscope: parse websocket url: %w", err)
	}
	jwt := u.Query().Get("jwt")
	if jwt == "" {
		return "", fmt.Errorf("browserscope: missing jwt query parameter")
	}
	return jwt, nil
}
