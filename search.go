// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
	"github.com/kernel/kernel-go-sdk/internal/paramutil"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
	"github.com/kernel/kernel-go-sdk/shared/constant"
)

// Search the web and retrieve content for selected results.
//
// SearchService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSearchService] method instead.
type SearchService struct {
	Options []option.RequestOption
	// Search the web and retrieve content for selected results.
	Contents SearchContentService
	// Search the web and retrieve content for selected results.
	Providers SearchProviderService
}

// NewSearchService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewSearchService(opts ...option.RequestOption) (r SearchService) {
	r = SearchService{}
	r.Options = opts
	r.Contents = NewSearchContentService(opts...)
	r.Providers = NewSearchProviderService(opts...)
	return
}

// Returns ranked results from one serving provider. The default strategy selects a
// provider that supports the requested options. The fallback strategy tries
// providers in the supplied order. Results are not blended across providers.
// Portable filters may be approximated or omitted according to provider
// capabilities; warnings describe those outcomes unless strict_params is true.
// Native options apply only to their selected provider.
func (r *SearchService) New(ctx context.Context, body SearchNewParams, opts ...option.RequestOption) (res *Search, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "search"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Returns the retained search resource exactly as it was returned by POST /search:
// results, attempts, warnings, usage, and expires_at. No provider is called and
// nothing is billed. Use it to look up a search by ID for debugging, cost review,
// or to recover result IDs before calling the contents endpoint. Inline content
// fetched at search time is included; content fetched later through the contents
// endpoint is not merged in. Missing, expired, or inaccessible searches
// return 404.
func (r *SearchService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *Search, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("search/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type Attempt struct {
	DurationMs int64 `json:"duration_ms" api:"required"`
	// Any of "success", "empty", "error", "timeout".
	Outcome   AttemptOutcome `json:"outcome" api:"required"`
	Provider  string         `json:"provider" api:"required"`
	ErrorCode string         `json:"error_code"`
	Retryable bool           `json:"retryable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DurationMs  respjson.Field
		Outcome     respjson.Field
		Provider    respjson.Field
		ErrorCode   respjson.Field
		Retryable   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Attempt) RawJSON() string { return r.JSON.raw }
func (r *Attempt) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AttemptOutcome string

const (
	AttemptOutcomeSuccess AttemptOutcome = "success"
	AttemptOutcomeEmpty   AttemptOutcome = "empty"
	AttemptOutcomeError   AttemptOutcome = "error"
	AttemptOutcomeTimeout AttemptOutcome = "timeout"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ProviderTargetUnionParam struct {
	OfBrave      *ProviderTargetBraveParam      `json:",omitzero,inline"`
	OfExa        *ProviderTargetExaParam        `json:",omitzero,inline"`
	OfPerplexity *ProviderTargetPerplexityParam `json:",omitzero,inline"`
	OfContext    *ProviderTargetContextParam    `json:",omitzero,inline"`
	OfParallel   *ProviderTargetParallelParam   `json:",omitzero,inline"`
	OfValyu      *ProviderTargetValyuParam      `json:",omitzero,inline"`
	OfOcten      *ProviderTargetOctenParam      `json:",omitzero,inline"`
	OfYou        *ProviderTargetYouParam        `json:",omitzero,inline"`
	OfTavily     *ProviderTargetTavilyParam     `json:",omitzero,inline"`
	OfSerpapi    *ProviderTargetSerpapiParam    `json:",omitzero,inline"`
	paramUnion
}

func (u ProviderTargetUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBrave,
		u.OfExa,
		u.OfPerplexity,
		u.OfContext,
		u.OfParallel,
		u.OfValyu,
		u.OfOcten,
		u.OfYou,
		u.OfTavily,
		u.OfSerpapi)
}
func (u *ProviderTargetUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ProviderTargetUnionParam) asAny() any {
	if !param.IsOmitted(u.OfBrave) {
		return u.OfBrave
	} else if !param.IsOmitted(u.OfExa) {
		return u.OfExa
	} else if !param.IsOmitted(u.OfPerplexity) {
		return u.OfPerplexity
	} else if !param.IsOmitted(u.OfContext) {
		return u.OfContext
	} else if !param.IsOmitted(u.OfParallel) {
		return u.OfParallel
	} else if !param.IsOmitted(u.OfValyu) {
		return u.OfValyu
	} else if !param.IsOmitted(u.OfOcten) {
		return u.OfOcten
	} else if !param.IsOmitted(u.OfYou) {
		return u.OfYou
	} else if !param.IsOmitted(u.OfTavily) {
		return u.OfTavily
	} else if !param.IsOmitted(u.OfSerpapi) {
		return u.OfSerpapi
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ProviderTargetUnionParam) GetProvider() *string {
	if vt := u.OfBrave; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfExa; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfPerplexity; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfContext; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfParallel; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfValyu; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfOcten; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfYou; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfTavily; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfSerpapi; vt != nil {
		return (*string)(&vt.Provider)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ProviderTargetUnionParam) GetOptions() (res providerTargetUnionParamOptions) {
	if vt := u.OfBrave; vt != nil {
		res.any = &vt.Options
	} else if vt := u.OfExa; vt != nil {
		res.any = &vt.Options
	} else if vt := u.OfPerplexity; vt != nil {
		res.any = &vt.Options
	} else if vt := u.OfContext; vt != nil {
		res.any = &vt.Options
	} else if vt := u.OfParallel; vt != nil {
		res.any = &vt.Options
	} else if vt := u.OfValyu; vt != nil {
		res.any = &vt.Options
	} else if vt := u.OfOcten; vt != nil {
		res.any = &vt.Options
	} else if vt := u.OfYou; vt != nil {
		res.any = &vt.Options
	} else if vt := u.OfTavily; vt != nil {
		res.any = &vt.Options
	} else if vt := u.OfSerpapi; vt != nil {
		res.any = &vt.Options
	}
	return
}

// Can have the runtime types [*ProviderTargetBraveOptionsParam],
// [*ProviderTargetExaOptionsParam], [*ProviderTargetPerplexityOptionsParam],
// [*ProviderTargetContextOptionsParam], [*ProviderTargetParallelOptionsParam],
// [*ProviderTargetValyuOptionsParam], [*ProviderTargetOctenOptionsParam],
// [*ProviderTargetYouOptionsParam], [*ProviderTargetTavilyOptionsParam],
// [*ProviderTargetSerpapiOptionsParam]
type providerTargetUnionParamOptions struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *kernel.ProviderTargetBraveOptionsParam:
//	case *kernel.ProviderTargetExaOptionsParam:
//	case *kernel.ProviderTargetPerplexityOptionsParam:
//	case *kernel.ProviderTargetContextOptionsParam:
//	case *kernel.ProviderTargetParallelOptionsParam:
//	case *kernel.ProviderTargetValyuOptionsParam:
//	case *kernel.ProviderTargetOctenOptionsParam:
//	case *kernel.ProviderTargetYouOptionsParam:
//	case *kernel.ProviderTargetTavilyOptionsParam:
//	case *kernel.ProviderTargetSerpapiOptionsParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u providerTargetUnionParamOptions) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetExtraSnippets() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return paramutil.AddrIfPresent(vt.ExtraSnippets)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetGoggles() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return paramutil.AddrIfPresent(vt.Goggles)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetGogglesID() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return paramutil.AddrIfPresent(vt.GogglesID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetIncludeFetchMetadata() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return paramutil.AddrIfPresent(vt.IncludeFetchMetadata)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetOperators() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return paramutil.AddrIfPresent(vt.Operators)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetResultFilter() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return paramutil.AddrIfPresent(vt.ResultFilter)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetSearchLang() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return paramutil.AddrIfPresent(vt.SearchLang)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetSpellcheck() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return paramutil.AddrIfPresent(vt.Spellcheck)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetUiLang() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return paramutil.AddrIfPresent(vt.UiLang)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetUnits() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return &vt.Units
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetCategory() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetExaOptionsParam:
		return paramutil.AddrIfPresent(vt.Category)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetCompliance() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetExaOptionsParam:
		return paramutil.AddrIfPresent(vt.Compliance)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetContents() *ProviderTargetExaOptionsContentsParam {
	switch vt := u.any.(type) {
	case *ProviderTargetExaOptionsParam:
		return &vt.Contents
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetMaxAgeHours() *float64 {
	switch vt := u.any.(type) {
	case *ProviderTargetExaOptionsParam:
		return paramutil.AddrIfPresent(vt.MaxAgeHours)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetType() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetExaOptionsParam:
		return &vt.Type
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetLastUpdatedAfterFilter() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetPerplexityOptionsParam:
		return paramutil.AddrIfPresent(vt.LastUpdatedAfterFilter)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetLastUpdatedBeforeFilter() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetPerplexityOptionsParam:
		return paramutil.AddrIfPresent(vt.LastUpdatedBeforeFilter)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetMaxTokens() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetPerplexityOptionsParam:
		return paramutil.AddrIfPresent(vt.MaxTokens)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetMaxTokensPerPage() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetPerplexityOptionsParam:
		return paramutil.AddrIfPresent(vt.MaxTokensPerPage)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetQuery() []string {
	switch vt := u.any.(type) {
	case *ProviderTargetPerplexityOptionsParam:
		return vt.Query
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetSearchContextSize() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetPerplexityOptionsParam:
		return &vt.SearchContextSize
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetSearchLanguageFilter() []string {
	switch vt := u.any.(type) {
	case *ProviderTargetPerplexityOptionsParam:
		return vt.SearchLanguageFilter
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetCountry() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetContextOptionsParam:
		return paramutil.AddrIfPresent(vt.Country)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetFreshness() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetContextOptionsParam:
		return &vt.Freshness
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetMarkdownOptions() *ProviderTargetContextOptionsMarkdownOptionsParam {
	switch vt := u.any.(type) {
	case *ProviderTargetContextOptionsParam:
		return &vt.MarkdownOptions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetQueryFanout() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetContextOptionsParam:
		return paramutil.AddrIfPresent(vt.QueryFanout)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetTags() []string {
	switch vt := u.any.(type) {
	case *ProviderTargetContextOptionsParam:
		return vt.Tags
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetTimeoutMs() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetContextOptionsParam:
		return paramutil.AddrIfPresent(vt.TimeoutMs)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetAdvancedSettings() *ProviderTargetParallelOptionsAdvancedSettingsParam {
	switch vt := u.any.(type) {
	case *ProviderTargetParallelOptionsParam:
		return &vt.AdvancedSettings
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetClientModel() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetParallelOptionsParam:
		return paramutil.AddrIfPresent(vt.ClientModel)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetMaxCharsTotal() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetParallelOptionsParam:
		return paramutil.AddrIfPresent(vt.MaxCharsTotal)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetMode() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetParallelOptionsParam:
		return &vt.Mode
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetObjective() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetParallelOptionsParam:
		return paramutil.AddrIfPresent(vt.Objective)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetSearchQueries() []string {
	switch vt := u.any.(type) {
	case *ProviderTargetParallelOptionsParam:
		return vt.SearchQueries
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetSessionID() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetParallelOptionsParam:
		return paramutil.AddrIfPresent(vt.SessionID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetFastMode() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return paramutil.AddrIfPresent(vt.FastMode)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetHistoricalCache() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return paramutil.AddrIfPresent(vt.HistoricalCache)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetIncludeAbstracts() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return paramutil.AddrIfPresent(vt.IncludeAbstracts)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetInstructions() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return paramutil.AddrIfPresent(vt.Instructions)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetIsToolCall() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return paramutil.AddrIfPresent(vt.IsToolCall)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetMaxNumResults() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return paramutil.AddrIfPresent(vt.MaxNumResults)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetMaxPrice() *float64 {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return paramutil.AddrIfPresent(vt.MaxPrice)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetRelevanceThreshold() *float64 {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return paramutil.AddrIfPresent(vt.RelevanceThreshold)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetResponseLength() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return &vt.ResponseLength
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetSearchType() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return paramutil.AddrIfPresent(vt.SearchType)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetSourceBiases() []string {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return vt.SourceBiases
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetURLOnly() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetValyuOptionsParam:
		return paramutil.AddrIfPresent(vt.URLOnly)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetEndTime() *time.Time {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return paramutil.AddrIfPresent(vt.EndTime)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetExcludeText() []string {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return vt.ExcludeText
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetFormat() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return &vt.Format
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetFullContent() *ProviderTargetOctenOptionsFullContentParam {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return &vt.FullContent
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetHighlight() *ProviderTargetOctenOptionsHighlightParam {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return &vt.Highlight
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetIncludeText() []string {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return vt.IncludeText
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetSafesearch() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return &vt.Safesearch
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetStartTime() *time.Time {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return paramutil.AddrIfPresent(vt.StartTime)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetTimeBasis() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return &vt.TimeBasis
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetTimeRange() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return &vt.TimeRange
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetBoostDomains() []string {
	switch vt := u.any.(type) {
	case *ProviderTargetYouOptionsParam:
		return vt.BoostDomains
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetCrawlTimeout() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetYouOptionsParam:
		return paramutil.AddrIfPresent(vt.CrawlTimeout)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetExtraction() *ProviderTargetYouOptionsExtractionParam {
	switch vt := u.any.(type) {
	case *ProviderTargetYouOptionsParam:
		return &vt.Extraction
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetKnowledge() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetYouOptionsParam:
		return &vt.Knowledge
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetAutoParameters() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetTavilyOptionsParam:
		return paramutil.AddrIfPresent(vt.AutoParameters)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetChunksPerSource() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetTavilyOptionsParam:
		return paramutil.AddrIfPresent(vt.ChunksPerSource)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetExactMatch() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetTavilyOptionsParam:
		return paramutil.AddrIfPresent(vt.ExactMatch)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetFilterByLanguage() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetTavilyOptionsParam:
		return paramutil.AddrIfPresent(vt.FilterByLanguage)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetIncludeAnswer() *ProviderTargetTavilyOptionsIncludeAnswerUnionParam {
	switch vt := u.any.(type) {
	case *ProviderTargetTavilyOptionsParam:
		return &vt.IncludeAnswer
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetIncludeDomainsMode() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetTavilyOptionsParam:
		return &vt.IncludeDomainsMode
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetIncludeFavicon() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetTavilyOptionsParam:
		return paramutil.AddrIfPresent(vt.IncludeFavicon)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetIncludeImageDescriptions() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetTavilyOptionsParam:
		return paramutil.AddrIfPresent(vt.IncludeImageDescriptions)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetIncludeRawContent() *ProviderTargetTavilyOptionsIncludeRawContentUnionParam {
	switch vt := u.any.(type) {
	case *ProviderTargetTavilyOptionsParam:
		return &vt.IncludeRawContent
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetSearchDepth() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetTavilyOptionsParam:
		return &vt.SearchDepth
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetEngine() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return &vt.Engine
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetDevice() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return &vt.Device
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetFilter() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return &vt.Filter
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetGl() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return paramutil.AddrIfPresent(vt.Gl)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetGoogleDomain() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return paramutil.AddrIfPresent(vt.GoogleDomain)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetHl() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return paramutil.AddrIfPresent(vt.Hl)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetLocation() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return paramutil.AddrIfPresent(vt.Location)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetNfpr() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return &vt.Nfpr
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetNoCache() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return paramutil.AddrIfPresent(vt.NoCache)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetNum() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return paramutil.AddrIfPresent(vt.Num)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetSafe() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return &vt.Safe
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetStart() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return paramutil.AddrIfPresent(vt.Start)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetTbm() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return paramutil.AddrIfPresent(vt.Tbm)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetTbs() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetSerpapiOptionsParam:
		return paramutil.AddrIfPresent(vt.Tbs)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetCount() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return paramutil.AddrIfPresent(vt.Count)
	case *ProviderTargetOctenOptionsParam:
		return paramutil.AddrIfPresent(vt.Count)
	case *ProviderTargetYouOptionsParam:
		return paramutil.AddrIfPresent(vt.Count)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetOffset() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetBraveOptionsParam:
		return paramutil.AddrIfPresent(vt.Offset)
	case *ProviderTargetYouOptionsParam:
		return paramutil.AddrIfPresent(vt.Offset)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetNumResults() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetExaOptionsParam:
		return paramutil.AddrIfPresent(vt.NumResults)
	case *ProviderTargetContextOptionsParam:
		return paramutil.AddrIfPresent(vt.NumResults)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetMaxResults() *int64 {
	switch vt := u.any.(type) {
	case *ProviderTargetPerplexityOptionsParam:
		return paramutil.AddrIfPresent(vt.MaxResults)
	case *ProviderTargetTavilyOptionsParam:
		return paramutil.AddrIfPresent(vt.MaxResults)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetIncludeImages() *bool {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return paramutil.AddrIfPresent(vt.IncludeImages)
	case *ProviderTargetTavilyOptionsParam:
		return paramutil.AddrIfPresent(vt.IncludeImages)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u providerTargetUnionParamOptions) GetTopic() *string {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		return (*string)(&vt.Topic)
	case *ProviderTargetTavilyOptionsParam:
		return (*string)(&vt.Topic)
	}
	return nil
}

// Returns a pointer to the underlying variant's ExcludeDomains property, if
// present.
func (u providerTargetUnionParamOptions) GetExcludeDomains() []string {
	switch vt := u.any.(type) {
	case *ProviderTargetContextOptionsParam:
		return vt.ExcludeDomains
	case *ProviderTargetOctenOptionsParam:
		return vt.ExcludeDomains
	}
	return nil
}

// Returns a pointer to the underlying variant's IncludeDomains property, if
// present.
func (u providerTargetUnionParamOptions) GetIncludeDomains() []string {
	switch vt := u.any.(type) {
	case *ProviderTargetContextOptionsParam:
		return vt.IncludeDomains
	case *ProviderTargetOctenOptionsParam:
		return vt.IncludeDomains
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u providerTargetUnionParamOptions) GetLanguage() (res providerTargetUnionParamOptionsLanguage) {
	switch vt := u.any.(type) {
	case *ProviderTargetOctenOptionsParam:
		res.any = &vt.Language
	case *ProviderTargetYouOptionsParam:
		res.any = paramutil.AddrIfPresent(vt.Language)
	case *ProviderTargetTavilyOptionsParam:
		res.any = paramutil.AddrIfPresent(vt.Language)
	}
	return res
}

// Can have the runtime types [*[]string], [*string]
type providerTargetUnionParamOptionsLanguage struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]string:
//	case *string:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u providerTargetUnionParamOptionsLanguage) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[ProviderTargetUnionParam](
		"provider",
		apijson.Discriminator[ProviderTargetBraveParam]("brave"),
		apijson.Discriminator[ProviderTargetExaParam]("exa"),
		apijson.Discriminator[ProviderTargetPerplexityParam]("perplexity"),
		apijson.Discriminator[ProviderTargetContextParam]("context"),
		apijson.Discriminator[ProviderTargetParallelParam]("parallel"),
		apijson.Discriminator[ProviderTargetValyuParam]("valyu"),
		apijson.Discriminator[ProviderTargetOctenParam]("octen"),
		apijson.Discriminator[ProviderTargetYouParam]("you"),
		apijson.Discriminator[ProviderTargetTavilyParam]("tavily"),
		apijson.Discriminator[ProviderTargetSerpapiParam]("serpapi"),
	)
}

// The property Provider is required.
type ProviderTargetBraveParam struct {
	Options ProviderTargetBraveOptionsParam `json:"options,omitzero"`
	// This field can be elided, and will marshal its zero value as "brave".
	Provider constant.Brave `json:"provider" default:"brave"`
	paramObj
}

func (r ProviderTargetBraveParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetBraveParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetBraveParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetBraveOptionsParam struct {
	// Provider-native count. Lower-only alias for max_results; cannot raise the
	// effective result cap.
	Count param.Opt[int64] `json:"count,omitzero"`
	// Request additional snippets from Brave.
	ExtraSnippets param.Opt[bool] `json:"extra_snippets,omitzero"`
	// Goggles re-ranking definition URL.
	Goggles param.Opt[string] `json:"goggles,omitzero"`
	// Deprecated Brave Goggle identifier. Prefer goggles.
	//
	// Deprecated: deprecated
	GogglesID param.Opt[string] `json:"goggles_id,omitzero"`
	// Include Brave's fetch metadata.
	IncludeFetchMetadata param.Opt[bool] `json:"include_fetch_metadata,omitzero"`
	// Page offset supported by Brave.
	Offset param.Opt[int64] `json:"offset,omitzero"`
	// Brave search operators.
	Operators param.Opt[string] `json:"operators,omitzero"`
	// Comma-separated result types to include, e.g. "web,news".
	ResultFilter param.Opt[string] `json:"result_filter,omitzero"`
	// Language of the search, e.g. "en".
	SearchLang param.Opt[string] `json:"search_lang,omitzero"`
	// Apply Brave's query spellcheck.
	Spellcheck param.Opt[bool] `json:"spellcheck,omitzero"`
	// Language for UI strings in the response.
	UiLang param.Opt[string] `json:"ui_lang,omitzero"`
	// Measurement units.
	//
	// Any of "metric", "imperial".
	Units string `json:"units,omitzero"`
	paramObj
}

func (r ProviderTargetBraveOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetBraveOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetBraveOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProviderTargetBraveOptionsParam](
		"units", "metric", "imperial",
	)
}

// The property Provider is required.
type ProviderTargetExaParam struct {
	Options ProviderTargetExaOptionsParam `json:"options,omitzero"`
	// This field can be elided, and will marshal its zero value as "exa".
	Provider constant.Exa `json:"provider" default:"exa"`
	paramObj
}

func (r ProviderTargetExaParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetExaParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetExaParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetExaOptionsParam struct {
	// Provider data-category hint.
	Category param.Opt[string] `json:"category,omitzero"`
	// Provider-native compliance controls. Requires support and authorization on the
	// provider account.
	Compliance param.Opt[string] `json:"compliance,omitzero"`
	// Provider-native cache-age control. Unlike content.max_age_hours, this retains
	// Exa semantics, including any native sentinel values. It does not imply a
	// cross-provider freshness guarantee.
	MaxAgeHours param.Opt[float64] `json:"maxAgeHours,omitzero"`
	// Provider-native count. Lower-only alias for max_results.
	NumResults param.Opt[int64] `json:"numResults,omitzero"`
	// Provider-native content retrieval. Available without requesting Kernel browser
	// retrieval; may incur provider retrieval charges.
	Contents ProviderTargetExaOptionsContentsParam `json:"contents,omitzero"`
	// Search mode supported by Exa.
	//
	// Any of "auto", "fast", "instant".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r ProviderTargetExaOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetExaOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetExaOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProviderTargetExaOptionsParam](
		"type", "auto", "fast", "instant",
	)
}

// Provider-native content retrieval. Available without requesting Kernel browser
// retrieval; may incur provider retrieval charges.
type ProviderTargetExaOptionsContentsParam struct {
	// Return query-relevant provider excerpts.
	Highlights param.Opt[bool] `json:"highlights,omitzero"`
	// Return provider page text.
	Text param.Opt[bool] `json:"text,omitzero"`
	paramObj
}

func (r ProviderTargetExaOptionsContentsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetExaOptionsContentsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetExaOptionsContentsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Provider is required.
type ProviderTargetPerplexityParam struct {
	Options ProviderTargetPerplexityOptionsParam `json:"options,omitzero"`
	// This field can be elided, and will marshal its zero value as "perplexity".
	Provider constant.Perplexity `json:"provider" default:"perplexity"`
	paramObj
}

func (r ProviderTargetPerplexityParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetPerplexityParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetPerplexityParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetPerplexityOptionsParam struct {
	// MM/DD/YYYY. Filters by last-updated date, not published date.
	LastUpdatedAfterFilter param.Opt[string] `json:"last_updated_after_filter,omitzero"`
	// MM/DD/YYYY upper bound on last-updated date.
	LastUpdatedBeforeFilter param.Opt[string] `json:"last_updated_before_filter,omitzero"`
	// Provider-native count. Lower-only alias for max_results.
	MaxResults param.Opt[int64] `json:"max_results,omitzero"`
	// Values outside the documented range are rejected.
	MaxTokens param.Opt[int64] `json:"max_tokens,omitzero"`
	// Per-page token cap.
	MaxTokensPerPage param.Opt[int64] `json:"max_tokens_per_page,omitzero"`
	// Provider-native multi-query form, applied only to Perplexity. Other fallback
	// providers receive the top-level query. Each query may incur a separate provider
	// charge.
	Query []string `json:"query,omitzero"`
	// Provider context size supported by the selected model.
	//
	// Any of "low", "medium", "high".
	SearchContextSize string `json:"search_context_size,omitzero"`
	// ISO 639-1 language codes, max 20.
	SearchLanguageFilter []string `json:"search_language_filter,omitzero"`
	paramObj
}

func (r ProviderTargetPerplexityOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetPerplexityOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetPerplexityOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProviderTargetPerplexityOptionsParam](
		"search_context_size", "low", "medium", "high",
	)
}

// The property Provider is required.
type ProviderTargetContextParam struct {
	Options ProviderTargetContextOptionsParam `json:"options,omitzero"`
	// This field can be elided, and will marshal its zero value as "context".
	Provider constant.Context `json:"provider" default:"context"`
	paramObj
}

func (r ProviderTargetContextParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetContextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetContextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetContextOptionsParam struct {
	// ISO 3166-1 alpha-2 country code.
	Country param.Opt[string] `json:"country,omitzero"`
	// Number of results to request from Context.dev.
	NumResults param.Opt[int64] `json:"numResults,omitzero"`
	// Expand the query into multiple parallel variants.
	QueryFanout param.Opt[bool] `json:"queryFanout,omitzero"`
	// Context.dev request timeout in milliseconds.
	TimeoutMs param.Opt[int64] `json:"timeoutMS,omitzero"`
	// Blocklist of result domains.
	ExcludeDomains []string `json:"excludeDomains,omitzero"`
	// Restrict results to content published within this window.
	//
	// Any of "last_24_hours", "last_week", "last_month", "last_year".
	Freshness string `json:"freshness,omitzero"`
	// Allowlist of result domains.
	IncludeDomains  []string                                         `json:"includeDomains,omitzero"`
	MarkdownOptions ProviderTargetContextOptionsMarkdownOptionsParam `json:"markdownOptions,omitzero"`
	// Usage tracking tags.
	Tags []string `json:"tags,omitzero"`
	paramObj
}

func (r ProviderTargetContextOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetContextOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetContextOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProviderTargetContextOptionsParam](
		"freshness", "last_24_hours", "last_week", "last_month", "last_year",
	)
}

type ProviderTargetContextOptionsMarkdownOptionsParam struct {
	Enabled             param.Opt[bool]                                     `json:"enabled,omitzero"`
	IncludeFrames       param.Opt[bool]                                     `json:"includeFrames,omitzero"`
	IncludeImages       param.Opt[bool]                                     `json:"includeImages,omitzero"`
	IncludeLinks        param.Opt[bool]                                     `json:"includeLinks,omitzero"`
	MaxAgeMs            param.Opt[int64]                                    `json:"maxAgeMs,omitzero"`
	ShortenBase64Images param.Opt[bool]                                     `json:"shortenBase64Images,omitzero"`
	TimeoutMs           param.Opt[int64]                                    `json:"timeoutMS,omitzero"`
	UseMainContentOnly  param.Opt[bool]                                     `json:"useMainContentOnly,omitzero"`
	WaitForMs           param.Opt[int64]                                    `json:"waitForMs,omitzero"`
	Pdf                 ProviderTargetContextOptionsMarkdownOptionsPdfParam `json:"pdf,omitzero"`
	paramObj
}

func (r ProviderTargetContextOptionsMarkdownOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetContextOptionsMarkdownOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetContextOptionsMarkdownOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetContextOptionsMarkdownOptionsPdfParam struct {
	ShouldParse param.Opt[bool] `json:"shouldParse,omitzero"`
	paramObj
}

func (r ProviderTargetContextOptionsMarkdownOptionsPdfParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetContextOptionsMarkdownOptionsPdfParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetContextOptionsMarkdownOptionsPdfParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Provider is required.
type ProviderTargetParallelParam struct {
	Options ProviderTargetParallelOptionsParam `json:"options,omitzero"`
	// This field can be elided, and will marshal its zero value as "parallel".
	Provider constant.Parallel `json:"provider" default:"parallel"`
	paramObj
}

func (r ProviderTargetParallelParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetParallelParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetParallelParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetParallelOptionsParam struct {
	// Client model hint.
	ClientModel param.Opt[string] `json:"client_model,omitzero"`
	// Cap total characters returned.
	MaxCharsTotal param.Opt[int64] `json:"max_chars_total,omitzero"`
	// The goal behind the search, stated separately from the query.
	Objective param.Opt[string] `json:"objective,omitzero"`
	// Group related searches.
	SessionID param.Opt[string] `json:"session_id,omitzero"`
	// Explicit search settings. Unified search parameters are re-applied to
	// overlapping settings; native result counts are lower-only.
	AdvancedSettings ProviderTargetParallelOptionsAdvancedSettingsParam `json:"advanced_settings,omitzero"`
	// Search mode. Basic is used when omitted; each mode can have different latency
	// and charges.
	//
	// Any of "turbo", "fast", "basic", "advanced".
	Mode string `json:"mode,omitzero"`
	// Provider-native multi-query search. Defaults to [query] for this provider.
	SearchQueries []string `json:"search_queries,omitzero"`
	paramObj
}

func (r ProviderTargetParallelOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetParallelOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetParallelOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProviderTargetParallelOptionsParam](
		"mode", "turbo", "fast", "basic", "advanced",
	)
}

// Explicit search settings. Unified search parameters are re-applied to
// overlapping settings; native result counts are lower-only.
type ProviderTargetParallelOptionsAdvancedSettingsParam struct {
	// Native ISO 3166-1 alpha-2 location preference.
	Location param.Opt[string] `json:"location,omitzero"`
	// Native count; lower-only alias for unified max_results.
	MaxResults      param.Opt[int64]                                                  `json:"max_results,omitzero"`
	ExcerptSettings ProviderTargetParallelOptionsAdvancedSettingsExcerptSettingsParam `json:"excerpt_settings,omitzero"`
	FetchPolicy     ProviderTargetParallelOptionsAdvancedSettingsFetchPolicyParam     `json:"fetch_policy,omitzero"`
	SourcePolicy    ProviderTargetParallelOptionsAdvancedSettingsSourcePolicyParam    `json:"source_policy,omitzero"`
	paramObj
}

func (r ProviderTargetParallelOptionsAdvancedSettingsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetParallelOptionsAdvancedSettingsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetParallelOptionsAdvancedSettingsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetParallelOptionsAdvancedSettingsExcerptSettingsParam struct {
	MaxCharsPerResult param.Opt[int64] `json:"max_chars_per_result,omitzero"`
	paramObj
}

func (r ProviderTargetParallelOptionsAdvancedSettingsExcerptSettingsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetParallelOptionsAdvancedSettingsExcerptSettingsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetParallelOptionsAdvancedSettingsExcerptSettingsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetParallelOptionsAdvancedSettingsFetchPolicyParam struct {
	// Native live-fetch trigger; minimum 600 seconds. Not the unified hard-freshness
	// control.
	MaxAgeSeconds param.Opt[int64] `json:"max_age_seconds,omitzero"`
	// Native live-fetch timeout, bounded by the remaining overall deadline.
	TimeoutSeconds param.Opt[float64] `json:"timeout_seconds,omitzero"`
	// When false, the provider may return cached content after live fetching fails.
	DisableCacheFallback param.Opt[bool] `json:"disable_cache_fallback,omitzero"`
	paramObj
}

func (r ProviderTargetParallelOptionsAdvancedSettingsFetchPolicyParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetParallelOptionsAdvancedSettingsFetchPolicyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetParallelOptionsAdvancedSettingsFetchPolicyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetParallelOptionsAdvancedSettingsSourcePolicyParam struct {
	// Native publication-date lower bound.
	AfterDate param.Opt[time.Time] `json:"after_date,omitzero" format:"date"`
	// Native exclusions; the provider ignores these when native include_domains is
	// non-empty.
	ExcludeDomains []string `json:"exclude_domains,omitzero"`
	// Native domain/path restrictions. Explicit unified domain parameters take
	// precedence. Combined include/exclude native lists cannot exceed 200 entries.
	IncludeDomains []string `json:"include_domains,omitzero"`
	paramObj
}

func (r ProviderTargetParallelOptionsAdvancedSettingsSourcePolicyParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetParallelOptionsAdvancedSettingsSourcePolicyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetParallelOptionsAdvancedSettingsSourcePolicyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Provider is required.
type ProviderTargetValyuParam struct {
	Options ProviderTargetValyuOptionsParam `json:"options,omitzero"`
	// This field can be elided, and will marshal its zero value as "valyu".
	Provider constant.Valyu `json:"provider" default:"valyu"`
	paramObj
}

func (r ProviderTargetValyuParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetValyuParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetValyuParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetValyuOptionsParam struct {
	// Trade depth for latency.
	FastMode param.Opt[bool] `json:"fast_mode,omitzero"`
	// Allow historical cached results.
	HistoricalCache param.Opt[bool] `json:"historical_cache,omitzero"`
	// Include abstracts for academic sources.
	IncludeAbstracts param.Opt[bool] `json:"include_abstracts,omitzero"`
	// Natural-language retrieval guidance.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// Mark the search as an agent tool call.
	IsToolCall param.Opt[bool] `json:"is_tool_call,omitzero"`
	// Provider-native count. Lower-only alias for max_results.
	MaxNumResults param.Opt[int64] `json:"max_num_results,omitzero"`
	// Provider-native USD-per-thousand-results price ceiling. Forwarded to Valyu.
	MaxPrice param.Opt[float64] `json:"max_price,omitzero"`
	// Provider-native minimum relevance threshold. Not a normalized cross-provider
	// score.
	RelevanceThreshold param.Opt[float64] `json:"relevance_threshold,omitzero"`
	// Corpus selector, including all, web, proprietary, and news. Provider corpus
	// choice may change billing; Kernel does not force web-only searches.
	SearchType param.Opt[string] `json:"search_type,omitzero"`
	// Return URLs without content.
	URLOnly param.Opt[bool] `json:"url_only,omitzero"`
	// Provider result-content length preset.
	//
	// Any of "short", "medium", "large", "max".
	ResponseLength string `json:"response_length,omitzero"`
	// Bias retrieval toward these sources.
	SourceBiases []string `json:"source_biases,omitzero"`
	paramObj
}

func (r ProviderTargetValyuOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetValyuOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetValyuOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProviderTargetValyuOptionsParam](
		"response_length", "short", "medium", "large", "max",
	)
}

// The property Provider is required.
type ProviderTargetOctenParam struct {
	Options ProviderTargetOctenOptionsParam `json:"options,omitzero"`
	// This field can be elided, and will marshal its zero value as "octen".
	Provider constant.Octen `json:"provider" default:"octen"`
	paramObj
}

func (r ProviderTargetOctenParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetOctenParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetOctenParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetOctenOptionsParam struct {
	Count          param.Opt[int64]     `json:"count,omitzero"`
	EndTime        param.Opt[time.Time] `json:"end_time,omitzero" format:"date-time"`
	IncludeImages  param.Opt[bool]      `json:"include_images,omitzero"`
	StartTime      param.Opt[time.Time] `json:"start_time,omitzero" format:"date-time"`
	ExcludeDomains []string             `json:"exclude_domains,omitzero"`
	ExcludeText    []string             `json:"exclude_text,omitzero"`
	// Any of "markdown", "text".
	Format         string                                     `json:"format,omitzero"`
	FullContent    ProviderTargetOctenOptionsFullContentParam `json:"full_content,omitzero"`
	Highlight      ProviderTargetOctenOptionsHighlightParam   `json:"highlight,omitzero"`
	IncludeDomains []string                                   `json:"include_domains,omitzero"`
	IncludeText    []string                                   `json:"include_text,omitzero"`
	Language       []string                                   `json:"language,omitzero"`
	// Any of "off", "strict".
	Safesearch string `json:"safesearch,omitzero"`
	// Any of "auto", "published", "crawled".
	TimeBasis string `json:"time_basis,omitzero"`
	// Any of "day", "week", "month", "year", "d", "w", "m", "y".
	TimeRange string `json:"time_range,omitzero"`
	// Any of "general", "news".
	Topic string `json:"topic,omitzero"`
	paramObj
}

func (r ProviderTargetOctenOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetOctenOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetOctenOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProviderTargetOctenOptionsParam](
		"format", "markdown", "text",
	)
	apijson.RegisterFieldValidator[ProviderTargetOctenOptionsParam](
		"safesearch", "off", "strict",
	)
	apijson.RegisterFieldValidator[ProviderTargetOctenOptionsParam](
		"time_basis", "auto", "published", "crawled",
	)
	apijson.RegisterFieldValidator[ProviderTargetOctenOptionsParam](
		"time_range", "day", "week", "month", "year", "d", "w", "m", "y",
	)
	apijson.RegisterFieldValidator[ProviderTargetOctenOptionsParam](
		"topic", "general", "news",
	)
}

type ProviderTargetOctenOptionsFullContentParam struct {
	Enable    param.Opt[bool]  `json:"enable,omitzero"`
	MaxTokens param.Opt[int64] `json:"max_tokens,omitzero"`
	paramObj
}

func (r ProviderTargetOctenOptionsFullContentParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetOctenOptionsFullContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetOctenOptionsFullContentParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetOctenOptionsHighlightParam struct {
	Enable    param.Opt[bool]  `json:"enable,omitzero"`
	MaxTokens param.Opt[int64] `json:"max_tokens,omitzero"`
	paramObj
}

func (r ProviderTargetOctenOptionsHighlightParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetOctenOptionsHighlightParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetOctenOptionsHighlightParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Provider is required.
type ProviderTargetYouParam struct {
	Options ProviderTargetYouOptionsParam `json:"options,omitzero"`
	// This field can be elided, and will marshal its zero value as "you".
	Provider constant.You `json:"provider" default:"you"`
	paramObj
}

func (r ProviderTargetYouParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetYouParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetYouParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetYouOptionsParam struct {
	// Provider-native per-section count. Lower-only alias for max_results. Web and
	// news sections may produce more rows than Kernel returns.
	Count param.Opt[int64] `json:"count,omitzero"`
	// Native extraction timeout in seconds, bounded by the remaining Kernel request
	// deadline.
	CrawlTimeout param.Opt[int64] `json:"crawl_timeout,omitzero"`
	// BCP 47 result language from You.com's 51-value enum, e.g. "EN", "JA". Default
	// EN.
	Language param.Opt[string] `json:"language,omitzero"`
	// Page offset supported by You.com.
	Offset param.Opt[int64] `json:"offset,omitzero"`
	// Prefer these domains without excluding others. Cannot be combined with
	// include_domains if the provider does not accept the combination.
	BoostDomains []string `json:"boost_domains,omitzero"`
	// Provider-native page extraction. Both modes may incur per-row charges; full_page
	// may retrieve web and news rows.
	Extraction ProviderTargetYouOptionsExtractionParam `json:"extraction,omitzero"`
	// Request licensed-data output. URL-less knowledge entries are not converted into
	// web results; include_raw exposes the full provider response separately.
	//
	// Any of "core".
	Knowledge string `json:"knowledge,omitzero"`
	paramObj
}

func (r ProviderTargetYouOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetYouOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetYouOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProviderTargetYouOptionsParam](
		"knowledge", "core",
	)
}

// Provider-native page extraction. Both modes may incur per-row charges; full_page
// may retrieve web and news rows.
//
// The property ExtractionMode is required.
type ProviderTargetYouOptionsExtractionParam struct {
	// Any of "highlights", "full_page".
	ExtractionMode string                                          `json:"extraction_mode,omitzero" api:"required"`
	FullPage       ProviderTargetYouOptionsExtractionFullPageParam `json:"full_page,omitzero"`
	paramObj
}

func (r ProviderTargetYouOptionsExtractionParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetYouOptionsExtractionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetYouOptionsExtractionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProviderTargetYouOptionsExtractionParam](
		"extraction_mode", "highlights", "full_page",
	)
}

type ProviderTargetYouOptionsExtractionFullPageParam struct {
	// Any of "html", "markdown".
	ExtractionFormats []string `json:"extraction_formats,omitzero"`
	paramObj
}

func (r ProviderTargetYouOptionsExtractionFullPageParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetYouOptionsExtractionFullPageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetYouOptionsExtractionFullPageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Provider is required.
type ProviderTargetTavilyParam struct {
	Options ProviderTargetTavilyOptionsParam `json:"options,omitzero"`
	// This field can be elided, and will marshal its zero value as "tavily".
	Provider constant.Tavily `json:"provider" default:"tavily"`
	paramObj
}

func (r ProviderTargetTavilyParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetTavilyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetTavilyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProviderTargetTavilyOptionsParam struct {
	// Allow the provider to choose search parameters. May select a different billing
	// tier; explicit caller values retain provider-native precedence.
	AutoParameters param.Opt[bool] `json:"auto_parameters,omitzero"`
	// Provider excerpts per source, up to 500 characters each.
	ChunksPerSource param.Opt[int64] `json:"chunks_per_source,omitzero"`
	// Require the quoted phrases in the query verbatim, bypassing semantic matches.
	ExactMatch param.Opt[bool] `json:"exact_match,omitzero"`
	// Strictly filter non-matching languages. Requires `language`.
	FilterByLanguage param.Opt[bool] `json:"filter_by_language,omitzero"`
	// Favicon URL per result.
	IncludeFavicon param.Opt[bool] `json:"include_favicon,omitzero"`
	// Describe each image. Needs include_images.
	IncludeImageDescriptions param.Opt[bool] `json:"include_image_descriptions,omitzero"`
	// Query-related images plus per-result images.
	IncludeImages param.Opt[bool] `json:"include_images,omitzero"`
	// ISO 639-1 code or English language name. Ranking boost unless
	// filter_by_language.
	Language param.Opt[string] `json:"language,omitzero"`
	// Provider-native count. Lower-only alias for max_results.
	MaxResults param.Opt[int64] `json:"max_results,omitzero"`
	// Request the provider's generated answer. Returned as answer on the search
	// response, independently of include_raw.
	IncludeAnswer ProviderTargetTavilyOptionsIncludeAnswerUnionParam `json:"include_answer,omitzero"`
	// Native filter versus ranking boost semantics. Boost influences ranking rather
	// than restricting results to the listed domains. Requires include_domains.
	//
	// Any of "filter", "boost".
	IncludeDomainsMode string `json:"include_domains_mode,omitzero"`
	// Request native full-page content. Defaults to markdown when omitted for this
	// provider, False disables that native retrieval; it does not disable explicitly
	// requested Kernel browser retrieval.
	IncludeRawContent ProviderTargetTavilyOptionsIncludeRawContentUnionParam `json:"include_raw_content,omitzero"`
	// Provider relevance and latency tier. Some tiers cannot be combined with native
	// safe search; conflicts are described in warnings.
	//
	// Any of "advanced", "basic", "fast", "ultra-fast".
	SearchDepth string `json:"search_depth,omitzero"`
	// Provider corpus selector. Publication metadata depends on the selected topic.
	//
	// Any of "general", "news", "finance".
	Topic string `json:"topic,omitzero"`
	paramObj
}

func (r ProviderTargetTavilyOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetTavilyOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetTavilyOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProviderTargetTavilyOptionsParam](
		"include_domains_mode", "filter", "boost",
	)
	apijson.RegisterFieldValidator[ProviderTargetTavilyOptionsParam](
		"search_depth", "advanced", "basic", "fast", "ultra-fast",
	)
	apijson.RegisterFieldValidator[ProviderTargetTavilyOptionsParam](
		"topic", "general", "news", "finance",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ProviderTargetTavilyOptionsIncludeAnswerUnionParam struct {
	OfBool param.Opt[bool] `json:",omitzero,inline"`
	// Check if union is this variant with
	// !param.IsOmitted(union.OfProviderTargetTavilyOptionsIncludeAnswerString)
	OfProviderTargetTavilyOptionsIncludeAnswerString param.Opt[string] `json:",omitzero,inline"`
	paramUnion
}

func (u ProviderTargetTavilyOptionsIncludeAnswerUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBool, u.OfProviderTargetTavilyOptionsIncludeAnswerString)
}
func (u *ProviderTargetTavilyOptionsIncludeAnswerUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ProviderTargetTavilyOptionsIncludeAnswerUnionParam) asAny() any {
	if !param.IsOmitted(u.OfBool) {
		return &u.OfBool.Value
	} else if !param.IsOmitted(u.OfProviderTargetTavilyOptionsIncludeAnswerString) {
		return &u.OfProviderTargetTavilyOptionsIncludeAnswerString
	}
	return nil
}

type ProviderTargetTavilyOptionsIncludeAnswerString string

const (
	ProviderTargetTavilyOptionsIncludeAnswerStringBasic    ProviderTargetTavilyOptionsIncludeAnswerString = "basic"
	ProviderTargetTavilyOptionsIncludeAnswerStringAdvanced ProviderTargetTavilyOptionsIncludeAnswerString = "advanced"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ProviderTargetTavilyOptionsIncludeRawContentUnionParam struct {
	OfBool param.Opt[bool] `json:",omitzero,inline"`
	// Check if union is this variant with
	// !param.IsOmitted(union.OfProviderTargetTavilyOptionsIncludeRawContentString)
	OfProviderTargetTavilyOptionsIncludeRawContentString param.Opt[string] `json:",omitzero,inline"`
	paramUnion
}

func (u ProviderTargetTavilyOptionsIncludeRawContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBool, u.OfProviderTargetTavilyOptionsIncludeRawContentString)
}
func (u *ProviderTargetTavilyOptionsIncludeRawContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ProviderTargetTavilyOptionsIncludeRawContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfBool) {
		return &u.OfBool.Value
	} else if !param.IsOmitted(u.OfProviderTargetTavilyOptionsIncludeRawContentString) {
		return &u.OfProviderTargetTavilyOptionsIncludeRawContentString
	}
	return nil
}

type ProviderTargetTavilyOptionsIncludeRawContentString string

const (
	ProviderTargetTavilyOptionsIncludeRawContentStringMarkdown ProviderTargetTavilyOptionsIncludeRawContentString = "markdown"
	ProviderTargetTavilyOptionsIncludeRawContentStringText     ProviderTargetTavilyOptionsIncludeRawContentString = "text"
)

// The property Provider is required.
type ProviderTargetSerpapiParam struct {
	Options ProviderTargetSerpapiOptionsParam `json:"options,omitzero"`
	// This field can be elided, and will marshal its zero value as "serpapi".
	Provider constant.Serpapi `json:"provider" default:"serpapi"`
	paramObj
}

func (r ProviderTargetSerpapiParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetSerpapiParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetSerpapiParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Engine is required.
type ProviderTargetSerpapiOptionsParam struct {
	// SerpApi engine identifier. The Kernel integration currently supports google
	// only.
	Engine string `json:"engine" api:"required"`
	// Two-letter Google country code.
	Gl param.Opt[string] `json:"gl,omitzero"`
	// Google domain to search when using the google engine.
	GoogleDomain param.Opt[string] `json:"google_domain,omitzero"`
	// Interface language code.
	Hl param.Opt[string] `json:"hl,omitzero"`
	// Free-form geographic location used for localized results.
	Location param.Opt[string] `json:"location,omitzero"`
	// When true, bypass SerpApi cached results when supported.
	NoCache param.Opt[bool] `json:"no_cache,omitzero"`
	// Number of results requested from the search engine.
	Num param.Opt[int64] `json:"num,omitzero"`
	// Zero-based result offset for pagination.
	Start param.Opt[int64] `json:"start,omitzero"`
	// Google vertical search selector, such as images, video, news, or shopping.
	Tbm param.Opt[string] `json:"tbm,omitzero"`
	// Google time and search modifiers, including freshness filters.
	Tbs param.Opt[string] `json:"tbs,omitzero"`
	// Device profile used for the search.
	//
	// Any of "desktop", "mobile", "tablet".
	Device string `json:"device,omitzero"`
	// Google duplicate-content filter.
	//
	// Any of 0, 1.
	Filter int64 `json:"filter,omitzero"`
	// Google auto-correction filter.
	//
	// Any of 0, 1.
	Nfpr int64 `json:"nfpr,omitzero"`
	// Safe-search setting for engines that support it.
	//
	// Any of "active", "off".
	Safe string `json:"safe,omitzero"`
	paramObj
}

func (r ProviderTargetSerpapiOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ProviderTargetSerpapiOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProviderTargetSerpapiOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ProviderTargetSerpapiOptionsParam](
		"device", "desktop", "mobile", "tablet",
	)
	apijson.RegisterFieldValidator[ProviderTargetSerpapiOptionsParam](
		"filter", 0, 1,
	)
	apijson.RegisterFieldValidator[ProviderTargetSerpapiOptionsParam](
		"nfpr", 0, 1,
	)
	apijson.RegisterFieldValidator[ProviderTargetSerpapiOptionsParam](
		"safe", "active", "off",
	)
}

// The property Query is required.
type RequestParam struct {
	// Primary search query. A provider-native multi-query option applies only to that
	// provider; other providers in a fallback chain receive this query.
	Query string `json:"query" api:"required"`
	// ISO 3166-1 alpha-2 search locale preference.
	Country param.Opt[string] `json:"country,omitzero"`
	// Inclusive publication-date upper bound; must not precede start_date. If recency
	// is also supplied, recency takes precedence with a warning. Unsupported or
	// approximated filtering is reported, or rejected under strict_params.
	EndDate param.Opt[time.Time] `json:"end_date,omitzero" format:"date"`
	// Include untouched per-result payloads and the full serving-provider response in
	// raw fields. Off by default; native top-level outputs such as answer remain
	// available without it.
	IncludeRaw param.Opt[bool] `json:"include_raw,omitzero"`
	// BCP 47 search language preference.
	Language param.Opt[string] `json:"language,omitzero"`
	// Requested result count from 1 through 100. The effective count is clamped to the
	// serving provider's cap with a warning. Effective native counts are the lower of
	// this limit and supplied provider-native count aliases. Strict mode rejects
	// unsupported counts.
	MaxResults param.Opt[int64] `json:"max_results,omitzero"`
	// Inclusive publication-date lower bound. If recency is also supplied, recency
	// takes precedence with a warning. Provider date semantics, precision, and
	// unsupported filters are reported; unknown source dates are not fabricated or
	// universally post-filtered.
	StartDate param.Opt[time.Time] `json:"start_date,omitzero" format:"date"`
	// When false, unsupported portable parameters are omitted and approximations are
	// described in warnings. When true, every supplied portable parameter must be
	// honored exactly. Requests that cannot be served with those parameters are
	// rejected. This does not guarantee identical rankings or document timestamps
	// across indexes. Authentication and project isolation are always enforced.
	StrictParams param.Opt[bool] `json:"strict_params,omitzero"`
	// Overall deadline across search attempts and inline retrieval. No new attempt
	// starts after the deadline. Completed search results survive inline retrieval
	// timeouts.
	TimeoutMs param.Opt[int64] `json:"timeout_ms,omitzero"`
	// Optional portable content retrieval. Pass true for defaults or an options
	// object. Omission never starts Kernel browser work; provider-supplied content is
	// still returned when available, including when requested through native options.
	// Both inline and deferred retrieval use the same options schema.
	Content RequestContentUnionParam `json:"content,omitzero"`
	// Hostname exclusions, with the same best-effort/strict behavior as
	// include_domains. Provider-specific combinations that cannot be represented are
	// reported via warnings or rejected in strict mode.
	ExcludeDomains []string `json:"exclude_domains,omitzero"`
	// Hostname inclusion preference, matching a hostname and its subdomains. Empty
	// means unrestricted. Translated, emulated, or dropped with a warning according to
	// provider capability unless strict_params is true. Native boost modes remain
	// advisory and are identified in warnings.
	IncludeDomains []string `json:"include_domains,omitzero"`
	// Relative search window. Takes precedence over start_date/end_date with a warning
	// if both are set. Provider-native recency behavior is retained, including
	// documented hour-to-day widening. Unsupported filters are rejected only in strict
	// mode.
	//
	// Any of "hour", "day", "week", "month", "year".
	Recency RequestRecency `json:"recency,omitzero"`
	// Optional safety preference. Omit to use provider defaults. Unsupported values
	// are dropped with a warning unless strict_params is true. A search filter is not
	// an authorization boundary.
	//
	// Any of "off", "moderate", "strict".
	SafeSearch RequestSafeSearch `json:"safe_search,omitzero"`
	// Omitted strategy defaults to auto.
	Strategy StrategyUnionParam `json:"strategy,omitzero"`
	paramObj
}

func (r RequestParam) MarshalJSON() (data []byte, err error) {
	type shadow RequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RequestContentUnionParam struct {
	// Check if union is this variant with
	// !param.IsOmitted(union.OfRequestContentBoolean)
	OfRequestContentBoolean              param.Opt[bool]                          `json:",omitzero,inline"`
	OfRequestContentSearchContentOptions *RequestContentSearchContentOptionsParam `json:",omitzero,inline"`
	paramUnion
}

func (u RequestContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRequestContentBoolean, u.OfRequestContentSearchContentOptions)
}
func (u *RequestContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *RequestContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfRequestContentBoolean) {
		return &u.OfRequestContentBoolean
	} else if !param.IsOmitted(u.OfRequestContentSearchContentOptions) {
		return u.OfRequestContentSearchContentOptions
	}
	return nil
}

// Pass true to enable default portable content retrieval (auto source, markdown
// format, 10,000 char cap).
type RequestContentBoolean bool

const (
	RequestContentBooleanTrue RequestContentBoolean = true
)

type RequestContentSearchContentOptionsParam struct {
	// Maximum acceptable age of cached page content, measured from origin retrieval. 0
	// forces a live fetch. Governs the Kernel content cache, which is scoped to the
	// caller organization and project and separated by retrieval context; fetches
	// through a caller-supplied browser_id bypass that cache. Mapped to the provider
	// freshness control when source is provider and the provider supports one;
	// otherwise provider content age is reported as unknown via fetched_at.
	MaxAgeHours param.Opt[int64] `json:"max_age_hours,omitzero"`
	// Per-result Unicode character limit after extraction.
	MaxChars param.Opt[int64] `json:"max_chars,omitzero"`
	// Per-result deadline including capacity acquisition, retrieval, and extraction.
	// Also bounded by the overall request deadline.
	TimeoutMs param.Opt[int64] `json:"timeout_ms,omitzero"`
	// Invalid with source=provider. Supplying browser_id requires source=browser so
	// the chosen identity is not bypassed.
	Browser RequestContentSearchContentOptionsBrowserParam `json:"browser,omitzero"`
	// Any of "markdown", "text".
	Format string `json:"format,omitzero"`
	// provider uses the search provider's native content retrieval; browser fetches
	// each URL through a Kernel browser; auto prefers Kernel browser retrieval and
	// falls back to provider-native content when browser retrieval is unavailable or
	// unsuitable. Defaults to auto for both inline and deferred retrieval. Deferred
	// provider retrieval requires post_hoc capability; an explicit provider source
	// without it is a 400. Missing documents produce per-result unavailable outcomes,
	// not request failures.
	//
	// Any of "auto", "provider", "browser".
	Source string `json:"source,omitzero"`
	paramObj
}

func (r RequestContentSearchContentOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow RequestContentSearchContentOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RequestContentSearchContentOptionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RequestContentSearchContentOptionsParam](
		"format", "markdown", "text",
	)
	apijson.RegisterFieldValidator[RequestContentSearchContentOptionsParam](
		"source", "auto", "provider", "browser",
	)
}

// Invalid with source=provider. Supplying browser_id requires source=browser so
// the chosen identity is not bypassed.
type RequestContentSearchContentOptionsBrowserParam struct {
	// Existing browser session ID authorized for the caller and selected project.
	// Reuses its cookies, proxy, and browser configuration. Kernel does not delete a
	// caller-supplied browser. Render mode uses a temporary tab; website activity may
	// still change shared cookies and storage. When omitted, Kernel obtains isolated
	// browser capacity in the caller's account and releases it after retrieval. That
	// capacity is not retained for later interaction. Existing browser quotas apply.
	BrowserID param.Opt[string] `json:"browser_id,omitzero"`
	// Curl uses the browser HTTP stack without navigation or JavaScript execution.
	// Render navigates a temporary page and extracts from its DOM. The selected mode
	// is used for the retrieval.
	//
	// Any of "curl", "render".
	Mode string `json:"mode,omitzero"`
	paramObj
}

func (r RequestContentSearchContentOptionsBrowserParam) MarshalJSON() (data []byte, err error) {
	type shadow RequestContentSearchContentOptionsBrowserParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RequestContentSearchContentOptionsBrowserParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RequestContentSearchContentOptionsBrowserParam](
		"mode", "curl", "render",
	)
}

// Relative search window. Takes precedence over start_date/end_date with a warning
// if both are set. Provider-native recency behavior is retained, including
// documented hour-to-day widening. Unsupported filters are rejected only in strict
// mode.
type RequestRecency string

const (
	RequestRecencyHour  RequestRecency = "hour"
	RequestRecencyDay   RequestRecency = "day"
	RequestRecencyWeek  RequestRecency = "week"
	RequestRecencyMonth RequestRecency = "month"
	RequestRecencyYear  RequestRecency = "year"
)

// Optional safety preference. Omit to use provider defaults. Unsupported values
// are dropped with a warning unless strict_params is true. A search filter is not
// an authorization boundary.
type RequestSafeSearch string

const (
	RequestSafeSearchOff      RequestSafeSearch = "off"
	RequestSafeSearchModerate RequestSafeSearch = "moderate"
	RequestSafeSearchStrict   RequestSafeSearch = "strict"
)

type Result struct {
	// Kernel-generated identifier for this result. Stable only within the retained
	// search; not standardized across providers. Provider-native IDs, when available,
	// remain provider-specific raw fields.
	ID string `json:"id" api:"required"`
	// One-based position in the returned ranking.
	Rank int64 `json:"rank" api:"required"`
	// Provider-returned URL, not assumed canonical.
	URL                string   `json:"url" api:"required" format:"uri"`
	AdditionalSnippets []string `json:"additional_snippets"`
	// Portable retrieval outcome, or native content supplied by the search provider.
	// Identity fields remain on the enclosing result. Native excerpts are labeled
	// excerpt rather than full_page. Omission never triggers browser retrieval.
	Content ResultContent `json:"content"`
	// Provider-supplied date or timestamp, preserving available precision. No
	// publication date is fabricated. Retains the published field name.
	PublishedDate string `json:"published_date" api:"nullable"`
	// Original provider result, included only with include_raw=true. Provider
	// relevance scores are not normalized. Top-level provider data is available in
	// Search.raw.
	Raw     any    `json:"raw"`
	Snippet string `json:"snippet" api:"nullable"`
	// Provider source name or result URL hostname, when available.
	Source string `json:"source" api:"nullable"`
	Title  string `json:"title" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                 respjson.Field
		Rank               respjson.Field
		URL                respjson.Field
		AdditionalSnippets respjson.Field
		Content            respjson.Field
		PublishedDate      respjson.Field
		Raw                respjson.Field
		Snippet            respjson.Field
		Source             respjson.Field
		Title              respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Result) RawJSON() string { return r.JSON.raw }
func (r *Result) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Portable retrieval outcome, or native content supplied by the search provider.
// Identity fields remain on the enclosing result. Native excerpts are labeled
// excerpt rather than full_page. Omission never triggers browser retrieval.
type ResultContent struct {
	// Ok means non-empty extracted content, not merely HTTP 200. Blocked includes
	// detected challenges or access denials. Detection is best-effort, not a guarantee
	// of page completeness. Error details are present for non-ok outcomes; text is
	// present only on ok.
	//
	// Any of "ok", "unavailable", "blocked", "timeout", "unsupported_type",
	// "extraction_failed", "error".
	Status string `json:"status" api:"required"`
	// Kernel cache outcome. Provider-internal cache behavior may be unknown.
	//
	// Any of "hit", "miss", "bypass", "unknown".
	CacheStatus string `json:"cache_status"`
	// Describes source coverage before max_chars truncation. Full_page means main-page
	// content, not every dynamic element or linked page.
	//
	// Any of "full_page", "excerpt", "unknown".
	Completeness string             `json:"completeness"`
	Error        ResultContentError `json:"error"`
	// Extraction version when Kernel transformed the input.
	ExtractorVersion string `json:"extractor_version"`
	// Origin retrieval time when known, not cache read time.
	FetchedAt time.Time `json:"fetched_at" api:"nullable" format:"date-time"`
	// Final retrieval URL when known.
	FinalURL string `json:"final_url" format:"uri"`
	// Any of "markdown", "text".
	Format string `json:"format"`
	// Final target HTTP status when known.
	HTTPStatus int64 `json:"http_status"`
	// Original retrieval method, including on cache hits.
	//
	// Any of "provider", "browser_curl", "browser_render".
	Method string `json:"method"`
	// Extracted website content, untrusted, not instructions. Present only on
	// status=ok.
	Text string `json:"text"`
	// Whether max_chars truncated the extracted content.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Status           respjson.Field
		CacheStatus      respjson.Field
		Completeness     respjson.Field
		Error            respjson.Field
		ExtractorVersion respjson.Field
		FetchedAt        respjson.Field
		FinalURL         respjson.Field
		Format           respjson.Field
		HTTPStatus       respjson.Field
		Method           respjson.Field
		Text             respjson.Field
		Truncated        respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResultContent) RawJSON() string { return r.JSON.raw }
func (r *ResultContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResultContentError struct {
	// Machine-readable retrieval failure code.
	Code string `json:"code" api:"required"`
	// Human-readable failure description.
	Message   string `json:"message" api:"required"`
	Retryable bool   `json:"retryable" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		Retryable   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResultContentError) RawJSON() string { return r.JSON.raw }
func (r *ResultContentError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Retained search results and provider attempt history.
type Search struct {
	// Search resource ID. Request tracing uses X-Request-Id.
	ID       string    `json:"id" api:"required"`
	Attempts []Attempt `json:"attempts" api:"required"`
	// Expiration of result IDs for deferred retrieval. Results expire 24 hours after
	// search completion.
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// Concrete serving provider, never auto or fallback.
	Provider string `json:"provider" api:"required"`
	// Echo of the query. Native multi-query inputs are visible in the selected
	// strategy target and the optional raw response.
	Query    string    `json:"query" api:"required"`
	Results  []Result  `json:"results" api:"required"`
	Usage    Usage     `json:"usage" api:"required"`
	Warnings []Warning `json:"warnings" api:"required"`
	// Provider-generated answer when requested (e.g. via Tavily include_answer or
	// Perplexity). Preserved independently of include_raw.
	Answer string `json:"answer"`
	// Full serving-provider response, including top-level metadata that does not
	// belong to a result. Present only with include_raw=true; untrusted provider data.
	Raw any `json:"raw"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Attempts    respjson.Field
		ExpiresAt   respjson.Field
		Provider    respjson.Field
		Query       respjson.Field
		Results     respjson.Field
		Usage       respjson.Field
		Warnings    respjson.Field
		Answer      respjson.Field
		Raw         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Search) RawJSON() string { return r.JSON.raw }
func (r *Search) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func StrategyParamOfPinned[
	T ProviderTargetBraveParam | ProviderTargetExaParam | ProviderTargetPerplexityParam | ProviderTargetContextParam | ProviderTargetParallelParam | ProviderTargetValyuParam | ProviderTargetOctenParam | ProviderTargetYouParam | ProviderTargetTavilyParam | ProviderTargetSerpapiParam,
](provider T) StrategyUnionParam {
	var pinned StrategyPinnedParam
	switch v := any(provider).(type) {
	case ProviderTargetBraveParam:
		pinned.Provider.OfBrave = &v
	case ProviderTargetExaParam:
		pinned.Provider.OfExa = &v
	case ProviderTargetPerplexityParam:
		pinned.Provider.OfPerplexity = &v
	case ProviderTargetContextParam:
		pinned.Provider.OfContext = &v
	case ProviderTargetParallelParam:
		pinned.Provider.OfParallel = &v
	case ProviderTargetValyuParam:
		pinned.Provider.OfValyu = &v
	case ProviderTargetOctenParam:
		pinned.Provider.OfOcten = &v
	case ProviderTargetYouParam:
		pinned.Provider.OfYou = &v
	case ProviderTargetTavilyParam:
		pinned.Provider.OfTavily = &v
	case ProviderTargetSerpapiParam:
		pinned.Provider.OfSerpapi = &v
	}
	return StrategyUnionParam{OfPinned: &pinned}
}

func StrategyParamOfFallback(providers []ProviderTargetUnionParam) StrategyUnionParam {
	var fallback StrategyFallbackParam
	fallback.Providers = providers
	return StrategyUnionParam{OfFallback: &fallback}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type StrategyUnionParam struct {
	OfAuto     *StrategyAutoParam     `json:",omitzero,inline"`
	OfPinned   *StrategyPinnedParam   `json:",omitzero,inline"`
	OfFallback *StrategyFallbackParam `json:",omitzero,inline"`
	paramUnion
}

func (u StrategyUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAuto, u.OfPinned, u.OfFallback)
}
func (u *StrategyUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *StrategyUnionParam) asAny() any {
	if !param.IsOmitted(u.OfAuto) {
		return u.OfAuto
	} else if !param.IsOmitted(u.OfPinned) {
		return u.OfPinned
	} else if !param.IsOmitted(u.OfFallback) {
		return u.OfFallback
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u StrategyUnionParam) GetProviderOptions() []ProviderTargetUnionParam {
	if vt := u.OfAuto; vt != nil {
		return vt.ProviderOptions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u StrategyUnionParam) GetProvider() *ProviderTargetUnionParam {
	if vt := u.OfPinned; vt != nil {
		return &vt.Provider
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u StrategyUnionParam) GetProviders() []ProviderTargetUnionParam {
	if vt := u.OfFallback; vt != nil {
		return vt.Providers
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u StrategyUnionParam) GetType() *string {
	if vt := u.OfAuto; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfPinned; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFallback; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's FallbackOn property, if present.
func (u StrategyUnionParam) GetFallbackOn() []string {
	if vt := u.OfAuto; vt != nil {
		return vt.FallbackOn
	} else if vt := u.OfFallback; vt != nil {
		return vt.FallbackOn
	}
	return nil
}

func init() {
	apijson.RegisterUnion[StrategyUnionParam](
		"type",
		apijson.Discriminator[StrategyAutoParam]("auto"),
		apijson.Discriminator[StrategyPinnedParam]("pinned"),
		apijson.Discriminator[StrategyFallbackParam]("fallback"),
	)
}

// The property Type is required.
type StrategyAutoParam struct {
	// Conditions that advance to the next provider under auto routing or an explicit
	// providers chain. Ignored when provider pins a single provider. error means a
	// retryable provider failure, including rate limiting, not invalid caller input or
	// caller quotas. empty means zero results after required filtering. An empty list
	// disables fallback. If every attempt is empty or fails, the response is the first
	// valid empty response with the full attempt trail, or a 502 if none succeeded.
	//
	// Any of "error", "timeout", "empty".
	FallbackOn []string `json:"fallback_on,omitzero"`
	// Provider targets available to auto routing, each paired with typed native
	// options. Provider names must be unique.
	ProviderOptions []ProviderTargetUnionParam `json:"provider_options,omitzero"`
	// Let Kernel choose an eligible provider by capability fit.
	//
	// This field can be elided, and will marshal its zero value as "auto".
	Type constant.Auto `json:"type" default:"auto"`
	paramObj
}

func (r StrategyAutoParam) MarshalJSON() (data []byte, err error) {
	type shadow StrategyAutoParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *StrategyAutoParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Provider, Type are required.
type StrategyPinnedParam struct {
	// Provider name paired with its typed native options.
	Provider ProviderTargetUnionParam `json:"provider,omitzero" api:"required"`
	// Use exactly the selected provider with no cross-provider fallback.
	//
	// This field can be elided, and will marshal its zero value as "pinned".
	Type constant.Pinned `json:"type" default:"pinned"`
	paramObj
}

func (r StrategyPinnedParam) MarshalJSON() (data []byte, err error) {
	type shadow StrategyPinnedParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *StrategyPinnedParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Providers, Type are required.
type StrategyFallbackParam struct {
	// Ordered provider targets. Provider names must be unique.
	Providers []ProviderTargetUnionParam `json:"providers,omitzero" api:"required"`
	// Conditions that advance to the next provider under auto routing or an explicit
	// providers chain. Ignored when provider pins a single provider. error means a
	// retryable provider failure, including rate limiting, not invalid caller input or
	// caller quotas. empty means zero results after required filtering. An empty list
	// disables fallback. If every attempt is empty or fails, the response is the first
	// valid empty response with the full attempt trail, or a 502 if none succeeded.
	//
	// Any of "error", "timeout", "empty".
	FallbackOn []string `json:"fallback_on,omitzero"`
	// Try providers in order and advance when fallback_on matches the outcome.
	//
	// This field can be elided, and will marshal its zero value as "fallback".
	Type constant.Fallback `json:"type" default:"fallback"`
	paramObj
}

func (r StrategyFallbackParam) MarshalJSON() (data []byte, err error) {
	type shadow StrategyFallbackParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *StrategyFallbackParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Usage struct {
	// Number of result URLs for which a Kernel browser retrieval was attempted,
	// excluding cache-only hits.
	ContentFetches int64 `json:"content_fetches" api:"required"`
	// Number of result entries returned, including failed entries on the contents
	// endpoint.
	ResultsCount int64 `json:"results_count" api:"required"`
	// Total customer charge in USD when billing data is available.
	Cost float64 `json:"cost"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentFetches respjson.Field
		ResultsCount   respjson.Field
		Cost           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Usage) RawJSON() string { return r.JSON.raw }
func (r *Usage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Warning struct {
	// Examples: param_unsupported, preference_unsupported, max_results_clamped,
	// domains_truncated, recency_emulated, filter_emulated, date_filter_overridden,
	// provider_ineligible, fallback_failed, content_partial.
	Code     string `json:"code" api:"required"`
	Message  string `json:"message" api:"required"`
	Param    string `json:"param"`
	Provider string `json:"provider"`
	ResultID string `json:"result_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		Param       respjson.Field
		Provider    respjson.Field
		ResultID    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Warning) RawJSON() string { return r.JSON.raw }
func (r *Warning) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SearchNewParams struct {
	Request RequestParam
	paramObj
}

func (r SearchNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.Request)
}
func (r *SearchNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
