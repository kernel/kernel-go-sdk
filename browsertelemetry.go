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

// Stream live telemetry events from a browser session.
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

// An agent-driven HTTP call handled by the in-VM API server.
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
	// OpenAPI operationId of the matched route (e.g. processExec, takeScreenshot).
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
	// Captcha vendor family. Provider-specific task names are normalized into this
	// set; anything not covered is reported as other.
	//
	// Any of "hcaptcha", "recaptcha_v2", "recaptcha_v3", "turnstile", "geetest",
	// "other".
	CaptchaType string `json:"captcha_type" api:"required"`
	// Wall-clock duration from solve start to terminal outcome.
	DurationMs float64 `json:"duration_ms" api:"required"`
	// Terminal outcome. success: solver returned a usable solution. failure: solver
	// returned an error (see error_code). timeout: solver did not return within the
	// caller's wait budget. abandoned: caller cancelled or the page navigated away
	// mid-solve.
	//
	// Any of "success", "failure", "timeout", "abandoned".
	Status string `json:"status" api:"required"`
	// Solver-specific error code on failure (e.g. ERROR_CAPTCHA_UNSOLVABLE). Absent on
	// success.
	ErrorCode string `json:"error_code"`
	// Solver-assigned identifier. Opaque, useful for support cross-references.
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

// An external client (e.g. customer SDK, Playwright, Puppeteer) connected to the
// CDP WebSocket proxy on this VM.
type BrowserCdpConnectEvent struct {
	Category constant.Connection `json:"category" default:"connection"`
	// Provenance metadata identifying which producer emitted the event.
	Source BrowserEventSource `json:"source" api:"required"`
	// Event timestamp in Unix microseconds.
	Ts   int64               `json:"ts" api:"required"`
	Type constant.CdpConnect `json:"type" default:"cdp_connect"`
	// True if the data field was truncated due to size limits.
	Truncated bool `json:"truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Category    respjson.Field
		Source      respjson.Field
		Ts          respjson.Field
		Type        respjson.Field
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
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DurationMs   respjson.Field
		MessageCount respjson.Field
		Reason       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
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
	// Error message text. Present in both source paths.
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
// (console, network, page, interaction) and screenshot are off by default and are
// opt-in: set enabled=true to capture them.
type BrowserTelemetryCategoriesConfig struct {
	// Captcha solve attempt outcomes. On by default.
	Captcha BrowserTelemetryCategoryConfig `json:"captcha"`
	// Client attach/detach lifecycle for the CDP proxy and live view. On by default.
	Connection BrowserTelemetryCategoryConfig `json:"connection"`
	// Console output (log, warn, error) and uncaught exceptions. CDP category; off by
	// default.
	Console BrowserTelemetryCategoryConfig `json:"console"`
	// Agent-driven actions against the browser, such as inbound calls to the in-VM
	// API. On by default.
	Control BrowserTelemetryCategoryConfig `json:"control"`
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
// (console, network, page, interaction) and screenshot are off by default and are
// opt-in: set enabled=true to capture them.
type BrowserTelemetryCategoriesConfigParam struct {
	// Captcha solve attempt outcomes. On by default.
	Captcha BrowserTelemetryCategoryConfigParam `json:"captcha,omitzero"`
	// Client attach/detach lifecycle for the CDP proxy and live view. On by default.
	Connection BrowserTelemetryCategoryConfigParam `json:"connection,omitzero"`
	// Console output (log, warn, error) and uncaught exceptions. CDP category; off by
	// default.
	Console BrowserTelemetryCategoryConfigParam `json:"console,omitzero"`
	// Agent-driven actions against the browser, such as inbound calls to the in-VM
	// API. On by default.
	Control BrowserTelemetryCategoryConfigParam `json:"control,omitzero"`
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
	// network, page, interaction) and screenshot default to false; set true to opt in.
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
	// network, page, interaction) and screenshot default to false; set true to opt in.
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

// Active telemetry configuration for a browser session.
type BrowserTelemetryConfig struct {
	// Per-category enable/disable flags.
	Browser BrowserTelemetryCategoriesConfig `json:"browser"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Browser     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserTelemetryConfig) RawJSON() string { return r.JSON.raw }
func (r *BrowserTelemetryConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BrowserTelemetryEventUnion contains all possible properties and values from
// [BrowserConsoleLogEvent], [BrowserConsoleErrorEvent],
// [BrowserNetworkRequestEvent], [BrowserNetworkResponseEvent],
// [BrowserNetworkLoadingFailedEvent], [BrowserNetworkIdleEvent],
// [BrowserPageNavigationEvent], [BrowserPageDomContentLoadedEvent],
// [BrowserPageLoadEvent], [BrowserPageTabOpenedEvent],
// [BrowserPageLayoutShiftEvent], [BrowserPageLcpEvent],
// [BrowserPageLayoutSettledEvent], [BrowserPageNavigationSettledEvent],
// [BrowserInteractionClickEvent], [BrowserInteractionKeyEvent],
// [BrowserInteractionScrollSettledEvent], [BrowserMonitorScreenshotEvent],
// [BrowserMonitorDisconnectedEvent], [BrowserMonitorReconnectedEvent],
// [BrowserMonitorReconnectFailedEvent], [BrowserMonitorInitFailedEvent],
// [BrowserAPICallEvent], [BrowserCdpConnectEvent], [BrowserCdpDisconnectEvent],
// [BrowserLiveViewConnectEvent], [BrowserLiveViewDisconnectEvent],
// [BrowserCaptchaSolveResultEvent], [BrowserSystemOomKillEvent],
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
	// "network_loading_failed", "network_idle", "page_navigation",
	// "page_dom_content_loaded", "page_load", "page_tab_opened", "page_layout_shift",
	// "page_lcp", "page_layout_settled", "page_navigation_settled",
	// "interaction_click", "interaction_key", "interaction_scroll_settled",
	// "monitor_screenshot", "monitor_disconnected", "monitor_reconnected",
	// "monitor_reconnect_failed", "monitor_init_failed", "api_call", "cdp_connect",
	// "cdp_disconnect", "live_view_connect", "live_view_disconnect",
	// "captcha_solve_result", "system_oom_kill", "service_crashed".
	Type string `json:"type"`
	// This field is a union of [BrowserConsoleLogEventData],
	// [BrowserConsoleErrorEventData], [BrowserNetworkRequestEventData],
	// [BrowserNetworkResponseEventData], [BrowserNetworkLoadingFailedEventData],
	// [BrowserEventContext], [BrowserPageNavigationEventData],
	// [BrowserPageDomContentLoadedEventData], [BrowserPageLoadEventData],
	// [BrowserPageTabOpenedEventData], [BrowserPageLayoutShiftEventData],
	// [BrowserPageLcpEventData], [BrowserInteractionClickEventData],
	// [BrowserInteractionKeyEventData], [BrowserInteractionScrollSettledEventData],
	// [BrowserMonitorScreenshotEventData], [BrowserMonitorDisconnectedEventData],
	// [BrowserMonitorReconnectedEventData], [BrowserMonitorReconnectFailedEventData],
	// [BrowserMonitorInitFailedEventData], [BrowserAPICallEventData],
	// [BrowserCdpDisconnectEventData], [BrowserLiveViewConnectEventData],
	// [BrowserLiveViewDisconnectEventData], [BrowserCaptchaSolveResultEventData],
	// [BrowserSystemOomKillEventData], [BrowserServiceCrashedEventData]
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
func (BrowserPageNavigationEvent) implBrowserTelemetryEventUnion()           {}
func (BrowserPageDomContentLoadedEvent) implBrowserTelemetryEventUnion()     {}
func (BrowserPageLoadEvent) implBrowserTelemetryEventUnion()                 {}
func (BrowserPageTabOpenedEvent) implBrowserTelemetryEventUnion()            {}
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
func (BrowserCdpConnectEvent) implBrowserTelemetryEventUnion()               {}
func (BrowserCdpDisconnectEvent) implBrowserTelemetryEventUnion()            {}
func (BrowserLiveViewConnectEvent) implBrowserTelemetryEventUnion()          {}
func (BrowserLiveViewDisconnectEvent) implBrowserTelemetryEventUnion()       {}
func (BrowserCaptchaSolveResultEvent) implBrowserTelemetryEventUnion()       {}
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
//	case kernel.BrowserPageNavigationEvent:
//	case kernel.BrowserPageDomContentLoadedEvent:
//	case kernel.BrowserPageLoadEvent:
//	case kernel.BrowserPageTabOpenedEvent:
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
//	case kernel.BrowserCdpConnectEvent:
//	case kernel.BrowserCdpDisconnectEvent:
//	case kernel.BrowserLiveViewConnectEvent:
//	case kernel.BrowserLiveViewDisconnectEvent:
//	case kernel.BrowserCaptchaSolveResultEvent:
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
	case "page_navigation":
		return u.AsPageNavigation()
	case "page_dom_content_loaded":
		return u.AsPageDomContentLoaded()
	case "page_load":
		return u.AsPageLoad()
	case "page_tab_opened":
		return u.AsPageTabOpened()
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
	case "cdp_connect":
		return u.AsCdpConnect()
	case "cdp_disconnect":
		return u.AsCdpDisconnect()
	case "live_view_connect":
		return u.AsLiveViewConnect()
	case "live_view_disconnect":
		return u.AsLiveViewDisconnect()
	case "captcha_solve_result":
		return u.AsCaptchaSolveResult()
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

func (u BrowserTelemetryEventUnion) AsCaptchaSolveResult() (v BrowserCaptchaSolveResultEvent) {
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
	// [BrowserEventContext], [BrowserPageDomContentLoadedEventData],
	// [BrowserPageLoadEventData], [BrowserPageLayoutShiftEventData],
	// [BrowserPageLcpEventData], [BrowserInteractionClickEventData],
	// [BrowserInteractionKeyEventData], [BrowserInteractionScrollSettledEventData].
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
	// This field is a union of [int64], [int64], [string]
	Status BrowserTelemetryEventUnionDataStatus `json:"status"`
	// This field is from variant [BrowserNetworkResponseEventData].
	StatusText string `json:"status_text"`
	// This field is from variant [BrowserNetworkLoadingFailedEventData].
	Canceled bool `json:"canceled"`
	// This field is from variant [BrowserNetworkLoadingFailedEventData].
	ErrorText string `json:"error_text"`
	// This field is from variant [BrowserPageNavigationEventData].
	ParentFrameID string  `json:"parent_frame_id"`
	CdpTimestamp  float64 `json:"cdp_timestamp"`
	// This field is from variant [BrowserPageTabOpenedEventData].
	OpenerID string `json:"opener_id"`
	// This field is from variant [BrowserPageTabOpenedEventData].
	Title string `json:"title"`
	// This field is from variant [BrowserPageLayoutShiftEventData].
	Duration float64 `json:"duration"`
	// This field is from variant [BrowserPageLayoutShiftEventData].
	LayoutShiftDetails BrowserPageLayoutShiftEventDataLayoutShiftDetails `json:"layout_shift_details"`
	SourceFrameID      string                                            `json:"source_frame_id"`
	Time               float64                                           `json:"time"`
	// This field is from variant [BrowserPageLcpEventData].
	LcpDetails BrowserPageLcpEventDataLcpDetails `json:"lcp_details"`
	Selector   string                            `json:"selector"`
	Tag        string                            `json:"tag"`
	// This field is from variant [BrowserInteractionClickEventData].
	X int64 `json:"x"`
	// This field is from variant [BrowserInteractionClickEventData].
	Y int64 `json:"y"`
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
	Step       string  `json:"step"`
	DurationMs float64 `json:"duration_ms"`
	// This field is from variant [BrowserAPICallEventData].
	OperationID string `json:"operation_id"`
	// This field is from variant [BrowserCdpDisconnectEventData].
	MessageCount int64 `json:"message_count"`
	// This field is from variant [BrowserCaptchaSolveResultEventData].
	CaptchaType string `json:"captcha_type"`
	// This field is from variant [BrowserCaptchaSolveResultEventData].
	ErrorCode string `json:"error_code"`
	// This field is from variant [BrowserCaptchaSolveResultEventData].
	TaskID string `json:"task_id"`
	// This field is from variant [BrowserCaptchaSolveResultEventData].
	WebsiteHost string `json:"website_host"`
	// This field is from variant [BrowserCaptchaSolveResultEventData].
	WebsitePath string `json:"website_path"`
	Pid         int64  `json:"pid"`
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
		FrameID             respjson.Field
		LoaderID            respjson.Field
		NavSeq              respjson.Field
		SessionID           respjson.Field
		TargetID            respjson.Field
		TargetType          respjson.Field
		URL                 respjson.Field
		Args                respjson.Field
		Level               respjson.Field
		StackTrace          respjson.Field
		Text                respjson.Field
		Column              respjson.Field
		Line                respjson.Field
		SourceURL           respjson.Field
		DocumentURL         respjson.Field
		Headers             respjson.Field
		InitiatorType       respjson.Field
		IsRedirect          respjson.Field
		Method              respjson.Field
		PostData            respjson.Field
		RedirectURL         respjson.Field
		RequestID           respjson.Field
		ResourceType        respjson.Field
		Body                respjson.Field
		MimeType            respjson.Field
		Status              respjson.Field
		StatusText          respjson.Field
		Canceled            respjson.Field
		ErrorText           respjson.Field
		ParentFrameID       respjson.Field
		CdpTimestamp        respjson.Field
		OpenerID            respjson.Field
		Title               respjson.Field
		Duration            respjson.Field
		LayoutShiftDetails  respjson.Field
		SourceFrameID       respjson.Field
		Time                respjson.Field
		LcpDetails          respjson.Field
		Selector            respjson.Field
		Tag                 respjson.Field
		X                   respjson.Field
		Y                   respjson.Field
		Key                 respjson.Field
		FromX               respjson.Field
		FromY               respjson.Field
		TargetSelector      respjson.Field
		ToX                 respjson.Field
		ToY                 respjson.Field
		Png                 respjson.Field
		Reason              respjson.Field
		ReconnectDurationMs respjson.Field
		Step                respjson.Field
		DurationMs          respjson.Field
		OperationID         respjson.Field
		MessageCount        respjson.Field
		CaptchaType         respjson.Field
		ErrorCode           respjson.Field
		TaskID              respjson.Field
		WebsiteHost         respjson.Field
		WebsitePath         respjson.Field
		Pid                 respjson.Field
		ProcessName         respjson.Field
		RssKB               respjson.Field
		Constraint          respjson.Field
		MemFreeKB           respjson.Field
		MemTotalKB          respjson.Field
		TopTasks            respjson.Field
		TriggerPid          respjson.Field
		TriggerProcessName  respjson.Field
		Phase               respjson.Field
		ServiceName         respjson.Field
		raw                 string
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
// will be valid: OfInt OfBrowserCaptchaSolveResultEventDataStatus]
type BrowserTelemetryEventUnionDataStatus struct {
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	// This field will be present if the value is a [string] instead of an object.
	OfBrowserCaptchaSolveResultEventDataStatus string `json:",inline"`
	JSON                                       struct {
		OfInt                                      respjson.Field
		OfBrowserCaptchaSolveResultEventDataStatus respjson.Field
		raw                                        string
	} `json:"-"`
}

func (r *BrowserTelemetryEventUnionDataStatus) UnmarshalJSON(data []byte) error {
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
	// Any of "console", "network", "page", "interaction", "control", "connection",
	// "system", "screenshot", "captcha", "monitor".
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
