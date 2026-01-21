// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apiform"
	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
	"github.com/kernel/kernel-go-sdk/shared"
)

// BrowserService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserService] method instead.
type BrowserService struct {
	Options    []option.RequestOption
	Replays    BrowserReplayService
	Fs         BrowserFService
	Process    BrowserProcessService
	Logs       BrowserLogService
	Computer   BrowserComputerService
	Playwright BrowserPlaywrightService
}

// NewBrowserService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewBrowserService(opts ...option.RequestOption) (r BrowserService) {
	r = BrowserService{}
	r.Options = opts
	r.Replays = NewBrowserReplayService(opts...)
	r.Fs = NewBrowserFService(opts...)
	r.Process = NewBrowserProcessService(opts...)
	r.Logs = NewBrowserLogService(opts...)
	r.Computer = NewBrowserComputerService(opts...)
	r.Playwright = NewBrowserPlaywrightService(opts...)
	return
}

// Create a new browser session from within an action.
func (r *BrowserService) New(ctx context.Context, body BrowserNewParams, opts ...option.RequestOption) (res *BrowserNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "browsers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get information about a browser session.
func (r *BrowserService) Get(ctx context.Context, id string, query BrowserGetParams, opts ...option.RequestOption) (res *BrowserGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Update a browser session.
func (r *BrowserService) Update(ctx context.Context, id string, body BrowserUpdateParams, opts ...option.RequestOption) (res *BrowserUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// List all browser sessions with pagination support. Use status parameter to
// filter by session state.
func (r *BrowserService) List(ctx context.Context, query BrowserListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[BrowserListResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "browsers"
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

// List all browser sessions with pagination support. Use status parameter to
// filter by session state.
func (r *BrowserService) ListAutoPaging(ctx context.Context, query BrowserListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[BrowserListResponse] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// DEPRECATED: Use DELETE /browsers/{id} instead. Delete a persistent browser
// session by its persistent_id.
//
// Deprecated: deprecated
func (r *BrowserService) Delete(ctx context.Context, body BrowserDeleteParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	path := "browsers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, nil, opts...)
	return
}

// Delete a browser session by ID
func (r *BrowserService) DeleteByID(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return
}

// Loads one or more unpacked extensions and restarts Chromium on the browser
// instance.
func (r *BrowserService) LoadExtensions(ctx context.Context, id string, body BrowserLoadExtensionsParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/extensions", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return
}

// DEPRECATED: Use timeout_seconds (up to 72 hours) and Profiles instead.
//
// Deprecated: deprecated
type BrowserPersistence struct {
	// DEPRECATED: Unique identifier for the persistent browser session.
	ID string `json:"id,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPersistence) RawJSON() string { return r.JSON.raw }
func (r *BrowserPersistence) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BrowserPersistence to a BrowserPersistenceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BrowserPersistenceParam.Overrides()
func (r BrowserPersistence) ToParam() BrowserPersistenceParam {
	return param.Override[BrowserPersistenceParam](json.RawMessage(r.RawJSON()))
}

// DEPRECATED: Use timeout_seconds (up to 72 hours) and Profiles instead.
//
// Deprecated: deprecated
//
// The property ID is required.
type BrowserPersistenceParam struct {
	// DEPRECATED: Unique identifier for the persistent browser session.
	ID string `json:"id,required"`
	paramObj
}

func (r BrowserPersistenceParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserPersistenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserPersistenceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser profile metadata.
type Profile struct {
	// Unique identifier for the profile
	ID string `json:"id,required"`
	// Timestamp when the profile was created
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// Timestamp when the profile was last used
	LastUsedAt time.Time `json:"last_used_at" format:"date-time"`
	// Optional, easier-to-reference name for the profile
	Name string `json:"name,nullable"`
	// Timestamp when the profile was last updated
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		LastUsedAt  respjson.Field
		Name        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Profile) RawJSON() string { return r.JSON.raw }
func (r *Profile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserNewResponse struct {
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
func (r BrowserNewResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserGetResponse struct {
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
func (r BrowserGetResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserUpdateResponse struct {
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
func (r BrowserUpdateResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserUpdateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserListResponse struct {
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
func (r BrowserListResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserNewParams struct {
	// If true, launches the browser using a headless image (no VNC/GUI). Defaults to
	// false.
	Headless param.Opt[bool] `json:"headless,omitzero"`
	// action invocation ID
	InvocationID param.Opt[string] `json:"invocation_id,omitzero"`
	// If true, launches the browser in kiosk mode to hide address bar and tabs in live
	// view.
	KioskMode param.Opt[bool] `json:"kiosk_mode,omitzero"`
	// Optional proxy to associate to the browser session. Must reference a proxy
	// belonging to the caller's org.
	ProxyID param.Opt[string] `json:"proxy_id,omitzero"`
	// If true, launches the browser in stealth mode to reduce detection by anti-bot
	// mechanisms.
	Stealth param.Opt[bool] `json:"stealth,omitzero"`
	// The number of seconds of inactivity before the browser session is terminated.
	// Activity includes CDP connections and live view connections. Defaults to 60
	// seconds. Minimum allowed is 10 seconds. Maximum allowed is 259200 (72 hours). We
	// check for inactivity every 5 seconds, so the actual timeout behavior you will
	// see is +/- 5 seconds around the specified value.
	TimeoutSeconds param.Opt[int64] `json:"timeout_seconds,omitzero"`
	// List of browser extensions to load into the session. Provide each by id or name.
	Extensions []shared.BrowserExtensionParam `json:"extensions,omitzero"`
	// DEPRECATED: Use timeout_seconds (up to 72 hours) and Profiles instead.
	Persistence BrowserPersistenceParam `json:"persistence,omitzero"`
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

func (r BrowserNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserGetParams struct {
	// When true, includes soft-deleted browser sessions in the lookup.
	IncludeDeleted param.Opt[bool] `query:"include_deleted,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserGetParams]'s query parameters as `url.Values`.
func (r BrowserGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BrowserUpdateParams struct {
	// ID of the proxy to use. Omit to leave unchanged, set to empty string to remove
	// proxy.
	ProxyID param.Opt[string] `json:"proxy_id,omitzero"`
	paramObj
}

func (r BrowserUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserListParams struct {
	// Deprecated: Use status=all instead. When true, includes soft-deleted browser
	// sessions in the results alongside active sessions.
	IncludeDeleted param.Opt[bool] `query:"include_deleted,omitzero" json:"-"`
	// Maximum number of results to return. Defaults to 20, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Number of results to skip. Defaults to 0.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Filter sessions by status. "active" returns only active sessions (default),
	// "deleted" returns only soft-deleted sessions, "all" returns both.
	//
	// Any of "active", "deleted", "all".
	Status BrowserListParamsStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserListParams]'s query parameters as `url.Values`.
func (r BrowserListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter sessions by status. "active" returns only active sessions (default),
// "deleted" returns only soft-deleted sessions, "all" returns both.
type BrowserListParamsStatus string

const (
	BrowserListParamsStatusActive  BrowserListParamsStatus = "active"
	BrowserListParamsStatusDeleted BrowserListParamsStatus = "deleted"
	BrowserListParamsStatusAll     BrowserListParamsStatus = "all"
)

type BrowserDeleteParams struct {
	// Persistent browser identifier
	PersistentID string `query:"persistent_id,required" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserDeleteParams]'s query parameters as `url.Values`.
func (r BrowserDeleteParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BrowserLoadExtensionsParams struct {
	// List of extensions to upload and activate
	Extensions []BrowserLoadExtensionsParamsExtension `json:"extensions,omitzero,required"`
	paramObj
}

func (r BrowserLoadExtensionsParams) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err == nil {
		err = apiform.WriteExtras(writer, r.ExtraFields())
	}
	if err != nil {
		writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

// The properties Name, ZipFile are required.
type BrowserLoadExtensionsParamsExtension struct {
	// Folder name to place the extension under /home/kernel/extensions/<name>
	Name string `json:"name,required"`
	// Zip archive containing an unpacked Chromium extension (must include
	// manifest.json)
	ZipFile io.Reader `json:"zip_file,omitzero,required" format:"binary"`
	paramObj
}

func (r BrowserLoadExtensionsParamsExtension) MarshalJSON() (data []byte, err error) {
	type shadow BrowserLoadExtensionsParamsExtension
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserLoadExtensionsParamsExtension) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
