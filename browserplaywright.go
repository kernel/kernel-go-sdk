// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/onkernel/kernel-go-sdk/internal/apijson"
	"github.com/onkernel/kernel-go-sdk/internal/requestconfig"
	"github.com/onkernel/kernel-go-sdk/option"
	"github.com/onkernel/kernel-go-sdk/packages/param"
	"github.com/onkernel/kernel-go-sdk/packages/respjson"
)

// BrowserPlaywrightService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserPlaywrightService] method instead.
type BrowserPlaywrightService struct {
	Options []option.RequestOption
}

// NewBrowserPlaywrightService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBrowserPlaywrightService(opts ...option.RequestOption) (r BrowserPlaywrightService) {
	r = BrowserPlaywrightService{}
	r.Options = opts
	return
}

// Execute arbitrary Playwright code in a fresh execution context against the
// browser. The code runs in the same VM as the browser, minimizing latency and
// maximizing throughput. It has access to 'page', 'context', and 'browser'
// variables. It can `return` a value, and this value is returned in the response.
func (r *BrowserPlaywrightService) Execute(ctx context.Context, id string, body BrowserPlaywrightExecuteParams, opts ...option.RequestOption) (res *BrowserPlaywrightExecuteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/playwright/execute", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Result of Playwright code execution
type BrowserPlaywrightExecuteResponse struct {
	// Whether the code executed successfully
	Success bool `json:"success,required"`
	// Error message if execution failed
	Error string `json:"error"`
	// The value returned by the code (if any)
	Result any `json:"result"`
	// Standard error from the execution
	Stderr string `json:"stderr"`
	// Standard output from the execution
	Stdout string `json:"stdout"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Success     respjson.Field
		Error       respjson.Field
		Result      respjson.Field
		Stderr      respjson.Field
		Stdout      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPlaywrightExecuteResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserPlaywrightExecuteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPlaywrightExecuteParams struct {
	// TypeScript/JavaScript code to execute. The code has access to 'page', 'context',
	// and 'browser' variables. It runs within a function, so you can use a return
	// statement at the end to return a value. This value is returned as the `result`
	// property in the response. Example: "await page.goto('https://example.com');
	// return await page.title();"
	Code string `json:"code,required"`
	// Maximum execution time in seconds. Default is 60.
	TimeoutSec param.Opt[int64] `json:"timeout_sec,omitzero"`
	paramObj
}

func (r BrowserPlaywrightExecuteParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPlaywrightExecuteParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPlaywrightExecuteParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
