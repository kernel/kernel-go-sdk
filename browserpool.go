// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
	"github.com/kernel/kernel-go-sdk/shared"
)

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

// Create a new browser pool with the specified configuration and size.
func (r *BrowserPoolService) New(ctx context.Context, body BrowserPoolNewParams, opts ...option.RequestOption) (res *BrowserPool, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "browser_pools"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieve details for a single browser pool by its ID or name.
func (r *BrowserPoolService) Get(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *BrowserPool, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("browser_pools/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Updates the configuration used to create browsers in the pool.
func (r *BrowserPoolService) Update(ctx context.Context, idOrName string, body BrowserPoolUpdateParams, opts ...option.RequestOption) (res *BrowserPool, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("browser_pools/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// List browser pools owned by the caller's organization.
func (r *BrowserPoolService) List(ctx context.Context, opts ...option.RequestOption) (res *[]BrowserPool, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "browser_pools"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Delete a browser pool and all browsers in it. By default, deletion is blocked if
// browsers are currently leased. Use force=true to terminate leased browsers.
func (r *BrowserPoolService) Delete(ctx context.Context, idOrName string, body BrowserPoolDeleteParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("browser_pools/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, nil, opts...)
	return
}

// Long-polling endpoint to acquire a browser from the pool. Returns immediately
// when a browser is available, or returns 204 No Content when the poll times out.
// The client should retry the request to continue waiting for a browser. The
// acquired browser will use the pool's timeout_seconds for its idle timeout.
func (r *BrowserPoolService) Acquire(ctx context.Context, idOrName string, body BrowserPoolAcquireParams, opts ...option.RequestOption) (res *BrowserPoolAcquireResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("browser_pools/%s/acquire", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Destroys all idle browsers in the pool; leased browsers are not affected.
func (r *BrowserPoolService) Flush(ctx context.Context, idOrName string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("browser_pools/%s/flush", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, nil, opts...)
	return
}

// Release a browser back to the pool, optionally recreating the browser instance.
func (r *BrowserPoolService) Release(ctx context.Context, idOrName string, body BrowserPoolReleaseParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("browser_pools/%s/release", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return
}

// A browser pool containing multiple identically configured browsers.
type BrowserPool struct {
	// Unique identifier for the browser pool
	ID string `json:"id,required"`
	// Number of browsers currently acquired from the pool
	AcquiredCount int64 `json:"acquired_count,required"`
	// Number of browsers currently available in the pool
	AvailableCount int64 `json:"available_count,required"`
	// Configuration used to create all browsers in this pool
	BrowserPoolConfig BrowserPoolBrowserPoolConfig `json:"browser_pool_config,required"`
	// Timestamp when the browser pool was created
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// Browser pool name, if set
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		AcquiredCount     respjson.Field
		AvailableCount    respjson.Field
		BrowserPoolConfig respjson.Field
		CreatedAt         respjson.Field
		Name              respjson.Field
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
	// Number of browsers to maintain in the pool. The maximum size is determined by
	// your organization's pooled sessions limit (the sum of all pool sizes cannot
	// exceed your limit).
	Size int64 `json:"size,required"`
	// List of browser extensions to load into the session. Provide each by id or name.
	Extensions []shared.BrowserExtension `json:"extensions"`
	// Percentage of the pool to fill per minute. Defaults to 10%.
	FillRatePerMinute int64 `json:"fill_rate_per_minute"`
	// If true, launches the browser using a headless image. Defaults to false.
	Headless bool `json:"headless"`
	// If true, launches the browser in kiosk mode to hide address bar and tabs in live
	// view.
	KioskMode bool `json:"kiosk_mode"`
	// Optional name for the browser pool. Must be unique within the organization.
	Name string `json:"name"`
	// Profile selection for the browser session. Provide either id or name. If
	// specified, the matching profile will be loaded into the browser session.
	// Profiles must be created beforehand.
	Profile shared.BrowserProfile `json:"profile"`
	// Optional proxy to associate to the browser session. Must reference a proxy
	// belonging to the caller's org.
	ProxyID string `json:"proxy_id"`
	// If true, launches the browser in stealth mode to reduce detection by anti-bot
	// mechanisms.
	Stealth bool `json:"stealth"`
	// Default idle timeout in seconds for browsers acquired from this pool before they
	// are destroyed. Defaults to 600 seconds if not specified
	TimeoutSeconds int64 `json:"timeout_seconds"`
	// Initial browser window size in pixels with optional refresh rate. If omitted,
	// image defaults apply (1920x1080@25). Only specific viewport configurations are
	// supported. The server will reject unsupported combinations. Supported
	// resolutions are: 2560x1440@10, 1920x1080@25, 1920x1200@25, 1440x900@25,
	// 1024x768@60, 1200x800@60 If refresh_rate is not provided, it will be
	// automatically determined from the width and height if they match a supported
	// configuration exactly. Note: Higher resolutions may affect the responsiveness of
	// live view browser
	Viewport shared.BrowserViewport `json:"viewport"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Size              respjson.Field
		Extensions        respjson.Field
		FillRatePerMinute respjson.Field
		Headless          respjson.Field
		KioskMode         respjson.Field
		Name              respjson.Field
		Profile           respjson.Field
		ProxyID           respjson.Field
		Stealth           respjson.Field
		TimeoutSeconds    respjson.Field
		Viewport          respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPoolBrowserPoolConfig) RawJSON() string { return r.JSON.raw }
func (r *BrowserPoolBrowserPoolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPoolAcquireResponse struct {
	// Websocket URL for Chrome DevTools Protocol connections to the browser session
	CdpWsURL string `json:"cdp_ws_url,required"`
	// When the browser session was created.
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// Whether the browser session is running in headless mode.
	Headless bool `json:"headless,required"`
	// Unique identifier for the browser session
	SessionID string `json:"session_id,required"`
	// Whether the browser session is running in stealth mode.
	Stealth bool `json:"stealth,required"`
	// The number of seconds of inactivity before the browser session is terminated.
	TimeoutSeconds int64 `json:"timeout_seconds,required"`
	// Remote URL for live viewing the browser session. Only available for non-headless
	// browsers.
	BrowserLiveViewURL string `json:"browser_live_view_url"`
	// When the browser session was soft-deleted. Only present for deleted sessions.
	DeletedAt time.Time `json:"deleted_at" format:"date-time"`
	// Whether the browser session is running in kiosk mode.
	KioskMode bool `json:"kiosk_mode"`
	// DEPRECATED: Use timeout_seconds (up to 72 hours) and Profiles instead.
	//
	// Deprecated: deprecated
	Persistence BrowserPersistence `json:"persistence"`
	// Browser profile metadata.
	Profile Profile `json:"profile"`
	// ID of the proxy associated with this browser session, if any.
	ProxyID string `json:"proxy_id"`
	// Initial browser window size in pixels with optional refresh rate. If omitted,
	// image defaults apply (1920x1080@25). Only specific viewport configurations are
	// supported. The server will reject unsupported combinations. Supported
	// resolutions are: 2560x1440@10, 1920x1080@25, 1920x1200@25, 1440x900@25,
	// 1024x768@60, 1200x800@60 If refresh_rate is not provided, it will be
	// automatically determined from the width and height if they match a supported
	// configuration exactly. Note: Higher resolutions may affect the responsiveness of
	// live view browser
	Viewport shared.BrowserViewport `json:"viewport"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CdpWsURL           respjson.Field
		CreatedAt          respjson.Field
		Headless           respjson.Field
		SessionID          respjson.Field
		Stealth            respjson.Field
		TimeoutSeconds     respjson.Field
		BrowserLiveViewURL respjson.Field
		DeletedAt          respjson.Field
		KioskMode          respjson.Field
		Persistence        respjson.Field
		Profile            respjson.Field
		ProxyID            respjson.Field
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
	Size int64 `json:"size,required"`
	// Percentage of the pool to fill per minute. Defaults to 10%.
	FillRatePerMinute param.Opt[int64] `json:"fill_rate_per_minute,omitzero"`
	// If true, launches the browser using a headless image. Defaults to false.
	Headless param.Opt[bool] `json:"headless,omitzero"`
	// If true, launches the browser in kiosk mode to hide address bar and tabs in live
	// view.
	KioskMode param.Opt[bool] `json:"kiosk_mode,omitzero"`
	// Optional name for the browser pool. Must be unique within the organization.
	Name param.Opt[string] `json:"name,omitzero"`
	// Optional proxy to associate to the browser session. Must reference a proxy
	// belonging to the caller's org.
	ProxyID param.Opt[string] `json:"proxy_id,omitzero"`
	// If true, launches the browser in stealth mode to reduce detection by anti-bot
	// mechanisms.
	Stealth param.Opt[bool] `json:"stealth,omitzero"`
	// Default idle timeout in seconds for browsers acquired from this pool before they
	// are destroyed. Defaults to 600 seconds if not specified
	TimeoutSeconds param.Opt[int64] `json:"timeout_seconds,omitzero"`
	// List of browser extensions to load into the session. Provide each by id or name.
	Extensions []shared.BrowserExtensionParam `json:"extensions,omitzero"`
	// Profile selection for the browser session. Provide either id or name. If
	// specified, the matching profile will be loaded into the browser session.
	// Profiles must be created beforehand.
	Profile shared.BrowserProfileParam `json:"profile,omitzero"`
	// Initial browser window size in pixels with optional refresh rate. If omitted,
	// image defaults apply (1920x1080@25). Only specific viewport configurations are
	// supported. The server will reject unsupported combinations. Supported
	// resolutions are: 2560x1440@10, 1920x1080@25, 1920x1200@25, 1440x900@25,
	// 1024x768@60, 1200x800@60 If refresh_rate is not provided, it will be
	// automatically determined from the width and height if they match a supported
	// configuration exactly. Note: Higher resolutions may affect the responsiveness of
	// live view browser
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

type BrowserPoolUpdateParams struct {
	// Number of browsers to maintain in the pool. The maximum size is determined by
	// your organization's pooled sessions limit (the sum of all pool sizes cannot
	// exceed your limit).
	Size int64 `json:"size,required"`
	// Whether to discard all idle browsers and rebuild the pool immediately. Defaults
	// to false.
	DiscardAllIdle param.Opt[bool] `json:"discard_all_idle,omitzero"`
	// Percentage of the pool to fill per minute. Defaults to 10%.
	FillRatePerMinute param.Opt[int64] `json:"fill_rate_per_minute,omitzero"`
	// If true, launches the browser using a headless image. Defaults to false.
	Headless param.Opt[bool] `json:"headless,omitzero"`
	// If true, launches the browser in kiosk mode to hide address bar and tabs in live
	// view.
	KioskMode param.Opt[bool] `json:"kiosk_mode,omitzero"`
	// Optional name for the browser pool. Must be unique within the organization.
	Name param.Opt[string] `json:"name,omitzero"`
	// Optional proxy to associate to the browser session. Must reference a proxy
	// belonging to the caller's org.
	ProxyID param.Opt[string] `json:"proxy_id,omitzero"`
	// If true, launches the browser in stealth mode to reduce detection by anti-bot
	// mechanisms.
	Stealth param.Opt[bool] `json:"stealth,omitzero"`
	// Default idle timeout in seconds for browsers acquired from this pool before they
	// are destroyed. Defaults to 600 seconds if not specified
	TimeoutSeconds param.Opt[int64] `json:"timeout_seconds,omitzero"`
	// List of browser extensions to load into the session. Provide each by id or name.
	Extensions []shared.BrowserExtensionParam `json:"extensions,omitzero"`
	// Profile selection for the browser session. Provide either id or name. If
	// specified, the matching profile will be loaded into the browser session.
	// Profiles must be created beforehand.
	Profile shared.BrowserProfileParam `json:"profile,omitzero"`
	// Initial browser window size in pixels with optional refresh rate. If omitted,
	// image defaults apply (1920x1080@25). Only specific viewport configurations are
	// supported. The server will reject unsupported combinations. Supported
	// resolutions are: 2560x1440@10, 1920x1080@25, 1920x1200@25, 1440x900@25,
	// 1024x768@60, 1200x800@60 If refresh_rate is not provided, it will be
	// automatically determined from the width and height if they match a supported
	// configuration exactly. Note: Higher resolutions may affect the responsiveness of
	// live view browser
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
	paramObj
}

func (r BrowserPoolAcquireParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPoolAcquireParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPoolAcquireParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPoolReleaseParams struct {
	// Browser session ID to release back to the pool
	SessionID string `json:"session_id,required"`
	// Whether to reuse the browser instance or destroy it and create a new one.
	// Defaults to true.
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
