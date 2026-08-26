// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
	"github.com/kernel/kernel-go-sdk/shared"
	"github.com/kernel/kernel-go-sdk/shared/constant"
)

// Resolve browser and proxy recommendations for bot-protected sites.
//
// ConfigRegistryService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewConfigRegistryService] method instead.
type ConfigRegistryService struct {
	Options []option.RequestOption
	// Resolve browser and proxy recommendations for bot-protected sites.
	Analyses ConfigRegistryAnalysisService
}

// NewConfigRegistryService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewConfigRegistryService(opts ...option.RequestOption) (r ConfigRegistryService) {
	r = ConfigRegistryService{}
	r.Options = opts
	r.Analyses = NewConfigRegistryAnalysisService(opts...)
	return
}

// Lists unique domains previously analyzed by the selected project with their
// current domain-level recommendations.
func (r *ConfigRegistryService) List(ctx context.Context, query ConfigRegistryListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[RecommendationSummary], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "config-registry"
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

// Lists unique domains previously analyzed by the selected project with their
// current domain-level recommendations.
func (r *ConfigRegistryService) ListAutoPaging(ctx context.Context, query ConfigRegistryListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[RecommendationSummary] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Returns current global knowledge without resolving DNS, creating an analysis, or
// updating config registry data.
func (r *ConfigRegistryService) Lookup(ctx context.Context, body ConfigRegistryLookupParams, opts ...option.RequestOption) (res *LookupResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "config-registry/lookup"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Explicitly starts or retries a project-scoped background analysis while
// preserving current global knowledge when available. Use
// `/config-registry/lookup` for side-effect-free reads.
func (r *ConfigRegistryService) Resolve(ctx context.Context, body ConfigRegistryResolveParams, opts ...option.RequestOption) (res *ConfigRegistryResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "config-registry/resolve"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type Analysis struct {
	// Discovery run ID used to poll analysis status.
	ID string `json:"id" api:"required"`
	// Time the analysis was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Present for failed or canceled analyses. Messages contain safe retry guidance
	// rather than internal workflow errors.
	Failure shared.ErrorModel `json:"failure" api:"required"`
	// Time the analysis reached a terminal status. Null while it is running.
	FinishedAt time.Time `json:"finished_at" api:"required" format:"date-time"`
	// Lifecycle status of a background analysis.
	//
	// Any of "running", "completed", "failed", "canceled".
	Status AnalysisStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Failure     respjson.Field
		FinishedAt  respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Analysis) RawJSON() string { return r.JSON.raw }
func (r *Analysis) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Lifecycle status of a background analysis.
type AnalysisStatus string

const (
	AnalysisStatusRunning   AnalysisStatus = "running"
	AnalysisStatusCompleted AnalysisStatus = "completed"
	AnalysisStatusFailed    AnalysisStatus = "failed"
	AnalysisStatusCanceled  AnalysisStatus = "canceled"
)

type AnalysisSummary struct {
	Analysis Analysis `json:"analysis" api:"required"`
	Target   Target   `json:"target" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Analysis    respjson.Field
		Target      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AnalysisSummary) RawJSON() string { return r.JSON.raw }
func (r *AnalysisSummary) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser settings that can be passed directly to `POST /browsers`.
type Browser struct {
	GPU      bool `json:"gpu" api:"required"`
	Headless bool `json:"headless" api:"required"`
	Stealth  bool `json:"stealth" api:"required"`
	// Initial browser window size in pixels with optional refresh rate. If omitted,
	// image defaults apply (1920x1080@25). For GPU images, the default is
	// 1920x1080@60. Arbitrary viewport dimensions and refresh rates are accepted.
	// Known-good presets include: 2560x1440@10, 1920x1080@25, 1920x1200@25,
	// 1440x900@25, 1280x800@60, 1024x768@60, 1200x800@60, 768x1024@60, 390x844@60. For
	// GPU images, recommended presets use one of these resolutions with refresh rates
	// 60, 30, 25, or 10: 800x600, 960x720, 1024x576, 1024x768, 1152x648, 1200x800,
	// 1280x720, 1368x768, 1440x900, 1600x900, 1920x1080, 1920x1200, 390x844, 360x250,
	// 768x1024, 800x1600. Viewports outside this list may exhibit unstable live view
	// or recording behavior. If refresh_rate is not provided, it will be automatically
	// determined based on the resolution (higher resolutions use lower refresh rates
	// to keep bandwidth reasonable).
	Viewport shared.BrowserViewport `json:"viewport" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		GPU         respjson.Field
		Headless    respjson.Field
		Stealth     respjson.Field
		Viewport    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Browser) RawJSON() string { return r.JSON.raw }
func (r *Browser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ConfigRegistryResponse struct {
	// Pollable analysis after workflow submission is acknowledged. Null when no
	// refresh was submitted.
	Analysis Analysis `json:"analysis" api:"required"`
	// A recommendation or a structured no-recommendation result.
	Recommendation RecommendationResultUnion `json:"recommendation" api:"required"`
	Target         Target                    `json:"target" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Analysis       respjson.Field
		Recommendation respjson.Field
		Target         respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConfigRegistryResponse) RawJSON() string { return r.JSON.raw }
func (r *ConfigRegistryResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Evidence struct {
	Accessed     int64 `json:"accessed" api:"required"`
	Blocked      int64 `json:"blocked" api:"required"`
	Inconclusive int64 `json:"inconclusive" api:"required"`
	// Most recent contributing observation. Recommendations remain eligible regardless
	// of age and can be returned while a new analysis refreshes them.
	LastObservedAt time.Time `json:"last_observed_at" api:"required" format:"date-time"`
	RunCount       int64     `json:"run_count" api:"required"`
	// Number of judged trials.
	SampleSize int64 `json:"sample_size" api:"required"`
	// Accessed trials divided by judged trials. Inconclusive trials are excluded.
	SuccessRate float64 `json:"success_rate" api:"required"`
	// Most recent contributing run where this config met the success threshold.
	// Omitted for knowledge assembled from runs that did not independently meet the
	// threshold.
	LastVerifiedAt time.Time `json:"last_verified_at" api:"nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Accessed       respjson.Field
		Blocked        respjson.Field
		Inconclusive   respjson.Field
		LastObservedAt respjson.Field
		RunCount       respjson.Field
		SampleSize     respjson.Field
		SuccessRate    respjson.Field
		LastVerifiedAt respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Evidence) RawJSON() string { return r.JSON.raw }
func (r *Evidence) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property URL is required.
type LookupRequestParam struct {
	// Public HTTP(S) URL to look up.
	URL string `json:"url" api:"required" format:"uri"`
	// ISO 3166 country codes Kernel may use when returning a proxy configuration. When
	// omitted, Kernel uses its default country selection.
	AllowedProxyCountries []string `json:"allowed_proxy_countries,omitzero"`
	paramObj
}

func (r LookupRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow LookupRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *LookupRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type LookupResponse struct {
	Recommendation Recommendation `json:"recommendation" api:"required"`
	Target         Target         `json:"target" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Recommendation respjson.Field
		Target         respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LookupResponse) RawJSON() string { return r.JSON.raw }
func (r *LookupResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type NoRecommendation struct {
	// Machine-readable reason Kernel cannot currently provide a config recommendation.
	//
	// Any of "proxy_restricted", "no_working_configuration", "inconclusive".
	Code NoRecommendationCode `json:"code" api:"required"`
	// Human-readable explanation suitable for display.
	Message string                    `json:"message" api:"required"`
	Type    constant.NoRecommendation `json:"type" default:"no_recommendation"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r NoRecommendation) RawJSON() string { return r.JSON.raw }
func (r *NoRecommendation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Machine-readable reason Kernel cannot currently provide a config recommendation.
type NoRecommendationCode string

const (
	NoRecommendationCodeProxyRestricted        NoRecommendationCode = "proxy_restricted"
	NoRecommendationCodeNoWorkingConfiguration NoRecommendationCode = "no_working_configuration"
	NoRecommendationCodeInconclusive           NoRecommendationCode = "inconclusive"
)

// ProxyUnion contains all possible properties and values from [ProxyDirect],
// [ProxyManaged].
//
// Use the [ProxyUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ProxyUnion struct {
	// Any of "direct", "managed".
	Mode string `json:"mode"`
	// This field is from variant [ProxyManaged].
	Create ProxyManagedCreate `json:"create"`
	JSON   struct {
		Mode   respjson.Field
		Create respjson.Field
		raw    string
	} `json:"-"`
}

// anyProxy is implemented by each variant of [ProxyUnion] to add type safety for
// the return type of [ProxyUnion.AsAny]
type anyProxy interface {
	implProxyUnion()
}

func (ProxyDirect) implProxyUnion()  {}
func (ProxyManaged) implProxyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ProxyUnion.AsAny().(type) {
//	case kernel.ProxyDirect:
//	case kernel.ProxyManaged:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ProxyUnion) AsAny() anyProxy {
	switch u.Mode {
	case "direct":
		return u.AsDirect()
	case "managed":
		return u.AsManaged()
	}
	return nil
}

func (u ProxyUnion) AsDirect() (v ProxyDirect) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyUnion) AsManaged() (v ProxyManaged) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProxyUnion) RawJSON() string { return u.JSON.raw }

func (r *ProxyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Direct egress recipe. Pass `{ "mode": "direct" }` as the browser's `proxy`.
type ProxyDirect struct {
	Mode constant.Direct `json:"mode" default:"direct"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Mode        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyDirect) RawJSON() string { return r.JSON.raw }
func (r *ProxyDirect) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Managed proxy recipe. `create` is a non-idempotent `POST /proxies` payload:
// create the resource once, retain its ID, and reuse that ID as the browser's
// `proxy.id`. Do not submit this recipe before every browser session.
type ProxyManaged struct {
	// Configuration for routing traffic through a proxy.
	Create ProxyManagedCreate `json:"create" api:"required"`
	Mode   constant.Managed   `json:"mode" default:"managed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Create      respjson.Field
		Mode        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyManaged) RawJSON() string { return r.JSON.raw }
func (r *ProxyManaged) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for routing traffic through a proxy.
type ProxyManagedCreate struct {
	// Proxy type to use. In terms of quality for avoiding bot-detection, from best to
	// worst: `mobile` > `residential` > `isp` > `datacenter`.
	//
	// Any of "datacenter", "isp", "residential", "mobile", "custom".
	Type string `json:"type" api:"required"`
	// Hostnames that should bypass the parent proxy and connect directly.
	BypassHosts []string `json:"bypass_hosts"`
	// Configuration specific to the selected proxy `type`.
	Config ProxyManagedCreateConfigUnion `json:"config"`
	// Readable name of the proxy.
	Name string `json:"name"`
	// Protocol to use for the proxy connection.
	//
	// Any of "http", "https".
	Protocol string `json:"protocol"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		BypassHosts respjson.Field
		Config      respjson.Field
		Name        respjson.Field
		Protocol    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyManagedCreate) RawJSON() string { return r.JSON.raw }
func (r *ProxyManagedCreate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ProxyManagedCreateConfigUnion contains all possible properties and values from
// [ProxyManagedCreateConfigDatacenter], [ProxyManagedCreateConfigIsp],
// [ProxyManagedCreateConfigResidential], [ProxyManagedCreateConfigMobile],
// [ProxyManagedCreateConfigCustom].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ProxyManagedCreateConfigUnion struct {
	Country string `json:"country"`
	// This field is from variant [ProxyManagedCreateConfigResidential].
	Asn  string `json:"asn"`
	City string `json:"city"`
	// This field is from variant [ProxyManagedCreateConfigResidential].
	Os    string `json:"os"`
	State string `json:"state"`
	// This field is from variant [ProxyManagedCreateConfigResidential].
	Zip string `json:"zip"`
	// This field is from variant [ProxyManagedCreateConfigCustom].
	Host string `json:"host"`
	// This field is from variant [ProxyManagedCreateConfigCustom].
	Port int64 `json:"port"`
	// This field is from variant [ProxyManagedCreateConfigCustom].
	CaBundle string `json:"ca_bundle"`
	// This field is from variant [ProxyManagedCreateConfigCustom].
	Password string `json:"password"`
	// This field is from variant [ProxyManagedCreateConfigCustom].
	Username string `json:"username"`
	JSON     struct {
		Country  respjson.Field
		Asn      respjson.Field
		City     respjson.Field
		Os       respjson.Field
		State    respjson.Field
		Zip      respjson.Field
		Host     respjson.Field
		Port     respjson.Field
		CaBundle respjson.Field
		Password respjson.Field
		Username respjson.Field
		raw      string
	} `json:"-"`
}

func (u ProxyManagedCreateConfigUnion) AsDatacenter() (v ProxyManagedCreateConfigDatacenter) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyManagedCreateConfigUnion) AsIsp() (v ProxyManagedCreateConfigIsp) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyManagedCreateConfigUnion) AsResidential() (v ProxyManagedCreateConfigResidential) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyManagedCreateConfigUnion) AsMobile() (v ProxyManagedCreateConfigMobile) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProxyManagedCreateConfigUnion) AsCustom() (v ProxyManagedCreateConfigCustom) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProxyManagedCreateConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ProxyManagedCreateConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a datacenter proxy.
type ProxyManagedCreateConfigDatacenter struct {
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
func (r ProxyManagedCreateConfigDatacenter) RawJSON() string { return r.JSON.raw }
func (r *ProxyManagedCreateConfigDatacenter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for an ISP proxy.
type ProxyManagedCreateConfigIsp struct {
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
func (r ProxyManagedCreateConfigIsp) RawJSON() string { return r.JSON.raw }
func (r *ProxyManagedCreateConfigIsp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for residential proxies.
type ProxyManagedCreateConfigResidential struct {
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
func (r ProxyManagedCreateConfigResidential) RawJSON() string { return r.JSON.raw }
func (r *ProxyManagedCreateConfigResidential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for mobile proxies.
type ProxyManagedCreateConfigMobile struct {
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
func (r ProxyManagedCreateConfigMobile) RawJSON() string { return r.JSON.raw }
func (r *ProxyManagedCreateConfigMobile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for a custom proxy (e.g., private proxy server).
type ProxyManagedCreateConfigCustom struct {
	// Proxy host address or IP.
	Host string `json:"host" api:"required"`
	// Proxy port.
	Port int64 `json:"port" api:"required"`
	// PEM-encoded CA certificate bundle the proxy re-signs upstream TLS with. Provide
	// when the proxy terminates TLS (MITM) so the browser trusts its certificates. May
	// contain multiple concatenated certificates.
	CaBundle string `json:"ca_bundle"`
	// Password for proxy authentication.
	Password string `json:"password"`
	// Username for proxy authentication.
	Username string `json:"username"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Host        respjson.Field
		Port        respjson.Field
		CaBundle    respjson.Field
		Password    respjson.Field
		Username    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProxyManagedCreateConfigCustom) RawJSON() string { return r.JSON.raw }
func (r *ProxyManagedCreateConfigCustom) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Recommendation struct {
	// Browser settings that can be passed directly to `POST /browsers`.
	Browser  Browser  `json:"browser" api:"required"`
	Evidence Evidence `json:"evidence" api:"required"`
	// Specificity of knowledge matched for this recommendation.
	//
	// Any of "exact", "host", "domain".
	MatchScope RecommendationMatchScope `json:"match_scope" api:"required"`
	// Target value that supplied the recommendation.
	MatchedTarget string `json:"matched_target" api:"required"`
	// Proxy recipe for the recommended browser.
	Proxy ProxyUnion              `json:"proxy" api:"required"`
	Type  constant.Recommendation `json:"type" default:"recommendation"`
	// Exact matches meet the evidence threshold; host and domain fallbacks are
	// inferred. Check evidence.last_verified_at for successful verification age and
	// last_observed_at for the latest evidence.
	//
	// Any of "verified", "inferred".
	Verification RecommendationVerification `json:"verification" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Browser       respjson.Field
		Evidence      respjson.Field
		MatchScope    respjson.Field
		MatchedTarget respjson.Field
		Proxy         respjson.Field
		Type          respjson.Field
		Verification  respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Recommendation) RawJSON() string { return r.JSON.raw }
func (r *Recommendation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Specificity of knowledge matched for this recommendation.
type RecommendationMatchScope string

const (
	RecommendationMatchScopeExact  RecommendationMatchScope = "exact"
	RecommendationMatchScopeHost   RecommendationMatchScope = "host"
	RecommendationMatchScopeDomain RecommendationMatchScope = "domain"
)

// Exact matches meet the evidence threshold; host and domain fallbacks are
// inferred. Check evidence.last_verified_at for successful verification age and
// last_observed_at for the latest evidence.
type RecommendationVerification string

const (
	RecommendationVerificationVerified RecommendationVerification = "verified"
	RecommendationVerificationInferred RecommendationVerification = "inferred"
)

// RecommendationResultUnion contains all possible properties and values from
// [Recommendation], [NoRecommendation].
//
// Use the [RecommendationResultUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type RecommendationResultUnion struct {
	// This field is from variant [Recommendation].
	Browser Browser `json:"browser"`
	// This field is from variant [Recommendation].
	Evidence Evidence `json:"evidence"`
	// This field is from variant [Recommendation].
	MatchScope RecommendationMatchScope `json:"match_scope"`
	// This field is from variant [Recommendation].
	MatchedTarget string `json:"matched_target"`
	// This field is from variant [Recommendation].
	Proxy ProxyUnion `json:"proxy"`
	// Any of "recommendation", "no_recommendation".
	Type string `json:"type"`
	// This field is from variant [Recommendation].
	Verification RecommendationVerification `json:"verification"`
	// This field is from variant [NoRecommendation].
	Code NoRecommendationCode `json:"code"`
	// This field is from variant [NoRecommendation].
	Message string `json:"message"`
	JSON    struct {
		Browser       respjson.Field
		Evidence      respjson.Field
		MatchScope    respjson.Field
		MatchedTarget respjson.Field
		Proxy         respjson.Field
		Type          respjson.Field
		Verification  respjson.Field
		Code          respjson.Field
		Message       respjson.Field
		raw           string
	} `json:"-"`
}

// anyRecommendationResult is implemented by each variant of
// [RecommendationResultUnion] to add type safety for the return type of
// [RecommendationResultUnion.AsAny]
type anyRecommendationResult interface {
	implRecommendationResultUnion()
}

func (Recommendation) implRecommendationResultUnion()   {}
func (NoRecommendation) implRecommendationResultUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := RecommendationResultUnion.AsAny().(type) {
//	case kernel.Recommendation:
//	case kernel.NoRecommendation:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u RecommendationResultUnion) AsAny() anyRecommendationResult {
	switch u.Type {
	case "recommendation":
		return u.AsRecommendation()
	case "no_recommendation":
		return u.AsNoRecommendation()
	}
	return nil
}

func (u RecommendationResultUnion) AsRecommendation() (v Recommendation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RecommendationResultUnion) AsNoRecommendation() (v NoRecommendation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RecommendationResultUnion) RawJSON() string { return u.JSON.raw }

func (r *RecommendationResultUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RecommendationSummary struct {
	// ID of the most recently requested analysis for this domain.
	AnalysisID string `json:"analysis_id" api:"required"`
	// Lifecycle status of the most recently requested analysis for this domain.
	//
	// Any of "running", "completed", "failed", "canceled".
	AnalysisStatus RecommendationSummaryAnalysisStatus `json:"analysis_status" api:"required"`
	// Most recent time the selected project requested an analysis for this domain.
	LastRequestedAt time.Time `json:"last_requested_at" api:"required" format:"date-time"`
	// Current domain-level recommendation. Null when no eligible knowledge exists.
	Recommendation Recommendation `json:"recommendation" api:"required"`
	// Display label for the recommended browser configuration.
	RecommendedConfigLabel string `json:"recommended_config_label" api:"required"`
	// Success rate for the recommended configuration. Null when no eligible knowledge
	// exists.
	SuccessRate float64 `json:"success_rate" api:"required"`
	// Registrable domain previously analyzed by the selected project.
	Target string `json:"target" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AnalysisID             respjson.Field
		AnalysisStatus         respjson.Field
		LastRequestedAt        respjson.Field
		Recommendation         respjson.Field
		RecommendedConfigLabel respjson.Field
		SuccessRate            respjson.Field
		Target                 respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RecommendationSummary) RawJSON() string { return r.JSON.raw }
func (r *RecommendationSummary) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Lifecycle status of the most recently requested analysis for this domain.
type RecommendationSummaryAnalysisStatus string

const (
	RecommendationSummaryAnalysisStatusRunning   RecommendationSummaryAnalysisStatus = "running"
	RecommendationSummaryAnalysisStatusCompleted RecommendationSummaryAnalysisStatus = "completed"
	RecommendationSummaryAnalysisStatusFailed    RecommendationSummaryAnalysisStatus = "failed"
	RecommendationSummaryAnalysisStatusCanceled  RecommendationSummaryAnalysisStatus = "canceled"
)

// The property URL is required.
type ResolveRequestParam struct {
	// Public HTTP(S) URL to refresh.
	URL string `json:"url" api:"required" format:"uri"`
	// ISO 3166 country codes Kernel may use when searching for or returning a proxy
	// configuration. Kernel may test a subset of allowed countries. When omitted,
	// Kernel uses its default country selection.
	AllowedProxyCountries []string `json:"allowed_proxy_countries,omitzero"`
	paramObj
}

func (r ResolveRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow ResolveRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResolveRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Target struct {
	// Registrable domain.
	Domain string `json:"domain" api:"required"`
	// Full hostname, including subdomain.
	Host string `json:"host" api:"required"`
	// Exact normalized scheme, host, port, and path used for lookup.
	Normalized string `json:"normalized" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Domain      respjson.Field
		Host        respjson.Field
		Normalized  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Target) RawJSON() string { return r.JSON.raw }
func (r *Target) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ConfigRegistryListParams struct {
	Limit  param.Opt[int64] `query:"limit,omitzero" json:"-"`
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Case-insensitive domain search. Full URLs are reduced to their registrable
	// domain.
	Search param.Opt[string] `query:"search,omitzero" json:"-"`
	// Any of "target", "analysis_status", "recommended_config", "last_requested_at",
	// "success_rate".
	SortBy ConfigRegistryListParamsSortBy `query:"sort_by,omitzero" json:"-"`
	// Any of "asc", "desc".
	SortOrder ConfigRegistryListParamsSortOrder `query:"sort_order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ConfigRegistryListParams]'s query parameters as
// `url.Values`.
func (r ConfigRegistryListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type ConfigRegistryListParamsSortBy string

const (
	ConfigRegistryListParamsSortByTarget            ConfigRegistryListParamsSortBy = "target"
	ConfigRegistryListParamsSortByAnalysisStatus    ConfigRegistryListParamsSortBy = "analysis_status"
	ConfigRegistryListParamsSortByRecommendedConfig ConfigRegistryListParamsSortBy = "recommended_config"
	ConfigRegistryListParamsSortByLastRequestedAt   ConfigRegistryListParamsSortBy = "last_requested_at"
	ConfigRegistryListParamsSortBySuccessRate       ConfigRegistryListParamsSortBy = "success_rate"
)

type ConfigRegistryListParamsSortOrder string

const (
	ConfigRegistryListParamsSortOrderAsc  ConfigRegistryListParamsSortOrder = "asc"
	ConfigRegistryListParamsSortOrderDesc ConfigRegistryListParamsSortOrder = "desc"
)

type ConfigRegistryLookupParams struct {
	LookupRequest LookupRequestParam
	paramObj
}

func (r ConfigRegistryLookupParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.LookupRequest)
}
func (r *ConfigRegistryLookupParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ConfigRegistryResolveParams struct {
	ResolveRequest ResolveRequestParam
	paramObj
}

func (r ConfigRegistryResolveParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.ResolveRequest)
}
func (r *ConfigRegistryResolveParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
