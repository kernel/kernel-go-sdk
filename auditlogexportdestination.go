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
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Read audit log records for the authenticated organization.
//
// AuditLogExportDestinationService contains methods and other services that help
// with interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAuditLogExportDestinationService] method instead.
type AuditLogExportDestinationService struct {
	Options []option.RequestOption
}

// NewAuditLogExportDestinationService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAuditLogExportDestinationService(opts ...option.RequestOption) (r AuditLogExportDestinationService) {
	r = AuditLogExportDestinationService{}
	r.Options = opts
	return
}

// Create a paused destination. Activate it with a status update once the
// destination test passes. Requires an active Enterprise plan.
func (r *AuditLogExportDestinationService) New(ctx context.Context, body AuditLogExportDestinationNewParams, opts ...option.RequestOption) (res *AuditLogExportDestination, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "audit-logs/export/destinations"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve details for a single audit log export destination by its ID.
func (r *AuditLogExportDestinationService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *AuditLogExportDestination, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("audit-logs/export/destinations/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Apply a partial update to a destination. Requires an active Enterprise plan.
// Returns 409 when the destination was changed concurrently, because the merged
// configuration this request validated is no longer the one that would be stored;
// retry against fresh state. Pausing prevents new delivery attempts, but an S3
// upload already in progress may complete after the response.
func (r *AuditLogExportDestinationService) Update(ctx context.Context, id string, body AuditLogExportDestinationUpdateParams, opts ...option.RequestOption) (res *AuditLogExportDestination, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("audit-logs/export/destinations/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List audit log export destinations for the organization with pagination support.
func (r *AuditLogExportDestinationService) List(ctx context.Context, query AuditLogExportDestinationListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[AuditLogExportDestination], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "audit-logs/export/destinations"
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

// List audit log export destinations for the organization with pagination support.
func (r *AuditLogExportDestinationService) ListAutoPaging(ctx context.Context, query AuditLogExportDestinationListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[AuditLogExportDestination] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Soft delete the destination and prevent new delivery attempts. An S3 upload
// already in progress may complete after the response.
func (r *AuditLogExportDestinationService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("audit-logs/export/destinations/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Verify the destination is writable by assuming the configured role and uploading
// a temporary probe object with the same request metadata as a real delivery.
// Requires an active Enterprise plan.
func (r *AuditLogExportDestinationService) Test(ctx context.Context, id string, opts ...option.RequestOption) (res *AuditLogExportDestinationTestResult, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("audit-logs/export/destinations/%s/test", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// An organization-scoped audit log export destination.
//
// Delivery is at-least-once for rows visible when their window is committed: a
// delivery that is retried rewrites the same object, and the same `event_id` can
// appear in more than one object, so consumers must deduplicate on `event_id`.
// Each event-time window is held for ten minutes before it commits; a row that
// becomes visible after its window is committed may not be delivered.
//
// Objects are written as
// `<prefix>/destination_id=<destination>/org_id=<org>/date=<YYYY-MM-DD>/hour=<HH>/<window>-<chunk>.jsonl.gz`,
// where `date` and `hour` are the UTC calendar hour that fully contains every row
// in the object, so the layout is safe to register as a Hive-partitioned table.
// The object name is derived from the rows it holds, so a retried delivery
// rewrites its own object.
type AuditLogExportDestination struct {
	ID                  string    `json:"id" api:"required"`
	Bucket              string    `json:"bucket" api:"required"`
	ConsecutiveFailures int64     `json:"consecutive_failures" api:"required"`
	CreatedAt           time.Time `json:"created_at" api:"required" format:"date-time"`
	ExternalID          string    `json:"external_id" api:"required"`
	// Any of "jsonl.gz".
	Format AuditLogExportDestinationFormat `json:"format" api:"required"`
	// The Kernel role that assumes `role_arn` in your account to deliver logs. Allow
	// this role as the principal in your role's trust policy, and require
	// `external_id` as the `sts:ExternalId` condition.
	//
	// Recreating a destination issues a new `external_id`, which the trust policy has
	// to be updated to match.
	KernelRoleArn string `json:"kernel_role_arn" api:"required"`
	Prefix        string `json:"prefix" api:"required"`
	Region        string `json:"region" api:"required"`
	RoleArn       string `json:"role_arn" api:"required"`
	// Pausing prevents new delivery attempts. An S3 upload already in progress may
	// complete after the pause response; its rows can appear again after the
	// destination is resumed.
	//
	// Any of "active", "paused".
	Status AuditLogExportDestinationStatus `json:"status" api:"required"`
	// Any of "s3".
	Type      AuditLogExportDestinationType `json:"type" api:"required"`
	UpdatedAt time.Time                     `json:"updated_at" api:"required" format:"date-time"`
	KmsKeyID  string                        `json:"kms_key_id"`
	// Sanitized description of the most recent delivery failure.
	LastError   string    `json:"last_error"`
	LastErrorAt time.Time `json:"last_error_at" format:"date-time"`
	// Opaque, versioned checkpoint for forward-only continuous export. This value is
	// not compatible with audit-log list page tokens.
	//
	// Delivery starts at the moment the destination is activated, so events recorded
	// before that are not delivered. Pausing stops delivery and resuming starts again
	// from the time of the resume: events recorded while a destination was paused are
	// never exported, and pausing is not a way to defer delivery.
	LastExportedCursor string    `json:"last_exported_cursor"`
	LastSuccessAt      time.Time `json:"last_success_at" format:"date-time"`
	NextAttemptAt      time.Time `json:"next_attempt_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		Bucket              respjson.Field
		ConsecutiveFailures respjson.Field
		CreatedAt           respjson.Field
		ExternalID          respjson.Field
		Format              respjson.Field
		KernelRoleArn       respjson.Field
		Prefix              respjson.Field
		Region              respjson.Field
		RoleArn             respjson.Field
		Status              respjson.Field
		Type                respjson.Field
		UpdatedAt           respjson.Field
		KmsKeyID            respjson.Field
		LastError           respjson.Field
		LastErrorAt         respjson.Field
		LastExportedCursor  respjson.Field
		LastSuccessAt       respjson.Field
		NextAttemptAt       respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuditLogExportDestination) RawJSON() string { return r.JSON.raw }
func (r *AuditLogExportDestination) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuditLogExportDestinationFormat string

const (
	AuditLogExportDestinationFormatJSONLGz AuditLogExportDestinationFormat = "jsonl.gz"
)

// Pausing prevents new delivery attempts. An S3 upload already in progress may
// complete after the pause response; its rows can appear again after the
// destination is resumed.
type AuditLogExportDestinationStatus string

const (
	AuditLogExportDestinationStatusActive AuditLogExportDestinationStatus = "active"
	AuditLogExportDestinationStatusPaused AuditLogExportDestinationStatus = "paused"
)

type AuditLogExportDestinationType string

const (
	AuditLogExportDestinationTypeS3 AuditLogExportDestinationType = "s3"
)

type AuditLogExportDestinationTestResult struct {
	// Any of "assume_role", "put_object", "complete".
	Stage   AuditLogExportDestinationTestResultStage `json:"stage" api:"required"`
	Success bool                                     `json:"success" api:"required"`
	Error   AuditLogExportDestinationTestResultError `json:"error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Stage       respjson.Field
		Success     respjson.Field
		Error       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuditLogExportDestinationTestResult) RawJSON() string { return r.JSON.raw }
func (r *AuditLogExportDestinationTestResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuditLogExportDestinationTestResultStage string

const (
	AuditLogExportDestinationTestResultStageAssumeRole AuditLogExportDestinationTestResultStage = "assume_role"
	AuditLogExportDestinationTestResultStagePutObject  AuditLogExportDestinationTestResultStage = "put_object"
	AuditLogExportDestinationTestResultStageComplete   AuditLogExportDestinationTestResultStage = "complete"
)

type AuditLogExportDestinationTestResultError struct {
	// Any of "assume_role_failed", "put_object_failed".
	Code    string `json:"code" api:"required"`
	Message string `json:"message" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuditLogExportDestinationTestResultError) RawJSON() string { return r.JSON.raw }
func (r *AuditLogExportDestinationTestResultError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Bucket, Format, Prefix, Region, RoleArn, Type are required.
type CreateAuditLogExportDestinationRequestParam struct {
	Bucket string `json:"bucket" api:"required"`
	// Any of "jsonl.gz".
	Format  CreateAuditLogExportDestinationRequestFormat `json:"format,omitzero" api:"required"`
	Prefix  string                                       `json:"prefix" api:"required"`
	Region  string                                       `json:"region" api:"required"`
	RoleArn string                                       `json:"role_arn" api:"required"`
	// Any of "s3".
	Type     CreateAuditLogExportDestinationRequestType `json:"type,omitzero" api:"required"`
	KmsKeyID param.Opt[string]                          `json:"kms_key_id,omitzero"`
	paramObj
}

func (r CreateAuditLogExportDestinationRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateAuditLogExportDestinationRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CreateAuditLogExportDestinationRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateAuditLogExportDestinationRequestFormat string

const (
	CreateAuditLogExportDestinationRequestFormatJSONLGz CreateAuditLogExportDestinationRequestFormat = "jsonl.gz"
)

type CreateAuditLogExportDestinationRequestType string

const (
	CreateAuditLogExportDestinationRequestTypeS3 CreateAuditLogExportDestinationRequestType = "s3"
)

type UpdateAuditLogExportDestinationRequestParam struct {
	Bucket param.Opt[string] `json:"bucket,omitzero"`
	// KMS key ID, alias, or ARN. Set to an empty string to remove the configured KMS
	// key; omit or send null to leave unchanged.
	KmsKeyID param.Opt[string] `json:"kms_key_id,omitzero"`
	Prefix   param.Opt[string] `json:"prefix,omitzero"`
	Region   param.Opt[string] `json:"region,omitzero"`
	RoleArn  param.Opt[string] `json:"role_arn,omitzero"`
	// Any of "active", "paused".
	Status UpdateAuditLogExportDestinationRequestStatus `json:"status,omitzero"`
	paramObj
}

func (r UpdateAuditLogExportDestinationRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow UpdateAuditLogExportDestinationRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *UpdateAuditLogExportDestinationRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type UpdateAuditLogExportDestinationRequestStatus string

const (
	UpdateAuditLogExportDestinationRequestStatusActive UpdateAuditLogExportDestinationRequestStatus = "active"
	UpdateAuditLogExportDestinationRequestStatusPaused UpdateAuditLogExportDestinationRequestStatus = "paused"
)

type AuditLogExportDestinationNewParams struct {
	CreateAuditLogExportDestinationRequest CreateAuditLogExportDestinationRequestParam
	paramObj
}

func (r AuditLogExportDestinationNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.CreateAuditLogExportDestinationRequest)
}
func (r *AuditLogExportDestinationNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuditLogExportDestinationUpdateParams struct {
	UpdateAuditLogExportDestinationRequest UpdateAuditLogExportDestinationRequestParam
	paramObj
}

func (r AuditLogExportDestinationUpdateParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.UpdateAuditLogExportDestinationRequest)
}
func (r *AuditLogExportDestinationUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuditLogExportDestinationListParams struct {
	// Limit the number of destinations to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Offset the number of destinations to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AuditLogExportDestinationListParams]'s query parameters as
// `url.Values`.
func (r AuditLogExportDestinationListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
