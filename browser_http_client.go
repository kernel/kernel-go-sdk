package kernel

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/lib/browserrouting"
	"github.com/kernel/kernel-go-sdk/option"
)

// HTTPClient returns an [http.Client] that performs HTTP requests through the
// browser VM's internal /curl/raw path using cached browser route data.
func (r *BrowserService) HTTPClient(id string, opts ...option.RequestOption) (*http.Client, error) {
	opts = slices.Concat(r.Options, opts)
	cache := browserRouteCacheFromOptions(opts)
	if cache == nil {
		return nil, fmt.Errorf("kernel: browser route cache is not configured")
	}

	route, ok := cache.Load(id)
	if !ok {
		return nil, fmt.Errorf("kernel: browser route cache does not contain session %s", id)
	}

	cfg, err := requestconfig.NewRequestConfig(context.Background(), http.MethodGet, "https://example.com", nil, nil, opts...)
	if err != nil {
		return nil, err
	}

	return browserrouting.NewHTTPClient(route.BaseURL, route.JWT, cfg.HTTPClient), nil
}
