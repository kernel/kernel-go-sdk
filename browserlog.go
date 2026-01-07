// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/onkernel/kernel-go-sdk/internal/apiquery"
	"github.com/onkernel/kernel-go-sdk/internal/requestconfig"
	"github.com/onkernel/kernel-go-sdk/option"
	"github.com/onkernel/kernel-go-sdk/packages/param"
	"github.com/onkernel/kernel-go-sdk/packages/ssestream"
	"github.com/onkernel/kernel-go-sdk/shared"
)

// BrowserLogService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserLogService] method instead.
type BrowserLogService struct {
	Options []option.RequestOption
}

// NewBrowserLogService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBrowserLogService(opts ...option.RequestOption) (r BrowserLogService) {
	r = BrowserLogService{}
	r.Options = opts
	return
}

// Stream log files on the browser instance via SSE
func (r *BrowserLogService) StreamStreaming(ctx context.Context, id string, query BrowserLogStreamParams, opts ...option.RequestOption) (stream *ssestream.Stream[shared.LogEvent]) {
	var (
		raw *http.Response
		err error
	)
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "text/event-stream")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/logs/stream", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &raw, opts...)
	return ssestream.NewStream[shared.LogEvent](ssestream.NewDecoder(raw), err)
}

type BrowserLogStreamParams struct {
	// Any of "path", "supervisor".
	Source BrowserLogStreamParamsSource `query:"source,omitzero,required" json:"-"`
	Follow param.Opt[bool]              `query:"follow,omitzero" json:"-"`
	// only required if source is path
	Path param.Opt[string] `query:"path,omitzero" json:"-"`
	// only required if source is supervisor
	SupervisorProcess param.Opt[string] `query:"supervisor_process,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserLogStreamParams]'s query parameters as `url.Values`.
func (r BrowserLogStreamParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BrowserLogStreamParamsSource string

const (
	BrowserLogStreamParamsSourcePath       BrowserLogStreamParamsSource = "path"
	BrowserLogStreamParamsSourceSupervisor BrowserLogStreamParamsSource = "supervisor"
)
