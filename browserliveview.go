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
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Create and manage browser sessions.
//
// BrowserLiveViewService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserLiveViewService] method instead.
type BrowserLiveViewService struct {
	Options []option.RequestOption
}

// NewBrowserLiveViewService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBrowserLiveViewService(opts ...option.RequestOption) (r BrowserLiveViewService) {
	r = BrowserLiveViewService{}
	r.Options = opts
	return
}

// Creates a live view URL for a headful browser that grants either read-only or
// control access. The URL works only for the live view: it cannot connect to CDP,
// WebDriver, or any other browser endpoint. A read-only URL shows the browser but
// cannot send input, regardless of query parameters. The URL stays valid until the
// browser is deleted, or until `ttl_seconds` elapse when set. For a browser
// acquired from a pool, it stops working when the browser is released. Nothing is
// stored, so a URL cannot be revoked before then.
func (r *BrowserLiveViewService) New(ctx context.Context, idOrName string, body BrowserLiveViewNewParams, opts ...option.RequestOption) (res *BrowserLiveView, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("browsers/%s/live_view", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type BrowserLiveView struct {
	// `read` shows the browser without accepting input. `control` also lets the viewer
	// drive the browser through the live view.
	//
	// Any of "read", "control".
	Scope BrowserLiveViewScope `json:"scope" api:"required"`
	// Live view URL. Treat it as a credential.
	URL string `json:"url" api:"required"`
	// When the URL expires. Absent when it stays valid until the browser is deleted.
	ExpiresAt time.Time `json:"expires_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Scope       respjson.Field
		URL         respjson.Field
		ExpiresAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserLiveView) RawJSON() string { return r.JSON.raw }
func (r *BrowserLiveView) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserLiveViewRequestParam struct {
	// Seconds until the URL expires. Omit to keep it valid until the browser is
	// deleted.
	TtlSeconds param.Opt[int64] `json:"ttl_seconds,omitzero"`
	// `read` shows the browser without accepting input. `control` also lets the viewer
	// drive the browser through the live view.
	//
	// Any of "read", "control".
	Scope BrowserLiveViewScope `json:"scope,omitzero"`
	paramObj
}

func (r BrowserLiveViewRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserLiveViewRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserLiveViewRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `read` shows the browser without accepting input. `control` also lets the viewer
// drive the browser through the live view.
type BrowserLiveViewScope string

const (
	BrowserLiveViewScopeRead    BrowserLiveViewScope = "read"
	BrowserLiveViewScopeControl BrowserLiveViewScope = "control"
)

type BrowserLiveViewNewParams struct {
	BrowserLiveViewRequest BrowserLiveViewRequestParam
	paramObj
}

func (r BrowserLiveViewNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.BrowserLiveViewRequest)
}
func (r *BrowserLiveViewNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
