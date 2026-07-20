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
	auditLogDownloadAttempts      = 7
	auditLogDownloadMaxRetryDelay = 8 * time.Second
)

var auditLogDownloadRetryBaseDelay = time.Second

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
	onProgress     func(AuditLogDownloadProgress)
	requestOptions []option.RequestOption
}

// WithAuditLogDownloadProgress registers a callback that runs after each chunk
// is verified and written.
func WithAuditLogDownloadProgress(fn func(AuditLogDownloadProgress)) AuditLogDownloadOption {
	return func(config *auditLogDownloadConfig) {
		config.onProgress = fn
	}
}

// WithAuditLogDownloadRequestOptions applies request options to every chunk
// request. Download controls retries itself so it can also retry truncated or
// corrupt response bodies.
func WithAuditLogDownloadRequestOptions(opts ...option.RequestOption) AuditLogDownloadOption {
	return func(config *auditLogDownloadConfig) {
		config.requestOptions = append(config.requestOptions, opts...)
	}
}

// Download writes a complete audit log export to dst. It requests chunks until
// the export is complete, verifies every chunk checksum, and retries transient
// transfer failures. Download does not close dst.
func (r *AuditLogService) Download(ctx context.Context, params AuditLogDownloadParams, dst io.Writer, opts ...AuditLogDownloadOption) (AuditLogDownloadResult, error) {
	if dst == nil {
		return AuditLogDownloadResult{}, fmt.Errorf("audit log download destination is nil")
	}

	config := auditLogDownloadConfig{}
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
	for {
		if cursor != "" {
			query.Cursor = String(cursor)
		}
		body, header, err := r.downloadAuditLogChunk(ctx, query, config.requestOptions)
		if err != nil {
			return result, err
		}
		chunkRows, nextCursor, hasMore, err := parseAuditLogDownloadHeaders(header, cursor)
		if err != nil {
			return result, err
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

func (r *AuditLogService) downloadAuditLogChunk(ctx context.Context, query AuditLogExportChunkParams, opts []option.RequestOption) ([]byte, http.Header, error) {
	for attempt := 1; ; attempt++ {
		body, header, err := r.downloadAuditLogChunkOnce(ctx, query, opts)
		if err == nil || attempt == auditLogDownloadAttempts || !retryableAuditLogDownloadError(err) {
			return body, header, err
		}
		delay := min(auditLogDownloadRetryBaseDelay<<(attempt-1), auditLogDownloadMaxRetryDelay)
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-time.After(delay):
		}
	}
}

func (r *AuditLogService) downloadAuditLogChunkOnce(ctx context.Context, query AuditLogExportChunkParams, opts []option.RequestOption) ([]byte, http.Header, error) {
	requestOpts := make([]option.RequestOption, 0, len(opts)+1)
	requestOpts = append(requestOpts, opts...)
	requestOpts = append(requestOpts, option.WithMaxRetries(0))
	res, err := r.ExportChunk(ctx, query, requestOpts...)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read audit log chunk: %w", err)
	}
	want := res.Header.Get("X-Content-Sha256")
	if want == "" {
		return nil, nil, fmt.Errorf("response missing X-Content-Sha256 header")
	}
	sum := sha256.Sum256(body)
	if got := hex.EncodeToString(sum[:]); got != want {
		return nil, nil, fmt.Errorf("audit log chunk checksum mismatch (got %s, want %s)", got, want)
	}
	return body, res.Header, nil
}

func retryableAuditLogDownloadError(err error) bool {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests || apiErr.StatusCode >= 500
	}
	return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
}

func parseAuditLogDownloadHeaders(header http.Header, currentCursor string) (int64, string, bool, error) {
	hasMore, err := strconv.ParseBool(header.Get("X-Has-More"))
	if err != nil {
		return 0, "", false, fmt.Errorf("response missing or invalid X-Has-More header")
	}
	rows, err := strconv.ParseInt(header.Get("X-Row-Count"), 10, 64)
	if err != nil || rows < 0 {
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
