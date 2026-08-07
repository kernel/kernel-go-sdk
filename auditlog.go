// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
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

// Read audit log records for the authenticated organization.
//
// AuditLogService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAuditLogService] method instead.
type AuditLogService struct {
	Options []option.RequestOption
	// Read audit log records for the authenticated organization.
	ExportDestinations AuditLogExportDestinationService
}

// NewAuditLogService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewAuditLogService(opts ...option.RequestOption) (r AuditLogService) {
	r = AuditLogService{}
	r.Options = opts
	r.ExportDestinations = NewAuditLogExportDestinationService(opts...)
	return
}

// API for searching audit logs. Limited to at most 30 day search, returns up to
// 100 records per page. Not recommended for bulk export.
func (r *AuditLogService) List(ctx context.Context, query AuditLogListParams, opts ...option.RequestOption) (res *pagination.PageTokenPagination[AuditLogEntry], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "audit-logs"
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

// API for searching audit logs. Limited to at most 30 day search, returns up to
// 100 records per page. Not recommended for bulk export.
func (r *AuditLogService) ListAutoPaging(ctx context.Context, query AuditLogListParams, opts ...option.RequestOption) *pagination.PageTokenPaginationAutoPager[AuditLogEntry] {
	return pagination.NewPageTokenPaginationAutoPager(r.List(ctx, query, opts...))
}

// Download an organization's audit log records for a time range as a file, for
// archival, compliance, or offline analysis. For interactive browsing, use GET
// /audit-logs.
func (r *AuditLogService) ExportChunk(ctx context.Context, query AuditLogExportChunkParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/octet-stream")}, opts...)
	path := "audit-logs/export/chunk"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type AuditLogEntry struct {
	// Authentication strategy used for the request.
	AuthStrategy string `json:"auth_strategy" api:"required"`
	// Client IP address.
	ClientIP string `json:"client_ip" api:"required"`
	// Request host.
	Domain string `json:"domain" api:"required"`
	// Request duration in milliseconds.
	DurationMs int64 `json:"duration_ms" api:"required"`
	// Email of the authenticated user at request time, if any.
	Email string `json:"email" api:"required"`
	// HTTP method.
	Method string `json:"method" api:"required"`
	// Request path.
	Path string `json:"path" api:"required"`
	// Matched API route pattern, if available.
	Route string `json:"route" api:"required"`
	// HTTP response status code.
	Status int64 `json:"status" api:"required"`
	// UTC time when the request was received.
	Timestamp time.Time `json:"timestamp" api:"required" format:"date-time"`
	// User agent header.
	UserAgent string `json:"user_agent" api:"required"`
	// ID of the authenticated user, if any.
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AuthStrategy respjson.Field
		ClientIP     respjson.Field
		Domain       respjson.Field
		DurationMs   respjson.Field
		Email        respjson.Field
		Method       respjson.Field
		Path         respjson.Field
		Route        respjson.Field
		Status       respjson.Field
		Timestamp    respjson.Field
		UserAgent    respjson.Field
		UserID       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuditLogEntry) RawJSON() string { return r.JSON.raw }
func (r *AuditLogEntry) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuditLogListParams struct {
	// Upper bound (exclusive) for the audit record timestamp.
	End time.Time `query:"end" api:"required" format:"date-time" json:"-"`
	// Lower bound (inclusive) for the audit record timestamp.
	Start time.Time `query:"start" api:"required" format:"date-time" json:"-"`
	// Filter by authentication strategy.
	AuthStrategy param.Opt[string] `query:"auth_strategy,omitzero" json:"-"`
	// Maximum number of results to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter by HTTP method.
	Method param.Opt[string] `query:"method,omitzero" json:"-"`
	// Opaque page token from X-Next-Page-Token for the next page of older records.
	PageToken param.Opt[string] `query:"page_token,omitzero" json:"-"`
	// Free-text search over path, user ID, email, client IP, and status.
	Search param.Opt[string] `query:"search,omitzero" json:"-"`
	// Filter by service name.
	Service param.Opt[string] `query:"service,omitzero" json:"-"`
	// Filter out results by HTTP method.
	ExcludeMethod []string `query:"exclude_method,omitzero" json:"-"`
	// Additional user IDs to OR into free-text search.
	SearchUserID []string `query:"search_user_id,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AuditLogListParams]'s query parameters as `url.Values`.
func (r AuditLogListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AuditLogExportChunkParams struct {
	// Upper bound (exclusive) for the audit record timestamp.
	End time.Time `query:"end" api:"required" format:"date-time" json:"-"`
	// Lower bound (inclusive) for the audit record timestamp.
	Start time.Time `query:"start" api:"required" format:"date-time" json:"-"`
	// Filter by authentication strategy.
	AuthStrategy param.Opt[string] `query:"auth_strategy,omitzero" json:"-"`
	// Opaque cursor from X-Next-Cursor for the next chunk of older records.
	Cursor param.Opt[string] `query:"cursor,omitzero" json:"-"`
	// Maximum number of records to return in this chunk.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter by HTTP method.
	Method param.Opt[string] `query:"method,omitzero" json:"-"`
	// Free-text search over path, user ID, email, client IP, and status.
	Search param.Opt[string] `query:"search,omitzero" json:"-"`
	// Filter by service name.
	Service param.Opt[string] `query:"service,omitzero" json:"-"`
	// Filter out results by HTTP method.
	ExcludeMethod []string `query:"exclude_method,omitzero" json:"-"`
	// Encoding for the returned chunk.
	//
	// Any of "jsonl", "jsonl.gz".
	Format AuditLogExportChunkParamsFormat `query:"format,omitzero" json:"-"`
	// Additional user IDs to OR into free-text search.
	SearchUserID []string `query:"search_user_id,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AuditLogExportChunkParams]'s query parameters as
// `url.Values`.
func (r AuditLogExportChunkParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Encoding for the returned chunk.
type AuditLogExportChunkParamsFormat string

const (
	AuditLogExportChunkParamsFormatJSONL   AuditLogExportChunkParamsFormat = "jsonl"
	AuditLogExportChunkParamsFormatJSONLGz AuditLogExportChunkParamsFormat = "jsonl.gz"
)
