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
	"github.com/kernel/kernel-go-sdk/shared"
)

// Create and manage browser pools for acquiring and releasing browsers.
//
// BrowserPoolService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserPoolService] method instead.
type BrowserPoolService struct {
	Options []option.RequestOption
}

// NewBrowserPoolService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBrowserPoolService(opts ...option.RequestOption) (r BrowserPoolService) {
	r = BrowserPoolService{}
	r.Options = opts
	return
}

// Create a new browser pool with the specified configuration and size. Pooled
// browsers load their profile read-only: any save_changes on the profile is
// ignored (not rejected), so pooled browsers never persist changes back to the
// profile.
func (r *BrowserPoolService) New(ctx context.Context, body BrowserPoolNewParams, opts ...option.RequestOption) (res *BrowserPool, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "browser_pools"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve details for a single browser pool by its ID or name.
func (r *BrowserPoolService) Get(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *BrowserPool, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browser_pools/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Updates the configuration used to create browsers in the pool. As with creation,
// save_changes on the pool profile is ignored (not rejected); pooled browsers
// never persist changes back to the profile. To clear the profile reference, send
// `profile: { "id": "" }`. Clearing the profile also disables
// `refresh_on_profile_update`.
func (r *BrowserPoolService) Update(ctx context.Context, idOrName string, body BrowserPoolUpdateParams, opts ...option.RequestOption) (res *BrowserPool, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browser_pools/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List browser pools in the resolved project.
func (r *BrowserPoolService) List(ctx context.Context, query BrowserPoolListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[BrowserPool], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "browser_pools"
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

// List browser pools in the resolved project.
func (r *BrowserPoolService) ListAutoPaging(ctx context.Context, query BrowserPoolListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[BrowserPool] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete a browser pool and all browsers in it. By default, deletion is blocked if
// browsers are currently leased. Use force=true to terminate leased browsers.
func (r *BrowserPoolService) Delete(ctx context.Context, idOrName string, body BrowserPoolDeleteParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return err
	}
	path := fmt.Sprintf("browser_pools/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, nil, opts...)
	return err
}

// Long-polling endpoint to acquire a browser from the pool. Returns immediately
// when a browser is available, or returns 204 No Content when the poll times out.
// The client should retry the request to continue waiting for a browser. The
// acquired browser will use the pool's timeout_seconds for its idle timeout.
func (r *BrowserPoolService) Acquire(ctx context.Context, idOrName string, body BrowserPoolAcquireParams, opts ...option.RequestOption) (res *BrowserPoolAcquireResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browser_pools/%s/acquire", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Destroys all idle browsers in the pool; leased browsers are not affected.
func (r *BrowserPoolService) Flush(ctx context.Context, idOrName string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return err
	}
	path := fmt.Sprintf("browser_pools/%s/flush", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, nil, opts...)
	return err
}

// Release a browser back to the pool, optionally recreating the browser instance.
func (r *BrowserPoolService) Release(ctx context.Context, idOrName string, body BrowserPoolReleaseParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return err
	}
	path := fmt.Sprintf("browser_pools/%s/release", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

// A browser pool containing multiple identically configured browsers.
type BrowserPool struct {
	// Unique identifier for the browser pool
	ID string `json:"id" api:"required"`
	// Number of browsers currently acquired from the pool
	AcquiredCount int64 `json:"acquired_count" api:"required"`
	// Number of browsers currently available in the pool
	AvailableCount int64 `json:"available_count" api:"required"`
	// Configuration used to create all browsers in this pool
	BrowserPoolConfig BrowserPoolBrowserPoolConfig `json:"browser_pool_config" api:"required"`
	// Timestamp when the browser pool was created
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Resolved extension IDs attached to the pool, in configured load order. Empty
	// when no extensions are attached. Authoritative for programmatic consumers; the
	// extensions inside `browser_pool_config` reflect the configured selector (echoed
	// as sent on create).
	ExtensionIDs []string `json:"extension_ids" api:"required"`
	// Browser pool name, if set
	Name string `json:"name"`
	// Resolved profile ID the pool is attached to. Omitted when no profile is
	// attached. Authoritative for programmatic consumers; the profile inside
	// `browser_pool_config` reflects the configured selector (echoed as sent on
	// create).
	ProfileID string `json:"profile_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		AcquiredCount     respjson.Field
		AvailableCount    respjson.Field
		BrowserPoolConfig respjson.Field
		CreatedAt         respjson.Field
		ExtensionIDs      respjson.Field
		Name              respjson.Field
		ProfileID         respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPool) RawJSON() string { return r.JSON.raw }
func (r *BrowserPool) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration used to create all browsers in this pool
type BrowserPoolBrowserPoolConfig struct {
	// Number of browsers maintained in the pool. The maximum size is determined by
	// your organization's pooled sessions limit (the sum of all pool sizes cannot
	// exceed your limit).
	Size int64 `json:"size" api:"required"`
	// Custom Chrome enterprise policy overrides applied to all browsers in this pool.
	// Keys are Chrome enterprise policy names; values must match their expected types.
	// Blocked: kernel-managed policies (extensions, proxy, CDP/automation). See
	// https://chromeenterprise.google/policies/ The serialized JSON payload is capped
	// at 5 MiB.
	ChromePolicy map[string]any `json:"chrome_policy"`
	// List of browser extensions to load into the session. Provide each by id or name.
	Extensions []shared.BrowserExtension `json:"extensions"`
	// Percentage of the pool to fill per minute. The cap is 25 for most organizations
	// but can be raised per-organization, so only the lower bound is enforced here.
	FillRatePerMinute int64 `json:"fill_rate_per_minute"`
	// If true, launches the browser using a headless image.
	Headless bool `json:"headless"`
	// If true, launches the browser in kiosk mode to hide address bar and tabs in live
	// view.
	KioskMode bool `json:"kiosk_mode"`
	// Optional name for the browser pool. Must be unique within the project.
	Name string `json:"name"`
	// Profile configuration for browsers in a pool. Provide either id or name.
	// Profiles must be created beforehand. Unlike single browser sessions, pools load
	// the profile read-only and never persist changes back to it, so save_changes is
	// omitted here. Any save_changes value sent on a pool profile is silently ignored
	// rather than rejected.
	Profile BrowserPoolBrowserPoolConfigProfile `json:"profile"`
	// Optional proxy associated to the browser session. References a proxy in the same
	// project as the browser session.
	ProxyID string `json:"proxy_id"`
	// When true, flush idle browsers when the profile the pool uses is updated, so
	// pool browsers pick up the latest profile data. When a profile is provided during
	// creation, this defaults to true. Requires a profile to be set on the pool.
	RefreshOnProfileUpdate bool `json:"refresh_on_profile_update"`
	// Optional URL to navigate to when a new browser is warmed into the pool.
	// Best-effort: failures to navigate do not fail pool fill. Only applied to
	// newly-warmed browsers; browsers reused via release/acquire keep whatever URL the
	// previous lease left them on. Accepts any URL Chromium can resolve, including
	// chrome:// pages.
	StartURL string `json:"start_url"`
	// If true, launches the browser in stealth mode to reduce detection by anti-bot
	// mechanisms.
	Stealth bool `json:"stealth"`
	// Active telemetry configuration applied to browsers warmed into this pool, if
	// any.
	Telemetry BrowserTelemetryConfig `json:"telemetry" api:"nullable"`
	// Default idle timeout in seconds for browsers acquired from this pool before they
	// are destroyed. Minimum 10, maximum 259200 (72 hours).
	TimeoutSeconds int64 `json:"timeout_seconds"`
	// Initial browser window size in pixels with optional refresh rate. If omitted,
	// image defaults apply (1920x1080@25). For GPU images, the default is
	// 1920x1080@60. Arbitrary viewport dimensions and refresh rates are accepted.
	// Known-good presets include: 2560x1440@10, 1920x1080@25, 1920x1200@25,
	// 1440x900@25, 1280x800@60, 1024x768@60, 1200x800@60, 768x1024@60, 390x844@60. For
	// GPU images, recommended presets use one of these resolutions with refresh rates
	// 60, 30, 25, or 10: 800x600, 960x720, 1024x576, 1024x768, 1152x648, 1200x800,
	// 1280x720, 1368x768, 1440x900, 1600x900, 1920x1080, 1920x1200, 390x844, 360x250,
	// 768x1024, 800x1600. Viewports outside this list may exhibit unstable live view
	// or recording behavior. If refresh_rate is not provided, it will be automatically
	// determined based on the resolution (higher resolutions use lower refresh rates
	// to keep bandwidth reasonable).
	Viewport shared.BrowserViewport `json:"viewport"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Size                   respjson.Field
		ChromePolicy           respjson.Field
		Extensions             respjson.Field
		FillRatePerMinute      respjson.Field
		Headless               respjson.Field
		KioskMode              respjson.Field
		Name                   respjson.Field
		Profile                respjson.Field
		ProxyID                respjson.Field
		RefreshOnProfileUpdate respjson.Field
		StartURL               respjson.Field
		Stealth                respjson.Field
		Telemetry              respjson.Field
		TimeoutSeconds         respjson.Field
		Viewport               respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPoolBrowserPoolConfig) RawJSON() string { return r.JSON.raw }
func (r *BrowserPoolBrowserPoolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Profile configuration for browsers in a pool. Provide either id or name.
// Profiles must be created beforehand. Unlike single browser sessions, pools load
// the profile read-only and never persist changes back to it, so save_changes is
// omitted here. Any save_changes value sent on a pool profile is silently ignored
// rather than rejected.
type BrowserPoolBrowserPoolConfigProfile struct {
	// Profile ID to load for browsers in this pool
	ID string `json:"id"`
	// Profile name to load for browsers in this pool (instead of id). Must be 1-255
	// characters, using letters, numbers, dots, underscores, or hyphens.
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
func (r BrowserPoolBrowserPoolConfigProfile) RawJSON() string { return r.JSON.raw }
func (r *BrowserPoolBrowserPoolConfigProfile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPoolAcquireResponse struct {
	// Websocket URL for Chrome DevTools Protocol connections to the browser session
	CdpWsURL string `json:"cdp_ws_url" api:"required"`
	// When the browser session was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Whether the browser session is running in headless mode.
	Headless bool `json:"headless" api:"required"`
	// Unique identifier for the browser session
	SessionID string `json:"session_id" api:"required"`
	// Whether the browser session is running in stealth mode.
	Stealth bool `json:"stealth" api:"required"`
	// The number of seconds of inactivity before the browser session is terminated.
	TimeoutSeconds int64 `json:"timeout_seconds" api:"required"`
	// Websocket URL for WebDriver BiDi connections to the browser session
	WebdriverWsURL string `json:"webdriver_ws_url" api:"required"`
	// Metro-API HTTP base URL for this browser session.
	BaseURL string `json:"base_url"`
	// Remote URL for live viewing the browser session. Only available for non-headless
	// browsers.
	BrowserLiveViewURL string `json:"browser_live_view_url"`
	// Custom Chrome enterprise policy overrides that were applied to this browser
	// session, if any. Echoed back for verification. Keys are Chrome enterprise policy
	// names.
	ChromePolicy map[string]any `json:"chrome_policy"`
	// When the browser session was soft-deleted. Only present for deleted sessions.
	DeletedAt time.Time `json:"deleted_at" format:"date-time"`
	// Whether GPU acceleration is enabled for the browser session (only supported for
	// headful sessions).
	GPU bool `json:"gpu"`
	// Whether the browser session is running in kiosk mode.
	KioskMode bool `json:"kiosk_mode"`
	// Human-readable name of the browser session, if one was set at creation.
	Name string `json:"name"`
	// Browser pool this session was acquired from, if any.
	Pool BrowserPoolRef `json:"pool"`
	// Browser profile metadata.
	Profile Profile `json:"profile"`
	// Whether changes made during this browser session are saved back to its profile
	// when the session ends. Omitted when no profile is attached.
	ProfileSaveChanges bool `json:"profile_save_changes"`
	// ID of the proxy associated with this browser session, if any.
	ProxyID string `json:"proxy_id"`
	// URL the session was asked to navigate to on creation, if any. Recorded for
	// debugging. Navigation is fire-and-forget — the URL is dispatched to the browser
	// without waiting for it to load, and any errors (DNS failure, bad status,
	// timeout) are silently dropped. Captures what was requested, not what the browser
	// actually loaded.
	StartURL string `json:"start_url"`
	// User-defined key-value tags that were set on this browser session, if any.
	// Echoed back when present.
	Tags Tags `json:"tags"`
	// Active telemetry configuration for the session, if any.
	Telemetry BrowserTelemetryConfig `json:"telemetry" api:"nullable"`
	// Session usage metrics.
	Usage BrowserUsage `json:"usage"`
	// Initial browser window size in pixels with optional refresh rate. If omitted,
	// image defaults apply (1920x1080@25). For GPU images, the default is
	// 1920x1080@60. Arbitrary viewport dimensions and refresh rates are accepted.
	// Known-good presets include: 2560x1440@10, 1920x1080@25, 1920x1200@25,
	// 1440x900@25, 1280x800@60, 1024x768@60, 1200x800@60, 768x1024@60, 390x844@60. For
	// GPU images, recommended presets use one of these resolutions with refresh rates
	// 60, 30, 25, or 10: 800x600, 960x720, 1024x576, 1024x768, 1152x648, 1200x800,
	// 1280x720, 1368x768, 1440x900, 1600x900, 1920x1080, 1920x1200, 390x844, 360x250,
	// 768x1024, 800x1600. Viewports outside this list may exhibit unstable live view
	// or recording behavior. If refresh_rate is not provided, it will be automatically
	// determined based on the resolution (higher resolutions use lower refresh rates
	// to keep bandwidth reasonable).
	Viewport shared.BrowserViewport `json:"viewport"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CdpWsURL           respjson.Field
		CreatedAt          respjson.Field
		Headless           respjson.Field
		SessionID          respjson.Field
		Stealth            respjson.Field
		TimeoutSeconds     respjson.Field
		WebdriverWsURL     respjson.Field
		BaseURL            respjson.Field
		BrowserLiveViewURL respjson.Field
		ChromePolicy       respjson.Field
		DeletedAt          respjson.Field
		GPU                respjson.Field
		KioskMode          respjson.Field
		Name               respjson.Field
		Pool               respjson.Field
		Profile            respjson.Field
		ProfileSaveChanges respjson.Field
		ProxyID            respjson.Field
		StartURL           respjson.Field
		Tags               respjson.Field
		Telemetry          respjson.Field
		Usage              respjson.Field
		Viewport           respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPoolAcquireResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserPoolAcquireResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPoolNewParams struct {
	// Number of browsers to maintain in the pool. The maximum size is determined by
	// your organization's pooled sessions limit (the sum of all pool sizes cannot
	// exceed your limit).
	Size int64 `json:"size" api:"required"`
	// Percentage of the pool to fill per minute. Defaults to 25. The cap is 25 for
	// most organizations but can be raised per-organization, so only the lower bound
	// is enforced here.
	FillRatePerMinute param.Opt[int64] `json:"fill_rate_per_minute,omitzero"`
	// If true, launches the browser using a headless image. Defaults to false.
	Headless param.Opt[bool] `json:"headless,omitzero"`
	// If true, launches the browser in kiosk mode to hide address bar and tabs in live
	// view. Defaults to false.
	KioskMode param.Opt[bool] `json:"kiosk_mode,omitzero"`
	// Optional name for the browser pool. Must be unique within the project.
	Name param.Opt[string] `json:"name,omitzero"`
	// Optional proxy to associate to the browser session. Must reference a proxy in
	// the same project as the browser session.
	ProxyID param.Opt[string] `json:"proxy_id,omitzero"`
	// When true, flush idle browsers when the profile the pool uses is updated, so
	// pool browsers pick up the latest profile data. When a profile is provided during
	// creation, this defaults to true. Requires a profile to be set on the pool.
	RefreshOnProfileUpdate param.Opt[bool] `json:"refresh_on_profile_update,omitzero"`
	// Optional URL to navigate to when a new browser is warmed into the pool.
	// Best-effort: failures to navigate do not fail pool fill. Only applied to
	// newly-warmed browsers; browsers reused via release/acquire keep whatever URL the
	// previous lease left them on. Accepts any URL Chromium can resolve, including
	// chrome:// pages.
	StartURL param.Opt[string] `json:"start_url,omitzero"`
	// If true, launches the browser in stealth mode to reduce detection by anti-bot
	// mechanisms. Defaults to false.
	Stealth param.Opt[bool] `json:"stealth,omitzero"`
	// Default idle timeout in seconds for browsers acquired from this pool before they
	// are destroyed. Defaults to 600 seconds. Minimum 10, maximum 259200 (72 hours).
	TimeoutSeconds param.Opt[int64] `json:"timeout_seconds,omitzero"`
	// Telemetry configuration applied to browsers warmed into this pool. Set enabled
	// to true to start capture using the default set, or provide browser category
	// settings. If omitted, null, set to an empty object ({}), set to enabled: false
	// without browser category settings, or all four CDP categories are explicitly
	// disabled, no telemetry is configured on the pool. Only applied to newly-warmed
	// browsers.
	Telemetry BrowserPoolNewParamsTelemetry `json:"telemetry,omitzero"`
	// Custom Chrome enterprise policy overrides applied to all browsers in this pool.
	// Keys are Chrome enterprise policy names; values must match their expected types.
	// Blocked: kernel-managed policies (extensions, proxy, CDP/automation). See
	// https://chromeenterprise.google/policies/ The serialized JSON payload is capped
	// at 5 MiB.
	ChromePolicy map[string]any `json:"chrome_policy,omitzero"`
	// List of browser extensions to load into the session. Provide each by id or name.
	Extensions []shared.BrowserExtensionParam `json:"extensions,omitzero"`
	// Profile configuration for browsers in a pool. Provide either id or name.
	// Profiles must be created beforehand. Unlike single browser sessions, pools load
	// the profile read-only and never persist changes back to it, so save_changes is
	// omitted here. Any save_changes value sent on a pool profile is silently ignored
	// rather than rejected.
	Profile BrowserPoolNewParamsProfile `json:"profile,omitzero"`
	// Initial browser window size in pixels with optional refresh rate. If omitted,
	// image defaults apply (1920x1080@25). For GPU images, the default is
	// 1920x1080@60. Arbitrary viewport dimensions and refresh rates are accepted.
	// Known-good presets include: 2560x1440@10, 1920x1080@25, 1920x1200@25,
	// 1440x900@25, 1280x800@60, 1024x768@60, 1200x800@60, 768x1024@60, 390x844@60. For
	// GPU images, recommended presets use one of these resolutions with refresh rates
	// 60, 30, 25, or 10: 800x600, 960x720, 1024x576, 1024x768, 1152x648, 1200x800,
	// 1280x720, 1368x768, 1440x900, 1600x900, 1920x1080, 1920x1200, 390x844, 360x250,
	// 768x1024, 800x1600. Viewports outside this list may exhibit unstable live view
	// or recording behavior. If refresh_rate is not provided, it will be automatically
	// determined based on the resolution (higher resolutions use lower refresh rates
	// to keep bandwidth reasonable).
	Viewport shared.BrowserViewportParam `json:"viewport,omitzero"`
	paramObj
}

func (r BrowserPoolNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Profile configuration for browsers in a pool. Provide either id or name.
// Profiles must be created beforehand. Unlike single browser sessions, pools load
// the profile read-only and never persist changes back to it, so save_changes is
// omitted here. Any save_changes value sent on a pool profile is silently ignored
// rather than rejected.
type BrowserPoolNewParamsProfile struct {
	// Profile ID to load for browsers in this pool
	ID param.Opt[string] `json:"id,omitzero"`
	// Profile name to load for browsers in this pool (instead of id). Must be 1-255
	// characters, using letters, numbers, dots, underscores, or hyphens.
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r BrowserPoolNewParamsProfile) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolNewParamsProfile
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolNewParamsProfile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Telemetry configuration applied to browsers warmed into this pool. Set enabled
// to true to start capture using the default set, or provide browser category
// settings. If omitted, null, set to an empty object ({}), set to enabled: false
// without browser category settings, or all four CDP categories are explicitly
// disabled, no telemetry is configured on the pool. Only applied to newly-warmed
// browsers.
type BrowserPoolNewParamsTelemetry struct {
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
	// interaction) and screenshot are off by default; set enabled=true to opt in. On
	// create, provided categories layer onto the default set. On update, provided
	// categories merge onto the session's current config; when no telemetry is active
	// this falls back to the default set (matching create). If browser is omitted or
	// empty, the default set is used. A browser config that disables every category
	// stops capture on update and starts no capture on create.
	Browser BrowserTelemetryCategoriesConfigParam `json:"browser,omitzero"`
	// Where to export this session's captured telemetry. Omit to capture without
	// exporting.
	Export BrowserPoolNewParamsTelemetryExport `json:"export,omitzero"`
	paramObj
}

func (r BrowserPoolNewParamsTelemetry) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolNewParamsTelemetry
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolNewParamsTelemetry) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where to export this session's captured telemetry. Omit to capture without
// exporting.
type BrowserPoolNewParamsTelemetryExport struct {
	// Export captured telemetry over OTLP to one of the org's configured destinations.
	Otlp BrowserPoolNewParamsTelemetryExportOtlp `json:"otlp,omitzero"`
	paramObj
}

func (r BrowserPoolNewParamsTelemetryExport) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolNewParamsTelemetryExport
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolNewParamsTelemetryExport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Export captured telemetry over OTLP to one of the org's configured destinations.
type BrowserPoolNewParamsTelemetryExportOtlp struct {
	// Whether to export captured telemetry over OTLP. Setting destination implies
	// enabled=true, so this only needs to be set explicitly to disable export
	// (enabled=false with a destination is rejected).
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// OTLP destination to export this session's captured telemetry to. Provide either
	// id or name. Requires telemetry capture to be enabled.
	Destination BrowserPoolNewParamsTelemetryExportOtlpDestination `json:"destination,omitzero"`
	paramObj
}

func (r BrowserPoolNewParamsTelemetryExportOtlp) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolNewParamsTelemetryExportOtlp
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolNewParamsTelemetryExportOtlp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OTLP destination to export this session's captured telemetry to. Provide either
// id or name. Requires telemetry capture to be enabled.
type BrowserPoolNewParamsTelemetryExportOtlpDestination struct {
	// OTLP destination ID
	ID param.Opt[string] `json:"id,omitzero"`
	// OTLP destination name
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r BrowserPoolNewParamsTelemetryExportOtlpDestination) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolNewParamsTelemetryExportOtlpDestination
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolNewParamsTelemetryExportOtlpDestination) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPoolUpdateParams struct {
	// Whether to discard all idle browsers and rebuild them immediately with the new
	// configuration. Defaults to false. Only browsers that are idle when the update
	// runs are rebuilt. A browser that is in use during the update keeps its original
	// configuration, and if it is later released with `reuse: true` it returns to the
	// pool with that stale configuration until it is discarded (by this flag on a
	// later update, or by flushing the pool).
	DiscardAllIdle param.Opt[bool] `json:"discard_all_idle,omitzero"`
	// If provided, replaces the percentage of the pool to fill per minute. The cap is
	// 25 for most organizations but can be raised per-organization, so only the lower
	// bound is enforced here.
	FillRatePerMinute param.Opt[int64] `json:"fill_rate_per_minute,omitzero"`
	// If provided, replaces whether browsers launch using a headless image.
	Headless param.Opt[bool] `json:"headless,omitzero"`
	// If provided, replaces whether browsers launch in kiosk mode.
	KioskMode param.Opt[bool] `json:"kiosk_mode,omitzero"`
	// If provided, replaces the pool name. Empty string is a no-op; the pool name
	// cannot be cleared or reset to empty once assigned.
	Name param.Opt[string] `json:"name,omitzero"`
	// Empty string clears the previously-selected proxy. Omit this field to leave the
	// proxy unchanged.
	ProxyID param.Opt[string] `json:"proxy_id,omitzero"`
	// If provided, replaces whether idle browsers are flushed when the profile the
	// pool uses is updated. When the pool's profile reference is changed (including
	// newly attached) and this field is omitted, it defaults to true. Re-sending the
	// same profile reference leaves this setting unchanged. Clearing the profile also
	// disables this setting. Requires a profile to be set on the pool.
	RefreshOnProfileUpdate param.Opt[bool] `json:"refresh_on_profile_update,omitzero"`
	// If provided, replaces the number of browsers to maintain in the pool. The
	// maximum size is determined by your organization's pooled sessions limit (the sum
	// of all pool sizes cannot exceed your limit).
	Size param.Opt[int64] `json:"size,omitzero"`
	// If provided, replaces the URL to navigate to when a new browser is warmed into
	// the pool. Empty string clears the previously-set URL. Omit this field to leave
	// it unchanged.
	StartURL param.Opt[string] `json:"start_url,omitzero"`
	// If provided, replaces whether browsers launch in stealth mode.
	Stealth param.Opt[bool] `json:"stealth,omitzero"`
	// If provided, replaces the default idle timeout in seconds for browsers acquired
	// from this pool before they are destroyed. Minimum 10, maximum 259200 (72 hours).
	TimeoutSeconds param.Opt[int64] `json:"timeout_seconds,omitzero"`
	// If provided, updates the pool's telemetry configuration. Omit, set to null, or
	// set to an empty object ({}) to leave the existing configuration unchanged. Set
	// enabled to true to enable capture using the default set. Set enabled to false to
	// clear the pool's telemetry. Provide browser category settings for per-category
	// updates, merged onto the pool's current configuration. Only applied to browsers
	// warmed after the update; browsers already in the pool keep their configuration
	// until discarded.
	Telemetry BrowserPoolUpdateParamsTelemetry `json:"telemetry,omitzero"`
	// If provided, replaces the custom Chrome enterprise policy overrides applied to
	// all browsers in this pool. Empty object clears any previously-set policy. Keys
	// are Chrome enterprise policy names; values must match their expected types.
	// Blocked: kernel-managed policies (extensions, proxy, CDP/automation). See
	// https://chromeenterprise.google/policies/ The serialized JSON payload is capped
	// at 5 MiB.
	ChromePolicy map[string]any `json:"chrome_policy,omitzero"`
	// If provided, replaces the extension list. Empty array clears all
	// previously-selected extensions. Omit this field to leave extensions unchanged.
	Extensions []shared.BrowserExtensionParam `json:"extensions,omitzero"`
	// Profile configuration for browsers in a pool. Provide either id or name.
	// Profiles must be created beforehand. Unlike single browser sessions, pools load
	// the profile read-only and never persist changes back to it, so save_changes is
	// omitted here. Any save_changes value sent on a pool profile is silently ignored
	// rather than rejected.
	Profile BrowserPoolUpdateParamsProfile `json:"profile,omitzero"`
	// Initial browser window size in pixels with optional refresh rate. If omitted,
	// image defaults apply (1920x1080@25). For GPU images, the default is
	// 1920x1080@60. Arbitrary viewport dimensions and refresh rates are accepted.
	// Known-good presets include: 2560x1440@10, 1920x1080@25, 1920x1200@25,
	// 1440x900@25, 1280x800@60, 1024x768@60, 1200x800@60, 768x1024@60, 390x844@60. For
	// GPU images, recommended presets use one of these resolutions with refresh rates
	// 60, 30, 25, or 10: 800x600, 960x720, 1024x576, 1024x768, 1152x648, 1200x800,
	// 1280x720, 1368x768, 1440x900, 1600x900, 1920x1080, 1920x1200, 390x844, 360x250,
	// 768x1024, 800x1600. Viewports outside this list may exhibit unstable live view
	// or recording behavior. If refresh_rate is not provided, it will be automatically
	// determined based on the resolution (higher resolutions use lower refresh rates
	// to keep bandwidth reasonable).
	Viewport shared.BrowserViewportParam `json:"viewport,omitzero"`
	paramObj
}

func (r BrowserPoolUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Profile configuration for browsers in a pool. Provide either id or name.
// Profiles must be created beforehand. Unlike single browser sessions, pools load
// the profile read-only and never persist changes back to it, so save_changes is
// omitted here. Any save_changes value sent on a pool profile is silently ignored
// rather than rejected.
type BrowserPoolUpdateParamsProfile struct {
	// Profile ID to load for browsers in this pool
	ID param.Opt[string] `json:"id,omitzero"`
	// Profile name to load for browsers in this pool (instead of id). Must be 1-255
	// characters, using letters, numbers, dots, underscores, or hyphens.
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r BrowserPoolUpdateParamsProfile) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolUpdateParamsProfile
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolUpdateParamsProfile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// If provided, updates the pool's telemetry configuration. Omit, set to null, or
// set to an empty object ({}) to leave the existing configuration unchanged. Set
// enabled to true to enable capture using the default set. Set enabled to false to
// clear the pool's telemetry. Provide browser category settings for per-category
// updates, merged onto the pool's current configuration. Only applied to browsers
// warmed after the update; browsers already in the pool keep their configuration
// until discarded.
type BrowserPoolUpdateParamsTelemetry struct {
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
	// interaction) and screenshot are off by default; set enabled=true to opt in. On
	// create, provided categories layer onto the default set. On update, provided
	// categories merge onto the session's current config; when no telemetry is active
	// this falls back to the default set (matching create). If browser is omitted or
	// empty, the default set is used. A browser config that disables every category
	// stops capture on update and starts no capture on create.
	Browser BrowserTelemetryCategoriesConfigParam `json:"browser,omitzero"`
	// Where to export this session's captured telemetry. Omit to capture without
	// exporting.
	Export BrowserPoolUpdateParamsTelemetryExport `json:"export,omitzero"`
	paramObj
}

func (r BrowserPoolUpdateParamsTelemetry) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolUpdateParamsTelemetry
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolUpdateParamsTelemetry) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where to export this session's captured telemetry. Omit to capture without
// exporting.
type BrowserPoolUpdateParamsTelemetryExport struct {
	// Export captured telemetry over OTLP to one of the org's configured destinations.
	Otlp BrowserPoolUpdateParamsTelemetryExportOtlp `json:"otlp,omitzero"`
	paramObj
}

func (r BrowserPoolUpdateParamsTelemetryExport) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolUpdateParamsTelemetryExport
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolUpdateParamsTelemetryExport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Export captured telemetry over OTLP to one of the org's configured destinations.
type BrowserPoolUpdateParamsTelemetryExportOtlp struct {
	// Whether to export captured telemetry over OTLP. Setting destination implies
	// enabled=true, so this only needs to be set explicitly to disable export
	// (enabled=false with a destination is rejected).
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// OTLP destination to export this session's captured telemetry to. Provide either
	// id or name. Requires telemetry capture to be enabled.
	Destination BrowserPoolUpdateParamsTelemetryExportOtlpDestination `json:"destination,omitzero"`
	paramObj
}

func (r BrowserPoolUpdateParamsTelemetryExportOtlp) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolUpdateParamsTelemetryExportOtlp
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolUpdateParamsTelemetryExportOtlp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OTLP destination to export this session's captured telemetry to. Provide either
// id or name. Requires telemetry capture to be enabled.
type BrowserPoolUpdateParamsTelemetryExportOtlpDestination struct {
	// OTLP destination ID
	ID param.Opt[string] `json:"id,omitzero"`
	// OTLP destination name
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r BrowserPoolUpdateParamsTelemetryExportOtlpDestination) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolUpdateParamsTelemetryExportOtlpDestination
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolUpdateParamsTelemetryExportOtlpDestination) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPoolListParams struct {
	// Limit the number of browser pools to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Exact-match filter on browser pool name using the database collation. In
	// production, matching is case- and accent-insensitive. During the default-project
	// migration, unscoped requests prefer a concrete default-project browser pool over
	// a legacy unscoped browser pool with the same name.
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	// Offset the number of browser pools to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Case-insensitive substring match against browser pool name. IDs match by exact
	// value.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserPoolListParams]'s query parameters as `url.Values`.
func (r BrowserPoolListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BrowserPoolDeleteParams struct {
	// If true, force delete even if browsers are currently leased. Leased browsers
	// will be terminated.
	Force param.Opt[bool] `json:"force,omitzero"`
	paramObj
}

func (r BrowserPoolDeleteParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolDeleteParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolDeleteParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPoolAcquireParams struct {
	// Maximum number of seconds to wait for a browser to be available. Defaults to the
	// calculated time it would take to fill the pool at the currently configured fill
	// rate.
	AcquireTimeoutSeconds param.Opt[int64] `json:"acquire_timeout_seconds,omitzero"`
	// Optional human-readable name for the acquired browser session, used to find it
	// later in the dashboard. Must be unique among active sessions within the pool's
	// project. Applies to this lease only and is cleared when the browser is released
	// back to the pool.
	Name param.Opt[string] `json:"name,omitzero"`
	// Optional URL to navigate the acquired browser to. Overrides the pool's start_url
	// for this acquire only. Best-effort: failures to navigate do not fail the
	// acquire.
	StartURL param.Opt[string] `json:"start_url,omitzero"`
	// Telemetry override for the acquired browser, applied to this lease only. Merges
	// onto the browser's current (pool-inherited) telemetry using the same
	// per-category semantics as PATCH /browsers: provided categories override the
	// current configuration, omitted categories are inherited. Set enabled to true to
	// resolve the config fresh from the default set, or enabled to false to stop
	// capture. When the browser is released back to the pool with reuse, its telemetry
	// is reset to the pool's baseline, so the override does not carry over to the next
	// lease.
	Telemetry BrowserPoolAcquireParamsTelemetry `json:"telemetry,omitzero"`
	// Optional user-defined key-value tags for the acquired browser session, used to
	// find and group sessions later. Applies to this lease only and are cleared when
	// the browser is released back to the pool. Up to 50 pairs.
	Tags Tags `json:"tags,omitzero"`
	paramObj
}

func (r BrowserPoolAcquireParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolAcquireParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolAcquireParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Telemetry override for the acquired browser, applied to this lease only. Merges
// onto the browser's current (pool-inherited) telemetry using the same
// per-category semantics as PATCH /browsers: provided categories override the
// current configuration, omitted categories are inherited. Set enabled to true to
// resolve the config fresh from the default set, or enabled to false to stop
// capture. When the browser is released back to the pool with reuse, its telemetry
// is reset to the pool's baseline, so the override does not carry over to the next
// lease.
type BrowserPoolAcquireParamsTelemetry struct {
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
	// interaction) and screenshot are off by default; set enabled=true to opt in. On
	// create, provided categories layer onto the default set. On update, provided
	// categories merge onto the session's current config; when no telemetry is active
	// this falls back to the default set (matching create). If browser is omitted or
	// empty, the default set is used. A browser config that disables every category
	// stops capture on update and starts no capture on create.
	Browser BrowserTelemetryCategoriesConfigParam `json:"browser,omitzero"`
	// Where to export this session's captured telemetry. Omit to capture without
	// exporting.
	Export BrowserPoolAcquireParamsTelemetryExport `json:"export,omitzero"`
	paramObj
}

func (r BrowserPoolAcquireParamsTelemetry) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolAcquireParamsTelemetry
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolAcquireParamsTelemetry) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where to export this session's captured telemetry. Omit to capture without
// exporting.
type BrowserPoolAcquireParamsTelemetryExport struct {
	// Export captured telemetry over OTLP to one of the org's configured destinations.
	Otlp BrowserPoolAcquireParamsTelemetryExportOtlp `json:"otlp,omitzero"`
	paramObj
}

func (r BrowserPoolAcquireParamsTelemetryExport) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolAcquireParamsTelemetryExport
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolAcquireParamsTelemetryExport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Export captured telemetry over OTLP to one of the org's configured destinations.
type BrowserPoolAcquireParamsTelemetryExportOtlp struct {
	// Whether to export captured telemetry over OTLP. Setting destination implies
	// enabled=true, so this only needs to be set explicitly to disable export
	// (enabled=false with a destination is rejected).
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// OTLP destination to export this session's captured telemetry to. Provide either
	// id or name. Requires telemetry capture to be enabled.
	Destination BrowserPoolAcquireParamsTelemetryExportOtlpDestination `json:"destination,omitzero"`
	paramObj
}

func (r BrowserPoolAcquireParamsTelemetryExportOtlp) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolAcquireParamsTelemetryExportOtlp
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolAcquireParamsTelemetryExportOtlp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OTLP destination to export this session's captured telemetry to. Provide either
// id or name. Requires telemetry capture to be enabled.
type BrowserPoolAcquireParamsTelemetryExportOtlpDestination struct {
	// OTLP destination ID
	ID param.Opt[string] `json:"id,omitzero"`
	// OTLP destination name
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r BrowserPoolAcquireParamsTelemetryExportOtlpDestination) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolAcquireParamsTelemetryExportOtlpDestination
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolAcquireParamsTelemetryExportOtlpDestination) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPoolReleaseParams struct {
	// Browser session ID to release back to the pool
	SessionID string `json:"session_id" api:"required"`
	// Whether to reuse the browser instance or destroy it and create a new one.
	// Defaults to true. A reused browser keeps the configuration it was created with,
	// so it does not pick up pool configuration changes made while it was in use.
	// Release with `reuse: false`, or flush the pool afterward, to rebuild it with the
	// current configuration.
	Reuse param.Opt[bool] `json:"reuse,omitzero"`
	paramObj
}

func (r BrowserPoolReleaseParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolReleaseParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolReleaseParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
