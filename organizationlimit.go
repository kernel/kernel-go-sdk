// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"net/http"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Read and manage organization-level limits.
//
// OrganizationLimitService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationLimitService] method instead.
type OrganizationLimitService struct {
	Options []option.RequestOption
}

// NewOrganizationLimitService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewOrganizationLimitService(opts ...option.RequestOption) (r OrganizationLimitService) {
	r = OrganizationLimitService{}
	r.Options = opts
	return
}

// Get the organization's effective limits and managed auth and vault usage.
func (r *OrganizationLimitService) Get(ctx context.Context, opts ...option.RequestOption) (res *OrgLimits, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "org/limits"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Set the default per-project concurrency cap applied to projects without an
// explicit override. Set the value to 0 to remove the default; omit to leave it
// unchanged. The default cannot exceed the organization's concurrency limit.
func (r *OrganizationLimitService) Update(ctx context.Context, body OrganizationLimitUpdateParams, opts ...option.RequestOption) (res *OrgLimits, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "org/limits"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

type OrgLimits struct {
	// The organization's current non-deleted managed auth connections, counted
	// org-wide across every project. Compare against max_auth_connections to show
	// remaining capacity before a create is rejected with 403 insufficient_plan.
	AuthConnectionsUsed int64 `json:"auth_connections_used" api:"required"`
	// Maximum managed auth connections the organization's plan allows. Null means
	// unlimited. Counted org-wide, so it cannot be multiplied across projects.
	MaxAuthConnections int64 `json:"max_auth_connections" api:"required"`
	// Maximum non-deleted vaults allowed org-wide across all projects. Null means
	// unlimited.
	MaxVaults int64 `json:"max_vaults" api:"required"`
	// Smallest health_check_interval the organization's plan accepts on a managed auth
	// connection. Requests below this are rejected with 400. Existing connections
	// stored below the floor are grandfathered until edited.
	MinHealthCheckIntervalSeconds int64 `json:"min_health_check_interval_seconds" api:"required"`
	// Current non-deleted vault count across all projects in the organization.
	VaultsUsed int64 `json:"vaults_used" api:"required"`
	// Default maximum concurrent browsers applied to every project that has no
	// explicit per-project override. Null means no org-level default, so such projects
	// are uncapped (only the org-wide limit applies). Applies to existing and newly
	// created projects.
	DefaultProjectMaxConcurrentSessions int64 `json:"default_project_max_concurrent_sessions" api:"nullable"`
	// The organization's effective concurrency limit — the maximum browsers running at
	// once, covering both on-demand sessions and browser pool reservations — from its
	// plan or an override. Read-only and shared across all projects in the org; a
	// per-project default cannot exceed it.
	MaxConcurrentSessions int64 `json:"max_concurrent_sessions"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AuthConnectionsUsed                 respjson.Field
		MaxAuthConnections                  respjson.Field
		MaxVaults                           respjson.Field
		MinHealthCheckIntervalSeconds       respjson.Field
		VaultsUsed                          respjson.Field
		DefaultProjectMaxConcurrentSessions respjson.Field
		MaxConcurrentSessions               respjson.Field
		ExtraFields                         map[string]respjson.Field
		raw                                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgLimits) RawJSON() string { return r.JSON.raw }
func (r *OrgLimits) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type UpdateOrgLimitsRequestParam struct {
	// Default maximum concurrent browsers for projects without an explicit override.
	// Set to 0 to remove the default; omit to leave unchanged. Cannot exceed the
	// organization's concurrency limit.
	DefaultProjectMaxConcurrentSessions param.Opt[int64] `json:"default_project_max_concurrent_sessions,omitzero"`
	paramObj
}

func (r UpdateOrgLimitsRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow UpdateOrgLimitsRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *UpdateOrgLimitsRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationLimitUpdateParams struct {
	UpdateOrgLimitsRequest UpdateOrgLimitsRequestParam
	paramObj
}

func (r OrganizationLimitUpdateParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.UpdateOrgLimitsRequest)
}
func (r *OrganizationLimitUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
