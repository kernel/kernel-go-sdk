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

// Create and manage API keys for organization and project-scoped access.
//
// APIKeyService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAPIKeyService] method instead.
type APIKeyService struct {
	Options []option.RequestOption
}

// NewAPIKeyService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewAPIKeyService(opts ...option.RequestOption) (r APIKeyService) {
	r = APIKeyService{}
	r.Options = opts
	return
}

// Create a personal API key, or link a key to an organization-owned service
// account. Service-account keys require an organization administrator, explicit
// project and nonempty policy. An optional `policy` scopes the key to its listed
// actions and resources, and each statement may carry a `condition`. A `source_ip`
// condition restricts the statement to requests whose client address falls inside
// its CIDR networks, so a key limited to your egress IP ranges only works from
// your network, even if its secret leaks.
func (r *APIKeyService) New(ctx context.Context, body APIKeyNewParams, opts ...option.RequestOption) (res *CreatedAPIKey, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "org/api_keys"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve an API key by ID for the authenticated organization. API keys are
// masked.
func (r *APIKeyService) Get(ctx context.Context, id string, query APIKeyGetParams, opts ...option.RequestOption) (res *APIKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/api_keys/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Update an API key's name, or replace its permissions by attaching a new policy
// in place: the key's credential never changes, so consumers keep the same secret.
// The replacement policy may carry `source_ip` conditions, so you can add, change,
// or remove IP allowlisting on an existing key without rotating its credential. A
// present `policy` fully replaces the key's policy and takes effect for new
// requests within the authorization cache; requests from outside a condition's
// networks get `403` and are recorded in the audit log with the client address.
// Requires api_keys:write to rename; api_keys:set_policy to replace permissions.
func (r *APIKeyService) Update(ctx context.Context, id string, body APIKeyUpdateParams, opts ...option.RequestOption) (res *APIKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/api_keys/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List API keys for the authenticated organization. API keys are masked.
func (r *APIKeyService) List(ctx context.Context, query APIKeyListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[APIKey], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "org/api_keys"
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

// List API keys for the authenticated organization. API keys are masked.
func (r *APIKeyService) ListAutoPaging(ctx context.Context, query APIKeyListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[APIKey] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete an API key and asynchronously end browser sessions attributed to it,
// discarding unsaved profile changes. A key cannot delete itself; use a different
// key to delete this one. If cleanup cannot be registered, the key is not deleted.
func (r *APIKeyService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("org/api_keys/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Rotate an API key. Issues a new key that copies the name and project of the
// rotated key, and schedules the rotated key to expire after a grace period so
// in-flight callers can swap over. The new plaintext key is returned once. With
// expire_in_days set to 0, browser sessions attributed to the old key are ended
// asynchronously and unsaved profile changes are discarded. Grace-period and
// natural expiry do not end existing browsers. If cleanup cannot be registered, an
// immediate rotation is not committed. Service-account keys retain their service
// principal and policy; the rotating administrator is recorded separately for
// audit.
func (r *APIKeyService) Rotate(ctx context.Context, id string, body APIKeyRotateParams, opts ...option.RequestOption) (res *CreatedAPIKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/api_keys/%s/rotate", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type APIKey struct {
	// Unique API key identifier
	ID string `json:"id" api:"required"`
	// When the API key was created
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Human audit creator, or null after creator deletion. Not the service
	// credential's principal.
	CreatedBy APIKeyCreatedBy `json:"created_by" api:"required"`
	// When the API key was deleted (soft-deleted). Null for keys that have not been
	// deleted.
	DeletedAt time.Time `json:"deleted_at" api:"required" format:"date-time"`
	// When the API key expires
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// When the API key last authenticated a request, accurate to within about 20
	// minutes. Null if the key has not been used since last-use tracking began.
	LastUsedAt time.Time `json:"last_used_at" api:"required" format:"date-time"`
	// Masked version of the API key
	MaskedKey string `json:"masked_key" api:"required"`
	// Label for the API key. API keys are not addressable by name; use the ID or key
	// identifier for stable references.
	Name string `json:"name" api:"required"`
	// Project identifier for project-scoped API keys. Null means org-wide.
	ProjectID string `json:"project_id" api:"required"`
	// Project name for project-scoped API keys. Null means the key is org-wide or the
	// project name is unavailable.
	ProjectName string `json:"project_name" api:"required"`
	// Permissions granted to a scoped API key. Absent for an ordinary key, which keeps
	// unrestricted access within its organization or project.
	Policy    APIKeyPolicy    `json:"policy" api:"nullable"`
	Principal APIKeyPrincipal `json:"principal"`
	// Organization-owned service principal, when linked.
	ServiceAccountID string `json:"service_account_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		CreatedAt        respjson.Field
		CreatedBy        respjson.Field
		DeletedAt        respjson.Field
		ExpiresAt        respjson.Field
		LastUsedAt       respjson.Field
		MaskedKey        respjson.Field
		Name             respjson.Field
		ProjectID        respjson.Field
		ProjectName      respjson.Field
		Policy           respjson.Field
		Principal        respjson.Field
		ServiceAccountID respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKey) RawJSON() string { return r.JSON.raw }
func (r *APIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Human audit creator, or null after creator deletion. Not the service
// credential's principal.
type APIKeyCreatedBy struct {
	// Kernel user ID of the creator.
	ID string `json:"id" api:"required"`
	// Email address of the creator.
	Email string `json:"email" api:"required" format:"email"`
	// Display name of the creator, if available.
	Name string `json:"name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Email       respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKeyCreatedBy) RawJSON() string { return r.JSON.raw }
func (r *APIKeyCreatedBy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Permissions granted to a scoped API key. Absent for an ordinary key, which keeps
// unrestricted access within its organization or project.
type APIKeyPolicy struct {
	// Allow statements. Every action and resource not listed is denied.
	Statements []APIKeyPolicyStatement `json:"statements" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Statements  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKeyPolicy) RawJSON() string { return r.JSON.raw }
func (r *APIKeyPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type APIKeyPolicyStatement struct {
	// Permitted action, such as browsers:use or vaults:read. When creating an API key
	// or replacing its policy, `*` may replace whole components as shorthand:
	// browsers:_ names every browsers action, _:delete every delete action the
	// resource can carry, and _:_ every action the resource can carry. The stored
	// policy lists the concrete actions, so actions added later are not granted, and
	// api_keys:issue_beyond_expiry must always be listed explicitly. Partial
	// components such as brow*:use and a bare * are rejected. Other policies, such as
	// workload identity rules and deployment runtime policies, accept only exact
	// actions.
	Action string `json:"action" api:"required"`
	// Resource the action applies to, in the format projects/{project_id}/{kind}/{id}.
	// Wildcards are supported.
	Resource string `json:"resource" api:"required"`
	// A condition narrows a policy statement to requests matching its constraint: the
	// request attribute, a comparison operator, and that operator's values.
	Condition APIKeyPolicyStatementCondition `json:"condition"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Action      respjson.Field
		Resource    respjson.Field
		Condition   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKeyPolicyStatement) RawJSON() string { return r.JSON.raw }
func (r *APIKeyPolicyStatement) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A condition narrows a policy statement to requests matching its constraint: the
// request attribute, a comparison operator, and that operator's values.
type APIKeyPolicyStatementCondition struct {
	// Request attribute the condition constrains. `source_ip` restricts the statement
	// to requests whose client address falls inside one of `values`.
	//
	// Any of "source_ip".
	Key string `json:"key" api:"required"`
	// Comparison the condition applies between the request attribute and values.
	//
	// Any of "in".
	Operator string `json:"operator" api:"required"`
	// Values the request attribute must match. For `source_ip`, CIDR networks: a bare
	// address without a mask is rejected, and so are an entry with host bits set and
	// an IPv4 network written in IPv6 notation.
	Values []string `json:"values" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Key         respjson.Field
		Operator    respjson.Field
		Values      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKeyPolicyStatementCondition) RawJSON() string { return r.JSON.raw }
func (r *APIKeyPolicyStatementCondition) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type APIKeyPrincipal struct {
	ID string `json:"id" api:"required"`
	// Any of "user", "service_account".
	Type APIKeyPrincipalType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKeyPrincipal) RawJSON() string { return r.JSON.raw }
func (r *APIKeyPrincipal) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type APIKeyPrincipalType string

const (
	APIKeyPrincipalTypeUser           APIKeyPrincipalType = "user"
	APIKeyPrincipalTypeServiceAccount APIKeyPrincipalType = "service_account"
)

// API key returned immediately after creation. Includes the plaintext key once.
type CreatedAPIKey struct {
	// Plaintext API key. Only returned once when the key is created.
	Key string `json:"key" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Key         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
	APIKey
}

// Returns the unmodified JSON received from the API
func (r CreatedAPIKey) RawJSON() string { return r.JSON.raw }
func (r *CreatedAPIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type APIKeyNewParams struct {
	// Label for the API key (1-255 characters). API keys are not addressable by name.
	Name string `json:"name" api:"required"`
	// Number of days until expiry, up to 3650. Use null for never. Cannot be combined
	// with expires_at.
	DaysToExpire param.Opt[int64] `json:"days_to_expire,omitzero"`
	// Exact expiry time. Must be in the future and no later than the issuing
	// credential's own expiry. Cannot be combined with days_to_expire.
	ExpiresAt param.Opt[time.Time] `json:"expires_at,omitzero" format:"date-time"`
	// Unique project identifier
	ProjectID param.Opt[string] `json:"project_id,omitzero"`
	// Link to an active same-organization service account. Requires an organization
	// administrator, explicit project and nonempty policy. Omit for a personal key.
	ServiceAccountID param.Opt[string] `json:"service_account_id,omitzero"`
	// Issue a scoped key that may only perform the listed actions on the listed
	// resources. A statement may carry a condition; a `source_ip` condition restricts
	// the statement to requests from inside its CIDR networks.
	Policy APIKeyNewParamsPolicy `json:"policy,omitzero"`
	paramObj
}

func (r APIKeyNewParams) MarshalJSON() (data []byte, err error) {
	type shadow APIKeyNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *APIKeyNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Issue a scoped key that may only perform the listed actions on the listed
// resources. A statement may carry a condition; a `source_ip` condition restricts
// the statement to requests from inside its CIDR networks.
//
// The property Statements is required.
type APIKeyNewParamsPolicy struct {
	// Allow statements. Every action and resource not listed is denied.
	Statements []APIKeyNewParamsPolicyStatement `json:"statements,omitzero" api:"required"`
	paramObj
}

func (r APIKeyNewParamsPolicy) MarshalJSON() (data []byte, err error) {
	type shadow APIKeyNewParamsPolicy
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *APIKeyNewParamsPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Action, Resource are required.
type APIKeyNewParamsPolicyStatement struct {
	// Permitted action, such as browsers:use or vaults:read. When creating an API key
	// or replacing its policy, `*` may replace whole components as shorthand:
	// browsers:_ names every browsers action, _:delete every delete action the
	// resource can carry, and _:_ every action the resource can carry. The stored
	// policy lists the concrete actions, so actions added later are not granted, and
	// api_keys:issue_beyond_expiry must always be listed explicitly. Partial
	// components such as brow*:use and a bare * are rejected. Other policies, such as
	// workload identity rules and deployment runtime policies, accept only exact
	// actions.
	Action string `json:"action" api:"required"`
	// Resource the action applies to, in the format projects/{project_id}/{kind}/{id}.
	// Wildcards are supported.
	Resource string `json:"resource" api:"required"`
	// A condition narrows a policy statement to requests matching its constraint: the
	// request attribute, a comparison operator, and that operator's values.
	Condition APIKeyNewParamsPolicyStatementCondition `json:"condition,omitzero"`
	paramObj
}

func (r APIKeyNewParamsPolicyStatement) MarshalJSON() (data []byte, err error) {
	type shadow APIKeyNewParamsPolicyStatement
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *APIKeyNewParamsPolicyStatement) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A condition narrows a policy statement to requests matching its constraint: the
// request attribute, a comparison operator, and that operator's values.
//
// The properties Key, Operator, Values are required.
type APIKeyNewParamsPolicyStatementCondition struct {
	// Request attribute the condition constrains. `source_ip` restricts the statement
	// to requests whose client address falls inside one of `values`.
	//
	// Any of "source_ip".
	Key string `json:"key,omitzero" api:"required"`
	// Comparison the condition applies between the request attribute and values.
	//
	// Any of "in".
	Operator string `json:"operator,omitzero" api:"required"`
	// Values the request attribute must match. For `source_ip`, CIDR networks: a bare
	// address without a mask is rejected, and so are an entry with host bits set and
	// an IPv4 network written in IPv6 notation.
	Values []string `json:"values,omitzero" api:"required"`
	paramObj
}

func (r APIKeyNewParamsPolicyStatementCondition) MarshalJSON() (data []byte, err error) {
	type shadow APIKeyNewParamsPolicyStatementCondition
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *APIKeyNewParamsPolicyStatementCondition) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[APIKeyNewParamsPolicyStatementCondition](
		"key", "source_ip",
	)
	apijson.RegisterFieldValidator[APIKeyNewParamsPolicyStatementCondition](
		"operator", "in",
	)
}

type APIKeyGetParams struct {
	// When true, return the API key even if it has been deleted (soft-deleted), for
	// audit purposes. Defaults to false, which returns 404 for a deleted key.
	IncludeDeleted param.Opt[bool] `query:"include_deleted,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [APIKeyGetParams]'s query parameters as `url.Values`.
func (r APIKeyGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type APIKeyUpdateParams struct {
	// New API key name
	Name param.Opt[string] `json:"name,omitzero"`
	// Permissions that fully replace the key's policy. Attaching a policy to a key
	// without one makes the key scoped: every action and resource the policy does not
	// list is denied. The key's credential, name, project, and expiry are unchanged,
	// so consumers never need the key redistributed. A `source_ip` condition adds or
	// changes IP allowlisting without rotating the credential. Requires
	// api_keys:set_policy on the key.
	Policy APIKeyUpdateParamsPolicy `json:"policy,omitzero"`
	paramObj
}

func (r APIKeyUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow APIKeyUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *APIKeyUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Permissions that fully replace the key's policy. Attaching a policy to a key
// without one makes the key scoped: every action and resource the policy does not
// list is denied. The key's credential, name, project, and expiry are unchanged,
// so consumers never need the key redistributed. A `source_ip` condition adds or
// changes IP allowlisting without rotating the credential. Requires
// api_keys:set_policy on the key.
//
// The property Statements is required.
type APIKeyUpdateParamsPolicy struct {
	// Allow statements. Every action and resource not listed is denied.
	Statements []APIKeyUpdateParamsPolicyStatement `json:"statements,omitzero" api:"required"`
	paramObj
}

func (r APIKeyUpdateParamsPolicy) MarshalJSON() (data []byte, err error) {
	type shadow APIKeyUpdateParamsPolicy
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *APIKeyUpdateParamsPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Action, Resource are required.
type APIKeyUpdateParamsPolicyStatement struct {
	// Permitted action, such as browsers:use or vaults:read. When creating an API key
	// or replacing its policy, `*` may replace whole components as shorthand:
	// browsers:_ names every browsers action, _:delete every delete action the
	// resource can carry, and _:_ every action the resource can carry. The stored
	// policy lists the concrete actions, so actions added later are not granted, and
	// api_keys:issue_beyond_expiry must always be listed explicitly. Partial
	// components such as brow*:use and a bare * are rejected. Other policies, such as
	// workload identity rules and deployment runtime policies, accept only exact
	// actions.
	Action string `json:"action" api:"required"`
	// Resource the action applies to, in the format projects/{project_id}/{kind}/{id}.
	// Wildcards are supported.
	Resource string `json:"resource" api:"required"`
	// A condition narrows a policy statement to requests matching its constraint: the
	// request attribute, a comparison operator, and that operator's values.
	Condition APIKeyUpdateParamsPolicyStatementCondition `json:"condition,omitzero"`
	paramObj
}

func (r APIKeyUpdateParamsPolicyStatement) MarshalJSON() (data []byte, err error) {
	type shadow APIKeyUpdateParamsPolicyStatement
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *APIKeyUpdateParamsPolicyStatement) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A condition narrows a policy statement to requests matching its constraint: the
// request attribute, a comparison operator, and that operator's values.
//
// The properties Key, Operator, Values are required.
type APIKeyUpdateParamsPolicyStatementCondition struct {
	// Request attribute the condition constrains. `source_ip` restricts the statement
	// to requests whose client address falls inside one of `values`.
	//
	// Any of "source_ip".
	Key string `json:"key,omitzero" api:"required"`
	// Comparison the condition applies between the request attribute and values.
	//
	// Any of "in".
	Operator string `json:"operator,omitzero" api:"required"`
	// Values the request attribute must match. For `source_ip`, CIDR networks: a bare
	// address without a mask is rejected, and so are an entry with host bits set and
	// an IPv4 network written in IPv6 notation.
	Values []string `json:"values,omitzero" api:"required"`
	paramObj
}

func (r APIKeyUpdateParamsPolicyStatementCondition) MarshalJSON() (data []byte, err error) {
	type shadow APIKeyUpdateParamsPolicyStatementCondition
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *APIKeyUpdateParamsPolicyStatementCondition) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[APIKeyUpdateParamsPolicyStatementCondition](
		"key", "source_ip",
	)
	apijson.RegisterFieldValidator[APIKeyUpdateParamsPolicyStatementCondition](
		"operator", "in",
	)
}

type APIKeyListParams struct {
	// Return only keys created by this user ID. This is the same person the key's
	// created_by field describes.
	CreatedBy param.Opt[string] `query:"created_by,omitzero" json:"-"`
	// Deprecated: use status=all instead. When true, include deleted (soft-deleted)
	// API keys in the results for audit purposes.
	IncludeDeleted param.Opt[bool] `query:"include_deleted,omitzero" json:"-"`
	// Maximum number of results to return
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Exact-match filter on API key name using the database collation. In production,
	// matching is case- and accent-insensitive. Names are not required to be unique,
	// so multiple keys may match. When status=all or include_deleted=true is set,
	// soft-deleted keys with the same name may also match.
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	// Number of results to skip
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Return only keys scoped to this project ID. The value "org" returns keys that
	// are not scoped to a project.
	ProjectID param.Opt[string] `query:"project_id,omitzero" json:"-"`
	// Case-insensitive substring match against API key name, creator, and project. API
	// key identifiers and masked keys match by exact value or prefix.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
	// Return only keys that belong to this service account.
	ServiceAccountID param.Opt[string] `query:"service_account_id,omitzero" json:"-"`
	// Return only keys that act as a user, or only keys that act as a service account.
	//
	// Any of "user", "service_account".
	PrincipalType APIKeyListParamsPrincipalType `query:"principal_type,omitzero" json:"-"`
	// Field to sort API keys by. Keys with no expiry or no recorded use sort last for
	// expires_at and last_used_at in either direction.
	//
	// Any of "created_at", "name", "expires_at", "last_used_at".
	SortBy APIKeyListParamsSortBy `query:"sort_by,omitzero" json:"-"`
	// Sort direction for API keys.
	//
	// Any of "asc", "desc".
	SortDirection APIKeyListParamsSortDirection `query:"sort_direction,omitzero" json:"-"`
	// Filter API keys by status. "active" returns keys that are neither deleted nor
	// expired (default), "expired" returns keys that are not deleted and whose expiry
	// has passed, "deleted" returns only soft-deleted keys, "all" returns every key.
	//
	// Any of "active", "expired", "deleted", "all".
	Status APIKeyListParamsStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [APIKeyListParams]'s query parameters as `url.Values`.
func (r APIKeyListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Return only keys that act as a user, or only keys that act as a service account.
type APIKeyListParamsPrincipalType string

const (
	APIKeyListParamsPrincipalTypeUser           APIKeyListParamsPrincipalType = "user"
	APIKeyListParamsPrincipalTypeServiceAccount APIKeyListParamsPrincipalType = "service_account"
)

// Field to sort API keys by. Keys with no expiry or no recorded use sort last for
// expires_at and last_used_at in either direction.
type APIKeyListParamsSortBy string

const (
	APIKeyListParamsSortByCreatedAt  APIKeyListParamsSortBy = "created_at"
	APIKeyListParamsSortByName       APIKeyListParamsSortBy = "name"
	APIKeyListParamsSortByExpiresAt  APIKeyListParamsSortBy = "expires_at"
	APIKeyListParamsSortByLastUsedAt APIKeyListParamsSortBy = "last_used_at"
)

// Sort direction for API keys.
type APIKeyListParamsSortDirection string

const (
	APIKeyListParamsSortDirectionAsc  APIKeyListParamsSortDirection = "asc"
	APIKeyListParamsSortDirectionDesc APIKeyListParamsSortDirection = "desc"
)

// Filter API keys by status. "active" returns keys that are neither deleted nor
// expired (default), "expired" returns keys that are not deleted and whose expiry
// has passed, "deleted" returns only soft-deleted keys, "all" returns every key.
type APIKeyListParamsStatus string

const (
	APIKeyListParamsStatusActive  APIKeyListParamsStatus = "active"
	APIKeyListParamsStatusExpired APIKeyListParamsStatus = "expired"
	APIKeyListParamsStatusDeleted APIKeyListParamsStatus = "deleted"
	APIKeyListParamsStatusAll     APIKeyListParamsStatus = "all"
)

type APIKeyRotateParams struct {
	// Lifetime in days for the new key, up to 3650. Omit to reuse the rotated key's
	// original lifetime, or never-expires if it had none.
	DaysToExpire param.Opt[int64] `json:"days_to_expire,omitzero"`
	// Grace period in days before the rotated key expires. Use 0 to expire it
	// immediately. Omit for the default grace period of 7 days.
	ExpireInDays param.Opt[int64] `json:"expire_in_days,omitzero"`
	paramObj
}

func (r APIKeyRotateParams) MarshalJSON() (data []byte, err error) {
	type shadow APIKeyRotateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *APIKeyRotateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
