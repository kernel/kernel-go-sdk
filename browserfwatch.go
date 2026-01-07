// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/onkernel/kernel-go-sdk/internal/apijson"
	"github.com/onkernel/kernel-go-sdk/internal/requestconfig"
	"github.com/onkernel/kernel-go-sdk/option"
	"github.com/onkernel/kernel-go-sdk/packages/param"
	"github.com/onkernel/kernel-go-sdk/packages/respjson"
	"github.com/onkernel/kernel-go-sdk/packages/ssestream"
)

// BrowserFWatchService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserFWatchService] method instead.
type BrowserFWatchService struct {
	Options []option.RequestOption
}

// NewBrowserFWatchService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBrowserFWatchService(opts ...option.RequestOption) (r BrowserFWatchService) {
	r = BrowserFWatchService{}
	r.Options = opts
	return
}

// Stream filesystem events for a watch
func (r *BrowserFWatchService) EventsStreaming(ctx context.Context, watchID string, query BrowserFWatchEventsParams, opts ...option.RequestOption) (stream *ssestream.Stream[BrowserFWatchEventsResponse]) {
	var (
		raw *http.Response
		err error
	)
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "text/event-stream")}, opts...)
	if query.ID == "" {
		err = errors.New("missing required id parameter")
		return
	}
	if watchID == "" {
		err = errors.New("missing required watch_id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/watch/%s/events", query.ID, watchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &raw, opts...)
	return ssestream.NewStream[BrowserFWatchEventsResponse](ssestream.NewDecoder(raw), err)
}

// Watch a directory for changes
func (r *BrowserFWatchService) Start(ctx context.Context, id string, body BrowserFWatchStartParams, opts ...option.RequestOption) (res *BrowserFWatchStartResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/watch", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Stop watching a directory
func (r *BrowserFWatchService) Stop(ctx context.Context, watchID string, body BrowserFWatchStopParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if body.ID == "" {
		err = errors.New("missing required id parameter")
		return
	}
	if watchID == "" {
		err = errors.New("missing required watch_id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/watch/%s", body.ID, watchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return
}

// Filesystem change event.
type BrowserFWatchEventsResponse struct {
	// Absolute path of the file or directory.
	Path string `json:"path,required"`
	// Event type.
	//
	// Any of "CREATE", "WRITE", "DELETE", "RENAME".
	Type BrowserFWatchEventsResponseType `json:"type,required"`
	// Whether the affected path is a directory.
	IsDir bool `json:"is_dir"`
	// Base name of the file or directory affected.
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Path        respjson.Field
		Type        respjson.Field
		IsDir       respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserFWatchEventsResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserFWatchEventsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event type.
type BrowserFWatchEventsResponseType string

const (
	BrowserFWatchEventsResponseTypeCreate BrowserFWatchEventsResponseType = "CREATE"
	BrowserFWatchEventsResponseTypeWrite  BrowserFWatchEventsResponseType = "WRITE"
	BrowserFWatchEventsResponseTypeDelete BrowserFWatchEventsResponseType = "DELETE"
	BrowserFWatchEventsResponseTypeRename BrowserFWatchEventsResponseType = "RENAME"
)

type BrowserFWatchStartResponse struct {
	// Unique identifier for the directory watch
	WatchID string `json:"watch_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		WatchID     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserFWatchStartResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserFWatchStartResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserFWatchEventsParams struct {
	ID string `path:"id,required" json:"-"`
	paramObj
}

type BrowserFWatchStartParams struct {
	// Directory to watch.
	Path string `json:"path,required"`
	// Whether to watch recursively.
	Recursive param.Opt[bool] `json:"recursive,omitzero"`
	paramObj
}

func (r BrowserFWatchStartParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserFWatchStartParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserFWatchStartParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserFWatchStopParams struct {
	ID string `path:"id,required" json:"-"`
	paramObj
}
