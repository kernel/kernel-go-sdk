// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// WorkloadIdentityAuthenticationEventService contains methods and other services
// that help with interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewWorkloadIdentityAuthenticationEventService] method instead.
type WorkloadIdentityAuthenticationEventService struct {
	Options []option.RequestOption
}

// NewWorkloadIdentityAuthenticationEventService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewWorkloadIdentityAuthenticationEventService(opts ...option.RequestOption) (r WorkloadIdentityAuthenticationEventService) {
	r = WorkloadIdentityAuthenticationEventService{}
	r.Options = opts
	return
}

// Any signed-in member of the organization can read this. API keys cannot.
func (r *WorkloadIdentityAuthenticationEventService) List(ctx context.Context, query WorkloadIdentityAuthenticationEventListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[WorkloadIdentityAuthenticationEvent], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "org/workload_identity_authentication_events"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// Any signed-in member of the organization can read this. API keys cannot.
func (r *WorkloadIdentityAuthenticationEventService) ListAutoPaging(ctx context.Context, query WorkloadIdentityAuthenticationEventListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[WorkloadIdentityAuthenticationEvent] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

type WorkloadIdentityAuthenticationEvent struct {
	ID               string    `json:"id" api:"required"`
	CorrelationID    string    `json:"correlation_id" api:"required"`
	CreatedAt        time.Time `json:"created_at" api:"required" format:"date-time"`
	CredentialID     string    `json:"credential_id" api:"required"`
	FederationRuleID string    `json:"federation_rule_id" api:"required"`
	IssuerID         string    `json:"issuer_id" api:"required"`
	OrgID            string    `json:"org_id" api:"required"`
	// Any of "success", "failure".
	Outcome          WorkloadIdentityAuthenticationEventOutcome `json:"outcome" api:"required"`
	ProjectID        string                                     `json:"project_id" api:"required"`
	Reason           string                                     `json:"reason" api:"required"`
	ServiceAccountID string                                     `json:"service_account_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		CorrelationID    respjson.Field
		CreatedAt        respjson.Field
		CredentialID     respjson.Field
		FederationRuleID respjson.Field
		IssuerID         respjson.Field
		OrgID            respjson.Field
		Outcome          respjson.Field
		ProjectID        respjson.Field
		Reason           respjson.Field
		ServiceAccountID respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityAuthenticationEvent) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityAuthenticationEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityAuthenticationEventOutcome string

const (
	WorkloadIdentityAuthenticationEventOutcomeSuccess WorkloadIdentityAuthenticationEventOutcome = "success"
	WorkloadIdentityAuthenticationEventOutcomeFailure WorkloadIdentityAuthenticationEventOutcome = "failure"
)

type WorkloadIdentityAuthenticationEventListParams struct {
	EndTime          param.Opt[time.Time] `query:"end_time,omitzero" format:"date-time" json:"-"`
	FederationRuleID param.Opt[string]    `query:"federation_rule_id,omitzero" json:"-"`
	IssuerID         param.Opt[string]    `query:"issuer_id,omitzero" json:"-"`
	Limit            param.Opt[int64]     `query:"limit,omitzero" json:"-"`
	Offset           param.Opt[int64]     `query:"offset,omitzero" json:"-"`
	ProjectID        param.Opt[string]    `query:"project_id,omitzero" json:"-"`
	ServiceAccountID param.Opt[string]    `query:"service_account_id,omitzero" json:"-"`
	StartTime        param.Opt[time.Time] `query:"start_time,omitzero" format:"date-time" json:"-"`
	// Any of "success", "failure".
	Outcome WorkloadIdentityAuthenticationEventListParamsOutcome `query:"outcome,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [WorkloadIdentityAuthenticationEventListParams]'s query
// parameters as `url.Values`.
func (r WorkloadIdentityAuthenticationEventListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type WorkloadIdentityAuthenticationEventListParamsOutcome string

const (
	WorkloadIdentityAuthenticationEventListParamsOutcomeSuccess WorkloadIdentityAuthenticationEventListParamsOutcome = "success"
	WorkloadIdentityAuthenticationEventListParamsOutcomeFailure WorkloadIdentityAuthenticationEventListParamsOutcome = "failure"
)
