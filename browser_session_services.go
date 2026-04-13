package kernel

import (
	"context"
	"io"
	"net/http"

	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/ssestream"
	"github.com/kernel/kernel-go-sdk/shared"
)

// BrowserSessionProcessService exposes process APIs without passing browser id.
type BrowserSessionProcessService struct {
	inner BrowserProcessService
	id    string
}

func (s BrowserSessionProcessService) Exec(ctx context.Context, body BrowserProcessExecParams, opts ...option.RequestOption) (*BrowserProcessExecResponse, error) {
	return s.inner.Exec(ctx, s.id, body, opts...)
}

func (s BrowserSessionProcessService) Kill(ctx context.Context, processID string, params BrowserProcessKillParams, opts ...option.RequestOption) (*BrowserProcessKillResponse, error) {
	params.ID = s.id
	return s.inner.Kill(ctx, processID, params, opts...)
}

func (s BrowserSessionProcessService) Resize(ctx context.Context, processID string, params BrowserProcessResizeParams, opts ...option.RequestOption) (*BrowserProcessResizeResponse, error) {
	params.ID = s.id
	return s.inner.Resize(ctx, processID, params, opts...)
}

func (s BrowserSessionProcessService) Spawn(ctx context.Context, body BrowserProcessSpawnParams, opts ...option.RequestOption) (*BrowserProcessSpawnResponse, error) {
	return s.inner.Spawn(ctx, s.id, body, opts...)
}

func (s BrowserSessionProcessService) Status(ctx context.Context, processID string, query BrowserProcessStatusParams, opts ...option.RequestOption) (*BrowserProcessStatusResponse, error) {
	query.ID = s.id
	return s.inner.Status(ctx, processID, query, opts...)
}

func (s BrowserSessionProcessService) Stdin(ctx context.Context, processID string, params BrowserProcessStdinParams, opts ...option.RequestOption) (*BrowserProcessStdinResponse, error) {
	params.ID = s.id
	return s.inner.Stdin(ctx, processID, params, opts...)
}

func (s BrowserSessionProcessService) StdoutStreamStreaming(ctx context.Context, processID string, query BrowserProcessStdoutStreamParams, opts ...option.RequestOption) *ssestream.Stream[BrowserProcessStdoutStreamResponse] {
	query.ID = s.id
	return s.inner.StdoutStreamStreaming(ctx, processID, query, opts...)
}

// BrowserSessionComputerService exposes computer APIs without passing browser id.
type BrowserSessionComputerService struct {
	inner BrowserComputerService
	id    string
}

func (s BrowserSessionComputerService) Batch(ctx context.Context, body BrowserComputerBatchParams, opts ...option.RequestOption) error {
	return s.inner.Batch(ctx, s.id, body, opts...)
}

func (s BrowserSessionComputerService) CaptureScreenshot(ctx context.Context, body BrowserComputerCaptureScreenshotParams, opts ...option.RequestOption) (*http.Response, error) {
	return s.inner.CaptureScreenshot(ctx, s.id, body, opts...)
}

func (s BrowserSessionComputerService) ClickMouse(ctx context.Context, body BrowserComputerClickMouseParams, opts ...option.RequestOption) error {
	return s.inner.ClickMouse(ctx, s.id, body, opts...)
}

func (s BrowserSessionComputerService) DragMouse(ctx context.Context, body BrowserComputerDragMouseParams, opts ...option.RequestOption) error {
	return s.inner.DragMouse(ctx, s.id, body, opts...)
}

func (s BrowserSessionComputerService) GetMousePosition(ctx context.Context, opts ...option.RequestOption) (*BrowserComputerGetMousePositionResponse, error) {
	return s.inner.GetMousePosition(ctx, s.id, opts...)
}

func (s BrowserSessionComputerService) MoveMouse(ctx context.Context, body BrowserComputerMoveMouseParams, opts ...option.RequestOption) error {
	return s.inner.MoveMouse(ctx, s.id, body, opts...)
}

func (s BrowserSessionComputerService) PressKey(ctx context.Context, body BrowserComputerPressKeyParams, opts ...option.RequestOption) error {
	return s.inner.PressKey(ctx, s.id, body, opts...)
}

func (s BrowserSessionComputerService) ReadClipboard(ctx context.Context, opts ...option.RequestOption) (*BrowserComputerReadClipboardResponse, error) {
	return s.inner.ReadClipboard(ctx, s.id, opts...)
}

func (s BrowserSessionComputerService) Scroll(ctx context.Context, body BrowserComputerScrollParams, opts ...option.RequestOption) error {
	return s.inner.Scroll(ctx, s.id, body, opts...)
}

func (s BrowserSessionComputerService) SetCursorVisibility(ctx context.Context, body BrowserComputerSetCursorVisibilityParams, opts ...option.RequestOption) (*BrowserComputerSetCursorVisibilityResponse, error) {
	return s.inner.SetCursorVisibility(ctx, s.id, body, opts...)
}

func (s BrowserSessionComputerService) TypeText(ctx context.Context, body BrowserComputerTypeTextParams, opts ...option.RequestOption) error {
	return s.inner.TypeText(ctx, s.id, body, opts...)
}

func (s BrowserSessionComputerService) WriteClipboard(ctx context.Context, body BrowserComputerWriteClipboardParams, opts ...option.RequestOption) error {
	return s.inner.WriteClipboard(ctx, s.id, body, opts...)
}

// BrowserSessionFService exposes filesystem APIs without passing browser id.
type BrowserSessionFService struct {
	inner BrowserFService
	id    string
	Watch BrowserSessionFWatchService
}

func (s BrowserSessionFService) NewDirectory(ctx context.Context, body BrowserFNewDirectoryParams, opts ...option.RequestOption) error {
	return s.inner.NewDirectory(ctx, s.id, body, opts...)
}

func (s BrowserSessionFService) DeleteDirectory(ctx context.Context, body BrowserFDeleteDirectoryParams, opts ...option.RequestOption) error {
	return s.inner.DeleteDirectory(ctx, s.id, body, opts...)
}

func (s BrowserSessionFService) DeleteFile(ctx context.Context, body BrowserFDeleteFileParams, opts ...option.RequestOption) error {
	return s.inner.DeleteFile(ctx, s.id, body, opts...)
}

func (s BrowserSessionFService) DownloadDirZip(ctx context.Context, query BrowserFDownloadDirZipParams, opts ...option.RequestOption) (*http.Response, error) {
	return s.inner.DownloadDirZip(ctx, s.id, query, opts...)
}

func (s BrowserSessionFService) FileInfo(ctx context.Context, query BrowserFFileInfoParams, opts ...option.RequestOption) (*BrowserFFileInfoResponse, error) {
	return s.inner.FileInfo(ctx, s.id, query, opts...)
}

func (s BrowserSessionFService) ListFiles(ctx context.Context, query BrowserFListFilesParams, opts ...option.RequestOption) (*[]BrowserFListFilesResponse, error) {
	return s.inner.ListFiles(ctx, s.id, query, opts...)
}

func (s BrowserSessionFService) Move(ctx context.Context, body BrowserFMoveParams, opts ...option.RequestOption) error {
	return s.inner.Move(ctx, s.id, body, opts...)
}

func (s BrowserSessionFService) ReadFile(ctx context.Context, query BrowserFReadFileParams, opts ...option.RequestOption) (*http.Response, error) {
	return s.inner.ReadFile(ctx, s.id, query, opts...)
}

func (s BrowserSessionFService) SetFilePermissions(ctx context.Context, body BrowserFSetFilePermissionsParams, opts ...option.RequestOption) error {
	return s.inner.SetFilePermissions(ctx, s.id, body, opts...)
}

func (s BrowserSessionFService) Upload(ctx context.Context, body BrowserFUploadParams, opts ...option.RequestOption) error {
	return s.inner.Upload(ctx, s.id, body, opts...)
}

func (s BrowserSessionFService) UploadZip(ctx context.Context, body BrowserFUploadZipParams, opts ...option.RequestOption) error {
	return s.inner.UploadZip(ctx, s.id, body, opts...)
}

func (s BrowserSessionFService) WriteFile(ctx context.Context, contents io.Reader, params BrowserFWriteFileParams, opts ...option.RequestOption) error {
	return s.inner.WriteFile(ctx, s.id, contents, params, opts...)
}

// BrowserSessionFWatchService exposes fs watch APIs without passing browser id.
type BrowserSessionFWatchService struct {
	inner BrowserFWatchService
	id    string
}

func (s BrowserSessionFWatchService) EventsStreaming(ctx context.Context, watchID string, query BrowserFWatchEventsParams, opts ...option.RequestOption) *ssestream.Stream[BrowserFWatchEventsResponse] {
	query.ID = s.id
	return s.inner.EventsStreaming(ctx, watchID, query, opts...)
}

func (s BrowserSessionFWatchService) Start(ctx context.Context, body BrowserFWatchStartParams, opts ...option.RequestOption) (*BrowserFWatchStartResponse, error) {
	return s.inner.Start(ctx, s.id, body, opts...)
}

func (s BrowserSessionFWatchService) Stop(ctx context.Context, watchID string, body BrowserFWatchStopParams, opts ...option.RequestOption) error {
	body.ID = s.id
	return s.inner.Stop(ctx, watchID, body, opts...)
}

// BrowserSessionReplayService exposes replay APIs without passing browser id.
type BrowserSessionReplayService struct {
	inner BrowserReplayService
	id    string
}

func (s BrowserSessionReplayService) List(ctx context.Context, opts ...option.RequestOption) (*[]BrowserReplayListResponse, error) {
	return s.inner.List(ctx, s.id, opts...)
}

func (s BrowserSessionReplayService) Download(ctx context.Context, replayID string, query BrowserReplayDownloadParams, opts ...option.RequestOption) (*http.Response, error) {
	query.ID = s.id
	return s.inner.Download(ctx, replayID, query, opts...)
}

func (s BrowserSessionReplayService) Start(ctx context.Context, body BrowserReplayStartParams, opts ...option.RequestOption) (*BrowserReplayStartResponse, error) {
	return s.inner.Start(ctx, s.id, body, opts...)
}

func (s BrowserSessionReplayService) Stop(ctx context.Context, replayID string, body BrowserReplayStopParams, opts ...option.RequestOption) error {
	body.ID = s.id
	return s.inner.Stop(ctx, replayID, body, opts...)
}

// BrowserSessionLogService exposes log streaming without passing browser id.
type BrowserSessionLogService struct {
	inner BrowserLogService
	id    string
}

func (s BrowserSessionLogService) StreamStreaming(ctx context.Context, query BrowserLogStreamParams, opts ...option.RequestOption) *ssestream.Stream[shared.LogEvent] {
	return s.inner.StreamStreaming(ctx, s.id, query, opts...)
}

// BrowserSessionPlaywrightService exposes Playwright execution without passing browser id.
type BrowserSessionPlaywrightService struct {
	inner BrowserPlaywrightService
	id    string
}

func (s BrowserSessionPlaywrightService) Execute(ctx context.Context, body BrowserPlaywrightExecuteParams, opts ...option.RequestOption) (*BrowserPlaywrightExecuteResponse, error) {
	return s.inner.Execute(ctx, s.id, body, opts...)
}
