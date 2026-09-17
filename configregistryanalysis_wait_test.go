package kernel_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/option"
)

func configRegistryAnalysisResponse(id, status string, finished bool) string {
	finishedAt := "null"
	if finished {
		finishedAt = `"2026-09-16T00:01:00Z"`
	}
	return fmt.Sprintf(`{
		"analysis": {
			"id": %q,
			"created_at": "2026-09-16T00:00:00Z",
			"expires_at": "2026-09-16T00:45:00Z",
			"failure": null,
			"finished_at": %s,
			"status": %q,
			"intent": null
		},
		"recommendation": null,
		"target": {"domain":"example.com","host":"example.com","normalized":"https://example.com/"},
		"working_configurations": [],
		"guidance": null,
		"workload_outcome": null
	}`, id, finishedAt, status)
}

func TestConfigRegistryAnalysisWaitForResultPollsUnknownUnfinishedStatus(t *testing.T) {
	statuses := []string{"running", "queued", "archived"}
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index := int(requests.Add(1)) - 1
		if r.Header.Get("X-Test") != "preserved" {
			t.Errorf("X-Test = %q, want preserved", r.Header.Get("X-Test"))
		}
		if r.Header.Get("X-Stainless-Poll-Helper") != "true" {
			t.Errorf("X-Stainless-Poll-Helper = %q, want true", r.Header.Get("X-Stainless-Poll-Helper"))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, configRegistryAnalysisResponse("analysis-1", statuses[index], index == 2))
	}))
	defer server.Close()

	client := kernel.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))
	result, err := client.ConfigRegistry.Analyses.WaitForResult(
		context.Background(),
		"analysis-1",
		kernel.WithConfigRegistryAnalysisPollInterval(time.Millisecond),
		kernel.WithConfigRegistryAnalysisWaitRequestOptions(
			option.WithHeader("X-Test", "preserved"),
			option.WithHeader("X-Stainless-Poll-Helper", "caller"),
		),
	)
	if err != nil {
		t.Fatalf("WaitForResult() error = %v", err)
	}
	if result.Analysis.Status != kernel.AnalysisStatus("archived") {
		t.Fatalf("status = %q, want archived", result.Analysis.Status)
	}
	if got := requests.Load(); got != 3 {
		t.Fatalf("requests = %d, want 3", got)
	}
}

func TestConfigRegistryAnalysisWaitForResultReturnsKnownTerminalStatusWithoutFinishedAt(t *testing.T) {
	for _, status := range []string{"completed", "failed", "canceled", "expired"} {
		t.Run(status, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, configRegistryAnalysisResponse("analysis-1", status, false))
			}))
			defer server.Close()

			client := kernel.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))
			result, err := client.ConfigRegistry.Analyses.WaitForResult(context.Background(), "analysis-1")
			if err != nil {
				t.Fatalf("WaitForResult() error = %v", err)
			}
			if string(result.Analysis.Status) != status {
				t.Fatalf("status = %q, want %q", result.Analysis.Status, status)
			}
			if got := requests.Load(); got != 1 {
				t.Fatalf("requests = %d, want 1", got)
			}
		})
	}
}

func TestConfigRegistryAnalysisWaitForResultZeroMaxWaitReadsOnce(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, configRegistryAnalysisResponse("analysis-1", "running", false))
	}))
	defer server.Close()

	client := kernel.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))
	_, err := client.ConfigRegistry.Analyses.WaitForResult(
		context.Background(),
		"analysis-1",
		kernel.WithConfigRegistryAnalysisMaxWait(0),
	)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("WaitForResult() error = %v, want context deadline exceeded", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("requests = %d, want 1", got)
	}
}

func TestConfigRegistryAnalysisWaitForResultRejectsInvalidResponse(t *testing.T) {
	for _, test := range []struct {
		name string
		body string
		want string
	}{
		{name: "null analysis", body: `{"analysis":null}`, want: "missing an analysis"},
		{name: "mismatched id", body: configRegistryAnalysisResponse("analysis-2", "running", false), want: "analysis-2"},
		{
			name: "missing finished_at",
			body: strings.Replace(configRegistryAnalysisResponse("analysis-1", "running", false), `"finished_at": null,`, "", 1),
			want: "finished_at",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, test.body)
			}))
			defer server.Close()

			client := kernel.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))
			_, err := client.ConfigRegistry.Analyses.WaitForResult(context.Background(), "analysis-1")
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("WaitForResult() error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestConfigRegistryAnalysisWaitForResultStopsDuringPollSleep(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, configRegistryAnalysisResponse("analysis-1", "running", false))
	}))
	defer server.Close()

	client := kernel.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(10*time.Millisecond, cancel)
	_, err := client.ConfigRegistry.Analyses.WaitForResult(
		ctx,
		"analysis-1",
		kernel.WithConfigRegistryAnalysisPollInterval(time.Minute),
	)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("WaitForResult() error = %v, want context canceled", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("requests = %d, want 1", got)
	}
}

func TestConfigRegistryAnalysisWaitOptionsRejectInvalidDurations(t *testing.T) {
	for name, fn := range map[string]func(){
		"zero poll interval":     func() { kernel.WithConfigRegistryAnalysisPollInterval(0) },
		"negative poll interval": func() { kernel.WithConfigRegistryAnalysisPollInterval(-time.Second) },
		"negative max wait":      func() { kernel.WithConfigRegistryAnalysisMaxWait(-time.Second) },
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("option did not panic")
				}
			}()
			fn()
		})
	}
}
