// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Create and manage projects for resource isolation within an organization.
//
// ProjectLimitService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewProjectLimitService] method instead.
type ProjectLimitService struct {
	Options []option.RequestOption
}

// NewProjectLimitService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewProjectLimitService(opts ...option.RequestOption) (r ProjectLimitService) {
	r = ProjectLimitService{}
	r.Options = opts
	return
}

// Get the resource limit overrides for a project. Null values mean no
// project-level cap (org limit applies).
func (r *ProjectLimitService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *ProjectLimits, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/projects/%s/limits", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update resource limit overrides for a project. Only fields present in the
// request are modified. Set a field to 0 to remove that limit cap; omit a field to
// leave it unchanged.
func (r *ProjectLimitService) Update(ctx context.Context, id string, body ProjectLimitUpdateParams, opts ...option.RequestOption) (res *ProjectLimits, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/projects/%s/limits", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

type ProjectLimits struct {
	// Maximum concurrent app invocations for this project. Null means no project-level
	// cap.
	MaxConcurrentInvocations int64 `json:"max_concurrent_invocations" api:"nullable"`
	// Maximum concurrent browsers for this project, covering both on-demand sessions
	// (`browsers.create()`) and browser pool reservations. Null means no project-level
	// cap.
	MaxConcurrentSessions int64 `json:"max_concurrent_sessions" api:"nullable"`
	// Deprecated: pooled browsers now count toward `max_concurrent_sessions`. Always
	// null once the unified concurrency limit is enabled for your organization.
	//
	// Deprecated: deprecated
	MaxPooledSessions int64 `json:"max_pooled_sessions" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MaxConcurrentInvocations respjson.Field
		MaxConcurrentSessions    respjson.Field
		MaxPooledSessions        respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectLimits) RawJSON() string { return r.JSON.raw }
func (r *ProjectLimits) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type UpdateProjectLimitsRequestParam struct {
	// Maximum concurrent app invocations for this project. Set to 0 to remove the cap;
	// omit to leave unchanged.
	MaxConcurrentInvocations param.Opt[int64] `json:"max_concurrent_invocations,omitzero"`
	// Maximum concurrent browsers for this project, covering both on-demand sessions
	// and browser pool reservations. Set to 0 to remove the cap; omit to leave
	// unchanged.
	MaxConcurrentSessions param.Opt[int64] `json:"max_concurrent_sessions,omitzero"`
	// Deprecated: pooled browsers now count toward `max_concurrent_sessions`. Requests
	// that set this field are rejected with a 400 once the unified concurrency limit
	// is enabled for your organization.
	//
	// Deprecated: deprecated
	MaxPooledSessions param.Opt[int64] `json:"max_pooled_sessions,omitzero"`
	paramObj
}

func (r UpdateProjectLimitsRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow UpdateProjectLimitsRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *UpdateProjectLimitsRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProjectLimitUpdateParams struct {
	UpdateProjectLimitsRequest UpdateProjectLimitsRequestParam
	paramObj
}

func (r ProjectLimitUpdateParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.UpdateProjectLimitsRequest)
}
func (r *ProjectLimitUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
