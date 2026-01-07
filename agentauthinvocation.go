// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/onkernel/kernel-go-sdk/internal/apijson"
	shimjson "github.com/onkernel/kernel-go-sdk/internal/encoding/json"
	"github.com/onkernel/kernel-go-sdk/internal/requestconfig"
	"github.com/onkernel/kernel-go-sdk/option"
	"github.com/onkernel/kernel-go-sdk/packages/param"
	"github.com/onkernel/kernel-go-sdk/packages/respjson"
)

// AgentAuthInvocationService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAgentAuthInvocationService] method instead.
type AgentAuthInvocationService struct {
	Options []option.RequestOption
}

// NewAgentAuthInvocationService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAgentAuthInvocationService(opts ...option.RequestOption) (r AgentAuthInvocationService) {
	r = AgentAuthInvocationService{}
	r.Options = opts
	return
}

// Creates a new authentication invocation for the specified auth agent. This
// starts the auth flow and returns a hosted URL for the user to complete
// authentication.
func (r *AgentAuthInvocationService) New(ctx context.Context, body AgentAuthInvocationNewParams, opts ...option.RequestOption) (res *AuthAgentInvocationCreateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "agents/auth/invocations"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Returns invocation details including status, app_name, and domain. Supports both
// API key and JWT (from exchange endpoint) authentication.
func (r *AgentAuthInvocationService) Get(ctx context.Context, invocationID string, opts ...option.RequestOption) (res *AgentAuthInvocationResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if invocationID == "" {
		err = errors.New("missing required invocation_id parameter")
		return
	}
	path := fmt.Sprintf("agents/auth/invocations/%s", invocationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Validates the handoff code and returns a JWT token for subsequent requests. No
// authentication required (the handoff code serves as the credential).
func (r *AgentAuthInvocationService) Exchange(ctx context.Context, invocationID string, body AgentAuthInvocationExchangeParams, opts ...option.RequestOption) (res *AgentAuthInvocationExchangeResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if invocationID == "" {
		err = errors.New("missing required invocation_id parameter")
		return
	}
	path := fmt.Sprintf("agents/auth/invocations/%s/exchange", invocationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Submits field values for the discovered login form. Returns immediately after
// submission is accepted. Poll the invocation endpoint to track progress and get
// results.
func (r *AgentAuthInvocationService) Submit(ctx context.Context, invocationID string, body AgentAuthInvocationSubmitParams, opts ...option.RequestOption) (res *AgentAuthSubmitResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if invocationID == "" {
		err = errors.New("missing required invocation_id parameter")
		return
	}
	path := fmt.Sprintf("agents/auth/invocations/%s/submit", invocationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Response from exchange endpoint
type AgentAuthInvocationExchangeResponse struct {
	// Invocation ID
	InvocationID string `json:"invocation_id,required"`
	// JWT token with invocation_id claim (30 minute TTL)
	Jwt string `json:"jwt,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InvocationID respjson.Field
		Jwt          respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentAuthInvocationExchangeResponse) RawJSON() string { return r.JSON.raw }
func (r *AgentAuthInvocationExchangeResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AgentAuthInvocationNewParams struct {
	// Request to create an invocation for an existing auth agent
	AuthAgentInvocationCreateRequest AuthAgentInvocationCreateRequestParam
	paramObj
}

func (r AgentAuthInvocationNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.AuthAgentInvocationCreateRequest)
}
func (r *AgentAuthInvocationNewParams) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.AuthAgentInvocationCreateRequest)
}

type AgentAuthInvocationExchangeParams struct {
	// Handoff code from start endpoint
	Code string `json:"code,required"`
	paramObj
}

func (r AgentAuthInvocationExchangeParams) MarshalJSON() (data []byte, err error) {
	type shadow AgentAuthInvocationExchangeParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentAuthInvocationExchangeParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AgentAuthInvocationSubmitParams struct {

	//
	// Request body variants
	//

	// This field is a request body variant, only one variant field can be set.
	OfFieldValues *AgentAuthInvocationSubmitParamsBodyFieldValues `json:",inline"`
	// This field is a request body variant, only one variant field can be set.
	OfSSOButton *AgentAuthInvocationSubmitParamsBodySSOButton `json:",inline"`

	paramObj
}

func (u AgentAuthInvocationSubmitParams) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFieldValues, u.OfSSOButton)
}
func (r *AgentAuthInvocationSubmitParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property FieldValues is required.
type AgentAuthInvocationSubmitParamsBodyFieldValues struct {
	// Values for the discovered login fields
	FieldValues map[string]string `json:"field_values,omitzero,required"`
	paramObj
}

func (r AgentAuthInvocationSubmitParamsBodyFieldValues) MarshalJSON() (data []byte, err error) {
	type shadow AgentAuthInvocationSubmitParamsBodyFieldValues
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentAuthInvocationSubmitParamsBodyFieldValues) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property SSOButton is required.
type AgentAuthInvocationSubmitParamsBodySSOButton struct {
	// Selector of SSO button to click
	SSOButton string `json:"sso_button,required"`
	paramObj
}

func (r AgentAuthInvocationSubmitParamsBodySSOButton) MarshalJSON() (data []byte, err error) {
	type shadow AgentAuthInvocationSubmitParamsBodySSOButton
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentAuthInvocationSubmitParamsBodySSOButton) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
