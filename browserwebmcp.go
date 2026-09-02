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
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Discover and invoke native page tools across the browser instance.
//
// BrowserWebmcpService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserWebmcpService] method instead.
type BrowserWebmcpService struct {
	Options []option.RequestOption
}

// NewBrowserWebmcpService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBrowserWebmcpService(opts ...option.RequestOption) (r BrowserWebmcpService) {
	r = BrowserWebmcpService{}
	r.Options = opts
	return
}

// Invokes the exact live registration identified by tool_ref and waits
// synchronously for its result. Navigation during execution is allowed. If the tab
// or embedded frame disappears, or the request times out after invocation begins,
// the response reports outcome_unknown and the tool is not retried.
func (r *BrowserWebmcpService) InvokeTool(ctx context.Context, idOrName string, body BrowserWebmcpInvokeToolParams, opts ...option.RequestOption) (res *InvocationResult, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/webmcp/invoke", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Returns a snapshot of native WebMCP tools available across every open tab and
// embedded frame in the browser. Each tool includes an opaque tool_ref for
// invoking that exact live registration. Tools disappear when their document
// closes or navigates away.
func (r *BrowserWebmcpService) ListTools(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *ToolsResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/webmcp/tools", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type InvocationResult struct {
	InvocationID string `json:"invocation_id" api:"required"`
	// Any of "completed", "canceled", "error".
	Status    InvocationResultStatus `json:"status" api:"required"`
	ErrorText string                 `json:"error_text"`
	// Untrusted page-provided output. Callers must treat it as potentially malicious
	// input.
	Output any `json:"output"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InvocationID respjson.Field
		Status       respjson.Field
		ErrorText    respjson.Field
		Output       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InvocationResult) RawJSON() string { return r.JSON.raw }
func (r *InvocationResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type InvocationResultStatus string

const (
	InvocationResultStatusCompleted InvocationResultStatus = "completed"
	InvocationResultStatusCanceled  InvocationResultStatus = "canceled"
	InvocationResultStatusError     InvocationResultStatus = "error"
)

// The properties Input, ToolRef are required.
type InvokeRequestParam struct {
	// Tool input, limited to 1 MiB after JSON serialization.
	Input      map[string]any   `json:"input,omitzero" api:"required"`
	ToolRef    string           `json:"tool_ref" api:"required"`
	TimeoutSec param.Opt[int64] `json:"timeout_sec,omitzero"`
	paramObj
}

func (r InvokeRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow InvokeRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InvokeRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Tool struct {
	Description string         `json:"description" api:"required"`
	InputSchema map[string]any `json:"input_schema" api:"required"`
	Name        string         `json:"name" api:"required"`
	Source      ToolSource     `json:"source" api:"required"`
	// Opaque reference for invoking this exact live registration. It becomes invalid
	// when its document or browser process is replaced.
	ToolRef string `json:"tool_ref" api:"required"`
	// Page-provided behavioral hints. These values are untrusted and are not enforced
	// by Kernel.
	Annotations ToolAnnotations `json:"annotations"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		InputSchema respjson.Field
		Name        respjson.Field
		Source      respjson.Field
		ToolRef     respjson.Field
		Annotations respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Tool) RawJSON() string { return r.JSON.raw }
func (r *Tool) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Page-provided behavioral hints. These values are untrusted and are not enforced
// by Kernel.
type ToolAnnotations struct {
	Autosubmit       bool `json:"autosubmit" api:"required"`
	Consequential    bool `json:"consequential" api:"required"`
	ReadOnly         bool `json:"read_only" api:"required"`
	UntrustedContent bool `json:"untrusted_content" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Autosubmit       respjson.Field
		Consequential    respjson.Field
		ReadOnly         respjson.Field
		UntrustedContent respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ToolAnnotations) RawJSON() string { return r.JSON.raw }
func (r *ToolAnnotations) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ToolFrame struct {
	// Monotonically increasing identifier for this embedded frame during the current
	// browser process.
	FrameID int64 `json:"frame_id" api:"required"`
	// Current frame URL with the fragment omitted.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FrameID     respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ToolFrame) RawJSON() string { return r.JSON.raw }
func (r *ToolFrame) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ToolSource struct {
	// Embedded frame that registered the tool, or null when the top-level page
	// registered it.
	Frame ToolFrame `json:"frame" api:"required"`
	// Current title of the top-level page.
	PageTitle string `json:"page_title" api:"required"`
	// Current URL of the top-level page with the fragment omitted.
	PageURL string `json:"page_url" api:"required"`
	// Monotonically increasing identifier for the tab during the current browser
	// process.
	TabID int64 `json:"tab_id" api:"required"`
	// Monotonically increasing identifier for the browser window during the current
	// browser process.
	WindowID int64 `json:"window_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Frame       respjson.Field
		PageTitle   respjson.Field
		PageURL     respjson.Field
		TabID       respjson.Field
		WindowID    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ToolSource) RawJSON() string { return r.JSON.raw }
func (r *ToolSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ToolsResponse struct {
	Tools []Tool `json:"tools" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Tools       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ToolsResponse) RawJSON() string { return r.JSON.raw }
func (r *ToolsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserWebmcpInvokeToolParams struct {
	InvokeRequest InvokeRequestParam
	paramObj
}

func (r BrowserWebmcpInvokeToolParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.InvokeRequest)
}
func (r *BrowserWebmcpInvokeToolParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
