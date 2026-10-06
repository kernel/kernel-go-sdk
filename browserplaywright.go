// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Execute Playwright code against the browser instance and manage the executors it
// runs in.
//
// BrowserPlaywrightService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserPlaywrightService] method instead.
type BrowserPlaywrightService struct {
	Options []option.RequestOption
	// Execute Playwright code against the browser instance and manage the executors it
	// runs in.
	Executors BrowserPlaywrightExecutorService
}

// NewBrowserPlaywrightService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBrowserPlaywrightService(opts ...option.RequestOption) (r BrowserPlaywrightService) {
	r = BrowserPlaywrightService{}
	r.Options = opts
	r.Executors = NewBrowserPlaywrightExecutorService(opts...)
	return
}

// Execute arbitrary Playwright code in a fresh execution context against the
// browser. The code runs in the same VM as the browser, minimizing latency and
// maximizing throughput. It has access to 'page', 'context', 'browser', and
// 'webmcp' variables. Use 'webmcp.listTools()' to discover browser-wide WebMCP
// tools and 'webmcp.invokeTool(toolRef, input?, { timeoutSec? })' to invoke an
// exact registration. It can `return` a value, and this value is returned in the
// response.
//
// Every call runs in an executor: a dedicated Node.js process with its own browser
// connection. Calls on different executors run concurrently; calls on the same
// executor run one at a time. A timeout, crash, or blocked event loop in one
// executor does not affect other executors. After a timeout the executor keeps its
// process and drops its browser connection, so code abandoned by the timeout
// cannot keep driving the browser. After a crash or a blocked event loop, the next
// call on that executor starts a fresh process.
//
// Calls without 'executor' run in the executor named 'default', which always
// exists and is the same as passing 'executor: "default"'. In the default
// executor, 'page' is bound to an active tab reported by Chrome. In single-window
// sessions, this is the foreground tab. When multiple browser windows are open,
// Chrome reports one active tab per window and the selected window is unspecified.
// 'context' is the BrowserContext that owns the selected page. Use
// 'browser.contexts()' to select a context or page explicitly.
//
// Pass any other name to run the call in a named executor. The first call with a
// new name creates it. Each named executor owns a tab: its first call opens a new
// background tab in the default browser context, and 'page' is bound to that tab
// on every later call while it stays open. Opening it does not change the active
// tab of an existing window. If the tab is closed, the next call opens a new one
// and reports 'tab.created: true'. Executor code can still reach other tabs
// through 'context' and 'browser'; ownership only decides what 'page' is bound to.
// Use named executors to drive several tabs of one browser in parallel.
//
// A browser can have at most 8 named executors; the default executor does not
// count. A call that would create another returns 409 with the current executors;
// delete one with DELETE /browsers/{id_or_name}/playwright/executors/{name}. Named
// executors are not removed automatically while the browser runs; when it shuts
// down, they are removed and their tabs closed.
//
// A named call to a browser whose image predates executors fails with 400 instead
// of running on the active tab; calls without 'executor' keep working on every
// image.
func (r *BrowserPlaywrightService) Execute(ctx context.Context, idOrName string, body BrowserPlaywrightExecuteParams, opts ...option.RequestOption) (res *BrowserPlaywrightExecuteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/playwright/execute", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// The tab 'page' was bound to for this call. Absent if the call failed before
// binding a tab.
type Tab struct {
	// Whether this call opened the tab. For a named executor this happens on its first
	// call and after its previous tab was closed. For the default executor it happens
	// only when the browser had no open page.
	Created bool `json:"created" api:"required"`
	// CDP page target ID of the tab
	TargetID string `json:"target_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Created     respjson.Field
		TargetID    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Tab) RawJSON() string { return r.JSON.raw }
func (r *Tab) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Result of Playwright code execution
type BrowserPlaywrightExecuteResponse struct {
	// Whether the code executed successfully
	Success bool `json:"success" api:"required"`
	// Error message if execution failed
	Error string `json:"error"`
	// The value returned by the code (if any)
	Result any `json:"result"`
	// Standard error from the execution
	Stderr string `json:"stderr"`
	// Standard output from the execution
	Stdout string `json:"stdout"`
	// The tab 'page' was bound to for this call. Absent if the call failed before
	// binding a tab.
	Tab Tab `json:"tab"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Success     respjson.Field
		Error       respjson.Field
		Result      respjson.Field
		Stderr      respjson.Field
		Stdout      respjson.Field
		Tab         respjson.Field
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
	Code string `json:"code" api:"required"`
	// Name of a Playwright executor. Calls with the same name run in the same
	// executor, one at a time; the first call with a new name creates it. Calls on
	// different executors run concurrently. 'default' names the executor that runs
	// calls without a name.
	Executor param.Opt[string] `json:"executor,omitzero"`
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
