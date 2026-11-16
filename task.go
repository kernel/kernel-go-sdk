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

// Define tasks that agents work toward and link the browsers they use.
//
// TaskService contains methods and other services that help with interacting with
// the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewTaskService] method instead.
type TaskService struct {
	Options []option.RequestOption
}

// NewTaskService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewTaskService(opts ...option.RequestOption) (r TaskService) {
	r = TaskService{}
	r.Options = opts
	return
}

// Define a task before or while an agent works on it, optionally linking the
// browsers it uses. Submit reports against the returned task ID.
func (r *TaskService) New(ctx context.Context, body TaskNewParams, opts ...option.RequestOption) (res *Task, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "tasks"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a task and the IDs of its linked browsers.
func (r *TaskService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *Task, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("tasks/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List tasks in the request's project, newest first.
func (r *TaskService) List(ctx context.Context, query TaskListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[Task], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "tasks"
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

// List tasks in the request's project, newest first.
func (r *TaskService) ListAutoPaging(ctx context.Context, query TaskListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[Task] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Link browsers to an existing task. Linking a browser that is already linked has
// no effect, so the call is safe to retry. Returns the updated task.
func (r *TaskService) LinkBrowsers(ctx context.Context, id string, body TaskLinkBrowsersParams, opts ...option.RequestOption) (res *Task, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("tasks/%s/browsers", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// A goal an agent works toward. Reports are submitted against a task, and browsers
// can be linked to it. Tasks are immutable apart from added browser links.
type Task struct {
	// Unique identifier for the task
	ID string `json:"id" api:"required"`
	// When the task was created
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// What the task is meant to accomplish
	Description string `json:"description" api:"required"`
	// Caller-defined string key-value pairs for grouping and filtering tasks, such as
	// a user or thread ID. At most 20 keys; keys are 1-64 characters, values at most
	// 256 characters, and keys and values together at most 4 KB.
	Metadata TaskMetadata `json:"metadata" api:"required"`
	// IDs of browsers linked to the task, oldest link first. Linked browsers may since
	// have been deleted.
	SessionIDs []string `json:"session_ids" api:"required"`
	// When work on the task started
	StartedAt time.Time `json:"started_at" api:"required" format:"date-time"`
	// Hostname of the website the task targets
	TargetWebsite string `json:"target_website" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		Description   respjson.Field
		Metadata      respjson.Field
		SessionIDs    respjson.Field
		StartedAt     respjson.Field
		TargetWebsite respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Task) RawJSON() string { return r.JSON.raw }
func (r *Task) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TaskMetadata map[string]string

type TaskNewParams struct {
	// What the task is meant to accomplish, at most 4 KB.
	Description string `json:"description" api:"required"`
	// When work on the task started. Defaults to the time the task is created. Must be
	// no more than 5 minutes in the future or 30 days in the past.
	StartedAt param.Opt[time.Time] `json:"started_at,omitzero" format:"date-time"`
	// Hostname of the website the task targets, without a scheme or path.
	TargetWebsite param.Opt[string] `json:"target_website,omitzero"`
	// Caller-defined string key-value pairs for grouping and filtering tasks, such as
	// a user or thread ID. At most 20 keys; keys are 1-64 characters, values at most
	// 256 characters, and keys and values together at most 4 KB.
	Metadata TaskMetadata `json:"metadata,omitzero"`
	// IDs of browsers to link to the task. Each must belong to the request's project;
	// deleted browsers can be linked.
	SessionIDs []string `json:"session_ids,omitzero"`
	paramObj
}

func (r TaskNewParams) MarshalJSON() (data []byte, err error) {
	type shadow TaskNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TaskNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TaskListParams struct {
	// Limit the number of tasks to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Offset the number of tasks to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Only return tasks linked to this browser ID.
	SessionID param.Opt[string] `query:"session_id,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [TaskListParams]'s query parameters as `url.Values`.
func (r TaskListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type TaskLinkBrowsersParams struct {
	// IDs of browsers to link. Each must belong to the task's project; deleted
	// browsers can be linked. Browsers already linked are ignored. A task can have at
	// most 100 linked browsers.
	SessionIDs []string `json:"session_ids,omitzero" api:"required"`
	paramObj
}

func (r TaskLinkBrowsersParams) MarshalJSON() (data []byte, err error) {
	type shadow TaskLinkBrowsersParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TaskLinkBrowsersParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
