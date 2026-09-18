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

type cacheLifecycle struct {
	sniffResponse  bool
	evictSessionID string
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

// DeleteIfJWT removes the cached route only when its JWT still matches.
func (c *RouteCache) DeleteIfJWT(sessionID, jwt string) bool {
	if c == nil {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	route, ok := c.routes[sessionID]
	if !ok || route.JWT != jwt {
		return false
	}
	delete(c.routes, sessionID)
	return true
}

// DirectVMRoutingMiddleware rewrites allowlisted browser subresource requests to
// the browser VM using cached base_url and jwt data.
func DirectVMRoutingMiddleware(cache *RouteCache, subresources []string) option.Middleware {
	// allowPrefixes are path prefixes (relative to browsers/{id}/) eligible for
	// direct-to-VM routing, e.g. "curl" or "telemetry/stream". Matching is
	// segment-boundary aware (see matchesDirectVMPrefix), so "telemetry/stream"
	// covers "telemetry/stream[/...]" but NOT "telemetry/events" (a control-plane
	// historical read served from S2) or "telemetry/streamfoo".
	allowPrefixes := make([]string, 0, len(subresources))
	for _, s := range subresources {
		if trimmed := strings.Trim(strings.TrimSpace(s), "/"); trimmed != "" {
			allowPrefixes = append(allowPrefixes, trimmed)
		}
	}

	return func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		lifecycle, err := parseCacheLifecycle(req)
		if err != nil {
			return nil, err
		}
		origURL := cloneURL(req.URL)
		origHost := req.Host
		origAuth := req.Header.Get("Authorization")
		sessionID, subresource, suffix, ok := parseDirectVMPath(req.URL.Path)
		routed := false
		if ok {
			if matchesDirectVMPrefix(subresource+suffix, allowPrefixes) {
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
					routed = true
				}
			}
		}

		res, err := next(req)
		if err != nil {
			return res, err
		}
		if routed && isStaleDirectVMAuthResponse(res, req) {
			if sessionID != "" {
				cache.DeleteIfJWT(sessionID, req.URL.Query().Get("jwt"))
			}
			// Without a replayable body the fallback would send a truncated request,
			// so surface the auth failure instead; the route is already evicted, so a
			// later call from the caller goes to the control plane.
			if !prepareControlPlaneFallback(req, origURL, origHost, origAuth) {
				return res, nil
			}
			if res.Body != nil {
				_ = res.Body.Close()
			}
			res, err = next(req)
			if err != nil {
				return res, err
			}
		}
		return finalizeResponse(res, cache, lifecycle)
	}
}

func parseCacheLifecycle(req *http.Request) (cacheLifecycle, error) {
	if req == nil || req.URL == nil {
		return cacheLifecycle{}, nil
	}

	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	for i := 0; i < len(parts); i++ {
		switch parts[i] {
		case "browsers":
			return parseBrowserCacheLifecycle(req.Method, parts, i), nil
		case "browser_pools":
			return parseBrowserPoolCacheLifecycle(req, parts, i)
		}
	}
	return cacheLifecycle{}, nil
}

func parseBrowserCacheLifecycle(method string, parts []string, index int) cacheLifecycle {
	switch len(parts) - index {
	case 1:
		return cacheLifecycle{sniffResponse: true}
	case 2:
		if parts[index+1] == "" {
			return cacheLifecycle{}
		}
		lifecycle := cacheLifecycle{sniffResponse: true}
		if method == http.MethodDelete {
			lifecycle.evictSessionID = parts[index+1]
		}
		return lifecycle
	default:
		return cacheLifecycle{}
	}
}

func parseBrowserPoolCacheLifecycle(req *http.Request, parts []string, index int) (cacheLifecycle, error) {
	switch len(parts) - index {
	case 3:
		if parts[index+1] == "" || parts[index+2] == "" {
			return cacheLifecycle{}, nil
		}
		switch parts[index+2] {
		case "acquire":
			if req.Method != http.MethodPost {
				return cacheLifecycle{}, nil
			}
			return cacheLifecycle{sniffResponse: true}, nil
		case "release":
			if req.Method != http.MethodPost {
				return cacheLifecycle{}, nil
			}
			sessionID, err := parseBrowserPoolReleaseSessionID(req)
			if err != nil {
				return cacheLifecycle{}, err
			}
			return cacheLifecycle{evictSessionID: sessionID}, nil
		default:
			return cacheLifecycle{}, nil
		}
	default:
		return cacheLifecycle{}, nil
	}
}

func parseBrowserPoolReleaseSessionID(req *http.Request) (string, error) {
	if req == nil || req.Body == nil {
		return "", nil
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	_ = req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.ContentLength = int64(len(body))

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", nil
	}
	sessionID, _ := payload["session_id"].(string)
	return strings.TrimSpace(sessionID), nil
}

func finalizeResponse(res *http.Response, cache *RouteCache, lifecycle cacheLifecycle) (*http.Response, error) {
	if lifecycle.sniffResponse {
		if err := sniffAndPopulateCache(res, cache); err != nil {
			return nil, err
		}
	}
	if lifecycle.evictSessionID != "" && isSuccessfulResponse(res) {
		cache.Delete(lifecycle.evictSessionID)
	}
	return res, nil
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

// matchesDirectVMPrefix reports whether tail (the request path after
// browsers/{id}/) is covered by an allow prefix, matching on segment boundaries:
// "telemetry/stream" matches "telemetry/stream" and "telemetry/stream/...", but
// not "telemetry/events" or "telemetry/streamfoo". This keeps historical
// control-plane reads (e.g. telemetry/events, served from S2) off the VM.
func matchesDirectVMPrefix(tail string, prefixes []string) bool {
	tail = strings.Trim(tail, "/")
	for _, p := range prefixes {
		if tail == p || strings.HasPrefix(tail, p+"/") {
			return true
		}
	}
	return false
}

func prepareControlPlaneFallback(req *http.Request, origURL *url.URL, origHost, origAuth string) bool {
	if req.Body != nil && req.GetBody == nil {
		return false
	}
	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return false
		}
		req.Body = body
	}
	req.URL = origURL
	req.Host = origHost
	if origAuth != "" {
		req.Header.Set("Authorization", origAuth)
	} else {
		req.Header.Del("Authorization")
	}
	q := req.URL.Query()
	q.Del("jwt")
	req.URL.RawQuery = q.Encode()
	return true
}

func isStaleDirectVMAuthResponse(res *http.Response, req *http.Request) bool {
	if res == nil || req == nil || req.URL == nil {
		return false
	}
	if res.StatusCode != http.StatusUnauthorized && res.StatusCode != http.StatusForbidden {
		return false
	}
	return req.URL.Query().Get("jwt") != ""
}

func cloneURL(u *url.URL) *url.URL {
	if u == nil {
		return nil
	}
	c := *u
	if u.User != nil {
		user := *u.User
		c.User = &user
	}
	return &c
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
