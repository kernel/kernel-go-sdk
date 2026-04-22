package kernel

import (
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/lib/browserscope"
	"github.com/kernel/kernel-go-sdk/option"
)

// BrowserRoutingConfig controls which browser subresources route directly to the browser VM.
type BrowserRoutingConfig struct {
	Enabled                bool
	DirectToVMSubresources []string
}

type browserRoutingOption struct {
	cache  *browserscope.RouteCache
	config BrowserRoutingConfig
}

type browserRouteCacheOption struct {
	cache *browserscope.RouteCache
}

// WithBrowserRouting enables direct-to-VM routing for the configured browser subresources.
func WithBrowserRouting(config BrowserRoutingConfig) option.RequestOption {
	return &browserRoutingOption{config: config}
}

func (o *browserRoutingOption) Apply(r *requestconfig.RequestConfig) error {
	if !o.config.Enabled {
		return nil
	}
	r.Middlewares = append(r.Middlewares, browserscope.DirectVMRoutingMiddleware(o.cache, o.config.DirectToVMSubresources))
	return nil
}

func (o *browserRoutingOption) browserRouteCache() *browserscope.RouteCache {
	return o.cache
}

func (o *browserRouteCacheOption) Apply(*requestconfig.RequestConfig) error {
	return nil
}

func (o *browserRouteCacheOption) browserRouteCache() *browserscope.RouteCache {
	return o.cache
}

func withBrowserRouteCache(cache *browserscope.RouteCache) option.RequestOption {
	return &browserRouteCacheOption{cache: cache}
}

func browserRouteCacheFromOptions(opts []option.RequestOption) *browserscope.RouteCache {
	for _, opt := range opts {
		if carrier, ok := opt.(interface{ browserRouteCache() *browserscope.RouteCache }); ok {
			if cache := carrier.browserRouteCache(); cache != nil {
				return cache
			}
		}
	}
	return nil
}

func primeBrowserRouteCache(opts []option.RequestOption, refs ...browserscope.Ref) {
	cache := browserRouteCacheFromOptions(opts)
	for _, ref := range refs {
		if cache != nil {
			_ = cache.Prime(ref)
		}
	}
}
