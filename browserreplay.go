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
)

// BrowserReplayService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserReplayService] method instead.
type BrowserReplayService struct {
	Options []option.RequestOption
}

// NewBrowserReplayService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBrowserReplayService(opts ...option.RequestOption) (r BrowserReplayService) {
	r = BrowserReplayService{}
	r.Options = opts
	return
}

// List all replays for the specified browser session.
func (r *BrowserReplayService) List(ctx context.Context, id string, opts ...option.RequestOption) (res *[]BrowserReplayListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/replays", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Download or stream the specified replay recording.
func (r *BrowserReplayService) Download(ctx context.Context, replayID string, query BrowserReplayDownloadParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "video/mp4")}, opts...)
	if query.ID == "" {
		err = errors.New("missing required id parameter")
		return
	}
	if replayID == "" {
		err = errors.New("missing required replay_id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/replays/%s", query.ID, replayID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Start recording the browser session and return a replay ID.
func (r *BrowserReplayService) Start(ctx context.Context, id string, body BrowserReplayStartParams, opts ...option.RequestOption) (res *BrowserReplayStartResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/replays", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Stop the specified replay recording and persist the video.
func (r *BrowserReplayService) Stop(ctx context.Context, replayID string, body BrowserReplayStopParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if body.ID == "" {
		err = errors.New("missing required id parameter")
		return
	}
	if replayID == "" {
		err = errors.New("missing required replay_id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/replays/%s/stop", body.ID, replayID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, nil, opts...)
	return
}

// Information about a browser replay recording.
type BrowserReplayListResponse struct {
	// Unique identifier for the replay recording.
	ReplayID string `json:"replay_id,required"`
	// Timestamp when replay finished
	FinishedAt time.Time `json:"finished_at,nullable" format:"date-time"`
	// URL for viewing the replay recording.
	ReplayViewURL string `json:"replay_view_url"`
	// Timestamp when replay started
	StartedAt time.Time `json:"started_at,nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ReplayID      respjson.Field
		FinishedAt    respjson.Field
		ReplayViewURL respjson.Field
		StartedAt     respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserReplayListResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserReplayListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Information about a browser replay recording.
type BrowserReplayStartResponse struct {
	// Unique identifier for the replay recording.
	ReplayID string `json:"replay_id,required"`
	// Timestamp when replay finished
	FinishedAt time.Time `json:"finished_at,nullable" format:"date-time"`
	// URL for viewing the replay recording.
	ReplayViewURL string `json:"replay_view_url"`
	// Timestamp when replay started
	StartedAt time.Time `json:"started_at,nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ReplayID      respjson.Field
		FinishedAt    respjson.Field
		ReplayViewURL respjson.Field
		StartedAt     respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserReplayStartResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserReplayStartResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserReplayDownloadParams struct {
	ID string `path:"id,required" json:"-"`
	paramObj
}

type BrowserReplayStartParams struct {
	// Recording framerate in fps.
	Framerate param.Opt[int64] `json:"framerate,omitzero"`
	// Maximum recording duration in seconds.
	MaxDurationInSeconds param.Opt[int64] `json:"max_duration_in_seconds,omitzero"`
	paramObj
}

func (r BrowserReplayStartParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserReplayStartParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserReplayStartParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserReplayStopParams struct {
	ID string `path:"id,required" json:"-"`
	paramObj
}
