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
	"github.com/kernel/kernel-go-sdk/packages/ssestream"
)

// Execute and manage processes on the browser instance.
//
// BrowserProcessService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserProcessService] method instead.
type BrowserProcessService struct {
	Options []option.RequestOption
}

// NewBrowserProcessService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBrowserProcessService(opts ...option.RequestOption) (r BrowserProcessService) {
	r = BrowserProcessService{}
	r.Options = opts
	return
}

// Execute a command synchronously
func (r *BrowserProcessService) Exec(ctx context.Context, idOrName string, body BrowserProcessExecParams, opts ...option.RequestOption) (res *BrowserProcessExecResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/process/exec", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Send signal to process
func (r *BrowserProcessService) Kill(ctx context.Context, processID string, params BrowserProcessKillParams, opts ...option.RequestOption) (res *BrowserProcessKillResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if params.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	if processID == "" {
		err = errors.New("missing required process_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/process/%s/kill", params.IDOrName, processID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Resize a PTY-backed process terminal
func (r *BrowserProcessService) Resize(ctx context.Context, processID string, params BrowserProcessResizeParams, opts ...option.RequestOption) (res *BrowserProcessResizeResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if params.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	if processID == "" {
		err = errors.New("missing required process_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/process/%s/resize", params.IDOrName, processID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Execute a command asynchronously
func (r *BrowserProcessService) Spawn(ctx context.Context, idOrName string, body BrowserProcessSpawnParams, opts ...option.RequestOption) (res *BrowserProcessSpawnResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/process/spawn", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get process status
func (r *BrowserProcessService) Status(ctx context.Context, processID string, query BrowserProcessStatusParams, opts ...option.RequestOption) (res *BrowserProcessStatusResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if query.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	if processID == "" {
		err = errors.New("missing required process_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/process/%s/status", query.IDOrName, processID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Write to process stdin
func (r *BrowserProcessService) Stdin(ctx context.Context, processID string, params BrowserProcessStdinParams, opts ...option.RequestOption) (res *BrowserProcessStdinResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if params.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	if processID == "" {
		err = errors.New("missing required process_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/process/%s/stdin", params.IDOrName, processID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Stream process stdout via SSE
func (r *BrowserProcessService) StdoutStreamStreaming(ctx context.Context, processID string, query BrowserProcessStdoutStreamParams, opts ...option.RequestOption) (stream *ssestream.Stream[BrowserProcessStdoutStreamResponse]) {
	var (
		raw *http.Response
		err error
	)
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "text/event-stream")}, opts...)
	if query.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return ssestream.NewStream[BrowserProcessStdoutStreamResponse](nil, err)
	}
	if processID == "" {
		err = errors.New("missing required process_id parameter")
		return ssestream.NewStream[BrowserProcessStdoutStreamResponse](nil, err)
	}
	path := fmt.Sprintf("browsers/%s/process/%s/stdout/stream", query.IDOrName, processID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &raw, opts...)
	return ssestream.NewStream[BrowserProcessStdoutStreamResponse](ssestream.NewDecoder(raw), err)
}

// Result of a synchronous command execution.
type BrowserProcessExecResponse struct {
	// Execution duration in milliseconds.
	DurationMs int64 `json:"duration_ms"`
	// Process exit code.
	ExitCode int64 `json:"exit_code"`
	// Base64-encoded stderr buffer.
	StderrB64 string `json:"stderr_b64"`
	// Base64-encoded stdout buffer.
	StdoutB64 string `json:"stdout_b64"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DurationMs  respjson.Field
		ExitCode    respjson.Field
		StderrB64   respjson.Field
		StdoutB64   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserProcessExecResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserProcessExecResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Generic OK response.
type BrowserProcessKillResponse struct {
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
func (r BrowserProcessKillResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserProcessKillResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Generic OK response.
type BrowserProcessResizeResponse struct {
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
func (r BrowserProcessResizeResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserProcessResizeResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Information about a spawned process.
type BrowserProcessSpawnResponse struct {
	// OS process ID.
	Pid int64 `json:"pid"`
	// Server-assigned identifier for the process.
	ProcessID string `json:"process_id" format:"uuid"`
	// Timestamp when the process started.
	StartedAt time.Time `json:"started_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Pid         respjson.Field
		ProcessID   respjson.Field
		StartedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserProcessSpawnResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserProcessSpawnResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current status of a process.
type BrowserProcessStatusResponse struct {
	// Estimated CPU usage percentage.
	CPUPct float64 `json:"cpu_pct"`
	// Exit code if the process has exited.
	ExitCode int64 `json:"exit_code" api:"nullable"`
	// Estimated resident memory usage in bytes.
	MemBytes int64 `json:"mem_bytes"`
	// Process state.
	//
	// Any of "running", "exited".
	State BrowserProcessStatusResponseState `json:"state"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CPUPct      respjson.Field
		ExitCode    respjson.Field
		MemBytes    respjson.Field
		State       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserProcessStatusResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserProcessStatusResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Process state.
type BrowserProcessStatusResponseState string

const (
	BrowserProcessStatusResponseStateRunning BrowserProcessStatusResponseState = "running"
	BrowserProcessStatusResponseStateExited  BrowserProcessStatusResponseState = "exited"
)

// Result of writing to stdin.
type BrowserProcessStdinResponse struct {
	// Number of bytes written.
	WrittenBytes int64 `json:"written_bytes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		WrittenBytes respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserProcessStdinResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserProcessStdinResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SSE payload representing process output or lifecycle events.
type BrowserProcessStdoutStreamResponse struct {
	// Base64-encoded data from the process stream.
	DataB64 string `json:"data_b64"`
	// Lifecycle event type.
	//
	// Any of "exit".
	Event BrowserProcessStdoutStreamResponseEvent `json:"event"`
	// Exit code when the event is "exit".
	ExitCode int64 `json:"exit_code"`
	// Source stream of the data chunk.
	//
	// Any of "stdout", "stderr".
	Stream BrowserProcessStdoutStreamResponseStream `json:"stream"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DataB64     respjson.Field
		Event       respjson.Field
		ExitCode    respjson.Field
		Stream      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserProcessStdoutStreamResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserProcessStdoutStreamResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Lifecycle event type.
type BrowserProcessStdoutStreamResponseEvent string

const (
	BrowserProcessStdoutStreamResponseEventExit BrowserProcessStdoutStreamResponseEvent = "exit"
)

// Source stream of the data chunk.
type BrowserProcessStdoutStreamResponseStream string

const (
	BrowserProcessStdoutStreamResponseStreamStdout BrowserProcessStdoutStreamResponseStream = "stdout"
	BrowserProcessStdoutStreamResponseStreamStderr BrowserProcessStdoutStreamResponseStream = "stderr"
)

type BrowserProcessExecParams struct {
	// Executable or shell command to run.
	Command string `json:"command" api:"required"`
	// Run the process as this user.
	AsUser param.Opt[string] `json:"as_user,omitzero"`
	// Working directory (absolute path) to run the command in.
	Cwd param.Opt[string] `json:"cwd,omitzero"`
	// Maximum execution time in seconds.
	TimeoutSec param.Opt[int64] `json:"timeout_sec,omitzero"`
	// Run the process with root privileges.
	AsRoot param.Opt[bool] `json:"as_root,omitzero"`
	// Command arguments.
	Args []string `json:"args,omitzero"`
	// Environment variables to set for the process.
	Env map[string]string `json:"env,omitzero"`
	paramObj
}

func (r BrowserProcessExecParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserProcessExecParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserProcessExecParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserProcessKillParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`
	// Signal to send.
	//
	// Any of "TERM", "KILL", "INT", "HUP".
	Signal BrowserProcessKillParamsSignal `json:"signal,omitzero" api:"required"`
	paramObj
}

func (r BrowserProcessKillParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserProcessKillParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserProcessKillParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Signal to send.
type BrowserProcessKillParamsSignal string

const (
	BrowserProcessKillParamsSignalTerm BrowserProcessKillParamsSignal = "TERM"
	BrowserProcessKillParamsSignalKill BrowserProcessKillParamsSignal = "KILL"
	BrowserProcessKillParamsSignalInt  BrowserProcessKillParamsSignal = "INT"
	BrowserProcessKillParamsSignalHup  BrowserProcessKillParamsSignal = "HUP"
)

type BrowserProcessResizeParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`
	// New terminal columns.
	Cols int64 `json:"cols" api:"required"`
	// New terminal rows.
	Rows int64 `json:"rows" api:"required"`
	paramObj
}

func (r BrowserProcessResizeParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserProcessResizeParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserProcessResizeParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserProcessSpawnParams struct {
	// Executable or shell command to run.
	Command string `json:"command" api:"required"`
	// Run the process as this user.
	AsUser param.Opt[string] `json:"as_user,omitzero"`
	// Working directory (absolute path) to run the command in.
	Cwd param.Opt[string] `json:"cwd,omitzero"`
	// Maximum execution time in seconds.
	TimeoutSec param.Opt[int64] `json:"timeout_sec,omitzero"`
	// Allocate a pseudo-terminal (PTY) for interactive shells.
	AllocateTty param.Opt[bool] `json:"allocate_tty,omitzero"`
	// Run the process with root privileges.
	AsRoot param.Opt[bool] `json:"as_root,omitzero"`
	// Initial terminal columns. Only used when allocate_tty is true.
	Cols param.Opt[int64] `json:"cols,omitzero"`
	// Initial terminal rows. Only used when allocate_tty is true.
	Rows param.Opt[int64] `json:"rows,omitzero"`
	// Command arguments.
	Args []string `json:"args,omitzero"`
	// Environment variables to set for the process.
	Env map[string]string `json:"env,omitzero"`
	paramObj
}

func (r BrowserProcessSpawnParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserProcessSpawnParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserProcessSpawnParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserProcessStatusParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`
	paramObj
}

type BrowserProcessStdinParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`
	// Base64-encoded data to write.
	DataB64 string `json:"data_b64" api:"required"`
	paramObj
}

func (r BrowserProcessStdinParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserProcessStdinParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserProcessStdinParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserProcessStdoutStreamParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`
	paramObj
}
