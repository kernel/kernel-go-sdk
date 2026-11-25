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
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// WorkloadIdentityFederationRuleService contains methods and other services that
// help with interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewWorkloadIdentityFederationRuleService] method instead.
type WorkloadIdentityFederationRuleService struct {
	Options []option.RequestOption
}

// NewWorkloadIdentityFederationRuleService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewWorkloadIdentityFederationRuleService(opts ...option.RequestOption) (r WorkloadIdentityFederationRuleService) {
	r = WorkloadIdentityFederationRuleService{}
	r.Options = opts
	return
}

// Organization credential administrators only. No external namespace-ownership or
// staff approval is required. New trust resources default disabled.
func (r *WorkloadIdentityFederationRuleService) New(ctx context.Context, body WorkloadIdentityFederationRuleNewParams, opts ...option.RequestOption) (res *WorkloadIdentityFederationRule, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "org/workload_identity_federation_rules"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Any signed-in member of the organization can read this. API keys cannot.
func (r *WorkloadIdentityFederationRuleService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *WorkloadIdentityFederationRule, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/workload_identity_federation_rules/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a workload identity federation rule
func (r *WorkloadIdentityFederationRuleService) Update(ctx context.Context, id string, body WorkloadIdentityFederationRuleUpdateParams, opts ...option.RequestOption) (res *WorkloadIdentityFederationRule, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/workload_identity_federation_rules/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// Any signed-in member of the organization can read this. API keys cannot.
func (r *WorkloadIdentityFederationRuleService) List(ctx context.Context, query WorkloadIdentityFederationRuleListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[WorkloadIdentityFederationRule], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "org/workload_identity_federation_rules"
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

// Any signed-in member of the organization can read this. API keys cannot.
func (r *WorkloadIdentityFederationRuleService) ListAutoPaging(ctx context.Context, query WorkloadIdentityFederationRuleListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[WorkloadIdentityFederationRule] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Retains history and invalidates associated authority within the documented
// status bound. Does not delete created browser resources.
func (r *WorkloadIdentityFederationRuleService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("org/workload_identity_federation_rules/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type WorkloadIdentityAuthorization struct {
	// Allow statements for resources in this rule's project, in the same format as a
	// scoped API key policy. API key actions are not allowed.
	Policy WorkloadIdentityAuthorizationPolicy `json:"policy" api:"required"`
	// The one project credentials from this rule are limited to.
	ProjectID          string `json:"project_id" api:"required"`
	TokenAudience      string `json:"token_audience" api:"required"`
	MaxTokenTtlSeconds int64  `json:"max_token_ttl_seconds"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Policy             respjson.Field
		ProjectID          respjson.Field
		TokenAudience      respjson.Field
		MaxTokenTtlSeconds respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityAuthorization) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityAuthorization) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this WorkloadIdentityAuthorization to a
// WorkloadIdentityAuthorizationParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// WorkloadIdentityAuthorizationParam.Overrides()
func (r WorkloadIdentityAuthorization) ToParam() WorkloadIdentityAuthorizationParam {
	return param.Override[WorkloadIdentityAuthorizationParam](json.RawMessage(r.RawJSON()))
}

// Allow statements for resources in this rule's project, in the same format as a
// scoped API key policy. API key actions are not allowed.
type WorkloadIdentityAuthorizationPolicy struct {
	// Allow statements. Every action and resource not listed is denied.
	Statements []WorkloadIdentityAuthorizationPolicyStatement `json:"statements" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Statements  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityAuthorizationPolicy) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityAuthorizationPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityAuthorizationPolicyStatement struct {
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
	Condition WorkloadIdentityAuthorizationPolicyStatementCondition `json:"condition"`
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
func (r WorkloadIdentityAuthorizationPolicyStatement) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityAuthorizationPolicyStatement) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A condition narrows a policy statement to requests matching its constraint: the
// request attribute, a comparison operator, and that operator's values.
type WorkloadIdentityAuthorizationPolicyStatementCondition struct {
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
func (r WorkloadIdentityAuthorizationPolicyStatementCondition) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityAuthorizationPolicyStatementCondition) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Policy, ProjectID, TokenAudience are required.
type WorkloadIdentityAuthorizationParam struct {
	// Allow statements for resources in this rule's project, in the same format as a
	// scoped API key policy. API key actions are not allowed.
	Policy WorkloadIdentityAuthorizationPolicyParam `json:"policy,omitzero" api:"required"`
	// The one project credentials from this rule are limited to.
	ProjectID          string           `json:"project_id" api:"required"`
	TokenAudience      string           `json:"token_audience" api:"required"`
	MaxTokenTtlSeconds param.Opt[int64] `json:"max_token_ttl_seconds,omitzero"`
	paramObj
}

func (r WorkloadIdentityAuthorizationParam) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityAuthorizationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityAuthorizationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Allow statements for resources in this rule's project, in the same format as a
// scoped API key policy. API key actions are not allowed.
//
// The property Statements is required.
type WorkloadIdentityAuthorizationPolicyParam struct {
	// Allow statements. Every action and resource not listed is denied.
	Statements []WorkloadIdentityAuthorizationPolicyStatementParam `json:"statements,omitzero" api:"required"`
	paramObj
}

func (r WorkloadIdentityAuthorizationPolicyParam) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityAuthorizationPolicyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityAuthorizationPolicyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Action, Resource are required.
type WorkloadIdentityAuthorizationPolicyStatementParam struct {
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
	Condition WorkloadIdentityAuthorizationPolicyStatementConditionParam `json:"condition,omitzero"`
	paramObj
}

func (r WorkloadIdentityAuthorizationPolicyStatementParam) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityAuthorizationPolicyStatementParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityAuthorizationPolicyStatementParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A condition narrows a policy statement to requests matching its constraint: the
// request attribute, a comparison operator, and that operator's values.
//
// The properties Key, Operator, Values are required.
type WorkloadIdentityAuthorizationPolicyStatementConditionParam struct {
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

func (r WorkloadIdentityAuthorizationPolicyStatementConditionParam) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityAuthorizationPolicyStatementConditionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityAuthorizationPolicyStatementConditionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[WorkloadIdentityAuthorizationPolicyStatementConditionParam](
		"key", "source_ip",
	)
	apijson.RegisterFieldValidator[WorkloadIdentityAuthorizationPolicyStatementConditionParam](
		"operator", "in",
	)
}

type WorkloadIdentityFederationRule struct {
	ID               string                        `json:"id" api:"required"`
	ArchivedAt       time.Time                     `json:"archived_at" api:"required" format:"date-time"`
	Authorization    WorkloadIdentityAuthorization `json:"authorization" api:"required"`
	CreatedAt        time.Time                     `json:"created_at" api:"required" format:"date-time"`
	CreatedByUserID  string                        `json:"created_by_user_id" api:"required"`
	Description      string                        `json:"description" api:"required"`
	Enabled          bool                          `json:"enabled" api:"required"`
	ExpectedAudience string                        `json:"expected_audience" api:"required"`
	IssuerID         string                        `json:"issuer_id" api:"required"`
	LastExchangeAt   time.Time                     `json:"last_exchange_at" api:"required" format:"date-time"`
	Match            WorkloadIdentityMatch         `json:"match" api:"required"`
	Name             string                        `json:"name" api:"required"`
	OrgID            string                        `json:"org_id" api:"required"`
	// Authority revision. Cosmetic updates do not widen permissions.
	SecurityVersion int64                  `json:"security_version" api:"required"`
	Target          WorkloadIdentityTarget `json:"target" api:"required"`
	UpdatedAt       time.Time              `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		ArchivedAt       respjson.Field
		Authorization    respjson.Field
		CreatedAt        respjson.Field
		CreatedByUserID  respjson.Field
		Description      respjson.Field
		Enabled          respjson.Field
		ExpectedAudience respjson.Field
		IssuerID         respjson.Field
		LastExchangeAt   respjson.Field
		Match            respjson.Field
		Name             respjson.Field
		OrgID            respjson.Field
		SecurityVersion  respjson.Field
		Target           respjson.Field
		UpdatedAt        respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityFederationRule) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityFederationRule) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Authorization, ExpectedAudience, IssuerID, Match, Name, Target
// are required.
type WorkloadIdentityFederationRuleConfigParam struct {
	Authorization    WorkloadIdentityAuthorizationParam `json:"authorization,omitzero" api:"required"`
	ExpectedAudience string                             `json:"expected_audience" api:"required"`
	IssuerID         string                             `json:"issuer_id" api:"required"`
	Match            WorkloadIdentityMatchParam         `json:"match,omitzero" api:"required"`
	Name             string                             `json:"name" api:"required"`
	Target           WorkloadIdentityTargetParam        `json:"target,omitzero" api:"required"`
	Description      param.Opt[string]                  `json:"description,omitzero"`
	Enabled          param.Opt[bool]                    `json:"enabled,omitzero"`
	paramObj
}

func (r WorkloadIdentityFederationRuleConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityFederationRuleConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityFederationRuleConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityMatch struct {
	// Must be the issuer's preset_type.
	//
	// Any of "github_actions", "aws", "google_cloud", "microsoft_entra", "kubernetes",
	// "custom_oidc".
	Type WorkloadIdentityMatchType `json:"type" api:"required"`
	// Exact string conditions keyed by RFC 6901 JSON Pointers. All conditions are
	// ANDed.
	Claims         map[string]string `json:"claims"`
	Environment    string            `json:"environment"`
	EventNames     []string          `json:"event_names"`
	JobWorkflowRef string            `json:"job_workflow_ref"`
	Namespace      string            `json:"namespace"`
	Ref            string            `json:"ref"`
	// Friendly display name. Immutable IDs define the repository identity.
	Repository        string `json:"repository"`
	RepositoryID      string `json:"repository_id"`
	RepositoryOwnerID string `json:"repository_owner_id"`
	ServiceAccount    string `json:"service_account"`
	Subject           string `json:"subject"`
	SubjectPrefix     string `json:"subject_prefix"`
	Workflow          string `json:"workflow"`
	WorkflowRef       string `json:"workflow_ref"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type              respjson.Field
		Claims            respjson.Field
		Environment       respjson.Field
		EventNames        respjson.Field
		JobWorkflowRef    respjson.Field
		Namespace         respjson.Field
		Ref               respjson.Field
		Repository        respjson.Field
		RepositoryID      respjson.Field
		RepositoryOwnerID respjson.Field
		ServiceAccount    respjson.Field
		Subject           respjson.Field
		SubjectPrefix     respjson.Field
		Workflow          respjson.Field
		WorkflowRef       respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityMatch) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityMatch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this WorkloadIdentityMatch to a WorkloadIdentityMatchParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// WorkloadIdentityMatchParam.Overrides()
func (r WorkloadIdentityMatch) ToParam() WorkloadIdentityMatchParam {
	return param.Override[WorkloadIdentityMatchParam](json.RawMessage(r.RawJSON()))
}

// Must be the issuer's preset_type.
type WorkloadIdentityMatchType string

const (
	WorkloadIdentityMatchTypeGitHubActions  WorkloadIdentityMatchType = "github_actions"
	WorkloadIdentityMatchTypeAws            WorkloadIdentityMatchType = "aws"
	WorkloadIdentityMatchTypeGoogleCloud    WorkloadIdentityMatchType = "google_cloud"
	WorkloadIdentityMatchTypeMicrosoftEntra WorkloadIdentityMatchType = "microsoft_entra"
	WorkloadIdentityMatchTypeKubernetes     WorkloadIdentityMatchType = "kubernetes"
	WorkloadIdentityMatchTypeCustomOidc     WorkloadIdentityMatchType = "custom_oidc"
)

// The property Type is required.
type WorkloadIdentityMatchParam struct {
	// Must be the issuer's preset_type.
	//
	// Any of "github_actions", "aws", "google_cloud", "microsoft_entra", "kubernetes",
	// "custom_oidc".
	Type           WorkloadIdentityMatchType `json:"type,omitzero" api:"required"`
	Environment    param.Opt[string]         `json:"environment,omitzero"`
	JobWorkflowRef param.Opt[string]         `json:"job_workflow_ref,omitzero"`
	Namespace      param.Opt[string]         `json:"namespace,omitzero"`
	Ref            param.Opt[string]         `json:"ref,omitzero"`
	// Friendly display name. Immutable IDs define the repository identity.
	Repository        param.Opt[string] `json:"repository,omitzero"`
	RepositoryID      param.Opt[string] `json:"repository_id,omitzero"`
	RepositoryOwnerID param.Opt[string] `json:"repository_owner_id,omitzero"`
	ServiceAccount    param.Opt[string] `json:"service_account,omitzero"`
	Subject           param.Opt[string] `json:"subject,omitzero"`
	SubjectPrefix     param.Opt[string] `json:"subject_prefix,omitzero"`
	Workflow          param.Opt[string] `json:"workflow,omitzero"`
	WorkflowRef       param.Opt[string] `json:"workflow_ref,omitzero"`
	// Exact string conditions keyed by RFC 6901 JSON Pointers. All conditions are
	// ANDed.
	Claims     map[string]string `json:"claims,omitzero"`
	EventNames []string          `json:"event_names,omitzero"`
	paramObj
}

func (r WorkloadIdentityMatchParam) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityMatchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityMatchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityTarget struct {
	ServiceAccountID string `json:"service_account_id" api:"required"`
	// Any of "service_account".
	Type WorkloadIdentityTargetType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ServiceAccountID respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityTarget) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityTarget) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this WorkloadIdentityTarget to a WorkloadIdentityTargetParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// WorkloadIdentityTargetParam.Overrides()
func (r WorkloadIdentityTarget) ToParam() WorkloadIdentityTargetParam {
	return param.Override[WorkloadIdentityTargetParam](json.RawMessage(r.RawJSON()))
}

type WorkloadIdentityTargetType string

const (
	WorkloadIdentityTargetTypeServiceAccount WorkloadIdentityTargetType = "service_account"
)

// The properties ServiceAccountID, Type are required.
type WorkloadIdentityTargetParam struct {
	ServiceAccountID string `json:"service_account_id" api:"required"`
	// Any of "service_account".
	Type WorkloadIdentityTargetType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r WorkloadIdentityTargetParam) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityTargetParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityTargetParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityFederationRuleNewParams struct {
	WorkloadIdentityFederationRuleConfig WorkloadIdentityFederationRuleConfigParam
	paramObj
}

func (r WorkloadIdentityFederationRuleNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.WorkloadIdentityFederationRuleConfig)
}
func (r *WorkloadIdentityFederationRuleNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityFederationRuleUpdateParams struct {
	ExpectedVersion  int64                              `json:"expected_version" api:"required"`
	Description      param.Opt[string]                  `json:"description,omitzero"`
	Enabled          param.Opt[bool]                    `json:"enabled,omitzero"`
	ExpectedAudience param.Opt[string]                  `json:"expected_audience,omitzero"`
	IssuerID         param.Opt[string]                  `json:"issuer_id,omitzero"`
	Name             param.Opt[string]                  `json:"name,omitzero"`
	Authorization    WorkloadIdentityAuthorizationParam `json:"authorization,omitzero"`
	Match            WorkloadIdentityMatchParam         `json:"match,omitzero"`
	Target           WorkloadIdentityTargetParam        `json:"target,omitzero"`
	paramObj
}

func (r WorkloadIdentityFederationRuleUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityFederationRuleUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityFederationRuleUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityFederationRuleListParams struct {
	IssuerID  param.Opt[string] `query:"issuer_id,omitzero" json:"-"`
	Limit     param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	Offset    param.Opt[int64]  `query:"offset,omitzero" json:"-"`
	ProjectID param.Opt[string] `query:"project_id,omitzero" json:"-"`
	// Case-insensitive substring match against name or description; exact ID match.
	Query            param.Opt[string] `query:"query,omitzero" json:"-"`
	ServiceAccountID param.Opt[string] `query:"service_account_id,omitzero" json:"-"`
	// Any of "active", "archived", "all".
	Status WorkloadIdentityFederationRuleListParamsStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [WorkloadIdentityFederationRuleListParams]'s query
// parameters as `url.Values`.
func (r WorkloadIdentityFederationRuleListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type WorkloadIdentityFederationRuleListParamsStatus string

const (
	WorkloadIdentityFederationRuleListParamsStatusActive   WorkloadIdentityFederationRuleListParamsStatus = "active"
	WorkloadIdentityFederationRuleListParamsStatusArchived WorkloadIdentityFederationRuleListParamsStatus = "archived"
	WorkloadIdentityFederationRuleListParamsStatusAll      WorkloadIdentityFederationRuleListParamsStatus = "all"
)
