package browserrouting

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/kernel/kernel-go-sdk/option"
)

// Route identifies a cached direct-to-VM transport for one browser session.
type Route struct {
	SessionID string
	BaseURL   string
	JWT       string
}

// RouteCache stores browser session transport details keyed by session_id.
type RouteCache struct {
	mu     sync.RWMutex
	routes map[string]Route
}

// NewRouteCache returns an empty browser route cache.
func NewRouteCache() *RouteCache {
	return &RouteCache{routes: map[string]Route{}}
}

// Load returns the cached route for the given session id.
func (c *RouteCache) Load(sessionID string) (Route, bool) {
	if c == nil {
		return Route{}, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	route, ok := c.routes[sessionID]
	return route, ok
}

// Store normalizes and caches the given route.
func (c *RouteCache) Store(route Route) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.routes[strings.TrimSpace(route.SessionID)] = Route{
		SessionID: strings.TrimSpace(route.SessionID),
		BaseURL:   strings.TrimSpace(route.BaseURL),
		JWT:       strings.TrimSpace(route.JWT),
	}
}

// Delete removes a cached route.
func (c *RouteCache) Delete(sessionID string) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.routes, sessionID)
}

// DirectVMRoutingMiddleware rewrites allowlisted browser subresource requests to
// the browser VM using cached base_url and jwt data.
func DirectVMRoutingMiddleware(cache *RouteCache, subresources []string) option.Middleware {
	allowed := map[string]struct{}{}
	for _, subresource := range subresources {
		if trimmed := strings.TrimSpace(subresource); trimmed != "" {
			allowed[trimmed] = struct{}{}
		}
	}

	return func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		cacheSessionID, cacheablePath := parseBrowserMetadataPath(req.URL.Path)
		sessionID, subresource, suffix, ok := parseDirectVMPath(req.URL.Path)
		if !ok {
			res, err := next(req)
			if err != nil {
				return res, err
			}
			if req.Method == http.MethodDelete && cacheSessionID != "" && isSuccessfulResponse(res) {
				cache.Delete(cacheSessionID)
			}
			if cacheablePath {
				if err := sniffAndPopulateCache(res, cache); err != nil {
					return nil, err
				}
			}
			return res, nil
		}
		if _, ok := allowed[subresource]; ok {
			route, ok := cache.Load(sessionID)
			if ok {
				base, err := url.Parse(route.BaseURL)
				if err != nil {
					return nil, err
				}
				req.Header.Del("Authorization")
				if route.JWT != "" {
					q := req.URL.Query()
					if q.Get("jwt") == "" {
						q.Set("jwt", route.JWT)
						req.URL.RawQuery = q.Encode()
					}
				}

				req.URL.Scheme = base.Scheme
				req.URL.Host = base.Host
				req.Host = base.Host
				req.URL.Path = joinURLPath(base.Path, subresource, suffix)
				req.URL.RawPath = ""
			}
		}

		res, err := next(req)
		if err != nil {
			return res, err
		}
		if req.Method == http.MethodDelete && cacheSessionID != "" && isSuccessfulResponse(res) {
			cache.Delete(cacheSessionID)
		}
		if cacheablePath {
			if err := sniffAndPopulateCache(res, cache); err != nil {
				return nil, err
			}
		}
		return res, nil
	}
}

func parseBrowserMetadataPath(path string) (sessionID string, ok bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i := 0; i < len(parts); i++ {
		if parts[i] != "browsers" {
			continue
		}
		switch len(parts) - i {
		case 1:
			return "", true
		case 2:
			if parts[i+1] == "" {
				return "", false
			}
			return parts[i+1], true
		}
	}
	return "", false
}

func sniffAndPopulateCache(res *http.Response, cache *RouteCache) error {
	if res == nil || res.Body == nil || cache == nil || !isSuccessfulResponse(res) || !isJSONResponse(res.Header) {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	_ = res.Body.Close()
	res.Body = io.NopCloser(bytes.NewReader(body))
	res.ContentLength = int64(len(body))

	var value any
	if err := json.Unmarshal(body, &value); err != nil {
		return nil
	}
	populateCache(value, cache)
	return nil
}

func isSuccessfulResponse(res *http.Response) bool {
	return res != nil && res.StatusCode >= 200 && res.StatusCode < 300
}

func isJSONResponse(header http.Header) bool {
	mediaType, _, _ := mime.ParseMediaType(header.Get("Content-Type"))
	return strings.Contains(mediaType, "application/json") || strings.HasSuffix(mediaType, "+json")
}

func populateCache(value any, cache *RouteCache) {
	if route, ok := routeFromValue(value); ok {
		cache.Store(route)
	}

	switch value := value.(type) {
	case []any:
		for _, item := range value {
			populateCache(item, cache)
		}
	case map[string]any:
		for _, child := range value {
			if child != nil {
				populateCache(child, cache)
			}
		}
	}
}

func routeFromValue(value any) (Route, bool) {
	record, ok := value.(map[string]any)
	if !ok {
		return Route{}, false
	}

	sessionID, _ := record["session_id"].(string)
	baseURL, _ := record["base_url"].(string)
	jwt, _ := record["jwt"].(string)
	cdpWsURL, _ := record["cdp_ws_url"].(string)
	ref, err := (Ref{
		SessionID: sessionID,
		BaseURL:   baseURL,
		JWT:       jwt,
		CdpWsURL:  cdpWsURL,
	}).Normalize()
	if err != nil {
		return Route{}, false
	}

	return Route{
		SessionID: ref.SessionID,
		BaseURL:   ref.BaseURL,
		JWT:       ref.JWT,
	}, true
}

func parseDirectVMPath(path string) (sessionID, subresource, suffix string, ok bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i := 0; i+2 < len(parts); i++ {
		if parts[i] != "browsers" {
			continue
		}
		sessionID = parts[i+1]
		subresource = parts[i+2]
		if sessionID == "" || subresource == "" {
			return "", "", "", false
		}
		if i+3 < len(parts) {
			suffix = "/" + strings.Join(parts[i+3:], "/")
		}
		return sessionID, subresource, suffix, true
	}
	return "", "", "", false
}

func joinURLPath(basePath, subresource, suffix string) string {
	base := "/" + strings.Trim(strings.TrimSpace(basePath), "/")
	if base == "/" {
		base = ""
	}
	out := base + "/" + strings.TrimPrefix(subresource, "/")
	if suffix != "" {
		out += "/" + strings.TrimPrefix(suffix, "/")
	}
	return out
}
