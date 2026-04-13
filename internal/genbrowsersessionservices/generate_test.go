package main

import (
	"bytes"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGenerateDeterministic(t *testing.T) {
	root := moduleRoot(t)
	a, err := Generate(root)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Generate(root)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("Generate is not deterministic for identical inputs")
	}
}

func TestGenerateIncludesExpectedServices(t *testing.T) {
	root := moduleRoot(t)
	src, err := Generate(root)
	if err != nil {
		t.Fatal(err)
	}
	text := string(src)
	for _, want := range []string{
		"type BrowserSessionClient struct",
		"func newBrowserSessionClient(sessionID, kernelBase, jwt string, scoped []option.RequestOption) *BrowserSessionClient",
		"\tReplays    BrowserSessionReplayService",
		"\tinnerFs := NewBrowserFService(scoped...)",
		"type BrowserSessionComputerService struct",
		"type BrowserSessionFService struct",
		"type BrowserSessionFWatchService struct",
		"type BrowserSessionLogService struct",
		"type BrowserSessionPlaywrightService struct",
		"type BrowserSessionProcessService struct",
		"type BrowserSessionReplayService struct",
		"func (s BrowserSessionProcessService) Kill(",
		"params.ID = s.id",
		"func (s BrowserSessionFService) WriteFile(ctx context.Context, contents io.Reader, params BrowserFWriteFileParams, opts ...option.RequestOption)",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("generated source missing %q", want)
		}
	}
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller")
	}
	// internal/genbrowsersessionservices/generate_test.go -> repo root
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
