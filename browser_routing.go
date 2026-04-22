package kernel

import (
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/lib/browserrouting"
	"github.com/kernel/kernel-go-sdk/option"
)

// BrowserRoutingConfig controls which browser subresources route directly to the browser VM.
type BrowserRoutingConfig struct {
	Enabled      bool
	Subresources []string
}

type browserRoutingOption struct {
	cache  *browserrouting.RouteCache
	config BrowserRoutingConfig
}

type browserRouteCacheOption struct {
	cache *browserrouting.RouteCache
}

// WithBrowserRouting enables direct-to-VM routing for the configured browser subresources.
func WithBrowserRouting(config BrowserRoutingConfig) option.RequestOption {
	return &browserRoutingOption{config: config}
}

func (o *browserRoutingOption) Apply(r *requestconfig.RequestConfig) error {
	if !o.config.Enabled {
		return nil
	}
	r.Middlewares = append(r.Middlewares, browserrouting.DirectVMRoutingMiddleware(o.cache, o.config.Subresources))
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
