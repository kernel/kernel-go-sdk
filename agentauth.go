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

	"github.com/onkernel/kernel-go-sdk/internal/apijson"
	"github.com/onkernel/kernel-go-sdk/internal/apiquery"
	shimjson "github.com/onkernel/kernel-go-sdk/internal/encoding/json"
	"github.com/onkernel/kernel-go-sdk/internal/requestconfig"
	"github.com/onkernel/kernel-go-sdk/option"
	"github.com/onkernel/kernel-go-sdk/packages/pagination"
	"github.com/onkernel/kernel-go-sdk/packages/param"
	"github.com/onkernel/kernel-go-sdk/packages/respjson"
)

// AgentAuthService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAgentAuthService] method instead.
type AgentAuthService struct {
	Options     []option.RequestOption
	Invocations AgentAuthInvocationService
}

// NewAgentAuthService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewAgentAuthService(opts ...option.RequestOption) (r AgentAuthService) {
	r = AgentAuthService{}
	r.Options = opts
	r.Invocations = NewAgentAuthInvocationService(opts...)
	return
}

// Creates a new auth agent for the specified domain and profile combination, or
// returns an existing one if it already exists. This is idempotent - calling with
// the same domain and profile will return the same agent. Does NOT start an
// invocation - use POST /agents/auth/invocations to start an auth flow.
func (r *AgentAuthService) New(ctx context.Context, body AgentAuthNewParams, opts ...option.RequestOption) (res *AuthAgent, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "agents/auth"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieve an auth agent by its ID. Returns the current authentication status of
// the managed profile.
func (r *AgentAuthService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *AuthAgent, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("agents/auth/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List auth agents with optional filters for profile_name and domain.
func (r *AgentAuthService) List(ctx context.Context, query AgentAuthListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[AuthAgent], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "agents/auth"
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

// List auth agents with optional filters for profile_name and domain.
func (r *AgentAuthService) ListAutoPaging(ctx context.Context, query AgentAuthListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[AuthAgent] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Deletes an auth agent and terminates its workflow. This will:
//
// - Soft delete the auth agent record
// - Gracefully terminate the agent's Temporal workflow
// - Cancel any in-progress invocations
func (r *AgentAuthService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("agents/auth/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return
}

// Response from get invocation endpoint
type AgentAuthInvocationResponse struct {
	// App name (org name at time of invocation creation)
	AppName string `json:"app_name,required"`
	// Domain for authentication
	Domain string `json:"domain,required"`
	// When the handoff code expires
	ExpiresAt time.Time `json:"expires_at,required" format:"date-time"`
	// Invocation status
	//
	// Any of "IN_PROGRESS", "SUCCESS", "EXPIRED", "CANCELED", "FAILED".
	Status AgentAuthInvocationResponseStatus `json:"status,required"`
	// Current step in the invocation workflow
	//
	// Any of "initialized", "discovering", "awaiting_input",
	// "awaiting_external_action", "submitting", "completed", "expired".
	Step AgentAuthInvocationResponseStep `json:"step,required"`
	// The invocation type:
	//
	// - login: First-time authentication
	// - reauth: Re-authentication for previously authenticated agents
	// - auto_login: Legacy type (no longer created, kept for backward compatibility)
	//
	// Any of "login", "auto_login", "reauth".
	Type AgentAuthInvocationResponseType `json:"type,required"`
	// Error message explaining why the invocation failed (present when status=FAILED)
	ErrorMessage string `json:"error_message,nullable"`
	// Instructions for user when external action is required (present when
	// step=awaiting_external_action)
	ExternalActionMessage string `json:"external_action_message,nullable"`
	// Browser live view URL for debugging the invocation
	LiveViewURL string `json:"live_view_url,nullable"`
	// Fields currently awaiting input (present when step=awaiting_input)
	PendingFields []DiscoveredField `json:"pending_fields,nullable"`
	// SSO buttons available on the page (present when step=awaiting_input)
	PendingSSOButtons []AgentAuthInvocationResponsePendingSSOButton `json:"pending_sso_buttons,nullable"`
	// Names of fields that have been submitted (present when step=submitting or later)
	SubmittedFields []string `json:"submitted_fields,nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AppName               respjson.Field
		Domain                respjson.Field
		ExpiresAt             respjson.Field
		Status                respjson.Field
		Step                  respjson.Field
		Type                  respjson.Field
		ErrorMessage          respjson.Field
		ExternalActionMessage respjson.Field
		LiveViewURL           respjson.Field
		PendingFields         respjson.Field
		PendingSSOButtons     respjson.Field
		SubmittedFields       respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentAuthInvocationResponse) RawJSON() string { return r.JSON.raw }
func (r *AgentAuthInvocationResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Invocation status
type AgentAuthInvocationResponseStatus string

const (
	AgentAuthInvocationResponseStatusInProgress AgentAuthInvocationResponseStatus = "IN_PROGRESS"
	AgentAuthInvocationResponseStatusSuccess    AgentAuthInvocationResponseStatus = "SUCCESS"
	AgentAuthInvocationResponseStatusExpired    AgentAuthInvocationResponseStatus = "EXPIRED"
	AgentAuthInvocationResponseStatusCanceled   AgentAuthInvocationResponseStatus = "CANCELED"
	AgentAuthInvocationResponseStatusFailed     AgentAuthInvocationResponseStatus = "FAILED"
)

// Current step in the invocation workflow
type AgentAuthInvocationResponseStep string

const (
	AgentAuthInvocationResponseStepInitialized            AgentAuthInvocationResponseStep = "initialized"
	AgentAuthInvocationResponseStepDiscovering            AgentAuthInvocationResponseStep = "discovering"
	AgentAuthInvocationResponseStepAwaitingInput          AgentAuthInvocationResponseStep = "awaiting_input"
	AgentAuthInvocationResponseStepAwaitingExternalAction AgentAuthInvocationResponseStep = "awaiting_external_action"
	AgentAuthInvocationResponseStepSubmitting             AgentAuthInvocationResponseStep = "submitting"
	AgentAuthInvocationResponseStepCompleted              AgentAuthInvocationResponseStep = "completed"
	AgentAuthInvocationResponseStepExpired                AgentAuthInvocationResponseStep = "expired"
)

// The invocation type:
//
// - login: First-time authentication
// - reauth: Re-authentication for previously authenticated agents
// - auto_login: Legacy type (no longer created, kept for backward compatibility)
type AgentAuthInvocationResponseType string

const (
	AgentAuthInvocationResponseTypeLogin     AgentAuthInvocationResponseType = "login"
	AgentAuthInvocationResponseTypeAutoLogin AgentAuthInvocationResponseType = "auto_login"
	AgentAuthInvocationResponseTypeReauth    AgentAuthInvocationResponseType = "reauth"
)

// An SSO button for signing in with an external identity provider
type AgentAuthInvocationResponsePendingSSOButton struct {
	// Visible button text
	Label string `json:"label,required"`
	// Identity provider name
	Provider string `json:"provider,required"`
	// XPath selector for the button
	Selector string `json:"selector,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Label       respjson.Field
		Provider    respjson.Field
		Selector    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentAuthInvocationResponsePendingSSOButton) RawJSON() string { return r.JSON.raw }
func (r *AgentAuthInvocationResponsePendingSSOButton) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Response from submit endpoint - returns immediately after submission is accepted
type AgentAuthSubmitResponse struct {
	// Whether the submission was accepted for processing
	Accepted bool `json:"accepted,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Accepted    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentAuthSubmitResponse) RawJSON() string { return r.JSON.raw }
func (r *AgentAuthSubmitResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An auth agent that manages authentication for a specific domain and profile
// combination
type AuthAgent struct {
	// Unique identifier for the auth agent
	ID string `json:"id,required"`
	// Target domain for authentication
	Domain string `json:"domain,required"`
	// Name of the profile associated with this auth agent
	ProfileName string `json:"profile_name,required"`
	// Current authentication status of the managed profile
	//
	// Any of "AUTHENTICATED", "NEEDS_AUTH".
	Status AuthAgentStatus `json:"status,required"`
	// Additional domains that are valid for this auth agent's authentication flow
	// (besides the primary domain). Useful when login pages redirect to different
	// domains.
	AllowedDomains []string `json:"allowed_domains"`
	// Whether automatic re-authentication is possible (has credential_id, selectors,
	// and login_url)
	CanReauth bool `json:"can_reauth"`
	// ID of the linked credential for automatic re-authentication
	CredentialID string `json:"credential_id"`
	// Name of the linked credential for automatic re-authentication
	CredentialName string `json:"credential_name"`
	// Whether this auth agent has stored selectors for deterministic re-authentication
	HasSelectors bool `json:"has_selectors"`
	// When the last authentication check was performed
	LastAuthCheckAt time.Time `json:"last_auth_check_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		Domain          respjson.Field
		ProfileName     respjson.Field
		Status          respjson.Field
		AllowedDomains  respjson.Field
		CanReauth       respjson.Field
		CredentialID    respjson.Field
		CredentialName  respjson.Field
		HasSelectors    respjson.Field
		LastAuthCheckAt respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthAgent) RawJSON() string { return r.JSON.raw }
func (r *AuthAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current authentication status of the managed profile
type AuthAgentStatus string

const (
	AuthAgentStatusAuthenticated AuthAgentStatus = "AUTHENTICATED"
	AuthAgentStatusNeedsAuth     AuthAgentStatus = "NEEDS_AUTH"
)

// Request to create or find an auth agent
//
// The properties Domain, ProfileName are required.
type AuthAgentCreateRequestParam struct {
	// Domain for authentication
	Domain string `json:"domain,required"`
	// Name of the profile to use for this auth agent
	ProfileName string `json:"profile_name,required"`
	// Optional name of an existing credential to use for this auth agent. If provided,
	// the credential will be linked to the agent and its values will be used to
	// auto-fill the login form on invocation.
	CredentialName param.Opt[string] `json:"credential_name,omitzero"`
	// Optional login page URL. If provided, will be stored on the agent and used to
	// skip discovery in future invocations.
	LoginURL param.Opt[string] `json:"login_url,omitzero" format:"uri"`
	// Additional domains that are valid for this auth agent's authentication flow
	// (besides the primary domain). Useful when login pages redirect to different
	// domains.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	// Optional proxy configuration
	Proxy AuthAgentCreateRequestProxyParam `json:"proxy,omitzero"`
	paramObj
}

func (r AuthAgentCreateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow AuthAgentCreateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AuthAgentCreateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Optional proxy configuration
type AuthAgentCreateRequestProxyParam struct {
	// ID of the proxy to use
	ProxyID param.Opt[string] `json:"proxy_id,omitzero"`
	paramObj
}

func (r AuthAgentCreateRequestProxyParam) MarshalJSON() (data []byte, err error) {
	type shadow AuthAgentCreateRequestProxyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AuthAgentCreateRequestProxyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request to create an invocation for an existing auth agent
//
// The property AuthAgentID is required.
type AuthAgentInvocationCreateRequestParam struct {
	// ID of the auth agent to create an invocation for
	AuthAgentID string `json:"auth_agent_id,required"`
	// If provided, saves the submitted credentials under this name upon successful
	// login. The credential will be linked to the auth agent for automatic
	// re-authentication.
	SaveCredentialAs param.Opt[string] `json:"save_credential_as,omitzero"`
	paramObj
}

func (r AuthAgentInvocationCreateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow AuthAgentInvocationCreateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AuthAgentInvocationCreateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Response from creating an invocation. Always returns an invocation_id.
type AuthAgentInvocationCreateResponse struct {
	// When the handoff code expires.
	ExpiresAt time.Time `json:"expires_at,required" format:"date-time"`
	// One-time code for handoff.
	HandoffCode string `json:"handoff_code,required"`
	// URL to redirect user to.
	HostedURL string `json:"hosted_url,required" format:"uri"`
	// Unique identifier for the invocation.
	InvocationID string `json:"invocation_id,required"`
	// The invocation type:
	//
	// - login: First-time authentication
	// - reauth: Re-authentication for previously authenticated agents
	// - auto_login: Legacy type (no longer created, kept for backward compatibility)
	//
	// Any of "login", "auto_login", "reauth".
	Type AuthAgentInvocationCreateResponseType `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExpiresAt    respjson.Field
		HandoffCode  respjson.Field
		HostedURL    respjson.Field
		InvocationID respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthAgentInvocationCreateResponse) RawJSON() string { return r.JSON.raw }
func (r *AuthAgentInvocationCreateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The invocation type:
//
// - login: First-time authentication
// - reauth: Re-authentication for previously authenticated agents
// - auto_login: Legacy type (no longer created, kept for backward compatibility)
type AuthAgentInvocationCreateResponseType string

const (
	AuthAgentInvocationCreateResponseTypeLogin     AuthAgentInvocationCreateResponseType = "login"
	AuthAgentInvocationCreateResponseTypeAutoLogin AuthAgentInvocationCreateResponseType = "auto_login"
	AuthAgentInvocationCreateResponseTypeReauth    AuthAgentInvocationCreateResponseType = "reauth"
)

// A discovered form field
type DiscoveredField struct {
	// Field label
	Label string `json:"label,required"`
	// Field name
	Name string `json:"name,required"`
	// CSS selector for the field
	Selector string `json:"selector,required"`
	// Field type
	//
	// Any of "text", "email", "password", "tel", "number", "url", "code", "totp".
	Type DiscoveredFieldType `json:"type,required"`
	// Field placeholder
	Placeholder string `json:"placeholder"`
	// Whether field is required
	Required bool `json:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Label       respjson.Field
		Name        respjson.Field
		Selector    respjson.Field
		Type        respjson.Field
		Placeholder respjson.Field
		Required    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DiscoveredField) RawJSON() string { return r.JSON.raw }
func (r *DiscoveredField) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Field type
type DiscoveredFieldType string

const (
	DiscoveredFieldTypeText     DiscoveredFieldType = "text"
	DiscoveredFieldTypeEmail    DiscoveredFieldType = "email"
	DiscoveredFieldTypePassword DiscoveredFieldType = "password"
	DiscoveredFieldTypeTel      DiscoveredFieldType = "tel"
	DiscoveredFieldTypeNumber   DiscoveredFieldType = "number"
	DiscoveredFieldTypeURL      DiscoveredFieldType = "url"
	DiscoveredFieldTypeCode     DiscoveredFieldType = "code"
	DiscoveredFieldTypeTotp     DiscoveredFieldType = "totp"
)

type AgentAuthNewParams struct {
	// Request to create or find an auth agent
	AuthAgentCreateRequest AuthAgentCreateRequestParam
	paramObj
}

func (r AgentAuthNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.AuthAgentCreateRequest)
}
func (r *AgentAuthNewParams) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.AuthAgentCreateRequest)
}

type AgentAuthListParams struct {
	// Filter by domain
	Domain param.Opt[string] `query:"domain,omitzero" json:"-"`
	// Maximum number of results to return
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Number of results to skip
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Filter by profile name
	ProfileName param.Opt[string] `query:"profile_name,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AgentAuthListParams]'s query parameters as `url.Values`.
func (r AgentAuthListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
