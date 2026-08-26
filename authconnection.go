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

// Returns a chronological timeline of events for an auth connection — login
// attempts, automatic re-auth attempts, and health checks. Events are returned
// newest-first.
func (r *AuthConnectionService) Timeline(ctx context.Context, id string, query AuthConnectionTimelineParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[ManagedAuthTimelineEvent], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("auth/connections/%s/timeline", id)
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

// Returns a chronological timeline of events for an auth connection — login
// attempts, automatic re-auth attempts, and health checks. Events are returned
// newest-first.
func (r *AuthConnectionService) TimelineAutoPaging(ctx context.Context, id string, query AuthConnectionTimelineParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[ManagedAuthTimelineEvent] {
	return pagination.NewOffsetPaginationAutoPager(r.Timeline(ctx, id, query, opts...))
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
	// Additional hostname roots valid for this auth flow, besides the primary domain.
	// Each value allows credential entry on that exact hostname and its subdomains.
	// Leading `www.` and `*.` labels are normalized away. When omitted or empty,
	// credential entry is unrestricted.
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
	// Default browser configuration for login, reauthentication, and health-check
	// sessions.
	Browser ManagedAuthBrowserConfig `json:"browser"`
	// ID of the underlying browser session driving the current flow (present when flow
	// in progress). Use this to inspect or terminate the browser session via the
	// `/browsers` API.
	BrowserSessionID string `json:"browser_session_id" api:"nullable"`
	// Deprecated. Use browser.telemetry. Retained during migration for existing
	// clients.
	//
	// Deprecated: deprecated
	BrowserTelemetry ManagedAuthBrowserTelemetry `json:"browser_telemetry" api:"nullable"`
	// Whether Kernel can automatically re-authenticate this connection when the
	// session expires. Requires a prior successful login plus either a Kernel
	// credential or an external credential reference. See `can_reauth_reason` for the
	// specific outcome.
	CanReauth bool `json:"can_reauth"`
	// Machine-readable reason for the current value of `can_reauth`. Affirmative
	// values (re-auth is possible):
	//
	// - `external_credential` — an external credential provider is attached
	// - `cua_has_credential` — CUA flow with a stored credential
	// - `has_credential` — Kernel credential is attached (optimistic; plan viability
	//   not checked)
	// - `viable_plans_found` — at least one stored login plan can be replayed
	// - `no_requirements_recorded` — no recorded credential requirements to fail
	//   against
	// - `requirements_satisfiable` — recorded requirements can be met by the attached
	//   credential
	//
	// Negative values (a human must complete the login flow):
	//
	// - `no_prior_successful_login` — connection has never completed a successful
	//   login
	// - `no_credential` — no Kernel or external credential attached
	// - `no_viable_plans` — credential attached but no replayable login plan exists
	//   yet
	// - `viable_plans_require_external_action` — stored plans need an external step
	//   (email link, push, etc.)
	// - `requires_external_action` — recorded requirements include an external step
	// - `requires_totp_without_secret` — flow needs a TOTP code but no TOTP secret is
	//   stored
	// - `requires_sms_code` — flow needs an SMS code that cannot be received
	//   automatically
	// - `requires_email_code` — flow needs an email code that cannot be received
	//   automatically
	// - `requires_customer_input` — flow needs another field or choice that is
	//   unavailable during unattended re-authentication
	//
	// Any of "external_credential", "cua_has_credential", "has_credential",
	// "viable_plans_found", "no_requirements_recorded", "requirements_satisfiable",
	// "no_prior_successful_login", "no_credential", "no_viable_plans",
	// "viable_plans_require_external_action", "requires_external_action",
	// "requires_totp_without_secret", "requires_sms_code", "requires_email_code",
	// "requires_customer_input".
	CanReauthReason ManagedAuthCanReauthReason `json:"can_reauth_reason"`
	// Canonical choices awaiting selection. Prefer this over pending_sso_buttons,
	// mfa_options, and sign_in_options when present.
	Choices []ManagedAuthChoice `json:"choices" api:"nullable"`
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
	// Canonical fields awaiting input. Prefer this over discovered_fields when
	// present.
	Fields []ManagedAuthField `json:"fields" api:"nullable"`
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
	// if needed. Maximum is 86400 (24 hours). Default is 3600 (1 hour) or your plan
	// minimum, whichever is larger. The minimum depends on your plan: Enterprise: 300
	// (5 minutes), Startup: 1200 (20 minutes), Hobbyist: 3600 (1 hour), Free: 21600 (6
	// hours).
	HealthCheckInterval int64 `json:"health_check_interval" api:"nullable"`
	// Whether periodic health checks are enabled for this connection. When false, the
	// system will not automatically verify authentication status, and `auto_reauth`
	// has no effect on the automatic flow (since re-auth is only triggered by a failed
	// scheduled health check). Manually triggering a health check via the API still
	// works regardless of this setting.
	HealthChecks bool `json:"health_checks"`
	// URL to redirect user to for hosted login (present when flow in progress)
	HostedURL string `json:"hosted_url" api:"nullable" format:"uri"`
	// Opaque identifier for the current canonical interaction. Required when
	// submitting fields or choices and changes for each new actionable pause.
	InteractionID string `json:"interaction_id" api:"nullable"`
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
	// Deprecated. Read browser.proxy instead. Retained during migration for existing
	// clients.
	//
	// Deprecated: deprecated
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
		Browser               respjson.Field
		BrowserSessionID      respjson.Field
		BrowserTelemetry      respjson.Field
		CanReauth             respjson.Field
		CanReauthReason       respjson.Field
		Choices               respjson.Field
		Credential            respjson.Field
		DiscoveredFields      respjson.Field
		ErrorCode             respjson.Field
		ErrorMessage          respjson.Field
		ExternalActionMessage respjson.Field
		Fields                respjson.Field
		FlowExpiresAt         respjson.Field
		FlowStatus            respjson.Field
		FlowStep              respjson.Field
		FlowType              respjson.Field
		HealthCheckInterval   respjson.Field
		HealthChecks          respjson.Field
		HostedURL             respjson.Field
		InteractionID         respjson.Field
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

// Deprecated. Use browser.telemetry. Retained during migration for existing
// clients.
//
// Deprecated: deprecated
type ManagedAuthBrowserTelemetry struct {
	// Per-category capture flags. The operational categories (control, connection,
	// system, captcha) are captured whenever telemetry is enabled; set one to
	// enabled=false to opt out. The CDP categories (console, network, page,
	// interaction), screenshot and platform are off by default; set enabled=true to
	// opt in. On create, provided categories layer onto the default set. On update,
	// provided categories merge onto the session's current config; when no telemetry
	// is active this falls back to the default set (matching create). If browser is
	// omitted or empty, the default set is used. A browser config that disables every
	// category stops capture on update and starts no capture on create.
	Browser BrowserTelemetryCategoriesConfig `json:"browser"`
	// Request shortcut for browser telemetry capture. True enables capture; with no
	// browser category settings it captures the default set (control, connection,
	// system, captcha), and any browser category settings are layered onto that
	// default set. On update, enabled=true resolves the config fresh from the default
	// set plus any provided categories, replacing the session's current selection
	// rather than merging onto it; omit enabled to merge categories onto the current
	// selection instead. False stops capture on update and starts no capture on
	// create. enabled=false cannot be combined with browser category settings.
	Enabled bool `json:"enabled"`
	// Where to export this session's captured telemetry. Omit to capture without
	// exporting.
	Export ManagedAuthBrowserTelemetryExport `json:"export"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Browser     respjson.Field
		Enabled     respjson.Field
		Export      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthBrowserTelemetry) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthBrowserTelemetry) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where to export this session's captured telemetry. Omit to capture without
// exporting.
type ManagedAuthBrowserTelemetryExport struct {
	// Export captured telemetry over OTLP to one of the org's configured destinations.
	Otlp ManagedAuthBrowserTelemetryExportOtlp `json:"otlp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Otlp        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthBrowserTelemetryExport) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthBrowserTelemetryExport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Export captured telemetry over OTLP to one of the org's configured destinations.
type ManagedAuthBrowserTelemetryExportOtlp struct {
	// OTLP destination to export this session's captured telemetry to. Provide either
	// id or name. Requires telemetry capture to be enabled.
	Destination ManagedAuthBrowserTelemetryExportOtlpDestination `json:"destination"`
	// Whether to export captured telemetry over OTLP. Setting destination implies
	// enabled=true, so this only needs to be set explicitly to disable export
	// (enabled=false with a destination is rejected).
	Enabled bool `json:"enabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Destination respjson.Field
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthBrowserTelemetryExportOtlp) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthBrowserTelemetryExportOtlp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OTLP destination to export this session's captured telemetry to. Provide either
// id or name. Requires telemetry capture to be enabled.
type ManagedAuthBrowserTelemetryExportOtlpDestination struct {
	// OTLP destination ID
	ID string `json:"id"`
	// OTLP destination name
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthBrowserTelemetryExportOtlpDestination) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthBrowserTelemetryExportOtlpDestination) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Machine-readable reason for the current value of `can_reauth`. Affirmative
// values (re-auth is possible):
//
//   - `external_credential` — an external credential provider is attached
//   - `cua_has_credential` — CUA flow with a stored credential
//   - `has_credential` — Kernel credential is attached (optimistic; plan viability
//     not checked)
//   - `viable_plans_found` — at least one stored login plan can be replayed
//   - `no_requirements_recorded` — no recorded credential requirements to fail
//     against
//   - `requirements_satisfiable` — recorded requirements can be met by the attached
//     credential
//
// Negative values (a human must complete the login flow):
//
//   - `no_prior_successful_login` — connection has never completed a successful
//     login
//   - `no_credential` — no Kernel or external credential attached
//   - `no_viable_plans` — credential attached but no replayable login plan exists
//     yet
//   - `viable_plans_require_external_action` — stored plans need an external step
//     (email link, push, etc.)
//   - `requires_external_action` — recorded requirements include an external step
//   - `requires_totp_without_secret` — flow needs a TOTP code but no TOTP secret is
//     stored
//   - `requires_sms_code` — flow needs an SMS code that cannot be received
//     automatically
//   - `requires_email_code` — flow needs an email code that cannot be received
//     automatically
//   - `requires_customer_input` — flow needs another field or choice that is
//     unavailable during unattended re-authentication
type ManagedAuthCanReauthReason string

const (
	ManagedAuthCanReauthReasonExternalCredential               ManagedAuthCanReauthReason = "external_credential"
	ManagedAuthCanReauthReasonCuaHasCredential                 ManagedAuthCanReauthReason = "cua_has_credential"
	ManagedAuthCanReauthReasonHasCredential                    ManagedAuthCanReauthReason = "has_credential"
	ManagedAuthCanReauthReasonViablePlansFound                 ManagedAuthCanReauthReason = "viable_plans_found"
	ManagedAuthCanReauthReasonNoRequirementsRecorded           ManagedAuthCanReauthReason = "no_requirements_recorded"
	ManagedAuthCanReauthReasonRequirementsSatisfiable          ManagedAuthCanReauthReason = "requirements_satisfiable"
	ManagedAuthCanReauthReasonNoPriorSuccessfulLogin           ManagedAuthCanReauthReason = "no_prior_successful_login"
	ManagedAuthCanReauthReasonNoCredential                     ManagedAuthCanReauthReason = "no_credential"
	ManagedAuthCanReauthReasonNoViablePlans                    ManagedAuthCanReauthReason = "no_viable_plans"
	ManagedAuthCanReauthReasonViablePlansRequireExternalAction ManagedAuthCanReauthReason = "viable_plans_require_external_action"
	ManagedAuthCanReauthReasonRequiresExternalAction           ManagedAuthCanReauthReason = "requires_external_action"
	ManagedAuthCanReauthReasonRequiresTotpWithoutSecret        ManagedAuthCanReauthReason = "requires_totp_without_secret"
	ManagedAuthCanReauthReasonRequiresSMSCode                  ManagedAuthCanReauthReason = "requires_sms_code"
	ManagedAuthCanReauthReasonRequiresEmailCode                ManagedAuthCanReauthReason = "requires_email_code"
	ManagedAuthCanReauthReasonRequiresCustomerInput            ManagedAuthCanReauthReason = "requires_customer_input"
)

// Canonical auth-flow choice awaiting user selection.
type ManagedAuthChoice struct {
	// Stable choice identifier for canonical submit.
	ID string `json:"id" api:"required"`
	// Human-readable choice label.
	Label string `json:"label" api:"required"`
	// Choice type.
	//
	// Any of "mfa_method", "sso_provider", "sign_in_method", "auth_method",
	// "identifier_method", "account", "other".
	Type string `json:"type" api:"required"`
	// Context captured for a choice.
	Context string `json:"context" api:"nullable"`
	// Additional context for the choice.
	Description string `json:"description" api:"nullable"`
	// Display text captured for a choice.
	DisplayText string `json:"display_text" api:"nullable"`
	// Masked phone number or email address shown for an MFA choice.
	MaskedDestination string `json:"masked_destination" api:"nullable"`
	// Semantic MFA method. Choice id remains the stable identity of the exact option
	// selected.
	//
	// Any of "sms", "call", "email", "totp", "push", "password", "passkey", "switch",
	// "other".
	MfaType string `json:"mfa_type" api:"nullable"`
	// Selector for the visible choice, when available.
	ObservedSelector string `json:"observed_selector" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		Label             respjson.Field
		Type              respjson.Field
		Context           respjson.Field
		Description       respjson.Field
		DisplayText       respjson.Field
		MaskedDestination respjson.Field
		MfaType           respjson.Field
		ObservedSelector  respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthChoice) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthChoice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

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

// Canonical field awaiting user input.
type ManagedAuthField struct {
	// Stable field identifier for canonical submit.
	ID string `json:"id" api:"required"`
	// Why the field requires user input.
	//
	// Any of "missing", "rejected".
	Reason string `json:"reason" api:"required"`
	// Credential reference name to store the submitted value under.
	Ref string `json:"ref" api:"required"`
	// Managed-auth field type.
	//
	// Any of "identifier", "password", "code", "totp_code", "totp_secret", "text".
	Type string `json:"type" api:"required"`
	// Context shown near the field, including a masked code destination.
	Hint string `json:"hint"`
	// Human-readable label shown to the user.
	Label string `json:"label"`
	// Selector for the visible field, when available.
	ObservedSelector string `json:"observed_selector" api:"nullable"`
	// Whether this field is required.
	Required bool `json:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		Reason           respjson.Field
		Ref              respjson.Field
		Type             respjson.Field
		Hint             respjson.Field
		Label            respjson.Field
		ObservedSelector respjson.Field
		Required         respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthField) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthField) UnmarshalJSON(data []byte) error {
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

// Browser configuration applied to browser sessions created for a managed auth
// connection. Managed auth controls the profile, headless mode, timeout, start
// URL, kiosk mode, and viewport.
type ManagedAuthBrowserConfig struct {
	// Proxy configuration for managed auth browser sessions. Omit on create to derive
	// the default from stealth, or on update and login to preserve or inherit the
	// connection default.
	Proxy BrowserProxyConfig `json:"proxy"`
	// Whether managed auth browser sessions use stealth mode. Defaults to true when
	// omitted.
	Stealth bool `json:"stealth"`
	// Browser telemetry configuration using the same semantics as browser create.
	Telemetry ManagedAuthBrowserConfigTelemetry `json:"telemetry" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Proxy       respjson.Field
		Stealth     respjson.Field
		Telemetry   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthBrowserConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthBrowserConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAuthBrowserConfig to a
// ManagedAuthBrowserConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAuthBrowserConfigParam.Overrides()
func (r ManagedAuthBrowserConfig) ToParam() ManagedAuthBrowserConfigParam {
	return param.Override[ManagedAuthBrowserConfigParam](json.RawMessage(r.RawJSON()))
}

// Browser telemetry configuration using the same semantics as browser create.
type ManagedAuthBrowserConfigTelemetry struct {
	// Per-category capture flags. The operational categories (control, connection,
	// system, captcha) are captured whenever telemetry is enabled; set one to
	// enabled=false to opt out. The CDP categories (console, network, page,
	// interaction), screenshot and platform are off by default; set enabled=true to
	// opt in. On create, provided categories layer onto the default set. On update,
	// provided categories merge onto the session's current config; when no telemetry
	// is active this falls back to the default set (matching create). If browser is
	// omitted or empty, the default set is used. A browser config that disables every
	// category stops capture on update and starts no capture on create.
	Browser BrowserTelemetryCategoriesConfig `json:"browser"`
	// Request shortcut for browser telemetry capture. True enables capture; with no
	// browser category settings it captures the default set (control, connection,
	// system, captcha), and any browser category settings are layered onto that
	// default set. On update, enabled=true resolves the config fresh from the default
	// set plus any provided categories, replacing the session's current selection
	// rather than merging onto it; omit enabled to merge categories onto the current
	// selection instead. False stops capture on update and starts no capture on
	// create. enabled=false cannot be combined with browser category settings.
	Enabled bool `json:"enabled"`
	// Where to export this session's captured telemetry. Omit to capture without
	// exporting.
	Export ManagedAuthBrowserConfigTelemetryExport `json:"export"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Browser     respjson.Field
		Enabled     respjson.Field
		Export      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthBrowserConfigTelemetry) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthBrowserConfigTelemetry) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where to export this session's captured telemetry. Omit to capture without
// exporting.
type ManagedAuthBrowserConfigTelemetryExport struct {
	// Export captured telemetry over OTLP to one of the org's configured destinations.
	Otlp ManagedAuthBrowserConfigTelemetryExportOtlp `json:"otlp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Otlp        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthBrowserConfigTelemetryExport) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthBrowserConfigTelemetryExport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Export captured telemetry over OTLP to one of the org's configured destinations.
type ManagedAuthBrowserConfigTelemetryExportOtlp struct {
	// OTLP destination to export this session's captured telemetry to. Provide either
	// id or name. Requires telemetry capture to be enabled.
	Destination ManagedAuthBrowserConfigTelemetryExportOtlpDestination `json:"destination"`
	// Whether to export captured telemetry over OTLP. Setting destination implies
	// enabled=true, so this only needs to be set explicitly to disable export
	// (enabled=false with a destination is rejected).
	Enabled bool `json:"enabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Destination respjson.Field
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthBrowserConfigTelemetryExportOtlp) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthBrowserConfigTelemetryExportOtlp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OTLP destination to export this session's captured telemetry to. Provide either
// id or name. Requires telemetry capture to be enabled.
type ManagedAuthBrowserConfigTelemetryExportOtlpDestination struct {
	// OTLP destination ID
	ID string `json:"id"`
	// OTLP destination name
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthBrowserConfigTelemetryExportOtlpDestination) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthBrowserConfigTelemetryExportOtlpDestination) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser configuration applied to browser sessions created for a managed auth
// connection. Managed auth controls the profile, headless mode, timeout, start
// URL, kiosk mode, and viewport.
type ManagedAuthBrowserConfigParam struct {
	// Whether managed auth browser sessions use stealth mode. Defaults to true when
	// omitted.
	Stealth param.Opt[bool] `json:"stealth,omitzero"`
	// Browser telemetry configuration using the same semantics as browser create.
	Telemetry ManagedAuthBrowserConfigTelemetryParam `json:"telemetry,omitzero"`
	// Proxy configuration for managed auth browser sessions. Omit on create to derive
	// the default from stealth, or on update and login to preserve or inherit the
	// connection default.
	Proxy BrowserProxyConfigParam `json:"proxy,omitzero"`
	paramObj
}

func (r ManagedAuthBrowserConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthBrowserConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthBrowserConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser telemetry configuration using the same semantics as browser create.
type ManagedAuthBrowserConfigTelemetryParam struct {
	// Request shortcut for browser telemetry capture. True enables capture; with no
	// browser category settings it captures the default set (control, connection,
	// system, captcha), and any browser category settings are layered onto that
	// default set. On update, enabled=true resolves the config fresh from the default
	// set plus any provided categories, replacing the session's current selection
	// rather than merging onto it; omit enabled to merge categories onto the current
	// selection instead. False stops capture on update and starts no capture on
	// create. enabled=false cannot be combined with browser category settings.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Per-category capture flags. The operational categories (control, connection,
	// system, captcha) are captured whenever telemetry is enabled; set one to
	// enabled=false to opt out. The CDP categories (console, network, page,
	// interaction), screenshot and platform are off by default; set enabled=true to
	// opt in. On create, provided categories layer onto the default set. On update,
	// provided categories merge onto the session's current config; when no telemetry
	// is active this falls back to the default set (matching create). If browser is
	// omitted or empty, the default set is used. A browser config that disables every
	// category stops capture on update and starts no capture on create.
	Browser BrowserTelemetryCategoriesConfigParam `json:"browser,omitzero"`
	// Where to export this session's captured telemetry. Omit to capture without
	// exporting.
	Export ManagedAuthBrowserConfigTelemetryExportParam `json:"export,omitzero"`
	paramObj
}

func (r ManagedAuthBrowserConfigTelemetryParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthBrowserConfigTelemetryParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthBrowserConfigTelemetryParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where to export this session's captured telemetry. Omit to capture without
// exporting.
type ManagedAuthBrowserConfigTelemetryExportParam struct {
	// Export captured telemetry over OTLP to one of the org's configured destinations.
	Otlp ManagedAuthBrowserConfigTelemetryExportOtlpParam `json:"otlp,omitzero"`
	paramObj
}

func (r ManagedAuthBrowserConfigTelemetryExportParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthBrowserConfigTelemetryExportParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthBrowserConfigTelemetryExportParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Export captured telemetry over OTLP to one of the org's configured destinations.
type ManagedAuthBrowserConfigTelemetryExportOtlpParam struct {
	// Whether to export captured telemetry over OTLP. Setting destination implies
	// enabled=true, so this only needs to be set explicitly to disable export
	// (enabled=false with a destination is rejected).
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// OTLP destination to export this session's captured telemetry to. Provide either
	// id or name. Requires telemetry capture to be enabled.
	Destination ManagedAuthBrowserConfigTelemetryExportOtlpDestinationParam `json:"destination,omitzero"`
	paramObj
}

func (r ManagedAuthBrowserConfigTelemetryExportOtlpParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthBrowserConfigTelemetryExportOtlpParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthBrowserConfigTelemetryExportOtlpParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OTLP destination to export this session's captured telemetry to. Provide either
// id or name. Requires telemetry capture to be enabled.
type ManagedAuthBrowserConfigTelemetryExportOtlpDestinationParam struct {
	// OTLP destination ID
	ID param.Opt[string] `json:"id,omitzero"`
	// OTLP destination name
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r ManagedAuthBrowserConfigTelemetryExportOtlpDestinationParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthBrowserConfigTelemetryExportOtlpDestinationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthBrowserConfigTelemetryExportOtlpDestinationParam) UnmarshalJSON(data []byte) error {
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
	// if needed. Maximum is 86400 (24 hours). Default is 3600 (1 hour) or your plan
	// minimum, whichever is larger. The minimum depends on your plan: Enterprise: 300
	// (5 minutes), Startup: 1200 (20 minutes), Hobbyist: 3600 (1 hour), Free: 21600 (6
	// hours).
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
	// Deprecated. Use browser.telemetry. Retained during migration for existing
	// clients.
	//
	// Deprecated: deprecated
	BrowserTelemetry ManagedAuthCreateRequestBrowserTelemetryParam `json:"browser_telemetry,omitzero"`
	// Additional hostname roots valid for this auth flow, besides the primary domain.
	// Each value allows credential entry on that exact hostname and its subdomains.
	// Leading `www.` and `*.` labels are normalized away. When omitted or empty,
	// credential entry is unrestricted.
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
	// Default browser configuration for login, reauthentication, and health-check
	// sessions.
	Browser ManagedAuthBrowserConfigParam `json:"browser,omitzero"`
	// Reference to credentials for the auth connection. Use one of:
	//
	// - { name } for Kernel credentials
	// - { provider, path } for external provider item
	// - { provider, auto: true } for external provider domain lookup
	Credential ManagedAuthCreateRequestCredentialParam `json:"credential,omitzero"`
	// Deprecated. Use browser.proxy. Retained during migration for existing clients.
	//
	// Deprecated: deprecated
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

// Deprecated. Use browser.telemetry. Retained during migration for existing
// clients.
//
// Deprecated: deprecated
type ManagedAuthCreateRequestBrowserTelemetryParam struct {
	// Request shortcut for browser telemetry capture. True enables capture; with no
	// browser category settings it captures the default set (control, connection,
	// system, captcha), and any browser category settings are layered onto that
	// default set. On update, enabled=true resolves the config fresh from the default
	// set plus any provided categories, replacing the session's current selection
	// rather than merging onto it; omit enabled to merge categories onto the current
	// selection instead. False stops capture on update and starts no capture on
	// create. enabled=false cannot be combined with browser category settings.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Per-category capture flags. The operational categories (control, connection,
	// system, captcha) are captured whenever telemetry is enabled; set one to
	// enabled=false to opt out. The CDP categories (console, network, page,
	// interaction), screenshot and platform are off by default; set enabled=true to
	// opt in. On create, provided categories layer onto the default set. On update,
	// provided categories merge onto the session's current config; when no telemetry
	// is active this falls back to the default set (matching create). If browser is
	// omitted or empty, the default set is used. A browser config that disables every
	// category stops capture on update and starts no capture on create.
	Browser BrowserTelemetryCategoriesConfigParam `json:"browser,omitzero"`
	// Where to export this session's captured telemetry. Omit to capture without
	// exporting.
	Export ManagedAuthCreateRequestBrowserTelemetryExportParam `json:"export,omitzero"`
	paramObj
}

func (r ManagedAuthCreateRequestBrowserTelemetryParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthCreateRequestBrowserTelemetryParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthCreateRequestBrowserTelemetryParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where to export this session's captured telemetry. Omit to capture without
// exporting.
type ManagedAuthCreateRequestBrowserTelemetryExportParam struct {
	// Export captured telemetry over OTLP to one of the org's configured destinations.
	Otlp ManagedAuthCreateRequestBrowserTelemetryExportOtlpParam `json:"otlp,omitzero"`
	paramObj
}

func (r ManagedAuthCreateRequestBrowserTelemetryExportParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthCreateRequestBrowserTelemetryExportParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthCreateRequestBrowserTelemetryExportParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Export captured telemetry over OTLP to one of the org's configured destinations.
type ManagedAuthCreateRequestBrowserTelemetryExportOtlpParam struct {
	// Whether to export captured telemetry over OTLP. Setting destination implies
	// enabled=true, so this only needs to be set explicitly to disable export
	// (enabled=false with a destination is rejected).
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// OTLP destination to export this session's captured telemetry to. Provide either
	// id or name. Requires telemetry capture to be enabled.
	Destination ManagedAuthCreateRequestBrowserTelemetryExportOtlpDestinationParam `json:"destination,omitzero"`
	paramObj
}

func (r ManagedAuthCreateRequestBrowserTelemetryExportOtlpParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthCreateRequestBrowserTelemetryExportOtlpParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthCreateRequestBrowserTelemetryExportOtlpParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OTLP destination to export this session's captured telemetry to. Provide either
// id or name. Requires telemetry capture to be enabled.
type ManagedAuthCreateRequestBrowserTelemetryExportOtlpDestinationParam struct {
	// OTLP destination ID
	ID param.Opt[string] `json:"id,omitzero"`
	// OTLP destination name
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r ManagedAuthCreateRequestBrowserTelemetryExportOtlpDestinationParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthCreateRequestBrowserTelemetryExportOtlpDestinationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthCreateRequestBrowserTelemetryExportOtlpDestinationParam) UnmarshalJSON(data []byte) error {
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

// Deprecated. Use browser.proxy. Retained during migration for existing clients.
//
// Deprecated: deprecated
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

// A single event in an auth connection's history — a login attempt, an automatic
// re-auth attempt, or a health check.
type ManagedAuthTimelineEvent struct {
	// Identifier of the underlying login/reauth session or health check.
	ID string `json:"id" api:"required"`
	// Outcome of the event. For login/reauth events this is the flow status
	// (IN_PROGRESS, SUCCESS, EXPIRED, CANCELED, FAILED). For health_check events it is
	// the observed session state (AUTHENTICATED, NEEDS_AUTH).
	//
	// Any of "IN_PROGRESS", "SUCCESS", "EXPIRED", "CANCELED", "FAILED",
	// "AUTHENTICATED", "NEEDS_AUTH".
	Status ManagedAuthTimelineEventStatus `json:"status" api:"required"`
	// When the event occurred.
	Timestamp time.Time `json:"timestamp" api:"required" format:"date-time"`
	// The kind of event. "login" and "reauth" are authentication attempts;
	// "health_check" is a periodic session-validity check.
	//
	// Any of "login", "reauth", "health_check".
	Type ManagedAuthTimelineEventType `json:"type" api:"required"`
	// Browser session that produced the event, if one was created.
	BrowserSessionID string `json:"browser_session_id"`
	// Machine-readable error code. Present when a login/reauth event failed.
	ErrorCode string `json:"error_code"`
	// Human-readable error message. Present when a login/reauth event failed.
	ErrorMessage string `json:"error_message"`
	// The session state observed before this event. Present for health_check events
	// that recorded a prior state.
	//
	// Any of "AUTHENTICATED", "NEEDS_AUTH".
	PreviousStatus ManagedAuthTimelineEventPreviousStatus `json:"previous_status"`
	// Replay recording ID for the event's browser session, if session recording was
	// enabled.
	ReplayID string `json:"replay_id"`
	// The step the flow reached. Present for login/reauth events.
	//
	// Any of "INITIALIZED", "DISCOVERING", "AWAITING_INPUT",
	// "AWAITING_EXTERNAL_ACTION", "AWAITING_HUMAN_INTERVENTION", "SUBMITTING",
	// "COMPLETED", "EXPIRED".
	Step ManagedAuthTimelineEventStep `json:"step"`
	// Whether browser telemetry capture started for this event's browser session.
	TelemetryCaptured bool `json:"telemetry_captured"`
	// When the event was last updated. Present for login/reauth events.
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	// Visible error message from the website (e.g., 'Incorrect password'). Present
	// when the website displayed an error during the attempt.
	WebsiteError string `json:"website_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		Status            respjson.Field
		Timestamp         respjson.Field
		Type              respjson.Field
		BrowserSessionID  respjson.Field
		ErrorCode         respjson.Field
		ErrorMessage      respjson.Field
		PreviousStatus    respjson.Field
		ReplayID          respjson.Field
		Step              respjson.Field
		TelemetryCaptured respjson.Field
		UpdatedAt         respjson.Field
		WebsiteError      respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAuthTimelineEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAuthTimelineEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Outcome of the event. For login/reauth events this is the flow status
// (IN_PROGRESS, SUCCESS, EXPIRED, CANCELED, FAILED). For health_check events it is
// the observed session state (AUTHENTICATED, NEEDS_AUTH).
type ManagedAuthTimelineEventStatus string

const (
	ManagedAuthTimelineEventStatusInProgress    ManagedAuthTimelineEventStatus = "IN_PROGRESS"
	ManagedAuthTimelineEventStatusSuccess       ManagedAuthTimelineEventStatus = "SUCCESS"
	ManagedAuthTimelineEventStatusExpired       ManagedAuthTimelineEventStatus = "EXPIRED"
	ManagedAuthTimelineEventStatusCanceled      ManagedAuthTimelineEventStatus = "CANCELED"
	ManagedAuthTimelineEventStatusFailed        ManagedAuthTimelineEventStatus = "FAILED"
	ManagedAuthTimelineEventStatusAuthenticated ManagedAuthTimelineEventStatus = "AUTHENTICATED"
	ManagedAuthTimelineEventStatusNeedsAuth     ManagedAuthTimelineEventStatus = "NEEDS_AUTH"
)

// The kind of event. "login" and "reauth" are authentication attempts;
// "health_check" is a periodic session-validity check.
type ManagedAuthTimelineEventType string

const (
	ManagedAuthTimelineEventTypeLogin       ManagedAuthTimelineEventType = "login"
	ManagedAuthTimelineEventTypeReauth      ManagedAuthTimelineEventType = "reauth"
	ManagedAuthTimelineEventTypeHealthCheck ManagedAuthTimelineEventType = "health_check"
)

// The session state observed before this event. Present for health_check events
// that recorded a prior state.
type ManagedAuthTimelineEventPreviousStatus string

const (
	ManagedAuthTimelineEventPreviousStatusAuthenticated ManagedAuthTimelineEventPreviousStatus = "AUTHENTICATED"
	ManagedAuthTimelineEventPreviousStatusNeedsAuth     ManagedAuthTimelineEventPreviousStatus = "NEEDS_AUTH"
)

// The step the flow reached. Present for login/reauth events.
type ManagedAuthTimelineEventStep string

const (
	ManagedAuthTimelineEventStepInitialized               ManagedAuthTimelineEventStep = "INITIALIZED"
	ManagedAuthTimelineEventStepDiscovering               ManagedAuthTimelineEventStep = "DISCOVERING"
	ManagedAuthTimelineEventStepAwaitingInput             ManagedAuthTimelineEventStep = "AWAITING_INPUT"
	ManagedAuthTimelineEventStepAwaitingExternalAction    ManagedAuthTimelineEventStep = "AWAITING_EXTERNAL_ACTION"
	ManagedAuthTimelineEventStepAwaitingHumanIntervention ManagedAuthTimelineEventStep = "AWAITING_HUMAN_INTERVENTION"
	ManagedAuthTimelineEventStepSubmitting                ManagedAuthTimelineEventStep = "SUBMITTING"
	ManagedAuthTimelineEventStepCompleted                 ManagedAuthTimelineEventStep = "COMPLETED"
	ManagedAuthTimelineEventStepExpired                   ManagedAuthTimelineEventStep = "EXPIRED"
)

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
	// Deprecated. Use browser.telemetry. Retained during migration for existing
	// clients.
	//
	// Deprecated: deprecated
	BrowserTelemetry ManagedAuthUpdateRequestBrowserTelemetryParam `json:"browser_telemetry,omitzero"`
	// Additional hostname roots valid for this auth flow. Each value allows credential
	// entry on that exact hostname and its subdomains; leading `www.` and `*.` labels
	// are normalized away. An empty list leaves credential entry unrestricted.
	// Replaces the existing list.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	// Browser configuration updates for future login, reauthentication, and
	// health-check sessions. Omitted properties remain unchanged.
	Browser ManagedAuthBrowserConfigParam `json:"browser,omitzero"`
	// Reference to credentials for the auth connection. Use one of:
	//
	// - { name } for Kernel credentials
	// - { provider, path } for external provider item
	// - { provider, auto: true } for external provider domain lookup
	Credential ManagedAuthUpdateRequestCredentialParam `json:"credential,omitzero"`
	// Deprecated. Use browser.proxy. Retained during migration for existing clients.
	//
	// Deprecated: deprecated
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

// Deprecated. Use browser.telemetry. Retained during migration for existing
// clients.
//
// Deprecated: deprecated
type ManagedAuthUpdateRequestBrowserTelemetryParam struct {
	// Request shortcut for browser telemetry capture. True enables capture; with no
	// browser category settings it captures the default set (control, connection,
	// system, captcha), and any browser category settings are layered onto that
	// default set. On update, enabled=true resolves the config fresh from the default
	// set plus any provided categories, replacing the session's current selection
	// rather than merging onto it; omit enabled to merge categories onto the current
	// selection instead. False stops capture on update and starts no capture on
	// create. enabled=false cannot be combined with browser category settings.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Per-category capture flags. The operational categories (control, connection,
	// system, captcha) are captured whenever telemetry is enabled; set one to
	// enabled=false to opt out. The CDP categories (console, network, page,
	// interaction), screenshot and platform are off by default; set enabled=true to
	// opt in. On create, provided categories layer onto the default set. On update,
	// provided categories merge onto the session's current config; when no telemetry
	// is active this falls back to the default set (matching create). If browser is
	// omitted or empty, the default set is used. A browser config that disables every
	// category stops capture on update and starts no capture on create.
	Browser BrowserTelemetryCategoriesConfigParam `json:"browser,omitzero"`
	// Where to export this session's captured telemetry. Omit to capture without
	// exporting.
	Export ManagedAuthUpdateRequestBrowserTelemetryExportParam `json:"export,omitzero"`
	paramObj
}

func (r ManagedAuthUpdateRequestBrowserTelemetryParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthUpdateRequestBrowserTelemetryParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthUpdateRequestBrowserTelemetryParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where to export this session's captured telemetry. Omit to capture without
// exporting.
type ManagedAuthUpdateRequestBrowserTelemetryExportParam struct {
	// Export captured telemetry over OTLP to one of the org's configured destinations.
	Otlp ManagedAuthUpdateRequestBrowserTelemetryExportOtlpParam `json:"otlp,omitzero"`
	paramObj
}

func (r ManagedAuthUpdateRequestBrowserTelemetryExportParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthUpdateRequestBrowserTelemetryExportParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthUpdateRequestBrowserTelemetryExportParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Export captured telemetry over OTLP to one of the org's configured destinations.
type ManagedAuthUpdateRequestBrowserTelemetryExportOtlpParam struct {
	// Whether to export captured telemetry over OTLP. Setting destination implies
	// enabled=true, so this only needs to be set explicitly to disable export
	// (enabled=false with a destination is rejected).
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// OTLP destination to export this session's captured telemetry to. Provide either
	// id or name. Requires telemetry capture to be enabled.
	Destination ManagedAuthUpdateRequestBrowserTelemetryExportOtlpDestinationParam `json:"destination,omitzero"`
	paramObj
}

func (r ManagedAuthUpdateRequestBrowserTelemetryExportOtlpParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthUpdateRequestBrowserTelemetryExportOtlpParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthUpdateRequestBrowserTelemetryExportOtlpParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OTLP destination to export this session's captured telemetry to. Provide either
// id or name. Requires telemetry capture to be enabled.
type ManagedAuthUpdateRequestBrowserTelemetryExportOtlpDestinationParam struct {
	// OTLP destination ID
	ID param.Opt[string] `json:"id,omitzero"`
	// OTLP destination name
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r ManagedAuthUpdateRequestBrowserTelemetryExportOtlpDestinationParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAuthUpdateRequestBrowserTelemetryExportOtlpDestinationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAuthUpdateRequestBrowserTelemetryExportOtlpDestinationParam) UnmarshalJSON(data []byte) error {
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

// Deprecated. Use browser.proxy. Retained during migration for existing clients.
//
// Deprecated: deprecated
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
// select a sign-in option. Prefer canonical selected_choice_id/field_values when
// the API returns fields/choices; legacy
// fields/sso_button_selector/sso_provider/mfa_option_id/sign_in_option_id remain
// supported during deprecation.
type SubmitFieldsRequestParam struct {
	// Opaque interaction ID returned with canonical fields and choices. Required for
	// canonical submissions.
	InteractionID param.Opt[string] `json:"interaction_id,omitzero"`
	// The MFA method type to select (when mfa_options were returned)
	MfaOptionID param.Opt[string] `json:"mfa_option_id,omitzero"`
	// Canonical choice ID selected by the user.
	SelectedChoiceID param.Opt[string] `json:"selected_choice_id,omitzero"`
	// The sign-in option ID to select (when sign_in_options were returned)
	SignInOptionID param.Opt[string] `json:"sign_in_option_id,omitzero"`
	// XPath selector for the SSO button to click (ODA). Use sso_provider instead for
	// CUA.
	SSOButtonSelector param.Opt[string] `json:"sso_button_selector,omitzero"`
	// SSO provider to click, matching the provider field from pending_sso_buttons
	// (e.g., "google", "github"). Cannot be used with sso_button_selector.
	SSOProvider param.Opt[string] `json:"sso_provider,omitzero"`
	// Canonical map of field ID to submitted value.
	FieldValues map[string]string `json:"field_values,omitzero"`
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
	Choices []AuthConnectionFollowResponseManagedAuthStateChoice `json:"choices"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	DiscoveredFields []AuthConnectionFollowResponseManagedAuthStateDiscoveredField `json:"discovered_fields"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	ErrorCode string `json:"error_code"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	ErrorMessage string `json:"error_message"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	ExternalActionMessage string `json:"external_action_message"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	Fields []AuthConnectionFollowResponseManagedAuthStateField `json:"fields"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	FlowType string `json:"flow_type"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	HostedURL string `json:"hosted_url"`
	// This field is from variant [AuthConnectionFollowResponseManagedAuthState].
	InteractionID string `json:"interaction_id"`
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
		Choices               respjson.Field
		DiscoveredFields      respjson.Field
		ErrorCode             respjson.Field
		ErrorMessage          respjson.Field
		ExternalActionMessage respjson.Field
		Fields                respjson.Field
		FlowType              respjson.Field
		HostedURL             respjson.Field
		InteractionID         respjson.Field
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
	// Canonical choices awaiting selection. Prefer this over pending_sso_buttons,
	// mfa_options, and sign_in_options when present.
	Choices []AuthConnectionFollowResponseManagedAuthStateChoice `json:"choices"`
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
	// Canonical fields awaiting input. Prefer this over discovered_fields when
	// present.
	Fields []AuthConnectionFollowResponseManagedAuthStateField `json:"fields"`
	// Type of the current flow.
	//
	// Any of "LOGIN", "REAUTH".
	FlowType string `json:"flow_type"`
	// URL to redirect user to for hosted login.
	HostedURL string `json:"hosted_url" format:"uri"`
	// Opaque identifier for the current canonical interaction. Required when
	// submitting fields or choices and changes for each new actionable pause.
	InteractionID string `json:"interaction_id"`
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
		Choices               respjson.Field
		DiscoveredFields      respjson.Field
		ErrorCode             respjson.Field
		ErrorMessage          respjson.Field
		ExternalActionMessage respjson.Field
		Fields                respjson.Field
		FlowType              respjson.Field
		HostedURL             respjson.Field
		InteractionID         respjson.Field
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

// Canonical auth-flow choice awaiting user selection.
type AuthConnectionFollowResponseManagedAuthStateChoice struct {
	// Stable choice identifier for canonical submit.
	ID string `json:"id" api:"required"`
	// Human-readable choice label.
	Label string `json:"label" api:"required"`
	// Choice type.
	//
	// Any of "mfa_method", "sso_provider", "sign_in_method", "auth_method",
	// "identifier_method", "account", "other".
	Type string `json:"type" api:"required"`
	// Context captured for a choice.
	Context string `json:"context" api:"nullable"`
	// Additional context for the choice.
	Description string `json:"description" api:"nullable"`
	// Display text captured for a choice.
	DisplayText string `json:"display_text" api:"nullable"`
	// Masked phone number or email address shown for an MFA choice.
	MaskedDestination string `json:"masked_destination" api:"nullable"`
	// Semantic MFA method. Choice id remains the stable identity of the exact option
	// selected.
	//
	// Any of "sms", "call", "email", "totp", "push", "password", "passkey", "switch",
	// "other".
	MfaType string `json:"mfa_type" api:"nullable"`
	// Selector for the visible choice, when available.
	ObservedSelector string `json:"observed_selector" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		Label             respjson.Field
		Type              respjson.Field
		Context           respjson.Field
		Description       respjson.Field
		DisplayText       respjson.Field
		MaskedDestination respjson.Field
		MfaType           respjson.Field
		ObservedSelector  respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthConnectionFollowResponseManagedAuthStateChoice) RawJSON() string { return r.JSON.raw }
func (r *AuthConnectionFollowResponseManagedAuthStateChoice) UnmarshalJSON(data []byte) error {
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

// Canonical field awaiting user input.
type AuthConnectionFollowResponseManagedAuthStateField struct {
	// Stable field identifier for canonical submit.
	ID string `json:"id" api:"required"`
	// Why the field requires user input.
	//
	// Any of "missing", "rejected".
	Reason string `json:"reason" api:"required"`
	// Credential reference name to store the submitted value under.
	Ref string `json:"ref" api:"required"`
	// Managed-auth field type.
	//
	// Any of "identifier", "password", "code", "totp_code", "totp_secret", "text".
	Type string `json:"type" api:"required"`
	// Context shown near the field, including a masked code destination.
	Hint string `json:"hint"`
	// Human-readable label shown to the user.
	Label string `json:"label"`
	// Selector for the visible field, when available.
	ObservedSelector string `json:"observed_selector" api:"nullable"`
	// Whether this field is required.
	Required bool `json:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		Reason           respjson.Field
		Ref              respjson.Field
		Type             respjson.Field
		Hint             respjson.Field
		Label            respjson.Field
		ObservedSelector respjson.Field
		Required         respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthConnectionFollowResponseManagedAuthStateField) RawJSON() string { return r.JSON.raw }
func (r *AuthConnectionFollowResponseManagedAuthStateField) UnmarshalJSON(data []byte) error {
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
	// Search auth connections by ID, domain, or profile name.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
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
	// Deprecated. Use browser.telemetry. Retained during migration for existing
	// clients.
	BrowserTelemetry AuthConnectionLoginParamsBrowserTelemetry `json:"browser_telemetry,omitzero"`
	// Browser configuration override for this login. Omitted properties inherit the
	// connection defaults.
	Browser ManagedAuthBrowserConfigParam `json:"browser,omitzero"`
	// Deprecated. Use browser.proxy. Retained during migration for existing clients.
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

// Deprecated. Use browser.telemetry. Retained during migration for existing
// clients.
//
// Deprecated: deprecated
type AuthConnectionLoginParamsBrowserTelemetry struct {
	// Request shortcut for browser telemetry capture. True enables capture; with no
	// browser category settings it captures the default set (control, connection,
	// system, captcha), and any browser category settings are layered onto that
	// default set. On update, enabled=true resolves the config fresh from the default
	// set plus any provided categories, replacing the session's current selection
	// rather than merging onto it; omit enabled to merge categories onto the current
	// selection instead. False stops capture on update and starts no capture on
	// create. enabled=false cannot be combined with browser category settings.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Per-category capture flags. The operational categories (control, connection,
	// system, captcha) are captured whenever telemetry is enabled; set one to
	// enabled=false to opt out. The CDP categories (console, network, page,
	// interaction), screenshot and platform are off by default; set enabled=true to
	// opt in. On create, provided categories layer onto the default set. On update,
	// provided categories merge onto the session's current config; when no telemetry
	// is active this falls back to the default set (matching create). If browser is
	// omitted or empty, the default set is used. A browser config that disables every
	// category stops capture on update and starts no capture on create.
	Browser BrowserTelemetryCategoriesConfigParam `json:"browser,omitzero"`
	// Where to export this session's captured telemetry. Omit to capture without
	// exporting.
	Export AuthConnectionLoginParamsBrowserTelemetryExport `json:"export,omitzero"`
	paramObj
}

func (r AuthConnectionLoginParamsBrowserTelemetry) MarshalJSON() (data []byte, err error) {
	type shadow AuthConnectionLoginParamsBrowserTelemetry
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AuthConnectionLoginParamsBrowserTelemetry) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where to export this session's captured telemetry. Omit to capture without
// exporting.
type AuthConnectionLoginParamsBrowserTelemetryExport struct {
	// Export captured telemetry over OTLP to one of the org's configured destinations.
	Otlp AuthConnectionLoginParamsBrowserTelemetryExportOtlp `json:"otlp,omitzero"`
	paramObj
}

func (r AuthConnectionLoginParamsBrowserTelemetryExport) MarshalJSON() (data []byte, err error) {
	type shadow AuthConnectionLoginParamsBrowserTelemetryExport
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AuthConnectionLoginParamsBrowserTelemetryExport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Export captured telemetry over OTLP to one of the org's configured destinations.
type AuthConnectionLoginParamsBrowserTelemetryExportOtlp struct {
	// Whether to export captured telemetry over OTLP. Setting destination implies
	// enabled=true, so this only needs to be set explicitly to disable export
	// (enabled=false with a destination is rejected).
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// OTLP destination to export this session's captured telemetry to. Provide either
	// id or name. Requires telemetry capture to be enabled.
	Destination AuthConnectionLoginParamsBrowserTelemetryExportOtlpDestination `json:"destination,omitzero"`
	paramObj
}

func (r AuthConnectionLoginParamsBrowserTelemetryExportOtlp) MarshalJSON() (data []byte, err error) {
	type shadow AuthConnectionLoginParamsBrowserTelemetryExportOtlp
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AuthConnectionLoginParamsBrowserTelemetryExportOtlp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OTLP destination to export this session's captured telemetry to. Provide either
// id or name. Requires telemetry capture to be enabled.
type AuthConnectionLoginParamsBrowserTelemetryExportOtlpDestination struct {
	// OTLP destination ID
	ID param.Opt[string] `json:"id,omitzero"`
	// OTLP destination name
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r AuthConnectionLoginParamsBrowserTelemetryExportOtlpDestination) MarshalJSON() (data []byte, err error) {
	type shadow AuthConnectionLoginParamsBrowserTelemetryExportOtlpDestination
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AuthConnectionLoginParamsBrowserTelemetryExportOtlpDestination) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Deprecated. Use browser.proxy. Retained during migration for existing clients.
//
// Deprecated: deprecated
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
	// select a sign-in option. Prefer canonical selected_choice_id/field_values when
	// the API returns fields/choices; legacy
	// fields/sso_button_selector/sso_provider/mfa_option_id/sign_in_option_id remain
	// supported during deprecation.
	SubmitFieldsRequest SubmitFieldsRequestParam
	paramObj
}

func (r AuthConnectionSubmitParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.SubmitFieldsRequest)
}
func (r *AuthConnectionSubmitParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthConnectionTimelineParams struct {
	// Maximum number of events to return
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Number of events to skip
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Filter the timeline to a single event type.
	//
	// Any of "login", "reauth", "health_check".
	Type AuthConnectionTimelineParamsType `query:"type,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AuthConnectionTimelineParams]'s query parameters as
// `url.Values`.
func (r AuthConnectionTimelineParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter the timeline to a single event type.
type AuthConnectionTimelineParamsType string

const (
	AuthConnectionTimelineParamsTypeLogin       AuthConnectionTimelineParamsType = "login"
	AuthConnectionTimelineParamsTypeReauth      AuthConnectionTimelineParamsType = "reauth"
	AuthConnectionTimelineParamsTypeHealthCheck AuthConnectionTimelineParamsType = "health_check"
)
