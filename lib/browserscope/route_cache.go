package browserscope

import (
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
func DirectVMRoutingMiddleware(cache *RouteCache, directToVMSubresources []string) option.Middleware {
	allowed := map[string]struct{}{}
	for _, subresource := range directToVMSubresources {
		if trimmed := strings.TrimSpace(subresource); trimmed != "" {
			allowed[trimmed] = struct{}{}
		}
	}

	return func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		sessionID, subresource, suffix, ok := parseDirectVMPath(req.URL.Path)
		if !ok {
			return next(req)
		}
		if _, ok := allowed[subresource]; !ok {
			return next(req)
		}
		route, ok := cache.Load(sessionID)
		if !ok {
			return next(req)
		}

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

		return next(req)
	}
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
