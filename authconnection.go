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

// Create and manage auth connections for automated credential capture and login.
//
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

// Creates an auth connection for a profile and domain combination. If the provided
// profile_name does not exist, it is created automatically. Returns 409 Conflict
// if an auth connection already exists for the given profile and domain.
func (r *AuthConnectionService) New(ctx context.Context, body AuthConnectionNewParams, opts ...option.RequestOption) (res *ManagedAuth, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "auth/connections"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve an auth connection by its ID. Includes current flow state if a login is
// in progress.
func (r *AuthConnectionService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *ManagedAuth, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("auth/connections/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update an auth connection's configuration. Only the fields provided will be
// updated.
func (r *AuthConnectionService) Update(ctx context.Context, id string, body AuthConnectionUpdateParams, opts ...option.RequestOption) (res *ManagedAuth, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("auth/connections/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List auth connections with optional filters for profile_name and domain.
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

// List auth connections with optional filters for profile_name and domain.
func (r *AuthConnectionService) ListAutoPaging(ctx context.Context, query AuthConnectionListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[ManagedAuth] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Deletes an auth connection and terminates its workflow. This will:
//
// - Delete the auth connection record
// - Terminate the Temporal workflow
// - Cancel any in-progress login flows
func (r *AuthConnectionService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("auth/connections/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
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
		return ssestream.NewStream[AuthConnectionFollowResponseUnion](nil, err)
	}
	path := fmt.Sprintf("auth/connections/%s/events", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &raw, opts...)
	return ssestream.NewStream[AuthConnectionFollowResponseUnion](ssestream.NewDecoder(raw), err)
}

// Starts a login flow for the auth connection. Returns immediately with a hosted
// URL for the user to complete authentication, or triggers automatic re-auth if
// credentials are stored.
func (r *AuthConnectionService) Login(ctx context.Context, id string, body AuthConnectionLoginParams, opts ...option.RequestOption) (res *LoginResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("auth/connections/%s/login", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Submits field values for the login form. Poll the auth connection to track
// progress and get results.
func (r *AuthConnectionService) Submit(ctx context.Context, id string, body AuthConnectionSubmitParams, opts ...option.RequestOption) (res *SubmitFieldsResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("auth/connections/%s/submit", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Response from starting a login flow
type LoginResponse struct {
	// Auth connection ID
	ID string `json:"id" api:"required"`
	// When the login flow expires
	FlowExpiresAt time.Time `json:"flow_expires_at" api:"required" format:"date-time"`
	// Type of login flow started
	//
	// Any of "LOGIN", "REAUTH".
	FlowType LoginResponseFlowType `json:"flow_type" api:"required"`
	// URL to redirect user to for login
	HostedURL string `json:"hosted_url" api:"required" format:"uri"`
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
	// Unique identifier for the auth connection
	ID string `json:"id" api:"required"`
	// Target domain for authentication
	Domain string `json:"domain" api:"required"`
	// Name of the profile associated with this auth connection
	ProfileName string `json:"profile_name" api:"required"`
	// Whether to record browser session replays for this connection by default. Useful
	// for debugging login flows. Can be overridden per-login.
	RecordSession bool `json:"record_session" api:"required"`
	// Whether credentials are saved after every successful login. One-time codes
	// (TOTP, SMS, etc.) are not saved.
	SaveCredentials bool `json:"save_credentials" api:"required"`
	// Current authentication status of the managed profile
	//
	// Any of "AUTHENTICATED", "NEEDS_AUTH".
	Status ManagedAuthStatus `json:"status" api:"required"`
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
	// Whether automatic re-authentication is permitted for this connection. This is an
	// opt-in flag only — it does not check whether re-auth is actually feasible. Even
	// when true, re-auth only runs when the system has what it needs to perform it
	// (for example, saved credentials for the required login fields), and only after a
	// scheduled health check detects an expired session — so this flag has no effect
	// when `health_checks` is false. When false, expired sessions detected by a health
	// check are marked as `NEEDS_AUTH` instead of attempting re-auth.
	AutoReauth bool `json:"auto_reauth"`
	// ID of the underlying browser session driving the current flow (present when flow
	// in progress). Use this to inspect or terminate the browser session via the
	// `/browsers` API.
	BrowserSessionID string `json:"browser_session_id" api:"nullable"`
	// Whether automatic re-authentication is possible (has credential, selectors, and
	// login_url)
	CanReauth bool `json:"can_reauth"`
	// Reason why automatic re-authentication is or is not possible
	CanReauthReason string `json:"can_reauth_reason"`
	// Reference to credentials for the auth connection. Use one of:
	//
	// - { name } for Kernel credentials
	// - { provider, path } for external provider item
	// - { provider, auto: true } for external provider domain lookup
	Credential ManagedAuthCredential `json:"credential"`
	// Fields awaiting input (present when flow_step=awaiting_input; may also be
	// present with awaiting_external_action as fallback actions)
	DiscoveredFields []ManagedAuthDiscoveredField `json:"discovered_fields" api:"nullable"`
	// Machine-readable error code (present when flow_status=failed)
	ErrorCode string `json:"error_code" api:"nullable"`
	// Error message (present when flow_status=failed)
	ErrorMessage string `json:"error_message" api:"nullable"`
	// Instructions for external action (present when
	// flow_step=awaiting_external_action)
	ExternalActionMessage string `json:"external_action_message" api:"nullable"`
	// When the current flow expires (null when no flow in progress). A flow past this
	// timestamp is no longer valid and its `flow_status` will be `EXPIRED`. Clients
	// may start a new login to supersede a stale `IN_PROGRESS` flow past this
	// timestamp.
	FlowExpiresAt time.Time `json:"flow_expires_at" api:"nullable" format:"date-time"`
	// Current flow status (null when no flow in progress)
	//
	// Any of "IN_PROGRESS", "SUCCESS", "FAILED", "EXPIRED", "CANCELED".
	FlowStatus ManagedAuthFlowStatus `json:"flow_status" api:"nullable"`
	// Current step in the flow (null when no flow in progress)
	//
	// Any of "DISCOVERING", "AWAITING_INPUT", "AWAITING_EXTERNAL_ACTION",
	// "SUBMITTING", "COMPLETED".
	FlowStep ManagedAuthFlowStep `json:"flow_step" api:"nullable"`
	// Type of the current flow (null when no flow in progress)
	//
	// Any of "LOGIN", "REAUTH".
	FlowType ManagedAuthFlowType `json:"flow_type" api:"nullable"`
	// Interval in seconds between automatic health checks. When set, the system
	// periodically verifies the authentication status and triggers re-authentication
	// if needed. Maximum is 86400 (24 hours). Default is 3600 (1 hour). The minimum
	// depends on your plan: Enterprise: 300 (5 minutes), Startup: 1200 (20 minutes),
	// Hobbyist: 3600 (1 hour).
	HealthCheckInterval int64 `json:"health_check_interval" api:"nullable"`
	// Whether periodic health checks are enabled for this connection. When false, the
	// system will not automatically verify authentication status, and `auto_reauth`
	// has no effect on the automatic flow (since re-auth is only triggered by a failed
	// scheduled health check). Manually triggering a health check via the API still
	// works regardless of this setting.
	HealthChecks bool `json:"health_checks"`
	// URL to redirect user to for hosted login (present when flow in progress)
	HostedURL string `json:"hosted_url" api:"nullable" format:"uri"`
	// Deprecated alias for `last_auth_check_at`. Despite the name, this is the last
	// health-check timestamp, not the last successful authentication. Use
	// `last_auth_check_at` instead.
	//
	// Deprecated: deprecated
	LastAuthAt time.Time `json:"last_auth_at" format:"date-time"`
	// When the most recent auth health check ran for this connection, regardless of
	// outcome. Updated on every health check and does not by itself indicate that the
	// profile is currently authenticated - use `status` for that. May be newer than
	// `flow_expires_at` when a flow is still in progress because health checks
	// continue to run in parallel.
	LastAuthCheckAt time.Time `json:"last_auth_check_at" format:"date-time"`
	// Browser live view URL for debugging (present when flow in progress)
	LiveViewURL string `json:"live_view_url" api:"nullable" format:"uri"`
	// Optional login page URL to skip discovery
	LoginURL string `json:"login_url" format:"uri"`
	// MFA method options (present when flow_step=awaiting_input; may also be present
	// with awaiting_external_action as fallback actions)
	MfaOptions []ManagedAuthMfaOption `json:"mfa_options" api:"nullable"`
	// SSO buttons available (present when flow_step=awaiting_input; may also be
	// present with awaiting_external_action as fallback actions)
	PendingSSOButtons []ManagedAuthPendingSSOButton `json:"pending_sso_buttons" api:"nullable"`
	// URL where the browser landed after successful login
	PostLoginURL string `json:"post_login_url" format:"uri"`
	// ID of the proxy associated with this connection, if any.
	ProxyID string `json:"proxy_id"`
	// Non-MFA choices presented during the auth flow, such as account selection or org
	// pickers (present when flow_step=awaiting_input; may also be present with
	// awaiting_external_action as fallback actions).
	SignInOptions []ManagedAuthSignInOption `json:"sign_in_options" api:"nullable"`
	// SSO provider being used (e.g., google, github, microsoft)
	SSOProvider string `json:"sso_provider" api:"nullable"`
	// Visible error message from the website (e.g., 'Incorrect password'). Present
	// when the website displays an error during login.
	WebsiteError string `json:"website_error" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                    respjson.Field
		Domain                respjson.Field
		ProfileName           respjson.Field
		RecordSession         respjson.Field
		SaveCredentials       respjson.Field
		Status                respjson.Field
		AllowedDomains        respjson.Field
		AutoReauth            respjson.Field
		BrowserSessionID      respjson.Field
		CanReauth             respjson.Field
		CanReauthReason       respjson.Field
		Credential            respjson.Field
		DiscoveredFields      respjson.Field
		ErrorCode             respjson.Field
		ErrorMessage          respjson.Field
		ExternalActionMessage respjson.Field
		FlowExpiresAt         respjson.Field
		FlowStatus            respjson.Field
		FlowStep              respjson.Field
		FlowType              respjson.Field
		HealthCheckInterval   respjson.Field
		HealthChecks          respjson.Field
		HostedURL             respjson.Field
		LastAuthAt            respjson.Field
		LastAuthCheckAt       respjson.Field
		LiveViewURL           respjson.Field
		LoginURL              respjson.Field
		MfaOptions            respjson.Field
		PendingSSOButtons     respjson.Field
		PostLoginURL          respjson.Field
		ProxyID               respjson.Field
		SignInOptions         respjson.Field
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

// Reference to credentials for the auth connection. Use one of:
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

// A discovered form field
type ManagedAuthDiscoveredField struct {
	// Field label
	Label string `json:"label" api:"required"`
	// Field name
	Name string `json:"name" api:"required"`
	// CSS selector for the field
	Selector string `json:"selector" api:"required"`
	// Field type
	//
	// Any of "text", "email", "password", "tel", "number", "url", "code", "totp".
	Type string `json:"type" api:"required"`
	// Contextual help text near the field that tells the user what to enter (e.g.,
	// "Enter the phone ending in (**_) _**-\*\*92")
	Hint string `json:"hint"`
	// If this field is associated with an MFA option, the type of that option (e.g.,
	// password field linked to "Enter password" option)
	//
	// Any of "sms", "call", "email", "totp", "push", "password", "switch".
	LinkedMfaType string `json:"linked_mfa_type" api:"nullable"`
	// Field placeholder
	Placeholder string `json:"placeholder"`
	// Whether field is required
	Required bool `json:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Label         respjson.Field
		Name          respjson.Field
		Selector      respjson.Field
		Type          respjson.Field
		Hint          respjson.Field
		LinkedMfaType respjson.Field
		Placeholder   respjson.Field
		Required      respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthDiscoveredField) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthDiscoveredField) UnmarshalJSON(data []byte) error {
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
	Label string `json:"label" api:"required"`
	// The MFA delivery method type. Includes 'password' for auth method selection
	// pages and 'switch' for generic method-switcher links like "Use another method"
	// that do not name a specific method.
	//
	// Any of "sms", "call", "email", "totp", "push", "password", "switch".
	Type string `json:"type" api:"required"`
	// Additional instructions from the site
	Description string `json:"description" api:"nullable"`
	// The masked destination (phone/email) if shown
	Target string `json:"target" api:"nullable"`
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
	Label string `json:"label" api:"required"`
	// Identity provider name
	Provider string `json:"provider" api:"required"`
	// XPath selector for the button
	Selector string `json:"selector" api:"required"`
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

// A non-MFA choice presented during the auth flow (e.g. account selection, org
// picker)
type ManagedAuthSignInOption struct {
	// Unique identifier for this option (used to submit selection back)
	ID string `json:"id" api:"required"`
	// Display text for the option
	Label string `json:"label" api:"required"`
	// Additional context such as email address or org name
	Description string `json:"description" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Label       respjson.Field
		Description respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthSignInOption) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthSignInOption) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request to create an auth connection for a profile and domain
//
// The properties Domain, ProfileName are required.
type ManagedAuthCreateRequestParam struct {
	// Domain for authentication
	Domain string `json:"domain" api:"required"`
	// Name of the profile to manage authentication for. If the profile does not exist,
	// it is created automatically.
	ProfileName string `json:"profile_name" api:"required"`
	// Whether to permit automatic re-authentication when a scheduled health check
	// detects an expired session. This is an opt-in flag only — it does not check
	// whether re-auth is actually feasible. Even when true, re-auth only runs when the
	// system has what it needs to perform it (for example, saved credentials for the
	// required login fields), and only after a scheduled health check detects an
	// expired session — so this flag has no effect when `health_checks` is false. When
	// false, expired sessions are marked as `NEEDS_AUTH` instead of attempting
	// re-auth. Defaults to true.
	AutoReauth param.Opt[bool] `json:"auto_reauth,omitzero"`
	// Interval in seconds between automatic health checks. When set, the system
	// periodically verifies the authentication status and triggers re-authentication
	// if needed. Maximum is 86400 (24 hours). Default is 3600 (1 hour). The minimum
	// depends on your plan: Enterprise: 300 (5 minutes), Startup: 1200 (20 minutes),
	// Hobbyist: 3600 (1 hour).
	HealthCheckInterval param.Opt[int64] `json:"health_check_interval,omitzero"`
	// Whether to enable periodic health checks. When false, the system will not
	// automatically verify authentication status, and `auto_reauth` has no effect on
	// the automatic flow (since re-auth is only triggered by a failed scheduled health
	// check). Defaults to true.
	HealthChecks param.Opt[bool] `json:"health_checks,omitzero"`
	// Optional login page URL to skip discovery
	LoginURL param.Opt[string] `json:"login_url,omitzero" format:"uri"`
	// Whether to record browser sessions for this connection by default. Useful for
	// debugging. Can be overridden per-login. Defaults to false.
	RecordSession param.Opt[bool] `json:"record_session,omitzero"`
	// Whether to save credentials after every successful login. Defaults to true.
	// One-time codes (TOTP, SMS, etc.) are not saved.
	SaveCredentials param.Opt[bool] `json:"save_credentials,omitzero"`
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
	// Reference to credentials for the auth connection. Use one of:
	//
	// - { name } for Kernel credentials
	// - { provider, path } for external provider item
	// - { provider, auto: true } for external provider domain lookup
	Credential ManagedAuthCreateRequestCredentialParam `json:"credential,omitzero"`
	// Proxy selection. Provide either id or name. The proxy must belong to the
	// caller's org.
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

// Reference to credentials for the auth connection. Use one of:
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

// Proxy selection. Provide either id or name. The proxy must belong to the
// caller's org.
type ManagedAuthCreateRequestProxyParam struct {
	// Proxy ID
	ID param.Opt[string] `json:"id,omitzero"`
	// Proxy name
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r ManagedAuthCreateRequestProxyParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthCreateRequestProxyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthCreateRequestProxyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request to update an auth connection's configuration
type ManagedAuthUpdateRequestParam struct {
	// Whether automatic re-authentication is permitted for this connection. This is an
	// opt-in flag only — it does not check whether re-auth is actually feasible. Even
	// when true, re-auth only runs when the system has what it needs to perform it
	// (for example, saved credentials for the required login fields), and only after a
	// scheduled health check detects an expired session — so this flag has no effect
	// when `health_checks` is false. When false, expired sessions detected by a health
	// check are marked as `NEEDS_AUTH` instead of attempting re-auth.
	AutoReauth param.Opt[bool] `json:"auto_reauth,omitzero"`
	// Interval in seconds between automatic health checks
	HealthCheckInterval param.Opt[int64] `json:"health_check_interval,omitzero"`
	// Whether periodic health checks are enabled. When set to false, the system will
	// not automatically verify authentication status, and `auto_reauth` has no effect
	// on the automatic flow (since re-auth is only triggered by a failed scheduled
	// health check).
	HealthChecks param.Opt[bool] `json:"health_checks,omitzero"`
	// Login page URL. Set to empty string to clear.
	LoginURL param.Opt[string] `json:"login_url,omitzero" format:"uri"`
	// Whether to record browser sessions for this connection by default
	RecordSession param.Opt[bool] `json:"record_session,omitzero"`
	// Whether to save credentials after every successful login
	SaveCredentials param.Opt[bool] `json:"save_credentials,omitzero"`
	// Additional domains valid for this auth flow (replaces existing list)
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	// Reference to credentials for the auth connection. Use one of:
	//
	// - { name } for Kernel credentials
	// - { provider, path } for external provider item
	// - { provider, auto: true } for external provider domain lookup
	Credential ManagedAuthUpdateRequestCredentialParam `json:"credential,omitzero"`
	// Proxy selection. Provide either id or name. The proxy must belong to the
	// caller's org.
	Proxy ManagedAuthUpdateRequestProxyParam `json:"proxy,omitzero"`
	paramObj
}

func (r ManagedAuthUpdateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthUpdateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthUpdateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Reference to credentials for the auth connection. Use one of:
//
// - { name } for Kernel credentials
// - { provider, path } for external provider item
// - { provider, auto: true } for external provider domain lookup
type ManagedAuthUpdateRequestCredentialParam struct {
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

func (r ManagedAuthUpdateRequestCredentialParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthUpdateRequestCredentialParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthUpdateRequestCredentialParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Proxy selection. Provide either id or name. The proxy must belong to the
// caller's org.
type ManagedAuthUpdateRequestProxyParam struct {
	// Proxy ID
	ID param.Opt[string] `json:"id,omitzero"`
	// Proxy name
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r ManagedAuthUpdateRequestProxyParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthUpdateRequestProxyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthUpdateRequestProxyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request to submit field values, click an SSO button, select an MFA method, or
// select a sign-in option. Provide exactly one of fields, sso_button_selector,
// sso_provider, mfa_option_id, or sign_in_option_id.
type SubmitFieldsRequestParam struct {
	// The MFA method type to select (when mfa_options were returned)
	MfaOptionID param.Opt[string] `json:"mfa_option_id,omitzero"`
	// The sign-in option ID to select (when sign_in_options were returned)
	SignInOptionID param.Opt[string] `json:"sign_in_option_id,omitzero"`
	// XPath selector for the SSO button to click (ODA). Use sso_provider instead for
	// CUA.
	SSOButtonSelector param.Opt[string] `json:"sso_button_selector,omitzero"`
	// SSO provider to click, matching the provider field from pending_sso_buttons
	// (e.g., "google", "github"). Cannot be used with sso_button_selector.
	SSOProvider param.Opt[string] `json:"sso_provider,omitzero"`
	// Map of field name to value
	Fields map[string]string `json:"fields,omitzero"`
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
	Accepted bool `json:"accepted" api:"required"`
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
	DiscoveredFields []AuthConnectionFollowResponseManagedAuthStateDiscoveredField `json:"discovered_fields"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	ErrorCode string `json:"error_code"`
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
	SignInOptions []AuthConnectionFollowResponseManagedAuthStateSignInOption `json:"sign_in_options"`
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
		ErrorCode             respjson.Field
		ErrorMessage          respjson.Field
		ExternalActionMessage respjson.Field
		FlowType              respjson.Field
		HostedURL             respjson.Field
		LiveViewURL           respjson.Field
		MfaOptions            respjson.Field
		PendingSSOButtons     respjson.Field
		PostLoginURL          respjson.Field
		SignInOptions         respjson.Field
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
	Event constant.ManagedAuthState `json:"event" default:"managed_auth_state"`
	// Current flow status.
	//
	// Any of "IN_PROGRESS", "SUCCESS", "FAILED", "EXPIRED", "CANCELED".
	FlowStatus string `json:"flow_status" api:"required"`
	// Current step in the flow.
	//
	// Any of "DISCOVERING", "AWAITING_INPUT", "AWAITING_EXTERNAL_ACTION",
	// "SUBMITTING", "COMPLETED".
	FlowStep string `json:"flow_step" api:"required"`
	// Time the state was reported.
	Timestamp time.Time `json:"timestamp" api:"required" format:"date-time"`
	// Fields awaiting input (present when flow_step=AWAITING_INPUT; may also be
	// present with AWAITING_EXTERNAL_ACTION as fallback actions).
	DiscoveredFields []AuthConnectionFollowResponseManagedAuthStateDiscoveredField `json:"discovered_fields"`
	// Machine-readable error code (present when flow_status=FAILED).
	ErrorCode string `json:"error_code"`
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
	// MFA method options (present when flow_step=AWAITING_INPUT; may also be present
	// with AWAITING_EXTERNAL_ACTION as fallback actions).
	MfaOptions []AuthConnectionFollowResponseManagedAuthStateMfaOption `json:"mfa_options"`
	// SSO buttons available (present when flow_step=AWAITING_INPUT; may also be
	// present with AWAITING_EXTERNAL_ACTION as fallback actions).
	PendingSSOButtons []AuthConnectionFollowResponseManagedAuthStatePendingSSOButton `json:"pending_sso_buttons"`
	// URL where the browser landed after successful login.
	PostLoginURL string `json:"post_login_url" format:"uri"`
	// Non-MFA choices presented during the auth flow, such as account selection or org
	// pickers (present when flow_step=AWAITING_INPUT; may also be present with
	// AWAITING_EXTERNAL_ACTION as fallback actions).
	SignInOptions []AuthConnectionFollowResponseManagedAuthStateSignInOption `json:"sign_in_options"`
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
		ErrorCode             respjson.Field
		ErrorMessage          respjson.Field
		ExternalActionMessage respjson.Field
		FlowType              respjson.Field
		HostedURL             respjson.Field
		LiveViewURL           respjson.Field
		MfaOptions            respjson.Field
		PendingSSOButtons     respjson.Field
		PostLoginURL          respjson.Field
		SignInOptions         respjson.Field
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

// A discovered form field
type AuthConnectionFollowResponseManagedAuthStateDiscoveredField struct {
	// Field label
	Label string `json:"label" api:"required"`
	// Field name
	Name string `json:"name" api:"required"`
	// CSS selector for the field
	Selector string `json:"selector" api:"required"`
	// Field type
	//
	// Any of "text", "email", "password", "tel", "number", "url", "code", "totp".
	Type string `json:"type" api:"required"`
	// Contextual help text near the field that tells the user what to enter (e.g.,
	// "Enter the phone ending in (**_) _**-\*\*92")
	Hint string `json:"hint"`
	// If this field is associated with an MFA option, the type of that option (e.g.,
	// password field linked to "Enter password" option)
	//
	// Any of "sms", "call", "email", "totp", "push", "password", "switch".
	LinkedMfaType string `json:"linked_mfa_type" api:"nullable"`
	// Field placeholder
	Placeholder string `json:"placeholder"`
	// Whether field is required
	Required bool `json:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Label         respjson.Field
		Name          respjson.Field
		Selector      respjson.Field
		Type          respjson.Field
		Hint          respjson.Field
		LinkedMfaType respjson.Field
		Placeholder   respjson.Field
		Required      respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthConnectionFollowResponseManagedAuthStateDiscoveredField) RawJSON() string {
	return r.JSON.raw
}
func (r *AuthConnectionFollowResponseManagedAuthStateDiscoveredField) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An MFA method option for verification
type AuthConnectionFollowResponseManagedAuthStateMfaOption struct {
	// The visible option text
	Label string `json:"label" api:"required"`
	// The MFA delivery method type. Includes 'password' for auth method selection
	// pages and 'switch' for generic method-switcher links like "Use another method"
	// that do not name a specific method.
	//
	// Any of "sms", "call", "email", "totp", "push", "password", "switch".
	Type string `json:"type" api:"required"`
	// Additional instructions from the site
	Description string `json:"description" api:"nullable"`
	// The masked destination (phone/email) if shown
	Target string `json:"target" api:"nullable"`
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
	Label string `json:"label" api:"required"`
	// Identity provider name
	Provider string `json:"provider" api:"required"`
	// XPath selector for the button
	Selector string `json:"selector" api:"required"`
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

// A non-MFA choice presented during the auth flow (e.g. account selection, org
// picker)
type AuthConnectionFollowResponseManagedAuthStateSignInOption struct {
	// Unique identifier for this option (used to submit selection back)
	ID string `json:"id" api:"required"`
	// Display text for the option
	Label string `json:"label" api:"required"`
	// Additional context such as email address or org name
	Description string `json:"description" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Label       respjson.Field
		Description respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthConnectionFollowResponseManagedAuthStateSignInOption) RawJSON() string { return r.JSON.raw }
func (r *AuthConnectionFollowResponseManagedAuthStateSignInOption) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthConnectionNewParams struct {
	// Request to create an auth connection for a profile and domain
	ManagedAuthCreateRequest ManagedAuthCreateRequestParam
	paramObj
}

func (r AuthConnectionNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.ManagedAuthCreateRequest)
}
func (r *AuthConnectionNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthConnectionUpdateParams struct {
	// Request to update an auth connection's configuration
	ManagedAuthUpdateRequest ManagedAuthUpdateRequestParam
	paramObj
}

func (r AuthConnectionUpdateParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.ManagedAuthUpdateRequest)
}
func (r *AuthConnectionUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
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
	// Override the connection's default for recording this login's browser session.
	// When omitted, the connection's record_session default is used.
	RecordSession param.Opt[bool] `json:"record_session,omitzero"`
	// Proxy selection. Provide either id or name. The proxy must belong to the
	// caller's org.
	Proxy AuthConnectionLoginParamsProxy `json:"proxy,omitzero"`
	paramObj
}

func (r AuthConnectionLoginParams) MarshalJSON() (data []byte, err error) {
	type shadow AuthConnectionLoginParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AuthConnectionLoginParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Proxy selection. Provide either id or name. The proxy must belong to the
// caller's org.
type AuthConnectionLoginParamsProxy struct {
	// Proxy ID
	ID param.Opt[string] `json:"id,omitzero"`
	// Proxy name
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r AuthConnectionLoginParamsProxy) MarshalJSON() (data []byte, err error) {
	type shadow AuthConnectionLoginParamsProxy
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AuthConnectionLoginParamsProxy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthConnectionSubmitParams struct {
	// Request to submit field values, click an SSO button, select an MFA method, or
	// select a sign-in option. Provide exactly one of fields, sso_button_selector,
	// sso_provider, mfa_option_id, or sign_in_option_id.
	SubmitFieldsRequest SubmitFieldsRequestParam
	paramObj
}

func (r AuthConnectionSubmitParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.SubmitFieldsRequest)
}
func (r *AuthConnectionSubmitParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
