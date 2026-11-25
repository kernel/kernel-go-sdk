// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

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
// ReportTypeService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewReportTypeService] method instead.
type ReportTypeService struct {
	Options []option.RequestOption
	// Discover report types and submit typed reports against tasks.
	Versions ReportTypeVersionService
}

// NewReportTypeService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewReportTypeService(opts ...option.RequestOption) (r ReportTypeService) {
	r = ReportTypeService{}
	r.Options = opts
	r.Versions = NewReportTypeVersionService(opts...)
	return
}

// Retrieve a report type, its versions, and the JSON Schema of its latest version.
func (r *ReportTypeService) Get(ctx context.Context, type_ string, opts ...option.RequestOption) (res *ReportType, err error) {
	opts = slices.Concat(r.Options, opts)
	if type_ == "" {
		err = errors.New("missing required type parameter")
		return nil, err
	}
	path := fmt.Sprintf("report_types/%s", type_)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List the report types agents can submit for a task. Retrieve a type for the JSON
// Schema of its data.
func (r *ReportTypeService) List(ctx context.Context, query ReportTypeListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[ReportTypeSummary], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "report_types"
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

// List the report types agents can submit for a task. Retrieve a type for the JSON
// Schema of its data.
func (r *ReportTypeService) ListAutoPaging(ctx context.Context, query ReportTypeListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[ReportTypeSummary] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// A kind of report agents can submit against a task, with its versions and the
// JSON Schema of its latest version. A shipped version never changes; changes ship
// as new versions.
type ReportType struct {
	// What the report type records and when to submit it
	Description string `json:"description" api:"required"`
	// Highest version that is not retired. A report submitted without a schema_version
	// uses this version. It is 0 when every version is retired.
	LatestVersion int64 `json:"latest_version" api:"required"`
	// JSON Schema (draft 2020-12) of the latest version, which a report's data must
	// match. Its field descriptions and examples describe what to send. Null when
	// every version is retired.
	Schema map[string]any `json:"schema" api:"required"`
	// Report type name, such as task.result
	Type string `json:"type" api:"required"`
	// Every version of the type, oldest first, including retired versions.
	Versions []ReportTypeVersionSummary `json:"versions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description   respjson.Field
		LatestVersion respjson.Field
		Schema        respjson.Field
		Type          respjson.Field
		Versions      respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ReportType) RawJSON() string { return r.JSON.raw }
func (r *ReportType) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A kind of report agents can submit against a task. Retrieve the type for the
// JSON Schema of its data.
type ReportTypeSummary struct {
	// What the report type records and when to submit it
	Description string `json:"description" api:"required"`
	// Highest version that is not retired. A report submitted without a schema_version
	// uses this version. It is 0 when every version is retired.
	LatestVersion int64 `json:"latest_version" api:"required"`
	// Report type name, such as task.result
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description   respjson.Field
		LatestVersion respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ReportTypeSummary) RawJSON() string { return r.JSON.raw }
func (r *ReportTypeSummary) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// One schema version of a report type, without its schema.
type ReportTypeVersionSummary struct {
	// Whether the version is retired. Retired versions are no longer accepted for new
	// reports; existing reports keep them.
	Retired bool `json:"retired" api:"required"`
	// Schema version number, starting at 1.
	Version int64 `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Retired     respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ReportTypeVersionSummary) RawJSON() string { return r.JSON.raw }
func (r *ReportTypeVersionSummary) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ReportTypeListParams struct {
	// Task the agent is finding a report type for.
	TaskID string `query:"task_id" api:"required" json:"-"`
	// Limit the number of report types to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Offset the number of report types to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// What the agent wants to report, in its own words. Types are ranked by how
	// closely their meaning matches, most relevant first; they need not share any
	// words. Only this text affects the ranking, not the task. Without it, or if
	// ranking is unavailable, types are listed by name.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ReportTypeListParams]'s query parameters as `url.Values`.
func (r ReportTypeListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
