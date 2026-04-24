package kernel

import (
	"os"
	"strings"

	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/lib/browserrouting"
	"github.com/kernel/kernel-go-sdk/option"
)

const browserRoutingSubresourcesEnv = "KERNEL_BROWSER_ROUTING_SUBRESOURCES"

type browserRoutingOption struct {
	cache        *browserrouting.RouteCache
	subresources []string
}

type browserRouteCacheOption struct {
	cache *browserrouting.RouteCache
}

func withBrowserRoutingSubresources(cache *browserrouting.RouteCache, subresources []string) option.RequestOption {
	return &browserRoutingOption{cache: cache, subresources: subresources}
}

func (o *browserRoutingOption) Apply(r *requestconfig.RequestConfig) error {
	r.Middlewares = append(r.Middlewares, browserrouting.DirectVMRoutingMiddleware(o.cache, o.subresources))
	return nil
}

func (o *browserRoutingOption) browserRouteCache() *browserrouting.RouteCache {
	return o.cache
}

func (o *browserRouteCacheOption) Apply(*requestconfig.RequestConfig) error {
	return nil
}

func (o *browserRouteCacheOption) browserRouteCache() *browserrouting.RouteCache {
	return o.cache
}

func withBrowserRouteCache(cache *browserrouting.RouteCache) option.RequestOption {
	return &browserRouteCacheOption{cache: cache}
}

func browserRouteCacheFromOptions(opts []option.RequestOption) *browserrouting.RouteCache {
	for _, opt := range opts {
		if carrier, ok := opt.(interface{ browserRouteCache() *browserrouting.RouteCache }); ok {
			if cache := carrier.browserRouteCache(); cache != nil {
				return cache
			}
		}
	}
	return nil
}

func storeBrowserRouteCache(opts []option.RequestOption, refs ...browserrouting.Ref) {
	cache := browserRouteCacheFromOptions(opts)
	for _, ref := range refs {
		route, ok := browserRouteFromRef(ref)
		if cache != nil && ok {
			cache.Store(route)
		}
	}
}

func browserRouteFromRef(ref browserrouting.Ref) (browserrouting.Route, bool) {
	norm, err := ref.Normalize()
	if err != nil {
		return browserrouting.Route{}, false
	}
	return browserrouting.Route{
		SessionID: norm.SessionID,
		BaseURL:   norm.BaseURL,
		JWT:       norm.JWT,
	}, true
}

func browserRoutingSubresourcesFromEnv() []string {
	raw, ok := os.LookupEnv(browserRoutingSubresourcesEnv)
	if !ok {
		return []string{"curl"}
	}
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}

	subresources := make([]string, 0)
	for _, part := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			subresources = append(subresources, trimmed)
		}
	}
	return subresources
}
