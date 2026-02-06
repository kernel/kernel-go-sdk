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
	"github.com/kernel/kernel-go-sdk/packages/ssestream"
	"github.com/kernel/kernel-go-sdk/shared"
	"github.com/kernel/kernel-go-sdk/shared/constant"
)

// AuthConnectionService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAuthConnectionService] method instead.
type AuthConnectionService struct {
	Options []option.RequestOption
}

// NewAuthConnectionService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewAuthConnectionService(opts ...option.RequestOption) (r AuthConnectionService) {
	r = AuthConnectionService{}
	r.Options = opts
	return
}

// Creates managed authentication for a profile and domain combination. Returns 409
// Conflict if managed auth already exists for the given profile and domain.
func (r *AuthConnectionService) New(ctx context.Context, body AuthConnectionNewParams, opts ...option.RequestOption) (res *ManagedAuth, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "auth/connections"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieve managed auth by its ID. Includes current flow state if a login is in
// progress.
func (r *AuthConnectionService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *ManagedAuth, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("auth/connections/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List managed auths with optional filters for profile_name and domain.
func (r *AuthConnectionService) List(ctx context.Context, query AuthConnectionListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[ManagedAuth], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "auth/connections"
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

// List managed auths with optional filters for profile_name and domain.
func (r *AuthConnectionService) ListAutoPaging(ctx context.Context, query AuthConnectionListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[ManagedAuth] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Deletes managed auth and terminates its workflow. This will:
//
// - Delete the managed auth record
// - Terminate the Temporal workflow
// - Cancel any in-progress login flows
func (r *AuthConnectionService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("auth/connections/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return
}

// Establishes a Server-Sent Events (SSE) stream that delivers real-time login flow
// state updates. The stream terminates automatically once the flow reaches a
// terminal state (SUCCESS, FAILED, EXPIRED, CANCELED).
func (r *AuthConnectionService) FollowStreaming(ctx context.Context, id string, opts ...option.RequestOption) (stream *ssestream.Stream[AuthConnectionFollowResponseUnion]) {
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
	path := fmt.Sprintf("auth/connections/%s/events", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &raw, opts...)
	return ssestream.NewStream[AuthConnectionFollowResponseUnion](ssestream.NewDecoder(raw), err)
}

// Starts a login flow for the managed auth. Returns immediately with a hosted URL
// for the user to complete authentication, or triggers automatic re-auth if
// credentials are stored.
func (r *AuthConnectionService) Login(ctx context.Context, id string, body AuthConnectionLoginParams, opts ...option.RequestOption) (res *LoginResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("auth/connections/%s/login", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Submits field values for the login form. Poll the managed auth to track progress
// and get results.
func (r *AuthConnectionService) Submit(ctx context.Context, id string, body AuthConnectionSubmitParams, opts ...option.RequestOption) (res *SubmitFieldsResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("auth/connections/%s/submit", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Request to start a login flow
type LoginRequestParam struct {
	// If provided, saves credentials under this name upon successful login
	SaveCredentialAs param.Opt[string] `json:"save_credential_as,omitzero"`
	paramObj
}

func (r LoginRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow LoginRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *LoginRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Response from starting a login flow
type LoginResponse struct {
	// Managed auth ID
	ID string `json:"id,required"`
	// When the login flow expires
	FlowExpiresAt time.Time `json:"flow_expires_at,required" format:"date-time"`
	// Type of login flow started
	//
	// Any of "LOGIN", "REAUTH".
	FlowType LoginResponseFlowType `json:"flow_type,required"`
	// URL to redirect user to for login
	HostedURL string `json:"hosted_url,required" format:"uri"`
	// One-time code for handoff (internal use)
	HandoffCode string `json:"handoff_code"`
	// Browser live view URL for watching the login flow
	LiveViewURL string `json:"live_view_url" format:"uri"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		FlowExpiresAt respjson.Field
		FlowType      respjson.Field
		HostedURL     respjson.Field
		HandoffCode   respjson.Field
		LiveViewURL   respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LoginResponse) RawJSON() string { return r.JSON.raw }
func (r *LoginResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Type of login flow started
type LoginResponseFlowType string

const (
	LoginResponseFlowTypeLogin  LoginResponseFlowType = "LOGIN"
	LoginResponseFlowTypeReauth LoginResponseFlowType = "REAUTH"
)

// Managed authentication that keeps a profile logged into a specific domain. Flow
// fields (flow_status, flow_step, discovered_fields, mfa_options) reflect the most
// recent login flow and are null when no flow has been initiated.
type ManagedAuth struct {
	// Unique identifier for the managed auth
	ID string `json:"id,required"`
	// Target domain for authentication
	Domain string `json:"domain,required"`
	// Name of the profile associated with this managed auth
	ProfileName string `json:"profile_name,required"`
	// Current authentication status of the managed profile
	//
	// Any of "AUTHENTICATED", "NEEDS_AUTH".
	Status ManagedAuthStatus `json:"status,required"`
	// Additional domains that are valid for this auth flow (besides the primary
	// domain). Useful when login pages redirect to different domains.
	//
	// The following SSO/OAuth provider domains are automatically allowed by default
	// and do not need to be specified:
	//
	// - Google: accounts.google.com
	// - Microsoft/Azure AD: login.microsoftonline.com, login.live.com
	// - Okta: _.okta.com, _.oktapreview.com
	// - Auth0: _.auth0.com, _.us.auth0.com, _.eu.auth0.com, _.au.auth0.com
	// - Apple: appleid.apple.com
	// - GitHub: github.com
	// - Facebook/Meta: www.facebook.com
	// - LinkedIn: www.linkedin.com
	// - Amazon Cognito: \*.amazoncognito.com
	// - OneLogin: \*.onelogin.com
	// - Ping Identity: _.pingone.com, _.pingidentity.com
	AllowedDomains []string `json:"allowed_domains"`
	// Whether automatic re-authentication is possible (has credential, selectors, and
	// login_url)
	CanReauth bool `json:"can_reauth"`
	// Reference to credentials for managed auth. Use one of:
	//
	// - { name } for Kernel credentials
	// - { provider, path } for external provider item
	// - { provider, auto: true } for external provider domain lookup
	Credential ManagedAuthCredential `json:"credential"`
	// Fields awaiting input (present when flow_step=awaiting_input)
	DiscoveredFields []DiscoveredField `json:"discovered_fields,nullable"`
	// Error message (present when flow_status=failed)
	ErrorMessage string `json:"error_message,nullable"`
	// Instructions for external action (present when
	// flow_step=awaiting_external_action)
	ExternalActionMessage string `json:"external_action_message,nullable"`
	// When the current flow expires (null when no flow in progress)
	FlowExpiresAt time.Time `json:"flow_expires_at,nullable" format:"date-time"`
	// Current flow status (null when no flow in progress)
	//
	// Any of "IN_PROGRESS", "SUCCESS", "FAILED", "EXPIRED", "CANCELED".
	FlowStatus ManagedAuthFlowStatus `json:"flow_status,nullable"`
	// Current step in the flow (null when no flow in progress)
	//
	// Any of "DISCOVERING", "AWAITING_INPUT", "AWAITING_EXTERNAL_ACTION",
	// "SUBMITTING", "COMPLETED".
	FlowStep ManagedAuthFlowStep `json:"flow_step,nullable"`
	// Type of the current flow (null when no flow in progress)
	//
	// Any of "LOGIN", "REAUTH".
	FlowType ManagedAuthFlowType `json:"flow_type,nullable"`
	// Interval in seconds between automatic health checks. When set, the system
	// periodically verifies the authentication status and triggers re-authentication
	// if needed. Must be between 300 (5 minutes) and 86400 (24 hours). Default is 3600
	// (1 hour).
	HealthCheckInterval int64 `json:"health_check_interval,nullable"`
	// URL to redirect user to for hosted login (present when flow in progress)
	HostedURL string `json:"hosted_url,nullable" format:"uri"`
	// When the profile was last successfully authenticated
	LastAuthAt time.Time `json:"last_auth_at" format:"date-time"`
	// Browser live view URL for debugging (present when flow in progress)
	LiveViewURL string `json:"live_view_url,nullable" format:"uri"`
	// MFA method options (present when flow_step=awaiting_input and MFA selection
	// required)
	MfaOptions []ManagedAuthMfaOption `json:"mfa_options,nullable"`
	// SSO buttons available (present when flow_step=awaiting_input)
	PendingSSOButtons []ManagedAuthPendingSSOButton `json:"pending_sso_buttons,nullable"`
	// URL where the browser landed after successful login
	PostLoginURL string `json:"post_login_url" format:"uri"`
	// SSO provider being used (e.g., google, github, microsoft)
	SSOProvider string `json:"sso_provider,nullable"`
	// Visible error message from the website (e.g., 'Incorrect password'). Present
	// when the website displays an error during login.
	WebsiteError string `json:"website_error,nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                    respjson.Field
		Domain                respjson.Field
		ProfileName           respjson.Field
		Status                respjson.Field
		AllowedDomains        respjson.Field
		CanReauth             respjson.Field
		Credential            respjson.Field
		DiscoveredFields      respjson.Field
		ErrorMessage          respjson.Field
		ExternalActionMessage respjson.Field
		FlowExpiresAt         respjson.Field
		FlowStatus            respjson.Field
		FlowStep              respjson.Field
		FlowType              respjson.Field
		HealthCheckInterval   respjson.Field
		HostedURL             respjson.Field
		LastAuthAt            respjson.Field
		LiveViewURL           respjson.Field
		MfaOptions            respjson.Field
		PendingSSOButtons     respjson.Field
		PostLoginURL          respjson.Field
		SSOProvider           respjson.Field
		WebsiteError          respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuth) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuth) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current authentication status of the managed profile
type ManagedAuthStatus string

const (
	ManagedAuthStatusAuthenticated ManagedAuthStatus = "AUTHENTICATED"
	ManagedAuthStatusNeedsAuth     ManagedAuthStatus = "NEEDS_AUTH"
)

// Reference to credentials for managed auth. Use one of:
//
// - { name } for Kernel credentials
// - { provider, path } for external provider item
// - { provider, auto: true } for external provider domain lookup
type ManagedAuthCredential struct {
	// If true, lookup by domain from the specified provider
	Auto bool `json:"auto"`
	// Kernel credential name
	Name string `json:"name"`
	// Provider-specific path (e.g., "VaultName/ItemName" for 1Password)
	Path string `json:"path"`
	// External provider name (e.g., "my-1p")
	Provider string `json:"provider"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Auto        respjson.Field
		Name        respjson.Field
		Path        respjson.Field
		Provider    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthCredential) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthCredential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current flow status (null when no flow in progress)
type ManagedAuthFlowStatus string

const (
	ManagedAuthFlowStatusInProgress ManagedAuthFlowStatus = "IN_PROGRESS"
	ManagedAuthFlowStatusSuccess    ManagedAuthFlowStatus = "SUCCESS"
	ManagedAuthFlowStatusFailed     ManagedAuthFlowStatus = "FAILED"
	ManagedAuthFlowStatusExpired    ManagedAuthFlowStatus = "EXPIRED"
	ManagedAuthFlowStatusCanceled   ManagedAuthFlowStatus = "CANCELED"
)

// Current step in the flow (null when no flow in progress)
type ManagedAuthFlowStep string

const (
	ManagedAuthFlowStepDiscovering            ManagedAuthFlowStep = "DISCOVERING"
	ManagedAuthFlowStepAwaitingInput          ManagedAuthFlowStep = "AWAITING_INPUT"
	ManagedAuthFlowStepAwaitingExternalAction ManagedAuthFlowStep = "AWAITING_EXTERNAL_ACTION"
	ManagedAuthFlowStepSubmitting             ManagedAuthFlowStep = "SUBMITTING"
	ManagedAuthFlowStepCompleted              ManagedAuthFlowStep = "COMPLETED"
)

// Type of the current flow (null when no flow in progress)
type ManagedAuthFlowType string

const (
	ManagedAuthFlowTypeLogin  ManagedAuthFlowType = "LOGIN"
	ManagedAuthFlowTypeReauth ManagedAuthFlowType = "REAUTH"
)

// An MFA method option for verification
type ManagedAuthMfaOption struct {
	// The visible option text
	Label string `json:"label,required"`
	// The MFA delivery method type (includes password for auth method selection pages)
	//
	// Any of "sms", "call", "email", "totp", "push", "password".
	Type string `json:"type,required"`
	// Additional instructions from the site
	Description string `json:"description,nullable"`
	// The masked destination (phone/email) if shown
	Target string `json:"target,nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Label       respjson.Field
		Type        respjson.Field
		Description respjson.Field
		Target      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthMfaOption) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthMfaOption) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An SSO button for signing in with an external identity provider
type ManagedAuthPendingSSOButton struct {
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
func (r ManagedAuthPendingSSOButton) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthPendingSSOButton) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request to create managed auth for a profile and domain
//
// The properties Domain, ProfileName are required.
type ManagedAuthCreateRequestParam struct {
	// Domain for authentication
	Domain string `json:"domain,required"`
	// Name of the profile to manage authentication for
	ProfileName string `json:"profile_name,required"`
	// Interval in seconds between automatic health checks. When set, the system
	// periodically verifies the authentication status and triggers re-authentication
	// if needed. Must be between 300 (5 minutes) and 86400 (24 hours). Default is 3600
	// (1 hour).
	HealthCheckInterval param.Opt[int64] `json:"health_check_interval,omitzero"`
	// Optional login page URL to skip discovery
	LoginURL param.Opt[string] `json:"login_url,omitzero" format:"uri"`
	// Additional domains valid for this auth flow (besides the primary domain). Useful
	// when login pages redirect to different domains.
	//
	// The following SSO/OAuth provider domains are automatically allowed by default
	// and do not need to be specified:
	//
	// - Google: accounts.google.com
	// - Microsoft/Azure AD: login.microsoftonline.com, login.live.com
	// - Okta: _.okta.com, _.oktapreview.com
	// - Auth0: _.auth0.com, _.us.auth0.com, _.eu.auth0.com, _.au.auth0.com
	// - Apple: appleid.apple.com
	// - GitHub: github.com
	// - Facebook/Meta: www.facebook.com
	// - LinkedIn: www.linkedin.com
	// - Amazon Cognito: \*.amazoncognito.com
	// - OneLogin: \*.onelogin.com
	// - Ping Identity: _.pingone.com, _.pingidentity.com
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	// Reference to credentials for managed auth. Use one of:
	//
	// - { name } for Kernel credentials
	// - { provider, path } for external provider item
	// - { provider, auto: true } for external provider domain lookup
	Credential ManagedAuthCreateRequestCredentialParam `json:"credential,omitzero"`
	// Optional proxy configuration
	Proxy ManagedAuthCreateRequestProxyParam `json:"proxy,omitzero"`
	paramObj
}

func (r ManagedAuthCreateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthCreateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthCreateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Reference to credentials for managed auth. Use one of:
//
// - { name } for Kernel credentials
// - { provider, path } for external provider item
// - { provider, auto: true } for external provider domain lookup
type ManagedAuthCreateRequestCredentialParam struct {
	// If true, lookup by domain from the specified provider
	Auto param.Opt[bool] `json:"auto,omitzero"`
	// Kernel credential name
	Name param.Opt[string] `json:"name,omitzero"`
	// Provider-specific path (e.g., "VaultName/ItemName" for 1Password)
	Path param.Opt[string] `json:"path,omitzero"`
	// External provider name (e.g., "my-1p")
	Provider param.Opt[string] `json:"provider,omitzero"`
	paramObj
}

func (r ManagedAuthCreateRequestCredentialParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthCreateRequestCredentialParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthCreateRequestCredentialParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Optional proxy configuration
type ManagedAuthCreateRequestProxyParam struct {
	// ID of the proxy to use
	ProxyID param.Opt[string] `json:"proxy_id,omitzero"`
	paramObj
}

func (r ManagedAuthCreateRequestProxyParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthCreateRequestProxyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthCreateRequestProxyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request to submit field values for login
//
// The property Fields is required.
type SubmitFieldsRequestParam struct {
	// Map of field name to value
	Fields map[string]string `json:"fields,omitzero,required"`
	// Optional MFA option ID if user selected an MFA method
	MfaOptionID param.Opt[string] `json:"mfa_option_id,omitzero"`
	// Optional XPath selector if user chose to click an SSO button instead
	SSOButtonSelector param.Opt[string] `json:"sso_button_selector,omitzero"`
	paramObj
}

func (r SubmitFieldsRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow SubmitFieldsRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SubmitFieldsRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Response from submitting field values
type SubmitFieldsResponse struct {
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
func (r SubmitFieldsResponse) RawJSON() string { return r.JSON.raw }
func (r *SubmitFieldsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AuthConnectionFollowResponseUnion contains all possible properties and values
// from [AuthConnectionFollowResponseManagedAuthState], [shared.ErrorEvent],
// [shared.HeartbeatEvent].
//
// Use the [AuthConnectionFollowResponseUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AuthConnectionFollowResponseUnion struct {
	// Any of "managed_auth_state", "error", "sse_heartbeat".
	Event string `json:"event"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	FlowStatus string `json:"flow_status"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	FlowStep  string    `json:"flow_step"`
	Timestamp time.Time `json:"timestamp"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	DiscoveredFields []DiscoveredField `json:"discovered_fields"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	ErrorMessage string `json:"error_message"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	ExternalActionMessage string `json:"external_action_message"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	FlowType string `json:"flow_type"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	HostedURL string `json:"hosted_url"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	LiveViewURL string `json:"live_view_url"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	MfaOptions []AuthConnectionFollowResponseManagedAuthStateMfaOption `json:"mfa_options"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	PendingSSOButtons []AuthConnectionFollowResponseManagedAuthStatePendingSSOButton `json:"pending_sso_buttons"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	PostLoginURL string `json:"post_login_url"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	WebsiteError string `json:"website_error"`
	// This field is from variant [shared.ErrorEvent].
	Error shared.ErrorModel `json:"error"`
	JSON  struct {
		Event                 respjson.Field
		FlowStatus            respjson.Field
		FlowStep              respjson.Field
		Timestamp             respjson.Field
		DiscoveredFields      respjson.Field
		ErrorMessage          respjson.Field
		ExternalActionMessage respjson.Field
		FlowType              respjson.Field
		HostedURL             respjson.Field
		LiveViewURL           respjson.Field
		MfaOptions            respjson.Field
		PendingSSOButtons     respjson.Field
		PostLoginURL          respjson.Field
		WebsiteError          respjson.Field
		Error                 respjson.Field
		raw                   string
	} `json:"-"`
}

// anyAuthConnectionFollowResponse is implemented by each variant of
// [AuthConnectionFollowResponseUnion] to add type safety for the return type of
// [AuthConnectionFollowResponseUnion.AsAny]
type anyAuthConnectionFollowResponse interface {
	ImplAuthConnectionFollowResponseUnion()
}

func (AuthConnectionFollowResponseManagedAuthState) ImplAuthConnectionFollowResponseUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AuthConnectionFollowResponseUnion.AsAny().(type) {
//	case kernel.AuthConnectionFollowResponseManagedAuthState:
//	case shared.ErrorEvent:
//	case shared.HeartbeatEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AuthConnectionFollowResponseUnion) AsAny() anyAuthConnectionFollowResponse {
	switch u.Event {
	case "managed_auth_state":
		return u.AsManagedAuthState()
	case "error":
		return u.AsError()
	case "sse_heartbeat":
		return u.AsSseHeartbeat()
	}
	return nil
}

func (u AuthConnectionFollowResponseUnion) AsManagedAuthState() (v AuthConnectionFollowResponseManagedAuthState) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AuthConnectionFollowResponseUnion) AsError() (v shared.ErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AuthConnectionFollowResponseUnion) AsSseHeartbeat() (v shared.HeartbeatEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AuthConnectionFollowResponseUnion) RawJSON() string { return u.JSON.raw }

func (r *AuthConnectionFollowResponseUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An event representing the current state of a managed auth flow.
type AuthConnectionFollowResponseManagedAuthState struct {
	// Event type identifier (always "managed_auth_state").
	Event constant.ManagedAuthState `json:"event,required"`
	// Current flow status.
	//
	// Any of "IN_PROGRESS", "SUCCESS", "FAILED", "EXPIRED", "CANCELED".
	FlowStatus string `json:"flow_status,required"`
	// Current step in the flow.
	//
	// Any of "DISCOVERING", "AWAITING_INPUT", "AWAITING_EXTERNAL_ACTION",
	// "SUBMITTING", "COMPLETED".
	FlowStep string `json:"flow_step,required"`
	// Time the state was reported.
	Timestamp time.Time `json:"timestamp,required" format:"date-time"`
	// Fields awaiting input (present when flow_step=AWAITING_INPUT).
	DiscoveredFields []DiscoveredField `json:"discovered_fields"`
	// Error message (present when flow_status=FAILED).
	ErrorMessage string `json:"error_message"`
	// Instructions for external action (present when
	// flow_step=AWAITING_EXTERNAL_ACTION).
	ExternalActionMessage string `json:"external_action_message"`
	// Type of the current flow.
	//
	// Any of "LOGIN", "REAUTH".
	FlowType string `json:"flow_type"`
	// URL to redirect user to for hosted login.
	HostedURL string `json:"hosted_url" format:"uri"`
	// Browser live view URL for debugging.
	LiveViewURL string `json:"live_view_url" format:"uri"`
	// MFA method options (present when flow_step=AWAITING_INPUT and MFA selection
	// required).
	MfaOptions []AuthConnectionFollowResponseManagedAuthStateMfaOption `json:"mfa_options"`
	// SSO buttons available (present when flow_step=AWAITING_INPUT).
	PendingSSOButtons []AuthConnectionFollowResponseManagedAuthStatePendingSSOButton `json:"pending_sso_buttons"`
	// URL where the browser landed after successful login.
	PostLoginURL string `json:"post_login_url" format:"uri"`
	// Visible error message from the website (e.g., 'Incorrect password'). Present
	// when the website displays an error during login.
	WebsiteError string `json:"website_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Event                 respjson.Field
		FlowStatus            respjson.Field
		FlowStep              respjson.Field
		Timestamp             respjson.Field
		DiscoveredFields      respjson.Field
		ErrorMessage          respjson.Field
		ExternalActionMessage respjson.Field
		FlowType              respjson.Field
		HostedURL             respjson.Field
		LiveViewURL           respjson.Field
		MfaOptions            respjson.Field
		PendingSSOButtons     respjson.Field
		PostLoginURL          respjson.Field
		WebsiteError          respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthConnectionFollowResponseManagedAuthState) RawJSON() string { return r.JSON.raw }
func (r *AuthConnectionFollowResponseManagedAuthState) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An MFA method option for verification
type AuthConnectionFollowResponseManagedAuthStateMfaOption struct {
	// The visible option text
	Label string `json:"label,required"`
	// The MFA delivery method type (includes password for auth method selection pages)
	//
	// Any of "sms", "call", "email", "totp", "push", "password".
	Type string `json:"type,required"`
	// Additional instructions from the site
	Description string `json:"description,nullable"`
	// The masked destination (phone/email) if shown
	Target string `json:"target,nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Label       respjson.Field
		Type        respjson.Field
		Description respjson.Field
		Target      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthConnectionFollowResponseManagedAuthStateMfaOption) RawJSON() string { return r.JSON.raw }
func (r *AuthConnectionFollowResponseManagedAuthStateMfaOption) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An SSO button for signing in with an external identity provider
type AuthConnectionFollowResponseManagedAuthStatePendingSSOButton struct {
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
func (r AuthConnectionFollowResponseManagedAuthStatePendingSSOButton) RawJSON() string {
	return r.JSON.raw
}
func (r *AuthConnectionFollowResponseManagedAuthStatePendingSSOButton) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthConnectionNewParams struct {
	// Request to create managed auth for a profile and domain
	ManagedAuthCreateRequest ManagedAuthCreateRequestParam
	paramObj
}

func (r AuthConnectionNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.ManagedAuthCreateRequest)
}
func (r *AuthConnectionNewParams) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.ManagedAuthCreateRequest)
}

type AuthConnectionListParams struct {
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

// URLQuery serializes [AuthConnectionListParams]'s query parameters as
// `url.Values`.
func (r AuthConnectionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AuthConnectionLoginParams struct {
	// Request to start a login flow
	LoginRequest LoginRequestParam
	paramObj
}

func (r AuthConnectionLoginParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.LoginRequest)
}
func (r *AuthConnectionLoginParams) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.LoginRequest)
}

type AuthConnectionSubmitParams struct {
	// Request to submit field values for login
	SubmitFieldsRequest SubmitFieldsRequestParam
	paramObj
}

func (r AuthConnectionSubmitParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.SubmitFieldsRequest)
}
func (r *AuthConnectionSubmitParams) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.SubmitFieldsRequest)
}
