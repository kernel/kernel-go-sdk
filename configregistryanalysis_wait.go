package kernel

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

const defaultConfigRegistryAnalysisPollInterval = 5 * time.Second

// ConfigRegistryAnalysisWaitOption configures WaitForResult.
type ConfigRegistryAnalysisWaitOption func(*configRegistryAnalysisWaitConfig)

type configRegistryAnalysisWaitConfig struct {
	pollInterval   time.Duration
	maxWait        time.Duration
	maxWaitSet     bool
	requestOptions []option.RequestOption
}

// WithConfigRegistryAnalysisPollInterval sets the delay between analysis
// retrievals. The default is five seconds.
func WithConfigRegistryAnalysisPollInterval(interval time.Duration) ConfigRegistryAnalysisWaitOption {
	if interval <= 0 {
		panic("kernel: config registry analysis poll interval must be positive")
	}
	return func(config *configRegistryAnalysisWaitConfig) {
		config.pollInterval = interval
	}
}

// WithConfigRegistryAnalysisMaxWait sets a soft polling deadline. An in-flight
// request and its normal retries may finish after the deadline. A zero duration
// performs one retrieval and then times out if the analysis is unfinished.
func WithConfigRegistryAnalysisMaxWait(maxWait time.Duration) ConfigRegistryAnalysisWaitOption {
	if maxWait < 0 {
		panic("kernel: config registry analysis max wait cannot be negative")
	}
	return func(config *configRegistryAnalysisWaitConfig) {
		config.maxWait = maxWait
		config.maxWaitSet = true
	}
}

// WithConfigRegistryAnalysisWaitRequestOptions applies request options to every
// analysis retrieval made by WaitForResult.
func WithConfigRegistryAnalysisWaitRequestOptions(opts ...option.RequestOption) ConfigRegistryAnalysisWaitOption {
	return func(config *configRegistryAnalysisWaitConfig) {
		config.requestOptions = append(config.requestOptions, opts...)
	}
}

// WaitForResult waits for an analysis to finish and returns its complete result.
// The first retrieval happens immediately. Timing out or canceling ctx stops the
// local wait without canceling the remote analysis.
func (r *ConfigRegistryAnalysisService) WaitForResult(ctx context.Context, id string, opts ...ConfigRegistryAnalysisWaitOption) (*ConfigRegistryResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("config registry analysis ID cannot be empty")
	}
	config := configRegistryAnalysisWaitConfig{pollInterval: defaultConfigRegistryAnalysisPollInterval}
	for _, opt := range opts {
		opt(&config)
	}

	startedAt := time.Now()
	var deadline time.Time
	if config.maxWaitSet {
		deadline = startedAt.Add(config.maxWait)
	}
	requestOptions := append([]option.RequestOption{}, config.requestOptions...)
	requestOptions = append(requestOptions, option.WithHeader("X-Stainless-Poll-Helper", "true"))
	polls := 0
	lastStatus := AnalysisStatus("")

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if polls > 0 && config.maxWaitSet && !time.Now().Before(deadline) {
			return nil, configRegistryAnalysisWaitTimeoutError(id, polls, lastStatus, startedAt)
		}

		response, err := r.Get(ctx, id, requestOptions...)
		if err != nil {
			return nil, err
		}
		polls++
		finished, status, err := configRegistryAnalysisFinished(response, id)
		if err != nil {
			return nil, err
		}
		lastStatus = status
		if finished {
			return response, nil
		}

		delay := time.Duration(float64(config.pollInterval) * (0.9 + rand.Float64()*0.2))
		if delay <= 0 {
			delay = time.Nanosecond
		}
		if config.maxWaitSet {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return nil, configRegistryAnalysisWaitTimeoutError(id, polls, lastStatus, startedAt)
			}
			delay = min(delay, remaining)
		}

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

func configRegistryAnalysisFinished(response *ConfigRegistryResponse, requestedID string) (bool, AnalysisStatus, error) {
	if response == nil || !response.JSON.Analysis.Valid() {
		return false, "", fmt.Errorf("config registry response for %q is missing an analysis", requestedID)
	}
	analysis := response.Analysis
	if !analysis.JSON.ID.Valid() || analysis.ID == "" {
		return false, "", fmt.Errorf("config registry response for %q has no valid analysis ID", requestedID)
	}
	if analysis.ID != requestedID {
		return false, "", fmt.Errorf("config registry response for %q returned analysis %q", requestedID, analysis.ID)
	}
	if !analysis.JSON.Status.Valid() || analysis.Status == "" {
		return false, "", fmt.Errorf("config registry analysis %q has no valid status", requestedID)
	}
	if analysis.JSON.FinishedAt.Raw() == respjson.Omitted {
		return false, "", fmt.Errorf("config registry analysis %q is missing finished_at", requestedID)
	}
	if raw := analysis.JSON.FinishedAt.Raw(); raw != respjson.Null && !analysis.JSON.FinishedAt.Valid() {
		return false, "", fmt.Errorf("config registry analysis %q has an invalid finished_at", requestedID)
	}

	if analysis.JSON.FinishedAt.Valid() {
		return true, analysis.Status, nil
	}
	switch analysis.Status {
	case AnalysisStatusCompleted, AnalysisStatusFailed, AnalysisStatusCanceled, AnalysisStatusExpired:
		return true, analysis.Status, nil
	default:
		return false, analysis.Status, nil
	}
}

func configRegistryAnalysisWaitTimeoutError(id string, polls int, lastStatus AnalysisStatus, startedAt time.Time) error {
	return fmt.Errorf(
		"timed out waiting for config registry analysis %q after %s and %d polls; last status was %q: %w",
		id,
		time.Since(startedAt).Round(time.Millisecond),
		polls,
		lastStatus,
		context.DeadlineExceeded,
	)
}
