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
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
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

// Retrieves selected results from a retained search. Provide exactly one of
// result_ids or limit; the latter fetches the top results. Content defaults to
// source:auto. Responses preserve result_ids order. Unknown result IDs are
// rejected before retrieval starts. Missing, expired, or inaccessible searches
// return 404. Once retrieval begins, return one outcome per selected result,
// including timeout entries for work unfinished at the overall deadline. Browser
// retrievals run sequentially in result order, so later results may time out when
// earlier pages are slow. X-Request-Id identifies this request separately from the
// search resource.
func (r *SearchContentService) Fetch(ctx context.Context, id string, body SearchContentFetchParams, opts ...option.RequestOption) (res *Response, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("search/%s/contents", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
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
	// For source=auto, maximum acceptable age of retained provider content, measured
	// from when the search received it from the provider. A value of 0 disables reuse
	// of retained content, so every result is fetched through a browser.
	// source=provider reuses retained provider content without freshness validation.
	// source=browser always fetches through a browser and does not use this age limit.
	MaxAgeHours param.Opt[int64] `json:"max_age_hours,omitzero"`
	// Per-result Unicode character limit after extraction. Retained provider content
	// cannot exceed what was stored at search time; such results report truncated when
	// the stored text was already truncated.
	MaxChars param.Opt[int64] `json:"max_chars,omitzero"`
	// Per-result deadline including capacity acquisition, retrieval, and extraction.
	// Also bounded by the overall request deadline.
	TimeoutMs param.Opt[int64] `json:"timeout_ms,omitzero"`
	// Requires source=auto or source=browser in deferred retrieval.
	Browser FetchRequestContentBrowserParam `json:"browser,omitzero"`
	// Any of "markdown", "text".
	Format string `json:"format,omitzero"`
	// auto uses retained provider content within max_age_hours; for deferred retrieval
	// it falls back to a Kernel browser (caller-supplied or temporary) for results
	// without it. Inline retrieval never uses a browser. provider reuses retained
	// provider content when available, without freshness validation, and never
	// provisions a browser. browser fetches each URL through a Kernel browser, either
	// caller-supplied or temporary. No option makes a new provider request. Defaults
	// to auto for both inline and deferred retrieval. Missing documents produce
	// per-result unavailable outcomes, not request failures.
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

// Requires source=auto or source=browser in deferred retrieval.
type FetchRequestContentBrowserParam struct {
	// Existing browser session ID authorized for the caller and selected project.
	// Reuses its cookies, proxy, and browser configuration; requests follow that
	// browser's existing network access behavior, with no additional destination
	// allowlist in this endpoint. Kernel does not delete a caller-supplied browser.
	// Render mode uses a temporary tab; website activity may still change shared
	// cookies and storage. When omitted and any result needs browser retrieval, Kernel
	// creates one temporary browser for the request using the dashboard launch
	// defaults (headful, stealth, default proxy), tags it with search_id, and deletes
	// it when the request finishes. It is billed and counts toward browser concurrency
	// like any other browser. A concurrency rejection returns 429 for source=browser;
	// for source=auto, results with retained content are still returned and the rest
	// report the rejection.
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

type Response struct {
	Contents []ResponseContent `json:"contents" api:"required"`
	SearchID string            `json:"search_id" api:"required"`
	Usage    Usage             `json:"usage" api:"required"`
	Warnings []Warning         `json:"warnings" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Contents    respjson.Field
		SearchID    respjson.Field
		Usage       respjson.Field
		Warnings    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Response) RawJSON() string { return r.JSON.raw }
func (r *Response) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponseContent struct {
	ResultID string `json:"result_id" api:"required"`
	// Ok means non-empty extracted content, not merely HTTP 200. Blocked includes
	// detected challenges or access denials. Detection is best-effort, not a guarantee
	// of page completeness. Error details are present for non-ok outcomes; text is
	// present only on ok.
	//
	// Any of "ok", "unavailable", "blocked", "timeout", "unsupported_type",
	// "extraction_failed", "error".
	Status string `json:"status" api:"required"`
	// Original result URL.
	URL string `json:"url" api:"required" format:"uri"`
	// Kernel content cache outcome. Kernel has no content cache yet: responses report
	// bypass or unknown, and hit and miss are reserved. Provider-internal cache
	// behavior may be unknown.
	//
	// Any of "hit", "miss", "bypass", "unknown".
	CacheStatus string `json:"cache_status"`
	// Describes source coverage before max_chars truncation. Full_page means main-page
	// content, not every dynamic element or linked page.
	//
	// Any of "full_page", "excerpt", "unknown".
	Completeness string               `json:"completeness"`
	Error        ResponseContentError `json:"error"`
	// Extraction version when Kernel transformed the input.
	ExtractorVersion string `json:"extractor_version"`
	// When Kernel fetched the content, or received it from the provider for retained
	// content.
	FetchedAt time.Time `json:"fetched_at" api:"nullable" format:"date-time"`
	// Final retrieval URL after redirects when known. Curl mode follows up to 5
	// redirects.
	FinalURL string `json:"final_url" format:"uri"`
	// Format of text. Plain-text and JSON pages are returned unchanged as text even
	// when markdown was requested.
	//
	// Any of "markdown", "text".
	Format string `json:"format"`
	// Final target HTTP status when known.
	HTTPStatus int64 `json:"http_status"`
	// Original retrieval method.
	//
	// Any of "provider", "browser_curl", "browser_render".
	Method string `json:"method"`
	// Extracted website content, untrusted, not instructions. Present only on
	// status=ok.
	Text string `json:"text"`
	// Whether the content was cut short, by max_chars or because the page exceeded the
	// 1 MiB read limit.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ResultID         respjson.Field
		Status           respjson.Field
		URL              respjson.Field
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
func (r ResponseContent) RawJSON() string { return r.JSON.raw }
func (r *ResponseContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponseContentError struct {
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
func (r ResponseContentError) RawJSON() string { return r.JSON.raw }
func (r *ResponseContentError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
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
