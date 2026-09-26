// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"net/http"
	"net/url"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Search the web and retrieve content for selected results.
//
// SearchProviderService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSearchProviderService] method instead.
type SearchProviderService struct {
	Options []option.RequestOption
}

// NewSearchProviderService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSearchProviderService(opts ...option.RequestOption) (r SearchProviderService) {
	r = SearchProviderService{}
	r.Options = opts
	return
}

// Lists providers, capabilities, and machine-readable native-option schemas. Auto
// and fallback are strategies, not provider entries. The list is not paginated and
// contains no latency benchmarks. X-Request-Id identifies the request.
func (r *SearchProviderService) List(ctx context.Context, query SearchProviderListParams, opts ...option.RequestOption) (res *[]Provider, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "search/providers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type Provider struct {
	Content         ProviderContent         `json:"content" api:"required"`
	MaxResultsCap   int64                   `json:"max_results_cap" api:"required"`
	Params          ProviderParams          `json:"params" api:"required"`
	ProviderOptions ProviderProviderOptions `json:"provider_options" api:"required"`
	Slug            string                  `json:"slug" api:"required"`
	// Provider-specific limitations, conditional filter support, and warnings about
	// search modes.
	Notes []string `json:"notes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content         respjson.Field
		MaxResultsCap   respjson.Field
		Params          respjson.Field
		ProviderOptions respjson.Field
		Slug            respjson.Field
		Notes           respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Provider) RawJSON() string { return r.JSON.raw }
func (r *Provider) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderContent struct {
	// Can enforce the requested maximum content age.
	FreshnessControl bool `json:"freshness_control" api:"required"`
	// Supports content retrieval with the search request.
	Inline bool `json:"inline" api:"required"`
	// Supports content retrieval after the search completes.
	PostHoc bool `json:"post_hoc" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FreshnessControl respjson.Field
		Inline           respjson.Field
		PostHoc          respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProviderContent) RawJSON() string { return r.JSON.raw }
func (r *ProviderContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderParams struct {
	Country        ProviderParamsCountry        `json:"country" api:"required"`
	EndDate        ProviderParamsEndDate        `json:"end_date" api:"required"`
	ExcludeDomains ProviderParamsExcludeDomains `json:"exclude_domains" api:"required"`
	IncludeDomains ProviderParamsIncludeDomains `json:"include_domains" api:"required"`
	Language       ProviderParamsLanguage       `json:"language" api:"required"`
	Recency        ProviderParamsRecency        `json:"recency" api:"required"`
	SafeSearch     ProviderParamsSafeSearch     `json:"safe_search" api:"required"`
	StartDate      ProviderParamsStartDate      `json:"start_date" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Country        respjson.Field
		EndDate        respjson.Field
		ExcludeDomains respjson.Field
		IncludeDomains respjson.Field
		Language       respjson.Field
		Recency        respjson.Field
		SafeSearch     respjson.Field
		StartDate      respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProviderParams) RawJSON() string { return r.JSON.raw }
func (r *ProviderParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderParamsCountry struct {
	// Any of "native", "emulated", "unsupported".
	Support string `json:"support" api:"required"`
	// Translation behavior, limitations, and precision.
	Notes string `json:"notes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Support     respjson.Field
		Notes       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProviderParamsCountry) RawJSON() string { return r.JSON.raw }
func (r *ProviderParamsCountry) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderParamsEndDate struct {
	// Any of "native", "emulated", "unsupported".
	Support string `json:"support" api:"required"`
	// Translation behavior, limitations, and precision.
	Notes string `json:"notes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Support     respjson.Field
		Notes       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProviderParamsEndDate) RawJSON() string { return r.JSON.raw }
func (r *ProviderParamsEndDate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderParamsExcludeDomains struct {
	// Any of "native", "emulated", "unsupported".
	Support string `json:"support" api:"required"`
	// Translation behavior, limitations, and precision.
	Notes string `json:"notes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Support     respjson.Field
		Notes       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProviderParamsExcludeDomains) RawJSON() string { return r.JSON.raw }
func (r *ProviderParamsExcludeDomains) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderParamsIncludeDomains struct {
	// Any of "native", "emulated", "unsupported".
	Support string `json:"support" api:"required"`
	// Translation behavior, limitations, and precision.
	Notes string `json:"notes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Support     respjson.Field
		Notes       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProviderParamsIncludeDomains) RawJSON() string { return r.JSON.raw }
func (r *ProviderParamsIncludeDomains) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderParamsLanguage struct {
	// Any of "native", "emulated", "unsupported".
	Support string `json:"support" api:"required"`
	// Translation behavior, limitations, and precision.
	Notes string `json:"notes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Support     respjson.Field
		Notes       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProviderParamsLanguage) RawJSON() string { return r.JSON.raw }
func (r *ProviderParamsLanguage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderParamsRecency struct {
	// Any of "native", "emulated", "unsupported".
	Support string `json:"support" api:"required"`
	// Translation behavior, limitations, and precision.
	Notes string `json:"notes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Support     respjson.Field
		Notes       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProviderParamsRecency) RawJSON() string { return r.JSON.raw }
func (r *ProviderParamsRecency) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderParamsSafeSearch struct {
	// Any of "native", "emulated", "unsupported".
	Support string `json:"support" api:"required"`
	// Translation behavior, limitations, and precision.
	Notes string `json:"notes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Support     respjson.Field
		Notes       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProviderParamsSafeSearch) RawJSON() string { return r.JSON.raw }
func (r *ProviderParamsSafeSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderParamsStartDate struct {
	// Any of "native", "emulated", "unsupported".
	Support string `json:"support" api:"required"`
	// Translation behavior, limitations, and precision.
	Notes string `json:"notes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Support     respjson.Field
		Notes       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProviderParamsStartDate) RawJSON() string { return r.JSON.raw }
func (r *ProviderParamsStartDate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderProviderOptions struct {
	// JSON Schema for the provider-native options accepted by POST /search.
	Schema map[string]any `json:"schema" api:"required"`
	// OpenAPI component name for the matching typed provider-options schema.
	SchemaRef string           `json:"schema_ref" api:"required"`
	Examples  []map[string]any `json:"examples"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Schema      respjson.Field
		SchemaRef   respjson.Field
		Examples    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProviderProviderOptions) RawJSON() string { return r.JSON.raw }
func (r *ProviderProviderOptions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SearchProviderListParams struct {
	// Optional concrete provider slug filter. Omit to list every provider. A slug that
	// is not listed returns an empty array.
	//
	// Any of "brave", "exa", "perplexity", "context", "parallel", "valyu", "octen",
	// "you", "tavily", "serpapi".
	Slug SearchProviderListParamsSlug `query:"slug,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [SearchProviderListParams]'s query parameters as
// `url.Values`.
func (r SearchProviderListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Optional concrete provider slug filter. Omit to list every provider. A slug that
// is not listed returns an empty array.
type SearchProviderListParamsSlug string

const (
	SearchProviderListParamsSlugBrave      SearchProviderListParamsSlug = "brave"
	SearchProviderListParamsSlugExa        SearchProviderListParamsSlug = "exa"
	SearchProviderListParamsSlugPerplexity SearchProviderListParamsSlug = "perplexity"
	SearchProviderListParamsSlugContext    SearchProviderListParamsSlug = "context"
	SearchProviderListParamsSlugParallel   SearchProviderListParamsSlug = "parallel"
	SearchProviderListParamsSlugValyu      SearchProviderListParamsSlug = "valyu"
	SearchProviderListParamsSlugOcten      SearchProviderListParamsSlug = "octen"
	SearchProviderListParamsSlugYou        SearchProviderListParamsSlug = "you"
	SearchProviderListParamsSlugTavily     SearchProviderListParamsSlug = "tavily"
	SearchProviderListParamsSlugSerpapi    SearchProviderListParamsSlug = "serpapi"
)
