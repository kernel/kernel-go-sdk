package kernel

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/kernel/kernel-go-sdk/option"
)

func auditLogDownloadResponse(w http.ResponseWriter, body, nextCursor string, rows int, hasMore bool) {
	sum := sha256.Sum256([]byte(body))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("X-Content-Sha256", hex.EncodeToString(sum[:]))
	w.Header().Set("X-Has-More", strconv.FormatBool(hasMore))
	w.Header().Set("X-Row-Count", strconv.Itoa(rows))
	if nextCursor != "" {
		w.Header().Set("X-Next-Cursor", nextCursor)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(body))
}

func auditLogDownloadParams() AuditLogDownloadParams {
	return AuditLogDownloadParams{
		Start:  time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2026, time.June, 2, 0, 0, 0, 0, time.UTC),
		Format: AuditLogExportChunkParamsFormatJSONLGz,
	}
}

func TestAuditLogDownloadWritesVerifiedChunks(t *testing.T) {
	var cursors []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cursors = append(cursors, r.URL.Query().Get("cursor"))
		if len(cursors) == 1 {
			auditLogDownloadResponse(w, "first", "next", 2, true)
			return
		}
		auditLogDownloadResponse(w, "second", "", 1, false)
	}))
	defer server.Close()

	client := NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))
	var dst bytes.Buffer
	var progress []AuditLogDownloadProgress
	result, err := client.AuditLogs.Download(
		context.Background(),
		auditLogDownloadParams(),
		&dst,
		WithAuditLogDownloadProgress(func(update AuditLogDownloadProgress) {
			progress = append(progress, update)
		}),
	)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}
	if got, want := dst.String(), "firstsecond"; got != want {
		t.Fatalf("download body = %q, want %q", got, want)
	}
	if len(cursors) != 2 || cursors[0] != "" || cursors[1] != "next" {
		t.Fatalf("cursors = %#v, want [\"\" \"next\"]", cursors)
	}
	if result != (AuditLogDownloadResult{BytesWritten: 11, Chunks: 2, Rows: 3}) {
		t.Fatalf("result = %#v", result)
	}
	if len(progress) != 2 || progress[1].AuditLogDownloadResult != result || progress[1].ChunkRows != 1 {
		t.Fatalf("progress = %#v", progress)
	}
}

func TestAuditLogDownloadRetriesChecksumMismatch(t *testing.T) {
	oldDelay := auditLogDownloadRetryBaseDelay
	auditLogDownloadRetryBaseDelay = 0
	t.Cleanup(func() { auditLogDownloadRetryBaseDelay = oldDelay })

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("X-Content-Sha256", "bad")
			w.Header().Set("X-Has-More", "false")
			w.Header().Set("X-Row-Count", "1")
			_, _ = w.Write([]byte("chunk"))
			return
		}
		auditLogDownloadResponse(w, "chunk", "", 1, false)
	}))
	defer server.Close()

	client := NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))
	var dst bytes.Buffer
	_, err := client.AuditLogs.Download(context.Background(), auditLogDownloadParams(), &dst)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
	if dst.String() != "chunk" {
		t.Fatalf("download body = %q", dst.String())
	}
}

func TestAuditLogDownloadRejectsInvalidCursorBeforeWriting(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auditLogDownloadResponse(w, "chunk", "", 1, true)
	}))
	defer server.Close()

	client := NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))
	var dst bytes.Buffer
	_, err := client.AuditLogs.Download(context.Background(), auditLogDownloadParams(), &dst)
	if err == nil || err.Error() != "response has invalid X-Next-Cursor header" {
		t.Fatalf("Download() error = %v", err)
	}
	if dst.Len() != 0 {
		t.Fatalf("wrote %d bytes before cursor validation", dst.Len())
	}
}

func TestAuditLogDownloadOwnsHTTPRetries(t *testing.T) {
	oldDelay := auditLogDownloadRetryBaseDelay
	auditLogDownloadRetryBaseDelay = 0
	t.Cleanup(func() { auditLogDownloadRetryBaseDelay = oldDelay })

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"temporary failure"}`))
	}))
	defer server.Close()

	client := NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))
	_, err := client.AuditLogs.Download(
		context.Background(),
		auditLogDownloadParams(),
		&bytes.Buffer{},
		WithAuditLogDownloadRequestOptions(option.WithMaxRetries(5)),
	)
	if err == nil {
		t.Fatal("Download() error = nil")
	}
	if attempts != auditLogDownloadAttempts {
		t.Fatalf("attempts = %d, want %d", attempts, auditLogDownloadAttempts)
	}
}

func TestAuditLogDownloadDoesNotRetryClientErrors(t *testing.T) {
	oldDelay := auditLogDownloadRetryBaseDelay
	auditLogDownloadRetryBaseDelay = 0
	t.Cleanup(func() { auditLogDownloadRetryBaseDelay = oldDelay })

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"bad request"}`))
	}))
	defer server.Close()

	client := NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))
	_, err := client.AuditLogs.Download(context.Background(), auditLogDownloadParams(), &bytes.Buffer{})
	if err == nil {
		t.Fatal("Download() error = nil")
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}

type failingAuditLogWriter struct{}

func (failingAuditLogWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestAuditLogDownloadStopsOnDestinationError(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		auditLogDownloadResponse(w, "chunk", "next", 1, true)
	}))
	defer server.Close()

	client := NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))
	_, err := client.AuditLogs.Download(context.Background(), auditLogDownloadParams(), failingAuditLogWriter{})
	if err == nil || err.Error() != "write audit log download: write failed" {
		t.Fatalf("Download() error = %v", err)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}
