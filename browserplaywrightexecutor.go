// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Execute Playwright code against the browser instance and manage the executors it
// runs in.
//
// BrowserPlaywrightExecutorService contains methods and other services that help
// with interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserPlaywrightExecutorService] method instead.
type BrowserPlaywrightExecutorService struct {
	Options []option.RequestOption
}

// NewBrowserPlaywrightExecutorService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBrowserPlaywrightExecutorService(opts ...option.RequestOption) (r BrowserPlaywrightExecutorService) {
	r = BrowserPlaywrightExecutorService{}
	r.Options = opts
	return
}

// Lists the default executor first, then the named executors created by POST
// /browsers/{id_or_name}/playwright/execute. Each entry reports whether a call is
// running on it and, for named executors, the target ID and URL of the tab it
// owns. Returns 404 for a browser whose image predates executors.
func (r *BrowserPlaywrightExecutorService) List(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *ExecutorList, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/playwright/executors", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Stops the executor's process. A call running on it fails with an error saying
// the executor was deleted. By default the executor's tab is closed too. The name
// can be reused; the next call with it creates a new executor and tab.
//
// The default executor is restarted instead of removed: its process is stopped, a
// call running on it fails with an error saying the executor was restarted, and
// queued and later calls run on a new process. It owns no tab, so 'close_tab' has
// no effect on it.
func (r *BrowserPlaywrightExecutorService) Delete(ctx context.Context, name string, params BrowserPlaywrightExecutorDeleteParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if params.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return err
	}
	if name == "" {
		err = errors.New("missing required name parameter")
		return err
	}
	path := fmt.Sprintf("browsers/%s/playwright/executors/%s", params.IDOrName, name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, params, nil, opts...)
	return err
}

// A Playwright executor on the browser
type Executor struct {
	// Whether a call is running on the executor
	Busy bool `json:"busy" api:"required"`
	// When the executor was created
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// When the most recent call on the executor started
	LastUsedAt time.Time `json:"last_used_at" api:"required" format:"date-time"`
	// Name of a Playwright executor. Calls with the same name run in the same
	// executor, one at a time; the first call with a new name creates it. Calls on
	// different executors run concurrently. 'default' names the executor that runs
	// calls without a name.
	Name string `json:"name" api:"required"`
	// CDP page target ID of the executor's tab, once it has one. The default executor
	// owns no tab.
	TargetID string `json:"target_id"`
	// Current URL of the executor's tab, when it is open
	URL string `json:"url"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Busy        respjson.Field
		CreatedAt   respjson.Field
		LastUsedAt  respjson.Field
		Name        respjson.Field
		TargetID    respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Executor) RawJSON() string { return r.JSON.raw }
func (r *Executor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The default executor followed by the named executors
type ExecutorList struct {
	Executors []Executor `json:"executors" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Executors   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExecutorList) RawJSON() string { return r.JSON.raw }
func (r *ExecutorList) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPlaywrightExecutorDeleteParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`
	// Close the executor's tab. Defaults to true.
	CloseTab param.Opt[bool] `query:"close_tab,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserPlaywrightExecutorDeleteParams]'s query parameters
// as `url.Values`.
func (r BrowserPlaywrightExecutorDeleteParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
