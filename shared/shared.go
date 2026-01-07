// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

import (
	"encoding/json"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
	"github.com/kernel/kernel-go-sdk/shared/constant"
)

// aliased to make [param.APIUnion] private when embedding
type paramUnion = param.APIUnion

// aliased to make [param.APIObject] private when embedding
type paramObj = param.APIObject

// An action available on the app
type AppAction struct {
	// Name of the action
	Name string `json:"name,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AppAction) RawJSON() string { return r.JSON.raw }
func (r *AppAction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Extension selection for the browser session. Provide either id or name of an
// extension uploaded to Kernel.
type BrowserExtension struct {
	// Extension ID to load for this browser session
	ID string `json:"id"`
	// Extension name to load for this browser session (instead of id). Must be 1-255
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
func (r BrowserExtension) RawJSON() string { return r.JSON.raw }
func (r *BrowserExtension) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BrowserExtension to a BrowserExtensionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BrowserExtensionParam.Overrides()
func (r BrowserExtension) ToParam() BrowserExtensionParam {
	return param.Override[BrowserExtensionParam](json.RawMessage(r.RawJSON()))
}

// Extension selection for the browser session. Provide either id or name of an
// extension uploaded to Kernel.
type BrowserExtensionParam struct {
	// Extension ID to load for this browser session
	ID param.Opt[string] `json:"id,omitzero"`
	// Extension name to load for this browser session (instead of id). Must be 1-255
	// characters, using letters, numbers, dots, underscores, or hyphens.
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r BrowserExtensionParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserExtensionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserExtensionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Profile selection for the browser session. Provide either id or name. If
// specified, the matching profile will be loaded into the browser session.
// Profiles must be created beforehand.
type BrowserProfile struct {
	// Profile ID to load for this browser session
	ID string `json:"id"`
	// Profile name to load for this browser session (instead of id). Must be 1-255
	// characters, using letters, numbers, dots, underscores, or hyphens.
	Name string `json:"name"`
	// If true, save changes made during the session back to the profile when the
	// session ends.
	SaveChanges bool `json:"save_changes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		SaveChanges respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserProfile) RawJSON() string { return r.JSON.raw }
func (r *BrowserProfile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BrowserProfile to a BrowserProfileParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BrowserProfileParam.Overrides()
func (r BrowserProfile) ToParam() BrowserProfileParam {
	return param.Override[BrowserProfileParam](json.RawMessage(r.RawJSON()))
}

// Profile selection for the browser session. Provide either id or name. If
// specified, the matching profile will be loaded into the browser session.
// Profiles must be created beforehand.
type BrowserProfileParam struct {
	// Profile ID to load for this browser session
	ID param.Opt[string] `json:"id,omitzero"`
	// Profile name to load for this browser session (instead of id). Must be 1-255
	// characters, using letters, numbers, dots, underscores, or hyphens.
	Name param.Opt[string] `json:"name,omitzero"`
	// If true, save changes made during the session back to the profile when the
	// session ends.
	SaveChanges param.Opt[bool] `json:"save_changes,omitzero"`
	paramObj
}

func (r BrowserProfileParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserProfileParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserProfileParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Initial browser window size in pixels with optional refresh rate. If omitted,
// image defaults apply (1920x1080@25). Only specific viewport configurations are
// supported. The server will reject unsupported combinations. Supported
// resolutions are: 2560x1440@10, 1920x1080@25, 1920x1200@25, 1440x900@25,
// 1024x768@60, 1200x800@60 If refresh_rate is not provided, it will be
// automatically determined from the width and height if they match a supported
// configuration exactly. Note: Higher resolutions may affect the responsiveness of
// live view browser
type BrowserViewport struct {
	// Browser window height in pixels.
	Height int64 `json:"height,required"`
	// Browser window width in pixels.
	Width int64 `json:"width,required"`
	// Display refresh rate in Hz. If omitted, automatically determined from width and
	// height.
	RefreshRate int64 `json:"refresh_rate"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Height      respjson.Field
		Width       respjson.Field
		RefreshRate respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserViewport) RawJSON() string { return r.JSON.raw }
func (r *BrowserViewport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BrowserViewport to a BrowserViewportParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BrowserViewportParam.Overrides()
func (r BrowserViewport) ToParam() BrowserViewportParam {
	return param.Override[BrowserViewportParam](json.RawMessage(r.RawJSON()))
}

// Initial browser window size in pixels with optional refresh rate. If omitted,
// image defaults apply (1920x1080@25). Only specific viewport configurations are
// supported. The server will reject unsupported combinations. Supported
// resolutions are: 2560x1440@10, 1920x1080@25, 1920x1200@25, 1440x900@25,
// 1024x768@60, 1200x800@60 If refresh_rate is not provided, it will be
// automatically determined from the width and height if they match a supported
// configuration exactly. Note: Higher resolutions may affect the responsiveness of
// live view browser
//
// The properties Height, Width are required.
type BrowserViewportParam struct {
	// Browser window height in pixels.
	Height int64 `json:"height,required"`
	// Browser window width in pixels.
	Width int64 `json:"width,required"`
	// Display refresh rate in Hz. If omitted, automatically determined from width and
	// height.
	RefreshRate param.Opt[int64] `json:"refresh_rate,omitzero"`
	paramObj
}

func (r BrowserViewportParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserViewportParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserViewportParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ErrorDetail struct {
	// Lower-level error code providing more specific detail
	Code string `json:"code"`
	// Further detail about the error
	Message string `json:"message"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ErrorDetail) RawJSON() string { return r.JSON.raw }
func (r *ErrorDetail) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An error event from the application.
type ErrorEvent struct {
	Error ErrorModel `json:"error,required"`
	// Event type identifier (always "error").
	Event constant.Error `json:"event,required"`
	// Time the error occurred.
	Timestamp time.Time `json:"timestamp,required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Error       respjson.Field
		Event       respjson.Field
		Timestamp   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ErrorEvent) RawJSON() string { return r.JSON.raw }
func (r *ErrorEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ErrorEvent) ImplInvocationFollowResponseUnion() {}

type ErrorModel struct {
	// Application-specific error code (machine-readable)
	Code string `json:"code,required"`
	// Human-readable error description for debugging
	Message string `json:"message,required"`
	// Additional error details (for multiple errors)
	Details    []ErrorDetail `json:"details"`
	InnerError ErrorDetail   `json:"inner_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		Details     respjson.Field
		InnerError  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ErrorModel) RawJSON() string { return r.JSON.raw }
func (r *ErrorModel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Heartbeat event sent periodically to keep SSE connection alive.
type HeartbeatEvent struct {
	// Event type identifier (always "sse_heartbeat").
	Event constant.SseHeartbeat `json:"event,required"`
	// Time the heartbeat was sent.
	Timestamp time.Time `json:"timestamp,required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Event       respjson.Field
		Timestamp   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r HeartbeatEvent) RawJSON() string { return r.JSON.raw }
func (r *HeartbeatEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (HeartbeatEvent) ImplInvocationFollowResponseUnion() {}

// A log entry from the application.
type LogEvent struct {
	// Event type identifier (always "log").
	Event constant.Log `json:"event,required"`
	// Log message text.
	Message string `json:"message,required"`
	// Time the log entry was produced.
	Timestamp time.Time `json:"timestamp,required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Event       respjson.Field
		Message     respjson.Field
		Timestamp   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LogEvent) RawJSON() string { return r.JSON.raw }
func (r *LogEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (LogEvent) ImplInvocationFollowResponseUnion() {}
