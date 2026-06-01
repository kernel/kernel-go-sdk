// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Control mouse, keyboard, and screen on the browser instance.
//
// BrowserComputerService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserComputerService] method instead.
type BrowserComputerService struct {
	Options []option.RequestOption
}

// NewBrowserComputerService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBrowserComputerService(opts ...option.RequestOption) (r BrowserComputerService) {
	r = BrowserComputerService{}
	r.Options = opts
	return
}

// Send an array of computer actions to execute in order on the browser instance.
// Execution stops on the first error. This reduces network latency compared to
// sending individual action requests.
func (r *BrowserComputerService) Batch(ctx context.Context, id string, body BrowserComputerBatchParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("browsers/%s/computer/batch", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

// Capture a screenshot of the browser instance
func (r *BrowserComputerService) CaptureScreenshot(ctx context.Context, id string, body BrowserComputerCaptureScreenshotParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "image/png")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/computer/screenshot", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Simulate a mouse click action on the browser instance
func (r *BrowserComputerService) ClickMouse(ctx context.Context, id string, body BrowserComputerClickMouseParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("browsers/%s/computer/click_mouse", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

// Drag the mouse along a path
func (r *BrowserComputerService) DragMouse(ctx context.Context, id string, body BrowserComputerDragMouseParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("browsers/%s/computer/drag_mouse", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

// Get the current mouse cursor position on the browser instance
func (r *BrowserComputerService) GetMousePosition(ctx context.Context, id string, opts ...option.RequestOption) (res *BrowserComputerGetMousePositionResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/computer/get_mouse_position", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Move the mouse cursor to the specified coordinates on the browser instance
func (r *BrowserComputerService) MoveMouse(ctx context.Context, id string, body BrowserComputerMoveMouseParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("browsers/%s/computer/move_mouse", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

// Press one or more keys on the host computer
func (r *BrowserComputerService) PressKey(ctx context.Context, id string, body BrowserComputerPressKeyParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("browsers/%s/computer/press_key", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

// Read text from the clipboard on the browser instance
func (r *BrowserComputerService) ReadClipboard(ctx context.Context, id string, opts ...option.RequestOption) (res *BrowserComputerReadClipboardResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/computer/clipboard/read", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Scroll the mouse wheel at a position on the host computer
func (r *BrowserComputerService) Scroll(ctx context.Context, id string, body BrowserComputerScrollParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("browsers/%s/computer/scroll", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

// Set cursor visibility
func (r *BrowserComputerService) SetCursorVisibility(ctx context.Context, id string, body BrowserComputerSetCursorVisibilityParams, opts ...option.RequestOption) (res *BrowserComputerSetCursorVisibilityResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/computer/cursor", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Type text on the browser instance
func (r *BrowserComputerService) TypeText(ctx context.Context, id string, body BrowserComputerTypeTextParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("browsers/%s/computer/type", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

// Write text to the clipboard on the browser instance
func (r *BrowserComputerService) WriteClipboard(ctx context.Context, id string, body BrowserComputerWriteClipboardParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("browsers/%s/computer/clipboard/write", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

type BrowserComputerGetMousePositionResponse struct {
	// X coordinate of the cursor
	X int64 `json:"x" api:"required"`
	// Y coordinate of the cursor
	Y int64 `json:"y" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		X           respjson.Field
		Y           respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserComputerGetMousePositionResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserComputerGetMousePositionResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserComputerReadClipboardResponse struct {
	// Current clipboard text content
	Text string `json:"text" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserComputerReadClipboardResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserComputerReadClipboardResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Generic OK response.
type BrowserComputerSetCursorVisibilityResponse struct {
	// Indicates success.
	Ok bool `json:"ok" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Ok          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserComputerSetCursorVisibilityResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserComputerSetCursorVisibilityResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserComputerBatchParams struct {
	// Ordered list of actions to execute. Execution stops on the first error.
	Actions []BrowserComputerBatchParamsAction `json:"actions,omitzero" api:"required"`
	paramObj
}

func (r BrowserComputerBatchParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerBatchParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerBatchParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A single computer action to execute as part of a batch. The `type` field selects
// which action to perform, and the corresponding field contains the action
// parameters. Exactly one action field matching the type must be provided.
//
// The property Type is required.
type BrowserComputerBatchParamsAction struct {
	// The type of action to perform.
	//
	// Any of "click_mouse", "move_mouse", "type_text", "press_key", "scroll",
	// "drag_mouse", "set_cursor", "sleep".
	Type       string                                     `json:"type,omitzero" api:"required"`
	ClickMouse BrowserComputerBatchParamsActionClickMouse `json:"click_mouse,omitzero"`
	DragMouse  BrowserComputerBatchParamsActionDragMouse  `json:"drag_mouse,omitzero"`
	MoveMouse  BrowserComputerBatchParamsActionMoveMouse  `json:"move_mouse,omitzero"`
	PressKey   BrowserComputerBatchParamsActionPressKey   `json:"press_key,omitzero"`
	Scroll     BrowserComputerBatchParamsActionScroll     `json:"scroll,omitzero"`
	SetCursor  BrowserComputerBatchParamsActionSetCursor  `json:"set_cursor,omitzero"`
	// Pause execution for a specified duration.
	Sleep    BrowserComputerBatchParamsActionSleep    `json:"sleep,omitzero"`
	TypeText BrowserComputerBatchParamsActionTypeText `json:"type_text,omitzero"`
	paramObj
}

func (r BrowserComputerBatchParamsAction) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerBatchParamsAction
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerBatchParamsAction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[BrowserComputerBatchParamsAction](
		"type", "click_mouse", "move_mouse", "type_text", "press_key", "scroll", "drag_mouse", "set_cursor", "sleep",
	)
}

// The properties X, Y are required.
type BrowserComputerBatchParamsActionClickMouse struct {
	// X coordinate of the click position
	X int64 `json:"x" api:"required"`
	// Y coordinate of the click position
	Y int64 `json:"y" api:"required"`
	// Number of times to repeat the click
	NumClicks param.Opt[int64] `json:"num_clicks,omitzero"`
	// Mouse button to interact with
	//
	// Any of "left", "right", "middle", "back", "forward".
	Button string `json:"button,omitzero"`
	// Type of click action
	//
	// Any of "down", "up", "click".
	ClickType string `json:"click_type,omitzero"`
	// Modifier keys to hold during the click
	HoldKeys []string `json:"hold_keys,omitzero"`
	paramObj
}

func (r BrowserComputerBatchParamsActionClickMouse) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerBatchParamsActionClickMouse
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerBatchParamsActionClickMouse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[BrowserComputerBatchParamsActionClickMouse](
		"button", "left", "right", "middle", "back", "forward",
	)
	apijson.RegisterFieldValidator[BrowserComputerBatchParamsActionClickMouse](
		"click_type", "down", "up", "click",
	)
}

// The property Path is required.
type BrowserComputerBatchParamsActionDragMouse struct {
	// Ordered list of [x, y] coordinate pairs to move through while dragging. Must
	// contain at least 2 points.
	Path [][]int64 `json:"path,omitzero" api:"required"`
	// Delay in milliseconds between button down and starting to move along the path.
	Delay param.Opt[int64] `json:"delay,omitzero"`
	// Target total duration in milliseconds for the entire drag movement when
	// smooth=true. Omit for automatic timing based on total path length.
	DurationMs param.Opt[int64] `json:"duration_ms,omitzero"`
	// Use human-like Bezier curves between path waypoints instead of linear
	// interpolation. When true, steps_per_segment and step_delay_ms are ignored.
	Smooth param.Opt[bool] `json:"smooth,omitzero"`
	// Delay in milliseconds between relative steps while dragging (not the initial
	// delay).
	StepDelayMs param.Opt[int64] `json:"step_delay_ms,omitzero"`
	// Number of relative move steps per segment in the path. Minimum 1.
	StepsPerSegment param.Opt[int64] `json:"steps_per_segment,omitzero"`
	// Mouse button to drag with
	//
	// Any of "left", "middle", "right".
	Button string `json:"button,omitzero"`
	// Modifier keys to hold during the drag
	HoldKeys []string `json:"hold_keys,omitzero"`
	paramObj
}

func (r BrowserComputerBatchParamsActionDragMouse) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerBatchParamsActionDragMouse
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerBatchParamsActionDragMouse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[BrowserComputerBatchParamsActionDragMouse](
		"button", "left", "middle", "right",
	)
}

// The properties X, Y are required.
type BrowserComputerBatchParamsActionMoveMouse struct {
	// X coordinate to move the cursor to
	X int64 `json:"x" api:"required"`
	// Y coordinate to move the cursor to
	Y int64 `json:"y" api:"required"`
	// Target total duration in milliseconds for the mouse movement when smooth=true.
	// Omit for automatic timing based on distance.
	DurationMs param.Opt[int64] `json:"duration_ms,omitzero"`
	// Use human-like Bezier curve path instead of instant mouse movement.
	Smooth param.Opt[bool] `json:"smooth,omitzero"`
	// Modifier keys to hold during the move
	HoldKeys []string `json:"hold_keys,omitzero"`
	paramObj
}

func (r BrowserComputerBatchParamsActionMoveMouse) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerBatchParamsActionMoveMouse
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerBatchParamsActionMoveMouse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Keys is required.
type BrowserComputerBatchParamsActionPressKey struct {
	// List of key symbols to press. Each item should be a key symbol supported by
	// xdotool (see X11 keysym definitions). Examples include "Return", "Shift",
	// "Ctrl", "Alt", "F5". Items in this list could also be combinations, e.g.
	// "Ctrl+t" or "Ctrl+Shift+Tab".
	Keys []string `json:"keys,omitzero" api:"required"`
	// Duration to hold the keys down in milliseconds. If omitted or 0, keys are
	// tapped.
	Duration param.Opt[int64] `json:"duration,omitzero"`
	// Optional modifier keys to hold during the key press sequence.
	HoldKeys []string `json:"hold_keys,omitzero"`
	paramObj
}

func (r BrowserComputerBatchParamsActionPressKey) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerBatchParamsActionPressKey
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerBatchParamsActionPressKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties X, Y are required.
type BrowserComputerBatchParamsActionScroll struct {
	// X coordinate at which to perform the scroll
	X int64 `json:"x" api:"required"`
	// Y coordinate at which to perform the scroll
	Y int64 `json:"y" api:"required"`
	// Horizontal scroll amount in xdotool "wheel units." Positive scrolls right,
	// negative scrolls left.
	DeltaX param.Opt[int64] `json:"delta_x,omitzero"`
	// Vertical scroll amount in xdotool "wheel units." Positive scrolls down, negative
	// scrolls up.
	DeltaY param.Opt[int64] `json:"delta_y,omitzero"`
	// Modifier keys to hold during the scroll
	HoldKeys []string `json:"hold_keys,omitzero"`
	paramObj
}

func (r BrowserComputerBatchParamsActionScroll) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerBatchParamsActionScroll
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerBatchParamsActionScroll) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Hidden is required.
type BrowserComputerBatchParamsActionSetCursor struct {
	// Whether the cursor should be hidden or visible
	Hidden bool `json:"hidden" api:"required"`
	paramObj
}

func (r BrowserComputerBatchParamsActionSetCursor) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerBatchParamsActionSetCursor
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerBatchParamsActionSetCursor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Pause execution for a specified duration.
//
// The property DurationMs is required.
type BrowserComputerBatchParamsActionSleep struct {
	// Duration to sleep in milliseconds.
	DurationMs int64 `json:"duration_ms" api:"required"`
	paramObj
}

func (r BrowserComputerBatchParamsActionSleep) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerBatchParamsActionSleep
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerBatchParamsActionSleep) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Text is required.
type BrowserComputerBatchParamsActionTypeText struct {
	// Text to type on the browser instance
	Text string `json:"text" api:"required"`
	// Delay in milliseconds between keystrokes
	Delay param.Opt[int64] `json:"delay,omitzero"`
	paramObj
}

func (r BrowserComputerBatchParamsActionTypeText) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerBatchParamsActionTypeText
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerBatchParamsActionTypeText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserComputerCaptureScreenshotParams struct {
	Region BrowserComputerCaptureScreenshotParamsRegion `json:"region,omitzero"`
	paramObj
}

func (r BrowserComputerCaptureScreenshotParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerCaptureScreenshotParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerCaptureScreenshotParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Height, Width, X, Y are required.
type BrowserComputerCaptureScreenshotParamsRegion struct {
	// Height of the region in pixels
	Height int64 `json:"height" api:"required"`
	// Width of the region in pixels
	Width int64 `json:"width" api:"required"`
	// X coordinate of the region's top-left corner
	X int64 `json:"x" api:"required"`
	// Y coordinate of the region's top-left corner
	Y int64 `json:"y" api:"required"`
	paramObj
}

func (r BrowserComputerCaptureScreenshotParamsRegion) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerCaptureScreenshotParamsRegion
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerCaptureScreenshotParamsRegion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserComputerClickMouseParams struct {
	// X coordinate of the click position
	X int64 `json:"x" api:"required"`
	// Y coordinate of the click position
	Y int64 `json:"y" api:"required"`
	// Number of times to repeat the click
	NumClicks param.Opt[int64] `json:"num_clicks,omitzero"`
	// Mouse button to interact with
	//
	// Any of "left", "right", "middle", "back", "forward".
	Button BrowserComputerClickMouseParamsButton `json:"button,omitzero"`
	// Type of click action
	//
	// Any of "down", "up", "click".
	ClickType BrowserComputerClickMouseParamsClickType `json:"click_type,omitzero"`
	// Modifier keys to hold during the click
	HoldKeys []string `json:"hold_keys,omitzero"`
	paramObj
}

func (r BrowserComputerClickMouseParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerClickMouseParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerClickMouseParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Mouse button to interact with
type BrowserComputerClickMouseParamsButton string

const (
	BrowserComputerClickMouseParamsButtonLeft    BrowserComputerClickMouseParamsButton = "left"
	BrowserComputerClickMouseParamsButtonRight   BrowserComputerClickMouseParamsButton = "right"
	BrowserComputerClickMouseParamsButtonMiddle  BrowserComputerClickMouseParamsButton = "middle"
	BrowserComputerClickMouseParamsButtonBack    BrowserComputerClickMouseParamsButton = "back"
	BrowserComputerClickMouseParamsButtonForward BrowserComputerClickMouseParamsButton = "forward"
)

// Type of click action
type BrowserComputerClickMouseParamsClickType string

const (
	BrowserComputerClickMouseParamsClickTypeDown  BrowserComputerClickMouseParamsClickType = "down"
	BrowserComputerClickMouseParamsClickTypeUp    BrowserComputerClickMouseParamsClickType = "up"
	BrowserComputerClickMouseParamsClickTypeClick BrowserComputerClickMouseParamsClickType = "click"
)

type BrowserComputerDragMouseParams struct {
	// Ordered list of [x, y] coordinate pairs to move through while dragging. Must
	// contain at least 2 points.
	Path [][]int64 `json:"path,omitzero" api:"required"`
	// Delay in milliseconds between button down and starting to move along the path.
	Delay param.Opt[int64] `json:"delay,omitzero"`
	// Target total duration in milliseconds for the entire drag movement when
	// smooth=true. Omit for automatic timing based on total path length.
	DurationMs param.Opt[int64] `json:"duration_ms,omitzero"`
	// Use human-like Bezier curves between path waypoints instead of linear
	// interpolation. When true, steps_per_segment and step_delay_ms are ignored.
	Smooth param.Opt[bool] `json:"smooth,omitzero"`
	// Delay in milliseconds between relative steps while dragging (not the initial
	// delay).
	StepDelayMs param.Opt[int64] `json:"step_delay_ms,omitzero"`
	// Number of relative move steps per segment in the path. Minimum 1.
	StepsPerSegment param.Opt[int64] `json:"steps_per_segment,omitzero"`
	// Mouse button to drag with
	//
	// Any of "left", "middle", "right".
	Button BrowserComputerDragMouseParamsButton `json:"button,omitzero"`
	// Modifier keys to hold during the drag
	HoldKeys []string `json:"hold_keys,omitzero"`
	paramObj
}

func (r BrowserComputerDragMouseParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerDragMouseParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerDragMouseParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Mouse button to drag with
type BrowserComputerDragMouseParamsButton string

const (
	BrowserComputerDragMouseParamsButtonLeft   BrowserComputerDragMouseParamsButton = "left"
	BrowserComputerDragMouseParamsButtonMiddle BrowserComputerDragMouseParamsButton = "middle"
	BrowserComputerDragMouseParamsButtonRight  BrowserComputerDragMouseParamsButton = "right"
)

type BrowserComputerMoveMouseParams struct {
	// X coordinate to move the cursor to
	X int64 `json:"x" api:"required"`
	// Y coordinate to move the cursor to
	Y int64 `json:"y" api:"required"`
	// Target total duration in milliseconds for the mouse movement when smooth=true.
	// Omit for automatic timing based on distance.
	DurationMs param.Opt[int64] `json:"duration_ms,omitzero"`
	// Use human-like Bezier curve path instead of instant mouse movement.
	Smooth param.Opt[bool] `json:"smooth,omitzero"`
	// Modifier keys to hold during the move
	HoldKeys []string `json:"hold_keys,omitzero"`
	paramObj
}

func (r BrowserComputerMoveMouseParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerMoveMouseParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerMoveMouseParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserComputerPressKeyParams struct {
	// List of key symbols to press. Each item should be a key symbol supported by
	// xdotool (see X11 keysym definitions). Examples include "Return", "Shift",
	// "Ctrl", "Alt", "F5". Items in this list could also be combinations, e.g.
	// "Ctrl+t" or "Ctrl+Shift+Tab".
	Keys []string `json:"keys,omitzero" api:"required"`
	// Duration to hold the keys down in milliseconds. If omitted or 0, keys are
	// tapped.
	Duration param.Opt[int64] `json:"duration,omitzero"`
	// Optional modifier keys to hold during the key press sequence.
	HoldKeys []string `json:"hold_keys,omitzero"`
	paramObj
}

func (r BrowserComputerPressKeyParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerPressKeyParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerPressKeyParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserComputerScrollParams struct {
	// X coordinate at which to perform the scroll
	X int64 `json:"x" api:"required"`
	// Y coordinate at which to perform the scroll
	Y int64 `json:"y" api:"required"`
	// Horizontal scroll amount in xdotool "wheel units." Positive scrolls right,
	// negative scrolls left.
	DeltaX param.Opt[int64] `json:"delta_x,omitzero"`
	// Vertical scroll amount in xdotool "wheel units." Positive scrolls down, negative
	// scrolls up.
	DeltaY param.Opt[int64] `json:"delta_y,omitzero"`
	// Modifier keys to hold during the scroll
	HoldKeys []string `json:"hold_keys,omitzero"`
	paramObj
}

func (r BrowserComputerScrollParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerScrollParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerScrollParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserComputerSetCursorVisibilityParams struct {
	// Whether the cursor should be hidden or visible
	Hidden bool `json:"hidden" api:"required"`
	paramObj
}

func (r BrowserComputerSetCursorVisibilityParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerSetCursorVisibilityParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerSetCursorVisibilityParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserComputerTypeTextParams struct {
	// Text to type on the browser instance
	Text string `json:"text" api:"required"`
	// Delay in milliseconds between keystrokes
	Delay param.Opt[int64] `json:"delay,omitzero"`
	paramObj
}

func (r BrowserComputerTypeTextParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerTypeTextParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerTypeTextParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserComputerWriteClipboardParams struct {
	// Text to write to the system clipboard
	Text string `json:"text" api:"required"`
	paramObj
}

func (r BrowserComputerWriteClipboardParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserComputerWriteClipboardParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserComputerWriteClipboardParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
