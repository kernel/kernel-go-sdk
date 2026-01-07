// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"encoding/json"
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
	"github.com/kernel/kernel-go-sdk/packages/ssestream"
	"github.com/kernel/kernel-go-sdk/shared"
	"github.com/kernel/kernel-go-sdk/shared/constant"
)

// InvocationService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewInvocationService] method instead.
type InvocationService struct {
	Options []option.RequestOption
}

// NewInvocationService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewInvocationService(opts ...option.RequestOption) (r InvocationService) {
	r = InvocationService{}
	r.Options = opts
	return
}

// Invoke an action.
func (r *InvocationService) New(ctx context.Context, body InvocationNewParams, opts ...option.RequestOption) (res *InvocationNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "invocations"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get details about an invocation's status and output.
func (r *InvocationService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *InvocationGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("invocations/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Update an invocation's status or output. This can be used to cancel an
// invocation by setting the status to "failed".
func (r *InvocationService) Update(ctx context.Context, id string, body InvocationUpdateParams, opts ...option.RequestOption) (res *InvocationUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("invocations/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// List invocations. Optionally filter by application name, action name, status,
// deployment ID, or start time.
func (r *InvocationService) List(ctx context.Context, query InvocationListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[InvocationListResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "invocations"
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

// List invocations. Optionally filter by application name, action name, status,
// deployment ID, or start time.
func (r *InvocationService) ListAutoPaging(ctx context.Context, query InvocationListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[InvocationListResponse] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete all browser sessions created within the specified invocation.
func (r *InvocationService) DeleteBrowsers(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("invocations/%s/browsers", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return
}

// Establishes a Server-Sent Events (SSE) stream that delivers real-time logs and
// status updates for an invocation. The stream terminates automatically once the
// invocation reaches a terminal state.
func (r *InvocationService) FollowStreaming(ctx context.Context, id string, query InvocationFollowParams, opts ...option.RequestOption) (stream *ssestream.Stream[InvocationFollowResponseUnion]) {
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
	path := fmt.Sprintf("invocations/%s/events", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &raw, opts...)
	return ssestream.NewStream[InvocationFollowResponseUnion](ssestream.NewDecoder(raw), err)
}

// An event representing the current state of an invocation.
type InvocationStateEvent struct {
	// Event type identifier (always "invocation_state").
	Event      constant.InvocationState       `json:"event,required"`
	Invocation InvocationStateEventInvocation `json:"invocation,required"`
	// Time the state was reported.
	Timestamp time.Time `json:"timestamp,required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Event       respjson.Field
		Invocation  respjson.Field
		Timestamp   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InvocationStateEvent) RawJSON() string { return r.JSON.raw }
func (r *InvocationStateEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type InvocationStateEventInvocation struct {
	// ID of the invocation
	ID string `json:"id,required"`
	// Name of the action invoked
	ActionName string `json:"action_name,required"`
	// Name of the application
	AppName string `json:"app_name,required"`
	// RFC 3339 Nanoseconds timestamp when the invocation started
	StartedAt time.Time `json:"started_at,required" format:"date-time"`
	// Status of the invocation
	//
	// Any of "queued", "running", "succeeded", "failed".
	Status string `json:"status,required"`
	// Version label for the application
	Version string `json:"version,required"`
	// RFC 3339 Nanoseconds timestamp when the invocation finished (null if still
	// running)
	FinishedAt time.Time `json:"finished_at,nullable" format:"date-time"`
	// Output produced by the action, rendered as a JSON string. This could be: string,
	// number, boolean, array, object, or null.
	Output string `json:"output"`
	// Payload provided to the invocation. This is a string that can be parsed as JSON.
	Payload string `json:"payload"`
	// Status reason
	StatusReason string `json:"status_reason"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		ActionName   respjson.Field
		AppName      respjson.Field
		StartedAt    respjson.Field
		Status       respjson.Field
		Version      respjson.Field
		FinishedAt   respjson.Field
		Output       respjson.Field
		Payload      respjson.Field
		StatusReason respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InvocationStateEventInvocation) RawJSON() string { return r.JSON.raw }
func (r *InvocationStateEventInvocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type InvocationNewResponse struct {
	// ID of the invocation
	ID string `json:"id,required"`
	// Name of the action invoked
	ActionName string `json:"action_name,required"`
	// Status of the invocation
	//
	// Any of "queued", "running", "succeeded", "failed".
	Status InvocationNewResponseStatus `json:"status,required"`
	// The return value of the action that was invoked, rendered as a JSON string. This
	// could be: string, number, boolean, array, object, or null.
	Output string `json:"output"`
	// Status reason
	StatusReason string `json:"status_reason"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		ActionName   respjson.Field
		Status       respjson.Field
		Output       respjson.Field
		StatusReason respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InvocationNewResponse) RawJSON() string { return r.JSON.raw }
func (r *InvocationNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of the invocation
type InvocationNewResponseStatus string

const (
	InvocationNewResponseStatusQueued    InvocationNewResponseStatus = "queued"
	InvocationNewResponseStatusRunning   InvocationNewResponseStatus = "running"
	InvocationNewResponseStatusSucceeded InvocationNewResponseStatus = "succeeded"
	InvocationNewResponseStatusFailed    InvocationNewResponseStatus = "failed"
)

type InvocationGetResponse struct {
	// ID of the invocation
	ID string `json:"id,required"`
	// Name of the action invoked
	ActionName string `json:"action_name,required"`
	// Name of the application
	AppName string `json:"app_name,required"`
	// RFC 3339 Nanoseconds timestamp when the invocation started
	StartedAt time.Time `json:"started_at,required" format:"date-time"`
	// Status of the invocation
	//
	// Any of "queued", "running", "succeeded", "failed".
	Status InvocationGetResponseStatus `json:"status,required"`
	// Version label for the application
	Version string `json:"version,required"`
	// RFC 3339 Nanoseconds timestamp when the invocation finished (null if still
	// running)
	FinishedAt time.Time `json:"finished_at,nullable" format:"date-time"`
	// Output produced by the action, rendered as a JSON string. This could be: string,
	// number, boolean, array, object, or null.
	Output string `json:"output"`
	// Payload provided to the invocation. This is a string that can be parsed as JSON.
	Payload string `json:"payload"`
	// Status reason
	StatusReason string `json:"status_reason"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		ActionName   respjson.Field
		AppName      respjson.Field
		StartedAt    respjson.Field
		Status       respjson.Field
		Version      respjson.Field
		FinishedAt   respjson.Field
		Output       respjson.Field
		Payload      respjson.Field
		StatusReason respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InvocationGetResponse) RawJSON() string { return r.JSON.raw }
func (r *InvocationGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of the invocation
type InvocationGetResponseStatus string

const (
	InvocationGetResponseStatusQueued    InvocationGetResponseStatus = "queued"
	InvocationGetResponseStatusRunning   InvocationGetResponseStatus = "running"
	InvocationGetResponseStatusSucceeded InvocationGetResponseStatus = "succeeded"
	InvocationGetResponseStatusFailed    InvocationGetResponseStatus = "failed"
)

type InvocationUpdateResponse struct {
	// ID of the invocation
	ID string `json:"id,required"`
	// Name of the action invoked
	ActionName string `json:"action_name,required"`
	// Name of the application
	AppName string `json:"app_name,required"`
	// RFC 3339 Nanoseconds timestamp when the invocation started
	StartedAt time.Time `json:"started_at,required" format:"date-time"`
	// Status of the invocation
	//
	// Any of "queued", "running", "succeeded", "failed".
	Status InvocationUpdateResponseStatus `json:"status,required"`
	// Version label for the application
	Version string `json:"version,required"`
	// RFC 3339 Nanoseconds timestamp when the invocation finished (null if still
	// running)
	FinishedAt time.Time `json:"finished_at,nullable" format:"date-time"`
	// Output produced by the action, rendered as a JSON string. This could be: string,
	// number, boolean, array, object, or null.
	Output string `json:"output"`
	// Payload provided to the invocation. This is a string that can be parsed as JSON.
	Payload string `json:"payload"`
	// Status reason
	StatusReason string `json:"status_reason"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		ActionName   respjson.Field
		AppName      respjson.Field
		StartedAt    respjson.Field
		Status       respjson.Field
		Version      respjson.Field
		FinishedAt   respjson.Field
		Output       respjson.Field
		Payload      respjson.Field
		StatusReason respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InvocationUpdateResponse) RawJSON() string { return r.JSON.raw }
func (r *InvocationUpdateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of the invocation
type InvocationUpdateResponseStatus string

const (
	InvocationUpdateResponseStatusQueued    InvocationUpdateResponseStatus = "queued"
	InvocationUpdateResponseStatusRunning   InvocationUpdateResponseStatus = "running"
	InvocationUpdateResponseStatusSucceeded InvocationUpdateResponseStatus = "succeeded"
	InvocationUpdateResponseStatusFailed    InvocationUpdateResponseStatus = "failed"
)

type InvocationListResponse struct {
	// ID of the invocation
	ID string `json:"id,required"`
	// Name of the action invoked
	ActionName string `json:"action_name,required"`
	// Name of the application
	AppName string `json:"app_name,required"`
	// RFC 3339 Nanoseconds timestamp when the invocation started
	StartedAt time.Time `json:"started_at,required" format:"date-time"`
	// Status of the invocation
	//
	// Any of "queued", "running", "succeeded", "failed".
	Status InvocationListResponseStatus `json:"status,required"`
	// Version label for the application
	Version string `json:"version,required"`
	// RFC 3339 Nanoseconds timestamp when the invocation finished (null if still
	// running)
	FinishedAt time.Time `json:"finished_at,nullable" format:"date-time"`
	// Output produced by the action, rendered as a JSON string. This could be: string,
	// number, boolean, array, object, or null.
	Output string `json:"output"`
	// Payload provided to the invocation. This is a string that can be parsed as JSON.
	Payload string `json:"payload"`
	// Status reason
	StatusReason string `json:"status_reason"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		ActionName   respjson.Field
		AppName      respjson.Field
		StartedAt    respjson.Field
		Status       respjson.Field
		Version      respjson.Field
		FinishedAt   respjson.Field
		Output       respjson.Field
		Payload      respjson.Field
		StatusReason respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InvocationListResponse) RawJSON() string { return r.JSON.raw }
func (r *InvocationListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of the invocation
type InvocationListResponseStatus string

const (
	InvocationListResponseStatusQueued    InvocationListResponseStatus = "queued"
	InvocationListResponseStatusRunning   InvocationListResponseStatus = "running"
	InvocationListResponseStatusSucceeded InvocationListResponseStatus = "succeeded"
	InvocationListResponseStatusFailed    InvocationListResponseStatus = "failed"
)

// InvocationFollowResponseUnion contains all possible properties and values from
// [shared.LogEvent], [InvocationStateEvent], [shared.ErrorEvent],
// [shared.HeartbeatEvent].
//
// Use the [InvocationFollowResponseUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type InvocationFollowResponseUnion struct {
	// Any of "log", "invocation_state", "error", "sse_heartbeat".
	Event string `json:"event"`
	// This field is from variant [shared.LogEvent].
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	// This field is from variant [InvocationStateEvent].
	Invocation InvocationStateEventInvocation `json:"invocation"`
	// This field is from variant [shared.ErrorEvent].
	Error shared.ErrorModel `json:"error"`
	JSON  struct {
		Event      respjson.Field
		Message    respjson.Field
		Timestamp  respjson.Field
		Invocation respjson.Field
		Error      respjson.Field
		raw        string
	} `json:"-"`
}

// anyInvocationFollowResponse is implemented by each variant of
// [InvocationFollowResponseUnion] to add type safety for the return type of
// [InvocationFollowResponseUnion.AsAny]
type anyInvocationFollowResponse interface {
	ImplInvocationFollowResponseUnion()
}

func (InvocationStateEvent) ImplInvocationFollowResponseUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := InvocationFollowResponseUnion.AsAny().(type) {
//	case shared.LogEvent:
//	case kernel.InvocationStateEvent:
//	case shared.ErrorEvent:
//	case shared.HeartbeatEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u InvocationFollowResponseUnion) AsAny() anyInvocationFollowResponse {
	switch u.Event {
	case "log":
		return u.AsLog()
	case "invocation_state":
		return u.AsInvocationState()
	case "error":
		return u.AsError()
	case "sse_heartbeat":
		return u.AsSseHeartbeat()
	}
	return nil
}

func (u InvocationFollowResponseUnion) AsLog() (v shared.LogEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u InvocationFollowResponseUnion) AsInvocationState() (v InvocationStateEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u InvocationFollowResponseUnion) AsError() (v shared.ErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u InvocationFollowResponseUnion) AsSseHeartbeat() (v shared.HeartbeatEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u InvocationFollowResponseUnion) RawJSON() string { return u.JSON.raw }

func (r *InvocationFollowResponseUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type InvocationNewParams struct {
	// Name of the action to invoke
	ActionName string `json:"action_name,required"`
	// Name of the application
	AppName string `json:"app_name,required"`
	// Version of the application
	Version string `json:"version,required"`
	// If true, invoke asynchronously. When set, the API responds 202 Accepted with
	// status "queued".
	Async param.Opt[bool] `json:"async,omitzero"`
	// Timeout in seconds for async invocations (min 10, max 3600). Only applies when
	// async is true.
	AsyncTimeoutSeconds param.Opt[int64] `json:"async_timeout_seconds,omitzero"`
	// Input data for the action, sent as a JSON string.
	Payload param.Opt[string] `json:"payload,omitzero"`
	paramObj
}

func (r InvocationNewParams) MarshalJSON() (data []byte, err error) {
	type shadow InvocationNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InvocationNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type InvocationUpdateParams struct {
	// New status for the invocation.
	//
	// Any of "succeeded", "failed".
	Status InvocationUpdateParamsStatus `json:"status,omitzero,required"`
	// Updated output of the invocation rendered as JSON string.
	Output param.Opt[string] `json:"output,omitzero"`
	paramObj
}

func (r InvocationUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow InvocationUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InvocationUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// New status for the invocation.
type InvocationUpdateParamsStatus string

const (
	InvocationUpdateParamsStatusSucceeded InvocationUpdateParamsStatus = "succeeded"
	InvocationUpdateParamsStatusFailed    InvocationUpdateParamsStatus = "failed"
)

type InvocationListParams struct {
	// Filter results by action name.
	ActionName param.Opt[string] `query:"action_name,omitzero" json:"-"`
	// Filter results by application name.
	AppName param.Opt[string] `query:"app_name,omitzero" json:"-"`
	// Filter results by deployment ID.
	DeploymentID param.Opt[string] `query:"deployment_id,omitzero" json:"-"`
	// Limit the number of invocations to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Offset the number of invocations to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Show invocations that have started since the given time (RFC timestamps or
	// durations like 5m).
	Since param.Opt[string] `query:"since,omitzero" json:"-"`
	// Filter results by application version.
	Version param.Opt[string] `query:"version,omitzero" json:"-"`
	// Filter results by invocation status.
	//
	// Any of "queued", "running", "succeeded", "failed".
	Status InvocationListParamsStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [InvocationListParams]'s query parameters as `url.Values`.
func (r InvocationListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter results by invocation status.
type InvocationListParamsStatus string

const (
	InvocationListParamsStatusQueued    InvocationListParamsStatus = "queued"
	InvocationListParamsStatusRunning   InvocationListParamsStatus = "running"
	InvocationListParamsStatusSucceeded InvocationListParamsStatus = "succeeded"
	InvocationListParamsStatusFailed    InvocationListParamsStatus = "failed"
)

type InvocationFollowParams struct {
	// Show logs since the given time (RFC timestamps or durations like 5m).
	Since param.Opt[string] `query:"since,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [InvocationFollowParams]'s query parameters as `url.Values`.
func (r InvocationFollowParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
