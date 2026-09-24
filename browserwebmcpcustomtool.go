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
// BrowserWebmcpCustomToolService contains methods and other services that help
// with interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserWebmcpCustomToolService] method instead.
type BrowserWebmcpCustomToolService struct {
	Options []option.RequestOption
}

// NewBrowserWebmcpCustomToolService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBrowserWebmcpCustomToolService(opts ...option.RequestOption) (r BrowserWebmcpCustomToolService) {
	r = BrowserWebmcpCustomToolService{}
	r.Options = opts
	return
}

// Returns every registered custom tool with its generated ID, namespace, matcher,
// and MCP tool metadata.
func (r *BrowserWebmcpCustomToolService) List(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *CustomToolsResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/webmcp/custom-tools", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Add a namespaced batch of custom tools. A custom tool can be page-backed or
// CDP-backed. Page-backed tools execute in the page via JavaScript. CDP-backed
// tools execute via CDP and can use all browser REPL tools (see `/repl`). The
// source must evaluate to a non-empty array of definitions with URL matchers, tool
// metadata (including an optional output schema), and execute functions. The batch
// is added atomically. Matchers apply to top-level documents and nested frames,
// including out-of-process iframes; each matching tool is exposed once on the
// tab's top-level document and appears in the WebMCP tool snapshot.
//
// To update one tool, list the tools, delete its ID, and add its replacement. Set
// force_overwrite_namespace to replace every existing tool in the namespace
// atomically; omitted or false adds tools without replacing existing ones.
// Existing invocations continue.
func (r *BrowserWebmcpCustomToolService) Add(ctx context.Context, idOrName string, body BrowserWebmcpCustomToolAddParams, opts ...option.RequestOption) (res *CustomToolsResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/webmcp/custom-tools", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Removes one custom tool by generated ID. An invocation already in progress is
// not canceled.
func (r *BrowserWebmcpCustomToolService) Remove(ctx context.Context, id string, body BrowserWebmcpCustomToolRemoveParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if body.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return err
	}
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("browsers/%s/webmcp/custom-tools/%s", body.IDOrName, id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// The properties Namespace, Source are required.
type AddRequestParam struct {
	Namespace string `json:"namespace" api:"required"`
	// JavaScript expression that evaluates to a non-empty array of custom tool
	// definitions. Limited to 8,000,000 bytes when UTF-8 encoded, so multi-byte
	// characters reduce the allowed character count.
	Source string `json:"source" api:"required"`
	// Atomically replace all existing tools in this namespace with this batch when
	// true.
	ForceOverwriteNamespace param.Opt[bool] `json:"force_overwrite_namespace,omitzero"`
	paramObj
}

func (r AddRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow AddRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AddRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CustomToolsResponse struct {
	Tools []Definition `json:"tools" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Tools       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CustomToolsResponse) RawJSON() string { return r.JSON.raw }
func (r *CustomToolsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Definition struct {
	ID        string `json:"id" api:"required"`
	Kind      string `json:"kind" api:"required"`
	Match     Match  `json:"match" api:"required"`
	Namespace string `json:"namespace" api:"required"`
	// Tool metadata follows the
	// [MCP Tool definition](https://modelcontextprotocol.io/specification/2025-11-25/server/tools#tool)
	// and the
	// [WebMCP RegisteredTool definition](https://webmachinelearning.github.io/webmcp/#dictdef-registeredtool).
	// outputSchema is optional for page and custom tools.
	Tool ToolMetadata `json:"tool" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Kind        respjson.Field
		Match       respjson.Field
		Namespace   respjson.Field
		Tool        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Definition) RawJSON() string { return r.JSON.raw }
func (r *Definition) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Match struct {
	URLPatterns []string `json:"url_patterns" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		URLPatterns respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Match) RawJSON() string { return r.JSON.raw }
func (r *Match) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserWebmcpCustomToolAddParams struct {
	AddRequest AddRequestParam
	paramObj
}

func (r BrowserWebmcpCustomToolAddParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.AddRequest)
}
func (r *BrowserWebmcpCustomToolAddParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserWebmcpCustomToolRemoveParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`
	paramObj
}
