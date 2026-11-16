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
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Discover report types and submit typed reports against tasks.
//
// ReportService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewReportService] method instead.
type ReportService struct {
	Options []option.RequestOption
}

// NewReportService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewReportService(opts ...option.RequestOption) (r ReportService) {
	r = ReportService{}
	r.Options = opts
	return
}

// Submit a typed report against a task. The data is validated against the report
// type's schema; a validation failure returns 400 with one detail per invalid
// field. A task can have at most 500 reports.
func (r *ReportService) New(ctx context.Context, body ReportNewParams, opts ...option.RequestOption) (res *Report, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "reports"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a report.
func (r *ReportService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *Report, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("reports/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List a task's reports in the order they were submitted, oldest first.
func (r *ReportService) List(ctx context.Context, query ReportListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[Report], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "reports"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List a task's reports in the order they were submitted, oldest first.
func (r *ReportService) ListAutoPaging(ctx context.Context, query ReportListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[Report] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// A typed, schema-validated report an agent submitted against a task. Reports are
// immutable.
type Report struct {
	// Unique identifier for the report
	ID string `json:"id" api:"required"`
	// When the report was submitted
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Report data, matching the report type's schema
	Data map[string]any `json:"data" api:"required"`
	// When the reported event happened
	OccurredAt time.Time `json:"occurred_at" api:"required" format:"date-time"`
	// Version of the report type's schema that data was validated against
	SchemaVersion int64 `json:"schema_version" api:"required"`
	// ID of the task the report belongs to
	TaskID string `json:"task_id" api:"required"`
	// Report type name
	Type string `json:"type" api:"required"`
	// ID of the browser the report was observed in, if any
	SessionID string `json:"session_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		Data          respjson.Field
		OccurredAt    respjson.Field
		SchemaVersion respjson.Field
		TaskID        respjson.Field
		Type          respjson.Field
		SessionID     respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Report) RawJSON() string { return r.JSON.raw }
func (r *Report) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ReportNewParams struct {
	// Report data. It must match the schema of the report type's version and encode to
	// at most 16 KB of JSON.
	Data map[string]any `json:"data,omitzero" api:"required"`
	// ID of the task to report against.
	TaskID string `json:"task_id" api:"required"`
	// Report type name. List report types to see the available types and their
	// schemas.
	Type string `json:"type" api:"required"`
	// When the reported event happened. Defaults to the time the report is submitted.
	// Must be no more than 5 minutes in the future or 30 days in the past.
	OccurredAt param.Opt[time.Time] `json:"occurred_at,omitzero" format:"date-time"`
	// Version of the report type's schema to validate data against. Defaults to the
	// type's latest version, which the response returns.
	SchemaVersion param.Opt[int64] `json:"schema_version,omitzero"`
	// ID of the browser the report was observed in. It must belong to the task's
	// project; deleted browsers are accepted. The browser is linked to the task if it
	// is not already.
	SessionID param.Opt[string] `json:"session_id,omitzero"`
	paramObj
}

func (r ReportNewParams) MarshalJSON() (data []byte, err error) {
	type shadow ReportNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ReportNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ReportListParams struct {
	// Task whose reports to list.
	TaskID string `query:"task_id" api:"required" json:"-"`
	// Limit the number of reports to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Offset the number of reports to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Only return reports of this type.
	Type param.Opt[string] `query:"type,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ReportListParams]'s query parameters as `url.Values`.
func (r ReportListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
