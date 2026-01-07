// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apiform"
	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
	"github.com/kernel/kernel-go-sdk/packages/ssestream"
	"github.com/kernel/kernel-go-sdk/shared"
	"github.com/kernel/kernel-go-sdk/shared/constant"
)

// DeploymentService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDeploymentService] method instead.
type DeploymentService struct {
	Options []option.RequestOption
}

// NewDeploymentService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDeploymentService(opts ...option.RequestOption) (r DeploymentService) {
	r = DeploymentService{}
	r.Options = opts
	return
}

// Create a new deployment.
func (r *DeploymentService) New(ctx context.Context, body DeploymentNewParams, opts ...option.RequestOption) (res *DeploymentNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "deployments"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get information about a deployment's status.
func (r *DeploymentService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *DeploymentGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("deployments/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List deployments. Optionally filter by application name.
func (r *DeploymentService) List(ctx context.Context, query DeploymentListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[DeploymentListResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "deployments"
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

// List deployments. Optionally filter by application name.
func (r *DeploymentService) ListAutoPaging(ctx context.Context, query DeploymentListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[DeploymentListResponse] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Establishes a Server-Sent Events (SSE) stream that delivers real-time logs and
// status updates for a deployment. The stream terminates automatically once the
// deployment reaches a terminal state.
func (r *DeploymentService) FollowStreaming(ctx context.Context, id string, query DeploymentFollowParams, opts ...option.RequestOption) (stream *ssestream.Stream[DeploymentFollowResponseUnion]) {
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
	path := fmt.Sprintf("deployments/%s/events", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &raw, opts...)
	return ssestream.NewStream[DeploymentFollowResponseUnion](ssestream.NewDecoder(raw), err)
}

// An event representing the current state of a deployment.
type DeploymentStateEvent struct {
	// Deployment record information.
	Deployment DeploymentStateEventDeployment `json:"deployment,required"`
	// Event type identifier (always "deployment_state").
	Event constant.DeploymentState `json:"event,required"`
	// Time the state was reported.
	Timestamp time.Time `json:"timestamp,required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Deployment  respjson.Field
		Event       respjson.Field
		Timestamp   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DeploymentStateEvent) RawJSON() string { return r.JSON.raw }
func (r *DeploymentStateEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Deployment record information.
type DeploymentStateEventDeployment struct {
	// Unique identifier for the deployment
	ID string `json:"id,required"`
	// Timestamp when the deployment was created
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// Deployment region code
	Region constant.AwsUsEast1a `json:"region,required"`
	// Current status of the deployment
	//
	// Any of "queued", "in_progress", "running", "failed", "stopped".
	Status string `json:"status,required"`
	// Relative path to the application entrypoint
	EntrypointRelPath string `json:"entrypoint_rel_path"`
	// Environment variables configured for this deployment
	EnvVars map[string]string `json:"env_vars"`
	// Status reason
	StatusReason string `json:"status_reason"`
	// Timestamp when the deployment was last updated
	UpdatedAt time.Time `json:"updated_at,nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		CreatedAt         respjson.Field
		Region            respjson.Field
		Status            respjson.Field
		EntrypointRelPath respjson.Field
		EnvVars           respjson.Field
		StatusReason      respjson.Field
		UpdatedAt         respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DeploymentStateEventDeployment) RawJSON() string { return r.JSON.raw }
func (r *DeploymentStateEventDeployment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Deployment record information.
type DeploymentNewResponse struct {
	// Unique identifier for the deployment
	ID string `json:"id,required"`
	// Timestamp when the deployment was created
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// Deployment region code
	Region constant.AwsUsEast1a `json:"region,required"`
	// Current status of the deployment
	//
	// Any of "queued", "in_progress", "running", "failed", "stopped".
	Status DeploymentNewResponseStatus `json:"status,required"`
	// Relative path to the application entrypoint
	EntrypointRelPath string `json:"entrypoint_rel_path"`
	// Environment variables configured for this deployment
	EnvVars map[string]string `json:"env_vars"`
	// Status reason
	StatusReason string `json:"status_reason"`
	// Timestamp when the deployment was last updated
	UpdatedAt time.Time `json:"updated_at,nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		CreatedAt         respjson.Field
		Region            respjson.Field
		Status            respjson.Field
		EntrypointRelPath respjson.Field
		EnvVars           respjson.Field
		StatusReason      respjson.Field
		UpdatedAt         respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DeploymentNewResponse) RawJSON() string { return r.JSON.raw }
func (r *DeploymentNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current status of the deployment
type DeploymentNewResponseStatus string

const (
	DeploymentNewResponseStatusQueued     DeploymentNewResponseStatus = "queued"
	DeploymentNewResponseStatusInProgress DeploymentNewResponseStatus = "in_progress"
	DeploymentNewResponseStatusRunning    DeploymentNewResponseStatus = "running"
	DeploymentNewResponseStatusFailed     DeploymentNewResponseStatus = "failed"
	DeploymentNewResponseStatusStopped    DeploymentNewResponseStatus = "stopped"
)

// Deployment record information.
type DeploymentGetResponse struct {
	// Unique identifier for the deployment
	ID string `json:"id,required"`
	// Timestamp when the deployment was created
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// Deployment region code
	Region constant.AwsUsEast1a `json:"region,required"`
	// Current status of the deployment
	//
	// Any of "queued", "in_progress", "running", "failed", "stopped".
	Status DeploymentGetResponseStatus `json:"status,required"`
	// Relative path to the application entrypoint
	EntrypointRelPath string `json:"entrypoint_rel_path"`
	// Environment variables configured for this deployment
	EnvVars map[string]string `json:"env_vars"`
	// Status reason
	StatusReason string `json:"status_reason"`
	// Timestamp when the deployment was last updated
	UpdatedAt time.Time `json:"updated_at,nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		CreatedAt         respjson.Field
		Region            respjson.Field
		Status            respjson.Field
		EntrypointRelPath respjson.Field
		EnvVars           respjson.Field
		StatusReason      respjson.Field
		UpdatedAt         respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DeploymentGetResponse) RawJSON() string { return r.JSON.raw }
func (r *DeploymentGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current status of the deployment
type DeploymentGetResponseStatus string

const (
	DeploymentGetResponseStatusQueued     DeploymentGetResponseStatus = "queued"
	DeploymentGetResponseStatusInProgress DeploymentGetResponseStatus = "in_progress"
	DeploymentGetResponseStatusRunning    DeploymentGetResponseStatus = "running"
	DeploymentGetResponseStatusFailed     DeploymentGetResponseStatus = "failed"
	DeploymentGetResponseStatusStopped    DeploymentGetResponseStatus = "stopped"
)

// Deployment record information.
type DeploymentListResponse struct {
	// Unique identifier for the deployment
	ID string `json:"id,required"`
	// Timestamp when the deployment was created
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// Deployment region code
	Region constant.AwsUsEast1a `json:"region,required"`
	// Current status of the deployment
	//
	// Any of "queued", "in_progress", "running", "failed", "stopped".
	Status DeploymentListResponseStatus `json:"status,required"`
	// Relative path to the application entrypoint
	EntrypointRelPath string `json:"entrypoint_rel_path"`
	// Environment variables configured for this deployment
	EnvVars map[string]string `json:"env_vars"`
	// Status reason
	StatusReason string `json:"status_reason"`
	// Timestamp when the deployment was last updated
	UpdatedAt time.Time `json:"updated_at,nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		CreatedAt         respjson.Field
		Region            respjson.Field
		Status            respjson.Field
		EntrypointRelPath respjson.Field
		EnvVars           respjson.Field
		StatusReason      respjson.Field
		UpdatedAt         respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DeploymentListResponse) RawJSON() string { return r.JSON.raw }
func (r *DeploymentListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current status of the deployment
type DeploymentListResponseStatus string

const (
	DeploymentListResponseStatusQueued     DeploymentListResponseStatus = "queued"
	DeploymentListResponseStatusInProgress DeploymentListResponseStatus = "in_progress"
	DeploymentListResponseStatusRunning    DeploymentListResponseStatus = "running"
	DeploymentListResponseStatusFailed     DeploymentListResponseStatus = "failed"
	DeploymentListResponseStatusStopped    DeploymentListResponseStatus = "stopped"
)

// DeploymentFollowResponseUnion contains all possible properties and values from
// [shared.LogEvent], [DeploymentStateEvent],
// [DeploymentFollowResponseAppVersionSummaryEvent], [shared.ErrorEvent],
// [shared.HeartbeatEvent].
//
// Use the [DeploymentFollowResponseUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DeploymentFollowResponseUnion struct {
	// Any of "log", "deployment_state", nil, nil, "sse_heartbeat".
	Event string `json:"event"`
	// This field is from variant [shared.LogEvent].
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	// This field is from variant [DeploymentStateEvent].
	Deployment DeploymentStateEventDeployment `json:"deployment"`
	// This field is from variant [DeploymentFollowResponseAppVersionSummaryEvent].
	ID string `json:"id"`
	// This field is from variant [DeploymentFollowResponseAppVersionSummaryEvent].
	Actions []shared.AppAction `json:"actions"`
	// This field is from variant [DeploymentFollowResponseAppVersionSummaryEvent].
	AppName string `json:"app_name"`
	// This field is from variant [DeploymentFollowResponseAppVersionSummaryEvent].
	Region constant.AwsUsEast1a `json:"region"`
	// This field is from variant [DeploymentFollowResponseAppVersionSummaryEvent].
	Version string `json:"version"`
	// This field is from variant [DeploymentFollowResponseAppVersionSummaryEvent].
	EnvVars map[string]string `json:"env_vars"`
	// This field is from variant [shared.ErrorEvent].
	Error shared.ErrorModel `json:"error"`
	JSON  struct {
		Event      respjson.Field
		Message    respjson.Field
		Timestamp  respjson.Field
		Deployment respjson.Field
		ID         respjson.Field
		Actions    respjson.Field
		AppName    respjson.Field
		Region     respjson.Field
		Version    respjson.Field
		EnvVars    respjson.Field
		Error      respjson.Field
		raw        string
	} `json:"-"`
}

func (u DeploymentFollowResponseUnion) AsLog() (v shared.LogEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DeploymentFollowResponseUnion) AsDeploymentState() (v DeploymentStateEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DeploymentFollowResponseUnion) AsDeploymentFollowResponseAppVersionSummaryEvent() (v DeploymentFollowResponseAppVersionSummaryEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DeploymentFollowResponseUnion) AsErrorEvent() (v shared.ErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DeploymentFollowResponseUnion) AsSseHeartbeat() (v shared.HeartbeatEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DeploymentFollowResponseUnion) RawJSON() string { return u.JSON.raw }

func (r *DeploymentFollowResponseUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Summary of an application version.
type DeploymentFollowResponseAppVersionSummaryEvent struct {
	// Unique identifier for the app version
	ID string `json:"id,required"`
	// List of actions available on the app
	Actions []shared.AppAction `json:"actions,required"`
	// Name of the application
	AppName string `json:"app_name,required"`
	// Event type identifier (always "app_version_summary").
	Event constant.AppVersionSummary `json:"event,required"`
	// Deployment region code
	Region constant.AwsUsEast1a `json:"region,required"`
	// Time the state was reported.
	Timestamp time.Time `json:"timestamp,required" format:"date-time"`
	// Version label for the application
	Version string `json:"version,required"`
	// Environment variables configured for this app version
	EnvVars map[string]string `json:"env_vars"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Actions     respjson.Field
		AppName     respjson.Field
		Event       respjson.Field
		Region      respjson.Field
		Timestamp   respjson.Field
		Version     respjson.Field
		EnvVars     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DeploymentFollowResponseAppVersionSummaryEvent) RawJSON() string { return r.JSON.raw }
func (r *DeploymentFollowResponseAppVersionSummaryEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DeploymentNewParams struct {
	// Relative path to the entrypoint of the application
	EntrypointRelPath param.Opt[string] `json:"entrypoint_rel_path,omitzero"`
	// Allow overwriting an existing app version
	Force param.Opt[bool] `json:"force,omitzero"`
	// Version of the application. Can be any string.
	Version param.Opt[string] `json:"version,omitzero"`
	// Map of environment variables to set for the deployed application. Each key-value
	// pair represents an environment variable.
	EnvVars map[string]string `json:"env_vars,omitzero"`
	// ZIP file containing the application source directory
	File io.Reader `json:"file,omitzero" format:"binary"`
	// Region for deployment. Currently we only support "aws.us-east-1a"
	//
	// Any of "aws.us-east-1a".
	Region DeploymentNewParamsRegion `json:"region,omitzero"`
	// Source from which to fetch application code.
	Source DeploymentNewParamsSource `json:"source,omitzero"`
	paramObj
}

func (r DeploymentNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err == nil {
		err = apiform.WriteExtras(writer, r.ExtraFields())
	}
	if err != nil {
		writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

// Region for deployment. Currently we only support "aws.us-east-1a"
type DeploymentNewParamsRegion string

const (
	DeploymentNewParamsRegionAwsUsEast1a DeploymentNewParamsRegion = "aws.us-east-1a"
)

// Source from which to fetch application code.
//
// The properties Entrypoint, Ref, Type, URL are required.
type DeploymentNewParamsSource struct {
	// Relative path to the application entrypoint within the selected path.
	Entrypoint string `json:"entrypoint,required"`
	// Git ref (branch, tag, or commit SHA) to fetch.
	Ref string `json:"ref,required"`
	// Source type identifier.
	//
	// Any of "github".
	Type string `json:"type,omitzero,required"`
	// Base repository URL (without blob/tree suffixes).
	URL string `json:"url,required"`
	// Path within the repo to deploy (omit to use repo root).
	Path param.Opt[string] `json:"path,omitzero"`
	// Authentication for private repositories.
	Auth DeploymentNewParamsSourceAuth `json:"auth,omitzero"`
	paramObj
}

func (r DeploymentNewParamsSource) MarshalJSON() (data []byte, err error) {
	type shadow DeploymentNewParamsSource
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DeploymentNewParamsSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DeploymentNewParamsSource](
		"type", "github",
	)
}

// Authentication for private repositories.
//
// The properties Token, Method are required.
type DeploymentNewParamsSourceAuth struct {
	// GitHub PAT or installation access token
	Token string `json:"token,required" format:"password"`
	// Auth method
	//
	// Any of "github_token".
	Method string `json:"method,omitzero,required"`
	paramObj
}

func (r DeploymentNewParamsSourceAuth) MarshalJSON() (data []byte, err error) {
	type shadow DeploymentNewParamsSourceAuth
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DeploymentNewParamsSourceAuth) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DeploymentNewParamsSourceAuth](
		"method", "github_token",
	)
}

type DeploymentListParams struct {
	// Filter results by application name.
	AppName param.Opt[string] `query:"app_name,omitzero" json:"-"`
	// Limit the number of deployments to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Offset the number of deployments to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [DeploymentListParams]'s query parameters as `url.Values`.
func (r DeploymentListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DeploymentFollowParams struct {
	// Show logs since the given time (RFC timestamps or durations like 5m).
	Since param.Opt[string] `query:"since,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [DeploymentFollowParams]'s query parameters as `url.Values`.
func (r DeploymentFollowParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
