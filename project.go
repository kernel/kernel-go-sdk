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

// Create and manage projects for resource isolation within an organization.
//
// ProjectService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewProjectService] method instead.
type ProjectService struct {
	Options []option.RequestOption
	// Create and manage projects for resource isolation within an organization.
	Limits ProjectLimitService
}

// NewProjectService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewProjectService(opts ...option.RequestOption) (r ProjectService) {
	r = ProjectService{}
	r.Options = opts
	r.Limits = NewProjectLimitService(opts...)
	return
}

// Create a new project within the authenticated organization. Requires the
// projects feature flag.
func (r *ProjectService) New(ctx context.Context, body ProjectNewParams, opts ...option.RequestOption) (res *Project, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "projects"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get a project by ID.
func (r *ProjectService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *Project, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("projects/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a project's name or status.
func (r *ProjectService) Update(ctx context.Context, id string, body ProjectUpdateParams, opts ...option.RequestOption) (res *Project, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("projects/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List projects for the authenticated organization.
func (r *ProjectService) List(ctx context.Context, query ProjectListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[Project], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "projects"
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

// List projects for the authenticated organization.
func (r *ProjectService) ListAutoPaging(ctx context.Context, query ProjectListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[Project] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Soft-delete a project. The project must be empty (no active resources).
func (r *ProjectService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("projects/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// The property Name is required.
type CreateProjectRequestParam struct {
	// Project name (1-255 characters)
	Name string `json:"name" api:"required"`
	paramObj
}

func (r CreateProjectRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateProjectRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CreateProjectRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Project struct {
	// Unique project identifier
	ID string `json:"id" api:"required"`
	// When the project was created
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Project name
	Name string `json:"name" api:"required"`
	// Project status
	//
	// Any of "active", "archived".
	Status ProjectStatus `json:"status" api:"required"`
	// When the project was last updated
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Status      respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Project) RawJSON() string { return r.JSON.raw }
func (r *Project) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Project status
type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusArchived ProjectStatus = "archived"
)

type UpdateProjectRequestParam struct {
	// New project name
	Name param.Opt[string] `json:"name,omitzero"`
	// New project status
	//
	// Any of "active", "archived".
	Status UpdateProjectRequestStatus `json:"status,omitzero"`
	paramObj
}

func (r UpdateProjectRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow UpdateProjectRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *UpdateProjectRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// New project status
type UpdateProjectRequestStatus string

const (
	UpdateProjectRequestStatusActive   UpdateProjectRequestStatus = "active"
	UpdateProjectRequestStatusArchived UpdateProjectRequestStatus = "archived"
)

type ProjectNewParams struct {
	CreateProjectRequest CreateProjectRequestParam
	paramObj
}

func (r ProjectNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.CreateProjectRequest)
}
func (r *ProjectNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProjectUpdateParams struct {
	UpdateProjectRequest UpdateProjectRequestParam
	paramObj
}

func (r ProjectUpdateParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.UpdateProjectRequest)
}
func (r *ProjectUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProjectListParams struct {
	// Maximum number of results to return
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Number of results to skip
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ProjectListParams]'s query parameters as `url.Values`.
func (r ProjectListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
