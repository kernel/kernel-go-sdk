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

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
	"github.com/kernel/kernel-go-sdk/packages/ssestream"
	"github.com/kernel/kernel-go-sdk/shared/constant"
)

// Stream live telemetry events from a browser session, and manage the destinations
// sessions export them to.
//
// BrowserTelemetryService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserTelemetryService] method instead.
type BrowserTelemetryService struct {
	Options []option.RequestOption
}

// NewBrowserTelemetryService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBrowserTelemetryService(opts ...option.RequestOption) (r BrowserTelemetryService) {
	r = BrowserTelemetryService{}
	r.Options = opts
	return
}

// Reads a page of telemetry events for the browser session. To page through
// results, pass the X-Next-Offset value from the previous response as offset and
// repeat while X-Has-More is true. Returns an empty list when telemetry data is
// unavailable.
func (r *BrowserTelemetryService) Events(ctx context.Context, id string, query BrowserTelemetryEventsParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[BrowserTelemetryEventsResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/telemetry/events", id)
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

// Reads a page of telemetry events for the browser session. To page through
// results, pass the X-Next-Offset value from the previous response as offset and
// repeat while X-Has-More is true. Returns an empty list when telemetry data is
// unavailable.
func (r *BrowserTelemetryService) EventsAutoPaging(ctx context.Context, id string, query BrowserTelemetryEventsParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[BrowserTelemetryEventsResponse] {
	return pagination.NewOffsetPaginationAutoPager(r.Events(ctx, id, query, opts...))
}

// Streams browser telemetry events as a server-sent events (SSE) stream. The
// stream closes when the browser session terminates. Each event frame includes an
// id: field containing a monotonically increasing sequence number; pass it as
// Last-Event-ID on reconnect to resume without gaps. The event: field is never
// set; all frames carry JSON in the data: field. A keepalive comment frame is sent
// every 15 seconds when no events arrive. Returns 404 if the browser session does
// not exist. If telemetry was not enabled on the session, the stream opens but no
// events are delivered. Fresh connections only see new events; pass replay=all to
// start from the oldest retained event instead.
func (r *BrowserTelemetryService) StreamStreaming(ctx context.Context, id string, params BrowserTelemetryStreamParams, opts ...option.RequestOption) (stream *ssestream.Stream[BrowserTelemetryStreamResponse]) {
	var (
		raw *http.Response
		err error
	)
	if !param.IsOmitted(params.LastEventID) {
		opts = append(opts, option.WithHeader("Last-Event-ID", fmt.Sprintf("%v", params.LastEventID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "text/event-stream")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return ssestream.NewStream[BrowserTelemetryStreamResponse](nil, err)
	}
	path := fmt.Sprintf("browsers/%s/telemetry/stream", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &raw, opts...)
	return ssestream.NewStream[BrowserTelemetryStreamResponse](ssestream.NewDecoder(raw), err)
}

// An agent-driven HTTP call that drives the browser, handled by the in-VM API
// server. Calls that manage the VM instead emit platform_api_call.
type BrowserAPICallEvent struct {
	Category constant.Control `json:"category" default:"control"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                   `json:"ts" api:"required"`
	Type constant.APICall        `json:"type" default:"api_call"`
	Data BrowserAPICallEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserAPICallEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserAPICallEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserAPICallEventData struct {
	// Wall-clock duration of the handler in milliseconds.
	DurationMs float64 `json:"duration_ms" api:"required"`
	// Matched route's operation, named as the in-VM API names its handler (e.g.
	// ProcessExec, TakeScreenshot).
	OperationID string `json:"operation_id" api:"required"`
	// Per-request identifier from the in-VM API request middleware.
	RequestID string `json:"request_id" api:"required"`
	// HTTP response status code.
	Status int64 `json:"status" api:"required"`
	// Source submitted to the Playwright code-execution endpoint, capped at 8192 bytes
	// like every other captured string. A capped value is cut on a character boundary
	// and ends in `...[truncated]`. Absent for every other operation.
	Code string `json:"code"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DurationMs  respjson.Field
		OperationID respjson.Field
		RequestID   respjson.Field
		Status      respjson.Field
		Code        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserAPICallEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserAPICallEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CDP Runtime.StackTrace representing the JavaScript call stack at the time of an
// event. Fields use CDP naming conventions rather than snake_case to match the
// Chrome DevTools Protocol wire format.
type BrowserCallStack struct {
	// Ordered list of call frames, outermost first.
	CallFrames []BrowserCallStackCallFrame `json:"callFrames" api:"required"`
	// Optional label for the stack trace (e.g. async cause).
	Description string `json:"description"`
	// Parent stack trace for async stacks.
	Parent *BrowserCallStack `json:"parent"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CallFrames  respjson.Field
		Description respjson.Field
		Parent      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCallStack) RawJSON() string { return r.JSON.raw }
func (r *BrowserCallStack) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserCallStackCallFrame struct {
	// Zero-based column number within the line.
	ColumnNumber int64 `json:"columnNumber" api:"required"`
	// JavaScript function name, or empty string for anonymous functions.
	FunctionName string `json:"functionName" api:"required"`
	// Zero-based line number within the script.
	LineNumber int64 `json:"lineNumber" api:"required"`
	// CDP script identifier.
	ScriptID string `json:"scriptId" api:"required"`
	// URL or name of the script file.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ColumnNumber respjson.Field
		FunctionName respjson.Field
		LineNumber   respjson.Field
		ScriptID     respjson.Field
		URL          respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCallStackCallFrame) RawJSON() string { return r.JSON.raw }
func (r *BrowserCallStackCallFrame) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A visible captcha challenge reached a terminal outcome.
type BrowserCaptchaChallengeResultEvent struct {
	Category constant.Captcha `json:"category" default:"captcha"`
	// Per-challenge payload. This event is emitted once per challenge and determines
	// its overall outcome; captcha_solve_started and captcha_solve_result describe
	// individual tasks and may occur multiple times within the challenge.
	Data BrowserCaptchaChallengeResultEventData `json:"data" api:"required"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                           `json:"ts" api:"required"`
	Type constant.CaptchaChallengeResult `json:"type" default:"captcha_challenge_result"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Data        respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCaptchaChallengeResultEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserCaptchaChallengeResultEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-challenge payload. This event is emitted once per challenge and determines
// its overall outcome; captcha_solve_started and captcha_solve_result describe
// individual tasks and may occur multiple times within the challenge.
type BrowserCaptchaChallengeResultEventData struct {
	// Captcha kind. Enterprise reCAPTCHA variants are grouped into their version
	// bucket (recaptcha_v2 or recaptcha_v3), press-and-hold challenges use
	// press_and_hold, and unlisted kinds use other.
	//
	// Any of "hcaptcha", "recaptcha_v2", "recaptcha_v3", "turnstile", "geetest",
	// "press_and_hold", "other".
	CaptchaType string `json:"captcha_type" api:"required"`
	// Opaque identifier shared by events for one visible challenge. An image-grid
	// captcha may create multiple task_id values for one challenge_id. The same value
	// may continue across a page reload when the challenge episode continues. It does
	// not indicate task ordering or challenge completion.
	ChallengeID string `json:"challenge_id" api:"required"`
	// Wall-clock duration from the challenge appearing to its terminal outcome,
	// covering every solver attempt in between.
	DurationMs float64 `json:"duration_ms" api:"required"`
	// Terminal outcome of the visible challenge. solved: the page observed the
	// challenge clear after a solver attempt. failure: a terminal solver failure
	// occurred, or all attempts ended while the challenge remained. timeout: the
	// challenge-level wait budget expired while the challenge remained. abandoned:
	// observation ended without an attributable terminal challenge outcome. This
	// includes a dismissed widget or page unload without a solved signal or terminal
	// solver outcome, and a token appearing while multiple same-provider challenges
	// are open, because the producer cannot attribute that token to this visible
	// challenge. A captcha_solve_result with the same challenge_id may therefore
	// report success while the challenge result reports abandoned. A solved challenge
	// does not prove the site accepted the token or that the guarded action succeeded.
	//
	// Any of "solved", "failure", "timeout", "abandoned".
	Status string `json:"status" api:"required"`
	// Host of the page where the challenge appeared.
	WebsiteHost string `json:"website_host"`
	// Path of the page where the challenge appeared. Query string excluded.
	WebsitePath string `json:"website_path"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CaptchaType respjson.Field
		ChallengeID respjson.Field
		DurationMs  respjson.Field
		Status      respjson.Field
		WebsiteHost respjson.Field
		WebsitePath respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCaptchaChallengeResultEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserCaptchaChallengeResultEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A captcha solve attempt reached a terminal outcome.
type BrowserCaptchaSolveResultEvent struct {
	Category constant.Captcha `json:"category" default:"captcha"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                              `json:"ts" api:"required"`
	Type constant.CaptchaSolveResult        `json:"type" default:"captcha_solve_result"`
	Data BrowserCaptchaSolveResultEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCaptchaSolveResultEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserCaptchaSolveResultEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserCaptchaSolveResultEventData struct {
	// Captcha kind. Enterprise reCAPTCHA variants are grouped into their version
	// bucket (recaptcha_v2 or recaptcha_v3), press-and-hold challenges use
	// press_and_hold, and unlisted kinds use other.
	//
	// Any of "hcaptcha", "recaptcha_v2", "recaptcha_v3", "turnstile", "geetest",
	// "press_and_hold", "other".
	CaptchaType string `json:"captcha_type" api:"required"`
	// Wall-clock duration from solve start to terminal outcome. Authoritative solve
	// timing; do not derive it from the gap to a captcha_solve_started event, whose
	// delivery and ordering are not guaranteed.
	DurationMs float64 `json:"duration_ms" api:"required"`
	// Terminal outcome. success: solver returned a usable solution. failure: solver
	// returned an error (see error_code). timeout: solver did not return within the
	// caller's wait budget. abandoned: caller cancelled or the page navigated away
	// mid-solve.
	//
	// Any of "success", "failure", "timeout", "abandoned".
	Status string `json:"status" api:"required"`
	// Opaque identifier shared by events for one visible challenge. An image-grid
	// captcha may create multiple task_id values for one challenge_id. The same value
	// may continue across a page reload when the challenge episode continues. It does
	// not indicate task ordering or challenge completion.
	ChallengeID string `json:"challenge_id"`
	// Solver-specific error code on failure (e.g. ERROR_CAPTCHA_UNSOLVABLE). Absent on
	// success.
	ErrorCode string `json:"error_code"`
	// Opaque identifier shared with the matching captcha_solve_started.
	TaskID string `json:"task_id"`
	// Host of the page where the captcha was solved.
	WebsiteHost string `json:"website_host"`
	// Path of the page where the captcha was solved. Query string excluded.
	WebsitePath string `json:"website_path"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CaptchaType respjson.Field
		DurationMs  respjson.Field
		Status      respjson.Field
		ChallengeID respjson.Field
		ErrorCode   respjson.Field
		TaskID      respjson.Field
		WebsiteHost respjson.Field
		WebsitePath respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCaptchaSolveResultEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserCaptchaSolveResultEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A captcha solver accepted a task.
type BrowserCaptchaSolveStartedEvent struct {
	Category constant.Captcha `json:"category" default:"captcha"`
	// Per-task payload. A visible challenge may create multiple tasks. When present,
	// task_id correlates this event with a captcha_solve_result, while challenge_id
	// groups tasks from the same challenge. Events may arrive out of order or be
	// absent, so their arrival does not indicate current solve state.
	Data BrowserCaptchaSolveStartedEventData `json:"data" api:"required"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                        `json:"ts" api:"required"`
	Type constant.CaptchaSolveStarted `json:"type" default:"captcha_solve_started"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Data        respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCaptchaSolveStartedEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserCaptchaSolveStartedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-task payload. A visible challenge may create multiple tasks. When present,
// task_id correlates this event with a captcha_solve_result, while challenge_id
// groups tasks from the same challenge. Events may arrive out of order or be
// absent, so their arrival does not indicate current solve state.
type BrowserCaptchaSolveStartedEventData struct {
	// Captcha kind. Enterprise reCAPTCHA variants are grouped into their version
	// bucket (recaptcha_v2 or recaptcha_v3), press-and-hold challenges use
	// press_and_hold, and unlisted kinds use other.
	//
	// Any of "hcaptcha", "recaptcha_v2", "recaptcha_v3", "turnstile", "geetest",
	// "press_and_hold", "other".
	CaptchaType string `json:"captcha_type" api:"required"`
	// Opaque identifier shared by events for one visible challenge. An image-grid
	// captcha may create multiple task_id values for one challenge_id. The same value
	// may continue across a page reload when the challenge episode continues. It does
	// not indicate task ordering or challenge completion.
	ChallengeID string `json:"challenge_id"`
	// Opaque identifier shared with the matching captcha_solve_result.
	TaskID string `json:"task_id"`
	// Host of the page where the captcha is being solved. May be empty for solver
	// tasks that carry no page URL.
	WebsiteHost string `json:"website_host"`
	// Path of the page where the captcha is being solved. Query string excluded.
	WebsitePath string `json:"website_path"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CaptchaType respjson.Field
		ChallengeID respjson.Field
		TaskID      respjson.Field
		WebsiteHost respjson.Field
		WebsitePath respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCaptchaSolveStartedEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserCaptchaSolveStartedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser-control command a client sent over the CDP WebSocket proxy: input
// gestures, navigation, dialog handling, file selection and screenshots.
// Configuration commands and the DOM/Runtime traffic a client library issues on
// the caller's behalf are not reported. One event per browser-control command that
// reached the browser. The command stream is not sampled, coalesced or reordered.
// An event is lost only when the method is excluded by telemetry configuration,
// when the command's arguments do not decode, or when classification cannot keep
// up. Exclusions are counted in `cdp_disconnect.telemetry_excluded`; the rest in
// `cdp_disconnect.telemetry_dropped`.
type BrowserCdpCommandEvent struct {
	Category constant.Control `json:"category" default:"control"`
	// Per-command payload for `cdp_command` events, discriminated by `method`. Each
	// variant carries only the arguments approved for that command: values that could
	// hold a secret — typed and composition text, URLs, referrers, scripts, templates,
	// file paths, drag contents and autofill values — are replaced by a length, a
	// count, a presence flag, an enum or a URL scheme and host.
	Data BrowserCdpCommandEventDataUnion `json:"data" api:"required"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64               `json:"ts" api:"required"`
	Type constant.CdpCommand `json:"type" default:"cdp_command"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Data        respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BrowserCdpCommandEventDataUnion contains all possible properties and values from
// [BrowserCdpCommandEventDataInputDispatchMouseEvent],
// [BrowserCdpCommandEventDataInputDispatchKeyEvent],
// [BrowserCdpCommandEventDataInputInsertText],
// [BrowserCdpCommandEventDataInputImeSetComposition],
// [BrowserCdpCommandEventDataInputDispatchTouchEvent],
// [BrowserCdpCommandEventDataInputDispatchDragEvent],
// [BrowserCdpCommandEventDataInputCancelDragging],
// [BrowserCdpCommandEventDataInputEmulateTouchFromMouseEvent],
// [BrowserCdpCommandEventDataInputSynthesizePinchGesture],
// [BrowserCdpCommandEventDataInputSynthesizeScrollGesture],
// [BrowserCdpCommandEventDataInputSynthesizeTapGesture],
// [BrowserCdpCommandEventDataDomSetFileInputFiles],
// [BrowserCdpCommandEventDataDomFocus],
// [BrowserCdpCommandEventDataDomScrollIntoViewIfNeeded],
// [BrowserCdpCommandEventDataPageBringToFront],
// [BrowserCdpCommandEventDataPageCaptureScreenshot],
// [BrowserCdpCommandEventDataPageCaptureSnapshot],
// [BrowserCdpCommandEventDataPageHandleJavaScriptDialog],
// [BrowserCdpCommandEventDataPageNavigate],
// [BrowserCdpCommandEventDataPageNavigateToHistoryEntry],
// [BrowserCdpCommandEventDataPageReload],
// [BrowserCdpCommandEventDataPagePrintToPdf],
// [BrowserCdpCommandEventDataPageStartScreencast],
// [BrowserCdpCommandEventDataPageStopScreencast],
// [BrowserCdpCommandEventDataPageStopLoading],
// [BrowserCdpCommandEventDataPageClose],
// [BrowserCdpCommandEventDataPageSetWebLifecycleState],
// [BrowserCdpCommandEventDataTargetActivateTarget],
// [BrowserCdpCommandEventDataTargetCloseTarget],
// [BrowserCdpCommandEventDataTargetCreateTarget],
// [BrowserCdpCommandEventDataTargetCreateBrowserContext],
// [BrowserCdpCommandEventDataTargetDisposeBrowserContext],
// [BrowserCdpCommandEventDataTargetOpenDevTools],
// [BrowserCdpCommandEventDataBrowserCancelDownload],
// [BrowserCdpCommandEventDataBrowserClose],
// [BrowserCdpCommandEventDataBrowserSetWindowBounds],
// [BrowserCdpCommandEventDataBrowserSetContentsSize],
// [BrowserCdpCommandEventDataAutofillTrigger].
//
// Use the [BrowserCdpCommandEventDataUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BrowserCdpCommandEventDataUnion struct {
	EventType string `json:"event_type"`
	// Any of "Input.dispatchMouseEvent", "Input.dispatchKeyEvent", "Input.insertText",
	// "Input.imeSetComposition", "Input.dispatchTouchEvent",
	// "Input.dispatchDragEvent", "Input.cancelDragging",
	// "Input.emulateTouchFromMouseEvent", "Input.synthesizePinchGesture",
	// "Input.synthesizeScrollGesture", "Input.synthesizeTapGesture",
	// "DOM.setFileInputFiles", "DOM.focus", "DOM.scrollIntoViewIfNeeded",
	// "Page.bringToFront", "Page.captureScreenshot", "Page.captureSnapshot",
	// "Page.handleJavaScriptDialog", "Page.navigate", "Page.navigateToHistoryEntry",
	// "Page.reload", "Page.printToPDF", "Page.startScreencast", "Page.stopScreencast",
	// "Page.stopLoading", "Page.close", "Page.setWebLifecycleState",
	// "Target.activateTarget", "Target.closeTarget", "Target.createTarget",
	// "Target.createBrowserContext", "Target.disposeBrowserContext",
	// "Target.openDevTools", "Browser.cancelDownload", "Browser.close",
	// "Browser.setWindowBounds", "Browser.setContentsSize", "Autofill.trigger".
	Method string `json:"method"`
	Button string `json:"button"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchMouseEvent].
	Buttons      int64   `json:"buttons"`
	ClickCount   int64   `json:"click_count"`
	CommandID    int64   `json:"command_id"`
	ConnectionID string  `json:"connection_id"`
	DeltaX       float64 `json:"delta_x"`
	DeltaY       float64 `json:"delta_y"`
	Force        float64 `json:"force"`
	Modifiers    int64   `json:"modifiers"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchMouseEvent].
	PointerType        string  `json:"pointer_type"`
	SessionID          string  `json:"session_id"`
	TangentialPressure float64 `json:"tangential_pressure"`
	TiltX              float64 `json:"tilt_x"`
	TiltY              float64 `json:"tilt_y"`
	Twist              int64   `json:"twist"`
	X                  float64 `json:"x"`
	Y                  float64 `json:"y"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchKeyEvent].
	AutoRepeat bool `json:"auto_repeat"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchKeyEvent].
	CommandCount int64 `json:"command_count"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchKeyEvent].
	IsKeypad bool `json:"is_keypad"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchKeyEvent].
	IsSystemKey bool `json:"is_system_key"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchKeyEvent].
	Location int64 `json:"location"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchKeyEvent].
	NamedKey   string `json:"named_key"`
	TextLength int64  `json:"text_length"`
	// This field is from variant [BrowserCdpCommandEventDataInputImeSetComposition].
	ReplacementEnd int64 `json:"replacement_end"`
	// This field is from variant [BrowserCdpCommandEventDataInputImeSetComposition].
	ReplacementStart int64 `json:"replacement_start"`
	// This field is from variant [BrowserCdpCommandEventDataInputImeSetComposition].
	SelectionEnd int64 `json:"selection_end"`
	// This field is from variant [BrowserCdpCommandEventDataInputImeSetComposition].
	SelectionStart int64 `json:"selection_start"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchTouchEvent].
	TouchPointCount int64 `json:"touch_point_count"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchTouchEvent].
	RadiusX float64 `json:"radius_x"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchTouchEvent].
	RadiusY float64 `json:"radius_y"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchTouchEvent].
	RotationAngle float64 `json:"rotation_angle"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchDragEvent].
	DragFileCount int64 `json:"drag_file_count"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchDragEvent].
	DragItemCount int64 `json:"drag_item_count"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchDragEvent].
	DragMimeCategories []string `json:"drag_mime_categories"`
	// This field is from variant [BrowserCdpCommandEventDataInputDispatchDragEvent].
	DragOperationsMask int64  `json:"drag_operations_mask"`
	GestureSourceType  string `json:"gesture_source_type"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizePinchGesture].
	RelativeSpeed int64 `json:"relative_speed"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizePinchGesture].
	ScaleFactor float64 `json:"scale_factor"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizeScrollGesture].
	PreventFling bool `json:"prevent_fling"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizeScrollGesture].
	RepeatCount int64 `json:"repeat_count"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizeScrollGesture].
	RepeatDelayMs int64 `json:"repeat_delay_ms"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizeScrollGesture].
	Speed int64 `json:"speed"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizeScrollGesture].
	XDistance float64 `json:"x_distance"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizeScrollGesture].
	XOverscroll float64 `json:"x_overscroll"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizeScrollGesture].
	YDistance float64 `json:"y_distance"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizeScrollGesture].
	YOverscroll float64 `json:"y_overscroll"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizeTapGesture].
	Duration int64 `json:"duration"`
	// This field is from variant
	// [BrowserCdpCommandEventDataInputSynthesizeTapGesture].
	TapCount int64 `json:"tap_count"`
	// This field is from variant [BrowserCdpCommandEventDataDomSetFileInputFiles].
	FileCount     int64  `json:"file_count"`
	BackendNodeID int64  `json:"backend_node_id"`
	NodeID        int64  `json:"node_id"`
	ObjectID      string `json:"object_id"`
	// This field is from variant
	// [BrowserCdpCommandEventDataDomScrollIntoViewIfNeeded].
	RectHeight float64 `json:"rect_height"`
	// This field is from variant
	// [BrowserCdpCommandEventDataDomScrollIntoViewIfNeeded].
	RectWidth float64 `json:"rect_width"`
	// This field is from variant
	// [BrowserCdpCommandEventDataDomScrollIntoViewIfNeeded].
	RectX float64 `json:"rect_x"`
	// This field is from variant
	// [BrowserCdpCommandEventDataDomScrollIntoViewIfNeeded].
	RectY float64 `json:"rect_y"`
	// This field is from variant [BrowserCdpCommandEventDataPageCaptureScreenshot].
	CaptureBeyondViewport bool `json:"capture_beyond_viewport"`
	// This field is from variant [BrowserCdpCommandEventDataPageCaptureScreenshot].
	ClipHeight float64 `json:"clip_height"`
	// This field is from variant [BrowserCdpCommandEventDataPageCaptureScreenshot].
	ClipScale float64 `json:"clip_scale"`
	// This field is from variant [BrowserCdpCommandEventDataPageCaptureScreenshot].
	ClipWidth float64 `json:"clip_width"`
	// This field is from variant [BrowserCdpCommandEventDataPageCaptureScreenshot].
	ClipX float64 `json:"clip_x"`
	// This field is from variant [BrowserCdpCommandEventDataPageCaptureScreenshot].
	ClipY  float64 `json:"clip_y"`
	Format string  `json:"format"`
	// This field is from variant [BrowserCdpCommandEventDataPageCaptureScreenshot].
	FromSurface bool `json:"from_surface"`
	// This field is from variant [BrowserCdpCommandEventDataPageCaptureScreenshot].
	OptimizeForSpeed bool  `json:"optimize_for_speed"`
	Quality          int64 `json:"quality"`
	// This field is from variant
	// [BrowserCdpCommandEventDataPageHandleJavaScriptDialog].
	Accept bool `json:"accept"`
	// This field is from variant
	// [BrowserCdpCommandEventDataPageHandleJavaScriptDialog].
	PromptTextLength int64  `json:"prompt_text_length"`
	FrameID          string `json:"frame_id"`
	// This field is from variant [BrowserCdpCommandEventDataPageNavigate].
	ReferrerPolicy string `json:"referrer_policy"`
	// This field is from variant [BrowserCdpCommandEventDataPageNavigate].
	ReferrerPresent bool `json:"referrer_present"`
	// This field is from variant [BrowserCdpCommandEventDataPageNavigate].
	TransitionType string `json:"transition_type"`
	URLScheme      string `json:"url_scheme"`
	// This field is from variant
	// [BrowserCdpCommandEventDataPageNavigateToHistoryEntry].
	EntryID int64 `json:"entry_id"`
	// This field is from variant [BrowserCdpCommandEventDataPageReload].
	IgnoreCache bool `json:"ignore_cache"`
	// This field is from variant [BrowserCdpCommandEventDataPageReload].
	LoaderID string `json:"loader_id"`
	// This field is from variant [BrowserCdpCommandEventDataPageReload].
	ScriptLength int64 `json:"script_length"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	DisplayHeaderFooter bool `json:"display_header_footer"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	FooterTemplatePresent bool `json:"footer_template_present"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	GenerateDocumentOutline bool `json:"generate_document_outline"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	GenerateTaggedPdf bool `json:"generate_tagged_pdf"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	HeaderTemplatePresent bool `json:"header_template_present"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	Landscape bool `json:"landscape"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	MarginBottom float64 `json:"margin_bottom"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	MarginLeft float64 `json:"margin_left"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	MarginRight float64 `json:"margin_right"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	MarginTop float64 `json:"margin_top"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	PageRangesPresent bool `json:"page_ranges_present"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	PaperHeight float64 `json:"paper_height"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	PaperWidth float64 `json:"paper_width"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	PreferCssPageSize bool `json:"prefer_css_page_size"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	PrintBackground bool `json:"print_background"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	Scale float64 `json:"scale"`
	// This field is from variant [BrowserCdpCommandEventDataPagePrintToPdf].
	TransferMode string `json:"transfer_mode"`
	// This field is from variant [BrowserCdpCommandEventDataPageStartScreencast].
	EveryNthFrame int64 `json:"every_nth_frame"`
	// This field is from variant [BrowserCdpCommandEventDataPageStartScreencast].
	MaxHeight int64 `json:"max_height"`
	// This field is from variant [BrowserCdpCommandEventDataPageStartScreencast].
	MaxWidth int64 `json:"max_width"`
	// This field is from variant [BrowserCdpCommandEventDataPageSetWebLifecycleState].
	State    string `json:"state"`
	TargetID string `json:"target_id"`
	// This field is from variant [BrowserCdpCommandEventDataTargetCreateTarget].
	Background       bool   `json:"background"`
	BrowserContextID string `json:"browser_context_id"`
	// This field is from variant [BrowserCdpCommandEventDataTargetCreateTarget].
	EnableBeginFrameControl bool `json:"enable_begin_frame_control"`
	// This field is from variant [BrowserCdpCommandEventDataTargetCreateTarget].
	Focus bool `json:"focus"`
	// This field is from variant [BrowserCdpCommandEventDataTargetCreateTarget].
	ForTab bool  `json:"for_tab"`
	Height int64 `json:"height"`
	// This field is from variant [BrowserCdpCommandEventDataTargetCreateTarget].
	Hidden bool  `json:"hidden"`
	Left   int64 `json:"left"`
	// This field is from variant [BrowserCdpCommandEventDataTargetCreateTarget].
	NewWindow   bool   `json:"new_window"`
	Top         int64  `json:"top"`
	Width       int64  `json:"width"`
	WindowState string `json:"window_state"`
	// This field is from variant
	// [BrowserCdpCommandEventDataTargetCreateBrowserContext].
	DisposeOnDetach bool `json:"dispose_on_detach"`
	// This field is from variant
	// [BrowserCdpCommandEventDataTargetCreateBrowserContext].
	ProxyBypassListPresent bool `json:"proxy_bypass_list_present"`
	// This field is from variant
	// [BrowserCdpCommandEventDataTargetCreateBrowserContext].
	ProxyServerPresent bool `json:"proxy_server_present"`
	// This field is from variant
	// [BrowserCdpCommandEventDataTargetCreateBrowserContext].
	UniversalNetworkAccessOriginCount int64 `json:"universal_network_access_origin_count"`
	// This field is from variant [BrowserCdpCommandEventDataTargetOpenDevTools].
	PanelID string `json:"panel_id"`
	// This field is from variant [BrowserCdpCommandEventDataBrowserCancelDownload].
	DownloadGuid string `json:"download_guid"`
	WindowID     int64  `json:"window_id"`
	// This field is from variant [BrowserCdpCommandEventDataAutofillTrigger].
	FieldID int64 `json:"field_id"`
	// This field is from variant [BrowserCdpCommandEventDataAutofillTrigger].
	AddressFieldCount int64 `json:"address_field_count"`
	// This field is from variant [BrowserCdpCommandEventDataAutofillTrigger].
	Mode string `json:"mode"`
	JSON struct {
		EventType                         respjson.Field
		Method                            respjson.Field
		Button                            respjson.Field
		Buttons                           respjson.Field
		ClickCount                        respjson.Field
		CommandID                         respjson.Field
		ConnectionID                      respjson.Field
		DeltaX                            respjson.Field
		DeltaY                            respjson.Field
		Force                             respjson.Field
		Modifiers                         respjson.Field
		PointerType                       respjson.Field
		SessionID                         respjson.Field
		TangentialPressure                respjson.Field
		TiltX                             respjson.Field
		TiltY                             respjson.Field
		Twist                             respjson.Field
		X                                 respjson.Field
		Y                                 respjson.Field
		AutoRepeat                        respjson.Field
		CommandCount                      respjson.Field
		IsKeypad                          respjson.Field
		IsSystemKey                       respjson.Field
		Location                          respjson.Field
		NamedKey                          respjson.Field
		TextLength                        respjson.Field
		ReplacementEnd                    respjson.Field
		ReplacementStart                  respjson.Field
		SelectionEnd                      respjson.Field
		SelectionStart                    respjson.Field
		TouchPointCount                   respjson.Field
		RadiusX                           respjson.Field
		RadiusY                           respjson.Field
		RotationAngle                     respjson.Field
		DragFileCount                     respjson.Field
		DragItemCount                     respjson.Field
		DragMimeCategories                respjson.Field
		DragOperationsMask                respjson.Field
		GestureSourceType                 respjson.Field
		RelativeSpeed                     respjson.Field
		ScaleFactor                       respjson.Field
		PreventFling                      respjson.Field
		RepeatCount                       respjson.Field
		RepeatDelayMs                     respjson.Field
		Speed                             respjson.Field
		XDistance                         respjson.Field
		XOverscroll                       respjson.Field
		YDistance                         respjson.Field
		YOverscroll                       respjson.Field
		Duration                          respjson.Field
		TapCount                          respjson.Field
		FileCount                         respjson.Field
		BackendNodeID                     respjson.Field
		NodeID                            respjson.Field
		ObjectID                          respjson.Field
		RectHeight                        respjson.Field
		RectWidth                         respjson.Field
		RectX                             respjson.Field
		RectY                             respjson.Field
		CaptureBeyondViewport             respjson.Field
		ClipHeight                        respjson.Field
		ClipScale                         respjson.Field
		ClipWidth                         respjson.Field
		ClipX                             respjson.Field
		ClipY                             respjson.Field
		Format                            respjson.Field
		FromSurface                       respjson.Field
		OptimizeForSpeed                  respjson.Field
		Quality                           respjson.Field
		Accept                            respjson.Field
		PromptTextLength                  respjson.Field
		FrameID                           respjson.Field
		ReferrerPolicy                    respjson.Field
		ReferrerPresent                   respjson.Field
		TransitionType                    respjson.Field
		URLScheme                         respjson.Field
		EntryID                           respjson.Field
		IgnoreCache                       respjson.Field
		LoaderID                          respjson.Field
		ScriptLength                      respjson.Field
		DisplayHeaderFooter               respjson.Field
		FooterTemplatePresent             respjson.Field
		GenerateDocumentOutline           respjson.Field
		GenerateTaggedPdf                 respjson.Field
		HeaderTemplatePresent             respjson.Field
		Landscape                         respjson.Field
		MarginBottom                      respjson.Field
		MarginLeft                        respjson.Field
		MarginRight                       respjson.Field
		MarginTop                         respjson.Field
		PageRangesPresent                 respjson.Field
		PaperHeight                       respjson.Field
		PaperWidth                        respjson.Field
		PreferCssPageSize                 respjson.Field
		PrintBackground                   respjson.Field
		Scale                             respjson.Field
		TransferMode                      respjson.Field
		EveryNthFrame                     respjson.Field
		MaxHeight                         respjson.Field
		MaxWidth                          respjson.Field
		State                             respjson.Field
		TargetID                          respjson.Field
		Background                        respjson.Field
		BrowserContextID                  respjson.Field
		EnableBeginFrameControl           respjson.Field
		Focus                             respjson.Field
		ForTab                            respjson.Field
		Height                            respjson.Field
		Hidden                            respjson.Field
		Left                              respjson.Field
		NewWindow                         respjson.Field
		Top                               respjson.Field
		Width                             respjson.Field
		WindowState                       respjson.Field
		DisposeOnDetach                   respjson.Field
		ProxyBypassListPresent            respjson.Field
		ProxyServerPresent                respjson.Field
		UniversalNetworkAccessOriginCount respjson.Field
		PanelID                           respjson.Field
		DownloadGuid                      respjson.Field
		WindowID                          respjson.Field
		FieldID                           respjson.Field
		AddressFieldCount                 respjson.Field
		Mode                              respjson.Field
		raw                               string
	} `json:"-"`
}

// anyBrowserCdpCommandEventData is implemented by each variant of
// [BrowserCdpCommandEventDataUnion] to add type safety for the return type of
// [BrowserCdpCommandEventDataUnion.AsAny]
type anyBrowserCdpCommandEventData interface {
	implBrowserCdpCommandEventDataUnion()
}

func (BrowserCdpCommandEventDataInputDispatchMouseEvent) implBrowserCdpCommandEventDataUnion() {}
func (BrowserCdpCommandEventDataInputDispatchKeyEvent) implBrowserCdpCommandEventDataUnion()   {}
func (BrowserCdpCommandEventDataInputInsertText) implBrowserCdpCommandEventDataUnion()         {}
func (BrowserCdpCommandEventDataInputImeSetComposition) implBrowserCdpCommandEventDataUnion()  {}
func (BrowserCdpCommandEventDataInputDispatchTouchEvent) implBrowserCdpCommandEventDataUnion() {}
func (BrowserCdpCommandEventDataInputDispatchDragEvent) implBrowserCdpCommandEventDataUnion()  {}
func (BrowserCdpCommandEventDataInputCancelDragging) implBrowserCdpCommandEventDataUnion()     {}
func (BrowserCdpCommandEventDataInputEmulateTouchFromMouseEvent) implBrowserCdpCommandEventDataUnion() {
}
func (BrowserCdpCommandEventDataInputSynthesizePinchGesture) implBrowserCdpCommandEventDataUnion()  {}
func (BrowserCdpCommandEventDataInputSynthesizeScrollGesture) implBrowserCdpCommandEventDataUnion() {}
func (BrowserCdpCommandEventDataInputSynthesizeTapGesture) implBrowserCdpCommandEventDataUnion()    {}
func (BrowserCdpCommandEventDataDomSetFileInputFiles) implBrowserCdpCommandEventDataUnion()         {}
func (BrowserCdpCommandEventDataDomFocus) implBrowserCdpCommandEventDataUnion()                     {}
func (BrowserCdpCommandEventDataDomScrollIntoViewIfNeeded) implBrowserCdpCommandEventDataUnion()    {}
func (BrowserCdpCommandEventDataPageBringToFront) implBrowserCdpCommandEventDataUnion()             {}
func (BrowserCdpCommandEventDataPageCaptureScreenshot) implBrowserCdpCommandEventDataUnion()        {}
func (BrowserCdpCommandEventDataPageCaptureSnapshot) implBrowserCdpCommandEventDataUnion()          {}
func (BrowserCdpCommandEventDataPageHandleJavaScriptDialog) implBrowserCdpCommandEventDataUnion()   {}
func (BrowserCdpCommandEventDataPageNavigate) implBrowserCdpCommandEventDataUnion()                 {}
func (BrowserCdpCommandEventDataPageNavigateToHistoryEntry) implBrowserCdpCommandEventDataUnion()   {}
func (BrowserCdpCommandEventDataPageReload) implBrowserCdpCommandEventDataUnion()                   {}
func (BrowserCdpCommandEventDataPagePrintToPdf) implBrowserCdpCommandEventDataUnion()               {}
func (BrowserCdpCommandEventDataPageStartScreencast) implBrowserCdpCommandEventDataUnion()          {}
func (BrowserCdpCommandEventDataPageStopScreencast) implBrowserCdpCommandEventDataUnion()           {}
func (BrowserCdpCommandEventDataPageStopLoading) implBrowserCdpCommandEventDataUnion()              {}
func (BrowserCdpCommandEventDataPageClose) implBrowserCdpCommandEventDataUnion()                    {}
func (BrowserCdpCommandEventDataPageSetWebLifecycleState) implBrowserCdpCommandEventDataUnion()     {}
func (BrowserCdpCommandEventDataTargetActivateTarget) implBrowserCdpCommandEventDataUnion()         {}
func (BrowserCdpCommandEventDataTargetCloseTarget) implBrowserCdpCommandEventDataUnion()            {}
func (BrowserCdpCommandEventDataTargetCreateTarget) implBrowserCdpCommandEventDataUnion()           {}
func (BrowserCdpCommandEventDataTargetCreateBrowserContext) implBrowserCdpCommandEventDataUnion()   {}
func (BrowserCdpCommandEventDataTargetDisposeBrowserContext) implBrowserCdpCommandEventDataUnion()  {}
func (BrowserCdpCommandEventDataTargetOpenDevTools) implBrowserCdpCommandEventDataUnion()           {}
func (BrowserCdpCommandEventDataBrowserCancelDownload) implBrowserCdpCommandEventDataUnion()        {}
func (BrowserCdpCommandEventDataBrowserClose) implBrowserCdpCommandEventDataUnion()                 {}
func (BrowserCdpCommandEventDataBrowserSetWindowBounds) implBrowserCdpCommandEventDataUnion()       {}
func (BrowserCdpCommandEventDataBrowserSetContentsSize) implBrowserCdpCommandEventDataUnion()       {}
func (BrowserCdpCommandEventDataAutofillTrigger) implBrowserCdpCommandEventDataUnion()              {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BrowserCdpCommandEventDataUnion.AsAny().(type) {
//	case kernel.BrowserCdpCommandEventDataInputDispatchMouseEvent:
//	case kernel.BrowserCdpCommandEventDataInputDispatchKeyEvent:
//	case kernel.BrowserCdpCommandEventDataInputInsertText:
//	case kernel.BrowserCdpCommandEventDataInputImeSetComposition:
//	case kernel.BrowserCdpCommandEventDataInputDispatchTouchEvent:
//	case kernel.BrowserCdpCommandEventDataInputDispatchDragEvent:
//	case kernel.BrowserCdpCommandEventDataInputCancelDragging:
//	case kernel.BrowserCdpCommandEventDataInputEmulateTouchFromMouseEvent:
//	case kernel.BrowserCdpCommandEventDataInputSynthesizePinchGesture:
//	case kernel.BrowserCdpCommandEventDataInputSynthesizeScrollGesture:
//	case kernel.BrowserCdpCommandEventDataInputSynthesizeTapGesture:
//	case kernel.BrowserCdpCommandEventDataDomSetFileInputFiles:
//	case kernel.BrowserCdpCommandEventDataDomFocus:
//	case kernel.BrowserCdpCommandEventDataDomScrollIntoViewIfNeeded:
//	case kernel.BrowserCdpCommandEventDataPageBringToFront:
//	case kernel.BrowserCdpCommandEventDataPageCaptureScreenshot:
//	case kernel.BrowserCdpCommandEventDataPageCaptureSnapshot:
//	case kernel.BrowserCdpCommandEventDataPageHandleJavaScriptDialog:
//	case kernel.BrowserCdpCommandEventDataPageNavigate:
//	case kernel.BrowserCdpCommandEventDataPageNavigateToHistoryEntry:
//	case kernel.BrowserCdpCommandEventDataPageReload:
//	case kernel.BrowserCdpCommandEventDataPagePrintToPdf:
//	case kernel.BrowserCdpCommandEventDataPageStartScreencast:
//	case kernel.BrowserCdpCommandEventDataPageStopScreencast:
//	case kernel.BrowserCdpCommandEventDataPageStopLoading:
//	case kernel.BrowserCdpCommandEventDataPageClose:
//	case kernel.BrowserCdpCommandEventDataPageSetWebLifecycleState:
//	case kernel.BrowserCdpCommandEventDataTargetActivateTarget:
//	case kernel.BrowserCdpCommandEventDataTargetCloseTarget:
//	case kernel.BrowserCdpCommandEventDataTargetCreateTarget:
//	case kernel.BrowserCdpCommandEventDataTargetCreateBrowserContext:
//	case kernel.BrowserCdpCommandEventDataTargetDisposeBrowserContext:
//	case kernel.BrowserCdpCommandEventDataTargetOpenDevTools:
//	case kernel.BrowserCdpCommandEventDataBrowserCancelDownload:
//	case kernel.BrowserCdpCommandEventDataBrowserClose:
//	case kernel.BrowserCdpCommandEventDataBrowserSetWindowBounds:
//	case kernel.BrowserCdpCommandEventDataBrowserSetContentsSize:
//	case kernel.BrowserCdpCommandEventDataAutofillTrigger:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BrowserCdpCommandEventDataUnion) AsAny() anyBrowserCdpCommandEventData {
	switch u.Method {
	case "Input.dispatchMouseEvent":
		return u.AsInputDispatchMouseEvent()
	case "Input.dispatchKeyEvent":
		return u.AsInputDispatchKeyEvent()
	case "Input.insertText":
		return u.AsInputInsertText()
	case "Input.imeSetComposition":
		return u.AsInputImeSetComposition()
	case "Input.dispatchTouchEvent":
		return u.AsInputDispatchTouchEvent()
	case "Input.dispatchDragEvent":
		return u.AsInputDispatchDragEvent()
	case "Input.cancelDragging":
		return u.AsInputCancelDragging()
	case "Input.emulateTouchFromMouseEvent":
		return u.AsInputEmulateTouchFromMouseEvent()
	case "Input.synthesizePinchGesture":
		return u.AsInputSynthesizePinchGesture()
	case "Input.synthesizeScrollGesture":
		return u.AsInputSynthesizeScrollGesture()
	case "Input.synthesizeTapGesture":
		return u.AsInputSynthesizeTapGesture()
	case "DOM.setFileInputFiles":
		return u.AsDomSetFileInputFiles()
	case "DOM.focus":
		return u.AsDomFocus()
	case "DOM.scrollIntoViewIfNeeded":
		return u.AsDomScrollIntoViewIfNeeded()
	case "Page.bringToFront":
		return u.AsPageBringToFront()
	case "Page.captureScreenshot":
		return u.AsPageCaptureScreenshot()
	case "Page.captureSnapshot":
		return u.AsPageCaptureSnapshot()
	case "Page.handleJavaScriptDialog":
		return u.AsPageHandleJavaScriptDialog()
	case "Page.navigate":
		return u.AsPageNavigate()
	case "Page.navigateToHistoryEntry":
		return u.AsPageNavigateToHistoryEntry()
	case "Page.reload":
		return u.AsPageReload()
	case "Page.printToPDF":
		return u.AsPagePrintToPdf()
	case "Page.startScreencast":
		return u.AsPageStartScreencast()
	case "Page.stopScreencast":
		return u.AsPageStopScreencast()
	case "Page.stopLoading":
		return u.AsPageStopLoading()
	case "Page.close":
		return u.AsPageClose()
	case "Page.setWebLifecycleState":
		return u.AsPageSetWebLifecycleState()
	case "Target.activateTarget":
		return u.AsTargetActivateTarget()
	case "Target.closeTarget":
		return u.AsTargetCloseTarget()
	case "Target.createTarget":
		return u.AsTargetCreateTarget()
	case "Target.createBrowserContext":
		return u.AsTargetCreateBrowserContext()
	case "Target.disposeBrowserContext":
		return u.AsTargetDisposeBrowserContext()
	case "Target.openDevTools":
		return u.AsTargetOpenDevTools()
	case "Browser.cancelDownload":
		return u.AsBrowserCancelDownload()
	case "Browser.close":
		return u.AsBrowserClose()
	case "Browser.setWindowBounds":
		return u.AsBrowserSetWindowBounds()
	case "Browser.setContentsSize":
		return u.AsBrowserSetContentsSize()
	case "Autofill.trigger":
		return u.AsAutofillTrigger()
	}
	return nil
}

func (u BrowserCdpCommandEventDataUnion) AsInputDispatchMouseEvent() (v BrowserCdpCommandEventDataInputDispatchMouseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsInputDispatchKeyEvent() (v BrowserCdpCommandEventDataInputDispatchKeyEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsInputInsertText() (v BrowserCdpCommandEventDataInputInsertText) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsInputImeSetComposition() (v BrowserCdpCommandEventDataInputImeSetComposition) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsInputDispatchTouchEvent() (v BrowserCdpCommandEventDataInputDispatchTouchEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsInputDispatchDragEvent() (v BrowserCdpCommandEventDataInputDispatchDragEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsInputCancelDragging() (v BrowserCdpCommandEventDataInputCancelDragging) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsInputEmulateTouchFromMouseEvent() (v BrowserCdpCommandEventDataInputEmulateTouchFromMouseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsInputSynthesizePinchGesture() (v BrowserCdpCommandEventDataInputSynthesizePinchGesture) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsInputSynthesizeScrollGesture() (v BrowserCdpCommandEventDataInputSynthesizeScrollGesture) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsInputSynthesizeTapGesture() (v BrowserCdpCommandEventDataInputSynthesizeTapGesture) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsDomSetFileInputFiles() (v BrowserCdpCommandEventDataDomSetFileInputFiles) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsDomFocus() (v BrowserCdpCommandEventDataDomFocus) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsDomScrollIntoViewIfNeeded() (v BrowserCdpCommandEventDataDomScrollIntoViewIfNeeded) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageBringToFront() (v BrowserCdpCommandEventDataPageBringToFront) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageCaptureScreenshot() (v BrowserCdpCommandEventDataPageCaptureScreenshot) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageCaptureSnapshot() (v BrowserCdpCommandEventDataPageCaptureSnapshot) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageHandleJavaScriptDialog() (v BrowserCdpCommandEventDataPageHandleJavaScriptDialog) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageNavigate() (v BrowserCdpCommandEventDataPageNavigate) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageNavigateToHistoryEntry() (v BrowserCdpCommandEventDataPageNavigateToHistoryEntry) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageReload() (v BrowserCdpCommandEventDataPageReload) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPagePrintToPdf() (v BrowserCdpCommandEventDataPagePrintToPdf) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageStartScreencast() (v BrowserCdpCommandEventDataPageStartScreencast) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageStopScreencast() (v BrowserCdpCommandEventDataPageStopScreencast) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageStopLoading() (v BrowserCdpCommandEventDataPageStopLoading) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageClose() (v BrowserCdpCommandEventDataPageClose) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsPageSetWebLifecycleState() (v BrowserCdpCommandEventDataPageSetWebLifecycleState) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsTargetActivateTarget() (v BrowserCdpCommandEventDataTargetActivateTarget) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsTargetCloseTarget() (v BrowserCdpCommandEventDataTargetCloseTarget) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsTargetCreateTarget() (v BrowserCdpCommandEventDataTargetCreateTarget) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsTargetCreateBrowserContext() (v BrowserCdpCommandEventDataTargetCreateBrowserContext) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsTargetDisposeBrowserContext() (v BrowserCdpCommandEventDataTargetDisposeBrowserContext) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsTargetOpenDevTools() (v BrowserCdpCommandEventDataTargetOpenDevTools) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsBrowserCancelDownload() (v BrowserCdpCommandEventDataBrowserCancelDownload) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsBrowserClose() (v BrowserCdpCommandEventDataBrowserClose) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsBrowserSetWindowBounds() (v BrowserCdpCommandEventDataBrowserSetWindowBounds) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsBrowserSetContentsSize() (v BrowserCdpCommandEventDataBrowserSetContentsSize) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserCdpCommandEventDataUnion) AsAutofillTrigger() (v BrowserCdpCommandEventDataAutofillTrigger) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BrowserCdpCommandEventDataUnion) RawJSON() string { return u.JSON.raw }

func (r *BrowserCdpCommandEventDataUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Input.dispatchMouseEvent` arguments. Canonical input:
// `Input.dispatchMouseEvent` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataInputDispatchMouseEvent struct {
	// Mouse event phase: `mousePressed`, `mouseReleased`, `mouseMoved` or
	// `mouseWheel`. A value the protocol does not define is reported as `other`.
	//
	// Any of "mousePressed", "mouseReleased", "mouseMoved", "mouseWheel", "other".
	EventType string                           `json:"event_type" api:"required"`
	Method    constant.InputDispatchMouseEvent `json:"method" default:"Input.dispatchMouseEvent"`
	// Button named by the command (`none`, `left`, `middle`, `right`, `back`,
	// `forward`). A value the protocol does not define is reported as `other`.
	//
	// Any of "none", "left", "middle", "right", "back", "forward", "other".
	Button string `json:"button"`
	// Bit field of buttons held down. Non-zero on a `mouseMoved` means the move is a
	// drag path.
	Buttons int64 `json:"buttons"`
	// Number of times the button was clicked (2 is a double click).
	ClickCount int64 `json:"click_count"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Horizontal scroll delta, for `mouseWheel`.
	DeltaX float64 `json:"delta_x"`
	// Vertical scroll delta, for `mouseWheel`.
	DeltaY float64 `json:"delta_y"`
	// Normalized pressure, 0 to 1.
	Force float64 `json:"force"`
	// Bit field of held modifier keys (1=Alt, 2=Ctrl, 4=Meta, 8=Shift).
	Modifiers int64 `json:"modifiers"`
	// Pointer that generated the event (`mouse` or `pen`). A value the protocol does
	// not define is reported as `other`.
	//
	// Any of "mouse", "pen", "other".
	PointerType string `json:"pointer_type"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Normalized tangential pressure, -1 to 1.
	TangentialPressure float64 `json:"tangential_pressure"`
	// Pen tilt from the Y-Z plane, in degrees.
	TiltX float64 `json:"tilt_x"`
	// Pen tilt from the X-Z plane, in degrees.
	TiltY float64 `json:"tilt_y"`
	// Pen clockwise rotation, in degrees.
	Twist int64 `json:"twist"`
	// Viewport x coordinate in CSS pixels.
	X float64 `json:"x"`
	// Viewport y coordinate in CSS pixels.
	Y float64 `json:"y"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventType          respjson.Field
		Method             respjson.Field
		Button             respjson.Field
		Buttons            respjson.Field
		ClickCount         respjson.Field
		CommandID          respjson.Field
		ConnectionID       respjson.Field
		DeltaX             respjson.Field
		DeltaY             respjson.Field
		Force              respjson.Field
		Modifiers          respjson.Field
		PointerType        respjson.Field
		SessionID          respjson.Field
		TangentialPressure respjson.Field
		TiltX              respjson.Field
		TiltY              respjson.Field
		Twist              respjson.Field
		X                  respjson.Field
		Y                  respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataInputDispatchMouseEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataInputDispatchMouseEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Input.dispatchKeyEvent` arguments. Canonical input:
// `Input.dispatchKeyEvent` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataInputDispatchKeyEvent struct {
	// Key event phase: `keyDown`, `keyUp`, `rawKeyDown` or `char`. A value the
	// protocol does not define is reported as `other`.
	//
	// Any of "keyDown", "keyUp", "rawKeyDown", "char", "other".
	EventType string                         `json:"event_type" api:"required"`
	Method    constant.InputDispatchKeyEvent `json:"method" default:"Input.dispatchKeyEvent"`
	// Whether the event was generated by key repeat.
	AutoRepeat bool `json:"auto_repeat"`
	// Number of editing commands (e.g. `selectAll`) carried by the event.
	CommandCount int64 `json:"command_count"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Whether the key is on the numeric keypad.
	IsKeypad bool `json:"is_keypad"`
	// Whether the event is a system key event.
	IsSystemKey bool `json:"is_system_key"`
	// Keyboard location (1=left, 2=right, 3=numpad).
	Location int64 `json:"location"`
	// Bit field of held modifier keys (1=Alt, 2=Ctrl, 4=Meta, 8=Shift).
	Modifiers int64 `json:"modifiers"`
	// Key that commands the page rather than typing into it (e.g. `Enter`, `Tab`,
	// `ArrowDown`, `F5`). Keys that produce a character are never captured; those are
	// counted by `text_length`.
	NamedKey string `json:"named_key"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Number of characters the command submitted. The text itself is never captured.
	TextLength int64 `json:"text_length"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventType    respjson.Field
		Method       respjson.Field
		AutoRepeat   respjson.Field
		CommandCount respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		IsKeypad     respjson.Field
		IsSystemKey  respjson.Field
		Location     respjson.Field
		Modifiers    respjson.Field
		NamedKey     respjson.Field
		SessionID    respjson.Field
		TextLength   respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataInputDispatchKeyEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataInputDispatchKeyEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Input.insertText` arguments. Canonical input: `Input.insertText` in
// devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataInputInsertText struct {
	Method constant.InputInsertText `json:"method" default:"Input.insertText"`
	// Number of characters inserted. The text itself is never captured.
	TextLength int64 `json:"text_length" api:"required"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		TextLength   respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataInputInsertText) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataInputInsertText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Input.imeSetComposition` arguments. Canonical input:
// `Input.imeSetComposition` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataInputImeSetComposition struct {
	Method constant.InputImeSetComposition `json:"method" default:"Input.imeSetComposition"`
	// Number of characters in the composition. The text itself is never captured.
	TextLength int64 `json:"text_length" api:"required"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Replacement range end offset.
	ReplacementEnd int64 `json:"replacement_end"`
	// Replacement range start offset.
	ReplacementStart int64 `json:"replacement_start"`
	// Selection end offset within the composition.
	SelectionEnd int64 `json:"selection_end"`
	// Selection start offset within the composition.
	SelectionStart int64 `json:"selection_start"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method           respjson.Field
		TextLength       respjson.Field
		CommandID        respjson.Field
		ConnectionID     respjson.Field
		ReplacementEnd   respjson.Field
		ReplacementStart respjson.Field
		SelectionEnd     respjson.Field
		SelectionStart   respjson.Field
		SessionID        respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataInputImeSetComposition) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataInputImeSetComposition) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Input.dispatchTouchEvent` arguments. Canonical input:
// `Input.dispatchTouchEvent` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataInputDispatchTouchEvent struct {
	// Touch event phase: `touchStart`, `touchEnd`, `touchMove` or `touchCancel`. A
	// value the protocol does not define is reported as `other`.
	//
	// Any of "touchStart", "touchEnd", "touchMove", "touchCancel", "other".
	EventType string                           `json:"event_type" api:"required"`
	Method    constant.InputDispatchTouchEvent `json:"method" default:"Input.dispatchTouchEvent"`
	// Number of active touch points the command carried.
	TouchPointCount int64 `json:"touch_point_count" api:"required"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Normalized pressure of the first touch point, 0 to 1.
	Force float64 `json:"force"`
	// Bit field of held modifier keys (1=Alt, 2=Ctrl, 4=Meta, 8=Shift).
	Modifiers int64 `json:"modifiers"`
	// Horizontal radius of the first touch point.
	RadiusX float64 `json:"radius_x"`
	// Vertical radius of the first touch point.
	RadiusY float64 `json:"radius_y"`
	// Rotation of the first touch point, in degrees.
	RotationAngle float64 `json:"rotation_angle"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Normalized tangential pressure of the first touch point, -1 to 1.
	TangentialPressure float64 `json:"tangential_pressure"`
	// Tilt of the first touch point from the Y-Z plane, in degrees.
	TiltX float64 `json:"tilt_x"`
	// Tilt of the first touch point from the X-Z plane, in degrees.
	TiltY float64 `json:"tilt_y"`
	// Clockwise rotation of the first touch point, in degrees.
	Twist int64 `json:"twist"`
	// Viewport x coordinate of the first touch point. Touch coordinates live inside
	// `touchPoints`, so this is the primary point rather than a command-level
	// argument.
	X float64 `json:"x"`
	// Viewport y coordinate of the first touch point.
	Y float64 `json:"y"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventType          respjson.Field
		Method             respjson.Field
		TouchPointCount    respjson.Field
		CommandID          respjson.Field
		ConnectionID       respjson.Field
		Force              respjson.Field
		Modifiers          respjson.Field
		RadiusX            respjson.Field
		RadiusY            respjson.Field
		RotationAngle      respjson.Field
		SessionID          respjson.Field
		TangentialPressure respjson.Field
		TiltX              respjson.Field
		TiltY              respjson.Field
		Twist              respjson.Field
		X                  respjson.Field
		Y                  respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataInputDispatchTouchEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataInputDispatchTouchEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Input.dispatchDragEvent` arguments. Canonical input:
// `Input.dispatchDragEvent` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataInputDispatchDragEvent struct {
	// Drag event phase: `dragEnter`, `dragOver`, `drop` or `dragCancel`. A value the
	// protocol does not define is reported as `other`.
	//
	// Any of "dragEnter", "dragOver", "drop", "dragCancel", "other".
	EventType string                          `json:"event_type" api:"required"`
	Method    constant.InputDispatchDragEvent `json:"method" default:"Input.dispatchDragEvent"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Number of files in the drag payload. File paths are never captured.
	DragFileCount int64 `json:"drag_file_count"`
	// Number of items in the drag payload. Item contents are never captured.
	DragItemCount int64 `json:"drag_item_count"`
	// Distinct top-level MIME categories of the drag items (e.g. `text`, `image`,
	// `application`). Subtypes and contents are never captured. A value the protocol
	// does not define is reported as `other`.
	//
	// Any of "text", "image", "audio", "video", "application", "font", "model",
	// "multipart", "message", "other".
	DragMimeCategories []string `json:"drag_mime_categories"`
	// Bit field of allowed drag operations (1=copy, 2=link, 16=move).
	DragOperationsMask int64 `json:"drag_operations_mask"`
	// Bit field of held modifier keys (1=Alt, 2=Ctrl, 4=Meta, 8=Shift).
	Modifiers int64 `json:"modifiers"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Viewport x coordinate in CSS pixels.
	X float64 `json:"x"`
	// Viewport y coordinate in CSS pixels.
	Y float64 `json:"y"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventType          respjson.Field
		Method             respjson.Field
		CommandID          respjson.Field
		ConnectionID       respjson.Field
		DragFileCount      respjson.Field
		DragItemCount      respjson.Field
		DragMimeCategories respjson.Field
		DragOperationsMask respjson.Field
		Modifiers          respjson.Field
		SessionID          respjson.Field
		X                  respjson.Field
		Y                  respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataInputDispatchDragEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataInputDispatchDragEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Input.cancelDragging` arguments. Canonical input:
// `Input.cancelDragging` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataInputCancelDragging struct {
	Method constant.InputCancelDragging `json:"method" default:"Input.cancelDragging"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataInputCancelDragging) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataInputCancelDragging) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Input.emulateTouchFromMouseEvent` arguments. Canonical input:
// `Input.emulateTouchFromMouseEvent` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataInputEmulateTouchFromMouseEvent struct {
	// Mouse event phase being emulated as touch. A value the protocol does not define
	// is reported as `other`.
	//
	// Any of "mousePressed", "mouseReleased", "mouseMoved", "mouseWheel", "other".
	EventType string                                   `json:"event_type" api:"required"`
	Method    constant.InputEmulateTouchFromMouseEvent `json:"method" default:"Input.emulateTouchFromMouseEvent"`
	// Button named by the command. A value the protocol does not define is reported as
	// `other`.
	//
	// Any of "none", "left", "middle", "right", "back", "forward", "other".
	Button string `json:"button"`
	// Number of times the button was clicked.
	ClickCount int64 `json:"click_count"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Horizontal scroll delta.
	DeltaX float64 `json:"delta_x"`
	// Vertical scroll delta.
	DeltaY float64 `json:"delta_y"`
	// Bit field of held modifier keys (1=Alt, 2=Ctrl, 4=Meta, 8=Shift).
	Modifiers int64 `json:"modifiers"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Viewport x coordinate in CSS pixels.
	X float64 `json:"x"`
	// Viewport y coordinate in CSS pixels.
	Y float64 `json:"y"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventType    respjson.Field
		Method       respjson.Field
		Button       respjson.Field
		ClickCount   respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		DeltaX       respjson.Field
		DeltaY       respjson.Field
		Modifiers    respjson.Field
		SessionID    respjson.Field
		X            respjson.Field
		Y            respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataInputEmulateTouchFromMouseEvent) RawJSON() string {
	return r.JSON.raw
}
func (r *BrowserCdpCommandEventDataInputEmulateTouchFromMouseEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Input.synthesizePinchGesture` arguments. Canonical input:
// `Input.synthesizePinchGesture` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataInputSynthesizePinchGesture struct {
	Method constant.InputSynthesizePinchGesture `json:"method" default:"Input.synthesizePinchGesture"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Input source the synthesized gesture emulates. A value the protocol does not
	// define is reported as `other`.
	//
	// Any of "default", "touch", "mouse", "other".
	GestureSourceType string `json:"gesture_source_type"`
	// Relative pointer speed, in pixels per second.
	RelativeSpeed int64 `json:"relative_speed"`
	// Relative scale of the pinch (>1 zooms in).
	ScaleFactor float64 `json:"scale_factor"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Viewport x coordinate in CSS pixels.
	X float64 `json:"x"`
	// Viewport y coordinate in CSS pixels.
	Y float64 `json:"y"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method            respjson.Field
		CommandID         respjson.Field
		ConnectionID      respjson.Field
		GestureSourceType respjson.Field
		RelativeSpeed     respjson.Field
		ScaleFactor       respjson.Field
		SessionID         respjson.Field
		X                 respjson.Field
		Y                 respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataInputSynthesizePinchGesture) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataInputSynthesizePinchGesture) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Input.synthesizeScrollGesture` arguments. Canonical input:
// `Input.synthesizeScrollGesture` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataInputSynthesizeScrollGesture struct {
	Method constant.InputSynthesizeScrollGesture `json:"method" default:"Input.synthesizeScrollGesture"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Input source the synthesized gesture emulates. A value the protocol does not
	// define is reported as `other`.
	//
	// Any of "default", "touch", "mouse", "other".
	GestureSourceType string `json:"gesture_source_type"`
	// Whether fling was suppressed.
	PreventFling bool `json:"prevent_fling"`
	// Number of additional repeats of the scroll.
	RepeatCount int64 `json:"repeat_count"`
	// Delay between repeats, in milliseconds.
	RepeatDelayMs int64 `json:"repeat_delay_ms"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Swipe speed in pixels per second.
	Speed int64 `json:"speed"`
	// Viewport x coordinate in CSS pixels.
	X float64 `json:"x"`
	// Horizontal scroll distance in CSS pixels; positive scrolls left.
	XDistance float64 `json:"x_distance"`
	// Additional horizontal distance scrolled past the end.
	XOverscroll float64 `json:"x_overscroll"`
	// Viewport y coordinate in CSS pixels.
	Y float64 `json:"y"`
	// Vertical scroll distance in CSS pixels; positive scrolls up.
	YDistance float64 `json:"y_distance"`
	// Additional vertical distance scrolled past the end.
	YOverscroll float64 `json:"y_overscroll"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method            respjson.Field
		CommandID         respjson.Field
		ConnectionID      respjson.Field
		GestureSourceType respjson.Field
		PreventFling      respjson.Field
		RepeatCount       respjson.Field
		RepeatDelayMs     respjson.Field
		SessionID         respjson.Field
		Speed             respjson.Field
		X                 respjson.Field
		XDistance         respjson.Field
		XOverscroll       respjson.Field
		Y                 respjson.Field
		YDistance         respjson.Field
		YOverscroll       respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataInputSynthesizeScrollGesture) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataInputSynthesizeScrollGesture) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Input.synthesizeTapGesture` arguments. Canonical input:
// `Input.synthesizeTapGesture` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataInputSynthesizeTapGesture struct {
	Method constant.InputSynthesizeTapGesture `json:"method" default:"Input.synthesizeTapGesture"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Duration between touchdown and touchup, in milliseconds.
	Duration int64 `json:"duration"`
	// Input source the synthesized gesture emulates. A value the protocol does not
	// define is reported as `other`.
	//
	// Any of "default", "touch", "mouse", "other".
	GestureSourceType string `json:"gesture_source_type"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Number of times to tap (2 is a double tap).
	TapCount int64 `json:"tap_count"`
	// Viewport x coordinate in CSS pixels.
	X float64 `json:"x"`
	// Viewport y coordinate in CSS pixels.
	Y float64 `json:"y"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method            respjson.Field
		CommandID         respjson.Field
		ConnectionID      respjson.Field
		Duration          respjson.Field
		GestureSourceType respjson.Field
		SessionID         respjson.Field
		TapCount          respjson.Field
		X                 respjson.Field
		Y                 respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataInputSynthesizeTapGesture) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataInputSynthesizeTapGesture) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `DOM.setFileInputFiles` arguments. Canonical input:
// `DOM.setFileInputFiles` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataDomSetFileInputFiles struct {
	// Number of files handed to the input. File paths are never captured.
	FileCount int64                         `json:"file_count" api:"required"`
	Method    constant.DomSetFileInputFiles `json:"method" default:"DOM.setFileInputFiles"`
	// Opaque backend DOM node identifier the command targeted.
	BackendNodeID int64 `json:"backend_node_id"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Opaque DOM node identifier the command targeted.
	NodeID int64 `json:"node_id"`
	// Opaque Runtime remote object identifier the command targeted. Clipped to 128
	// characters; a longer value is not a real identifier.
	ObjectID string `json:"object_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileCount     respjson.Field
		Method        respjson.Field
		BackendNodeID respjson.Field
		CommandID     respjson.Field
		ConnectionID  respjson.Field
		NodeID        respjson.Field
		ObjectID      respjson.Field
		SessionID     respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataDomSetFileInputFiles) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataDomSetFileInputFiles) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `DOM.focus` arguments. Canonical input: `DOM.focus` in
// devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataDomFocus struct {
	Method constant.DomFocus `json:"method" default:"DOM.focus"`
	// Opaque backend DOM node identifier the command targeted.
	BackendNodeID int64 `json:"backend_node_id"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Opaque DOM node identifier the command targeted.
	NodeID int64 `json:"node_id"`
	// Opaque Runtime remote object identifier the command targeted. Clipped to 128
	// characters; a longer value is not a real identifier.
	ObjectID string `json:"object_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method        respjson.Field
		BackendNodeID respjson.Field
		CommandID     respjson.Field
		ConnectionID  respjson.Field
		NodeID        respjson.Field
		ObjectID      respjson.Field
		SessionID     respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataDomFocus) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataDomFocus) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `DOM.scrollIntoViewIfNeeded` arguments. Canonical input:
// `DOM.scrollIntoViewIfNeeded` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataDomScrollIntoViewIfNeeded struct {
	Method constant.DomScrollIntoViewIfNeeded `json:"method" default:"DOM.scrollIntoViewIfNeeded"`
	// Opaque backend DOM node identifier the command targeted.
	BackendNodeID int64 `json:"backend_node_id"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Opaque DOM node identifier the command targeted.
	NodeID int64 `json:"node_id"`
	// Opaque Runtime remote object identifier the command targeted. Clipped to 128
	// characters; a longer value is not a real identifier.
	ObjectID string `json:"object_id"`
	// Height of the rect the command scrolled to.
	RectHeight float64 `json:"rect_height"`
	// Width of the rect the command scrolled to.
	RectWidth float64 `json:"rect_width"`
	// X offset of the rect the command scrolled to, relative to the node.
	RectX float64 `json:"rect_x"`
	// Y offset of the rect the command scrolled to, relative to the node.
	RectY float64 `json:"rect_y"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method        respjson.Field
		BackendNodeID respjson.Field
		CommandID     respjson.Field
		ConnectionID  respjson.Field
		NodeID        respjson.Field
		ObjectID      respjson.Field
		RectHeight    respjson.Field
		RectWidth     respjson.Field
		RectX         respjson.Field
		RectY         respjson.Field
		SessionID     respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataDomScrollIntoViewIfNeeded) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataDomScrollIntoViewIfNeeded) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.bringToFront` arguments. Canonical input: `Page.bringToFront` in
// devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageBringToFront struct {
	Method constant.PageBringToFront `json:"method" default:"Page.bringToFront"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageBringToFront) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageBringToFront) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.captureScreenshot` arguments. Canonical input:
// `Page.captureScreenshot` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageCaptureScreenshot struct {
	Method constant.PageCaptureScreenshot `json:"method" default:"Page.captureScreenshot"`
	// Whether the capture extended past the viewport.
	CaptureBeyondViewport bool `json:"capture_beyond_viewport"`
	// Clip region height in CSS pixels.
	ClipHeight float64 `json:"clip_height"`
	// Clip region page scale factor.
	ClipScale float64 `json:"clip_scale"`
	// Clip region width in CSS pixels.
	ClipWidth float64 `json:"clip_width"`
	// Clip region x offset in CSS pixels.
	ClipX float64 `json:"clip_x"`
	// Clip region y offset in CSS pixels.
	ClipY float64 `json:"clip_y"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Image format requested (`jpeg`, `png` or `webp`). A value the protocol does not
	// define is reported as `other`.
	//
	// Any of "jpeg", "png", "webp", "other".
	Format string `json:"format"`
	// Whether the capture was taken from the surface rather than the view.
	FromSurface bool `json:"from_surface"`
	// Whether encoding favored speed over size.
	OptimizeForSpeed bool `json:"optimize_for_speed"`
	// Compression quality, 0 to 100, for lossy formats.
	Quality int64 `json:"quality"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method                respjson.Field
		CaptureBeyondViewport respjson.Field
		ClipHeight            respjson.Field
		ClipScale             respjson.Field
		ClipWidth             respjson.Field
		ClipX                 respjson.Field
		ClipY                 respjson.Field
		CommandID             respjson.Field
		ConnectionID          respjson.Field
		Format                respjson.Field
		FromSurface           respjson.Field
		OptimizeForSpeed      respjson.Field
		Quality               respjson.Field
		SessionID             respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageCaptureScreenshot) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageCaptureScreenshot) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.captureSnapshot` arguments. Canonical input:
// `Page.captureSnapshot` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageCaptureSnapshot struct {
	Method constant.PageCaptureSnapshot `json:"method" default:"Page.captureSnapshot"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Snapshot format requested (`mhtml`). A value the protocol does not define is
	// reported as `other`.
	//
	// Any of "mhtml", "other".
	Format string `json:"format"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		Format       respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageCaptureSnapshot) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageCaptureSnapshot) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.handleJavaScriptDialog` arguments. Canonical input:
// `Page.handleJavaScriptDialog` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageHandleJavaScriptDialog struct {
	// Whether the dialog was accepted or dismissed.
	Accept bool                                `json:"accept" api:"required"`
	Method constant.PageHandleJavaScriptDialog `json:"method" default:"Page.handleJavaScriptDialog"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Number of characters entered into a prompt dialog. The text itself is never
	// captured.
	PromptTextLength int64 `json:"prompt_text_length"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Accept           respjson.Field
		Method           respjson.Field
		CommandID        respjson.Field
		ConnectionID     respjson.Field
		PromptTextLength respjson.Field
		SessionID        respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageHandleJavaScriptDialog) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageHandleJavaScriptDialog) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.navigate` arguments. Canonical input: `Page.navigate` in
// devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageNavigate struct {
	Method constant.PageNavigate `json:"method" default:"Page.navigate"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Opaque frame identifier. Clipped to 128 characters; a longer value is not a real
	// identifier.
	FrameID string `json:"frame_id"`
	// Referrer policy named by the command. A value the protocol does not define is
	// reported as `other`.
	//
	// Any of "noReferrer", "noReferrerWhenDowngrade", "origin",
	// "originWhenCrossOrigin", "sameOrigin", "strictOrigin",
	// "strictOriginWhenCrossOrigin", "unsafeUrl", "other".
	ReferrerPolicy string `json:"referrer_policy"`
	// Whether the command carried a referrer. The referrer itself is never captured.
	ReferrerPresent bool `json:"referrer_present"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Navigation reason reported by the caller (e.g. `link`, `typed`, `reload`). A
	// value the protocol does not define is reported as `other`.
	//
	// Any of "link", "typed", "address_bar", "auto_bookmark", "auto_subframe",
	// "manual_subframe", "generated", "auto_toplevel", "form_submit", "reload",
	// "keyword", "keyword_generated", "other".
	TransitionType string `json:"transition_type"`
	// Scheme of the destination URL (e.g. `https`, `about`, `data`). The rest of the
	// URL is never captured.
	URLScheme string `json:"url_scheme"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method          respjson.Field
		CommandID       respjson.Field
		ConnectionID    respjson.Field
		FrameID         respjson.Field
		ReferrerPolicy  respjson.Field
		ReferrerPresent respjson.Field
		SessionID       respjson.Field
		TransitionType  respjson.Field
		URLScheme       respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageNavigate) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageNavigate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.navigateToHistoryEntry` arguments. Canonical input:
// `Page.navigateToHistoryEntry` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageNavigateToHistoryEntry struct {
	// History entry the command navigated to.
	EntryID int64                               `json:"entry_id" api:"required"`
	Method  constant.PageNavigateToHistoryEntry `json:"method" default:"Page.navigateToHistoryEntry"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EntryID      respjson.Field
		Method       respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageNavigateToHistoryEntry) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageNavigateToHistoryEntry) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.reload` arguments. Canonical input: `Page.reload` in
// devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageReload struct {
	Method constant.PageReload `json:"method" default:"Page.reload"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Whether the reload bypassed the cache.
	IgnoreCache bool `json:"ignore_cache"`
	// Opaque document loader identifier. Clipped to 128 characters; a longer value is
	// not a real identifier.
	LoaderID string `json:"loader_id"`
	// Number of characters in the injected script.
	ScriptLength int64 `json:"script_length"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		IgnoreCache  respjson.Field
		LoaderID     respjson.Field
		ScriptLength respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageReload) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageReload) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.printToPDF` arguments. Canonical input: `Page.printToPDF` in
// devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPagePrintToPdf struct {
	Method constant.PagePrintToPdf `json:"method" default:"Page.printToPDF"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Whether a header and footer were rendered.
	DisplayHeaderFooter bool `json:"display_header_footer"`
	// Whether a footer template was supplied. The template itself is never captured.
	FooterTemplatePresent bool `json:"footer_template_present"`
	// Whether a document outline was embedded.
	GenerateDocumentOutline bool `json:"generate_document_outline"`
	// Whether a tagged (accessible) PDF was requested.
	GenerateTaggedPdf bool `json:"generate_tagged_pdf"`
	// Whether a header template was supplied. The template itself is never captured.
	HeaderTemplatePresent bool `json:"header_template_present"`
	// Whether the page was laid out in landscape.
	Landscape bool `json:"landscape"`
	// Bottom margin in inches.
	MarginBottom float64 `json:"margin_bottom"`
	// Left margin in inches.
	MarginLeft float64 `json:"margin_left"`
	// Right margin in inches.
	MarginRight float64 `json:"margin_right"`
	// Top margin in inches.
	MarginTop float64 `json:"margin_top"`
	// Whether a page range was supplied.
	PageRangesPresent bool `json:"page_ranges_present"`
	// Paper height in inches.
	PaperHeight float64 `json:"paper_height"`
	// Paper width in inches.
	PaperWidth float64 `json:"paper_width"`
	// Whether the CSS page size was preferred over the paper size.
	PreferCssPageSize bool `json:"prefer_css_page_size"`
	// Whether background graphics were printed.
	PrintBackground bool `json:"print_background"`
	// Page render scale.
	Scale float64 `json:"scale"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// How the PDF was returned (`ReturnAsBase64` or `ReturnAsStream`). A value the
	// protocol does not define is reported as `other`.
	//
	// Any of "ReturnAsBase64", "ReturnAsStream", "other".
	TransferMode string `json:"transfer_mode"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method                  respjson.Field
		CommandID               respjson.Field
		ConnectionID            respjson.Field
		DisplayHeaderFooter     respjson.Field
		FooterTemplatePresent   respjson.Field
		GenerateDocumentOutline respjson.Field
		GenerateTaggedPdf       respjson.Field
		HeaderTemplatePresent   respjson.Field
		Landscape               respjson.Field
		MarginBottom            respjson.Field
		MarginLeft              respjson.Field
		MarginRight             respjson.Field
		MarginTop               respjson.Field
		PageRangesPresent       respjson.Field
		PaperHeight             respjson.Field
		PaperWidth              respjson.Field
		PreferCssPageSize       respjson.Field
		PrintBackground         respjson.Field
		Scale                   respjson.Field
		SessionID               respjson.Field
		TransferMode            respjson.Field
		ExtraFields             map[string]respjson.Field
		raw                     string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPagePrintToPdf) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPagePrintToPdf) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.startScreencast` arguments. Canonical input:
// `Page.startScreencast` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageStartScreencast struct {
	Method constant.PageStartScreencast `json:"method" default:"Page.startScreencast"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Frame sampling interval.
	EveryNthFrame int64 `json:"every_nth_frame"`
	// Frame format requested (`jpeg` or `png`). A value the protocol does not define
	// is reported as `other`.
	//
	// Any of "jpeg", "png", "other".
	Format string `json:"format"`
	// Maximum frame height in pixels.
	MaxHeight int64 `json:"max_height"`
	// Maximum frame width in pixels.
	MaxWidth int64 `json:"max_width"`
	// Compression quality, 0 to 100.
	Quality int64 `json:"quality"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method        respjson.Field
		CommandID     respjson.Field
		ConnectionID  respjson.Field
		EveryNthFrame respjson.Field
		Format        respjson.Field
		MaxHeight     respjson.Field
		MaxWidth      respjson.Field
		Quality       respjson.Field
		SessionID     respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageStartScreencast) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageStartScreencast) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.stopScreencast` arguments. Canonical input:
// `Page.stopScreencast` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageStopScreencast struct {
	Method constant.PageStopScreencast `json:"method" default:"Page.stopScreencast"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageStopScreencast) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageStopScreencast) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.stopLoading` arguments. Canonical input: `Page.stopLoading` in
// devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageStopLoading struct {
	Method constant.PageStopLoading `json:"method" default:"Page.stopLoading"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageStopLoading) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageStopLoading) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.close` arguments. Canonical input: `Page.close` in
// devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageClose struct {
	Method constant.PageClose `json:"method" default:"Page.close"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageClose) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageClose) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Page.setWebLifecycleState` arguments. Canonical input:
// `Page.setWebLifecycleState` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataPageSetWebLifecycleState struct {
	Method constant.PageSetWebLifecycleState `json:"method" default:"Page.setWebLifecycleState"`
	// Lifecycle state applied (`frozen` or `active`). A value the protocol does not
	// define is reported as `other`.
	//
	// Any of "frozen", "active", "other".
	State string `json:"state" api:"required"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		State        respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataPageSetWebLifecycleState) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataPageSetWebLifecycleState) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Target.activateTarget` arguments. Canonical input:
// `Target.activateTarget` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataTargetActivateTarget struct {
	Method constant.TargetActivateTarget `json:"method" default:"Target.activateTarget"`
	// Opaque target identifier. Clipped to 128 characters; a longer value is not a
	// real identifier.
	TargetID string `json:"target_id" api:"required"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		TargetID     respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataTargetActivateTarget) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataTargetActivateTarget) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Target.closeTarget` arguments. Canonical input: `Target.closeTarget`
// in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataTargetCloseTarget struct {
	Method constant.TargetCloseTarget `json:"method" default:"Target.closeTarget"`
	// Opaque target identifier. Clipped to 128 characters; a longer value is not a
	// real identifier.
	TargetID string `json:"target_id" api:"required"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		TargetID     respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataTargetCloseTarget) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataTargetCloseTarget) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Target.createTarget` arguments. Canonical input:
// `Target.createTarget` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataTargetCreateTarget struct {
	Method constant.TargetCreateTarget `json:"method" default:"Target.createTarget"`
	// Whether the target was created in the background.
	Background bool `json:"background"`
	// Opaque browser context identifier. Clipped to 128 characters; a longer value is
	// not a real identifier.
	BrowserContextID string `json:"browser_context_id"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Whether BeginFrame control was enabled (headless only).
	EnableBeginFrameControl bool `json:"enable_begin_frame_control"`
	// Whether the new target was focused.
	Focus bool `json:"focus"`
	// Whether a tab target rather than a page target was created.
	ForTab bool `json:"for_tab"`
	// Window height in DIP.
	Height int64 `json:"height"`
	// Whether the target was created hidden.
	Hidden bool `json:"hidden"`
	// Window x position in screen coordinates.
	Left int64 `json:"left"`
	// Whether a new window was requested.
	NewWindow bool `json:"new_window"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Window y position in screen coordinates.
	Top int64 `json:"top"`
	// Scheme of the destination URL (e.g. `https`, `about`, `data`). The rest of the
	// URL is never captured.
	URLScheme string `json:"url_scheme"`
	// Window width in DIP.
	Width int64 `json:"width"`
	// Window state requested (`normal`, `minimized`, `maximized`, `fullscreen`). A
	// value the protocol does not define is reported as `other`.
	//
	// Any of "normal", "minimized", "maximized", "fullscreen", "other".
	WindowState string `json:"window_state"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method                  respjson.Field
		Background              respjson.Field
		BrowserContextID        respjson.Field
		CommandID               respjson.Field
		ConnectionID            respjson.Field
		EnableBeginFrameControl respjson.Field
		Focus                   respjson.Field
		ForTab                  respjson.Field
		Height                  respjson.Field
		Hidden                  respjson.Field
		Left                    respjson.Field
		NewWindow               respjson.Field
		SessionID               respjson.Field
		Top                     respjson.Field
		URLScheme               respjson.Field
		Width                   respjson.Field
		WindowState             respjson.Field
		ExtraFields             map[string]respjson.Field
		raw                     string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataTargetCreateTarget) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataTargetCreateTarget) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Target.createBrowserContext` arguments. Canonical input:
// `Target.createBrowserContext` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataTargetCreateBrowserContext struct {
	Method constant.TargetCreateBrowserContext `json:"method" default:"Target.createBrowserContext"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Whether the context is disposed when the debugging session detaches.
	DisposeOnDetach bool `json:"dispose_on_detach"`
	// Whether a proxy bypass list was configured.
	ProxyBypassListPresent bool `json:"proxy_bypass_list_present"`
	// Whether a proxy was configured. The proxy address is never captured.
	ProxyServerPresent bool `json:"proxy_server_present"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Number of origins granted universal network access. The origins themselves are
	// never captured.
	UniversalNetworkAccessOriginCount int64 `json:"universal_network_access_origin_count"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method                            respjson.Field
		CommandID                         respjson.Field
		ConnectionID                      respjson.Field
		DisposeOnDetach                   respjson.Field
		ProxyBypassListPresent            respjson.Field
		ProxyServerPresent                respjson.Field
		SessionID                         respjson.Field
		UniversalNetworkAccessOriginCount respjson.Field
		ExtraFields                       map[string]respjson.Field
		raw                               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataTargetCreateBrowserContext) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataTargetCreateBrowserContext) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Target.disposeBrowserContext` arguments. Canonical input:
// `Target.disposeBrowserContext` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataTargetDisposeBrowserContext struct {
	// Opaque browser context identifier. Clipped to 128 characters; a longer value is
	// not a real identifier.
	BrowserContextID string                               `json:"browser_context_id" api:"required"`
	Method           constant.TargetDisposeBrowserContext `json:"method" default:"Target.disposeBrowserContext"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BrowserContextID respjson.Field
		Method           respjson.Field
		CommandID        respjson.Field
		ConnectionID     respjson.Field
		SessionID        respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataTargetDisposeBrowserContext) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataTargetDisposeBrowserContext) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Target.openDevTools` arguments. Canonical input:
// `Target.openDevTools` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataTargetOpenDevTools struct {
	Method constant.TargetOpenDevTools `json:"method" default:"Target.openDevTools"`
	// Opaque target identifier. Clipped to 128 characters; a longer value is not a
	// real identifier.
	TargetID string `json:"target_id" api:"required"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// DevTools panel opened. Clipped to 128 characters; a longer value is not a real
	// identifier.
	PanelID string `json:"panel_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		TargetID     respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		PanelID      respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataTargetOpenDevTools) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataTargetOpenDevTools) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Browser.cancelDownload` arguments. Canonical input:
// `Browser.cancelDownload` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataBrowserCancelDownload struct {
	// Opaque identifier of the download that was cancelled. Clipped to 128 characters;
	// a longer value is not a real identifier.
	DownloadGuid string                         `json:"download_guid" api:"required"`
	Method       constant.BrowserCancelDownload `json:"method" default:"Browser.cancelDownload"`
	// Opaque browser context identifier. Clipped to 128 characters; a longer value is
	// not a real identifier.
	BrowserContextID string `json:"browser_context_id"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DownloadGuid     respjson.Field
		Method           respjson.Field
		BrowserContextID respjson.Field
		CommandID        respjson.Field
		ConnectionID     respjson.Field
		SessionID        respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataBrowserCancelDownload) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataBrowserCancelDownload) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Browser.close` arguments. Canonical input: `Browser.close` in
// devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataBrowserClose struct {
	Method constant.BrowserClose `json:"method" default:"Browser.close"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		SessionID    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataBrowserClose) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataBrowserClose) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Browser.setWindowBounds` arguments. Canonical input:
// `Browser.setWindowBounds` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataBrowserSetWindowBounds struct {
	Method constant.BrowserSetWindowBounds `json:"method" default:"Browser.setWindowBounds"`
	// Browser window identifier.
	WindowID int64 `json:"window_id" api:"required"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Window height in DIP.
	Height int64 `json:"height"`
	// Window x position in screen coordinates.
	Left int64 `json:"left"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Window y position in screen coordinates.
	Top int64 `json:"top"`
	// Window width in DIP.
	Width int64 `json:"width"`
	// Window state requested (`normal`, `minimized`, `maximized`, `fullscreen`). A
	// value the protocol does not define is reported as `other`.
	//
	// Any of "normal", "minimized", "maximized", "fullscreen", "other".
	WindowState string `json:"window_state"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		WindowID     respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		Height       respjson.Field
		Left         respjson.Field
		SessionID    respjson.Field
		Top          respjson.Field
		Width        respjson.Field
		WindowState  respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataBrowserSetWindowBounds) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataBrowserSetWindowBounds) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Browser.setContentsSize` arguments. Canonical input:
// `Browser.setContentsSize` in devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataBrowserSetContentsSize struct {
	Method constant.BrowserSetContentsSize `json:"method" default:"Browser.setContentsSize"`
	// Browser window identifier.
	WindowID int64 `json:"window_id" api:"required"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Contents height in DIP.
	Height int64 `json:"height"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// Contents width in DIP.
	Width int64 `json:"width"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Method       respjson.Field
		WindowID     respjson.Field
		CommandID    respjson.Field
		ConnectionID respjson.Field
		Height       respjson.Field
		SessionID    respjson.Field
		Width        respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataBrowserSetContentsSize) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataBrowserSetContentsSize) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sanitized `Autofill.trigger` arguments. Canonical input: `Autofill.trigger` in
// devtools-protocol@2d019e73, pinned at
// https://github.com/ChromeDevTools/devtools-protocol/blob/2d019e73eb371d1d6985d26d395d78bd8f8a22ba/json/browser_protocol.json.
// Every argument of this command has a retained or redacted decision in
// lib/devtoolsproxy/testdata/cdp_arguments.yaml.
type BrowserCdpCommandEventDataAutofillTrigger struct {
	// Opaque backend node identifier of the field that was autofilled.
	FieldID int64                    `json:"field_id" api:"required"`
	Method  constant.AutofillTrigger `json:"method" default:"Autofill.trigger"`
	// Number of address fields the command filled. Their names and values are never
	// captured.
	AddressFieldCount int64 `json:"address_field_count"`
	// The command's JSON-RPC id, so the command can be joined to the result the
	// browser returned for it. Absent when the client sent none.
	CommandID int64 `json:"command_id"`
	// Identifies the CDP proxy connection the command arrived on, matching
	// `cdp_connect` and `cdp_disconnect`. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Opaque frame identifier. Clipped to 128 characters; a longer value is not a real
	// identifier.
	FrameID string `json:"frame_id"`
	// What was filled: `card` or `address`. The values themselves are never captured.
	//
	// Any of "card", "address".
	Mode string `json:"mode"`
	// CDP session identifier the command was addressed to. Absent for browser-level
	// commands. Clipped to 128 characters.
	SessionID string `json:"session_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FieldID           respjson.Field
		Method            respjson.Field
		AddressFieldCount respjson.Field
		CommandID         respjson.Field
		ConnectionID      respjson.Field
		FrameID           respjson.Field
		Mode              respjson.Field
		SessionID         respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpCommandEventDataAutofillTrigger) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpCommandEventDataAutofillTrigger) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser-control CDP method the proxy reports. The set covers the commands an
// agent drives the browser with; configuration, DOM and Runtime bookkeeping, and
// Chrome-specific UI commands are outside it. Canonical definitions:
// devtools-protocol@2d019e73.
type BrowserCdpCommandMethod string

const (
	BrowserCdpCommandMethodInputDispatchMouseEvent         BrowserCdpCommandMethod = "Input.dispatchMouseEvent"
	BrowserCdpCommandMethodInputDispatchKeyEvent           BrowserCdpCommandMethod = "Input.dispatchKeyEvent"
	BrowserCdpCommandMethodInputInsertText                 BrowserCdpCommandMethod = "Input.insertText"
	BrowserCdpCommandMethodInputImeSetComposition          BrowserCdpCommandMethod = "Input.imeSetComposition"
	BrowserCdpCommandMethodInputDispatchTouchEvent         BrowserCdpCommandMethod = "Input.dispatchTouchEvent"
	BrowserCdpCommandMethodInputDispatchDragEvent          BrowserCdpCommandMethod = "Input.dispatchDragEvent"
	BrowserCdpCommandMethodInputCancelDragging             BrowserCdpCommandMethod = "Input.cancelDragging"
	BrowserCdpCommandMethodInputEmulateTouchFromMouseEvent BrowserCdpCommandMethod = "Input.emulateTouchFromMouseEvent"
	BrowserCdpCommandMethodInputSynthesizePinchGesture     BrowserCdpCommandMethod = "Input.synthesizePinchGesture"
	BrowserCdpCommandMethodInputSynthesizeScrollGesture    BrowserCdpCommandMethod = "Input.synthesizeScrollGesture"
	BrowserCdpCommandMethodInputSynthesizeTapGesture       BrowserCdpCommandMethod = "Input.synthesizeTapGesture"
	BrowserCdpCommandMethodDomSetFileInputFiles            BrowserCdpCommandMethod = "DOM.setFileInputFiles"
	BrowserCdpCommandMethodDomFocus                        BrowserCdpCommandMethod = "DOM.focus"
	BrowserCdpCommandMethodDomScrollIntoViewIfNeeded       BrowserCdpCommandMethod = "DOM.scrollIntoViewIfNeeded"
	BrowserCdpCommandMethodPageBringToFront                BrowserCdpCommandMethod = "Page.bringToFront"
	BrowserCdpCommandMethodPageCaptureScreenshot           BrowserCdpCommandMethod = "Page.captureScreenshot"
	BrowserCdpCommandMethodPageCaptureSnapshot             BrowserCdpCommandMethod = "Page.captureSnapshot"
	BrowserCdpCommandMethodPageHandleJavaScriptDialog      BrowserCdpCommandMethod = "Page.handleJavaScriptDialog"
	BrowserCdpCommandMethodPageNavigate                    BrowserCdpCommandMethod = "Page.navigate"
	BrowserCdpCommandMethodPageNavigateToHistoryEntry      BrowserCdpCommandMethod = "Page.navigateToHistoryEntry"
	BrowserCdpCommandMethodPageReload                      BrowserCdpCommandMethod = "Page.reload"
	BrowserCdpCommandMethodPagePrintToPdf                  BrowserCdpCommandMethod = "Page.printToPDF"
	BrowserCdpCommandMethodPageStartScreencast             BrowserCdpCommandMethod = "Page.startScreencast"
	BrowserCdpCommandMethodPageStopScreencast              BrowserCdpCommandMethod = "Page.stopScreencast"
	BrowserCdpCommandMethodPageStopLoading                 BrowserCdpCommandMethod = "Page.stopLoading"
	BrowserCdpCommandMethodPageClose                       BrowserCdpCommandMethod = "Page.close"
	BrowserCdpCommandMethodPageSetWebLifecycleState        BrowserCdpCommandMethod = "Page.setWebLifecycleState"
	BrowserCdpCommandMethodTargetActivateTarget            BrowserCdpCommandMethod = "Target.activateTarget"
	BrowserCdpCommandMethodTargetCloseTarget               BrowserCdpCommandMethod = "Target.closeTarget"
	BrowserCdpCommandMethodTargetCreateTarget              BrowserCdpCommandMethod = "Target.createTarget"
	BrowserCdpCommandMethodTargetCreateBrowserContext      BrowserCdpCommandMethod = "Target.createBrowserContext"
	BrowserCdpCommandMethodTargetDisposeBrowserContext     BrowserCdpCommandMethod = "Target.disposeBrowserContext"
	BrowserCdpCommandMethodTargetOpenDevTools              BrowserCdpCommandMethod = "Target.openDevTools"
	BrowserCdpCommandMethodBrowserCancelDownload           BrowserCdpCommandMethod = "Browser.cancelDownload"
	BrowserCdpCommandMethodBrowserClose                    BrowserCdpCommandMethod = "Browser.close"
	BrowserCdpCommandMethodBrowserSetWindowBounds          BrowserCdpCommandMethod = "Browser.setWindowBounds"
	BrowserCdpCommandMethodBrowserSetContentsSize          BrowserCdpCommandMethod = "Browser.setContentsSize"
	BrowserCdpCommandMethodAutofillTrigger                 BrowserCdpCommandMethod = "Autofill.trigger"
)

// An external client (e.g. customer SDK, Playwright, Puppeteer) connected to the
// CDP WebSocket proxy on this VM.
type BrowserCdpConnectEvent struct {
	Category constant.Connection `json:"category" default:"connection"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                      `json:"ts" api:"required"`
	Type constant.CdpConnect        `json:"type" default:"cdp_connect"`
	Data BrowserCdpConnectEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpConnectEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpConnectEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserCdpConnectEventData struct {
	// Identifies this CDP proxy connection, matching the connection_id on the
	// cdp_command events that arrived on it. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ConnectionID respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpConnectEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpConnectEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An external client disconnected from the CDP WebSocket proxy on this VM. Pair
// with the immediately preceding cdp_connect on the same stream.
type BrowserCdpDisconnectEvent struct {
	Category constant.Connection `json:"category" default:"connection"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                         `json:"ts" api:"required"`
	Type constant.CdpDisconnect        `json:"type" default:"cdp_disconnect"`
	Data BrowserCdpDisconnectEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpDisconnectEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpDisconnectEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserCdpDisconnectEventData struct {
	// Wall-clock duration of the connection in milliseconds.
	DurationMs float64 `json:"duration_ms" api:"required"`
	// Number of CDP messages relayed across the connection in either direction.
	MessageCount int64 `json:"message_count" api:"required"`
	// Why the connection ended. client_close: the client initiated the close.
	// upstream_changed: Chromium restarted mid-session and the proxy tore down so the
	// client could reconnect against the new upstream. upstream_error: upstream dial
	// or message pump errored. context_cancelled: the request context was cancelled
	// (typically server shutdown).
	//
	// Any of "client_close", "upstream_changed", "upstream_error",
	// "context_cancelled".
	Reason string `json:"reason" api:"required"`
	// Identifies this CDP proxy connection, matching the connection_id on the
	// cdp_command events that arrived on it. Two clients driving the same browser are
	// told apart by this.
	ConnectionID string `json:"connection_id"`
	// Number of forwarded client frames the classifier never saw, because it could not
	// keep up or because classification failed. An upper bound on lost commands rather
	// than a count: a saturated queue turns away whatever arrives next, which may be
	// library traffic that would have produced no event. Telemetry loss only; every
	// command was still relayed to the browser. Absent on events from a browser image
	// predating the field, which is not the same as zero.
	TelemetryDropped int64 `json:"telemetry_dropped"`
	// Number of forwarded client commands that produced no cdp_command event because
	// their method is listed in control.cdp.excluded_methods. Configuration rather
	// than loss, so it is counted apart from telemetry_dropped.
	TelemetryExcluded int64 `json:"telemetry_excluded"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DurationMs        respjson.Field
		MessageCount      respjson.Field
		Reason            respjson.Field
		ConnectionID      respjson.Field
		TelemetryDropped  respjson.Field
		TelemetryExcluded respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserCdpDisconnectEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserCdpDisconnectEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser console error or uncaught JavaScript exception event. Emitted from two
// distinct CDP sources with different data shapes. Runtime.consoleAPICalled
// (console.error calls) produces level, text, args, and stack_trace.
// Runtime.exceptionThrown (uncaught exceptions) produces text, line, column,
// source_url, and stack_trace. Fields not applicable to the source are absent.
type BrowserConsoleErrorEvent struct {
	Category constant.Console `json:"category" default:"console"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                 `json:"ts" api:"required"`
	Type constant.ConsoleError `json:"type" default:"console_error"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserConsoleErrorEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserConsoleErrorEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserConsoleErrorEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserConsoleErrorEventData struct {
	// Human-readable error text, as the browser console would display it. For
	// console.error() calls, the first argument coerced to a string. For uncaught
	// exceptions, the prefix and error message, e.g. "Uncaught Error: boom" or
	// "Uncaught (in promise) TypeError: x is not a function".
	Text string `json:"text" api:"required"`
	// All console arguments coerced to strings. Present only when sourced from
	// Runtime.consoleAPICalled.
	Args []string `json:"args"`
	// Column number in the script where the exception was thrown. Present only when
	// sourced from Runtime.exceptionThrown.
	Column int64 `json:"column"`
	// CDP console type value, always "error". Present only when sourced from
	// Runtime.consoleAPICalled.
	Level string `json:"level"`
	// Line number in the script where the exception was thrown. Present only when
	// sourced from Runtime.exceptionThrown.
	Line int64 `json:"line"`
	// URL of the script file that threw the exception. Present only when sourced from
	// Runtime.exceptionThrown.
	SourceURL string `json:"source_url"`
	// CDP Runtime.StackTrace representing the JavaScript call stack at the time of an
	// event. Fields use CDP naming conventions rather than snake_case to match the
	// Chrome DevTools Protocol wire format.
	StackTrace BrowserCallStack `json:"stack_trace"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Args        respjson.Field
		Column      respjson.Field
		Level       respjson.Field
		Line        respjson.Field
		SourceURL   respjson.Field
		StackTrace  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserConsoleErrorEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserConsoleErrorEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser console log event (console.log, console.info, console.warn, etc.).
type BrowserConsoleLogEvent struct {
	Category constant.Console `json:"category" default:"console"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64               `json:"ts" api:"required"`
	Type constant.ConsoleLog `json:"type" default:"console_log"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserConsoleLogEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserConsoleLogEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserConsoleLogEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserConsoleLogEventData struct {
	// All console arguments coerced to strings.
	Args []string `json:"args"`
	// CDP Runtime.consoleAPICalled type, passed through unfiltered from Chrome. error
	// is routed to console_error events instead; all other CDP console types appear
	// here. See CDP spec for the full enum.
	Level string `json:"level"`
	// CDP Runtime.StackTrace representing the JavaScript call stack at the time of an
	// event. Fields use CDP naming conventions rather than snake_case to match the
	// Chrome DevTools Protocol wire format.
	StackTrace BrowserCallStack `json:"stack_trace"`
	// First console argument coerced to string.
	Text string `json:"text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Args        respjson.Field
		Level       respjson.Field
		StackTrace  respjson.Field
		Text        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserConsoleLogEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserConsoleLogEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserEventContext struct {
	// CDP frame identifier within the target.
	FrameID string `json:"frame_id"`
	// CDP document loader identifier, reset on each navigation.
	LoaderID string `json:"loader_id"`
	// Monotonically increasing navigation sequence number, incremented on each
	// top-level navigation within the target.
	NavSeq int64 `json:"nav_seq"`
	// CDP session identifier for the target connection.
	SessionID string `json:"session_id"`
	// Browser target identifier (stable across navigations within a tab).
	TargetID string `json:"target_id"`
	// CDP target type of the page that produced the event.
	//
	// Any of "page", "background_page", "service_worker", "shared_worker", "other".
	TargetType BrowserEventContextTargetType `json:"target_type"`
	// URL relevant to this event — page URL for navigation and page events, request
	// URL for network events.
	URL string `json:"url"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FrameID     respjson.Field
		LoaderID    respjson.Field
		NavSeq      respjson.Field
		SessionID   respjson.Field
		TargetID    respjson.Field
		TargetType  respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserEventContext) RawJSON() string { return r.JSON.raw }
func (r *BrowserEventContext) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CDP target type of the page that produced the event.
type BrowserEventContextTargetType string

const (
	BrowserEventContextTargetTypePage           BrowserEventContextTargetType = "page"
	BrowserEventContextTargetTypeBackgroundPage BrowserEventContextTargetType = "background_page"
	BrowserEventContextTargetTypeServiceWorker  BrowserEventContextTargetType = "service_worker"
	BrowserEventContextTargetTypeSharedWorker   BrowserEventContextTargetType = "shared_worker"
	BrowserEventContextTargetTypeOther          BrowserEventContextTargetType = "other"
)

// Provenance metadata identifying which producer emitted the event.
type BrowserEventSource struct {
	// Event producer. cdp: Chrome DevTools Protocol events from the browser.
	// kernel_api: Kernel API server. extension: injected Chrome extension.
	// local_process: system process running alongside the browser.
	//
	// Any of "cdp", "kernel_api", "extension", "local_process".
	Kind BrowserEventSourceKind `json:"kind" api:"required"`
	// Producer-specific event name (e.g. Runtime.consoleAPICalled for CDP-sourced
	// console events, Runtime.exceptionThrown for uncaught exceptions).
	Event string `json:"event"`
	// Producer-specific context (e.g. CDP target/session/frame IDs).
	Metadata map[string]string `json:"metadata"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Kind        respjson.Field
		Event       respjson.Field
		Metadata    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserEventSource) RawJSON() string { return r.JSON.raw }
func (r *BrowserEventSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event producer. cdp: Chrome DevTools Protocol events from the browser.
// kernel_api: Kernel API server. extension: injected Chrome extension.
// local_process: system process running alongside the browser.
type BrowserEventSourceKind string

const (
	BrowserEventSourceKindCdp          BrowserEventSourceKind = "cdp"
	BrowserEventSourceKindKernelAPI    BrowserEventSourceKind = "kernel_api"
	BrowserEventSourceKindExtension    BrowserEventSourceKind = "extension"
	BrowserEventSourceKindLocalProcess BrowserEventSourceKind = "local_process"
)

type BrowserHTTPHeaders map[string]any

// A browser user click event captured via injected page script.
type BrowserInteractionClickEvent struct {
	Category constant.Interaction `json:"category" default:"interaction"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                     `json:"ts" api:"required"`
	Type constant.InteractionClick `json:"type" default:"interaction_click"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserInteractionClickEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserInteractionClickEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserInteractionClickEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserInteractionClickEventData struct {
	// CSS selector path to the clicked element.
	Selector string `json:"selector"`
	// HTML tag name of the clicked element in uppercase (e.g. BUTTON, A, DIV).
	Tag string `json:"tag"`
	// Visible text content of the clicked element, trimmed.
	Text string `json:"text"`
	// Viewport x-coordinate of the click in CSS pixels.
	X int64 `json:"x"`
	// Viewport y-coordinate of the click in CSS pixels.
	Y int64 `json:"y"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Selector    respjson.Field
		Tag         respjson.Field
		Text        respjson.Field
		X           respjson.Field
		Y           respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserInteractionClickEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserInteractionClickEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser keyboard event captured via injected page script.
type BrowserInteractionKeyEvent struct {
	Category constant.Interaction `json:"category" default:"interaction"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                   `json:"ts" api:"required"`
	Type constant.InteractionKey `json:"type" default:"interaction_key"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserInteractionKeyEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserInteractionKeyEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserInteractionKeyEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserInteractionKeyEventData struct {
	// Key value from the KeyboardEvent (e.g. Enter, Backspace, a).
	Key string `json:"key"`
	// CSS selector path to the element that had focus when the key was pressed.
	Selector string `json:"selector"`
	// HTML tag name of the focused element in uppercase (e.g. INPUT, TEXTAREA, DIV).
	Tag string `json:"tag"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Key         respjson.Field
		Selector    respjson.Field
		Tag         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserInteractionKeyEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserInteractionKeyEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser scroll settled event emitted after scroll position stops changing,
// captured via injected page script.
type BrowserInteractionScrollSettledEvent struct {
	Category constant.Interaction `json:"category" default:"interaction"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                             `json:"ts" api:"required"`
	Type constant.InteractionScrollSettled `json:"type" default:"interaction_scroll_settled"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserInteractionScrollSettledEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserInteractionScrollSettledEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserInteractionScrollSettledEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserInteractionScrollSettledEventData struct {
	// Scroll x-position at the start of the scroll gesture in CSS pixels.
	FromX int64 `json:"from_x"`
	// Scroll y-position at the start of the scroll gesture in CSS pixels.
	FromY int64 `json:"from_y"`
	// CSS selector path to the scrolled element.
	TargetSelector string `json:"target_selector"`
	// Final scroll x-position after the gesture settled in CSS pixels.
	ToX int64 `json:"to_x"`
	// Final scroll y-position after the gesture settled in CSS pixels.
	ToY int64 `json:"to_y"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FromX          respjson.Field
		FromY          respjson.Field
		TargetSelector respjson.Field
		ToX            respjson.Field
		ToY            respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserInteractionScrollSettledEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserInteractionScrollSettledEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A live view client connected to the headful browser's WebRTC server. Headful
// only; not emitted for headless images.
type BrowserLiveViewConnectEvent struct {
	Category constant.Connection `json:"category" default:"connection"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                           `json:"ts" api:"required"`
	Type constant.LiveViewConnect        `json:"type" default:"live_view_connect"`
	Data BrowserLiveViewConnectEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserLiveViewConnectEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserLiveViewConnectEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserLiveViewConnectEventData struct {
	// Live view session identifier. Stable across reconnects, so a transient network
	// blip can emit two events with the same session_id.
	SessionID string `json:"session_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SessionID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserLiveViewConnectEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserLiveViewConnectEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A live view client disconnected from the headful browser's WebRTC server. Pair
// with live_view_connect by session_id.
type BrowserLiveViewDisconnectEvent struct {
	Category constant.Connection `json:"category" default:"connection"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                              `json:"ts" api:"required"`
	Type constant.LiveViewDisconnect        `json:"type" default:"live_view_disconnect"`
	Data BrowserLiveViewDisconnectEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserLiveViewDisconnectEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserLiveViewDisconnectEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserLiveViewDisconnectEventData struct {
	// Wall-clock duration of the connection in milliseconds.
	DurationMs float64 `json:"duration_ms" api:"required"`
	// Live view session identifier; matches the corresponding live_view_connect event.
	SessionID string `json:"session_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DurationMs  respjson.Field
		SessionID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserLiveViewDisconnectEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserLiveViewDisconnectEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The CDP connection to Chrome was lost. Telemetry events may be dropped until
// monitor_reconnected arrives. Treat any in-progress computed state (network_idle,
// page_layout_settled) as unreliable until then.
type BrowserMonitorDisconnectedEvent struct {
	Category constant.Monitor `json:"category" default:"monitor"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                               `json:"ts" api:"required"`
	Type constant.MonitorDisconnected        `json:"type" default:"monitor_disconnected"`
	Data BrowserMonitorDisconnectedEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserMonitorDisconnectedEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserMonitorDisconnectedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserMonitorDisconnectedEventData struct {
	// Reason for the disconnection. chrome_restarted: Chrome process restarted.
	//
	// Any of "chrome_restarted".
	Reason string `json:"reason"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Reason      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserMonitorDisconnectedEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserMonitorDisconnectedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The CDP session could not be initialized.
type BrowserMonitorInitFailedEvent struct {
	Category constant.Monitor `json:"category" default:"monitor"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                             `json:"ts" api:"required"`
	Type constant.MonitorInitFailed        `json:"type" default:"monitor_init_failed"`
	Data BrowserMonitorInitFailedEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserMonitorInitFailedEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserMonitorInitFailedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserMonitorInitFailedEventData struct {
	// The CDP method or initialization step that failed (e.g. Target.setAutoAttach).
	Step string `json:"step"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Step        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserMonitorInitFailedEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserMonitorInitFailedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The CDP connection to Chrome could not be re-established after exhausting all
// reconnection attempts. No further telemetry events will arrive on this session.
type BrowserMonitorReconnectFailedEvent struct {
	Category constant.Monitor `json:"category" default:"monitor"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                                  `json:"ts" api:"required"`
	Type constant.MonitorReconnectFailed        `json:"type" default:"monitor_reconnect_failed"`
	Data BrowserMonitorReconnectFailedEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserMonitorReconnectFailedEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserMonitorReconnectFailedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserMonitorReconnectFailedEventData struct {
	// Reason for the reconnection failure. reconnect_exhausted: all retry attempts
	// were used up without successfully restoring the CDP connection.
	//
	// Any of "reconnect_exhausted".
	Reason string `json:"reason"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Reason      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserMonitorReconnectFailedEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserMonitorReconnectFailedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The CDP connection to Chrome was successfully re-established after a
// disconnection. Events emitted during the gap are lost. Computed state is reset,
// so navigation and network tracking restart fresh from this point.
type BrowserMonitorReconnectedEvent struct {
	Category constant.Monitor `json:"category" default:"monitor"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                              `json:"ts" api:"required"`
	Type constant.MonitorReconnected        `json:"type" default:"monitor_reconnected"`
	Data BrowserMonitorReconnectedEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserMonitorReconnectedEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserMonitorReconnectedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserMonitorReconnectedEventData struct {
	// Wall-clock time in milliseconds taken to reconnect after the disconnection.
	ReconnectDurationMs int64 `json:"reconnect_duration_ms"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ReconnectDurationMs respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserMonitorReconnectedEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserMonitorReconnectedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A periodic screenshot of the browser viewport.
type BrowserMonitorScreenshotEvent struct {
	Category constant.Screenshot `json:"category" default:"screenshot"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                             `json:"ts" api:"required"`
	Type constant.MonitorScreenshot        `json:"type" default:"monitor_screenshot"`
	Data BrowserMonitorScreenshotEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserMonitorScreenshotEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserMonitorScreenshotEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserMonitorScreenshotEventData struct {
	// Base64-encoded PNG screenshot of the browser viewport.
	Png string `json:"png"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Png         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserMonitorScreenshotEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserMonitorScreenshotEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser network idle event emitted after a 500ms quiet period with no
// in-flight HTTP requests.
type BrowserNetworkIdleEvent struct {
	Category constant.Network `json:"category" default:"network"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                `json:"ts" api:"required"`
	Type constant.NetworkIdle `json:"type" default:"network_idle"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserEventContext `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserNetworkIdleEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserNetworkIdleEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser network loading failed event. If the request was already in flight
// when CDP attached (no prior network_request was emitted for it), url, frame_id,
// loader_id, and resource_type are absent; BrowserEventContext is partially
// populated in that case.
type BrowserNetworkLoadingFailedEvent struct {
	Category constant.Network `json:"category" default:"network"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                         `json:"ts" api:"required"`
	Type constant.NetworkLoadingFailed `json:"type" default:"network_loading_failed"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserNetworkLoadingFailedEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserNetworkLoadingFailedEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserNetworkLoadingFailedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserNetworkLoadingFailedEventData struct {
	// True if the request was canceled by the browser or page script.
	Canceled bool `json:"canceled"`
	// Network error description (e.g. net::ERR_CONNECTION_REFUSED).
	ErrorText string `json:"error_text"`
	// CDP request identifier matching the originating network_request event.
	RequestID string `json:"request_id"`
	// CDP Network.ResourceType for the request, passed through as-is from Chrome.
	// Known values include Document, Fetch, XHR, Script, Stylesheet, Image, Media,
	// Font, TextTrack, EventSource, WebSocket, Manifest, Prefetch, Other, and more.
	ResourceType string `json:"resource_type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Canceled     respjson.Field
		ErrorText    respjson.Field
		RequestID    respjson.Field
		ResourceType respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserNetworkLoadingFailedEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserNetworkLoadingFailedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser network request sent event.
type BrowserNetworkRequestEvent struct {
	Category constant.Network `json:"category" default:"network"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                   `json:"ts" api:"required"`
	Type constant.NetworkRequest `json:"type" default:"network_request"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserNetworkRequestEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserNetworkRequestEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserNetworkRequestEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserNetworkRequestEventData struct {
	// URL of the document that initiated the request.
	DocumentURL string `json:"document_url"`
	// Request headers.
	Headers BrowserHTTPHeaders `json:"headers"`
	// CDP Initiator.type indicating what caused the request, passed through as-is from
	// Chrome. Known values include script, parser, preload, and other.
	InitiatorType string `json:"initiator_type"`
	// True if this request is the result of a redirect.
	IsRedirect bool `json:"is_redirect"`
	// HTTP method as sent on the wire (e.g. GET, POST).
	Method string `json:"method"`
	// Request body for POST/PUT requests, if available.
	PostData string `json:"post_data"`
	// Original URL before the redirect, present when is_redirect is true.
	RedirectURL string `json:"redirect_url"`
	// CDP request identifier, unique within the session.
	RequestID string `json:"request_id"`
	// CDP Network.ResourceType for the request, passed through as-is from Chrome.
	// Known values include Document, Fetch, XHR, Script, Stylesheet, Image, Media,
	// Font, TextTrack, EventSource, WebSocket, Manifest, Prefetch, Other, and more.
	ResourceType string `json:"resource_type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DocumentURL   respjson.Field
		Headers       respjson.Field
		InitiatorType respjson.Field
		IsRedirect    respjson.Field
		Method        respjson.Field
		PostData      respjson.Field
		RedirectURL   respjson.Field
		RequestID     respjson.Field
		ResourceType  respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserNetworkRequestEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserNetworkRequestEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser network response received event. Fired after the response body is
// fully received, not when headers arrive.
type BrowserNetworkResponseEvent struct {
	Category constant.Network `json:"category" default:"network"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                    `json:"ts" api:"required"`
	Type constant.NetworkResponse `json:"type" default:"network_response"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserNetworkResponseEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserNetworkResponseEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserNetworkResponseEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserNetworkResponseEventData struct {
	// Truncated response body, present only for text MIME types.
	Body string `json:"body"`
	// Response headers.
	Headers BrowserHTTPHeaders `json:"headers"`
	// HTTP method of the original request.
	Method string `json:"method"`
	// MIME type of the response (e.g. text/html, application/json).
	MimeType string `json:"mime_type"`
	// CDP request identifier matching the originating network_request event.
	RequestID string `json:"request_id"`
	// CDP Network.ResourceType for the request, passed through as-is from Chrome.
	// Known values include Document, Fetch, XHR, Script, Stylesheet, Image, Media,
	// Font, TextTrack, EventSource, WebSocket, Manifest, Prefetch, Other, and more.
	ResourceType string `json:"resource_type"`
	// HTTP response status code.
	Status int64 `json:"status"`
	// HTTP response status text (e.g. OK, Not Found).
	StatusText string `json:"status_text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Body         respjson.Field
		Headers      respjson.Field
		Method       respjson.Field
		MimeType     respjson.Field
		RequestID    respjson.Field
		ResourceType respjson.Field
		Status       respjson.Field
		StatusText   respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserNetworkResponseEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserNetworkResponseEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A page's renderer process crashed (an "Aw, Snap!" failure) while the browser
// process itself stayed alive. Reported on the crashed page's session, with the
// session and target ids on `source.metadata`. Captured only while the `page`
// category is enabled.
type BrowserPageCrashedEvent struct {
	Category constant.Page `json:"category" default:"page"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                       `json:"ts" api:"required"`
	Type constant.PageCrashed        `json:"type" default:"page_crashed"`
	Data BrowserPageCrashedEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageCrashedEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageCrashedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPageCrashedEventData struct {
	// CDP target identifier of the crashed page.
	TargetID string `json:"target_id" api:"required"`
	// CDP target type of the page that produced the event.
	//
	// Any of "page", "background_page", "service_worker", "shared_worker", "other".
	TargetType string `json:"target_type" api:"required"`
	// URL the page was on when its renderer process crashed.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		TargetID    respjson.Field
		TargetType  respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageCrashedEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageCrashedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser DOMContentLoaded event (CDP Page.domContentEventFired).
type BrowserPageDomContentLoadedEvent struct {
	Category constant.Page `json:"category" default:"page"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                         `json:"ts" api:"required"`
	Type constant.PageDomContentLoaded `json:"type" default:"page_dom_content_loaded"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserPageDomContentLoadedEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageDomContentLoadedEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageDomContentLoadedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserPageDomContentLoadedEventData struct {
	// Chrome monotonic clock value in seconds at which DOMContentLoaded fired,
	// relative to browser process start (not Unix epoch). Use ts for wall-clock time.
	CdpTimestamp float64 `json:"cdp_timestamp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CdpTimestamp respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserPageDomContentLoadedEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageDomContentLoadedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser layout settled event emitted 1 second after page load with no
// intervening layout shifts, indicating visual stability. Each layout shift resets
// the 1-second timer.
type BrowserPageLayoutSettledEvent struct {
	Category constant.Page `json:"category" default:"page"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                      `json:"ts" api:"required"`
	Type constant.PageLayoutSettled `json:"type" default:"page_layout_settled"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserEventContext `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageLayoutSettledEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageLayoutSettledEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser cumulative layout shift (CLS) event from the Performance Timeline API.
type BrowserPageLayoutShiftEvent struct {
	Category constant.Page `json:"category" default:"page"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                    `json:"ts" api:"required"`
	Type constant.PageLayoutShift `json:"type" default:"page_layout_shift"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserPageLayoutShiftEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageLayoutShiftEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageLayoutShiftEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserPageLayoutShiftEventData struct {
	// Duration of the layout shift entry in milliseconds (always 0 for layout shifts
	// per spec).
	Duration float64 `json:"duration"`
	// PerformanceLayoutShift attributes from the Performance Timeline entry.
	LayoutShiftDetails BrowserPageLayoutShiftEventDataLayoutShiftDetails `json:"layout_shift_details"`
	// CDP frame identifier of the frame where the layout shift occurred.
	SourceFrameID string `json:"source_frame_id"`
	// Performance Timeline timestamp of the layout shift in milliseconds.
	Time float64 `json:"time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Duration           respjson.Field
		LayoutShiftDetails respjson.Field
		SourceFrameID      respjson.Field
		Time               respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserPageLayoutShiftEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageLayoutShiftEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// PerformanceLayoutShift attributes from the Performance Timeline entry.
type BrowserPageLayoutShiftEventDataLayoutShiftDetails struct {
	// True if the layout shift was preceded by user input within 500ms, excluding it
	// from CLS.
	HadRecentInput bool `json:"had_recent_input"`
	// Layout shift score for this entry (contribution to CLS).
	Value float64 `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		HadRecentInput respjson.Field
		Value          respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageLayoutShiftEventDataLayoutShiftDetails) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageLayoutShiftEventDataLayoutShiftDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser Largest Contentful Paint (LCP) event from the Performance Timeline
// API.
type BrowserPageLcpEvent struct {
	Category constant.Page `json:"category" default:"page"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64            `json:"ts" api:"required"`
	Type constant.PageLcp `json:"type" default:"page_lcp"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserPageLcpEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageLcpEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageLcpEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserPageLcpEventData struct {
	// LargestContentfulPaint attributes from the Performance Timeline entry.
	LcpDetails BrowserPageLcpEventDataLcpDetails `json:"lcp_details"`
	// CDP frame identifier of the frame where the LCP element was rendered.
	SourceFrameID string `json:"source_frame_id"`
	// Performance Timeline timestamp of the LCP entry in milliseconds.
	Time float64 `json:"time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		LcpDetails    respjson.Field
		SourceFrameID respjson.Field
		Time          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserPageLcpEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageLcpEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// LargestContentfulPaint attributes from the Performance Timeline entry.
type BrowserPageLcpEventDataLcpDetails struct {
	// id attribute of the LCP element, if present.
	ElementID string `json:"element_id"`
	// Load time of the LCP element in milliseconds.
	LoadTime float64 `json:"load_time"`
	// CDP DOM node identifier of the LCP element.
	NodeID int64 `json:"node_id"`
	// Render time of the LCP element in milliseconds; 0 for cross-origin images
	// without Timing-Allow-Origin.
	RenderTime float64 `json:"render_time"`
	// Visible area of the LCP element in pixels squared.
	Size float64 `json:"size"`
	// URL of the LCP element for image or video elements.
	URL string `json:"url"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ElementID   respjson.Field
		LoadTime    respjson.Field
		NodeID      respjson.Field
		RenderTime  respjson.Field
		Size        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageLcpEventDataLcpDetails) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageLcpEventDataLcpDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser page load event (CDP Page.loadEventFired).
type BrowserPageLoadEvent struct {
	Category constant.Page `json:"category" default:"page"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64             `json:"ts" api:"required"`
	Type constant.PageLoad `json:"type" default:"page_load"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserPageLoadEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageLoadEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageLoadEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserPageLoadEventData struct {
	// Chrome monotonic clock value in seconds at which the load event fired, relative
	// to browser process start (not Unix epoch). Use ts for wall-clock time.
	CdpTimestamp float64 `json:"cdp_timestamp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CdpTimestamp respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserPageLoadEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageLoadEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser page navigation started event (CDP Page.frameNavigated). Carries nav
// context fields inline but not nav_seq, as this event resets the navigation
// epoch.
type BrowserPageNavigationEvent struct {
	Category constant.Page `json:"category" default:"page"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                          `json:"ts" api:"required"`
	Type constant.PageNavigation        `json:"type" default:"page_navigation"`
	Data BrowserPageNavigationEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageNavigationEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageNavigationEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPageNavigationEventData struct {
	// CDP frame identifier of the navigated frame.
	FrameID string `json:"frame_id"`
	// New CDP document loader identifier assigned for this navigation.
	LoaderID string `json:"loader_id"`
	// Parent frame identifier for subframe navigations; absent for top-level
	// navigations.
	ParentFrameID string `json:"parent_frame_id"`
	// CDP session identifier.
	SessionID string `json:"session_id"`
	// Browser target identifier.
	TargetID string `json:"target_id"`
	// CDP target type of the page that produced the event.
	//
	// Any of "page", "background_page", "service_worker", "shared_worker", "other".
	TargetType string `json:"target_type"`
	// URL navigated to.
	URL string `json:"url"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FrameID       respjson.Field
		LoaderID      respjson.Field
		ParentFrameID respjson.Field
		SessionID     respjson.Field
		TargetID      respjson.Field
		TargetType    respjson.Field
		URL           respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageNavigationEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageNavigationEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when page_dom_content_loaded and page_layout_settled have both fired for
// the same navigation, indicating the page is loaded and visually stable.
// Independent of network_idle; a single pending request does not block it.
type BrowserPageNavigationSettledEvent struct {
	Category constant.Page `json:"category" default:"page"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                          `json:"ts" api:"required"`
	Type constant.PageNavigationSettled `json:"type" default:"page_navigation_settled"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserEventContext `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageNavigationSettledEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageNavigationSettledEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A new browser tab or target was opened (CDP Target.attachedToTarget for page
// targets). Fires before a CDP session is attached to the new target, so
// session_id, frame_id, loader_id, and nav_seq are absent; this event does not
// compose BrowserEventContext. Consumers reading context fields generically should
// treat it as a special case.
type BrowserPageTabOpenedEvent struct {
	Category constant.Page `json:"category" default:"page"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                         `json:"ts" api:"required"`
	Type constant.PageTabOpened        `json:"type" default:"page_tab_opened"`
	Data BrowserPageTabOpenedEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageTabOpenedEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageTabOpenedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPageTabOpenedEventData struct {
	// Target identifier of the tab that opened this one, if any.
	OpenerID string `json:"opener_id"`
	// CDP target identifier for the newly opened tab.
	TargetID string `json:"target_id"`
	// CDP target type of the page that produced the event.
	//
	// Any of "page", "background_page", "service_worker", "shared_worker", "other".
	TargetType string `json:"target_type"`
	// Initial page title of the new tab.
	Title string `json:"title"`
	// Initial URL of the new tab.
	URL string `json:"url"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		OpenerID    respjson.Field
		TargetID    respjson.Field
		TargetType  respjson.Field
		Title       respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPageTabOpenedEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserPageTabOpenedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An HTTP call that manages the browser VM rather than driving the browser,
// handled by the in-VM API server — recording lifecycle, filesystem and process
// management, telemetry and browser configuration. Mostly platform-induced (e.g.
// profile save, replay capture) rather than agent actions.
type BrowserPlatformAPICallEvent struct {
	Category constant.Platform `json:"category" default:"platform"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                           `json:"ts" api:"required"`
	Type constant.PlatformAPICall        `json:"type" default:"platform_api_call"`
	Data BrowserPlatformAPICallEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPlatformAPICallEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserPlatformAPICallEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserPlatformAPICallEventData struct {
	// Wall-clock duration of the handler in milliseconds.
	DurationMs float64 `json:"duration_ms" api:"required"`
	// Matched route's operation, named as the in-VM API names its handler (e.g.
	// ProcessExec, StartRecording).
	OperationID string `json:"operation_id" api:"required"`
	// Per-request identifier from the in-VM API request middleware.
	RequestID string `json:"request_id" api:"required"`
	// HTTP response status code.
	Status int64 `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DurationMs  respjson.Field
		OperationID respjson.Field
		RequestID   respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserPlatformAPICallEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserPlatformAPICallEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A branded proxy-layer failure observed by the browser. Emitted when the metro
// egress host-proxy serves a branded 5xx error page whose response carries the
// X-Kernel-Proxy-Error header. Low-volume and carries a typed code. Its value is
// per-session and per-URL attribution for sessions that already capture the
// network stream: proxy failures are only observable while the CDP network
// collector is running, so this is an opt-in refinement of the raw network events
// rather than a default-on alerting signal.
type BrowserProxyErrorEvent struct {
	Category constant.Network `json:"category" default:"network"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64               `json:"ts" api:"required"`
	Type constant.ProxyError `json:"type" default:"proxy_error"`
	// Browser event context stamped by the browser monitor onto all CDP-sourced
	// events. Identifies the target, frame, and navigation epoch in which the event
	// occurred.
	Data BrowserProxyErrorEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserProxyErrorEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserProxyErrorEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Browser event context stamped by the browser monitor onto all CDP-sourced
// events. Identifies the target, frame, and navigation epoch in which the event
// occurred.
type BrowserProxyErrorEventData struct {
	// Proxy-layer error code: the X-Kernel-Proxy-Error response header value from a
	// branded 5xx error page served by the metro egress host-proxy. Values mirror what
	// the proxy emits: destination_blocked, provider_blacklisted,
	// provider_unreachable, proxy_unavailable, upstream_timeout, upstream_dns_failure,
	// upstream_connect_failed. Unknown header values are dropped.
	//
	// Any of "destination_blocked", "provider_blacklisted", "provider_unreachable",
	// "proxy_unavailable", "upstream_timeout", "upstream_dns_failure",
	// "upstream_connect_failed".
	Code string `json:"code" api:"required"`
	// CDP request identifier matching the originating request.
	RequestID string `json:"request_id" api:"required"`
	// HTTP response status of the branded error page (502).
	Status int64 `json:"status" api:"required"`
	// HTTP method of the failed request, when known.
	Method string `json:"method"`
	// CDP Network.ResourceType for the request, when known.
	ResourceType string `json:"resource_type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code         respjson.Field
		RequestID    respjson.Field
		Status       respjson.Field
		Method       respjson.Field
		ResourceType respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
	BrowserEventContext
}

// Returns the unmodified JSON received from the API
func (r BrowserProxyErrorEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserProxyErrorEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A managed service exited unexpectedly. Intentional stops do not produce this
// event; only unexpected exits and terminal restart-give-up transitions do.
type BrowserServiceCrashedEvent struct {
	Category constant.System `json:"category" default:"system"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                          `json:"ts" api:"required"`
	Type constant.ServiceCrashed        `json:"type" default:"service_crashed"`
	Data BrowserServiceCrashedEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserServiceCrashedEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserServiceCrashedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserServiceCrashedEventData struct {
	// Lifecycle phase the crash occurred in. startup: the process died before reaching
	// a healthy running state. running: a previously healthy process died
	// unexpectedly. gave_up: the process manager exhausted its restart attempts and
	// stopped trying.
	//
	// Any of "startup", "running", "gave_up".
	Phase string `json:"phase" api:"required"`
	// Program name of the crashed service (e.g. chromium, mutter, kernel-images-api).
	ServiceName string `json:"service_name" api:"required"`
	// PID of the crashed process. Absent when the process manager gave up after
	// exhausting restart attempts.
	Pid int64 `json:"pid"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Phase       respjson.Field
		ServiceName respjson.Field
		Pid         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserServiceCrashedEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserServiceCrashedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The Linux kernel OOM-killer terminated a process inside the VM. Fires for any
// process killed by the kernel due to memory exhaustion, including Chrome renderer
// subprocesses that are not supervised.
type BrowserSystemOomKillEvent struct {
	Category constant.System `json:"category" default:"system"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64                         `json:"ts" api:"required"`
	Type constant.SystemOomKill        `json:"type" default:"system_oom_kill"`
	Data BrowserSystemOomKillEventData `json:"data"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
		Data        respjson.Field
		Truncated   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserSystemOomKillEvent) RawJSON() string { return r.JSON.raw }
func (r *BrowserSystemOomKillEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserSystemOomKillEventData struct {
	// PID of the killed process.
	Pid int64 `json:"pid" api:"required"`
	// Comm of the killed process as reported by the kernel (max 15 chars, truncated by
	// the kernel).
	ProcessName string `json:"process_name" api:"required"`
	// Resident set size of the killed process in KiB (sum of anon-rss, file-rss, and
	// shmem-rss).
	RssKB int64 `json:"rss_kb" api:"required"`
	// Why the kernel decided to OOM-kill. none means global memory exhaustion; memcg
	// means a cgroup memory limit was hit; cpuset / memory_policy are
	// NUMA/policy-driven kills. Absent on kernels older than 5.0.
	//
	// Any of "none", "memcg", "cpuset", "memory_policy".
	Constraint string `json:"constraint"`
	// Free system memory in KiB at the time of the kill. Assumes a 4 KiB page size.
	// Does not include reclaimable caches. Absent if the kernel did not emit a
	// parseable Mem-Info section.
	MemFreeKB int64 `json:"mem_free_kb"`
	// Total system memory in KiB at the time of the kill. Assumes a 4 KiB page size.
	// Absent if the kernel did not emit a parseable Mem-Info section.
	MemTotalKB int64 `json:"mem_total_kb"`
	// Top processes by resident-set-size at the moment of the kill, sorted descending.
	// Empty if the kernel did not emit the Tasks state table. Capped at 5 entries.
	TopTasks []BrowserSystemOomKillEventDataTopTask `json:"top_tasks"`
	// PID of the triggering process. Absent if the kernel did not emit the standard
	// header line.
	TriggerPid int64 `json:"trigger_pid"`
	// Comm of the process whose allocation request caused the kernel to invoke the
	// OOM-killer. Often the same as process_name but can differ. Max 15 chars.
	TriggerProcessName string `json:"trigger_process_name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Pid                respjson.Field
		ProcessName        respjson.Field
		RssKB              respjson.Field
		Constraint         respjson.Field
		MemFreeKB          respjson.Field
		MemTotalKB         respjson.Field
		TopTasks           respjson.Field
		TriggerPid         respjson.Field
		TriggerProcessName respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserSystemOomKillEventData) RawJSON() string { return r.JSON.raw }
func (r *BrowserSystemOomKillEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserSystemOomKillEventDataTopTask struct {
	// Comm of the process (max 15 chars, truncated by the kernel).
	Name string `json:"name" api:"required"`
	// PID of the process.
	Pid int64 `json:"pid" api:"required"`
	// Resident set size in KiB at the moment of the kill.
	RssKB int64 `json:"rss_kb" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Pid         respjson.Field
		RssKB       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserSystemOomKillEventDataTopTask) RawJSON() string { return r.JSON.raw }
func (r *BrowserSystemOomKillEventDataTopTask) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-category telemetry capture settings layered onto the default set. The
// operational signals (control, connection, system, captcha) are on by default and
// are opt-out: set one to enabled=false to stop capturing it. The CDP categories
// (console, network, page, interaction), screenshot and platform are off by
// default and are opt-in: set enabled=true to capture them.
type BrowserTelemetryCategoriesConfig struct {
	// Captcha solver tasks and visible challenge outcomes. On by default.
	Captcha BrowserTelemetryCategoryConfig `json:"captcha"`
	// Client attach/detach lifecycle for the CDP proxy and live view. On by default.
	Connection BrowserTelemetryCategoryConfig `json:"connection"`
	// Console output (log, warn, error) and uncaught exceptions. CDP category; off by
	// default.
	Console BrowserTelemetryCategoryConfig `json:"console"`
	// Agent-driven actions against the browser — computer-control calls, Playwright
	// code execution, screenshots, clipboard access, and browser-control commands sent
	// over the CDP proxy. On by default.
	Control BrowserTelemetryControlConfig `json:"control"`
	// User interaction events including clicks, keydowns, and scroll-settled events.
	// CDP category; off by default.
	Interaction BrowserTelemetryCategoryConfig `json:"interaction"`
	// HTTP request and response metadata including URL, method, status code, and
	// timing. Request post data is forwarded as-is from CDP. Text response bodies are
	// truncated at 8 KB for structured types (JSON, XML, form data) and 4 KB for other
	// text types. Binary responses (images, fonts, media) are excluded. CDP category;
	// off by default.
	Network BrowserTelemetryCategoryConfig `json:"network"`
	// Page lifecycle events including navigation, DOMContentLoaded, load, layout
	// shifts, and LCP. CDP category; off by default.
	Page BrowserTelemetryCategoryConfig `json:"page"`
	// In-VM API calls that manage the browser VM rather than drive the browser
	// (recording, filesystem, process, telemetry and browser configuration). Mostly
	// platform-induced; off by default and must be opted into.
	Platform BrowserTelemetryCategoryConfig `json:"platform"`
	// Periodic base64-encoded viewport screenshots. High volume; off by default and
	// must be opted into.
	Screenshot BrowserTelemetryCategoryConfig `json:"screenshot"`
	// Browser VM health, such as out-of-memory kills and managed-service crashes. On
	// by default.
	System BrowserTelemetryCategoryConfig `json:"system"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Captcha     respjson.Field
		Connection  respjson.Field
		Console     respjson.Field
		Control     respjson.Field
		Interaction respjson.Field
		Network     respjson.Field
		Page        respjson.Field
		Platform    respjson.Field
		Screenshot  respjson.Field
		System      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserTelemetryCategoriesConfig) RawJSON() string { return r.JSON.raw }
func (r *BrowserTelemetryCategoriesConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BrowserTelemetryCategoriesConfig to a
// BrowserTelemetryCategoriesConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BrowserTelemetryCategoriesConfigParam.Overrides()
func (r BrowserTelemetryCategoriesConfig) ToParam() BrowserTelemetryCategoriesConfigParam {
	return param.Override[BrowserTelemetryCategoriesConfigParam](json.RawMessage(r.RawJSON()))
}

// Per-category telemetry capture settings layered onto the default set. The
// operational signals (control, connection, system, captcha) are on by default and
// are opt-out: set one to enabled=false to stop capturing it. The CDP categories
// (console, network, page, interaction), screenshot and platform are off by
// default and are opt-in: set enabled=true to capture them.
type BrowserTelemetryCategoriesConfigParam struct {
	// Captcha solver tasks and visible challenge outcomes. On by default.
	Captcha BrowserTelemetryCategoryConfigParam `json:"captcha,omitzero"`
	// Client attach/detach lifecycle for the CDP proxy and live view. On by default.
	Connection BrowserTelemetryCategoryConfigParam `json:"connection,omitzero"`
	// Console output (log, warn, error) and uncaught exceptions. CDP category; off by
	// default.
	Console BrowserTelemetryCategoryConfigParam `json:"console,omitzero"`
	// Agent-driven actions against the browser — computer-control calls, Playwright
	// code execution, screenshots, clipboard access, and browser-control commands sent
	// over the CDP proxy. On by default.
	Control BrowserTelemetryControlConfigParam `json:"control,omitzero"`
	// User interaction events including clicks, keydowns, and scroll-settled events.
	// CDP category; off by default.
	Interaction BrowserTelemetryCategoryConfigParam `json:"interaction,omitzero"`
	// HTTP request and response metadata including URL, method, status code, and
	// timing. Request post data is forwarded as-is from CDP. Text response bodies are
	// truncated at 8 KB for structured types (JSON, XML, form data) and 4 KB for other
	// text types. Binary responses (images, fonts, media) are excluded. CDP category;
	// off by default.
	Network BrowserTelemetryCategoryConfigParam `json:"network,omitzero"`
	// Page lifecycle events including navigation, DOMContentLoaded, load, layout
	// shifts, and LCP. CDP category; off by default.
	Page BrowserTelemetryCategoryConfigParam `json:"page,omitzero"`
	// In-VM API calls that manage the browser VM rather than drive the browser
	// (recording, filesystem, process, telemetry and browser configuration). Mostly
	// platform-induced; off by default and must be opted into.
	Platform BrowserTelemetryCategoryConfigParam `json:"platform,omitzero"`
	// Periodic base64-encoded viewport screenshots. High volume; off by default and
	// must be opted into.
	Screenshot BrowserTelemetryCategoryConfigParam `json:"screenshot,omitzero"`
	// Browser VM health, such as out-of-memory kills and managed-service crashes. On
	// by default.
	System BrowserTelemetryCategoryConfigParam `json:"system,omitzero"`
	paramObj
}

func (r BrowserTelemetryCategoriesConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserTelemetryCategoriesConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserTelemetryCategoriesConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-category telemetry configuration.
type BrowserTelemetryCategoryConfig struct {
	// Whether this category is captured. Operational categories (control, connection,
	// system, captcha) default to true; set false to opt out. CDP categories (console,
	// network, page, interaction), screenshot and platform default to false; set true
	// to opt in.
	Enabled bool `json:"enabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserTelemetryCategoryConfig) RawJSON() string { return r.JSON.raw }
func (r *BrowserTelemetryCategoryConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BrowserTelemetryCategoryConfig to a
// BrowserTelemetryCategoryConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BrowserTelemetryCategoryConfigParam.Overrides()
func (r BrowserTelemetryCategoryConfig) ToParam() BrowserTelemetryCategoryConfigParam {
	return param.Override[BrowserTelemetryCategoryConfigParam](json.RawMessage(r.RawJSON()))
}

// Per-category telemetry configuration.
type BrowserTelemetryCategoryConfigParam struct {
	// Whether this category is captured. Operational categories (control, connection,
	// system, captcha) default to true; set false to opt out. CDP categories (console,
	// network, page, interaction), screenshot and platform default to false; set true
	// to opt in.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	paramObj
}

func (r BrowserTelemetryCategoryConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserTelemetryCategoryConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserTelemetryCategoryConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Settings for the cdp_command events the CDP proxy reports.
type BrowserTelemetryCdpControlConfig struct {
	// Methods to leave out of the cdp_command stream. Omit the list to keep the
	// current one; send an empty list to report every supported method again.
	// Exclusion is a telemetry setting only: an excluded command is still relayed to
	// the browser unchanged, it simply produces no event. Use it to drop the
	// highest-volume methods — Input.dispatchMouseEvent during a humanized cursor
	// path, or Page.captureScreenshot under a screencast — without turning the whole
	// category off. Excluded commands are counted in
	// cdp_disconnect.telemetry_excluded.
	ExcludedMethods []BrowserCdpCommandMethod `json:"excluded_methods"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExcludedMethods respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserTelemetryCdpControlConfig) RawJSON() string { return r.JSON.raw }
func (r *BrowserTelemetryCdpControlConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BrowserTelemetryCdpControlConfig to a
// BrowserTelemetryCdpControlConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BrowserTelemetryCdpControlConfigParam.Overrides()
func (r BrowserTelemetryCdpControlConfig) ToParam() BrowserTelemetryCdpControlConfigParam {
	return param.Override[BrowserTelemetryCdpControlConfigParam](json.RawMessage(r.RawJSON()))
}

// Settings for the cdp_command events the CDP proxy reports.
type BrowserTelemetryCdpControlConfigParam struct {
	// Methods to leave out of the cdp_command stream. Omit the list to keep the
	// current one; send an empty list to report every supported method again.
	// Exclusion is a telemetry setting only: an excluded command is still relayed to
	// the browser unchanged, it simply produces no event. Use it to drop the
	// highest-volume methods — Input.dispatchMouseEvent during a humanized cursor
	// path, or Page.captureScreenshot under a screencast — without turning the whole
	// category off. Excluded commands are counted in
	// cdp_disconnect.telemetry_excluded.
	ExcludedMethods []BrowserCdpCommandMethod `json:"excluded_methods,omitzero"`
	paramObj
}

func (r BrowserTelemetryCdpControlConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserTelemetryCdpControlConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserTelemetryCdpControlConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Active telemetry configuration for a browser session.
type BrowserTelemetryConfig struct {
	// Per-category enable/disable flags.
	Browser BrowserTelemetryCategoriesConfig `json:"browser"`
	// Where the session's captured telemetry is being exported. Omitted when the
	// export state is unknown.
	Export BrowserTelemetryExportConfig `json:"export"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Browser     respjson.Field
		Export      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserTelemetryConfig) RawJSON() string { return r.JSON.raw }
func (r *BrowserTelemetryConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the control category. Same enabled semantics as any other
// category, plus settings for the browser-control commands the CDP proxy reports.
type BrowserTelemetryControlConfig struct {
	// Settings for the cdp_command events the CDP proxy reports. Merged independently
	// of enabled, so a later update that only sets enabled keeps the current exclusion
	// list.
	Cdp BrowserTelemetryCdpControlConfig `json:"cdp"`
	// Whether this category is captured. Control is on by default; set false to opt
	// out.
	Enabled bool `json:"enabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Cdp         respjson.Field
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserTelemetryControlConfig) RawJSON() string { return r.JSON.raw }
func (r *BrowserTelemetryControlConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BrowserTelemetryControlConfig to a
// BrowserTelemetryControlConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BrowserTelemetryControlConfigParam.Overrides()
func (r BrowserTelemetryControlConfig) ToParam() BrowserTelemetryControlConfigParam {
	return param.Override[BrowserTelemetryControlConfigParam](json.RawMessage(r.RawJSON()))
}

// Configuration for the control category. Same enabled semantics as any other
// category, plus settings for the browser-control commands the CDP proxy reports.
type BrowserTelemetryControlConfigParam struct {
	// Whether this category is captured. Control is on by default; set false to opt
	// out.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Settings for the cdp_command events the CDP proxy reports. Merged independently
	// of enabled, so a later update that only sets enabled keeps the current exclusion
	// list.
	Cdp BrowserTelemetryCdpControlConfigParam `json:"cdp,omitzero"`
	paramObj
}

func (r BrowserTelemetryControlConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserTelemetryControlConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserTelemetryControlConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BrowserTelemetryEventUnion contains all possible properties and values from
// [BrowserConsoleLogEvent], [BrowserConsoleErrorEvent],
// [BrowserNetworkRequestEvent], [BrowserNetworkResponseEvent],
// [BrowserNetworkLoadingFailedEvent], [BrowserNetworkIdleEvent],
// [BrowserProxyErrorEvent], [BrowserPageNavigationEvent],
// [BrowserPageDomContentLoadedEvent], [BrowserPageLoadEvent],
// [BrowserPageTabOpenedEvent], [BrowserPageCrashedEvent],
// [BrowserPageLayoutShiftEvent], [BrowserPageLcpEvent],
// [BrowserPageLayoutSettledEvent], [BrowserPageNavigationSettledEvent],
// [BrowserInteractionClickEvent], [BrowserInteractionKeyEvent],
// [BrowserInteractionScrollSettledEvent], [BrowserMonitorScreenshotEvent],
// [BrowserMonitorDisconnectedEvent], [BrowserMonitorReconnectedEvent],
// [BrowserMonitorReconnectFailedEvent], [BrowserMonitorInitFailedEvent],
// [BrowserAPICallEvent], [BrowserPlatformAPICallEvent], [BrowserCdpCommandEvent],
// [BrowserCdpConnectEvent], [BrowserCdpDisconnectEvent],
// [BrowserLiveViewConnectEvent], [BrowserLiveViewDisconnectEvent],
// [BrowserCaptchaSolveStartedEvent], [BrowserCaptchaSolveResultEvent],
// [BrowserCaptchaChallengeResultEvent], [BrowserSystemOomKillEvent],
// [BrowserServiceCrashedEvent].
//
// Use the [BrowserTelemetryEventUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BrowserTelemetryEventUnion struct {
	Category string `json:"category"`
	// This field is from variant [BrowserConsoleLogEvent].
	Source BrowserEventSource `json:"source"`
	Ts     int64              `json:"ts"`
	// Any of "console_log", "console_error", "network_request", "network_response",
	// "network_loading_failed", "network_idle", "proxy_error", "page_navigation",
	// "page_dom_content_loaded", "page_load", "page_tab_opened", "page_crashed",
	// "page_layout_shift", "page_lcp", "page_layout_settled",
	// "page_navigation_settled", "interaction_click", "interaction_key",
	// "interaction_scroll_settled", "monitor_screenshot", "monitor_disconnected",
	// "monitor_reconnected", "monitor_reconnect_failed", "monitor_init_failed",
	// "api_call", "platform_api_call", "cdp_command", "cdp_connect", "cdp_disconnect",
	// "live_view_connect", "live_view_disconnect", "captcha_solve_started",
	// "captcha_solve_result", "captcha_challenge_result", "system_oom_kill",
	// "service_crashed".
	Type string `json:"type"`
	// This field is a union of [BrowserConsoleLogEventData],
	// [BrowserConsoleErrorEventData], [BrowserNetworkRequestEventData],
	// [BrowserNetworkResponseEventData], [BrowserNetworkLoadingFailedEventData],
	// [BrowserEventContext], [BrowserProxyErrorEventData],
	// [BrowserPageNavigationEventData], [BrowserPageDomContentLoadedEventData],
	// [BrowserPageLoadEventData], [BrowserPageTabOpenedEventData],
	// [BrowserPageCrashedEventData], [BrowserPageLayoutShiftEventData],
	// [BrowserPageLcpEventData], [BrowserInteractionClickEventData],
	// [BrowserInteractionKeyEventData], [BrowserInteractionScrollSettledEventData],
	// [BrowserMonitorScreenshotEventData], [BrowserMonitorDisconnectedEventData],
	// [BrowserMonitorReconnectedEventData], [BrowserMonitorReconnectFailedEventData],
	// [BrowserMonitorInitFailedEventData], [BrowserAPICallEventData],
	// [BrowserPlatformAPICallEventData], [BrowserCdpCommandEventDataUnion],
	// [BrowserCdpConnectEventData], [BrowserCdpDisconnectEventData],
	// [BrowserLiveViewConnectEventData], [BrowserLiveViewDisconnectEventData],
	// [BrowserCaptchaSolveStartedEventData], [BrowserCaptchaSolveResultEventData],
	// [BrowserCaptchaChallengeResultEventData], [BrowserSystemOomKillEventData],
	// [BrowserServiceCrashedEventData]
	Data      BrowserTelemetryEventUnionData `json:"data"`
	Truncated bool                           `json:"truncated"`
	JSON      struct {
		Category  respjson.Field
		Source    respjson.Field
		Ts        respjson.Field
		Type      respjson.Field
		Data      respjson.Field
		Truncated respjson.Field
		raw       string
	} `json:"-"`
}

// anyBrowserTelemetryEvent is implemented by each variant of
// [BrowserTelemetryEventUnion] to add type safety for the return type of
// [BrowserTelemetryEventUnion.AsAny]
type anyBrowserTelemetryEvent interface {
	implBrowserTelemetryEventUnion()
}

func (BrowserConsoleLogEvent) implBrowserTelemetryEventUnion()               {}
func (BrowserConsoleErrorEvent) implBrowserTelemetryEventUnion()             {}
func (BrowserNetworkRequestEvent) implBrowserTelemetryEventUnion()           {}
func (BrowserNetworkResponseEvent) implBrowserTelemetryEventUnion()          {}
func (BrowserNetworkLoadingFailedEvent) implBrowserTelemetryEventUnion()     {}
func (BrowserNetworkIdleEvent) implBrowserTelemetryEventUnion()              {}
func (BrowserProxyErrorEvent) implBrowserTelemetryEventUnion()               {}
func (BrowserPageNavigationEvent) implBrowserTelemetryEventUnion()           {}
func (BrowserPageDomContentLoadedEvent) implBrowserTelemetryEventUnion()     {}
func (BrowserPageLoadEvent) implBrowserTelemetryEventUnion()                 {}
func (BrowserPageTabOpenedEvent) implBrowserTelemetryEventUnion()            {}
func (BrowserPageCrashedEvent) implBrowserTelemetryEventUnion()              {}
func (BrowserPageLayoutShiftEvent) implBrowserTelemetryEventUnion()          {}
func (BrowserPageLcpEvent) implBrowserTelemetryEventUnion()                  {}
func (BrowserPageLayoutSettledEvent) implBrowserTelemetryEventUnion()        {}
func (BrowserPageNavigationSettledEvent) implBrowserTelemetryEventUnion()    {}
func (BrowserInteractionClickEvent) implBrowserTelemetryEventUnion()         {}
func (BrowserInteractionKeyEvent) implBrowserTelemetryEventUnion()           {}
func (BrowserInteractionScrollSettledEvent) implBrowserTelemetryEventUnion() {}
func (BrowserMonitorScreenshotEvent) implBrowserTelemetryEventUnion()        {}
func (BrowserMonitorDisconnectedEvent) implBrowserTelemetryEventUnion()      {}
func (BrowserMonitorReconnectedEvent) implBrowserTelemetryEventUnion()       {}
func (BrowserMonitorReconnectFailedEvent) implBrowserTelemetryEventUnion()   {}
func (BrowserMonitorInitFailedEvent) implBrowserTelemetryEventUnion()        {}
func (BrowserAPICallEvent) implBrowserTelemetryEventUnion()                  {}
func (BrowserPlatformAPICallEvent) implBrowserTelemetryEventUnion()          {}
func (BrowserCdpCommandEvent) implBrowserTelemetryEventUnion()               {}
func (BrowserCdpConnectEvent) implBrowserTelemetryEventUnion()               {}
func (BrowserCdpDisconnectEvent) implBrowserTelemetryEventUnion()            {}
func (BrowserLiveViewConnectEvent) implBrowserTelemetryEventUnion()          {}
func (BrowserLiveViewDisconnectEvent) implBrowserTelemetryEventUnion()       {}
func (BrowserCaptchaSolveStartedEvent) implBrowserTelemetryEventUnion()      {}
func (BrowserCaptchaSolveResultEvent) implBrowserTelemetryEventUnion()       {}
func (BrowserCaptchaChallengeResultEvent) implBrowserTelemetryEventUnion()   {}
func (BrowserSystemOomKillEvent) implBrowserTelemetryEventUnion()            {}
func (BrowserServiceCrashedEvent) implBrowserTelemetryEventUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BrowserTelemetryEventUnion.AsAny().(type) {
//	case kernel.BrowserConsoleLogEvent:
//	case kernel.BrowserConsoleErrorEvent:
//	case kernel.BrowserNetworkRequestEvent:
//	case kernel.BrowserNetworkResponseEvent:
//	case kernel.BrowserNetworkLoadingFailedEvent:
//	case kernel.BrowserNetworkIdleEvent:
//	case kernel.BrowserProxyErrorEvent:
//	case kernel.BrowserPageNavigationEvent:
//	case kernel.BrowserPageDomContentLoadedEvent:
//	case kernel.BrowserPageLoadEvent:
//	case kernel.BrowserPageTabOpenedEvent:
//	case kernel.BrowserPageCrashedEvent:
//	case kernel.BrowserPageLayoutShiftEvent:
//	case kernel.BrowserPageLcpEvent:
//	case kernel.BrowserPageLayoutSettledEvent:
//	case kernel.BrowserPageNavigationSettledEvent:
//	case kernel.BrowserInteractionClickEvent:
//	case kernel.BrowserInteractionKeyEvent:
//	case kernel.BrowserInteractionScrollSettledEvent:
//	case kernel.BrowserMonitorScreenshotEvent:
//	case kernel.BrowserMonitorDisconnectedEvent:
//	case kernel.BrowserMonitorReconnectedEvent:
//	case kernel.BrowserMonitorReconnectFailedEvent:
//	case kernel.BrowserMonitorInitFailedEvent:
//	case kernel.BrowserAPICallEvent:
//	case kernel.BrowserPlatformAPICallEvent:
//	case kernel.BrowserCdpCommandEvent:
//	case kernel.BrowserCdpConnectEvent:
//	case kernel.BrowserCdpDisconnectEvent:
//	case kernel.BrowserLiveViewConnectEvent:
//	case kernel.BrowserLiveViewDisconnectEvent:
//	case kernel.BrowserCaptchaSolveStartedEvent:
//	case kernel.BrowserCaptchaSolveResultEvent:
//	case kernel.BrowserCaptchaChallengeResultEvent:
//	case kernel.BrowserSystemOomKillEvent:
//	case kernel.BrowserServiceCrashedEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BrowserTelemetryEventUnion) AsAny() anyBrowserTelemetryEvent {
	switch u.Type {
	case "console_log":
		return u.AsConsoleLog()
	case "console_error":
		return u.AsConsoleError()
	case "network_request":
		return u.AsNetworkRequest()
	case "network_response":
		return u.AsNetworkResponse()
	case "network_loading_failed":
		return u.AsNetworkLoadingFailed()
	case "network_idle":
		return u.AsNetworkIdle()
	case "proxy_error":
		return u.AsProxyError()
	case "page_navigation":
		return u.AsPageNavigation()
	case "page_dom_content_loaded":
		return u.AsPageDomContentLoaded()
	case "page_load":
		return u.AsPageLoad()
	case "page_tab_opened":
		return u.AsPageTabOpened()
	case "page_crashed":
		return u.AsPageCrashed()
	case "page_layout_shift":
		return u.AsPageLayoutShift()
	case "page_lcp":
		return u.AsPageLcp()
	case "page_layout_settled":
		return u.AsPageLayoutSettled()
	case "page_navigation_settled":
		return u.AsPageNavigationSettled()
	case "interaction_click":
		return u.AsInteractionClick()
	case "interaction_key":
		return u.AsInteractionKey()
	case "interaction_scroll_settled":
		return u.AsInteractionScrollSettled()
	case "monitor_screenshot":
		return u.AsMonitorScreenshot()
	case "monitor_disconnected":
		return u.AsMonitorDisconnected()
	case "monitor_reconnected":
		return u.AsMonitorReconnected()
	case "monitor_reconnect_failed":
		return u.AsMonitorReconnectFailed()
	case "monitor_init_failed":
		return u.AsMonitorInitFailed()
	case "api_call":
		return u.AsAPICall()
	case "platform_api_call":
		return u.AsPlatformAPICall()
	case "cdp_command":
		return u.AsCdpCommand()
	case "cdp_connect":
		return u.AsCdpConnect()
	case "cdp_disconnect":
		return u.AsCdpDisconnect()
	case "live_view_connect":
		return u.AsLiveViewConnect()
	case "live_view_disconnect":
		return u.AsLiveViewDisconnect()
	case "captcha_solve_started":
		return u.AsCaptchaSolveStarted()
	case "captcha_solve_result":
		return u.AsCaptchaSolveResult()
	case "captcha_challenge_result":
		return u.AsCaptchaChallengeResult()
	case "system_oom_kill":
		return u.AsSystemOomKill()
	case "service_crashed":
		return u.AsServiceCrashed()
	}
	return nil
}

func (u BrowserTelemetryEventUnion) AsConsoleLog() (v BrowserConsoleLogEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsConsoleError() (v BrowserConsoleErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsNetworkRequest() (v BrowserNetworkRequestEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsNetworkResponse() (v BrowserNetworkResponseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsNetworkLoadingFailed() (v BrowserNetworkLoadingFailedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsNetworkIdle() (v BrowserNetworkIdleEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsProxyError() (v BrowserProxyErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsPageNavigation() (v BrowserPageNavigationEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsPageDomContentLoaded() (v BrowserPageDomContentLoadedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsPageLoad() (v BrowserPageLoadEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsPageTabOpened() (v BrowserPageTabOpenedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsPageCrashed() (v BrowserPageCrashedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsPageLayoutShift() (v BrowserPageLayoutShiftEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsPageLcp() (v BrowserPageLcpEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsPageLayoutSettled() (v BrowserPageLayoutSettledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsPageNavigationSettled() (v BrowserPageNavigationSettledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsInteractionClick() (v BrowserInteractionClickEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsInteractionKey() (v BrowserInteractionKeyEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsInteractionScrollSettled() (v BrowserInteractionScrollSettledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsMonitorScreenshot() (v BrowserMonitorScreenshotEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsMonitorDisconnected() (v BrowserMonitorDisconnectedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsMonitorReconnected() (v BrowserMonitorReconnectedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsMonitorReconnectFailed() (v BrowserMonitorReconnectFailedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsMonitorInitFailed() (v BrowserMonitorInitFailedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsAPICall() (v BrowserAPICallEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsPlatformAPICall() (v BrowserPlatformAPICallEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsCdpCommand() (v BrowserCdpCommandEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsCdpConnect() (v BrowserCdpConnectEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsCdpDisconnect() (v BrowserCdpDisconnectEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsLiveViewConnect() (v BrowserLiveViewConnectEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsLiveViewDisconnect() (v BrowserLiveViewDisconnectEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsCaptchaSolveStarted() (v BrowserCaptchaSolveStartedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsCaptchaSolveResult() (v BrowserCaptchaSolveResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsCaptchaChallengeResult() (v BrowserCaptchaChallengeResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsSystemOomKill() (v BrowserSystemOomKillEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BrowserTelemetryEventUnion) AsServiceCrashed() (v BrowserServiceCrashedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BrowserTelemetryEventUnion) RawJSON() string { return u.JSON.raw }

func (r *BrowserTelemetryEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BrowserTelemetryEventUnionData is an implicit subunion of
// [BrowserTelemetryEventUnion]. BrowserTelemetryEventUnionData provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BrowserTelemetryEventUnion].
type BrowserTelemetryEventUnionData struct {
	FrameID  string `json:"frame_id"`
	LoaderID string `json:"loader_id"`
	// This field is from variant [BrowserConsoleLogEventData],
	// [BrowserConsoleErrorEventData], [BrowserNetworkRequestEventData],
	// [BrowserNetworkResponseEventData], [BrowserNetworkLoadingFailedEventData],
	// [BrowserEventContext], [BrowserProxyErrorEventData],
	// [BrowserPageDomContentLoadedEventData], [BrowserPageLoadEventData],
	// [BrowserPageLayoutShiftEventData], [BrowserPageLcpEventData],
	// [BrowserInteractionClickEventData], [BrowserInteractionKeyEventData],
	// [BrowserInteractionScrollSettledEventData].
	NavSeq     int64    `json:"nav_seq"`
	SessionID  string   `json:"session_id"`
	TargetID   string   `json:"target_id"`
	TargetType string   `json:"target_type"`
	URL        string   `json:"url"`
	Args       []string `json:"args"`
	Level      string   `json:"level"`
	// This field is from variant [BrowserConsoleLogEventData].
	StackTrace BrowserCallStack `json:"stack_trace"`
	Text       string           `json:"text"`
	// This field is from variant [BrowserConsoleErrorEventData].
	Column int64 `json:"column"`
	// This field is from variant [BrowserConsoleErrorEventData].
	Line int64 `json:"line"`
	// This field is from variant [BrowserConsoleErrorEventData].
	SourceURL string `json:"source_url"`
	// This field is from variant [BrowserNetworkRequestEventData].
	DocumentURL string `json:"document_url"`
	// This field is from variant [BrowserNetworkRequestEventData].
	Headers BrowserHTTPHeaders `json:"headers"`
	// This field is from variant [BrowserNetworkRequestEventData].
	InitiatorType string `json:"initiator_type"`
	// This field is from variant [BrowserNetworkRequestEventData].
	IsRedirect bool   `json:"is_redirect"`
	Method     string `json:"method"`
	// This field is from variant [BrowserNetworkRequestEventData].
	PostData string `json:"post_data"`
	// This field is from variant [BrowserNetworkRequestEventData].
	RedirectURL  string `json:"redirect_url"`
	RequestID    string `json:"request_id"`
	ResourceType string `json:"resource_type"`
	// This field is from variant [BrowserNetworkResponseEventData].
	Body string `json:"body"`
	// This field is from variant [BrowserNetworkResponseEventData].
	MimeType string `json:"mime_type"`
	// This field is a union of [int64], [int64], [int64], [int64], [string], [string]
	Status BrowserTelemetryEventUnionDataStatus `json:"status"`
	// This field is from variant [BrowserNetworkResponseEventData].
	StatusText string `json:"status_text"`
	// This field is from variant [BrowserNetworkLoadingFailedEventData].
	Canceled bool `json:"canceled"`
	// This field is from variant [BrowserNetworkLoadingFailedEventData].
	ErrorText string `json:"error_text"`
	Code      string `json:"code"`
	// This field is from variant [BrowserPageNavigationEventData].
	ParentFrameID string  `json:"parent_frame_id"`
	CdpTimestamp  float64 `json:"cdp_timestamp"`
	// This field is from variant [BrowserPageTabOpenedEventData].
	OpenerID string `json:"opener_id"`
	// This field is from variant [BrowserPageTabOpenedEventData].
	Title string `json:"title"`
	// This field is a union of [float64], [int64]
	Duration BrowserTelemetryEventUnionDataDuration `json:"duration"`
	// This field is from variant [BrowserPageLayoutShiftEventData].
	LayoutShiftDetails BrowserPageLayoutShiftEventDataLayoutShiftDetails `json:"layout_shift_details"`
	SourceFrameID      string                                            `json:"source_frame_id"`
	Time               float64                                           `json:"time"`
	// This field is from variant [BrowserPageLcpEventData].
	LcpDetails BrowserPageLcpEventDataLcpDetails `json:"lcp_details"`
	Selector   string                            `json:"selector"`
	Tag        string                            `json:"tag"`
	// This field is a union of [int64], [float64], [float64], [float64], [float64],
	// [float64], [float64], [float64]
	X BrowserTelemetryEventUnionDataX `json:"x"`
	// This field is a union of [int64], [float64], [float64], [float64], [float64],
	// [float64], [float64], [float64]
	Y BrowserTelemetryEventUnionDataY `json:"y"`
	// This field is from variant [BrowserInteractionKeyEventData].
	Key string `json:"key"`
	// This field is from variant [BrowserInteractionScrollSettledEventData].
	FromX int64 `json:"from_x"`
	// This field is from variant [BrowserInteractionScrollSettledEventData].
	FromY int64 `json:"from_y"`
	// This field is from variant [BrowserInteractionScrollSettledEventData].
	TargetSelector string `json:"target_selector"`
	// This field is from variant [BrowserInteractionScrollSettledEventData].
	ToX int64 `json:"to_x"`
	// This field is from variant [BrowserInteractionScrollSettledEventData].
	ToY int64 `json:"to_y"`
	// This field is from variant [BrowserMonitorScreenshotEventData].
	Png    string `json:"png"`
	Reason string `json:"reason"`
	// This field is from variant [BrowserMonitorReconnectedEventData].
	ReconnectDurationMs int64 `json:"reconnect_duration_ms"`
	// This field is from variant [BrowserMonitorInitFailedEventData].
	Step        string  `json:"step"`
	DurationMs  float64 `json:"duration_ms"`
	OperationID string  `json:"operation_id"`
	EventType   string  `json:"event_type"`
	Button      string  `json:"button"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	Buttons      int64   `json:"buttons"`
	ClickCount   int64   `json:"click_count"`
	CommandID    int64   `json:"command_id"`
	ConnectionID string  `json:"connection_id"`
	DeltaX       float64 `json:"delta_x"`
	DeltaY       float64 `json:"delta_y"`
	Force        float64 `json:"force"`
	Modifiers    int64   `json:"modifiers"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	PointerType        string  `json:"pointer_type"`
	TangentialPressure float64 `json:"tangential_pressure"`
	TiltX              float64 `json:"tilt_x"`
	TiltY              float64 `json:"tilt_y"`
	Twist              int64   `json:"twist"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	AutoRepeat bool `json:"auto_repeat"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	CommandCount int64 `json:"command_count"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	IsKeypad bool `json:"is_keypad"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	IsSystemKey bool `json:"is_system_key"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	Location int64 `json:"location"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	NamedKey   string `json:"named_key"`
	TextLength int64  `json:"text_length"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ReplacementEnd int64 `json:"replacement_end"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ReplacementStart int64 `json:"replacement_start"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	SelectionEnd int64 `json:"selection_end"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	SelectionStart int64 `json:"selection_start"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	TouchPointCount int64 `json:"touch_point_count"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	RadiusX float64 `json:"radius_x"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	RadiusY float64 `json:"radius_y"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	RotationAngle float64 `json:"rotation_angle"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	DragFileCount int64 `json:"drag_file_count"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	DragItemCount int64 `json:"drag_item_count"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	DragMimeCategories []string `json:"drag_mime_categories"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	DragOperationsMask int64  `json:"drag_operations_mask"`
	GestureSourceType  string `json:"gesture_source_type"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	RelativeSpeed int64 `json:"relative_speed"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ScaleFactor float64 `json:"scale_factor"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	PreventFling bool `json:"prevent_fling"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	RepeatCount int64 `json:"repeat_count"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	RepeatDelayMs int64 `json:"repeat_delay_ms"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	Speed int64 `json:"speed"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	XDistance float64 `json:"x_distance"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	XOverscroll float64 `json:"x_overscroll"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	YDistance float64 `json:"y_distance"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	YOverscroll float64 `json:"y_overscroll"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	TapCount int64 `json:"tap_count"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	FileCount     int64  `json:"file_count"`
	BackendNodeID int64  `json:"backend_node_id"`
	NodeID        int64  `json:"node_id"`
	ObjectID      string `json:"object_id"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	RectHeight float64 `json:"rect_height"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	RectWidth float64 `json:"rect_width"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	RectX float64 `json:"rect_x"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	RectY float64 `json:"rect_y"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	CaptureBeyondViewport bool `json:"capture_beyond_viewport"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ClipHeight float64 `json:"clip_height"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ClipScale float64 `json:"clip_scale"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ClipWidth float64 `json:"clip_width"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ClipX float64 `json:"clip_x"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ClipY  float64 `json:"clip_y"`
	Format string  `json:"format"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	FromSurface bool `json:"from_surface"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	OptimizeForSpeed bool  `json:"optimize_for_speed"`
	Quality          int64 `json:"quality"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	Accept bool `json:"accept"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	PromptTextLength int64 `json:"prompt_text_length"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ReferrerPolicy string `json:"referrer_policy"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ReferrerPresent bool `json:"referrer_present"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	TransitionType string `json:"transition_type"`
	URLScheme      string `json:"url_scheme"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	EntryID int64 `json:"entry_id"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	IgnoreCache bool `json:"ignore_cache"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ScriptLength int64 `json:"script_length"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	DisplayHeaderFooter bool `json:"display_header_footer"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	FooterTemplatePresent bool `json:"footer_template_present"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	GenerateDocumentOutline bool `json:"generate_document_outline"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	GenerateTaggedPdf bool `json:"generate_tagged_pdf"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	HeaderTemplatePresent bool `json:"header_template_present"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	Landscape bool `json:"landscape"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	MarginBottom float64 `json:"margin_bottom"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	MarginLeft float64 `json:"margin_left"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	MarginRight float64 `json:"margin_right"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	MarginTop float64 `json:"margin_top"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	PageRangesPresent bool `json:"page_ranges_present"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	PaperHeight float64 `json:"paper_height"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	PaperWidth float64 `json:"paper_width"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	PreferCssPageSize bool `json:"prefer_css_page_size"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	PrintBackground bool `json:"print_background"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	Scale float64 `json:"scale"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	TransferMode string `json:"transfer_mode"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	EveryNthFrame int64 `json:"every_nth_frame"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	MaxHeight int64 `json:"max_height"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	MaxWidth int64 `json:"max_width"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	State string `json:"state"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	Background       bool   `json:"background"`
	BrowserContextID string `json:"browser_context_id"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	EnableBeginFrameControl bool `json:"enable_begin_frame_control"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	Focus bool `json:"focus"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ForTab bool  `json:"for_tab"`
	Height int64 `json:"height"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	Hidden bool  `json:"hidden"`
	Left   int64 `json:"left"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	NewWindow   bool   `json:"new_window"`
	Top         int64  `json:"top"`
	Width       int64  `json:"width"`
	WindowState string `json:"window_state"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	DisposeOnDetach bool `json:"dispose_on_detach"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ProxyBypassListPresent bool `json:"proxy_bypass_list_present"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	ProxyServerPresent bool `json:"proxy_server_present"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	UniversalNetworkAccessOriginCount int64 `json:"universal_network_access_origin_count"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	PanelID string `json:"panel_id"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	DownloadGuid string `json:"download_guid"`
	WindowID     int64  `json:"window_id"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	FieldID int64 `json:"field_id"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	AddressFieldCount int64 `json:"address_field_count"`
	// This field is from variant [BrowserCdpCommandEventDataUnion].
	Mode string `json:"mode"`
	// This field is from variant [BrowserCdpDisconnectEventData].
	MessageCount int64 `json:"message_count"`
	// This field is from variant [BrowserCdpDisconnectEventData].
	TelemetryDropped int64 `json:"telemetry_dropped"`
	// This field is from variant [BrowserCdpDisconnectEventData].
	TelemetryExcluded int64  `json:"telemetry_excluded"`
	CaptchaType       string `json:"captcha_type"`
	ChallengeID       string `json:"challenge_id"`
	TaskID            string `json:"task_id"`
	WebsiteHost       string `json:"website_host"`
	WebsitePath       string `json:"website_path"`
	// This field is from variant [BrowserCaptchaSolveResultEventData].
	ErrorCode string `json:"error_code"`
	Pid       int64  `json:"pid"`
	// This field is from variant [BrowserSystemOomKillEventData].
	ProcessName string `json:"process_name"`
	// This field is from variant [BrowserSystemOomKillEventData].
	RssKB int64 `json:"rss_kb"`
	// This field is from variant [BrowserSystemOomKillEventData].
	Constraint string `json:"constraint"`
	// This field is from variant [BrowserSystemOomKillEventData].
	MemFreeKB int64 `json:"mem_free_kb"`
	// This field is from variant [BrowserSystemOomKillEventData].
	MemTotalKB int64 `json:"mem_total_kb"`
	// This field is from variant [BrowserSystemOomKillEventData].
	TopTasks []BrowserSystemOomKillEventDataTopTask `json:"top_tasks"`
	// This field is from variant [BrowserSystemOomKillEventData].
	TriggerPid int64 `json:"trigger_pid"`
	// This field is from variant [BrowserSystemOomKillEventData].
	TriggerProcessName string `json:"trigger_process_name"`
	// This field is from variant [BrowserServiceCrashedEventData].
	Phase string `json:"phase"`
	// This field is from variant [BrowserServiceCrashedEventData].
	ServiceName string `json:"service_name"`
	JSON        struct {
		FrameID                           respjson.Field
		LoaderID                          respjson.Field
		NavSeq                            respjson.Field
		SessionID                         respjson.Field
		TargetID                          respjson.Field
		TargetType                        respjson.Field
		URL                               respjson.Field
		Args                              respjson.Field
		Level                             respjson.Field
		StackTrace                        respjson.Field
		Text                              respjson.Field
		Column                            respjson.Field
		Line                              respjson.Field
		SourceURL                         respjson.Field
		DocumentURL                       respjson.Field
		Headers                           respjson.Field
		InitiatorType                     respjson.Field
		IsRedirect                        respjson.Field
		Method                            respjson.Field
		PostData                          respjson.Field
		RedirectURL                       respjson.Field
		RequestID                         respjson.Field
		ResourceType                      respjson.Field
		Body                              respjson.Field
		MimeType                          respjson.Field
		Status                            respjson.Field
		StatusText                        respjson.Field
		Canceled                          respjson.Field
		ErrorText                         respjson.Field
		Code                              respjson.Field
		ParentFrameID                     respjson.Field
		CdpTimestamp                      respjson.Field
		OpenerID                          respjson.Field
		Title                             respjson.Field
		Duration                          respjson.Field
		LayoutShiftDetails                respjson.Field
		SourceFrameID                     respjson.Field
		Time                              respjson.Field
		LcpDetails                        respjson.Field
		Selector                          respjson.Field
		Tag                               respjson.Field
		X                                 respjson.Field
		Y                                 respjson.Field
		Key                               respjson.Field
		FromX                             respjson.Field
		FromY                             respjson.Field
		TargetSelector                    respjson.Field
		ToX                               respjson.Field
		ToY                               respjson.Field
		Png                               respjson.Field
		Reason                            respjson.Field
		ReconnectDurationMs               respjson.Field
		Step                              respjson.Field
		DurationMs                        respjson.Field
		OperationID                       respjson.Field
		EventType                         respjson.Field
		Button                            respjson.Field
		Buttons                           respjson.Field
		ClickCount                        respjson.Field
		CommandID                         respjson.Field
		ConnectionID                      respjson.Field
		DeltaX                            respjson.Field
		DeltaY                            respjson.Field
		Force                             respjson.Field
		Modifiers                         respjson.Field
		PointerType                       respjson.Field
		TangentialPressure                respjson.Field
		TiltX                             respjson.Field
		TiltY                             respjson.Field
		Twist                             respjson.Field
		AutoRepeat                        respjson.Field
		CommandCount                      respjson.Field
		IsKeypad                          respjson.Field
		IsSystemKey                       respjson.Field
		Location                          respjson.Field
		NamedKey                          respjson.Field
		TextLength                        respjson.Field
		ReplacementEnd                    respjson.Field
		ReplacementStart                  respjson.Field
		SelectionEnd                      respjson.Field
		SelectionStart                    respjson.Field
		TouchPointCount                   respjson.Field
		RadiusX                           respjson.Field
		RadiusY                           respjson.Field
		RotationAngle                     respjson.Field
		DragFileCount                     respjson.Field
		DragItemCount                     respjson.Field
		DragMimeCategories                respjson.Field
		DragOperationsMask                respjson.Field
		GestureSourceType                 respjson.Field
		RelativeSpeed                     respjson.Field
		ScaleFactor                       respjson.Field
		PreventFling                      respjson.Field
		RepeatCount                       respjson.Field
		RepeatDelayMs                     respjson.Field
		Speed                             respjson.Field
		XDistance                         respjson.Field
		XOverscroll                       respjson.Field
		YDistance                         respjson.Field
		YOverscroll                       respjson.Field
		TapCount                          respjson.Field
		FileCount                         respjson.Field
		BackendNodeID                     respjson.Field
		NodeID                            respjson.Field
		ObjectID                          respjson.Field
		RectHeight                        respjson.Field
		RectWidth                         respjson.Field
		RectX                             respjson.Field
		RectY                             respjson.Field
		CaptureBeyondViewport             respjson.Field
		ClipHeight                        respjson.Field
		ClipScale                         respjson.Field
		ClipWidth                         respjson.Field
		ClipX                             respjson.Field
		ClipY                             respjson.Field
		Format                            respjson.Field
		FromSurface                       respjson.Field
		OptimizeForSpeed                  respjson.Field
		Quality                           respjson.Field
		Accept                            respjson.Field
		PromptTextLength                  respjson.Field
		ReferrerPolicy                    respjson.Field
		ReferrerPresent                   respjson.Field
		TransitionType                    respjson.Field
		URLScheme                         respjson.Field
		EntryID                           respjson.Field
		IgnoreCache                       respjson.Field
		ScriptLength                      respjson.Field
		DisplayHeaderFooter               respjson.Field
		FooterTemplatePresent             respjson.Field
		GenerateDocumentOutline           respjson.Field
		GenerateTaggedPdf                 respjson.Field
		HeaderTemplatePresent             respjson.Field
		Landscape                         respjson.Field
		MarginBottom                      respjson.Field
		MarginLeft                        respjson.Field
		MarginRight                       respjson.Field
		MarginTop                         respjson.Field
		PageRangesPresent                 respjson.Field
		PaperHeight                       respjson.Field
		PaperWidth                        respjson.Field
		PreferCssPageSize                 respjson.Field
		PrintBackground                   respjson.Field
		Scale                             respjson.Field
		TransferMode                      respjson.Field
		EveryNthFrame                     respjson.Field
		MaxHeight                         respjson.Field
		MaxWidth                          respjson.Field
		State                             respjson.Field
		Background                        respjson.Field
		BrowserContextID                  respjson.Field
		EnableBeginFrameControl           respjson.Field
		Focus                             respjson.Field
		ForTab                            respjson.Field
		Height                            respjson.Field
		Hidden                            respjson.Field
		Left                              respjson.Field
		NewWindow                         respjson.Field
		Top                               respjson.Field
		Width                             respjson.Field
		WindowState                       respjson.Field
		DisposeOnDetach                   respjson.Field
		ProxyBypassListPresent            respjson.Field
		ProxyServerPresent                respjson.Field
		UniversalNetworkAccessOriginCount respjson.Field
		PanelID                           respjson.Field
		DownloadGuid                      respjson.Field
		WindowID                          respjson.Field
		FieldID                           respjson.Field
		AddressFieldCount                 respjson.Field
		Mode                              respjson.Field
		MessageCount                      respjson.Field
		TelemetryDropped                  respjson.Field
		TelemetryExcluded                 respjson.Field
		CaptchaType                       respjson.Field
		ChallengeID                       respjson.Field
		TaskID                            respjson.Field
		WebsiteHost                       respjson.Field
		WebsitePath                       respjson.Field
		ErrorCode                         respjson.Field
		Pid                               respjson.Field
		ProcessName                       respjson.Field
		RssKB                             respjson.Field
		Constraint                        respjson.Field
		MemFreeKB                         respjson.Field
		MemTotalKB                        respjson.Field
		TopTasks                          respjson.Field
		TriggerPid                        respjson.Field
		TriggerProcessName                respjson.Field
		Phase                             respjson.Field
		ServiceName                       respjson.Field
		raw                               string
	} `json:"-"`
}

func (r *BrowserTelemetryEventUnionData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BrowserTelemetryEventUnionDataStatus is an implicit subunion of
// [BrowserTelemetryEventUnion]. BrowserTelemetryEventUnionDataStatus provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BrowserTelemetryEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfInt OfBrowserCaptchaChallengeResultEventDataStatus]
type BrowserTelemetryEventUnionDataStatus struct {
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	// This field will be present if the value is a [string] instead of an object.
	OfBrowserCaptchaChallengeResultEventDataStatus string `json:",inline"`
	JSON                                           struct {
		OfInt                                          respjson.Field
		OfBrowserCaptchaChallengeResultEventDataStatus respjson.Field
		raw                                            string
	} `json:"-"`
}

func (r *BrowserTelemetryEventUnionDataStatus) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BrowserTelemetryEventUnionDataDuration is an implicit subunion of
// [BrowserTelemetryEventUnion]. BrowserTelemetryEventUnionDataDuration provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BrowserTelemetryEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfFloat OfInt]
type BrowserTelemetryEventUnionDataDuration struct {
	// This field will be present if the value is a [float64] instead of an object.
	OfFloat float64 `json:",inline"`
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	JSON  struct {
		OfFloat respjson.Field
		OfInt   respjson.Field
		raw     string
	} `json:"-"`
}

func (r *BrowserTelemetryEventUnionDataDuration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BrowserTelemetryEventUnionDataX is an implicit subunion of
// [BrowserTelemetryEventUnion]. BrowserTelemetryEventUnionDataX provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BrowserTelemetryEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfInt OfFloat]
type BrowserTelemetryEventUnionDataX struct {
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	// This field will be present if the value is a [float64] instead of an object.
	OfFloat float64 `json:",inline"`
	JSON    struct {
		OfInt   respjson.Field
		OfFloat respjson.Field
		raw     string
	} `json:"-"`
}

func (r *BrowserTelemetryEventUnionDataX) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BrowserTelemetryEventUnionDataY is an implicit subunion of
// [BrowserTelemetryEventUnion]. BrowserTelemetryEventUnionDataY provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BrowserTelemetryEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfInt OfFloat]
type BrowserTelemetryEventUnionDataY struct {
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	// This field will be present if the value is a [float64] instead of an object.
	OfFloat float64 `json:",inline"`
	JSON    struct {
		OfInt   respjson.Field
		OfFloat respjson.Field
		raw     string
	} `json:"-"`
}

func (r *BrowserTelemetryEventUnionDataY) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Active export state for a session's captured telemetry, by protocol.
type BrowserTelemetryExportConfig struct {
	// Active OTLP export state.
	Otlp BrowserTelemetryOtlpExportConfig `json:"otlp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Otlp        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserTelemetryExportConfig) RawJSON() string { return r.JSON.raw }
func (r *BrowserTelemetryExportConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Active OTLP export state for a browser session.
type BrowserTelemetryOtlpExportConfig struct {
	// ID of the OTLP destination the session is bound to. Omitted when the session is
	// not exporting.
	Destination string `json:"destination"`
	// Whether the session is exporting captured telemetry over OTLP.
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
func (r BrowserTelemetryOtlpExportConfig) RawJSON() string { return r.JSON.raw }
func (r *BrowserTelemetryOtlpExportConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Envelope wrapping a browser telemetry event with its monotonic sequence number.
// Each SSE data: frame carries one envelope as JSON. The seq value is also emitted
// as the SSE id: field so clients can pass it as Last-Event-ID on reconnect.
type BrowserTelemetryEventsResponse struct {
	// Union type representing any browser telemetry event. Discriminated on `type`.
	// Each event's `category` determines when it is captured. The CDP collector-health
	// events (monitor_disconnected, monitor_reconnected, monitor_reconnect_failed,
	// monitor_init_failed) use the `monitor` category, which is not user-configurable:
	// it flows automatically whenever any CDP category (console, network, page,
	// interaction) is captured, and is silent otherwise. monitor_screenshot uses the
	// opt-in `screenshot` category. All other event types are controlled by their
	// per-category enable/disable flags.
	Event BrowserTelemetryEventUnion `json:"event" api:"required"`
	// Process-monotonic sequence number assigned by the browser VM. Pass as
	// Last-Event-ID on reconnect to resume without gaps. Gaps in received seq values
	// indicate dropped events.
	Seq int64 `json:"seq" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Event       respjson.Field
		Seq         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserTelemetryEventsResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserTelemetryEventsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Envelope wrapping a browser telemetry event with its monotonic sequence number.
// Each SSE data: frame carries one envelope as JSON. The seq value is also emitted
// as the SSE id: field so clients can pass it as Last-Event-ID on reconnect.
type BrowserTelemetryStreamResponse struct {
	// Union type representing any browser telemetry event. Discriminated on `type`.
	// Each event's `category` determines when it is captured. The CDP collector-health
	// events (monitor_disconnected, monitor_reconnected, monitor_reconnect_failed,
	// monitor_init_failed) use the `monitor` category, which is not user-configurable:
	// it flows automatically whenever any CDP category (console, network, page,
	// interaction) is captured, and is silent otherwise. monitor_screenshot uses the
	// opt-in `screenshot` category. All other event types are controlled by their
	// per-category enable/disable flags.
	Event BrowserTelemetryEventUnion `json:"event" api:"required"`
	// Process-monotonic sequence number assigned by the browser VM. Pass as
	// Last-Event-ID on reconnect to resume without gaps. Gaps in received seq values
	// indicate dropped events.
	Seq int64 `json:"seq" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Event       respjson.Field
		Seq         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserTelemetryStreamResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserTelemetryStreamResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserTelemetryEventsParams struct {
	// Maximum number of events per page. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor: pass the X-Next-Offset value from the previous
	// response to fetch the next page. When set, paging continues from this cursor and
	// since is ignored, while until still bounds the page. It is not an event's seq
	// field, so do not derive it from the response body.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Read direction. asc (default) reads oldest first, starting from since or the
	// offset cursor. desc reads newest first: each request returns one page of up to
	// limit records ending at the offset cursor (or until, or the newest archived
	// event); combining desc with since is rejected with a 400. In either direction
	// the category filter applies within the page, so a filtered page may be empty
	// while X-Has-More is true.
	Order param.Opt[string] `query:"order,omitzero" json:"-"`
	// Start of the window: an RFC-3339 timestamp, or a duration like 5m meaning that
	// long ago. Defaults to 5m. Ignored when offset is set.
	Since param.Opt[string] `query:"since,omitzero" json:"-"`
	// End of the window (exclusive): an RFC-3339 timestamp, or a duration like 5m
	// meaning that long ago.
	Until param.Opt[string] `query:"until,omitzero" json:"-"`
	// Restrict results to these event categories. Repeat the parameter for multiple
	// values.
	//
	// Any of "console", "network", "page", "interaction", "control", "platform",
	// "connection", "system", "screenshot", "captcha", "monitor".
	Category []string `query:"category,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserTelemetryEventsParams]'s query parameters as
// `url.Values`.
func (r BrowserTelemetryEventsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BrowserTelemetryStreamParams struct {
	// Pass `all` to start from the oldest retained event instead of only new events;
	// any other value is treated as from-now. The buffer is bounded, so the first
	// event id may be greater than 1 if older events were evicted.
	Replay      param.Opt[string] `query:"replay,omitzero" json:"-"`
	LastEventID param.Opt[string] `header:"Last-Event-ID,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserTelemetryStreamParams]'s query parameters as
// `url.Values`.
func (r BrowserTelemetryStreamParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
