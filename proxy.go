// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/onkernel/kernel-go-sdk/internal/apijson"
	"github.com/onkernel/kernel-go-sdk/internal/requestconfig"
	"github.com/onkernel/kernel-go-sdk/option"
	"github.com/onkernel/kernel-go-sdk/packages/param"
	"github.com/onkernel/kernel-go-sdk/packages/respjson"
)

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

// Create a new proxy configuration for the caller's organization.
func (r *ProxyService) New(ctx context.Context, body ProxyNewParams, opts ...option.RequestOption) (res *ProxyNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "proxies"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieve a proxy belonging to the caller's organization by ID.
func (r *ProxyService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *ProxyGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("proxies/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List proxies owned by the caller's organization.
func (r *ProxyService) List(ctx context.Context, opts ...option.RequestOption) (res *[]ProxyListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "proxies"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Soft delete a proxy. Sessions referencing it are not modified.
func (r *ProxyService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("proxies/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return
}

// Run a health check on the proxy to verify it's working.
func (r *ProxyService) Check(ctx context.Context, id string, opts ...option.RequestOption) (res *ProxyCheckResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("proxies/%s/check", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Configuration for routing traffic through a proxy.
type ProxyNewResponse struct {
	// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
	// worst: `mobile` > `residential` > `isp` > `datacenter`.
	//
	// Any of "datacenter", "isp", "residential", "mobile", "custom".
	Type ProxyNewResponseType `json:"type,required"`
	ID   string               `json:"id"`
	// Configuration specific to the selected proxy `type`.
	Config ProxyNewResponseConfigUnion `json:"config"`
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
		Config      respjson.Field
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
// [ProxyNewResponseConfigDatacenterProxyConfig],
// [ProxyNewResponseConfigIspProxyConfig],
// [ProxyNewResponseConfigResidentialProxyConfig],
// [ProxyNewResponseConfigMobileProxyConfig],
// [ProxyNewResponseConfigCustomProxyConfig].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ProxyNewResponseConfigUnion struct {
	Country string `json:"country"`
	Asn     string `json:"asn"`
	City    string `json:"city"`
	// This field is from variant [ProxyNewResponseConfigResidentialProxyConfig].
	Os    string `json:"os"`
	State string `json:"state"`
	Zip   string `json:"zip"`
	// This field is from variant [ProxyNewResponseConfigMobileProxyConfig].
	Carrier string `json:"carrier"`
	// This field is from variant [ProxyNewResponseConfigCustomProxyConfig].
	Host string `json:"host"`
	// This field is from variant [ProxyNewResponseConfigCustomProxyConfig].
	Port int64 `json:"port"`
	// This field is from variant [ProxyNewResponseConfigCustomProxyConfig].
	HasPassword bool `json:"has_password"`
	// This field is from variant [ProxyNewResponseConfigCustomProxyConfig].
	Username string `json:"username"`
	JSON     struct {
		Country     respjson.Field
		Asn         respjson.Field
		City        respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		Carrier     respjson.Field
		Host        respjson.Field
		Port        respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		raw         string
	} `json:"-"`
}

func (u ProxyNewResponseConfigUnion) AsProxyNewResponseConfigDatacenterProxyConfig() (v ProxyNewResponseConfigDatacenterProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyNewResponseConfigUnion) AsProxyNewResponseConfigIspProxyConfig() (v ProxyNewResponseConfigIspProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyNewResponseConfigUnion) AsProxyNewResponseConfigResidentialProxyConfig() (v ProxyNewResponseConfigResidentialProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyNewResponseConfigUnion) AsProxyNewResponseConfigMobileProxyConfig() (v ProxyNewResponseConfigMobileProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyNewResponseConfigUnion) AsProxyNewResponseConfigCustomProxyConfig() (v ProxyNewResponseConfigCustomProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProxyNewResponseConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ProxyNewResponseConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a datacenter proxy.
type ProxyNewResponseConfigDatacenterProxyConfig struct {
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
func (r ProxyNewResponseConfigDatacenterProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyNewResponseConfigDatacenterProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyNewResponseConfigIspProxyConfig struct {
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
func (r ProxyNewResponseConfigIspProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyNewResponseConfigIspProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyNewResponseConfigResidentialProxyConfig struct {
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
func (r ProxyNewResponseConfigResidentialProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyNewResponseConfigResidentialProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for mobile proxies.
type ProxyNewResponseConfigMobileProxyConfig struct {
	// Autonomous system number. See https://bgp.potaroo.net/cidr/autnums.html
	Asn string `json:"asn"`
	// Mobile carrier.
	//
	// Any of "a1", "aircel", "airtel", "att", "celcom", "chinamobile", "claro",
	// "comcast", "cox", "digi", "dt", "docomo", "dtac", "etisalat", "idea",
	// "kyivstar", "meo", "megafon", "mtn", "mtnza", "mts", "optus", "orange", "qwest",
	// "reliance_jio", "robi", "sprint", "telefonica", "telstra", "tmobile", "tigo",
	// "tim", "verizon", "vimpelcom", "vodacomza", "vodafone", "vivo", "zain",
	// "vivabo", "telenormyanmar", "kcelljsc", "swisscom", "singtel", "asiacell",
	// "windit", "cellc", "ooredoo", "drei", "umobile", "cableone", "proximus",
	// "tele2", "mobitel", "o2", "bouygues", "free", "sfr", "digicel".
	Carrier string `json:"carrier"`
	// City name (no spaces, e.g. `sanfrancisco`). If provided, `country` must also be
	// provided.
	City string `json:"city"`
	// ISO 3166 country code
	Country string `json:"country"`
	// Two-letter state code.
	State string `json:"state"`
	// US ZIP code.
	Zip string `json:"zip"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Asn         respjson.Field
		Carrier     respjson.Field
		City        respjson.Field
		Country     respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyNewResponseConfigMobileProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyNewResponseConfigMobileProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a custom proxy (e.g., private proxy server).
type ProxyNewResponseConfigCustomProxyConfig struct {
	// Proxy host address or IP.
	Host string `json:"host,required"`
	// Proxy port.
	Port int64 `json:"port,required"`
	// Whether the proxy has a password.
	HasPassword bool `json:"has_password"`
	// Username for proxy authentication.
	Username string `json:"username"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Host        respjson.Field
		Port        respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyNewResponseConfigCustomProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyNewResponseConfigCustomProxyConfig) UnmarshalJSON(data []byte) error {
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
	Type ProxyGetResponseType `json:"type,required"`
	ID   string               `json:"id"`
	// Configuration specific to the selected proxy `type`.
	Config ProxyGetResponseConfigUnion `json:"config"`
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
		Config      respjson.Field
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
// [ProxyGetResponseConfigDatacenterProxyConfig],
// [ProxyGetResponseConfigIspProxyConfig],
// [ProxyGetResponseConfigResidentialProxyConfig],
// [ProxyGetResponseConfigMobileProxyConfig],
// [ProxyGetResponseConfigCustomProxyConfig].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ProxyGetResponseConfigUnion struct {
	Country string `json:"country"`
	Asn     string `json:"asn"`
	City    string `json:"city"`
	// This field is from variant [ProxyGetResponseConfigResidentialProxyConfig].
	Os    string `json:"os"`
	State string `json:"state"`
	Zip   string `json:"zip"`
	// This field is from variant [ProxyGetResponseConfigMobileProxyConfig].
	Carrier string `json:"carrier"`
	// This field is from variant [ProxyGetResponseConfigCustomProxyConfig].
	Host string `json:"host"`
	// This field is from variant [ProxyGetResponseConfigCustomProxyConfig].
	Port int64 `json:"port"`
	// This field is from variant [ProxyGetResponseConfigCustomProxyConfig].
	HasPassword bool `json:"has_password"`
	// This field is from variant [ProxyGetResponseConfigCustomProxyConfig].
	Username string `json:"username"`
	JSON     struct {
		Country     respjson.Field
		Asn         respjson.Field
		City        respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		Carrier     respjson.Field
		Host        respjson.Field
		Port        respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		raw         string
	} `json:"-"`
}

func (u ProxyGetResponseConfigUnion) AsProxyGetResponseConfigDatacenterProxyConfig() (v ProxyGetResponseConfigDatacenterProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyGetResponseConfigUnion) AsProxyGetResponseConfigIspProxyConfig() (v ProxyGetResponseConfigIspProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyGetResponseConfigUnion) AsProxyGetResponseConfigResidentialProxyConfig() (v ProxyGetResponseConfigResidentialProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyGetResponseConfigUnion) AsProxyGetResponseConfigMobileProxyConfig() (v ProxyGetResponseConfigMobileProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyGetResponseConfigUnion) AsProxyGetResponseConfigCustomProxyConfig() (v ProxyGetResponseConfigCustomProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProxyGetResponseConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ProxyGetResponseConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a datacenter proxy.
type ProxyGetResponseConfigDatacenterProxyConfig struct {
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
func (r ProxyGetResponseConfigDatacenterProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyGetResponseConfigDatacenterProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyGetResponseConfigIspProxyConfig struct {
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
func (r ProxyGetResponseConfigIspProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyGetResponseConfigIspProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyGetResponseConfigResidentialProxyConfig struct {
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
func (r ProxyGetResponseConfigResidentialProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyGetResponseConfigResidentialProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for mobile proxies.
type ProxyGetResponseConfigMobileProxyConfig struct {
	// Autonomous system number. See https://bgp.potaroo.net/cidr/autnums.html
	Asn string `json:"asn"`
	// Mobile carrier.
	//
	// Any of "a1", "aircel", "airtel", "att", "celcom", "chinamobile", "claro",
	// "comcast", "cox", "digi", "dt", "docomo", "dtac", "etisalat", "idea",
	// "kyivstar", "meo", "megafon", "mtn", "mtnza", "mts", "optus", "orange", "qwest",
	// "reliance_jio", "robi", "sprint", "telefonica", "telstra", "tmobile", "tigo",
	// "tim", "verizon", "vimpelcom", "vodacomza", "vodafone", "vivo", "zain",
	// "vivabo", "telenormyanmar", "kcelljsc", "swisscom", "singtel", "asiacell",
	// "windit", "cellc", "ooredoo", "drei", "umobile", "cableone", "proximus",
	// "tele2", "mobitel", "o2", "bouygues", "free", "sfr", "digicel".
	Carrier string `json:"carrier"`
	// City name (no spaces, e.g. `sanfrancisco`). If provided, `country` must also be
	// provided.
	City string `json:"city"`
	// ISO 3166 country code
	Country string `json:"country"`
	// Two-letter state code.
	State string `json:"state"`
	// US ZIP code.
	Zip string `json:"zip"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Asn         respjson.Field
		Carrier     respjson.Field
		City        respjson.Field
		Country     respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyGetResponseConfigMobileProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyGetResponseConfigMobileProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a custom proxy (e.g., private proxy server).
type ProxyGetResponseConfigCustomProxyConfig struct {
	// Proxy host address or IP.
	Host string `json:"host,required"`
	// Proxy port.
	Port int64 `json:"port,required"`
	// Whether the proxy has a password.
	HasPassword bool `json:"has_password"`
	// Username for proxy authentication.
	Username string `json:"username"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Host        respjson.Field
		Port        respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyGetResponseConfigCustomProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyGetResponseConfigCustomProxyConfig) UnmarshalJSON(data []byte) error {
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
type ProxyListResponse struct {
	// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
	// worst: `mobile` > `residential` > `isp` > `datacenter`.
	//
	// Any of "datacenter", "isp", "residential", "mobile", "custom".
	Type ProxyListResponseType `json:"type,required"`
	ID   string                `json:"id"`
	// Configuration specific to the selected proxy `type`.
	Config ProxyListResponseConfigUnion `json:"config"`
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
		Config      respjson.Field
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
// [ProxyListResponseConfigDatacenterProxyConfig],
// [ProxyListResponseConfigIspProxyConfig],
// [ProxyListResponseConfigResidentialProxyConfig],
// [ProxyListResponseConfigMobileProxyConfig],
// [ProxyListResponseConfigCustomProxyConfig].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ProxyListResponseConfigUnion struct {
	Country string `json:"country"`
	Asn     string `json:"asn"`
	City    string `json:"city"`
	// This field is from variant [ProxyListResponseConfigResidentialProxyConfig].
	Os    string `json:"os"`
	State string `json:"state"`
	Zip   string `json:"zip"`
	// This field is from variant [ProxyListResponseConfigMobileProxyConfig].
	Carrier string `json:"carrier"`
	// This field is from variant [ProxyListResponseConfigCustomProxyConfig].
	Host string `json:"host"`
	// This field is from variant [ProxyListResponseConfigCustomProxyConfig].
	Port int64 `json:"port"`
	// This field is from variant [ProxyListResponseConfigCustomProxyConfig].
	HasPassword bool `json:"has_password"`
	// This field is from variant [ProxyListResponseConfigCustomProxyConfig].
	Username string `json:"username"`
	JSON     struct {
		Country     respjson.Field
		Asn         respjson.Field
		City        respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		Carrier     respjson.Field
		Host        respjson.Field
		Port        respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		raw         string
	} `json:"-"`
}

func (u ProxyListResponseConfigUnion) AsProxyListResponseConfigDatacenterProxyConfig() (v ProxyListResponseConfigDatacenterProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyListResponseConfigUnion) AsProxyListResponseConfigIspProxyConfig() (v ProxyListResponseConfigIspProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyListResponseConfigUnion) AsProxyListResponseConfigResidentialProxyConfig() (v ProxyListResponseConfigResidentialProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyListResponseConfigUnion) AsProxyListResponseConfigMobileProxyConfig() (v ProxyListResponseConfigMobileProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyListResponseConfigUnion) AsProxyListResponseConfigCustomProxyConfig() (v ProxyListResponseConfigCustomProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProxyListResponseConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ProxyListResponseConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a datacenter proxy.
type ProxyListResponseConfigDatacenterProxyConfig struct {
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
func (r ProxyListResponseConfigDatacenterProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyListResponseConfigDatacenterProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyListResponseConfigIspProxyConfig struct {
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
func (r ProxyListResponseConfigIspProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyListResponseConfigIspProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyListResponseConfigResidentialProxyConfig struct {
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
func (r ProxyListResponseConfigResidentialProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyListResponseConfigResidentialProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for mobile proxies.
type ProxyListResponseConfigMobileProxyConfig struct {
	// Autonomous system number. See https://bgp.potaroo.net/cidr/autnums.html
	Asn string `json:"asn"`
	// Mobile carrier.
	//
	// Any of "a1", "aircel", "airtel", "att", "celcom", "chinamobile", "claro",
	// "comcast", "cox", "digi", "dt", "docomo", "dtac", "etisalat", "idea",
	// "kyivstar", "meo", "megafon", "mtn", "mtnza", "mts", "optus", "orange", "qwest",
	// "reliance_jio", "robi", "sprint", "telefonica", "telstra", "tmobile", "tigo",
	// "tim", "verizon", "vimpelcom", "vodacomza", "vodafone", "vivo", "zain",
	// "vivabo", "telenormyanmar", "kcelljsc", "swisscom", "singtel", "asiacell",
	// "windit", "cellc", "ooredoo", "drei", "umobile", "cableone", "proximus",
	// "tele2", "mobitel", "o2", "bouygues", "free", "sfr", "digicel".
	Carrier string `json:"carrier"`
	// City name (no spaces, e.g. `sanfrancisco`). If provided, `country` must also be
	// provided.
	City string `json:"city"`
	// ISO 3166 country code
	Country string `json:"country"`
	// Two-letter state code.
	State string `json:"state"`
	// US ZIP code.
	Zip string `json:"zip"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Asn         respjson.Field
		Carrier     respjson.Field
		City        respjson.Field
		Country     respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyListResponseConfigMobileProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyListResponseConfigMobileProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a custom proxy (e.g., private proxy server).
type ProxyListResponseConfigCustomProxyConfig struct {
	// Proxy host address or IP.
	Host string `json:"host,required"`
	// Proxy port.
	Port int64 `json:"port,required"`
	// Whether the proxy has a password.
	HasPassword bool `json:"has_password"`
	// Username for proxy authentication.
	Username string `json:"username"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Host        respjson.Field
		Port        respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyListResponseConfigCustomProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyListResponseConfigCustomProxyConfig) UnmarshalJSON(data []byte) error {
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
	Type ProxyCheckResponseType `json:"type,required"`
	ID   string                 `json:"id"`
	// Configuration specific to the selected proxy `type`.
	Config ProxyCheckResponseConfigUnion `json:"config"`
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
		Config      respjson.Field
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
// [ProxyCheckResponseConfigDatacenterProxyConfig],
// [ProxyCheckResponseConfigIspProxyConfig],
// [ProxyCheckResponseConfigResidentialProxyConfig],
// [ProxyCheckResponseConfigMobileProxyConfig],
// [ProxyCheckResponseConfigCustomProxyConfig].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ProxyCheckResponseConfigUnion struct {
	Country string `json:"country"`
	Asn     string `json:"asn"`
	City    string `json:"city"`
	// This field is from variant [ProxyCheckResponseConfigResidentialProxyConfig].
	Os    string `json:"os"`
	State string `json:"state"`
	Zip   string `json:"zip"`
	// This field is from variant [ProxyCheckResponseConfigMobileProxyConfig].
	Carrier string `json:"carrier"`
	// This field is from variant [ProxyCheckResponseConfigCustomProxyConfig].
	Host string `json:"host"`
	// This field is from variant [ProxyCheckResponseConfigCustomProxyConfig].
	Port int64 `json:"port"`
	// This field is from variant [ProxyCheckResponseConfigCustomProxyConfig].
	HasPassword bool `json:"has_password"`
	// This field is from variant [ProxyCheckResponseConfigCustomProxyConfig].
	Username string `json:"username"`
	JSON     struct {
		Country     respjson.Field
		Asn         respjson.Field
		City        respjson.Field
		Os          respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		Carrier     respjson.Field
		Host        respjson.Field
		Port        respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		raw         string
	} `json:"-"`
}

func (u ProxyCheckResponseConfigUnion) AsProxyCheckResponseConfigDatacenterProxyConfig() (v ProxyCheckResponseConfigDatacenterProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyCheckResponseConfigUnion) AsProxyCheckResponseConfigIspProxyConfig() (v ProxyCheckResponseConfigIspProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyCheckResponseConfigUnion) AsProxyCheckResponseConfigResidentialProxyConfig() (v ProxyCheckResponseConfigResidentialProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyCheckResponseConfigUnion) AsProxyCheckResponseConfigMobileProxyConfig() (v ProxyCheckResponseConfigMobileProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyCheckResponseConfigUnion) AsProxyCheckResponseConfigCustomProxyConfig() (v ProxyCheckResponseConfigCustomProxyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProxyCheckResponseConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ProxyCheckResponseConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a datacenter proxy.
type ProxyCheckResponseConfigDatacenterProxyConfig struct {
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
func (r ProxyCheckResponseConfigDatacenterProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyCheckResponseConfigDatacenterProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyCheckResponseConfigIspProxyConfig struct {
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
func (r ProxyCheckResponseConfigIspProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyCheckResponseConfigIspProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyCheckResponseConfigResidentialProxyConfig struct {
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
func (r ProxyCheckResponseConfigResidentialProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyCheckResponseConfigResidentialProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for mobile proxies.
type ProxyCheckResponseConfigMobileProxyConfig struct {
	// Autonomous system number. See https://bgp.potaroo.net/cidr/autnums.html
	Asn string `json:"asn"`
	// Mobile carrier.
	//
	// Any of "a1", "aircel", "airtel", "att", "celcom", "chinamobile", "claro",
	// "comcast", "cox", "digi", "dt", "docomo", "dtac", "etisalat", "idea",
	// "kyivstar", "meo", "megafon", "mtn", "mtnza", "mts", "optus", "orange", "qwest",
	// "reliance_jio", "robi", "sprint", "telefonica", "telstra", "tmobile", "tigo",
	// "tim", "verizon", "vimpelcom", "vodacomza", "vodafone", "vivo", "zain",
	// "vivabo", "telenormyanmar", "kcelljsc", "swisscom", "singtel", "asiacell",
	// "windit", "cellc", "ooredoo", "drei", "umobile", "cableone", "proximus",
	// "tele2", "mobitel", "o2", "bouygues", "free", "sfr", "digicel".
	Carrier string `json:"carrier"`
	// City name (no spaces, e.g. `sanfrancisco`). If provided, `country` must also be
	// provided.
	City string `json:"city"`
	// ISO 3166 country code
	Country string `json:"country"`
	// Two-letter state code.
	State string `json:"state"`
	// US ZIP code.
	Zip string `json:"zip"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Asn         respjson.Field
		Carrier     respjson.Field
		City        respjson.Field
		Country     respjson.Field
		State       respjson.Field
		Zip         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyCheckResponseConfigMobileProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyCheckResponseConfigMobileProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a custom proxy (e.g., private proxy server).
type ProxyCheckResponseConfigCustomProxyConfig struct {
	// Proxy host address or IP.
	Host string `json:"host,required"`
	// Proxy port.
	Port int64 `json:"port,required"`
	// Whether the proxy has a password.
	HasPassword bool `json:"has_password"`
	// Username for proxy authentication.
	Username string `json:"username"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Host        respjson.Field
		Port        respjson.Field
		HasPassword respjson.Field
		Username    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyCheckResponseConfigCustomProxyConfig) RawJSON() string { return r.JSON.raw }
func (r *ProxyCheckResponseConfigCustomProxyConfig) UnmarshalJSON(data []byte) error {
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
	Type ProxyNewParamsType `json:"type,omitzero,required"`
	// Readable name of the proxy.
	Name param.Opt[string] `json:"name,omitzero"`
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
	OfProxyNewsConfigDatacenterProxyConfig   *ProxyNewParamsConfigDatacenterProxyConfig   `json:",omitzero,inline"`
	OfProxyNewsConfigIspProxyConfig          *ProxyNewParamsConfigIspProxyConfig          `json:",omitzero,inline"`
	OfProxyNewsConfigResidentialProxyConfig  *ProxyNewParamsConfigResidentialProxyConfig  `json:",omitzero,inline"`
	OfProxyNewsConfigMobileProxyConfig       *ProxyNewParamsConfigMobileProxyConfig       `json:",omitzero,inline"`
	OfProxyNewsConfigCreateCustomProxyConfig *ProxyNewParamsConfigCreateCustomProxyConfig `json:",omitzero,inline"`
	paramUnion
}

func (u ProxyNewParamsConfigUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfProxyNewsConfigDatacenterProxyConfig,
		u.OfProxyNewsConfigIspProxyConfig,
		u.OfProxyNewsConfigResidentialProxyConfig,
		u.OfProxyNewsConfigMobileProxyConfig,
		u.OfProxyNewsConfigCreateCustomProxyConfig)
}
func (u *ProxyNewParamsConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ProxyNewParamsConfigUnion) asAny() any {
	if !param.IsOmitted(u.OfProxyNewsConfigDatacenterProxyConfig) {
		return u.OfProxyNewsConfigDatacenterProxyConfig
	} else if !param.IsOmitted(u.OfProxyNewsConfigIspProxyConfig) {
		return u.OfProxyNewsConfigIspProxyConfig
	} else if !param.IsOmitted(u.OfProxyNewsConfigResidentialProxyConfig) {
		return u.OfProxyNewsConfigResidentialProxyConfig
	} else if !param.IsOmitted(u.OfProxyNewsConfigMobileProxyConfig) {
		return u.OfProxyNewsConfigMobileProxyConfig
	} else if !param.IsOmitted(u.OfProxyNewsConfigCreateCustomProxyConfig) {
		return u.OfProxyNewsConfigCreateCustomProxyConfig
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetOs() *string {
	if vt := u.OfProxyNewsConfigResidentialProxyConfig; vt != nil {
		return &vt.Os
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetCarrier() *string {
	if vt := u.OfProxyNewsConfigMobileProxyConfig; vt != nil {
		return &vt.Carrier
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetHost() *string {
	if vt := u.OfProxyNewsConfigCreateCustomProxyConfig; vt != nil {
		return &vt.Host
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetPort() *int64 {
	if vt := u.OfProxyNewsConfigCreateCustomProxyConfig; vt != nil {
		return &vt.Port
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetPassword() *string {
	if vt := u.OfProxyNewsConfigCreateCustomProxyConfig; vt != nil && vt.Password.Valid() {
		return &vt.Password.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetUsername() *string {
	if vt := u.OfProxyNewsConfigCreateCustomProxyConfig; vt != nil && vt.Username.Valid() {
		return &vt.Username.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetCountry() *string {
	if vt := u.OfProxyNewsConfigDatacenterProxyConfig; vt != nil && vt.Country.Valid() {
		return &vt.Country.Value
	} else if vt := u.OfProxyNewsConfigIspProxyConfig; vt != nil && vt.Country.Valid() {
		return &vt.Country.Value
	} else if vt := u.OfProxyNewsConfigResidentialProxyConfig; vt != nil && vt.Country.Valid() {
		return &vt.Country.Value
	} else if vt := u.OfProxyNewsConfigMobileProxyConfig; vt != nil && vt.Country.Valid() {
		return &vt.Country.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetAsn() *string {
	if vt := u.OfProxyNewsConfigResidentialProxyConfig; vt != nil && vt.Asn.Valid() {
		return &vt.Asn.Value
	} else if vt := u.OfProxyNewsConfigMobileProxyConfig; vt != nil && vt.Asn.Valid() {
		return &vt.Asn.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetCity() *string {
	if vt := u.OfProxyNewsConfigResidentialProxyConfig; vt != nil && vt.City.Valid() {
		return &vt.City.Value
	} else if vt := u.OfProxyNewsConfigMobileProxyConfig; vt != nil && vt.City.Valid() {
		return &vt.City.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetState() *string {
	if vt := u.OfProxyNewsConfigResidentialProxyConfig; vt != nil && vt.State.Valid() {
		return &vt.State.Value
	} else if vt := u.OfProxyNewsConfigMobileProxyConfig; vt != nil && vt.State.Valid() {
		return &vt.State.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProxyNewParamsConfigUnion) GetZip() *string {
	if vt := u.OfProxyNewsConfigResidentialProxyConfig; vt != nil && vt.Zip.Valid() {
		return &vt.Zip.Value
	} else if vt := u.OfProxyNewsConfigMobileProxyConfig; vt != nil && vt.Zip.Valid() {
		return &vt.Zip.Value
	}
	return nil
}

// Configuration for a datacenter proxy.
type ProxyNewParamsConfigDatacenterProxyConfig struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country param.Opt[string] `json:"country,omitzero"`
	paramObj
}

func (r ProxyNewParamsConfigDatacenterProxyConfig) MarshalJSON() (data []byte, err error) {
	type shadow ProxyNewParamsConfigDatacenterProxyConfig
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyNewParamsConfigDatacenterProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyNewParamsConfigIspProxyConfig struct {
	// ISO 3166 country code. Defaults to US if not provided.
	Country param.Opt[string] `json:"country,omitzero"`
	paramObj
}

func (r ProxyNewParamsConfigIspProxyConfig) MarshalJSON() (data []byte, err error) {
	type shadow ProxyNewParamsConfigIspProxyConfig
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyNewParamsConfigIspProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyNewParamsConfigResidentialProxyConfig struct {
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

func (r ProxyNewParamsConfigResidentialProxyConfig) MarshalJSON() (data []byte, err error) {
	type shadow ProxyNewParamsConfigResidentialProxyConfig
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyNewParamsConfigResidentialProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProxyNewParamsConfigResidentialProxyConfig](
		"os", "windows", "macos", "android",
	)
}

// Configuration for mobile proxies.
type ProxyNewParamsConfigMobileProxyConfig struct {
	// Autonomous system number. See https://bgp.potaroo.net/cidr/autnums.html
	Asn param.Opt[string] `json:"asn,omitzero"`
	// City name (no spaces, e.g. `sanfrancisco`). If provided, `country` must also be
	// provided.
	City param.Opt[string] `json:"city,omitzero"`
	// ISO 3166 country code
	Country param.Opt[string] `json:"country,omitzero"`
	// Two-letter state code.
	State param.Opt[string] `json:"state,omitzero"`
	// US ZIP code.
	Zip param.Opt[string] `json:"zip,omitzero"`
	// Mobile carrier.
	//
	// Any of "a1", "aircel", "airtel", "att", "celcom", "chinamobile", "claro",
	// "comcast", "cox", "digi", "dt", "docomo", "dtac", "etisalat", "idea",
	// "kyivstar", "meo", "megafon", "mtn", "mtnza", "mts", "optus", "orange", "qwest",
	// "reliance_jio", "robi", "sprint", "telefonica", "telstra", "tmobile", "tigo",
	// "tim", "verizon", "vimpelcom", "vodacomza", "vodafone", "vivo", "zain",
	// "vivabo", "telenormyanmar", "kcelljsc", "swisscom", "singtel", "asiacell",
	// "windit", "cellc", "ooredoo", "drei", "umobile", "cableone", "proximus",
	// "tele2", "mobitel", "o2", "bouygues", "free", "sfr", "digicel".
	Carrier string `json:"carrier,omitzero"`
	paramObj
}

func (r ProxyNewParamsConfigMobileProxyConfig) MarshalJSON() (data []byte, err error) {
	type shadow ProxyNewParamsConfigMobileProxyConfig
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyNewParamsConfigMobileProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProxyNewParamsConfigMobileProxyConfig](
		"carrier", "a1", "aircel", "airtel", "att", "celcom", "chinamobile", "claro", "comcast", "cox", "digi", "dt", "docomo", "dtac", "etisalat", "idea", "kyivstar", "meo", "megafon", "mtn", "mtnza", "mts", "optus", "orange", "qwest", "reliance_jio", "robi", "sprint", "telefonica", "telstra", "tmobile", "tigo", "tim", "verizon", "vimpelcom", "vodacomza", "vodafone", "vivo", "zain", "vivabo", "telenormyanmar", "kcelljsc", "swisscom", "singtel", "asiacell", "windit", "cellc", "ooredoo", "drei", "umobile", "cableone", "proximus", "tele2", "mobitel", "o2", "bouygues", "free", "sfr", "digicel",
	)
}

// Configuration for a custom proxy (e.g., private proxy server).
//
// The properties Host, Port are required.
type ProxyNewParamsConfigCreateCustomProxyConfig struct {
	// Proxy host address or IP.
	Host string `json:"host,required"`
	// Proxy port.
	Port int64 `json:"port,required"`
	// Password for proxy authentication.
	Password param.Opt[string] `json:"password,omitzero"`
	// Username for proxy authentication.
	Username param.Opt[string] `json:"username,omitzero"`
	paramObj
}

func (r ProxyNewParamsConfigCreateCustomProxyConfig) MarshalJSON() (data []byte, err error) {
	type shadow ProxyNewParamsConfigCreateCustomProxyConfig
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProxyNewParamsConfigCreateCustomProxyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Protocol to use for the proxy connection.
type ProxyNewParamsProtocol string

const (
	ProxyNewParamsProtocolHTTP  ProxyNewParamsProtocol = "http"
	ProxyNewParamsProtocolHTTPS ProxyNewParamsProtocol = "https"
)
