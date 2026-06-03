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

// fallbackEndpoint identifies a routed subresource path that is eligible for
// control-plane fallback, expressed against the parsed routed path
// (subresource + suffix).
type fallbackEndpoint struct {
	subresource string
	suffix      string
}

// fallbackEligibleEndpoints is the routing-layer registry of endpoints that
// may fall back to the control plane when the VM reports the browser is gone.
// Everything not listed here is fallback-OFF by default. Adding a future
// eligible endpoint is a one-line edit here.
var fallbackEligibleEndpoints = map[fallbackEndpoint]struct{}{
	// PROSPECTIVE: GET /browsers/{id}/telemetry/events. The telemetry pull
	// method does not exist yet; this pre-wires the opt-in so fallback works
	// the moment that method ships.
	{subresource: "telemetry", suffix: "/events"}: {},
}

// isFallbackEligible reports whether the parsed routed path opts into
// control-plane fallback.
func isFallbackEligible(subresource, suffix string) bool {
	_, ok := fallbackEligibleEndpoints[fallbackEndpoint{subresource: subresource, suffix: suffix}]
	return ok
}

// browserGoneCode is the JSON body code metro-api (kernel#2317) returns, with
// HTTP 404, when a routed request targets a deleted/gone browser. A live VM's
// own 404 does not carry this code, and transient/upstream failures stay 5xx.
const browserGoneCode = "browser_gone"

// originalRequest captures the control-plane-bound request before the routing
// middleware mutates it, so it can be replayed verbatim on fallback.
type originalRequest struct {
	url           *url.URL
	host          string
	authorization []string
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
		lifecycle, err := parseCacheLifecycle(req)
		if err != nil {
			return nil, err
		}

		var (
			routed     bool
			routedSess string
			routedSub  string
			routedSuf  string
			snapshot   originalRequest
		)
		sessionID, subresource, suffix, ok := parseDirectVMPath(req.URL.Path)
		if ok {
			if _, ok := allowed[subresource]; ok {
				route, ok := cache.Load(sessionID)
				if ok {
					base, err := url.Parse(route.BaseURL)
					if err != nil {
						return nil, err
					}

					// Snapshot the original control-plane-bound request before
					// mutating it, so fallback can replay it verbatim.
					snapshot = snapshotRequest(req)
					routed = true
					routedSess = sessionID
					routedSub = subresource
					routedSuf = suffix

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
		}

		res, err := next(req)
		if err != nil {
			return res, err
		}

		if routed && shouldFallbackToControlPlane(req.Method, routedSub, routedSuf, res) {
			return controlPlaneFallback(req, next, cache, routedSess, snapshot)
		}

		return finalizeResponse(res, cache, lifecycle)
	}
}

// snapshotRequest captures the control-plane-bound request state (URL, Host,
// Authorization) before the routing middleware rewrites it.
func snapshotRequest(req *http.Request) originalRequest {
	snap := originalRequest{host: req.Host}
	if req.URL != nil {
		urlCopy := *req.URL
		snap.url = &urlCopy
	}
	if auth, ok := req.Header["Authorization"]; ok {
		snap.authorization = append([]string(nil), auth...)
	}
	return snap
}

// shouldFallbackToControlPlane reports whether a routed VM response warrants a
// control-plane fallback. It is the single point that decides fallback, and
// only inspects the body on a 404. On a 404 the body is buffered and restored
// so callers that do NOT fall back still receive it intact.
func shouldFallbackToControlPlane(method, subresource, suffix string, res *http.Response) bool {
	if method != http.MethodGet {
		return false
	}
	if !isFallbackEligible(subresource, suffix) {
		return false
	}
	if res == nil || res.StatusCode != http.StatusNotFound {
		return false
	}
	return responseHasBrowserGoneCode(res)
}

// responseHasBrowserGoneCode buffers the response body, restores it for later
// callers, and reports whether its JSON body carries code == "browser_gone".
func responseHasBrowserGoneCode(res *http.Response) bool {
	if res == nil || res.Body == nil {
		return false
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false
	}
	_ = res.Body.Close()
	res.Body = io.NopCloser(bytes.NewReader(body))
	res.ContentLength = int64(len(body))

	var payload struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return false
	}
	return payload.Code == browserGoneCode
}

// controlPlaneFallback evicts the authoritatively-gone route and replays the
// original control-plane request exactly once, returning its response. It never
// loops back through VM routing.
func controlPlaneFallback(req *http.Request, next option.MiddlewareNext, cache *RouteCache, sessionID string, snapshot originalRequest) (*http.Response, error) {
	cache.Delete(sessionID)

	if snapshot.url != nil {
		urlCopy := *snapshot.url
		req.URL = &urlCopy
	}
	req.Host = snapshot.host
	req.Header.Del("Authorization")
	for _, value := range snapshot.authorization {
		req.Header.Add("Authorization", value)
	}

	return next(req)
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
