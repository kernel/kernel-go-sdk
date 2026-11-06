// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Create and manage proxy configurations for routing browser traffic.
//
// ProxyIPService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewProxyIPService] method instead.
type ProxyIPService struct {
	Options []option.RequestOption
}

// NewProxyIPService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewProxyIPService(opts ...option.RequestOption) (r ProxyIPService) {
	r = ProxyIPService{}
	r.Options = opts
	return
}

// Retrieve one reserved IP with the ISP proxies using it as their primary exit.
// proxy_count spans the organization; proxies lists only those in the current
// project. Returns 404 when the IP is not in the organization's allocation and 403
// with code reserved_isp_not_enabled when the organization has none. Scoped API
// keys require proxies:read on the whole project or organization, not on
// individual proxies; otherwise the operation returns 403 with code
// insufficient_scope.
func (r *ProxyIPService) Get(ctx context.Context, ipAddress string, opts ...option.RequestOption) (res *ProxyIPDetail, err error) {
	opts = slices.Concat(r.Options, opts)
	if ipAddress == "" {
		err = errors.New("missing required ip_address parameter")
		return nil, err
	}
	path := fmt.Sprintf("proxy_ips/%s", ipAddress)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List the IP addresses reserved for the organization's ISP proxies, with each
// IP's country and whether it can currently be assigned. Available only to
// organizations with a reserved ISP allocation; others receive a 403 with code
// reserved_isp_not_enabled. Served from the allocation itself, so IPs appear
// whether or not a proxy uses them. Scoped API keys require proxies:read on the
// whole project or organization, not on individual proxies; otherwise the
// operation returns 403 with code insufficient_scope.
func (r *ProxyIPService) List(ctx context.Context, query ProxyIPListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[ProxyIP], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "proxy_ips"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List the IP addresses reserved for the organization's ISP proxies, with each
// IP's country and whether it can currently be assigned. Available only to
// organizations with a reserved ISP allocation; others receive a 403 with code
// reserved_isp_not_enabled. Served from the allocation itself, so IPs appear
// whether or not a proxy uses them. Scoped API keys require proxies:read on the
// whole project or organization, not on individual proxies; otherwise the
// operation returns 403 with code insufficient_scope.
func (r *ProxyIPService) ListAutoPaging(ctx context.Context, query ProxyIPListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[ProxyIP] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// One IP address reserved for the organization's ISP proxies.
type ProxyIP struct {
	// ISO 3166 country code of the IP.
	Country string `json:"country" api:"required"`
	// Reserved exit IP address.
	IPAddress string `json:"ip_address" api:"required"`
	// Whether the IP can currently be assigned. An excluded IP is listed but cannot be
	// pinned or moved to; proxies already on it keep working.
	//
	// Any of "active", "excluded".
	Status ProxyIPStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		IPAddress   respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyIP) RawJSON() string { return r.JSON.raw }
func (r *ProxyIP) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Whether the IP can currently be assigned. An excluded IP is listed but cannot be
// pinned or moved to; proxies already on it keep working.
type ProxyIPStatus string

const (
	ProxyIPStatusActive   ProxyIPStatus = "active"
	ProxyIPStatusExcluded ProxyIPStatus = "excluded"
)

// A reserved IP with the ISP proxies using it as their primary exit. Reservations
// are organization-wide, so proxy_count covers every project while proxies lists
// only those in the current project.
type ProxyIPDetail struct {
	// ISO 3166 country code of the IP.
	Country string `json:"country" api:"required"`
	// Reserved exit IP address.
	IPAddress string `json:"ip_address" api:"required"`
	// ISP proxies in the current project whose primary exit IP is this address.
	Proxies []ProxyIPDetailProxy `json:"proxies" api:"required"`
	// Number of ISP proxies across the organization whose primary exit IP is this
	// address.
	ProxyCount int64 `json:"proxy_count" api:"required"`
	// Whether the IP can currently be assigned. An excluded IP is listed but cannot be
	// pinned or moved to; proxies already on it keep working.
	//
	// Any of "active", "excluded".
	Status ProxyIPDetailStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		IPAddress   respjson.Field
		Proxies     respjson.Field
		ProxyCount  respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyIPDetail) RawJSON() string { return r.JSON.raw }
func (r *ProxyIPDetail) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An ISP proxy whose primary exit is the reserved IP.
type ProxyIPDetailProxy struct {
	ID   string `json:"id" api:"required"`
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyIPDetailProxy) RawJSON() string { return r.JSON.raw }
func (r *ProxyIPDetailProxy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Whether the IP can currently be assigned. An excluded IP is listed but cannot be
// pinned or moved to; proxies already on it keep working.
type ProxyIPDetailStatus string

const (
	ProxyIPDetailStatusActive   ProxyIPDetailStatus = "active"
	ProxyIPDetailStatusExcluded ProxyIPDetailStatus = "excluded"
)

type ProxyIPListParams struct {
	// ISO 3166 country code to filter by. Must be a country ISP proxies support.
	Country param.Opt[string] `query:"country,omitzero" json:"-"`
	// Limit the number of IPs to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Offset the number of IPs to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Filter by whether the IP can currently be assigned.
	//
	// Any of "active", "excluded".
	Status ProxyIPListParamsStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ProxyIPListParams]'s query parameters as `url.Values`.
func (r ProxyIPListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter by whether the IP can currently be assigned.
type ProxyIPListParamsStatus string

const (
	ProxyIPListParamsStatusActive   ProxyIPListParamsStatus = "active"
	ProxyIPListParamsStatusExcluded ProxyIPListParamsStatus = "excluded"
)
