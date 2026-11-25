// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"net/http"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Inspect the identity and authorization context for the current request.
//
// AuthContextService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAuthContextService] method instead.
type AuthContextService struct {
	Options []option.RequestOption
}

// NewAuthContextService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewAuthContextService(opts ...option.RequestOption) (r AuthContextService) {
	r = AuthContextService{}
	r.Options = opts
	return
}

// Returns the authenticated principal, organization, credential scope, and
// effective request scope. The response is derived from the verified request
// context and does not expose credential secrets.
func (r *AuthContextService) Get(ctx context.Context, opts ...option.RequestOption) (res *AuthContext, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "auth/context"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// The identity and authorization context resolved for the current request.
type AuthContext struct {
	Authentication AuthContextAuthentication `json:"authentication" api:"required"`
	// The credential's maximum scope and the effective scope selected for this
	// request. Future permission data can be added without changing scope semantics.
	Authorization    AuthContextAuthorization    `json:"authorization" api:"required"`
	Organization     AuthContextOrganization     `json:"organization" api:"required"`
	Principal        AuthContextPrincipal        `json:"principal" api:"required"`
	WorkloadIdentity AuthContextWorkloadIdentity `json:"workload_identity"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Authentication   respjson.Field
		Authorization    respjson.Field
		Organization     respjson.Field
		Principal        respjson.Field
		WorkloadIdentity respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthContext) RawJSON() string { return r.JSON.raw }
func (r *AuthContext) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthContextAuthentication struct {
	// The API key ID when authenticated with an API key; null for session credentials.
	CredentialID string `json:"credential_id" api:"required"`
	// The credential format used to authenticate the request.
	//
	// Any of "api_key", "jwt".
	Method string `json:"method" api:"required"`
	// The source classification resolved by authentication middleware.
	//
	// Any of "api_key", "oauth", "dashboard".
	Source string `json:"source" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CredentialID respjson.Field
		Method       respjson.Field
		Source       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthContextAuthentication) RawJSON() string { return r.JSON.raw }
func (r *AuthContextAuthentication) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The credential's maximum scope and the effective scope selected for this
// request. Future permission data can be added without changing scope semantics.
type AuthContextAuthorization struct {
	// A scope within the authenticated organization. A null project_id represents
	// organization-wide scope.
	CredentialScope AuthContextAuthorizationCredentialScope `json:"credential_scope" api:"required"`
	// A scope within the authenticated organization. A null project_id represents
	// organization-wide scope.
	EffectiveScope AuthContextAuthorizationEffectiveScope `json:"effective_scope" api:"required"`
	// Current effective paired policy for a service-account principal. Omitted for
	// legacy human/client identities.
	EffectivePolicy AuthContextAuthorizationEffectivePolicy `json:"effective_policy"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CredentialScope respjson.Field
		EffectiveScope  respjson.Field
		EffectivePolicy respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthContextAuthorization) RawJSON() string { return r.JSON.raw }
func (r *AuthContextAuthorization) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A scope within the authenticated organization. A null project_id represents
// organization-wide scope.
type AuthContextAuthorizationCredentialScope struct {
	// The Kernel project ID, or null when the scope is organization-wide.
	ProjectID string `json:"project_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthContextAuthorizationCredentialScope) RawJSON() string { return r.JSON.raw }
func (r *AuthContextAuthorizationCredentialScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A scope within the authenticated organization. A null project_id represents
// organization-wide scope.
type AuthContextAuthorizationEffectiveScope struct {
	// The Kernel project ID, or null when the scope is organization-wide.
	ProjectID string `json:"project_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthContextAuthorizationEffectiveScope) RawJSON() string { return r.JSON.raw }
func (r *AuthContextAuthorizationEffectiveScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current effective paired policy for a service-account principal. Omitted for
// legacy human/client identities.
type AuthContextAuthorizationEffectivePolicy struct {
	// Allow statements. Every action and resource not listed is denied.
	Statements []AuthContextAuthorizationEffectivePolicyStatement `json:"statements" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Statements  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthContextAuthorizationEffectivePolicy) RawJSON() string { return r.JSON.raw }
func (r *AuthContextAuthorizationEffectivePolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthContextAuthorizationEffectivePolicyStatement struct {
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
	Condition AuthContextAuthorizationEffectivePolicyStatementCondition `json:"condition"`
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
func (r AuthContextAuthorizationEffectivePolicyStatement) RawJSON() string { return r.JSON.raw }
func (r *AuthContextAuthorizationEffectivePolicyStatement) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A condition narrows a policy statement to requests matching its constraint: the
// request attribute, a comparison operator, and that operator's values.
type AuthContextAuthorizationEffectivePolicyStatementCondition struct {
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
func (r AuthContextAuthorizationEffectivePolicyStatementCondition) RawJSON() string {
	return r.JSON.raw
}
func (r *AuthContextAuthorizationEffectivePolicyStatementCondition) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthContextOrganization struct {
	// The authenticated Kernel organization ID.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthContextOrganization) RawJSON() string { return r.JSON.raw }
func (r *AuthContextOrganization) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthContextPrincipal struct {
	// The authenticated API key, user or organization-owned service-account ID.
	ID string `json:"id" api:"required"`
	// The kind of principal authenticated for the request.
	//
	// Any of "api_key", "user", "service_account".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthContextPrincipal) RawJSON() string { return r.JSON.raw }
func (r *AuthContextPrincipal) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthContextWorkloadIdentity struct {
	// Digest of the verified external subject. Never the original assertion or
	// arbitrary claims.
	ExternalSubjectDigest string `json:"external_subject_digest" api:"required"`
	FederationRuleID      string `json:"federation_rule_id" api:"required"`
	IssuerID              string `json:"issuer_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExternalSubjectDigest respjson.Field
		FederationRuleID      respjson.Field
		IssuerID              respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthContextWorkloadIdentity) RawJSON() string { return r.JSON.raw }
func (r *AuthContextWorkloadIdentity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
