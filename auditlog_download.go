package kernel

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
)

const (
	auditLogDownloadMaxChunkRows       = 50_000
	auditLogDownloadMaxTransferRetries = 6
	auditLogDownloadMaxRetryDelay      = 8 * time.Second
	auditLogDownloadRetryBaseDelay     = time.Second
)

// AuditLogDownloadParams configures a complete audit log export.
type AuditLogDownloadParams struct {
	End           time.Time
	Start         time.Time
	AuthStrategy  param.Opt[string]
	Limit         param.Opt[int64]
	Method        param.Opt[string]
	Search        param.Opt[string]
	Service       param.Opt[string]
	ExcludeMethod []string
	Format        AuditLogExportChunkParamsFormat
	SearchUserID  []string
}

// AuditLogDownloadResult summarizes a completed audit log download.
type AuditLogDownloadResult struct {
	BytesWritten int64
	Chunks       int
	Rows         int64
}

// AuditLogDownloadProgress reports progress after a verified chunk is written.
type AuditLogDownloadProgress struct {
	AuditLogDownloadResult
	ChunkRows int64
}

// AuditLogDownloadOption configures an audit log download.
type AuditLogDownloadOption func(*auditLogDownloadConfig)

type auditLogDownloadConfig struct {
	onProgress         func(AuditLogDownloadProgress)
	requestOptions     []option.RequestOption
	maxTransferRetries int
}

// WithAuditLogDownloadProgress registers a callback that runs after each chunk
// is verified and written.
func WithAuditLogDownloadProgress(fn func(AuditLogDownloadProgress)) AuditLogDownloadOption {
	return func(config *auditLogDownloadConfig) {
		config.onProgress = fn
	}
}

// WithAuditLogDownloadRequestOptions applies request options, including the
// standard HTTP retry policy, to every chunk request.
func WithAuditLogDownloadRequestOptions(opts ...option.RequestOption) AuditLogDownloadOption {
	return func(config *auditLogDownloadConfig) {
		config.requestOptions = append(config.requestOptions, opts...)
	}
}

// WithAuditLogDownloadMaxTransferRetries sets the number of retries for body
// read and checksum failures. HTTP retries remain controlled by request options.
func WithAuditLogDownloadMaxTransferRetries(retries int) AuditLogDownloadOption {
	if retries < 0 {
		panic("kernel: audit log download cannot have fewer than 0 transfer retries")
	}
	return func(config *auditLogDownloadConfig) {
		config.maxTransferRetries = retries
	}
}

// Download writes a complete audit log export to dst. It requests chunks until
// the export is complete, verifies every chunk checksum, and retries transient
// transfer failures. Download does not close dst. If Download returns an error,
// dst may contain a partial export; use a temporary file and atomic rename when
// the completed export must be published atomically.
func (r *AuditLogService) Download(ctx context.Context, params AuditLogDownloadParams, dst io.Writer, opts ...AuditLogDownloadOption) (AuditLogDownloadResult, error) {
	if dst == nil {
		return AuditLogDownloadResult{}, fmt.Errorf("audit log download destination is nil")
	}

	config := auditLogDownloadConfig{maxTransferRetries: auditLogDownloadMaxTransferRetries}
	for _, opt := range opts {
		opt(&config)
	}

	query := AuditLogExportChunkParams{
		End:           params.End,
		Start:         params.Start,
		AuthStrategy:  params.AuthStrategy,
		Limit:         params.Limit,
		Method:        params.Method,
		Search:        params.Search,
		Service:       params.Service,
		ExcludeMethod: params.ExcludeMethod,
		Format:        params.Format,
		SearchUserID:  params.SearchUserID,
	}
	cursor := ""
	result := AuditLogDownloadResult{}
	seenCursors := make(map[string]struct{})
	for {
		if cursor != "" {
			query.Cursor = String(cursor)
		}
		body, header, err := r.downloadAuditLogChunk(ctx, query, config.requestOptions, config.maxTransferRetries)
		if err != nil {
			return result, err
		}
		chunkRows, nextCursor, hasMore, err := parseAuditLogDownloadHeaders(header, cursor)
		if err != nil {
			return result, err
		}
		if hasMore {
			if _, ok := seenCursors[nextCursor]; ok {
				return result, fmt.Errorf("response repeated X-Next-Cursor header")
			}
			seenCursors[nextCursor] = struct{}{}
		}
		if err := writeAuditLogChunk(dst, body); err != nil {
			return result, fmt.Errorf("write audit log download: %w", err)
		}

		cursor = nextCursor
		result.BytesWritten += int64(len(body))
		result.Chunks++
		result.Rows += chunkRows
		if config.onProgress != nil {
			config.onProgress(AuditLogDownloadProgress{
				AuditLogDownloadResult: result,
				ChunkRows:              chunkRows,
			})
		}
		if !hasMore {
			return result, nil
		}
	}
}

func (r *AuditLogService) downloadAuditLogChunk(ctx context.Context, query AuditLogExportChunkParams, opts []option.RequestOption, maxTransferRetries int) ([]byte, http.Header, error) {
	for attempt := 1; ; attempt++ {
		body, header, err := r.downloadAuditLogChunkOnce(ctx, query, opts)
		var transferErr *auditLogTransferError
		if err == nil || attempt > maxTransferRetries || !errors.As(err, &transferErr) {
			return body, header, err
		}
		delay := auditLogDownloadRetryDelay(attempt)
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-time.After(delay):
		}
	}
}

func (r *AuditLogService) downloadAuditLogChunkOnce(ctx context.Context, query AuditLogExportChunkParams, opts []option.RequestOption) ([]byte, http.Header, error) {
	res, err := r.ExportChunk(ctx, query, opts...)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, &auditLogTransferError{err: fmt.Errorf("read audit log chunk: %w", err)}
	}
	want := res.Header.Get("X-Content-Sha256")
	if want == "" {
		return nil, nil, &auditLogTransferError{err: fmt.Errorf("response missing X-Content-Sha256 header")}
	}
	sum := sha256.Sum256(body)
	if got := hex.EncodeToString(sum[:]); got != want {
		return nil, nil, &auditLogTransferError{err: fmt.Errorf("audit log chunk checksum mismatch (got %s, want %s)", got, want)}
	}
	return body, res.Header, nil
}

type auditLogTransferError struct {
	err error
}

func (e *auditLogTransferError) Error() string {
	return e.err.Error()
}

func (e *auditLogTransferError) Unwrap() error {
	return e.err
}

func auditLogDownloadRetryDelay(attempt int) time.Duration {
	delay := auditLogDownloadRetryBaseDelay
	for retry := 1; retry < attempt && delay < auditLogDownloadMaxRetryDelay; retry++ {
		delay *= 2
	}
	return min(delay, auditLogDownloadMaxRetryDelay)
}

func parseAuditLogDownloadHeaders(header http.Header, currentCursor string) (int64, string, bool, error) {
	var hasMore bool
	switch header.Get("X-Has-More") {
	case "true":
		hasMore = true
	case "false":
		hasMore = false
	default:
		return 0, "", false, fmt.Errorf("response missing or invalid X-Has-More header")
	}

	rowCount := header.Get("X-Row-Count")
	if !isAuditLogDownloadDecimal(rowCount) {
		return 0, "", false, fmt.Errorf("response missing or invalid X-Row-Count header")
	}
	rows, err := strconv.ParseInt(rowCount, 10, 64)
	if err != nil || rows > auditLogDownloadMaxChunkRows {
		return 0, "", false, fmt.Errorf("response missing or invalid X-Row-Count header")
	}
	nextCursor := header.Get("X-Next-Cursor")
	if hasMore && (nextCursor == "" || nextCursor == currentCursor) {
		return 0, "", false, fmt.Errorf("response has invalid X-Next-Cursor header")
	}
	if !hasMore && nextCursor != "" {
		return 0, "", false, fmt.Errorf("response returned a cursor after the final chunk")
	}
	return rows, nextCursor, hasMore, nil
}

func isAuditLogDownloadDecimal(value string) bool {
	if value == "" {
		return false
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func writeAuditLogChunk(dst io.Writer, body []byte) error {
	for len(body) > 0 {
		n, err := dst.Write(body)
		if err != nil {
			return err
		}
		if n <= 0 || n > len(body) {
			return io.ErrShortWrite
		}
		body = body[n:]
	}
	return nil
}
