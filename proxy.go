// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

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
// ProxyService contains methods and other services that help with interacting with
// the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewProxyService] method instead.
type ProxyService struct {
	Options []option.RequestOption
}

// NewProxyService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewProxyService(opts ...option.RequestOption) (r ProxyService) {
	r = ProxyService{}
	r.Options = opts
	return
}

// Create a new proxy configuration in the resolved project.
func (r *ProxyService) New(ctx context.Context, body ProxyNewParams, opts ...option.RequestOption) (res *ProxyNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "proxies"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a proxy in the resolved project by ID.
func (r *ProxyService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *ProxyGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("proxies/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a proxy's name. Proxy names are not unique and are not ID-or-name
// addressable on this endpoint; duplicate names are allowed. Name-based
// session-create lookups can remain ambiguous until callers resolve proxies by ID
// or the API adds a stronger uniqueness contract.
func (r *ProxyService) Update(ctx context.Context, id string, body ProxyUpdateParams, opts ...option.RequestOption) (res *ProxyUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("proxies/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List proxies in the resolved project.
func (r *ProxyService) List(ctx context.Context, query ProxyListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[ProxyListResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "proxies"
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

// List proxies in the resolved project.
func (r *ProxyService) ListAutoPaging(ctx context.Context, query ProxyListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[ProxyListResponse] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Soft delete a proxy. Sessions referencing it are not modified.
func (r *ProxyService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("proxies/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Run a health check on the proxy to verify it's working. Optionally specify a URL
// to test reachability against a specific target. For ISP and datacenter proxies,
// this reliably tests whether the target site is reachable from the proxy's stable
// exit IP. For residential and mobile proxies, the exit node varies between
// requests, so this validates proxy configuration and connectivity rather than
// guaranteeing site-specific reachability.
func (r *ProxyService) Check(ctx context.Context, id string, body ProxyCheckParams, opts ...option.RequestOption) (res *ProxyCheckResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("proxies/%s/check", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Configuration for routing traffic through a proxy.
type ProxyNewResponse struct {
	// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
	// worst: `mobile` > `residential` > `isp` > `datacenter`.
	//
	// Any of "datacenter", "isp", "residential", "mobile", "custom".
	Type ProxyNewResponseType `json:"type" api:"required"`
	ID   string               `json:"id"`
	// Hostnames that should bypass the parent proxy and connect directly.
	BypassHosts []string `json:"bypass_hosts"`
	// Configuration specific to the selected proxy `type`.
	Config ProxyNewResponseConfigUnion `json:"config"`
	// IP address that the proxy uses when making requests.
	IPAddress string `json:"ip_address"`
	// Timestamp of the last health check performed on this proxy.
	LastChecked time.Time `json:"last_checked" format:"date-time"`
	// Readable name of the proxy.
	Name string `json:"name"`
	// Protocol to use for the proxy connection.
	//
	// Any of "http", "https".
	Protocol ProxyNewResponseProtocol `json:"protocol"`
	// Current health status of the proxy.
	//
	// Any of "available", "unavailable".
	Status ProxyNewResponseStatus `json:"status"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ID          respjson.Field
		BypassHosts respjson.Field
		Config      respjson.Field
		IPAddress   respjson.Field
		LastChecked respjson.Field
		Name        respjson.Field
		Protocol    respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyNewResponse) RawJSON() string { return r.JSON.raw }
func (r *ProxyNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
// worst: `mobile` > `residential` > `isp` > `datacenter`.
type ProxyNewResponseType string

const (
	ProxyNewResponseTypeDatacenter  ProxyNewResponseType = "datacenter"
	ProxyNewResponseTypeIsp         ProxyNewResponseType = "isp"
	ProxyNewResponseTypeResidential ProxyNewResponseType = "residential"
	ProxyNewResponseTypeMobile      ProxyNewResponseType = "mobile"
	ProxyNewResponseTypeCustom      ProxyNewResponseType = "custom"
)

// ProxyNewResponseConfigUnion contains all possible properties and values from
// [ProxyNewResponseConfigDatacenter], [ProxyNewResponseConfigIsp],
// [ProxyNewResponseConfigResidential], [ProxyNewResponseConfigMobile],
// [ProxyNewResponseConfigCustom].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ProxyNewResponseConfigUnion struct {
	Country string `json:"country"`
	// This field is from variant [ProxyNewResponseConfigResidential].
	Asn  string `json:"asn"`
	City string `json:"city"`
	// This field is from variant [ProxyNewResponseConfigResidential].
	Os    string `json:"os"`
	State string `json:"state"`
	// This field is from variant [ProxyNewResponseConfigResidential].
	Zip string `json:"zip"`
	// This field is from variant [ProxyNewResponseConfigCustom].
	Host string `json:"host"`
	// This field is from variant [ProxyNewResponseConfigCustom].
	Port int64 `json:"port"`
	// This field is from variant [ProxyNewResponseConfigCustom].
	HasCaBundle bool `json:"has_ca_bundle"`
	// This field is from variant [ProxyNewResponseConfigCustom].
	HasPassword bool `json:"has_password"`
	// This field is from variant [ProxyNewResponseConfigCustom].
	Username string `json:"username"`
	JSON     struct {
		Country     respjson.Field
		Asn         respjson.Field
		City        respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		Host        respjson.Field
		Port        respjson.Field
		HasCaBundle respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		raw         string
	} `json:"-"`
}

func (u ProxyNewResponseConfigUnion) AsDatacenter() (v ProxyNewResponseConfigDatacenter) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyNewResponseConfigUnion) AsIsp() (v ProxyNewResponseConfigIsp) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyNewResponseConfigUnion) AsResidential() (v ProxyNewResponseConfigResidential) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyNewResponseConfigUnion) AsMobile() (v ProxyNewResponseConfigMobile) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyNewResponseConfigUnion) AsCustom() (v ProxyNewResponseConfigCustom) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProxyNewResponseConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ProxyNewResponseConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a datacenter proxy.
type ProxyNewResponseConfigDatacenter struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country string `json:"country"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyNewResponseConfigDatacenter) RawJSON() string { return r.JSON.raw }
func (r *ProxyNewResponseConfigDatacenter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyNewResponseConfigIsp struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country string `json:"country"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyNewResponseConfigIsp) RawJSON() string { return r.JSON.raw }
func (r *ProxyNewResponseConfigIsp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyNewResponseConfigResidential struct {
	// Autonomous system number. See https://bgp.potaroo.net/cidr/autnums.html
	Asn string `json:"asn"`
	// City name (no spaces, e.g. `sanfrancisco`). If provided, `country` must also be
	// provided.
	City string `json:"city"`
	// ISO 3166 country code.
	Country string `json:"country"`
	// Operating system of the residential device.
	//
	// Any of "windows", "macos", "android".
	//
	// Deprecated: deprecated
	Os string `json:"os"`
	// Two-letter state code.
	State string `json:"state"`
	// US ZIP code.
	Zip string `json:"zip"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Asn         respjson.Field
		City        respjson.Field
		Country     respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyNewResponseConfigResidential) RawJSON() string { return r.JSON.raw }
func (r *ProxyNewResponseConfigResidential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for mobile proxies.
type ProxyNewResponseConfigMobile struct {
	// Provider city alias. Mobile carrier routing can make observed geo vary.
	City string `json:"city"`
	// ISO 3166 country code
	Country string `json:"country"`
	// US-only state code. Mobile carrier routing can make observed geo vary.
	State string `json:"state"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		City        respjson.Field
		Country     respjson.Field
		State       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyNewResponseConfigMobile) RawJSON() string { return r.JSON.raw }
func (r *ProxyNewResponseConfigMobile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a custom proxy (e.g., private proxy server).
type ProxyNewResponseConfigCustom struct {
	// Proxy host address or IP.
	Host string `json:"host" api:"required"`
	// Proxy port.
	Port int64 `json:"port" api:"required"`
	// Whether the proxy has a custom CA bundle configured.
	HasCaBundle bool `json:"has_ca_bundle"`
	// Whether the proxy has a password.
	HasPassword bool `json:"has_password"`
	// Username for proxy authentication.
	Username string `json:"username"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Host        respjson.Field
		Port        respjson.Field
		HasCaBundle respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyNewResponseConfigCustom) RawJSON() string { return r.JSON.raw }
func (r *ProxyNewResponseConfigCustom) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Protocol to use for the proxy connection.
type ProxyNewResponseProtocol string

const (
	ProxyNewResponseProtocolHTTP  ProxyNewResponseProtocol = "http"
	ProxyNewResponseProtocolHTTPS ProxyNewResponseProtocol = "https"
)

// Current health status of the proxy.
type ProxyNewResponseStatus string

const (
	ProxyNewResponseStatusAvailable   ProxyNewResponseStatus = "available"
	ProxyNewResponseStatusUnavailable ProxyNewResponseStatus = "unavailable"
)

// Configuration for routing traffic through a proxy.
type ProxyGetResponse struct {
	// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
	// worst: `mobile` > `residential` > `isp` > `datacenter`.
	//
	// Any of "datacenter", "isp", "residential", "mobile", "custom".
	Type ProxyGetResponseType `json:"type" api:"required"`
	ID   string               `json:"id"`
	// Hostnames that should bypass the parent proxy and connect directly.
	BypassHosts []string `json:"bypass_hosts"`
	// Configuration specific to the selected proxy `type`.
	Config ProxyGetResponseConfigUnion `json:"config"`
	// IP address that the proxy uses when making requests.
	IPAddress string `json:"ip_address"`
	// Timestamp of the last health check performed on this proxy.
	LastChecked time.Time `json:"last_checked" format:"date-time"`
	// Readable name of the proxy.
	Name string `json:"name"`
	// Protocol to use for the proxy connection.
	//
	// Any of "http", "https".
	Protocol ProxyGetResponseProtocol `json:"protocol"`
	// Current health status of the proxy.
	//
	// Any of "available", "unavailable".
	Status ProxyGetResponseStatus `json:"status"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ID          respjson.Field
		BypassHosts respjson.Field
		Config      respjson.Field
		IPAddress   respjson.Field
		LastChecked respjson.Field
		Name        respjson.Field
		Protocol    respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyGetResponse) RawJSON() string { return r.JSON.raw }
func (r *ProxyGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
// worst: `mobile` > `residential` > `isp` > `datacenter`.
type ProxyGetResponseType string

const (
	ProxyGetResponseTypeDatacenter  ProxyGetResponseType = "datacenter"
	ProxyGetResponseTypeIsp         ProxyGetResponseType = "isp"
	ProxyGetResponseTypeResidential ProxyGetResponseType = "residential"
	ProxyGetResponseTypeMobile      ProxyGetResponseType = "mobile"
	ProxyGetResponseTypeCustom      ProxyGetResponseType = "custom"
)

// ProxyGetResponseConfigUnion contains all possible properties and values from
// [ProxyGetResponseConfigDatacenter], [ProxyGetResponseConfigIsp],
// [ProxyGetResponseConfigResidential], [ProxyGetResponseConfigMobile],
// [ProxyGetResponseConfigCustom].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ProxyGetResponseConfigUnion struct {
	Country string `json:"country"`
	// This field is from variant [ProxyGetResponseConfigResidential].
	Asn  string `json:"asn"`
	City string `json:"city"`
	// This field is from variant [ProxyGetResponseConfigResidential].
	Os    string `json:"os"`
	State string `json:"state"`
	// This field is from variant [ProxyGetResponseConfigResidential].
	Zip string `json:"zip"`
	// This field is from variant [ProxyGetResponseConfigCustom].
	Host string `json:"host"`
	// This field is from variant [ProxyGetResponseConfigCustom].
	Port int64 `json:"port"`
	// This field is from variant [ProxyGetResponseConfigCustom].
	HasCaBundle bool `json:"has_ca_bundle"`
	// This field is from variant [ProxyGetResponseConfigCustom].
	HasPassword bool `json:"has_password"`
	// This field is from variant [ProxyGetResponseConfigCustom].
	Username string `json:"username"`
	JSON     struct {
		Country     respjson.Field
		Asn         respjson.Field
		City        respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		Host        respjson.Field
		Port        respjson.Field
		HasCaBundle respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		raw         string
	} `json:"-"`
}

func (u ProxyGetResponseConfigUnion) AsDatacenter() (v ProxyGetResponseConfigDatacenter) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyGetResponseConfigUnion) AsIsp() (v ProxyGetResponseConfigIsp) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyGetResponseConfigUnion) AsResidential() (v ProxyGetResponseConfigResidential) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyGetResponseConfigUnion) AsMobile() (v ProxyGetResponseConfigMobile) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyGetResponseConfigUnion) AsCustom() (v ProxyGetResponseConfigCustom) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProxyGetResponseConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ProxyGetResponseConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a datacenter proxy.
type ProxyGetResponseConfigDatacenter struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country string `json:"country"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyGetResponseConfigDatacenter) RawJSON() string { return r.JSON.raw }
func (r *ProxyGetResponseConfigDatacenter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyGetResponseConfigIsp struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country string `json:"country"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyGetResponseConfigIsp) RawJSON() string { return r.JSON.raw }
func (r *ProxyGetResponseConfigIsp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyGetResponseConfigResidential struct {
	// Autonomous system number. See https://bgp.potaroo.net/cidr/autnums.html
	Asn string `json:"asn"`
	// City name (no spaces, e.g. `sanfrancisco`). If provided, `country` must also be
	// provided.
	City string `json:"city"`
	// ISO 3166 country code.
	Country string `json:"country"`
	// Operating system of the residential device.
	//
	// Any of "windows", "macos", "android".
	//
	// Deprecated: deprecated
	Os string `json:"os"`
	// Two-letter state code.
	State string `json:"state"`
	// US ZIP code.
	Zip string `json:"zip"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Asn         respjson.Field
		City        respjson.Field
		Country     respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyGetResponseConfigResidential) RawJSON() string { return r.JSON.raw }
func (r *ProxyGetResponseConfigResidential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for mobile proxies.
type ProxyGetResponseConfigMobile struct {
	// Provider city alias. Mobile carrier routing can make observed geo vary.
	City string `json:"city"`
	// ISO 3166 country code
	Country string `json:"country"`
	// US-only state code. Mobile carrier routing can make observed geo vary.
	State string `json:"state"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		City        respjson.Field
		Country     respjson.Field
		State       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyGetResponseConfigMobile) RawJSON() string { return r.JSON.raw }
func (r *ProxyGetResponseConfigMobile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a custom proxy (e.g., private proxy server).
type ProxyGetResponseConfigCustom struct {
	// Proxy host address or IP.
	Host string `json:"host" api:"required"`
	// Proxy port.
	Port int64 `json:"port" api:"required"`
	// Whether the proxy has a custom CA bundle configured.
	HasCaBundle bool `json:"has_ca_bundle"`
	// Whether the proxy has a password.
	HasPassword bool `json:"has_password"`
	// Username for proxy authentication.
	Username string `json:"username"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Host        respjson.Field
		Port        respjson.Field
		HasCaBundle respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyGetResponseConfigCustom) RawJSON() string { return r.JSON.raw }
func (r *ProxyGetResponseConfigCustom) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Protocol to use for the proxy connection.
type ProxyGetResponseProtocol string

const (
	ProxyGetResponseProtocolHTTP  ProxyGetResponseProtocol = "http"
	ProxyGetResponseProtocolHTTPS ProxyGetResponseProtocol = "https"
)

// Current health status of the proxy.
type ProxyGetResponseStatus string

const (
	ProxyGetResponseStatusAvailable   ProxyGetResponseStatus = "available"
	ProxyGetResponseStatusUnavailable ProxyGetResponseStatus = "unavailable"
)

// Configuration for routing traffic through a proxy.
type ProxyUpdateResponse struct {
	// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
	// worst: `mobile` > `residential` > `isp` > `datacenter`.
	//
	// Any of "datacenter", "isp", "residential", "mobile", "custom".
	Type ProxyUpdateResponseType `json:"type" api:"required"`
	ID   string                  `json:"id"`
	// Hostnames that should bypass the parent proxy and connect directly.
	BypassHosts []string `json:"bypass_hosts"`
	// Configuration specific to the selected proxy `type`.
	Config ProxyUpdateResponseConfigUnion `json:"config"`
	// IP address that the proxy uses when making requests.
	IPAddress string `json:"ip_address"`
	// Timestamp of the last health check performed on this proxy.
	LastChecked time.Time `json:"last_checked" format:"date-time"`
	// Readable name of the proxy.
	Name string `json:"name"`
	// Protocol to use for the proxy connection.
	//
	// Any of "http", "https".
	Protocol ProxyUpdateResponseProtocol `json:"protocol"`
	// Current health status of the proxy.
	//
	// Any of "available", "unavailable".
	Status ProxyUpdateResponseStatus `json:"status"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ID          respjson.Field
		BypassHosts respjson.Field
		Config      respjson.Field
		IPAddress   respjson.Field
		LastChecked respjson.Field
		Name        respjson.Field
		Protocol    respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyUpdateResponse) RawJSON() string { return r.JSON.raw }
func (r *ProxyUpdateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
// worst: `mobile` > `residential` > `isp` > `datacenter`.
type ProxyUpdateResponseType string

const (
	ProxyUpdateResponseTypeDatacenter  ProxyUpdateResponseType = "datacenter"
	ProxyUpdateResponseTypeIsp         ProxyUpdateResponseType = "isp"
	ProxyUpdateResponseTypeResidential ProxyUpdateResponseType = "residential"
	ProxyUpdateResponseTypeMobile      ProxyUpdateResponseType = "mobile"
	ProxyUpdateResponseTypeCustom      ProxyUpdateResponseType = "custom"
)

// ProxyUpdateResponseConfigUnion contains all possible properties and values from
// [ProxyUpdateResponseConfigDatacenter], [ProxyUpdateResponseConfigIsp],
// [ProxyUpdateResponseConfigResidential], [ProxyUpdateResponseConfigMobile],
// [ProxyUpdateResponseConfigCustom].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ProxyUpdateResponseConfigUnion struct {
	Country string `json:"country"`
	// This field is from variant [ProxyUpdateResponseConfigResidential].
	Asn  string `json:"asn"`
	City string `json:"city"`
	// This field is from variant [ProxyUpdateResponseConfigResidential].
	Os    string `json:"os"`
	State string `json:"state"`
	// This field is from variant [ProxyUpdateResponseConfigResidential].
	Zip string `json:"zip"`
	// This field is from variant [ProxyUpdateResponseConfigCustom].
	Host string `json:"host"`
	// This field is from variant [ProxyUpdateResponseConfigCustom].
	Port int64 `json:"port"`
	// This field is from variant [ProxyUpdateResponseConfigCustom].
	HasCaBundle bool `json:"has_ca_bundle"`
	// This field is from variant [ProxyUpdateResponseConfigCustom].
	HasPassword bool `json:"has_password"`
	// This field is from variant [ProxyUpdateResponseConfigCustom].
	Username string `json:"username"`
	JSON     struct {
		Country     respjson.Field
		Asn         respjson.Field
		City        respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		Host        respjson.Field
		Port        respjson.Field
		HasCaBundle respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		raw         string
	} `json:"-"`
}

func (u ProxyUpdateResponseConfigUnion) AsDatacenter() (v ProxyUpdateResponseConfigDatacenter) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyUpdateResponseConfigUnion) AsIsp() (v ProxyUpdateResponseConfigIsp) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyUpdateResponseConfigUnion) AsResidential() (v ProxyUpdateResponseConfigResidential) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyUpdateResponseConfigUnion) AsMobile() (v ProxyUpdateResponseConfigMobile) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyUpdateResponseConfigUnion) AsCustom() (v ProxyUpdateResponseConfigCustom) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProxyUpdateResponseConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ProxyUpdateResponseConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a datacenter proxy.
type ProxyUpdateResponseConfigDatacenter struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country string `json:"country"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyUpdateResponseConfigDatacenter) RawJSON() string { return r.JSON.raw }
func (r *ProxyUpdateResponseConfigDatacenter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyUpdateResponseConfigIsp struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country string `json:"country"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyUpdateResponseConfigIsp) RawJSON() string { return r.JSON.raw }
func (r *ProxyUpdateResponseConfigIsp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyUpdateResponseConfigResidential struct {
	// Autonomous system number. See https://bgp.potaroo.net/cidr/autnums.html
	Asn string `json:"asn"`
	// City name (no spaces, e.g. `sanfrancisco`). If provided, `country` must also be
	// provided.
	City string `json:"city"`
	// ISO 3166 country code.
	Country string `json:"country"`
	// Operating system of the residential device.
	//
	// Any of "windows", "macos", "android".
	//
	// Deprecated: deprecated
	Os string `json:"os"`
	// Two-letter state code.
	State string `json:"state"`
	// US ZIP code.
	Zip string `json:"zip"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Asn         respjson.Field
		City        respjson.Field
		Country     respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyUpdateResponseConfigResidential) RawJSON() string { return r.JSON.raw }
func (r *ProxyUpdateResponseConfigResidential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for mobile proxies.
type ProxyUpdateResponseConfigMobile struct {
	// Provider city alias. Mobile carrier routing can make observed geo vary.
	City string `json:"city"`
	// ISO 3166 country code
	Country string `json:"country"`
	// US-only state code. Mobile carrier routing can make observed geo vary.
	State string `json:"state"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		City        respjson.Field
		Country     respjson.Field
		State       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyUpdateResponseConfigMobile) RawJSON() string { return r.JSON.raw }
func (r *ProxyUpdateResponseConfigMobile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a custom proxy (e.g., private proxy server).
type ProxyUpdateResponseConfigCustom struct {
	// Proxy host address or IP.
	Host string `json:"host" api:"required"`
	// Proxy port.
	Port int64 `json:"port" api:"required"`
	// Whether the proxy has a custom CA bundle configured.
	HasCaBundle bool `json:"has_ca_bundle"`
	// Whether the proxy has a password.
	HasPassword bool `json:"has_password"`
	// Username for proxy authentication.
	Username string `json:"username"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Host        respjson.Field
		Port        respjson.Field
		HasCaBundle respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyUpdateResponseConfigCustom) RawJSON() string { return r.JSON.raw }
func (r *ProxyUpdateResponseConfigCustom) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Protocol to use for the proxy connection.
type ProxyUpdateResponseProtocol string

const (
	ProxyUpdateResponseProtocolHTTP  ProxyUpdateResponseProtocol = "http"
	ProxyUpdateResponseProtocolHTTPS ProxyUpdateResponseProtocol = "https"
)

// Current health status of the proxy.
type ProxyUpdateResponseStatus string

const (
	ProxyUpdateResponseStatusAvailable   ProxyUpdateResponseStatus = "available"
	ProxyUpdateResponseStatusUnavailable ProxyUpdateResponseStatus = "unavailable"
)

// Configuration for routing traffic through a proxy.
type ProxyListResponse struct {
	// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
	// worst: `mobile` > `residential` > `isp` > `datacenter`.
	//
	// Any of "datacenter", "isp", "residential", "mobile", "custom".
	Type ProxyListResponseType `json:"type" api:"required"`
	ID   string                `json:"id"`
	// Hostnames that should bypass the parent proxy and connect directly.
	BypassHosts []string `json:"bypass_hosts"`
	// Configuration specific to the selected proxy `type`.
	Config ProxyListResponseConfigUnion `json:"config"`
	// IP address that the proxy uses when making requests.
	IPAddress string `json:"ip_address"`
	// Timestamp of the last health check performed on this proxy.
	LastChecked time.Time `json:"last_checked" format:"date-time"`
	// Readable name of the proxy.
	Name string `json:"name"`
	// Protocol to use for the proxy connection.
	//
	// Any of "http", "https".
	Protocol ProxyListResponseProtocol `json:"protocol"`
	// Current health status of the proxy.
	//
	// Any of "available", "unavailable".
	Status ProxyListResponseStatus `json:"status"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ID          respjson.Field
		BypassHosts respjson.Field
		Config      respjson.Field
		IPAddress   respjson.Field
		LastChecked respjson.Field
		Name        respjson.Field
		Protocol    respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyListResponse) RawJSON() string { return r.JSON.raw }
func (r *ProxyListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
// worst: `mobile` > `residential` > `isp` > `datacenter`.
type ProxyListResponseType string

const (
	ProxyListResponseTypeDatacenter  ProxyListResponseType = "datacenter"
	ProxyListResponseTypeIsp         ProxyListResponseType = "isp"
	ProxyListResponseTypeResidential ProxyListResponseType = "residential"
	ProxyListResponseTypeMobile      ProxyListResponseType = "mobile"
	ProxyListResponseTypeCustom      ProxyListResponseType = "custom"
)

// ProxyListResponseConfigUnion contains all possible properties and values from
// [ProxyListResponseConfigDatacenter], [ProxyListResponseConfigIsp],
// [ProxyListResponseConfigResidential], [ProxyListResponseConfigMobile],
// [ProxyListResponseConfigCustom].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ProxyListResponseConfigUnion struct {
	Country string `json:"country"`
	// This field is from variant [ProxyListResponseConfigResidential].
	Asn  string `json:"asn"`
	City string `json:"city"`
	// This field is from variant [ProxyListResponseConfigResidential].
	Os    string `json:"os"`
	State string `json:"state"`
	// This field is from variant [ProxyListResponseConfigResidential].
	Zip string `json:"zip"`
	// This field is from variant [ProxyListResponseConfigCustom].
	Host string `json:"host"`
	// This field is from variant [ProxyListResponseConfigCustom].
	Port int64 `json:"port"`
	// This field is from variant [ProxyListResponseConfigCustom].
	HasCaBundle bool `json:"has_ca_bundle"`
	// This field is from variant [ProxyListResponseConfigCustom].
	HasPassword bool `json:"has_password"`
	// This field is from variant [ProxyListResponseConfigCustom].
	Username string `json:"username"`
	JSON     struct {
		Country     respjson.Field
		Asn         respjson.Field
		City        respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		Host        respjson.Field
		Port        respjson.Field
		HasCaBundle respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		raw         string
	} `json:"-"`
}

func (u ProxyListResponseConfigUnion) AsDatacenter() (v ProxyListResponseConfigDatacenter) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyListResponseConfigUnion) AsIsp() (v ProxyListResponseConfigIsp) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyListResponseConfigUnion) AsResidential() (v ProxyListResponseConfigResidential) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyListResponseConfigUnion) AsMobile() (v ProxyListResponseConfigMobile) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyListResponseConfigUnion) AsCustom() (v ProxyListResponseConfigCustom) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProxyListResponseConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ProxyListResponseConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a datacenter proxy.
type ProxyListResponseConfigDatacenter struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country string `json:"country"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyListResponseConfigDatacenter) RawJSON() string { return r.JSON.raw }
func (r *ProxyListResponseConfigDatacenter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyListResponseConfigIsp struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country string `json:"country"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyListResponseConfigIsp) RawJSON() string { return r.JSON.raw }
func (r *ProxyListResponseConfigIsp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyListResponseConfigResidential struct {
	// Autonomous system number. See https://bgp.potaroo.net/cidr/autnums.html
	Asn string `json:"asn"`
	// City name (no spaces, e.g. `sanfrancisco`). If provided, `country` must also be
	// provided.
	City string `json:"city"`
	// ISO 3166 country code.
	Country string `json:"country"`
	// Operating system of the residential device.
	//
	// Any of "windows", "macos", "android".
	//
	// Deprecated: deprecated
	Os string `json:"os"`
	// Two-letter state code.
	State string `json:"state"`
	// US ZIP code.
	Zip string `json:"zip"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Asn         respjson.Field
		City        respjson.Field
		Country     respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyListResponseConfigResidential) RawJSON() string { return r.JSON.raw }
func (r *ProxyListResponseConfigResidential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for mobile proxies.
type ProxyListResponseConfigMobile struct {
	// Provider city alias. Mobile carrier routing can make observed geo vary.
	City string `json:"city"`
	// ISO 3166 country code
	Country string `json:"country"`
	// US-only state code. Mobile carrier routing can make observed geo vary.
	State string `json:"state"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		City        respjson.Field
		Country     respjson.Field
		State       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyListResponseConfigMobile) RawJSON() string { return r.JSON.raw }
func (r *ProxyListResponseConfigMobile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a custom proxy (e.g., private proxy server).
type ProxyListResponseConfigCustom struct {
	// Proxy host address or IP.
	Host string `json:"host" api:"required"`
	// Proxy port.
	Port int64 `json:"port" api:"required"`
	// Whether the proxy has a custom CA bundle configured.
	HasCaBundle bool `json:"has_ca_bundle"`
	// Whether the proxy has a password.
	HasPassword bool `json:"has_password"`
	// Username for proxy authentication.
	Username string `json:"username"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Host        respjson.Field
		Port        respjson.Field
		HasCaBundle respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyListResponseConfigCustom) RawJSON() string { return r.JSON.raw }
func (r *ProxyListResponseConfigCustom) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Protocol to use for the proxy connection.
type ProxyListResponseProtocol string

const (
	ProxyListResponseProtocolHTTP  ProxyListResponseProtocol = "http"
	ProxyListResponseProtocolHTTPS ProxyListResponseProtocol = "https"
)

// Current health status of the proxy.
type ProxyListResponseStatus string

const (
	ProxyListResponseStatusAvailable   ProxyListResponseStatus = "available"
	ProxyListResponseStatusUnavailable ProxyListResponseStatus = "unavailable"
)

// Configuration for routing traffic through a proxy.
type ProxyCheckResponse struct {
	// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
	// worst: `mobile` > `residential` > `isp` > `datacenter`.
	//
	// Any of "datacenter", "isp", "residential", "mobile", "custom".
	Type ProxyCheckResponseType `json:"type" api:"required"`
	ID   string                 `json:"id"`
	// Hostnames that should bypass the parent proxy and connect directly.
	BypassHosts []string `json:"bypass_hosts"`
	// Configuration specific to the selected proxy `type`.
	Config ProxyCheckResponseConfigUnion `json:"config"`
	// IP address that the proxy uses when making requests.
	IPAddress string `json:"ip_address"`
	// Timestamp of the last health check performed on this proxy.
	LastChecked time.Time `json:"last_checked" format:"date-time"`
	// Readable name of the proxy.
	Name string `json:"name"`
	// Protocol to use for the proxy connection.
	//
	// Any of "http", "https".
	Protocol ProxyCheckResponseProtocol `json:"protocol"`
	// Current health status of the proxy.
	//
	// Any of "available", "unavailable".
	Status ProxyCheckResponseStatus `json:"status"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ID          respjson.Field
		BypassHosts respjson.Field
		Config      respjson.Field
		IPAddress   respjson.Field
		LastChecked respjson.Field
		Name        respjson.Field
		Protocol    respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyCheckResponse) RawJSON() string { return r.JSON.raw }
func (r *ProxyCheckResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
// worst: `mobile` > `residential` > `isp` > `datacenter`.
type ProxyCheckResponseType string

const (
	ProxyCheckResponseTypeDatacenter  ProxyCheckResponseType = "datacenter"
	ProxyCheckResponseTypeIsp         ProxyCheckResponseType = "isp"
	ProxyCheckResponseTypeResidential ProxyCheckResponseType = "residential"
	ProxyCheckResponseTypeMobile      ProxyCheckResponseType = "mobile"
	ProxyCheckResponseTypeCustom      ProxyCheckResponseType = "custom"
)

// ProxyCheckResponseConfigUnion contains all possible properties and values from
// [ProxyCheckResponseConfigDatacenter], [ProxyCheckResponseConfigIsp],
// [ProxyCheckResponseConfigResidential], [ProxyCheckResponseConfigMobile],
// [ProxyCheckResponseConfigCustom].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ProxyCheckResponseConfigUnion struct {
	Country string `json:"country"`
	// This field is from variant [ProxyCheckResponseConfigResidential].
	Asn  string `json:"asn"`
	City string `json:"city"`
	// This field is from variant [ProxyCheckResponseConfigResidential].
	Os    string `json:"os"`
	State string `json:"state"`
	// This field is from variant [ProxyCheckResponseConfigResidential].
	Zip string `json:"zip"`
	// This field is from variant [ProxyCheckResponseConfigCustom].
	Host string `json:"host"`
	// This field is from variant [ProxyCheckResponseConfigCustom].
	Port int64 `json:"port"`
	// This field is from variant [ProxyCheckResponseConfigCustom].
	HasCaBundle bool `json:"has_ca_bundle"`
	// This field is from variant [ProxyCheckResponseConfigCustom].
	HasPassword bool `json:"has_password"`
	// This field is from variant [ProxyCheckResponseConfigCustom].
	Username string `json:"username"`
	JSON     struct {
		Country     respjson.Field
		Asn         respjson.Field
		City        respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		Host        respjson.Field
		Port        respjson.Field
		HasCaBundle respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		raw         string
	} `json:"-"`
}

func (u ProxyCheckResponseConfigUnion) AsDatacenter() (v ProxyCheckResponseConfigDatacenter) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyCheckResponseConfigUnion) AsIsp() (v ProxyCheckResponseConfigIsp) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyCheckResponseConfigUnion) AsResidential() (v ProxyCheckResponseConfigResidential) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyCheckResponseConfigUnion) AsMobile() (v ProxyCheckResponseConfigMobile) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyCheckResponseConfigUnion) AsCustom() (v ProxyCheckResponseConfigCustom) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProxyCheckResponseConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ProxyCheckResponseConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a datacenter proxy.
type ProxyCheckResponseConfigDatacenter struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country string `json:"country"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyCheckResponseConfigDatacenter) RawJSON() string { return r.JSON.raw }
func (r *ProxyCheckResponseConfigDatacenter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyCheckResponseConfigIsp struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country string `json:"country"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyCheckResponseConfigIsp) RawJSON() string { return r.JSON.raw }
func (r *ProxyCheckResponseConfigIsp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyCheckResponseConfigResidential struct {
	// Autonomous system number. See https://bgp.potaroo.net/cidr/autnums.html
	Asn string `json:"asn"`
	// City name (no spaces, e.g. `sanfrancisco`). If provided, `country` must also be
	// provided.
	City string `json:"city"`
	// ISO 3166 country code.
	Country string `json:"country"`
	// Operating system of the residential device.
	//
	// Any of "windows", "macos", "android".
	//
	// Deprecated: deprecated
	Os string `json:"os"`
	// Two-letter state code.
	State string `json:"state"`
	// US ZIP code.
	Zip string `json:"zip"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Asn         respjson.Field
		City        respjson.Field
		Country     respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyCheckResponseConfigResidential) RawJSON() string { return r.JSON.raw }
func (r *ProxyCheckResponseConfigResidential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for mobile proxies.
type ProxyCheckResponseConfigMobile struct {
	// Provider city alias. Mobile carrier routing can make observed geo vary.
	City string `json:"city"`
	// ISO 3166 country code
	Country string `json:"country"`
	// US-only state code. Mobile carrier routing can make observed geo vary.
	State string `json:"state"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		City        respjson.Field
		Country     respjson.Field
		State       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyCheckResponseConfigMobile) RawJSON() string { return r.JSON.raw }
func (r *ProxyCheckResponseConfigMobile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a custom proxy (e.g., private proxy server).
type ProxyCheckResponseConfigCustom struct {
	// Proxy host address or IP.
	Host string `json:"host" api:"required"`
	// Proxy port.
	Port int64 `json:"port" api:"required"`
	// Whether the proxy has a custom CA bundle configured.
	HasCaBundle bool `json:"has_ca_bundle"`
	// Whether the proxy has a password.
	HasPassword bool `json:"has_password"`
	// Username for proxy authentication.
	Username string `json:"username"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Host        respjson.Field
		Port        respjson.Field
		HasCaBundle respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyCheckResponseConfigCustom) RawJSON() string { return r.JSON.raw }
func (r *ProxyCheckResponseConfigCustom) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Protocol to use for the proxy connection.
type ProxyCheckResponseProtocol string

const (
	ProxyCheckResponseProtocolHTTP  ProxyCheckResponseProtocol = "http"
	ProxyCheckResponseProtocolHTTPS ProxyCheckResponseProtocol = "https"
)

// Current health status of the proxy.
type ProxyCheckResponseStatus string

const (
	ProxyCheckResponseStatusAvailable   ProxyCheckResponseStatus = "available"
	ProxyCheckResponseStatusUnavailable ProxyCheckResponseStatus = "unavailable"
)

type ProxyNewParams struct {
	// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
	// worst: `mobile` > `residential` > `isp` > `datacenter`.
	//
	// Any of "datacenter", "isp", "residential", "mobile", "custom".
	Type ProxyNewParamsType `json:"type,omitzero" api:"required"`
	// Readable name of the proxy.
	Name param.Opt[string] `json:"name,omitzero"`
	// Hostnames that should bypass the parent proxy and connect directly.
	BypassHosts []string `json:"bypass_hosts,omitzero"`
	// Configuration specific to the selected proxy `type`.
	Config ProxyNewParamsConfigUnion `json:"config,omitzero"`
	// Protocol to use for the proxy connection.
	//
	// Any of "http", "https".
	Protocol ProxyNewParamsProtocol `json:"protocol,omitzero"`
	paramObj
}

func (r ProxyNewParams) MarshalJSON() (data []byte, err error) {
	type shadow ProxyNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
// worst: `mobile` > `residential` > `isp` > `datacenter`.
type ProxyNewParamsType string

const (
	ProxyNewParamsTypeDatacenter  ProxyNewParamsType = "datacenter"
	ProxyNewParamsTypeIsp         ProxyNewParamsType = "isp"
	ProxyNewParamsTypeResidential ProxyNewParamsType = "residential"
	ProxyNewParamsTypeMobile      ProxyNewParamsType = "mobile"
	ProxyNewParamsTypeCustom      ProxyNewParamsType = "custom"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ProxyNewParamsConfigUnion struct {
	OfDatacenter  *ProxyNewParamsConfigDatacenter  `json:",omitzero,inline"`
	OfIsp         *ProxyNewParamsConfigIsp         `json:",omitzero,inline"`
	OfResidential *ProxyNewParamsConfigResidential `json:",omitzero,inline"`
	OfMobile      *ProxyNewParamsConfigMobile      `json:",omitzero,inline"`
	OfCustom      *ProxyNewParamsConfigCustom      `json:",omitzero,inline"`
	paramUnion
}

func (u ProxyNewParamsConfigUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDatacenter,
		u.OfIsp,
		u.OfResidential,
		u.OfMobile,
		u.OfCustom)
}
func (u *ProxyNewParamsConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ProxyNewParamsConfigUnion) asAny() any {
	if !param.IsOmitted(u.OfDatacenter) {
		return u.OfDatacenter
	} else if !param.IsOmitted(u.OfIsp) {
		return u.OfIsp
	} else if !param.IsOmitted(u.OfResidential) {
		return u.OfResidential
	} else if !param.IsOmitted(u.OfMobile) {
		return u.OfMobile
	} else if !param.IsOmitted(u.OfCustom) {
		return u.OfCustom
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetAsn() *string {
	if vt := u.OfResidential; vt != nil && vt.Asn.Valid() {
		return &vt.Asn.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetOs() *string {
	if vt := u.OfResidential; vt != nil {
		return &vt.Os
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetZip() *string {
	if vt := u.OfResidential; vt != nil && vt.Zip.Valid() {
		return &vt.Zip.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetHost() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Host
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetPort() *int64 {
	if vt := u.OfCustom; vt != nil {
		return &vt.Port
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetCaBundle() *string {
	if vt := u.OfCustom; vt != nil && vt.CaBundle.Valid() {
		return &vt.CaBundle.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetPassword() *string {
	if vt := u.OfCustom; vt != nil && vt.Password.Valid() {
		return &vt.Password.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetUsername() *string {
	if vt := u.OfCustom; vt != nil && vt.Username.Valid() {
		return &vt.Username.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetCountry() *string {
	if vt := u.OfDatacenter; vt != nil && vt.Country.Valid() {
		return &vt.Country.Value
	} else if vt := u.OfIsp; vt != nil && vt.Country.Valid() {
		return &vt.Country.Value
	} else if vt := u.OfResidential; vt != nil && vt.Country.Valid() {
		return &vt.Country.Value
	} else if vt := u.OfMobile; vt != nil && vt.Country.Valid() {
		return &vt.Country.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetCity() *string {
	if vt := u.OfResidential; vt != nil && vt.City.Valid() {
		return &vt.City.Value
	} else if vt := u.OfMobile; vt != nil && vt.City.Valid() {
		return &vt.City.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetState() *string {
	if vt := u.OfResidential; vt != nil && vt.State.Valid() {
		return &vt.State.Value
	} else if vt := u.OfMobile; vt != nil && vt.State.Valid() {
		return &vt.State.Value
	}
	return nil
}

// Configuration for a datacenter proxy.
type ProxyNewParamsConfigDatacenter struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country param.Opt[string] `json:"country,omitzero"`
	paramObj
}

func (r ProxyNewParamsConfigDatacenter) MarshalJSON() (data []byte, err error) {
	type shadow ProxyNewParamsConfigDatacenter
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyNewParamsConfigDatacenter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyNewParamsConfigIsp struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country param.Opt[string] `json:"country,omitzero"`
	paramObj
}

func (r ProxyNewParamsConfigIsp) MarshalJSON() (data []byte, err error) {
	type shadow ProxyNewParamsConfigIsp
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyNewParamsConfigIsp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyNewParamsConfigResidential struct {
	// Autonomous system number. See https://bgp.potaroo.net/cidr/autnums.html
	Asn param.Opt[string] `json:"asn,omitzero"`
	// City name (no spaces, e.g. `sanfrancisco`). If provided, `country` must also be
	// provided.
	City param.Opt[string] `json:"city,omitzero"`
	// ISO 3166 country code.
	Country param.Opt[string] `json:"country,omitzero"`
	// Two-letter state code.
	State param.Opt[string] `json:"state,omitzero"`
	// US ZIP code.
	Zip param.Opt[string] `json:"zip,omitzero"`
	// Operating system of the residential device.
	//
	// Any of "windows", "macos", "android".
	//
	// Deprecated: deprecated
	Os string `json:"os,omitzero"`
	paramObj
}

func (r ProxyNewParamsConfigResidential) MarshalJSON() (data []byte, err error) {
	type shadow ProxyNewParamsConfigResidential
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyNewParamsConfigResidential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProxyNewParamsConfigResidential](
		"os", "windows", "macos", "android",
	)
}

// Configuration for mobile proxies.
type ProxyNewParamsConfigMobile struct {
	// Provider city alias. Mobile carrier routing can make observed geo vary.
	City param.Opt[string] `json:"city,omitzero"`
	// ISO 3166 country code
	Country param.Opt[string] `json:"country,omitzero"`
	// US-only state code. Mobile carrier routing can make observed geo vary.
	State param.Opt[string] `json:"state,omitzero"`
	paramObj
}

func (r ProxyNewParamsConfigMobile) MarshalJSON() (data []byte, err error) {
	type shadow ProxyNewParamsConfigMobile
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyNewParamsConfigMobile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a custom proxy (e.g., private proxy server).
//
// The properties Host, Port are required.
type ProxyNewParamsConfigCustom struct {
	// Proxy host address or IP.
	Host string `json:"host" api:"required"`
	// Proxy port.
	Port int64 `json:"port" api:"required"`
	// PEM-encoded CA certificate bundle the proxy re-signs upstream TLS with. Provide
	// when the proxy terminates TLS (MITM) so the browser trusts its certificates. May
	// contain multiple concatenated certificates.
	CaBundle param.Opt[string] `json:"ca_bundle,omitzero"`
	// Password for proxy authentication.
	Password param.Opt[string] `json:"password,omitzero"`
	// Username for proxy authentication.
	Username param.Opt[string] `json:"username,omitzero"`
	paramObj
}

func (r ProxyNewParamsConfigCustom) MarshalJSON() (data []byte, err error) {
	type shadow ProxyNewParamsConfigCustom
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyNewParamsConfigCustom) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Protocol to use for the proxy connection.
type ProxyNewParamsProtocol string

const (
	ProxyNewParamsProtocolHTTP  ProxyNewParamsProtocol = "http"
	ProxyNewParamsProtocolHTTPS ProxyNewParamsProtocol = "https"
)

type ProxyUpdateParams struct {
	// New proxy name. Proxy names are trimmed and length-checked only; duplicates are
	// allowed because proxies are updated by ID, not by name.
	Name string `json:"name" api:"required"`
	paramObj
}

func (r ProxyUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow ProxyUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProxyListParams struct {
	// Limit the number of proxies to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Exact-match filter on proxy name using the database collation. In production,
	// matching is case- and accent-insensitive. Names are not required to be unique,
	// so multiple proxies may match.
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	// Offset the number of proxies to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Case-insensitive substring match against proxy name, host, or IP address. IDs
	// match by exact value.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ProxyListParams]'s query parameters as `url.Values`.
func (r ProxyListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type ProxyCheckParams struct {
	// An optional URL to test reachability against. If provided, the proxy check will
	// test connectivity to this URL instead of the default test URLs. Only HTTP and
	// HTTPS schemes are allowed, and the URL must resolve to a public IP address. For
	// ISP and datacenter proxies, the exit IP is stable, so a successful check
	// reliably indicates that subsequent browser sessions will reach the target site
	// with the same IP. For residential and mobile proxies, the exit node changes
	// between requests, so a successful check validates proxy configuration but does
	// not guarantee that a subsequent browser session will use the same exit IP or
	// reach the same site — it is useful for verifying credentials and connectivity,
	// not for predicting site-specific behavior. When provided, the check result does
	// not update the proxy's health status, since a failure may indicate a problem
	// with the target site rather than the proxy itself.
	URL param.Opt[string] `json:"url,omitzero"`
	paramObj
}

func (r ProxyCheckParams) MarshalJSON() (data []byte, err error) {
	type shadow ProxyCheckParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyCheckParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
