// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
)

// Search the web and retrieve content for selected results.
//
// SearchContentService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSearchContentService] method instead.
type SearchContentService struct {
	Options []option.RequestOption
}

// NewSearchContentService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSearchContentService(opts ...option.RequestOption) (r SearchContentService) {
	r = SearchContentService{}
	r.Options = opts
	return
}

// Deferred result-content retrieval is reserved but not available in this release.
// Requests return 404 until the retrieval implementation is shipped. X-Request-Id
// identifies this request separately from the search resource.
func (r *SearchContentService) Fetch(ctx context.Context, id string, body SearchContentFetchParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("search/%s/contents", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

type FetchRequestParam struct {
	// Maximum number of search results to fetch when result_ids is omitted, starting
	// from rank 1. Mutually exclusive with result_ids.
	Limit param.Opt[int64] `json:"limit,omitzero"`
	// Overall deadline across all selected results.
	TimeoutMs param.Opt[int64] `json:"timeout_ms,omitzero"`
	// Defaults to source:auto when omitted.
	Content FetchRequestContentParam `json:"content,omitzero"`
	// Kernel-generated IDs from the referenced retained search, in desired response
	// order. They are not provider-standard IDs. Mutually exclusive with limit.
	ResultIDs []string `json:"result_ids,omitzero"`
	paramObj
}

func (r FetchRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow FetchRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *FetchRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Defaults to source:auto when omitted.
type FetchRequestContentParam struct {
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
	Browser FetchRequestContentBrowserParam `json:"browser,omitzero"`
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

func (r FetchRequestContentParam) MarshalJSON() (data []byte, err error) {
	type shadow FetchRequestContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *FetchRequestContentParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[FetchRequestContentParam](
		"format", "markdown", "text",
	)
	apijson.RegisterFieldValidator[FetchRequestContentParam](
		"source", "auto", "provider", "browser",
	)
}

// Invalid with source=provider. Supplying browser_id requires source=browser so
// the chosen identity is not bypassed.
type FetchRequestContentBrowserParam struct {
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

func (r FetchRequestContentBrowserParam) MarshalJSON() (data []byte, err error) {
	type shadow FetchRequestContentBrowserParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *FetchRequestContentBrowserParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[FetchRequestContentBrowserParam](
		"mode", "curl", "render",
	)
}

type SearchContentFetchParams struct {
	FetchRequest FetchRequestParam
	paramObj
}

func (r SearchContentFetchParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.FetchRequest)
}
func (r *SearchContentFetchParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
