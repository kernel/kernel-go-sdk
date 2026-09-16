// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
	"github.com/kernel/kernel-go-sdk/shared/constant"
)

// VaultItemService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewVaultItemService] method instead.
type VaultItemService struct {
	Options []option.RequestOption
}

// NewVaultItemService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewVaultItemService(opts ...option.RequestOption) (r VaultItemService) {
	r = VaultItemService{}
	r.Options = opts
	return
}

// The response advertises operations that are valid in the item's current state
// and live data that can be requested through `expand`. Read each operation's
// description before using it. Expanded data is fetched from the provider and is
// not persisted in the vault item. Requesting an unavailable expansion returns 409
// instead of a partial item. Pending credential items return a collection action.
// Kernel-hosted active collection links are renewed atomically on expiry for ready
// or pending items without changing the item version. Invoke collect to open a
// form for a ready item without clearing values. Sensitive credential values are
// never returned.
func (r *VaultItemService) Get(ctx context.Context, key string, params VaultItemGetParams, opts ...option.RequestOption) (res *VaultItemUnion, err error) {
	opts = slices.Concat(r.Options, opts)
	if params.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	if key == "" {
		err = errors.New("missing required key parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/items/%s", params.IDOrName, key)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Credential updates require type credential and the current version, and change
// only values or description; omitted values are preserved, nonempty strings
// replace, and null or empty strings clear supported fields. Clearing required
// text/email/password values returns pending_collection; browser forms still
// require nonempty required inputs. Card updates may omit type for compatibility
// with legacy requests. Requested cards accept a replacement specification.
// Pending issuance requests may update provider-supported fields on their existing
// request, subject to atomic provider approval checks; omitted optional fields
// remain unchanged and explicit empty lists clear them. Wallet/provider binding
// and unsupported fields cannot change after authorization starts. An uncertain
// update enters recovery_required and must not be retried. Checkout cards may be
// edited between authorizations.
func (r *VaultItemService) Update(ctx context.Context, key string, params VaultItemUpdateParams, opts ...option.RequestOption) (res *VaultItemUnion, err error) {
	opts = slices.Concat(r.Options, opts)
	if params.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	if key == "" {
		err = errors.New("missing required key parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/items/%s", params.IDOrName, key)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, params, &res, opts...)
	return res, err
}

// Credential entries include safe field metadata and non-sensitive values. Listing
// never creates or renews collection sessions; only an existing unexpired active
// session is included. Use single-item GET or collect to obtain a fresh link.
func (r *VaultItemService) List(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *[]VaultItemUnion, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/items", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Unresolved payment operations normally block deletion, including operations on
// child cards of a wallet. An AgentCard card in recovery_required whose checkout
// create response returned no authorization ID may be explicitly abandoned by
// deleting that card directly; deleting its wallet or vault remains blocked.
// Deleting or recreating an item is not proof that a payment did not occur.
func (r *VaultItemService) Delete(ctx context.Context, key string, body VaultItemDeleteParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if body.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return err
	}
	if key == "" {
		err = errors.New("missing required key parameter")
		return err
	}
	path := fmt.Sprintf("vaults/%s/items/%s", body.IDOrName, key)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// List immutable audit events for a vault item
func (r *VaultItemService) Events(ctx context.Context, key string, params VaultItemEventsParams, opts ...option.RequestOption) (res *[]VaultItemEvent, err error) {
	opts = slices.Concat(r.Options, opts)
	if params.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	if key == "" {
		err = errors.New("missing required key parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/items/%s/events", params.IDOrName, key)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Retrieve the item first and invoke only an operation listed in
// `available_operations`, following its natural-language description. Availability
// is rechecked at execution time; unavailable operations return 409. Authorization
// and preparation may call an external provider and return updated state. Link
// cards advertise authorize without checkout context. Eligible unused AgentCard
// cards advertise prepare_checkout, which requires checkout context and obtains
// device approval before native Square Pay. Keep the returned approval page open,
// poll until ready_to_submit, then submit before preparation.expires_at. Unused
// preparations expire automatically and cannot be reused. If spend-request
// creation is rate limited, returns HTTP 429 with code
// `spend_request_rate_limited`; stop and back off before retrying.
//
// Fill returns a value-free execution result. Validation failures before writing
// return 400 (invalid request or targets), 403 (access or destination denied), 404
// (resource not found), or 409 (item or browser not ready). Once writing starts,
// known partial failures and indeterminate field outcomes return 200 with status
// `failed` or `unknown`, not an automatic-retry signal. A transport error may
// leave the outcome unknown; do not automatically retry.
func (r *VaultItemService) PerformOperation(ctx context.Context, key string, params VaultItemPerformOperationParams, opts ...option.RequestOption) (res *VaultItemOperationResponseUnion, err error) {
	opts = slices.Concat(r.Options, []option.RequestOption{option.WithMaxRetries(0)}, opts)
	if params.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	if key == "" {
		err = errors.New("missing required key parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/items/%s/operations", params.IDOrName, key)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Create an item under a key unique within its vault, or retrieve the existing
// item when its specification matches. An identical card PUT returns the existing
// card in any lifecycle state without polling the provider, reauthorizing,
// replacing aliases, or resetting recovery. Conflicting specifications return 409.
// Provider-specific authorization requirements and retry behavior are described in
// the item's request schema. Do not use credential items to store, collect, or
// fill credit card data, including card numbers (PANs), security codes (CVV/CVC),
// or expiration dates. Use wallet and card item types for credit cards and payment
// checkout instead.
func (r *VaultItemService) Upsert(ctx context.Context, key string, params VaultItemUpsertParams, opts ...option.RequestOption) (res *VaultItemUnion, err error) {
	opts = slices.Concat(r.Options, opts)
	if params.IDOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	if key == "" {
		err = errors.New("missing required key parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/items/%s", params.IDOrName, key)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, params, &res, opts...)
	return res, err
}

// The in-flight or most recent checkout authorization. Present while a checkout is
// pending approval and after it settles.
type AgentcardCheckoutAuthorization struct {
	ID          string    `json:"id" api:"required"`
	AmountCents int64     `json:"amount_cents" api:"required"`
	CreatedAt   time.Time `json:"created_at" api:"required" format:"date-time"`
	Currency    string    `json:"currency" api:"required"`
	Merchant    string    `json:"merchant" api:"required"`
	Psp         string    `json:"psp" api:"required"`
	// Any of "awaiting_approval", "approved", "declined", "expired".
	Status      AgentcardCheckoutAuthorizationStatus `json:"status" api:"required"`
	ActualCents int64                                `json:"actual_cents"`
	// Display amount shown on the approval screen.
	Amount string `json:"amount"`
	// Any of "display_only", "stripe_payment_intent".
	AmountAuthority AgentcardCheckoutAuthorizationAmountAuthority `json:"amount_authority"`
	AmountVerified  bool                                          `json:"amount_verified"`
	ApprovalURL     string                                        `json:"approval_url" format:"uri"`
	// Browser session that submitted the checkout.
	BrowserID          string `json:"browser_id"`
	ChargedAmountCents int64  `json:"charged_amount_cents"`
	ChargedCurrency    string `json:"charged_currency"`
	// Any of "captured", "authorized", "none".
	ChargedKind     AgentcardCheckoutAuthorizationChargedKind `json:"charged_kind"`
	ExpectedCents   int64                                     `json:"expected_cents"`
	ExpiresAt       time.Time                                 `json:"expires_at" format:"date-time"`
	PspErrorCode    string                                    `json:"psp_error_code"`
	Reason          string                                    `json:"reason"`
	ReplayAttempted bool                                      `json:"replay_attempted"`
	// Whether the processor response was delivered to the browser.
	ReplayDelivered bool `json:"replay_delivered"`
	// HTTP status of the replayed processor response.
	ReplayStatus int64 `json:"replay_status"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                 respjson.Field
		AmountCents        respjson.Field
		CreatedAt          respjson.Field
		Currency           respjson.Field
		Merchant           respjson.Field
		Psp                respjson.Field
		Status             respjson.Field
		ActualCents        respjson.Field
		Amount             respjson.Field
		AmountAuthority    respjson.Field
		AmountVerified     respjson.Field
		ApprovalURL        respjson.Field
		BrowserID          respjson.Field
		ChargedAmountCents respjson.Field
		ChargedCurrency    respjson.Field
		ChargedKind        respjson.Field
		ExpectedCents      respjson.Field
		ExpiresAt          respjson.Field
		PspErrorCode       respjson.Field
		Reason             respjson.Field
		ReplayAttempted    respjson.Field
		ReplayDelivered    respjson.Field
		ReplayStatus       respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentcardCheckoutAuthorization) RawJSON() string { return r.JSON.raw }
func (r *AgentcardCheckoutAuthorization) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AgentcardCheckoutAuthorizationStatus string

const (
	AgentcardCheckoutAuthorizationStatusAwaitingApproval AgentcardCheckoutAuthorizationStatus = "awaiting_approval"
	AgentcardCheckoutAuthorizationStatusApproved         AgentcardCheckoutAuthorizationStatus = "approved"
	AgentcardCheckoutAuthorizationStatusDeclined         AgentcardCheckoutAuthorizationStatus = "declined"
	AgentcardCheckoutAuthorizationStatusExpired          AgentcardCheckoutAuthorizationStatus = "expired"
)

type AgentcardCheckoutAuthorizationAmountAuthority string

const (
	AgentcardCheckoutAuthorizationAmountAuthorityDisplayOnly         AgentcardCheckoutAuthorizationAmountAuthority = "display_only"
	AgentcardCheckoutAuthorizationAmountAuthorityStripePaymentIntent AgentcardCheckoutAuthorizationAmountAuthority = "stripe_payment_intent"
)

type AgentcardCheckoutAuthorizationChargedKind string

const (
	AgentcardCheckoutAuthorizationChargedKindCaptured   AgentcardCheckoutAuthorizationChargedKind = "captured"
	AgentcardCheckoutAuthorizationChargedKindAuthorized AgentcardCheckoutAuthorizationChargedKind = "authorized"
	AgentcardCheckoutAuthorizationChargedKindNone       AgentcardCheckoutAuthorizationChargedKind = "none"
)

// One-use processor-bound checkout preparation. Keep the approval page open
// through token handoff. The amount is display-only and does not constrain the
// merchant's eventual charge.
type AgentcardCheckoutPreparation struct {
	BrowserID string    `json:"browser_id" api:"required"`
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Any of "production", "sandbox", "shared".
	Environment    AgentcardCheckoutPreparationEnvironment `json:"environment" api:"required"`
	MerchantOrigin string                                  `json:"merchant_origin" api:"required"`
	// Any of "square", "braintree", "worldpay", "bambora", "mercado_pago".
	Psp AgentcardPreparedProcessor `json:"psp" api:"required"`
	// Preparation consumed means egress claimed the preparation and it cannot be
	// reused. It does not mean the attempt settled. Use the enclosing item's status as
	// the lifecycle indicator; item consumed means the attempt settled, not that an
	// order or charge succeeded.
	//
	// Any of "creating", "awaiting_approval", "ready", "consumed", "cancelled",
	// "expired", "unknown".
	Status      AgentcardCheckoutPreparationStatus `json:"status" api:"required"`
	ID          string                             `json:"id"`
	ApprovalURL string                             `json:"approval_url" format:"uri"`
	// When ready, the absolute deadline to submit the first native request; no later
	// than provider readiness expiry or 30 seconds after Kernel first observes
	// readiness. Polling never extends this deadline.
	ExpiresAt time.Time `json:"expires_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BrowserID      respjson.Field
		CreatedAt      respjson.Field
		Environment    respjson.Field
		MerchantOrigin respjson.Field
		Psp            respjson.Field
		Status         respjson.Field
		ID             respjson.Field
		ApprovalURL    respjson.Field
		ExpiresAt      respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentcardCheckoutPreparation) RawJSON() string { return r.JSON.raw }
func (r *AgentcardCheckoutPreparation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AgentcardCheckoutPreparationEnvironment string

const (
	AgentcardCheckoutPreparationEnvironmentProduction AgentcardCheckoutPreparationEnvironment = "production"
	AgentcardCheckoutPreparationEnvironmentSandbox    AgentcardCheckoutPreparationEnvironment = "sandbox"
	AgentcardCheckoutPreparationEnvironmentShared     AgentcardCheckoutPreparationEnvironment = "shared"
)

// Preparation consumed means egress claimed the preparation and it cannot be
// reused. It does not mean the attempt settled. Use the enclosing item's status as
// the lifecycle indicator; item consumed means the attempt settled, not that an
// order or charge succeeded.
type AgentcardCheckoutPreparationStatus string

const (
	AgentcardCheckoutPreparationStatusCreating         AgentcardCheckoutPreparationStatus = "creating"
	AgentcardCheckoutPreparationStatusAwaitingApproval AgentcardCheckoutPreparationStatus = "awaiting_approval"
	AgentcardCheckoutPreparationStatusReady            AgentcardCheckoutPreparationStatus = "ready"
	AgentcardCheckoutPreparationStatusConsumed         AgentcardCheckoutPreparationStatus = "consumed"
	AgentcardCheckoutPreparationStatusCancelled        AgentcardCheckoutPreparationStatus = "cancelled"
	AgentcardCheckoutPreparationStatusExpired          AgentcardCheckoutPreparationStatus = "expired"
	AgentcardCheckoutPreparationStatusUnknown          AgentcardCheckoutPreparationStatus = "unknown"
)

type AgentcardPreparedProcessor string

const (
	AgentcardPreparedProcessorSquare      AgentcardPreparedProcessor = "square"
	AgentcardPreparedProcessorBraintree   AgentcardPreparedProcessor = "braintree"
	AgentcardPreparedProcessorWorldpay    AgentcardPreparedProcessor = "worldpay"
	AgentcardPreparedProcessorBambora     AgentcardPreparedProcessor = "bambora"
	AgentcardPreparedProcessorMercadoPago AgentcardPreparedProcessor = "mercado_pago"
)

// Authorize a Link card using its existing purchase specification. Use only after
// explicit user approval and when the item advertises authorize. Do not
// automatically retry provider failures or indeterminate outcomes. Checkout
// context is not accepted.
//
// The property Type is required.
type AuthorizeVaultItemOperationRequestParam struct {
	// Any of "authorize".
	Type AuthorizeVaultItemOperationRequestType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r AuthorizeVaultItemOperationRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow AuthorizeVaultItemOperationRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AuthorizeVaultItemOperationRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthorizeVaultItemOperationRequestType string

const (
	AuthorizeVaultItemOperationRequestTypeAuthorize AuthorizeVaultItemOperationRequestType = "authorize"
)

// CardVaultItemSpecUnion contains all possible properties and values from
// [CardVaultItemSpecLink], [CardVaultItemSpecAgentcard].
//
// Use the [CardVaultItemSpecUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type CardVaultItemSpecUnion struct {
	Amount int64 `json:"amount"`
	// This field is from variant [CardVaultItemSpecLink].
	Context  string `json:"context"`
	Currency string `json:"currency"`
	// This field is from variant [CardVaultItemSpecLink].
	MerchantName string `json:"merchant_name"`
	// This field is from variant [CardVaultItemSpecLink].
	MerchantURL string `json:"merchant_url"`
	// This field is from variant [CardVaultItemSpecLink].
	PaymentMethodID string `json:"payment_method_id"`
	// Any of "link", "agentcard".
	Provider string `json:"provider"`
	Wallet   string `json:"wallet"`
	// This field is from variant [CardVaultItemSpecLink].
	ExpiresAt int64 `json:"expires_at"`
	// This field is from variant [CardVaultItemSpecLink].
	LineItems []CardVaultItemSpecLinkLineItem `json:"line_items"`
	// This field is from variant [CardVaultItemSpecLink].
	Metadata map[string]string `json:"metadata"`
	// This field is from variant [CardVaultItemSpecLink].
	Totals []CardVaultItemSpecLinkTotal `json:"totals"`
	// This field is from variant [CardVaultItemSpecAgentcard].
	Merchant string `json:"merchant"`
	// This field is from variant [CardVaultItemSpecAgentcard].
	CardID string `json:"card_id"`
	JSON   struct {
		Amount          respjson.Field
		Context         respjson.Field
		Currency        respjson.Field
		MerchantName    respjson.Field
		MerchantURL     respjson.Field
		PaymentMethodID respjson.Field
		Provider        respjson.Field
		Wallet          respjson.Field
		ExpiresAt       respjson.Field
		LineItems       respjson.Field
		Metadata        respjson.Field
		Totals          respjson.Field
		Merchant        respjson.Field
		CardID          respjson.Field
		raw             string
	} `json:"-"`
}

// anyCardVaultItemSpec is implemented by each variant of [CardVaultItemSpecUnion]
// to add type safety for the return type of [CardVaultItemSpecUnion.AsAny]
type anyCardVaultItemSpec interface {
	implCardVaultItemSpecUnion()
}

func (CardVaultItemSpecLink) implCardVaultItemSpecUnion()      {}
func (CardVaultItemSpecAgentcard) implCardVaultItemSpecUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := CardVaultItemSpecUnion.AsAny().(type) {
//	case kernel.CardVaultItemSpecLink:
//	case kernel.CardVaultItemSpecAgentcard:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u CardVaultItemSpecUnion) AsAny() anyCardVaultItemSpec {
	switch u.Provider {
	case "link":
		return u.AsLink()
	case "agentcard":
		return u.AsAgentcard()
	}
	return nil
}

func (u CardVaultItemSpecUnion) AsLink() (v CardVaultItemSpecLink) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CardVaultItemSpecUnion) AsAgentcard() (v CardVaultItemSpecAgentcard) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u CardVaultItemSpecUnion) RawJSON() string { return u.JSON.raw }

func (r *CardVaultItemSpecUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this CardVaultItemSpecUnion to a CardVaultItemSpecUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// CardVaultItemSpecUnionParam.Overrides()
func (r CardVaultItemSpecUnion) ToParam() CardVaultItemSpecUnionParam {
	return param.Override[CardVaultItemSpecUnionParam](json.RawMessage(r.RawJSON()))
}

// Live payment card. Test-mode card creation is not supported.
type CardVaultItemSpecLink struct {
	// Integer amount in minor currency units.
	Amount       int64  `json:"amount" api:"required"`
	Context      string `json:"context" api:"required"`
	Currency     string `json:"currency" api:"required"`
	MerchantName string `json:"merchant_name" api:"required"`
	MerchantURL  string `json:"merchant_url" api:"required" format:"uri"`
	// Payment-method ID returned by the referenced wallet's payment-method listing.
	// The provider decides whether the selected funding method can satisfy the card
	// request.
	PaymentMethodID string        `json:"payment_method_id" api:"required"`
	Provider        constant.Link `json:"provider" default:"link"`
	// Wallet item key used to mint this card.
	Wallet    string                          `json:"wallet" api:"required"`
	ExpiresAt int64                           `json:"expires_at"`
	LineItems []CardVaultItemSpecLinkLineItem `json:"line_items"`
	Metadata  map[string]string               `json:"metadata"`
	Totals    []CardVaultItemSpecLinkTotal    `json:"totals"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Amount          respjson.Field
		Context         respjson.Field
		Currency        respjson.Field
		MerchantName    respjson.Field
		MerchantURL     respjson.Field
		PaymentMethodID respjson.Field
		Provider        respjson.Field
		Wallet          respjson.Field
		ExpiresAt       respjson.Field
		LineItems       respjson.Field
		Metadata        respjson.Field
		Totals          respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CardVaultItemSpecLink) RawJSON() string { return r.JSON.raw }
func (r *CardVaultItemSpecLink) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CardVaultItemSpecLinkLineItem struct {
	Name        string                               `json:"name" api:"required"`
	Description string                               `json:"description"`
	ImageURL    string                               `json:"image_url"`
	ProductURL  string                               `json:"product_url"`
	Quantity    int64                                `json:"quantity"`
	SKU         string                               `json:"sku"`
	Totals      []CardVaultItemSpecLinkLineItemTotal `json:"totals"`
	// Unit amount in minor currency units.
	UnitAmount int64  `json:"unit_amount"`
	URL        string `json:"url"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Description respjson.Field
		ImageURL    respjson.Field
		ProductURL  respjson.Field
		Quantity    respjson.Field
		SKU         respjson.Field
		Totals      respjson.Field
		UnitAmount  respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CardVaultItemSpecLinkLineItem) RawJSON() string { return r.JSON.raw }
func (r *CardVaultItemSpecLinkLineItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CardVaultItemSpecLinkLineItemTotal struct {
	// Total amount in minor currency units.
	Amount      int64  `json:"amount" api:"required"`
	DisplayText string `json:"display_text" api:"required"`
	Type        string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Amount      respjson.Field
		DisplayText respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CardVaultItemSpecLinkLineItemTotal) RawJSON() string { return r.JSON.raw }
func (r *CardVaultItemSpecLinkLineItemTotal) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CardVaultItemSpecLinkTotal struct {
	// Total amount in minor currency units.
	Amount      int64  `json:"amount" api:"required"`
	DisplayText string `json:"display_text" api:"required"`
	Type        string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Amount      respjson.Field
		DisplayText respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CardVaultItemSpecLinkTotal) RawJSON() string { return r.JSON.raw }
func (r *CardVaultItemSpecLinkTotal) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentCard reusable live payment card. Test-mode card creation is not supported.
// Each checkout creates an approval-gated authorization for spec.merchant /
// spec.amount. The card stays ready after each authorization.
type CardVaultItemSpecAgentcard struct {
	// Integer amount in minor currency units.
	Amount   int64  `json:"amount" api:"required"`
	Currency string `json:"currency" api:"required"`
	// Merchant name shown on the cardholder's approval screen.
	Merchant string             `json:"merchant" api:"required"`
	Provider constant.Agentcard `json:"provider" default:"agentcard"`
	// Wallet item key used to authorize checkouts.
	Wallet string `json:"wallet" api:"required"`
	// AgentCard vaulted card to pay with. Omitted, the cardholder picks on the
	// approval screen.
	CardID string `json:"card_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Amount      respjson.Field
		Currency    respjson.Field
		Merchant    respjson.Field
		Provider    respjson.Field
		Wallet      respjson.Field
		CardID      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CardVaultItemSpecAgentcard) RawJSON() string { return r.JSON.raw }
func (r *CardVaultItemSpecAgentcard) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type CardVaultItemSpecUnionParam struct {
	OfLink      *CardVaultItemSpecLinkParam      `json:",omitzero,inline"`
	OfAgentcard *CardVaultItemSpecAgentcardParam `json:",omitzero,inline"`
	paramUnion
}

func (u CardVaultItemSpecUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLink, u.OfAgentcard)
}
func (u *CardVaultItemSpecUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *CardVaultItemSpecUnionParam) asAny() any {
	if !param.IsOmitted(u.OfLink) {
		return u.OfLink
	} else if !param.IsOmitted(u.OfAgentcard) {
		return u.OfAgentcard
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetContext() *string {
	if vt := u.OfLink; vt != nil {
		return &vt.Context
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetMerchantName() *string {
	if vt := u.OfLink; vt != nil {
		return &vt.MerchantName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetMerchantURL() *string {
	if vt := u.OfLink; vt != nil {
		return &vt.MerchantURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetPaymentMethodID() *string {
	if vt := u.OfLink; vt != nil {
		return &vt.PaymentMethodID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetExpiresAt() *int64 {
	if vt := u.OfLink; vt != nil && vt.ExpiresAt.Valid() {
		return &vt.ExpiresAt.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetLineItems() []CardVaultItemSpecLinkLineItemParam {
	if vt := u.OfLink; vt != nil {
		return vt.LineItems
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetMetadata() map[string]string {
	if vt := u.OfLink; vt != nil {
		return vt.Metadata
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetTotals() []CardVaultItemSpecLinkTotalParam {
	if vt := u.OfLink; vt != nil {
		return vt.Totals
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetMerchant() *string {
	if vt := u.OfAgentcard; vt != nil {
		return &vt.Merchant
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetCardID() *string {
	if vt := u.OfAgentcard; vt != nil && vt.CardID.Valid() {
		return &vt.CardID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetAmount() *int64 {
	if vt := u.OfLink; vt != nil {
		return (*int64)(&vt.Amount)
	} else if vt := u.OfAgentcard; vt != nil {
		return (*int64)(&vt.Amount)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetCurrency() *string {
	if vt := u.OfLink; vt != nil {
		return (*string)(&vt.Currency)
	} else if vt := u.OfAgentcard; vt != nil {
		return (*string)(&vt.Currency)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetProvider() *string {
	if vt := u.OfLink; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfAgentcard; vt != nil {
		return (*string)(&vt.Provider)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CardVaultItemSpecUnionParam) GetWallet() *string {
	if vt := u.OfLink; vt != nil {
		return (*string)(&vt.Wallet)
	} else if vt := u.OfAgentcard; vt != nil {
		return (*string)(&vt.Wallet)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[CardVaultItemSpecUnionParam](
		"provider",
		apijson.Discriminator[CardVaultItemSpecLinkParam]("link"),
		apijson.Discriminator[CardVaultItemSpecAgentcardParam]("agentcard"),
	)
}

// Live payment card. Test-mode card creation is not supported.
//
// The properties Amount, Context, Currency, MerchantName, MerchantURL,
// PaymentMethodID, Provider, Wallet are required.
type CardVaultItemSpecLinkParam struct {
	// Integer amount in minor currency units.
	Amount       int64  `json:"amount" api:"required"`
	Context      string `json:"context" api:"required"`
	Currency     string `json:"currency" api:"required"`
	MerchantName string `json:"merchant_name" api:"required"`
	MerchantURL  string `json:"merchant_url" api:"required" format:"uri"`
	// Payment-method ID returned by the referenced wallet's payment-method listing.
	// The provider decides whether the selected funding method can satisfy the card
	// request.
	PaymentMethodID string `json:"payment_method_id" api:"required"`
	// Wallet item key used to mint this card.
	Wallet    string                               `json:"wallet" api:"required"`
	ExpiresAt param.Opt[int64]                     `json:"expires_at,omitzero"`
	LineItems []CardVaultItemSpecLinkLineItemParam `json:"line_items,omitzero"`
	Metadata  map[string]string                    `json:"metadata,omitzero"`
	Totals    []CardVaultItemSpecLinkTotalParam    `json:"totals,omitzero"`
	// This field can be elided, and will marshal its zero value as "link".
	Provider constant.Link `json:"provider" default:"link"`
	paramObj
}

func (r CardVaultItemSpecLinkParam) MarshalJSON() (data []byte, err error) {
	type shadow CardVaultItemSpecLinkParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CardVaultItemSpecLinkParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Name is required.
type CardVaultItemSpecLinkLineItemParam struct {
	Name        string            `json:"name" api:"required"`
	Description param.Opt[string] `json:"description,omitzero"`
	ImageURL    param.Opt[string] `json:"image_url,omitzero"`
	ProductURL  param.Opt[string] `json:"product_url,omitzero"`
	Quantity    param.Opt[int64]  `json:"quantity,omitzero"`
	SKU         param.Opt[string] `json:"sku,omitzero"`
	// Unit amount in minor currency units.
	UnitAmount param.Opt[int64]                          `json:"unit_amount,omitzero"`
	URL        param.Opt[string]                         `json:"url,omitzero"`
	Totals     []CardVaultItemSpecLinkLineItemTotalParam `json:"totals,omitzero"`
	paramObj
}

func (r CardVaultItemSpecLinkLineItemParam) MarshalJSON() (data []byte, err error) {
	type shadow CardVaultItemSpecLinkLineItemParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CardVaultItemSpecLinkLineItemParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Amount, DisplayText, Type are required.
type CardVaultItemSpecLinkLineItemTotalParam struct {
	// Total amount in minor currency units.
	Amount      int64  `json:"amount" api:"required"`
	DisplayText string `json:"display_text" api:"required"`
	Type        string `json:"type" api:"required"`
	paramObj
}

func (r CardVaultItemSpecLinkLineItemTotalParam) MarshalJSON() (data []byte, err error) {
	type shadow CardVaultItemSpecLinkLineItemTotalParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CardVaultItemSpecLinkLineItemTotalParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Amount, DisplayText, Type are required.
type CardVaultItemSpecLinkTotalParam struct {
	// Total amount in minor currency units.
	Amount      int64  `json:"amount" api:"required"`
	DisplayText string `json:"display_text" api:"required"`
	Type        string `json:"type" api:"required"`
	paramObj
}

func (r CardVaultItemSpecLinkTotalParam) MarshalJSON() (data []byte, err error) {
	type shadow CardVaultItemSpecLinkTotalParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CardVaultItemSpecLinkTotalParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentCard reusable live payment card. Test-mode card creation is not supported.
// Each checkout creates an approval-gated authorization for spec.merchant /
// spec.amount. The card stays ready after each authorization.
//
// The properties Amount, Currency, Merchant, Provider, Wallet are required.
type CardVaultItemSpecAgentcardParam struct {
	// Integer amount in minor currency units.
	Amount   int64  `json:"amount" api:"required"`
	Currency string `json:"currency" api:"required"`
	// Merchant name shown on the cardholder's approval screen.
	Merchant string `json:"merchant" api:"required"`
	// Wallet item key used to authorize checkouts.
	Wallet string `json:"wallet" api:"required"`
	// AgentCard vaulted card to pay with. Omitted, the cardholder picks on the
	// approval screen.
	CardID param.Opt[string] `json:"card_id,omitzero"`
	// This field can be elided, and will marshal its zero value as "agentcard".
	Provider constant.Agentcard `json:"provider" default:"agentcard"`
	paramObj
}

func (r CardVaultItemSpecAgentcardParam) MarshalJSON() (data []byte, err error) {
	type shadow CardVaultItemSpecAgentcardParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CardVaultItemSpecAgentcardParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CardVaultItemStateUnion contains all possible properties and values from
// [CardVaultItemStateLink], [CardVaultItemStateAgentcard].
//
// Use the [CardVaultItemStateUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type CardVaultItemStateUnion struct {
	// Any of "link", "agentcard".
	Provider string `json:"provider"`
	Status   string `json:"status"`
	// This field is from variant [CardVaultItemStateLink].
	Aliases VaultCardAliases `json:"aliases"`
	// This field is from variant [CardVaultItemStateLink].
	Domains []string `json:"domains"`
	// This field is a union of [CardVaultItemStateLinkMasks],
	// [CardVaultItemStateAgentcardMasks]
	Masks        CardVaultItemStateUnionMasks `json:"masks"`
	StatusReason string                       `json:"status_reason"`
	// This field is from variant [CardVaultItemStateAgentcard].
	Authorization AgentcardCheckoutAuthorization `json:"authorization"`
	// This field is from variant [CardVaultItemStateAgentcard].
	Preparation AgentcardCheckoutPreparation `json:"preparation"`
	JSON        struct {
		Provider      respjson.Field
		Status        respjson.Field
		Aliases       respjson.Field
		Domains       respjson.Field
		Masks         respjson.Field
		StatusReason  respjson.Field
		Authorization respjson.Field
		Preparation   respjson.Field
		raw           string
	} `json:"-"`
}

// anyCardVaultItemState is implemented by each variant of
// [CardVaultItemStateUnion] to add type safety for the return type of
// [CardVaultItemStateUnion.AsAny]
type anyCardVaultItemState interface {
	implCardVaultItemStateUnion()
}

func (CardVaultItemStateLink) implCardVaultItemStateUnion()      {}
func (CardVaultItemStateAgentcard) implCardVaultItemStateUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := CardVaultItemStateUnion.AsAny().(type) {
//	case kernel.CardVaultItemStateLink:
//	case kernel.CardVaultItemStateAgentcard:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u CardVaultItemStateUnion) AsAny() anyCardVaultItemState {
	switch u.Provider {
	case "link":
		return u.AsLink()
	case "agentcard":
		return u.AsAgentcard()
	}
	return nil
}

func (u CardVaultItemStateUnion) AsLink() (v CardVaultItemStateLink) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CardVaultItemStateUnion) AsAgentcard() (v CardVaultItemStateAgentcard) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u CardVaultItemStateUnion) RawJSON() string { return u.JSON.raw }

func (r *CardVaultItemStateUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CardVaultItemStateUnionMasks is an implicit subunion of
// [CardVaultItemStateUnion]. CardVaultItemStateUnionMasks provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [CardVaultItemStateUnion].
type CardVaultItemStateUnionMasks struct {
	Brand string `json:"brand"`
	Last4 string `json:"last4"`
	JSON  struct {
		Brand respjson.Field
		Last4 respjson.Field
		raw   string
	} `json:"-"`
}

func (r *CardVaultItemStateUnionMasks) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CardVaultItemStateLink struct {
	Provider constant.Link `json:"provider" default:"link"`
	// recovery_required means an original provider operation has an unresolved
	// outcome. Do not retry, delete, or replace it. Known references may be observed
	// safely, but unknown creation without an ID and uncertain card-material retrieval
	// require manual reconciliation with the provider or support. There is no reset or
	// caller-asserted reconciliation operation.
	//
	// Any of "requested", "pending_authorization", "ready", "consumed", "expired",
	// "declined", "recovery_required".
	Status       string                      `json:"status" api:"required"`
	Aliases      VaultCardAliases            `json:"aliases"`
	Domains      []string                    `json:"domains"`
	Masks        CardVaultItemStateLinkMasks `json:"masks"`
	StatusReason string                      `json:"status_reason"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Provider     respjson.Field
		Status       respjson.Field
		Aliases      respjson.Field
		Domains      respjson.Field
		Masks        respjson.Field
		StatusReason respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CardVaultItemStateLink) RawJSON() string { return r.JSON.raw }
func (r *CardVaultItemStateLink) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CardVaultItemStateLinkMasks struct {
	Brand       string            `json:"brand"`
	Last4       string            `json:"last4"`
	ExtraFields map[string]string `json:"" api:"extrafields"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Brand       respjson.Field
		Last4       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CardVaultItemStateLinkMasks) RawJSON() string { return r.JSON.raw }
func (r *CardVaultItemStateLinkMasks) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CardVaultItemStateAgentcard struct {
	Provider constant.Agentcard `json:"provider" default:"agentcard"`
	// ready_to_submit is device readiness for at most 30 seconds. consumed means the
	// prepared attempt has settled, not that an order succeeded. stopped cannot be
	// reused. outcome_unknown requires merchant reconciliation and blocks new
	// requests. recovery_required means the original checkout outcome is unresolved.
	// Automatic reuse is blocked. Known authorization IDs must be reconciled through
	// provider observations or support. When no authorization ID was returned, an
	// explicitly confirmed item deletion may abandon the unresolved attempt so the
	// caller can create a replacement; deletion does not prove that the original
	// attempt failed. It does not mean declined or expired.
	//
	// Any of "requested", "ready", "preparing", "ready_to_submit", "pending_approval",
	// "consumed", "stopped", "outcome_unknown", "degraded", "recovery_required".
	Status  string           `json:"status" api:"required"`
	Aliases VaultCardAliases `json:"aliases"`
	// The in-flight or most recent checkout authorization. Present while a checkout is
	// pending approval and after it settles.
	Authorization AgentcardCheckoutAuthorization   `json:"authorization"`
	Masks         CardVaultItemStateAgentcardMasks `json:"masks"`
	// One-use processor-bound checkout preparation. Keep the approval page open
	// through token handoff. The amount is display-only and does not constrain the
	// merchant's eventual charge.
	Preparation  AgentcardCheckoutPreparation `json:"preparation"`
	StatusReason string                       `json:"status_reason"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Provider      respjson.Field
		Status        respjson.Field
		Aliases       respjson.Field
		Authorization respjson.Field
		Masks         respjson.Field
		Preparation   respjson.Field
		StatusReason  respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CardVaultItemStateAgentcard) RawJSON() string { return r.JSON.raw }
func (r *CardVaultItemStateAgentcard) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CardVaultItemStateAgentcardMasks struct {
	Brand       string            `json:"brand"`
	Last4       string            `json:"last4"`
	ExtraFields map[string]string `json:"" api:"extrafields"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Brand       respjson.Field
		Last4       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CardVaultItemStateAgentcardMasks) RawJSON() string { return r.JSON.raw }
func (r *CardVaultItemStateAgentcardMasks) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Return the credential item with its collection action. Supported for ready and
// pending_collection credential items. Always render the same form from every
// form-supported field; totp fields have no form input and are omitted. No
// caller-selected field subsets or form overrides are accepted. Reuse an active
// Kernel-hosted session or renew an expired session atomically. Customer-hosted
// forms use their own backend and ordinary item GET/PATCH. Opening the form does
// not clear values or change readiness or item version. To observe edits on a
// ready item, record its version and poll GET without wait until the version
// changes, then reconcile the returned state. Version changes may also come from
// PATCH; they do not identify a particular form submission. Customer-hosted apps
// use their own submission callback, including for unchanged forms. The wait
// parameter waits for readiness, not edits.
//
// The property Type is required.
type CollectVaultItemOperationRequestParam struct {
	// Any of "collect".
	Type CollectVaultItemOperationRequestType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r CollectVaultItemOperationRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow CollectVaultItemOperationRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CollectVaultItemOperationRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CollectVaultItemOperationRequestType string

const (
	CollectVaultItemOperationRequestTypeCollect CollectVaultItemOperationRequestType = "collect"
)

// One schema-derived form for the item, available in ready or pending_collection
// state. Render every form-supported field as editable; omit totp fields and
// preserve their stored seeds. Prefill non-sensitive values, and allow existing
// sensitive values to be preserved or replaced without ever revealing them. No
// field subsets or per-request form configuration exist. Validate required fields
// against the resulting values, including preserved secrets. Submit changed values
// only, using the version used to render the form. Scoped hosted submission
// rejects totp edits; seed writes require the ordinary authenticated item API.
// Customer forms likewise omit totp from their payloads. Save edits atomically. A
// successful hosted submission increments the version, marks ready, and consumes
// the session; an empty edit may complete collection while preserving values. A
// customer form uses PATCH for changed values and does not send an empty PATCH
// when nothing changed. Kernel-hosted bearer sessions require no Kernel account
// and are bound to the item version. Expired, superseded, consumed, or
// deleted-item sessions cannot submit. Authenticated item GET renews expired
// active sessions for ready or pending items; pending items always receive an
// action. A ready item with no active session omits the action until collect is
// invoked. Concurrent renewals return the same link. Renewal changes neither
// values nor item version. An expired link cannot renew itself. The hosted form
// handles its collection protocol; callers only open the returned URL and do not
// extract or submit its token through the public API. For customer-hosted forms,
// use @onkernel/vault-react and an authenticated customer backend calling the
// ordinary item GET/PATCH API. Kernel does not store customer collection URLs or
// authenticate the customer's end users. Treat URLs and submitted values as
// secrets and exclude them from logs, traces, and errors.
type CredentialCollectionAction struct {
	// Expiry of the Kernel-hosted collection link (30 minutes after issuance).
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// Any of "collect".
	Name CredentialCollectionActionName `json:"name" api:"required"`
	// Time-scoped hosted form URL (vault.kernel.sh in production). Open this URL as
	// returned; treat it as a secret.
	URL string `json:"url" api:"required" format:"uri"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExpiresAt   respjson.Field
		Name        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialCollectionAction) RawJSON() string { return r.JSON.raw }
func (r *CredentialCollectionAction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialCollectionActionName string

const (
	CredentialCollectionActionNameCollect CredentialCollectionActionName = "collect"
)

type CredentialVaultFieldDefinition struct {
	// Whether a nonempty value is required for readiness and form submission.
	Required bool `json:"required" api:"required"`
	// Whether the value is omitted from every item response. Reserve true for secrets
	// such as passwords, API tokens, and TOTP seeds. Ordinary usernames and email
	// addresses should be false so the form can display and prefill them.
	Sensitive bool `json:"sensitive" api:"required"`
	// Text, email, and password have form inputs; totp does not and is omitted from
	// both Kernel-hosted and customer React forms. Password and totp must be
	// sensitive. A totp value is an RFC 4648 Base32 generator seed (case-insensitive,
	// optional trailing padding), not an otpauth URI or current code. Reject invalid
	// or empty decoded seeds. Browser fill generates an RFC 6238 code at execution
	// time using HMAC-SHA1, 6 digits, and a 30-second period. Preserve leading zeros;
	// never fill the seed. Custom algorithms, digits, periods, and form enrollment are
	// unsupported.
	//
	// Any of "text", "email", "password", "totp".
	Type CredentialVaultFieldType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Required    respjson.Field
		Sensitive   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialVaultFieldDefinition) RawJSON() string { return r.JSON.raw }
func (r *CredentialVaultFieldDefinition) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Type is required.
type CredentialVaultFieldInputParam struct {
	// Text, email, and password have form inputs; totp does not and is omitted from
	// both Kernel-hosted and customer React forms. Password and totp must be
	// sensitive. A totp value is an RFC 4648 Base32 generator seed (case-insensitive,
	// optional trailing padding), not an otpauth URI or current code. Reject invalid
	// or empty decoded seeds. Browser fill generates an RFC 6238 code at execution
	// time using HMAC-SHA1, 6 digits, and a 30-second period. Preserve leading zeros;
	// never fill the seed. Custom algorithms, digits, periods, and form enrollment are
	// unsupported.
	//
	// Any of "text", "email", "password", "totp".
	Type     CredentialVaultFieldType `json:"type,omitzero" api:"required"`
	Required param.Opt[bool]          `json:"required,omitzero"`
	// Set false explicitly for ordinary usernames, email addresses, and other
	// non-secret identifiers. Reserve true for secrets such as passwords, API tokens,
	// and TOTP seeds. Password and totp fields must be true. Omission defaults to true
	// for safety; do not rely on that default for every field. False permits API reads
	// and form prefilling.
	Sensitive param.Opt[bool] `json:"sensitive,omitzero"`
	// Optional initial value satisfying the declared type, at most 16 KiB in UTF-8
	// bytes. Omit to leave unset; null and empty strings are rejected on creation.
	// Sensitive values are encrypted and never copied into the returned spec.
	Value param.Opt[string] `json:"value,omitzero"`
	paramObj
}

func (r CredentialVaultFieldInputParam) MarshalJSON() (data []byte, err error) {
	type shadow CredentialVaultFieldInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialVaultFieldInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialVaultFieldState struct {
	HasValue bool `json:"has_value" api:"required"`
	// Present exactly when has_value is true and the field is not sensitive. Reflects
	// the latest developer or human edit. For totp, has_value indicates a stored seed;
	// neither the seed nor a generated code is returned.
	Value string `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		HasValue    respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialVaultFieldState) RawJSON() string { return r.JSON.raw }
func (r *CredentialVaultFieldState) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Text, email, and password have form inputs; totp does not and is omitted from
// both Kernel-hosted and customer React forms. Password and totp must be
// sensitive. A totp value is an RFC 4648 Base32 generator seed (case-insensitive,
// optional trailing padding), not an otpauth URI or current code. Reject invalid
// or empty decoded seeds. Browser fill generates an RFC 6238 code at execution
// time using HMAC-SHA1, 6 digits, and a 30-second period. Preserve leading zeros;
// never fill the seed. Custom algorithms, digits, periods, and form enrollment are
// unsupported.
type CredentialVaultFieldType string

const (
	CredentialVaultFieldTypeText     CredentialVaultFieldType = "text"
	CredentialVaultFieldTypeEmail    CredentialVaultFieldType = "email"
	CredentialVaultFieldTypePassword CredentialVaultFieldType = "password"
	CredentialVaultFieldTypeTotp     CredentialVaultFieldType = "totp"
)

// The property Value is required.
type CredentialVaultFieldUpdateParam struct {
	// Replacement value (at most 16 KiB in UTF-8 bytes), or null or an empty string to
	// immediately clear the stored value. Clearing a required form-supported field
	// reopens collection; clearing an optional field does not prevent readiness.
	// Values must satisfy the declared field type. For totp, value is the generator
	// seed, never a current code. Clearing a required totp field returns 400 because
	// it cannot be collected in a form.
	Value param.Opt[string] `json:"value,omitzero" api:"required"`
	paramObj
}

func (r CredentialVaultFieldUpdateParam) MarshalJSON() (data []byte, err error) {
	type shadow CredentialVaultFieldUpdateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialVaultFieldUpdateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialVaultItem struct {
	ID                  string                                  `json:"id" api:"required"`
	AvailableExpansions []CredentialVaultItemAvailableExpansion `json:"available_expansions" api:"required"`
	// Advertises collect for ready and pending_collection items. Browser fill is
	// advertised only when separately implemented and eligible.
	AvailableOperations []CredentialVaultItemAvailableOperation `json:"available_operations" api:"required"`
	CreatedAt           time.Time                               `json:"created_at" api:"required" format:"date-time"`
	// Immutable item key assigned when the item is created.
	Key   string                   `json:"key" api:"required"`
	Spec  CredentialVaultItemSpec  `json:"spec" api:"required"`
	State CredentialVaultItemState `json:"state" api:"required"`
	// Any of "credential".
	Type      CredentialVaultItemType `json:"type" api:"required"`
	UpdatedAt time.Time               `json:"updated_at" api:"required" format:"date-time"`
	// Starts at 1 and increments on PATCH and successful hosted submission, but not
	// collection-link renewal.
	Version int64 `json:"version" api:"required"`
	// One schema-derived form for the item, available in ready or pending_collection
	// state. Render every form-supported field as editable; omit totp fields and
	// preserve their stored seeds. Prefill non-sensitive values, and allow existing
	// sensitive values to be preserved or replaced without ever revealing them. No
	// field subsets or per-request form configuration exist. Validate required fields
	// against the resulting values, including preserved secrets. Submit changed values
	// only, using the version used to render the form. Scoped hosted submission
	// rejects totp edits; seed writes require the ordinary authenticated item API.
	// Customer forms likewise omit totp from their payloads. Save edits atomically. A
	// successful hosted submission increments the version, marks ready, and consumes
	// the session; an empty edit may complete collection while preserving values. A
	// customer form uses PATCH for changed values and does not send an empty PATCH
	// when nothing changed. Kernel-hosted bearer sessions require no Kernel account
	// and are bound to the item version. Expired, superseded, consumed, or
	// deleted-item sessions cannot submit. Authenticated item GET renews expired
	// active sessions for ready or pending items; pending items always receive an
	// action. A ready item with no active session omits the action until collect is
	// invoked. Concurrent renewals return the same link. Renewal changes neither
	// values nor item version. An expired link cannot renew itself. The hosted form
	// handles its collection protocol; callers only open the returned URL and do not
	// extract or submit its token through the public API. For customer-hosted forms,
	// use @onkernel/vault-react and an authenticated customer backend calling the
	// ordinary item GET/PATCH API. Kernel does not store customer collection URLs or
	// authenticate the customer's end users. Treat URLs and submitted values as
	// secrets and exclude them from logs, traces, and errors.
	Action CredentialCollectionAction `json:"action"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		AvailableExpansions respjson.Field
		AvailableOperations respjson.Field
		CreatedAt           respjson.Field
		Key                 respjson.Field
		Spec                respjson.Field
		State               respjson.Field
		Type                respjson.Field
		UpdatedAt           respjson.Field
		Version             respjson.Field
		Action              respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialVaultItem) RawJSON() string { return r.JSON.raw }
func (r *CredentialVaultItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Live data that can currently be requested by passing its type to the item GET
// expand parameter.
type CredentialVaultItemAvailableExpansion struct {
	Description string `json:"description" api:"required"`
	// Any of "payment_methods".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialVaultItemAvailableExpansion) RawJSON() string { return r.JSON.raw }
func (r *CredentialVaultItemAvailableExpansion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An operation that is currently valid for this item. Read the description before
// invoking it through the item operations endpoint.
type CredentialVaultItemAvailableOperation struct {
	Description string `json:"description" api:"required"`
	// Any of "authorize", "collect", "prepare_checkout", "fill".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialVaultItemAvailableOperation) RawJSON() string { return r.JSON.raw }
func (r *CredentialVaultItemAvailableOperation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialVaultItemType string

const (
	CredentialVaultItemTypeCredential CredentialVaultItemType = "credential"
)

// Create a credential item without a wallet or external provider. Do not use
// credential items to store, collect, or fill credit card data, including card
// numbers (PANs), security codes (CVV/CVC), or expiration dates. Use wallet and
// card item types for credit cards and payment checkout instead. If all required
// fields have values, return ready without a collection action; collect can still
// open its form. Otherwise return pending_collection with a time-scoped
// Kernel-hosted collection action. Missing optional fields alone do not trigger
// collection. Repeating the original creation request returns the current item
// without overwriting later edits; a different request at the same key
// returns 409. Use PATCH for updates. Required totp fields must include a valid
// seed on creation; otherwise return 400 rather than opening a form that cannot
// collect it. Optional totp fields may be unset and populated later through PATCH.
//
// The properties Spec, Type are required.
type CredentialVaultItemRequestParam struct {
	// Credential fields are for login and other non-payment credentials. Do not store,
	// collect, or fill credit card data in credential items. Use wallet and card item
	// types for credit cards and payment checkout instead.
	Spec CredentialVaultItemSpecInputParam `json:"spec,omitzero" api:"required"`
	// Any of "credential".
	Type CredentialVaultItemRequestType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r CredentialVaultItemRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow CredentialVaultItemRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialVaultItemRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialVaultItemRequestType string

const (
	CredentialVaultItemRequestTypeCredential CredentialVaultItemRequestType = "credential"
)

type CredentialVaultItemSpec struct {
	Fields map[string]CredentialVaultFieldDefinition `json:"fields" api:"required"`
	// Recognizable site or service name displayed verbatim as the form title, without
	// suffixes such as sign-in credentials. Display text only, not an enforced
	// destination policy.
	Description string `json:"description"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Fields      respjson.Field
		Description respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialVaultItemSpec) RawJSON() string { return r.JSON.raw }
func (r *CredentialVaultItemSpec) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Credential fields are for login and other non-payment credentials. Do not store,
// collect, or fill credit card data in credential items. Use wallet and card item
// types for credit cards and payment checkout instead.
//
// The property Fields is required.
type CredentialVaultItemSpecInputParam struct {
	Fields map[string]CredentialVaultFieldInputParam `json:"fields,omitzero" api:"required"`
	// The site's recognizable display name, used verbatim as the user-facing form
	// title (for example, Hacker News). Use only the site or service name; do not
	// append sign-in, login, credentials, or task instructions. This is display text,
	// not an enforced destination policy. At most 16 KiB in UTF-8 bytes.
	Description param.Opt[string] `json:"description,omitzero"`
	paramObj
}

func (r CredentialVaultItemSpecInputParam) MarshalJSON() (data []byte, err error) {
	type shadow CredentialVaultItemSpecInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialVaultItemSpecInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialVaultItemSpecUpdateParam struct {
	// Recognizable site or service name used as the form title, without suffixes such
	// as sign-in credentials. An empty string clears it. Display text only, not an
	// enforced destination policy. The server also enforces a 16 KiB UTF-8 byte limit.
	Description param.Opt[string]                          `json:"description,omitzero"`
	Fields      map[string]CredentialVaultFieldUpdateParam `json:"fields,omitzero"`
	paramObj
}

func (r CredentialVaultItemSpecUpdateParam) MarshalJSON() (data []byte, err error) {
	type shadow CredentialVaultItemSpecUpdateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialVaultItemSpecUpdateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialVaultItemState struct {
	// Exactly one entry for each declared field.
	Fields map[string]CredentialVaultFieldState `json:"fields" api:"required"`
	// Ready means all required fields have values, not that a login succeeded.
	// Optional fields may remain unset.
	//
	// Any of "pending_collection", "ready".
	Status CredentialVaultItemStateStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Fields      respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialVaultItemState) RawJSON() string { return r.JSON.raw }
func (r *CredentialVaultItemState) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Ready means all required fields have values, not that a login succeeded.
// Optional fields may remain unset.
type CredentialVaultItemStateStatus string

const (
	CredentialVaultItemStateStatusPendingCollection CredentialVaultItemStateStatus = "pending_collection"
	CredentialVaultItemStateStatusReady             CredentialVaultItemStateStatus = "ready"
)

// Atomically update description and selected values. Omitted properties are
// preserved. Field names, types, required flags, and sensitivity cannot change.
// Unknown field names return 400; stale versions or mismatched item types return
// 409 without changing the item. A successful update increments version and
// invalidates outstanding Kernel-hosted collection sessions. If required values
// remain missing, return pending_collection and a fresh collection action.
// Otherwise return ready without an action; collect can open the form again
// without clearing values. Customer URLs have no Kernel-managed expiry.
//
// The properties Spec, Type, Version are required.
type CredentialVaultItemUpdateRequestParam struct {
	Spec CredentialVaultItemSpecUpdateParam `json:"spec,omitzero" api:"required"`
	// Any of "credential".
	Type CredentialVaultItemUpdateRequestType `json:"type,omitzero" api:"required"`
	// Expected current item version from the latest read.
	Version int64 `json:"version" api:"required"`
	// Optional immutable item ID precondition. Returns 409 if the key now identifies a
	// different item. Accepted writes target this immutable ID, preventing
	// replacement-key races. Supply this when submitting a form bound to a previously
	// read item.
	ExpectedItemID param.Opt[string] `json:"expected_item_id,omitzero"`
	paramObj
}

func (r CredentialVaultItemUpdateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow CredentialVaultItemUpdateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialVaultItemUpdateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialVaultItemUpdateRequestType string

const (
	CredentialVaultItemUpdateRequestTypeCredential CredentialVaultItemUpdateRequestType = "credential"
)

// Fill selected fields from one ready credential or ready, unexpired Link card
// into a browser linked to its vault. Only invoke when the item advertises `fill`.
// Browser and vault must belong to the same project. Kernel checks access and
// allowed destinations before filling; providing a page URL does not authorize a
// destination.
//
// Find exactly one open page matching `page_url`. Credential items may omit
// `page_url` to require exactly one open page; cards require an HTTPS page URL.
// Credentials have no destination allowlist. TOTP fields generate a current code
// immediately before writing; their seeds never enter the browser. For each
// selector, search the main frame and all descendant frames for editable inputs or
// selects matched directly or contained within matching elements. Each selector
// must resolve to one unique editable element across all frames; zero or multiple
// candidates fail. Count each element once, even if multiple matching containers
// contain it. Validate all bindings before filling. Select elements match an
// option by its value, not its label. If the page navigates or a target disappears
// during filling, stop rather than selecting a different page or element.
//
// Fill in request order and stop on the first failure. This operation is not
// atomic: previously filled fields are not rolled back. Never submit the form or
// click buttons, though input/change events may trigger site behavior. Fill is the
// preferred browser-checkout path. Aliases remain an alternative for explicitly
// chosen egress-substitution integrations. Do not automatically retry or fall back
// to aliases after a failed or indeterminate operation.
//
// Secret values are never returned or included in operation logs, traces, audit
// events, or error details. This does not prevent an agent with unrestricted
// browser access from reading values from the page or other browser observation
// surfaces.
//
// The properties BrowserID, Fields, Type are required.
type FillVaultItemOperationRequestParam struct {
	// Browser session ID, not a reusable browser name.
	BrowserID string `json:"browser_id" api:"required"`
	// Field bindings for this step. No two bindings may resolve to the same element.
	Fields []VaultFillFieldParam `json:"fields,omitzero" api:"required"`
	// Any of "fill".
	Type FillVaultItemOperationRequestType `json:"type,omitzero" api:"required"`
	// Exact current top-level page URL, including path, query, and fragment. Must
	// match exactly one open page in the browser; zero or multiple matches fail. No
	// prefix or glob matching. Required for cards, which must use HTTPS without
	// embedded credentials. Optional for credentials, where omission requires exactly
	// one open page.
	PageURL param.Opt[string] `json:"page_url,omitzero" format:"uri"`
	// Total operation deadline in milliseconds, not a per-field timeout.
	TimeoutMs param.Opt[int64] `json:"timeout_ms,omitzero"`
	paramObj
}

func (r FillVaultItemOperationRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow FillVaultItemOperationRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *FillVaultItemOperationRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FillVaultItemOperationRequestType string

const (
	FillVaultItemOperationRequestTypeFill FillVaultItemOperationRequestType = "fill"
)

type FillVaultItemOperationResult struct {
	// Exactly one result per request binding, in request order. After the first failed
	// or unknown field, all remaining fields are not_attempted.
	Fields []VaultFillFieldResult `json:"fields" api:"required"`
	// Completed only when all fields were filled. Failed when execution stopped with
	// known outcomes. Unknown when any field's outcome cannot be determined. None of
	// these statuses confirms payment or merchant acceptance.
	//
	// Any of "completed", "failed", "unknown".
	Status FillVaultItemOperationResultStatus `json:"status" api:"required"`
	// Any of "fill".
	Type FillVaultItemOperationResultType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Fields      respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FillVaultItemOperationResult) RawJSON() string { return r.JSON.raw }
func (r *FillVaultItemOperationResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Completed only when all fields were filled. Failed when execution stopped with
// known outcomes. Unknown when any field's outcome cannot be determined. None of
// these statuses confirms payment or merchant acceptance.
type FillVaultItemOperationResultStatus string

const (
	FillVaultItemOperationResultStatusCompleted FillVaultItemOperationResultStatus = "completed"
	FillVaultItemOperationResultStatusFailed    FillVaultItemOperationResultStatus = "failed"
	FillVaultItemOperationResultStatusUnknown   FillVaultItemOperationResultStatus = "unknown"
)

type FillVaultItemOperationResultType string

const (
	FillVaultItemOperationResultTypeFill FillVaultItemOperationResultType = "fill"
)

// Prepare an unused AgentCard card for a supported tokenization checkout. Deliver
// the returned approval URL and keep the approval page open. Poll the item until
// ready_to_submit, then submit native Pay before preparation.expires_at. Readiness
// lasts at most 30 seconds. Unused preparations expire automatically. Preparations
// are single-use even after failure or expiry; do not automatically retry and
// reconcile uncertain outcomes with the merchant.
//
// The properties Checkout, Type are required.
type PrepareCheckoutVaultItemOperationRequestParam struct {
	// Required when preparing an unused AgentCard card for a supported tokenization
	// processor. Consent is bound to this browser and declared merchant origin, not a
	// tab. Wait for the item's ready_to_submit status before native Pay and submit
	// within its readiness deadline. Unused preparations expire automatically; every
	// preparation is single-use, including after failure or expiry.
	Checkout VaultCheckoutContextParam `json:"checkout,omitzero" api:"required"`
	// Any of "prepare_checkout".
	Type PrepareCheckoutVaultItemOperationRequestType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r PrepareCheckoutVaultItemOperationRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow PrepareCheckoutVaultItemOperationRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PrepareCheckoutVaultItemOperationRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type PrepareCheckoutVaultItemOperationRequestType string

const (
	PrepareCheckoutVaultItemOperationRequestTypePrepareCheckout PrepareCheckoutVaultItemOperationRequestType = "prepare_checkout"
)

type VaultCardAliases struct {
	Cvc      string `json:"cvc" api:"required"`
	ExpMonth string `json:"exp_month" api:"required"`
	ExpYear  string `json:"exp_year" api:"required"`
	Number   string `json:"number" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Cvc         respjson.Field
		ExpMonth    respjson.Field
		ExpYear     respjson.Field
		Number      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultCardAliases) RawJSON() string { return r.JSON.raw }
func (r *VaultCardAliases) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Required when preparing an unused AgentCard card for a supported tokenization
// processor. Consent is bound to this browser and declared merchant origin, not a
// tab. Wait for the item's ready_to_submit status before native Pay and submit
// within its readiness deadline. Unused preparations expire automatically; every
// preparation is single-use, including after failure or expiry.
//
// The properties BrowserID, Environment, MerchantOrigin are required.
type VaultCheckoutContextParam struct {
	// Active browser session with this vault bound to it.
	BrowserID string `json:"browser_id" api:"required"`
	// Use production or sandbox for Square, Braintree and Worldpay; shared for Bambora
	// and Mercado Pago. Shared endpoints do not establish test mode. Merchant
	// credentials/configuration determine processor test mode, independently of the
	// AgentCard credential mode.
	//
	// Any of "production", "sandbox", "shared".
	Environment VaultCheckoutContextEnvironment `json:"environment,omitzero" api:"required"`
	// Canonical HTTPS origin of the top-level merchant document, not a processor
	// iframe. HTTP localhost is accepted for tests.
	MerchantOrigin string `json:"merchant_origin" api:"required"`
	// Tokenization processor. Omit for Square compatibility. Non-Square processors
	// require multi-processor preparation enablement.
	//
	// Any of "square", "braintree", "worldpay", "bambora", "mercado_pago".
	Psp AgentcardPreparedProcessor `json:"psp,omitzero"`
	paramObj
}

func (r VaultCheckoutContextParam) MarshalJSON() (data []byte, err error) {
	type shadow VaultCheckoutContextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultCheckoutContextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Use production or sandbox for Square, Braintree and Worldpay; shared for Bambora
// and Mercado Pago. Shared endpoints do not establish test mode. Merchant
// credentials/configuration determine processor test mode, independently of the
// AgentCard credential mode.
type VaultCheckoutContextEnvironment string

const (
	VaultCheckoutContextEnvironmentProduction VaultCheckoutContextEnvironment = "production"
	VaultCheckoutContextEnvironmentSandbox    VaultCheckoutContextEnvironment = "sandbox"
	VaultCheckoutContextEnvironmentShared     VaultCheckoutContextEnvironment = "shared"
)

// The properties Field, Selector are required.
type VaultFillFieldParam struct {
	// A declared credential field name or a supported card field. Unset credential
	// fields cannot be filled.
	Field    string `json:"field" api:"required"`
	Selector string `json:"selector" api:"required"`
	// Required only for a card's combined expiration field. Forbidden for other card
	// fields and all credential fields.
	//
	// Any of "MM/YY", "MM/YYYY".
	Format VaultFillFieldFormat `json:"format,omitzero"`
	paramObj
}

func (r VaultFillFieldParam) MarshalJSON() (data []byte, err error) {
	type shadow VaultFillFieldParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultFillFieldParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Required only for a card's combined expiration field. Forbidden for other card
// fields and all credential fields.
type VaultFillFieldFormat string

const (
	VaultFillFieldFormatMmYy   VaultFillFieldFormat = "MM/YY"
	VaultFillFieldFormatMmYyyy VaultFillFieldFormat = "MM/YYYY"
)

type VaultFillFieldResult struct {
	// Zero-based index into the request fields array.
	Index int64 `json:"index" api:"required"`
	// Filled means the fill action completed, not that the website retained or
	// accepted the value.
	//
	// Any of "filled", "failed", "not_attempted", "unknown".
	Status VaultFillFieldResultStatus `json:"status" api:"required"`
	// Present only for failed or unknown fields. Never includes secret values, DOM
	// content, or raw browser errors.
	//
	// Any of "target_changed", "element_not_found", "ambiguous_selector",
	// "element_not_editable", "option_not_found", "timeout", "execution_failed".
	ErrorCode VaultFillFieldResultErrorCode `json:"error_code"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Index       respjson.Field
		Status      respjson.Field
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultFillFieldResult) RawJSON() string { return r.JSON.raw }
func (r *VaultFillFieldResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Filled means the fill action completed, not that the website retained or
// accepted the value.
type VaultFillFieldResultStatus string

const (
	VaultFillFieldResultStatusFilled       VaultFillFieldResultStatus = "filled"
	VaultFillFieldResultStatusFailed       VaultFillFieldResultStatus = "failed"
	VaultFillFieldResultStatusNotAttempted VaultFillFieldResultStatus = "not_attempted"
	VaultFillFieldResultStatusUnknown      VaultFillFieldResultStatus = "unknown"
)

// Present only for failed or unknown fields. Never includes secret values, DOM
// content, or raw browser errors.
type VaultFillFieldResultErrorCode string

const (
	VaultFillFieldResultErrorCodeTargetChanged      VaultFillFieldResultErrorCode = "target_changed"
	VaultFillFieldResultErrorCodeElementNotFound    VaultFillFieldResultErrorCode = "element_not_found"
	VaultFillFieldResultErrorCodeAmbiguousSelector  VaultFillFieldResultErrorCode = "ambiguous_selector"
	VaultFillFieldResultErrorCodeElementNotEditable VaultFillFieldResultErrorCode = "element_not_editable"
	VaultFillFieldResultErrorCodeOptionNotFound     VaultFillFieldResultErrorCode = "option_not_found"
	VaultFillFieldResultErrorCodeTimeout            VaultFillFieldResultErrorCode = "timeout"
	VaultFillFieldResultErrorCodeExecutionFailed    VaultFillFieldResultErrorCode = "execution_failed"
)

// VaultItemUnion contains all possible properties and values from
// [VaultItemWallet], [VaultItemCard], [CredentialVaultItem].
//
// Use the [VaultItemUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type VaultItemUnion struct {
	ID string `json:"id"`
	// This field is a union of [[]VaultItemWalletAvailableExpansion],
	// [[]VaultItemCardAvailableExpansion], [[]CredentialVaultItemAvailableExpansion]
	AvailableExpansions VaultItemUnionAvailableExpansions `json:"available_expansions"`
	// This field is a union of [[]VaultItemWalletAvailableOperation],
	// [[]VaultItemCardAvailableOperation], [[]CredentialVaultItemAvailableOperation]
	AvailableOperations VaultItemUnionAvailableOperations `json:"available_operations"`
	CreatedAt           time.Time                         `json:"created_at"`
	Key                 string                            `json:"key"`
	// This field is a union of [WalletVaultItemSpecUnion], [CardVaultItemSpecUnion],
	// [CredentialVaultItemSpec]
	Spec VaultItemUnionSpec `json:"spec"`
	// This field is a union of [WalletVaultItemStateUnion], [CardVaultItemStateUnion],
	// [CredentialVaultItemState]
	State VaultItemUnionState `json:"state"`
	// Any of "wallet", "card", "credential".
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	// This field is a union of [VaultItemActionUnion], [CredentialCollectionAction]
	Action VaultItemUnionAction `json:"action"`
	// This field is from variant [VaultItemWallet].
	Expanded  VaultItemWalletExpanded `json:"expanded"`
	ExpiresAt time.Time               `json:"expires_at"`
	// This field is from variant [CredentialVaultItem].
	Version int64 `json:"version"`
	JSON    struct {
		ID                  respjson.Field
		AvailableExpansions respjson.Field
		AvailableOperations respjson.Field
		CreatedAt           respjson.Field
		Key                 respjson.Field
		Spec                respjson.Field
		State               respjson.Field
		Type                respjson.Field
		UpdatedAt           respjson.Field
		Action              respjson.Field
		Expanded            respjson.Field
		ExpiresAt           respjson.Field
		Version             respjson.Field
		raw                 string
	} `json:"-"`
}

// anyVaultItem is implemented by each variant of [VaultItemUnion] to add type
// safety for the return type of [VaultItemUnion.AsAny]
type anyVaultItem interface {
	implVaultItemUnion()
}

func (VaultItemWallet) implVaultItemUnion()     {}
func (VaultItemCard) implVaultItemUnion()       {}
func (CredentialVaultItem) implVaultItemUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := VaultItemUnion.AsAny().(type) {
//	case kernel.VaultItemWallet:
//	case kernel.VaultItemCard:
//	case kernel.CredentialVaultItem:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u VaultItemUnion) AsAny() anyVaultItem {
	switch u.Type {
	case "wallet":
		return u.AsWallet()
	case "card":
		return u.AsCard()
	case "credential":
		return u.AsCredential()
	}
	return nil
}

func (u VaultItemUnion) AsWallet() (v VaultItemWallet) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultItemUnion) AsCard() (v VaultItemCard) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultItemUnion) AsCredential() (v CredentialVaultItem) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u VaultItemUnion) RawJSON() string { return u.JSON.raw }

func (r *VaultItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemUnionAvailableExpansions is an implicit subunion of [VaultItemUnion].
// VaultItemUnionAvailableExpansions provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfVaultItemWalletAvailableExpansions
// OfVaultItemCardAvailableExpansions OfCredentialVaultItemAvailableExpansions]
type VaultItemUnionAvailableExpansions struct {
	// This field will be present if the value is a
	// [[]VaultItemWalletAvailableExpansion] instead of an object.
	OfVaultItemWalletAvailableExpansions []VaultItemWalletAvailableExpansion `json:",inline"`
	// This field will be present if the value is a [[]VaultItemCardAvailableExpansion]
	// instead of an object.
	OfVaultItemCardAvailableExpansions []VaultItemCardAvailableExpansion `json:",inline"`
	// This field will be present if the value is a
	// [[]CredentialVaultItemAvailableExpansion] instead of an object.
	OfCredentialVaultItemAvailableExpansions []CredentialVaultItemAvailableExpansion `json:",inline"`
	JSON                                     struct {
		OfVaultItemWalletAvailableExpansions     respjson.Field
		OfVaultItemCardAvailableExpansions       respjson.Field
		OfCredentialVaultItemAvailableExpansions respjson.Field
		raw                                      string
	} `json:"-"`
}

func (r *VaultItemUnionAvailableExpansions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemUnionAvailableOperations is an implicit subunion of [VaultItemUnion].
// VaultItemUnionAvailableOperations provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfVaultItemWalletAvailableOperations
// OfVaultItemCardAvailableOperations OfCredentialVaultItemAvailableOperations]
type VaultItemUnionAvailableOperations struct {
	// This field will be present if the value is a
	// [[]VaultItemWalletAvailableOperation] instead of an object.
	OfVaultItemWalletAvailableOperations []VaultItemWalletAvailableOperation `json:",inline"`
	// This field will be present if the value is a [[]VaultItemCardAvailableOperation]
	// instead of an object.
	OfVaultItemCardAvailableOperations []VaultItemCardAvailableOperation `json:",inline"`
	// This field will be present if the value is a
	// [[]CredentialVaultItemAvailableOperation] instead of an object.
	OfCredentialVaultItemAvailableOperations []CredentialVaultItemAvailableOperation `json:",inline"`
	JSON                                     struct {
		OfVaultItemWalletAvailableOperations     respjson.Field
		OfVaultItemCardAvailableOperations       respjson.Field
		OfCredentialVaultItemAvailableOperations respjson.Field
		raw                                      string
	} `json:"-"`
}

func (r *VaultItemUnionAvailableOperations) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemUnionSpec is an implicit subunion of [VaultItemUnion].
// VaultItemUnionSpec provides convenient access to the sub-properties of the
// union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemUnion].
type VaultItemUnionSpec struct {
	// This field is from variant [WalletVaultItemSpecUnion].
	Authorization WalletVaultItemSpecLinkAuthorization `json:"authorization"`
	Provider      string                               `json:"provider"`
	// This field is from variant [WalletVaultItemSpecUnion].
	ProviderConfig WalletVaultItemSpecAgentcardProviderConfig `json:"provider_config"`
	// This field is from variant [WalletVaultItemSpecUnion].
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
	// This field is from variant [CardVaultItemSpecUnion].
	Context  string `json:"context"`
	Currency string `json:"currency"`
	// This field is from variant [CardVaultItemSpecUnion].
	MerchantName string `json:"merchant_name"`
	// This field is from variant [CardVaultItemSpecUnion].
	MerchantURL string `json:"merchant_url"`
	// This field is from variant [CardVaultItemSpecUnion].
	PaymentMethodID string `json:"payment_method_id"`
	Wallet          string `json:"wallet"`
	// This field is from variant [CardVaultItemSpecUnion].
	ExpiresAt int64 `json:"expires_at"`
	// This field is from variant [CardVaultItemSpecUnion].
	LineItems []CardVaultItemSpecLinkLineItem `json:"line_items"`
	// This field is from variant [CardVaultItemSpecUnion].
	Metadata map[string]string `json:"metadata"`
	// This field is from variant [CardVaultItemSpecUnion].
	Totals []CardVaultItemSpecLinkTotal `json:"totals"`
	// This field is from variant [CardVaultItemSpecUnion].
	Merchant string `json:"merchant"`
	// This field is from variant [CardVaultItemSpecUnion].
	CardID string `json:"card_id"`
	// This field is from variant [CredentialVaultItemSpec].
	Fields map[string]CredentialVaultFieldDefinition `json:"fields"`
	// This field is from variant [CredentialVaultItemSpec].
	Description string `json:"description"`
	JSON        struct {
		Authorization   respjson.Field
		Provider        respjson.Field
		ProviderConfig  respjson.Field
		UserID          respjson.Field
		Amount          respjson.Field
		Context         respjson.Field
		Currency        respjson.Field
		MerchantName    respjson.Field
		MerchantURL     respjson.Field
		PaymentMethodID respjson.Field
		Wallet          respjson.Field
		ExpiresAt       respjson.Field
		LineItems       respjson.Field
		Metadata        respjson.Field
		Totals          respjson.Field
		Merchant        respjson.Field
		CardID          respjson.Field
		Fields          respjson.Field
		Description     respjson.Field
		raw             string
	} `json:"-"`
}

func (r *VaultItemUnionSpec) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemUnionState is an implicit subunion of [VaultItemUnion].
// VaultItemUnionState provides convenient access to the sub-properties of the
// union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemUnion].
type VaultItemUnionState struct {
	Provider     string `json:"provider"`
	Status       string `json:"status"`
	StatusReason string `json:"status_reason"`
	// This field is from variant [WalletVaultItemStateUnion].
	UserID string `json:"user_id"`
	// This field is from variant [CardVaultItemStateUnion].
	Aliases VaultCardAliases `json:"aliases"`
	// This field is from variant [CardVaultItemStateUnion].
	Domains []string `json:"domains"`
	// This field is a union of [CardVaultItemStateLinkMasks],
	// [CardVaultItemStateAgentcardMasks]
	Masks VaultItemUnionStateMasks `json:"masks"`
	// This field is from variant [CardVaultItemStateUnion].
	Authorization AgentcardCheckoutAuthorization `json:"authorization"`
	// This field is from variant [CardVaultItemStateUnion].
	Preparation AgentcardCheckoutPreparation `json:"preparation"`
	// This field is from variant [CredentialVaultItemState].
	Fields map[string]CredentialVaultFieldState `json:"fields"`
	JSON   struct {
		Provider      respjson.Field
		Status        respjson.Field
		StatusReason  respjson.Field
		UserID        respjson.Field
		Aliases       respjson.Field
		Domains       respjson.Field
		Masks         respjson.Field
		Authorization respjson.Field
		Preparation   respjson.Field
		Fields        respjson.Field
		raw           string
	} `json:"-"`
}

func (r *VaultItemUnionState) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemUnionStateMasks is an implicit subunion of [VaultItemUnion].
// VaultItemUnionStateMasks provides convenient access to the sub-properties of the
// union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemUnion].
type VaultItemUnionStateMasks struct {
	Brand string `json:"brand"`
	Last4 string `json:"last4"`
	JSON  struct {
		Brand respjson.Field
		Last4 respjson.Field
		raw   string
	} `json:"-"`
}

func (r *VaultItemUnionStateMasks) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemUnionAction is an implicit subunion of [VaultItemUnion].
// VaultItemUnionAction provides convenient access to the sub-properties of the
// union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemUnion].
type VaultItemUnionAction struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	// This field is from variant [CredentialCollectionAction].
	ExpiresAt time.Time `json:"expires_at"`
	JSON      struct {
		Name      respjson.Field
		URL       respjson.Field
		ExpiresAt respjson.Field
		raw       string
	} `json:"-"`
}

func (r *VaultItemUnionAction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemWallet struct {
	ID                  string                              `json:"id" api:"required"`
	AvailableExpansions []VaultItemWalletAvailableExpansion `json:"available_expansions" api:"required"`
	AvailableOperations []VaultItemWalletAvailableOperation `json:"available_operations" api:"required"`
	CreatedAt           time.Time                           `json:"created_at" api:"required" format:"date-time"`
	// Immutable item key assigned when the item is created.
	Key string `json:"key" api:"required"`
	// AgentCard wallet. Omit provider_config to use Kernel-managed credentials, or
	// select a customer-owned configuration. Mode (sandbox vs live) is determined by
	// the selected credential; there is no per-item test flag. Without user_id,
	// creation returns a hosted enrollment action and Kernel polls until the user
	// connects. user_id may only reference a user already enrolled by a wallet in this
	// organization under the same configuration.
	Spec      WalletVaultItemSpecUnion  `json:"spec" api:"required"`
	State     WalletVaultItemStateUnion `json:"state" api:"required"`
	Type      constant.Wallet           `json:"type" default:"wallet"`
	UpdatedAt time.Time                 `json:"updated_at" api:"required" format:"date-time"`
	Action    VaultItemActionUnion      `json:"action"`
	// Live, non-persisted data requested through the item GET expand parameter.
	Expanded  VaultItemWalletExpanded `json:"expanded"`
	ExpiresAt time.Time               `json:"expires_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		AvailableExpansions respjson.Field
		AvailableOperations respjson.Field
		CreatedAt           respjson.Field
		Key                 respjson.Field
		Spec                respjson.Field
		State               respjson.Field
		Type                respjson.Field
		UpdatedAt           respjson.Field
		Action              respjson.Field
		Expanded            respjson.Field
		ExpiresAt           respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemWallet) RawJSON() string { return r.JSON.raw }
func (r *VaultItemWallet) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Live data that can currently be requested by passing its type to the item GET
// expand parameter.
type VaultItemWalletAvailableExpansion struct {
	Description string `json:"description" api:"required"`
	// Any of "payment_methods".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemWalletAvailableExpansion) RawJSON() string { return r.JSON.raw }
func (r *VaultItemWalletAvailableExpansion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An operation that is currently valid for this item. Read the description before
// invoking it through the item operations endpoint.
type VaultItemWalletAvailableOperation struct {
	Description string `json:"description" api:"required"`
	// Any of "authorize", "collect", "prepare_checkout", "fill".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemWalletAvailableOperation) RawJSON() string { return r.JSON.raw }
func (r *VaultItemWalletAvailableOperation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Live, non-persisted data requested through the item GET expand parameter.
type VaultItemWalletExpanded struct {
	PaymentMethods []VaultPaymentMethod `json:"payment_methods"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		PaymentMethods respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemWalletExpanded) RawJSON() string { return r.JSON.raw }
func (r *VaultItemWalletExpanded) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemCard struct {
	ID                  string                            `json:"id" api:"required"`
	AvailableExpansions []VaultItemCardAvailableExpansion `json:"available_expansions" api:"required"`
	AvailableOperations []VaultItemCardAvailableOperation `json:"available_operations" api:"required"`
	CreatedAt           time.Time                         `json:"created_at" api:"required" format:"date-time"`
	// Immutable item key assigned when the item is created.
	Key string `json:"key" api:"required"`
	// Live payment card. Test-mode card creation is not supported.
	Spec      CardVaultItemSpecUnion  `json:"spec" api:"required"`
	State     CardVaultItemStateUnion `json:"state" api:"required"`
	Type      constant.Card           `json:"type" default:"card"`
	UpdatedAt time.Time               `json:"updated_at" api:"required" format:"date-time"`
	Action    VaultItemActionUnion    `json:"action"`
	ExpiresAt time.Time               `json:"expires_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		AvailableExpansions respjson.Field
		AvailableOperations respjson.Field
		CreatedAt           respjson.Field
		Key                 respjson.Field
		Spec                respjson.Field
		State               respjson.Field
		Type                respjson.Field
		UpdatedAt           respjson.Field
		Action              respjson.Field
		ExpiresAt           respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemCard) RawJSON() string { return r.JSON.raw }
func (r *VaultItemCard) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Live data that can currently be requested by passing its type to the item GET
// expand parameter.
type VaultItemCardAvailableExpansion struct {
	Description string `json:"description" api:"required"`
	// Any of "payment_methods".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemCardAvailableExpansion) RawJSON() string { return r.JSON.raw }
func (r *VaultItemCardAvailableExpansion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An operation that is currently valid for this item. Read the description before
// invoking it through the item operations endpoint.
type VaultItemCardAvailableOperation struct {
	Description string `json:"description" api:"required"`
	// Any of "authorize", "collect", "prepare_checkout", "fill".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemCardAvailableOperation) RawJSON() string { return r.JSON.raw }
func (r *VaultItemCardAvailableOperation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemActionUnion contains all possible properties and values from
// [VaultItemActionLinkOAuth], [VaultItemActionSpendApproval],
// [VaultItemActionPushApproval], [VaultItemActionCollect], [VaultItemActionMfa],
// [VaultItemActionEmbeddedCeremony], [VaultItemActionCardEnrollment].
//
// Use the [VaultItemActionUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type VaultItemActionUnion struct {
	// Any of "link_oauth", "spend_approval", "push_approval", "collect", "mfa",
	// "embedded_ceremony", "card_enrollment".
	Name string `json:"name"`
	URL  string `json:"url"`
	JSON struct {
		Name respjson.Field
		URL  respjson.Field
		raw  string
	} `json:"-"`
}

// anyVaultItemAction is implemented by each variant of [VaultItemActionUnion] to
// add type safety for the return type of [VaultItemActionUnion.AsAny]
type anyVaultItemAction interface {
	implVaultItemActionUnion()
}

func (VaultItemActionLinkOAuth) implVaultItemActionUnion()        {}
func (VaultItemActionSpendApproval) implVaultItemActionUnion()    {}
func (VaultItemActionPushApproval) implVaultItemActionUnion()     {}
func (VaultItemActionCollect) implVaultItemActionUnion()          {}
func (VaultItemActionMfa) implVaultItemActionUnion()              {}
func (VaultItemActionEmbeddedCeremony) implVaultItemActionUnion() {}
func (VaultItemActionCardEnrollment) implVaultItemActionUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := VaultItemActionUnion.AsAny().(type) {
//	case kernel.VaultItemActionLinkOAuth:
//	case kernel.VaultItemActionSpendApproval:
//	case kernel.VaultItemActionPushApproval:
//	case kernel.VaultItemActionCollect:
//	case kernel.VaultItemActionMfa:
//	case kernel.VaultItemActionEmbeddedCeremony:
//	case kernel.VaultItemActionCardEnrollment:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u VaultItemActionUnion) AsAny() anyVaultItemAction {
	switch u.Name {
	case "link_oauth":
		return u.AsLinkOAuth()
	case "spend_approval":
		return u.AsSpendApproval()
	case "push_approval":
		return u.AsPushApproval()
	case "collect":
		return u.AsCollect()
	case "mfa":
		return u.AsMfa()
	case "embedded_ceremony":
		return u.AsEmbeddedCeremony()
	case "card_enrollment":
		return u.AsCardEnrollment()
	}
	return nil
}

func (u VaultItemActionUnion) AsLinkOAuth() (v VaultItemActionLinkOAuth) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultItemActionUnion) AsSpendApproval() (v VaultItemActionSpendApproval) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultItemActionUnion) AsPushApproval() (v VaultItemActionPushApproval) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultItemActionUnion) AsCollect() (v VaultItemActionCollect) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultItemActionUnion) AsMfa() (v VaultItemActionMfa) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultItemActionUnion) AsEmbeddedCeremony() (v VaultItemActionEmbeddedCeremony) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultItemActionUnion) AsCardEnrollment() (v VaultItemActionCardEnrollment) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u VaultItemActionUnion) RawJSON() string { return u.JSON.raw }

func (r *VaultItemActionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemActionLinkOAuth struct {
	Name constant.LinkOAuth `json:"name" default:"link_oauth"`
	URL  string             `json:"url" api:"required" format:"uri"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemActionLinkOAuth) RawJSON() string { return r.JSON.raw }
func (r *VaultItemActionLinkOAuth) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemActionSpendApproval struct {
	Name constant.SpendApproval `json:"name" default:"spend_approval"`
	URL  string                 `json:"url" api:"required" format:"uri"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemActionSpendApproval) RawJSON() string { return r.JSON.raw }
func (r *VaultItemActionSpendApproval) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemActionPushApproval struct {
	Name constant.PushApproval `json:"name" default:"push_approval"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemActionPushApproval) RawJSON() string { return r.JSON.raw }
func (r *VaultItemActionPushApproval) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemActionCollect struct {
	Name constant.Collect `json:"name" default:"collect"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemActionCollect) RawJSON() string { return r.JSON.raw }
func (r *VaultItemActionCollect) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemActionMfa struct {
	Name constant.Mfa `json:"name" default:"mfa"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemActionMfa) RawJSON() string { return r.JSON.raw }
func (r *VaultItemActionMfa) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemActionEmbeddedCeremony struct {
	Name constant.EmbeddedCeremony `json:"name" default:"embedded_ceremony"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemActionEmbeddedCeremony) RawJSON() string { return r.JSON.raw }
func (r *VaultItemActionEmbeddedCeremony) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemActionCardEnrollment struct {
	Name constant.CardEnrollment `json:"name" default:"card_enrollment"`
	URL  string                  `json:"url" api:"required" format:"uri"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemActionCardEnrollment) RawJSON() string { return r.JSON.raw }
func (r *VaultItemActionCardEnrollment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemEvent struct {
	ID        string    `json:"id" api:"required"`
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	Name      string    `json:"name" api:"required"`
	// Browser session associated with the event, when applicable.
	BrowserID string         `json:"browser_id"`
	Data      map[string]any `json:"data"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		BrowserID   respjson.Field
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemEvent) RawJSON() string { return r.JSON.raw }
func (r *VaultItemEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemOperationResponseUnion contains all possible properties and values from
// [VaultItemOperationResponseWalletVaultItem],
// [VaultItemOperationResponseCardVaultItem], [CredentialVaultItem],
// [FillVaultItemOperationResult].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type VaultItemOperationResponseUnion struct {
	ID string `json:"id"`
	// This field is a union of
	// [[]VaultItemOperationResponseWalletVaultItemAvailableExpansion],
	// [[]VaultItemOperationResponseCardVaultItemAvailableExpansion],
	// [[]CredentialVaultItemAvailableExpansion]
	AvailableExpansions VaultItemOperationResponseUnionAvailableExpansions `json:"available_expansions"`
	// This field is a union of
	// [[]VaultItemOperationResponseWalletVaultItemAvailableOperation],
	// [[]VaultItemOperationResponseCardVaultItemAvailableOperation],
	// [[]CredentialVaultItemAvailableOperation]
	AvailableOperations VaultItemOperationResponseUnionAvailableOperations `json:"available_operations"`
	CreatedAt           time.Time                                          `json:"created_at"`
	Key                 string                                             `json:"key"`
	// This field is a union of [WalletVaultItemSpecUnion], [CardVaultItemSpecUnion],
	// [CredentialVaultItemSpec]
	Spec VaultItemOperationResponseUnionSpec `json:"spec"`
	// This field is a union of [WalletVaultItemStateUnion], [CardVaultItemStateUnion],
	// [CredentialVaultItemState]
	State     VaultItemOperationResponseUnionState `json:"state"`
	Type      string                               `json:"type"`
	UpdatedAt time.Time                            `json:"updated_at"`
	// This field is a union of [VaultItemActionUnion], [CredentialCollectionAction]
	Action VaultItemOperationResponseUnionAction `json:"action"`
	// This field is from variant [VaultItemOperationResponseWalletVaultItem].
	Expanded  VaultItemOperationResponseWalletVaultItemExpanded `json:"expanded"`
	ExpiresAt time.Time                                         `json:"expires_at"`
	// This field is from variant [CredentialVaultItem].
	Version int64 `json:"version"`
	// This field is from variant [FillVaultItemOperationResult].
	Fields []VaultFillFieldResult `json:"fields"`
	// This field is from variant [FillVaultItemOperationResult].
	Status FillVaultItemOperationResultStatus `json:"status"`
	JSON   struct {
		ID                  respjson.Field
		AvailableExpansions respjson.Field
		AvailableOperations respjson.Field
		CreatedAt           respjson.Field
		Key                 respjson.Field
		Spec                respjson.Field
		State               respjson.Field
		Type                respjson.Field
		UpdatedAt           respjson.Field
		Action              respjson.Field
		Expanded            respjson.Field
		ExpiresAt           respjson.Field
		Version             respjson.Field
		Fields              respjson.Field
		Status              respjson.Field
		raw                 string
	} `json:"-"`
}

func (u VaultItemOperationResponseUnion) AsVaultItemOperationResponseWalletVaultItem() (v VaultItemOperationResponseWalletVaultItem) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultItemOperationResponseUnion) AsVaultItemOperationResponseCardVaultItem() (v VaultItemOperationResponseCardVaultItem) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultItemOperationResponseUnion) AsCredentialVaultItem() (v CredentialVaultItem) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultItemOperationResponseUnion) AsFillVaultItemOperationResult() (v FillVaultItemOperationResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u VaultItemOperationResponseUnion) RawJSON() string { return u.JSON.raw }

func (r *VaultItemOperationResponseUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemOperationResponseUnionAvailableExpansions is an implicit subunion of
// [VaultItemOperationResponseUnion].
// VaultItemOperationResponseUnionAvailableExpansions provides convenient access to
// the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemOperationResponseUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfVaultItemOperationResponseWalletVaultItemAvailableExpansions
// OfVaultItemOperationResponseCardVaultItemAvailableExpansions
// OfCredentialVaultItemAvailableExpansions]
type VaultItemOperationResponseUnionAvailableExpansions struct {
	// This field will be present if the value is a
	// [[]VaultItemOperationResponseWalletVaultItemAvailableExpansion] instead of an
	// object.
	OfVaultItemOperationResponseWalletVaultItemAvailableExpansions []VaultItemOperationResponseWalletVaultItemAvailableExpansion `json:",inline"`
	// This field will be present if the value is a
	// [[]VaultItemOperationResponseCardVaultItemAvailableExpansion] instead of an
	// object.
	OfVaultItemOperationResponseCardVaultItemAvailableExpansions []VaultItemOperationResponseCardVaultItemAvailableExpansion `json:",inline"`
	// This field will be present if the value is a
	// [[]CredentialVaultItemAvailableExpansion] instead of an object.
	OfCredentialVaultItemAvailableExpansions []CredentialVaultItemAvailableExpansion `json:",inline"`
	JSON                                     struct {
		OfVaultItemOperationResponseWalletVaultItemAvailableExpansions respjson.Field
		OfVaultItemOperationResponseCardVaultItemAvailableExpansions   respjson.Field
		OfCredentialVaultItemAvailableExpansions                       respjson.Field
		raw                                                            string
	} `json:"-"`
}

func (r *VaultItemOperationResponseUnionAvailableExpansions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemOperationResponseUnionAvailableOperations is an implicit subunion of
// [VaultItemOperationResponseUnion].
// VaultItemOperationResponseUnionAvailableOperations provides convenient access to
// the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemOperationResponseUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfVaultItemOperationResponseWalletVaultItemAvailableOperations
// OfVaultItemOperationResponseCardVaultItemAvailableOperations
// OfCredentialVaultItemAvailableOperations]
type VaultItemOperationResponseUnionAvailableOperations struct {
	// This field will be present if the value is a
	// [[]VaultItemOperationResponseWalletVaultItemAvailableOperation] instead of an
	// object.
	OfVaultItemOperationResponseWalletVaultItemAvailableOperations []VaultItemOperationResponseWalletVaultItemAvailableOperation `json:",inline"`
	// This field will be present if the value is a
	// [[]VaultItemOperationResponseCardVaultItemAvailableOperation] instead of an
	// object.
	OfVaultItemOperationResponseCardVaultItemAvailableOperations []VaultItemOperationResponseCardVaultItemAvailableOperation `json:",inline"`
	// This field will be present if the value is a
	// [[]CredentialVaultItemAvailableOperation] instead of an object.
	OfCredentialVaultItemAvailableOperations []CredentialVaultItemAvailableOperation `json:",inline"`
	JSON                                     struct {
		OfVaultItemOperationResponseWalletVaultItemAvailableOperations respjson.Field
		OfVaultItemOperationResponseCardVaultItemAvailableOperations   respjson.Field
		OfCredentialVaultItemAvailableOperations                       respjson.Field
		raw                                                            string
	} `json:"-"`
}

func (r *VaultItemOperationResponseUnionAvailableOperations) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemOperationResponseUnionSpec is an implicit subunion of
// [VaultItemOperationResponseUnion]. VaultItemOperationResponseUnionSpec provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemOperationResponseUnion].
type VaultItemOperationResponseUnionSpec struct {
	// This field is from variant [WalletVaultItemSpecUnion].
	Authorization WalletVaultItemSpecLinkAuthorization `json:"authorization"`
	Provider      string                               `json:"provider"`
	// This field is from variant [WalletVaultItemSpecUnion].
	ProviderConfig WalletVaultItemSpecAgentcardProviderConfig `json:"provider_config"`
	// This field is from variant [WalletVaultItemSpecUnion].
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
	// This field is from variant [CardVaultItemSpecUnion].
	Context  string `json:"context"`
	Currency string `json:"currency"`
	// This field is from variant [CardVaultItemSpecUnion].
	MerchantName string `json:"merchant_name"`
	// This field is from variant [CardVaultItemSpecUnion].
	MerchantURL string `json:"merchant_url"`
	// This field is from variant [CardVaultItemSpecUnion].
	PaymentMethodID string `json:"payment_method_id"`
	Wallet          string `json:"wallet"`
	// This field is from variant [CardVaultItemSpecUnion].
	ExpiresAt int64 `json:"expires_at"`
	// This field is from variant [CardVaultItemSpecUnion].
	LineItems []CardVaultItemSpecLinkLineItem `json:"line_items"`
	// This field is from variant [CardVaultItemSpecUnion].
	Metadata map[string]string `json:"metadata"`
	// This field is from variant [CardVaultItemSpecUnion].
	Totals []CardVaultItemSpecLinkTotal `json:"totals"`
	// This field is from variant [CardVaultItemSpecUnion].
	Merchant string `json:"merchant"`
	// This field is from variant [CardVaultItemSpecUnion].
	CardID string `json:"card_id"`
	// This field is from variant [CredentialVaultItemSpec].
	Fields map[string]CredentialVaultFieldDefinition `json:"fields"`
	// This field is from variant [CredentialVaultItemSpec].
	Description string `json:"description"`
	JSON        struct {
		Authorization   respjson.Field
		Provider        respjson.Field
		ProviderConfig  respjson.Field
		UserID          respjson.Field
		Amount          respjson.Field
		Context         respjson.Field
		Currency        respjson.Field
		MerchantName    respjson.Field
		MerchantURL     respjson.Field
		PaymentMethodID respjson.Field
		Wallet          respjson.Field
		ExpiresAt       respjson.Field
		LineItems       respjson.Field
		Metadata        respjson.Field
		Totals          respjson.Field
		Merchant        respjson.Field
		CardID          respjson.Field
		Fields          respjson.Field
		Description     respjson.Field
		raw             string
	} `json:"-"`
}

func (r *VaultItemOperationResponseUnionSpec) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemOperationResponseUnionState is an implicit subunion of
// [VaultItemOperationResponseUnion]. VaultItemOperationResponseUnionState provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemOperationResponseUnion].
type VaultItemOperationResponseUnionState struct {
	Provider     string `json:"provider"`
	Status       string `json:"status"`
	StatusReason string `json:"status_reason"`
	// This field is from variant [WalletVaultItemStateUnion].
	UserID string `json:"user_id"`
	// This field is from variant [CardVaultItemStateUnion].
	Aliases VaultCardAliases `json:"aliases"`
	// This field is from variant [CardVaultItemStateUnion].
	Domains []string `json:"domains"`
	// This field is a union of [CardVaultItemStateLinkMasks],
	// [CardVaultItemStateAgentcardMasks]
	Masks VaultItemOperationResponseUnionStateMasks `json:"masks"`
	// This field is from variant [CardVaultItemStateUnion].
	Authorization AgentcardCheckoutAuthorization `json:"authorization"`
	// This field is from variant [CardVaultItemStateUnion].
	Preparation AgentcardCheckoutPreparation `json:"preparation"`
	// This field is from variant [CredentialVaultItemState].
	Fields map[string]CredentialVaultFieldState `json:"fields"`
	JSON   struct {
		Provider      respjson.Field
		Status        respjson.Field
		StatusReason  respjson.Field
		UserID        respjson.Field
		Aliases       respjson.Field
		Domains       respjson.Field
		Masks         respjson.Field
		Authorization respjson.Field
		Preparation   respjson.Field
		Fields        respjson.Field
		raw           string
	} `json:"-"`
}

func (r *VaultItemOperationResponseUnionState) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemOperationResponseUnionStateMasks is an implicit subunion of
// [VaultItemOperationResponseUnion]. VaultItemOperationResponseUnionStateMasks
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemOperationResponseUnion].
type VaultItemOperationResponseUnionStateMasks struct {
	Brand string `json:"brand"`
	Last4 string `json:"last4"`
	JSON  struct {
		Brand respjson.Field
		Last4 respjson.Field
		raw   string
	} `json:"-"`
}

func (r *VaultItemOperationResponseUnionStateMasks) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// VaultItemOperationResponseUnionAction is an implicit subunion of
// [VaultItemOperationResponseUnion]. VaultItemOperationResponseUnionAction
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [VaultItemOperationResponseUnion].
type VaultItemOperationResponseUnionAction struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	// This field is from variant [CredentialCollectionAction].
	ExpiresAt time.Time `json:"expires_at"`
	JSON      struct {
		Name      respjson.Field
		URL       respjson.Field
		ExpiresAt respjson.Field
		raw       string
	} `json:"-"`
}

func (r *VaultItemOperationResponseUnionAction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemOperationResponseWalletVaultItem struct {
	ID                  string                                                        `json:"id" api:"required"`
	AvailableExpansions []VaultItemOperationResponseWalletVaultItemAvailableExpansion `json:"available_expansions" api:"required"`
	AvailableOperations []VaultItemOperationResponseWalletVaultItemAvailableOperation `json:"available_operations" api:"required"`
	CreatedAt           time.Time                                                     `json:"created_at" api:"required" format:"date-time"`
	// Immutable item key assigned when the item is created.
	Key string `json:"key" api:"required"`
	// AgentCard wallet. Omit provider_config to use Kernel-managed credentials, or
	// select a customer-owned configuration. Mode (sandbox vs live) is determined by
	// the selected credential; there is no per-item test flag. Without user_id,
	// creation returns a hosted enrollment action and Kernel polls until the user
	// connects. user_id may only reference a user already enrolled by a wallet in this
	// organization under the same configuration.
	Spec  WalletVaultItemSpecUnion  `json:"spec" api:"required"`
	State WalletVaultItemStateUnion `json:"state" api:"required"`
	// Any of "wallet".
	Type      string               `json:"type" api:"required"`
	UpdatedAt time.Time            `json:"updated_at" api:"required" format:"date-time"`
	Action    VaultItemActionUnion `json:"action"`
	// Live, non-persisted data requested through the item GET expand parameter.
	Expanded  VaultItemOperationResponseWalletVaultItemExpanded `json:"expanded"`
	ExpiresAt time.Time                                         `json:"expires_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		AvailableExpansions respjson.Field
		AvailableOperations respjson.Field
		CreatedAt           respjson.Field
		Key                 respjson.Field
		Spec                respjson.Field
		State               respjson.Field
		Type                respjson.Field
		UpdatedAt           respjson.Field
		Action              respjson.Field
		Expanded            respjson.Field
		ExpiresAt           respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemOperationResponseWalletVaultItem) RawJSON() string { return r.JSON.raw }
func (r *VaultItemOperationResponseWalletVaultItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Live data that can currently be requested by passing its type to the item GET
// expand parameter.
type VaultItemOperationResponseWalletVaultItemAvailableExpansion struct {
	Description string `json:"description" api:"required"`
	// Any of "payment_methods".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemOperationResponseWalletVaultItemAvailableExpansion) RawJSON() string {
	return r.JSON.raw
}
func (r *VaultItemOperationResponseWalletVaultItemAvailableExpansion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An operation that is currently valid for this item. Read the description before
// invoking it through the item operations endpoint.
type VaultItemOperationResponseWalletVaultItemAvailableOperation struct {
	Description string `json:"description" api:"required"`
	// Any of "authorize", "collect", "prepare_checkout", "fill".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemOperationResponseWalletVaultItemAvailableOperation) RawJSON() string {
	return r.JSON.raw
}
func (r *VaultItemOperationResponseWalletVaultItemAvailableOperation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Live, non-persisted data requested through the item GET expand parameter.
type VaultItemOperationResponseWalletVaultItemExpanded struct {
	PaymentMethods []VaultPaymentMethod `json:"payment_methods"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		PaymentMethods respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemOperationResponseWalletVaultItemExpanded) RawJSON() string { return r.JSON.raw }
func (r *VaultItemOperationResponseWalletVaultItemExpanded) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemOperationResponseCardVaultItem struct {
	ID                  string                                                      `json:"id" api:"required"`
	AvailableExpansions []VaultItemOperationResponseCardVaultItemAvailableExpansion `json:"available_expansions" api:"required"`
	AvailableOperations []VaultItemOperationResponseCardVaultItemAvailableOperation `json:"available_operations" api:"required"`
	CreatedAt           time.Time                                                   `json:"created_at" api:"required" format:"date-time"`
	// Immutable item key assigned when the item is created.
	Key string `json:"key" api:"required"`
	// Live payment card. Test-mode card creation is not supported.
	Spec  CardVaultItemSpecUnion  `json:"spec" api:"required"`
	State CardVaultItemStateUnion `json:"state" api:"required"`
	// Any of "card".
	Type      string               `json:"type" api:"required"`
	UpdatedAt time.Time            `json:"updated_at" api:"required" format:"date-time"`
	Action    VaultItemActionUnion `json:"action"`
	ExpiresAt time.Time            `json:"expires_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		AvailableExpansions respjson.Field
		AvailableOperations respjson.Field
		CreatedAt           respjson.Field
		Key                 respjson.Field
		Spec                respjson.Field
		State               respjson.Field
		Type                respjson.Field
		UpdatedAt           respjson.Field
		Action              respjson.Field
		ExpiresAt           respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemOperationResponseCardVaultItem) RawJSON() string { return r.JSON.raw }
func (r *VaultItemOperationResponseCardVaultItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Live data that can currently be requested by passing its type to the item GET
// expand parameter.
type VaultItemOperationResponseCardVaultItemAvailableExpansion struct {
	Description string `json:"description" api:"required"`
	// Any of "payment_methods".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemOperationResponseCardVaultItemAvailableExpansion) RawJSON() string {
	return r.JSON.raw
}
func (r *VaultItemOperationResponseCardVaultItemAvailableExpansion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An operation that is currently valid for this item. Read the description before
// invoking it through the item operations endpoint.
type VaultItemOperationResponseCardVaultItemAvailableOperation struct {
	Description string `json:"description" api:"required"`
	// Any of "authorize", "collect", "prepare_checkout", "fill".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultItemOperationResponseCardVaultItemAvailableOperation) RawJSON() string {
	return r.JSON.raw
}
func (r *VaultItemOperationResponseCardVaultItemAvailableOperation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultPaymentMethod struct {
	ID string `json:"id" api:"required"`
	// Provider-reported advisory capabilities. A missing capability is unknown, not
	// ineligible; only eligible=false is an explicit negative signal.
	Capabilities VaultPaymentMethodCapabilities `json:"capabilities" api:"required"`
	Display      VaultPaymentMethodDisplay      `json:"display" api:"required"`
	IsDefault    bool                           `json:"is_default" api:"required"`
	// Provider that issued this payment-method ID.
	Provider string `json:"provider" api:"required"`
	// Provider-neutral payment-method type normalized to lowercase.
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		Capabilities respjson.Field
		Display      respjson.Field
		IsDefault    respjson.Field
		Provider     respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultPaymentMethod) RawJSON() string { return r.JSON.raw }
func (r *VaultPaymentMethod) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Provider-reported advisory capabilities. A missing capability is unknown, not
// ineligible; only eligible=false is an explicit negative signal.
type VaultPaymentMethodCapabilities struct {
	SingleUseCard VaultPaymentMethodCapabilitiesSingleUseCard `json:"single_use_card"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SingleUseCard respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultPaymentMethodCapabilities) RawJSON() string { return r.JSON.raw }
func (r *VaultPaymentMethodCapabilities) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultPaymentMethodCapabilitiesSingleUseCard struct {
	Eligible bool     `json:"eligible" api:"required"`
	Reasons  []string `json:"reasons" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Eligible    respjson.Field
		Reasons     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultPaymentMethodCapabilitiesSingleUseCard) RawJSON() string { return r.JSON.raw }
func (r *VaultPaymentMethodCapabilitiesSingleUseCard) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultPaymentMethodDisplay struct {
	Brand string `json:"brand"`
	Label string `json:"label"`
	Last4 string `json:"last4"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Brand       respjson.Field
		Label       respjson.Field
		Last4       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultPaymentMethodDisplay) RawJSON() string { return r.JSON.raw }
func (r *VaultPaymentMethodDisplay) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WalletVaultItemSpecUnion contains all possible properties and values from
// [WalletVaultItemSpecLink], [WalletVaultItemSpecAgentcard].
//
// Use the [WalletVaultItemSpecUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type WalletVaultItemSpecUnion struct {
	// This field is from variant [WalletVaultItemSpecLink].
	Authorization WalletVaultItemSpecLinkAuthorization `json:"authorization"`
	// Any of "link", "agentcard".
	Provider string `json:"provider"`
	// This field is from variant [WalletVaultItemSpecAgentcard].
	ProviderConfig WalletVaultItemSpecAgentcardProviderConfig `json:"provider_config"`
	// This field is from variant [WalletVaultItemSpecAgentcard].
	UserID string `json:"user_id"`
	JSON   struct {
		Authorization  respjson.Field
		Provider       respjson.Field
		ProviderConfig respjson.Field
		UserID         respjson.Field
		raw            string
	} `json:"-"`
}

// anyWalletVaultItemSpec is implemented by each variant of
// [WalletVaultItemSpecUnion] to add type safety for the return type of
// [WalletVaultItemSpecUnion.AsAny]
type anyWalletVaultItemSpec interface {
	implWalletVaultItemSpecUnion()
}

func (WalletVaultItemSpecLink) implWalletVaultItemSpecUnion()      {}
func (WalletVaultItemSpecAgentcard) implWalletVaultItemSpecUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := WalletVaultItemSpecUnion.AsAny().(type) {
//	case kernel.WalletVaultItemSpecLink:
//	case kernel.WalletVaultItemSpecAgentcard:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u WalletVaultItemSpecUnion) AsAny() anyWalletVaultItemSpec {
	switch u.Provider {
	case "link":
		return u.AsLink()
	case "agentcard":
		return u.AsAgentcard()
	}
	return nil
}

func (u WalletVaultItemSpecUnion) AsLink() (v WalletVaultItemSpecLink) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WalletVaultItemSpecUnion) AsAgentcard() (v WalletVaultItemSpecAgentcard) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u WalletVaultItemSpecUnion) RawJSON() string { return u.JSON.raw }

func (r *WalletVaultItemSpecUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WalletVaultItemSpecLink struct {
	Authorization WalletVaultItemSpecLinkAuthorization `json:"authorization" api:"required"`
	Provider      constant.Link                        `json:"provider" default:"link"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Authorization respjson.Field
		Provider      respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletVaultItemSpecLink) RawJSON() string { return r.JSON.raw }
func (r *WalletVaultItemSpecLink) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WalletVaultItemSpecLinkAuthorization struct {
	Client WalletVaultItemSpecLinkAuthorizationClientUnion `json:"client" api:"required"`
	// Any of "oauth".
	Method string `json:"method" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Client      respjson.Field
		Method      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletVaultItemSpecLinkAuthorization) RawJSON() string { return r.JSON.raw }
func (r *WalletVaultItemSpecLinkAuthorization) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WalletVaultItemSpecLinkAuthorizationClientUnion contains all possible properties
// and values from [WalletVaultItemSpecLinkAuthorizationClientKernelManaged],
// [WalletVaultItemSpecLinkAuthorizationClientCustomerManaged].
//
// Use the [WalletVaultItemSpecLinkAuthorizationClientUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type WalletVaultItemSpecLinkAuthorizationClientUnion struct {
	// Any of "kernel_managed", "customer_managed".
	Type string `json:"type"`
	// This field is from variant
	// [WalletVaultItemSpecLinkAuthorizationClientCustomerManaged].
	ProviderConfig WalletVaultItemSpecLinkAuthorizationClientCustomerManagedProviderConfig `json:"provider_config"`
	JSON           struct {
		Type           respjson.Field
		ProviderConfig respjson.Field
		raw            string
	} `json:"-"`
}

// anyWalletVaultItemSpecLinkAuthorizationClient is implemented by each variant of
// [WalletVaultItemSpecLinkAuthorizationClientUnion] to add type safety for the
// return type of [WalletVaultItemSpecLinkAuthorizationClientUnion.AsAny]
type anyWalletVaultItemSpecLinkAuthorizationClient interface {
	implWalletVaultItemSpecLinkAuthorizationClientUnion()
}

func (WalletVaultItemSpecLinkAuthorizationClientKernelManaged) implWalletVaultItemSpecLinkAuthorizationClientUnion() {
}
func (WalletVaultItemSpecLinkAuthorizationClientCustomerManaged) implWalletVaultItemSpecLinkAuthorizationClientUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := WalletVaultItemSpecLinkAuthorizationClientUnion.AsAny().(type) {
//	case kernel.WalletVaultItemSpecLinkAuthorizationClientKernelManaged:
//	case kernel.WalletVaultItemSpecLinkAuthorizationClientCustomerManaged:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u WalletVaultItemSpecLinkAuthorizationClientUnion) AsAny() anyWalletVaultItemSpecLinkAuthorizationClient {
	switch u.Type {
	case "kernel_managed":
		return u.AsKernelManaged()
	case "customer_managed":
		return u.AsCustomerManaged()
	}
	return nil
}

func (u WalletVaultItemSpecLinkAuthorizationClientUnion) AsKernelManaged() (v WalletVaultItemSpecLinkAuthorizationClientKernelManaged) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WalletVaultItemSpecLinkAuthorizationClientUnion) AsCustomerManaged() (v WalletVaultItemSpecLinkAuthorizationClientCustomerManaged) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u WalletVaultItemSpecLinkAuthorizationClientUnion) RawJSON() string { return u.JSON.raw }

func (r *WalletVaultItemSpecLinkAuthorizationClientUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WalletVaultItemSpecLinkAuthorizationClientKernelManaged struct {
	Type constant.KernelManaged `json:"type" default:"kernel_managed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletVaultItemSpecLinkAuthorizationClientKernelManaged) RawJSON() string { return r.JSON.raw }
func (r *WalletVaultItemSpecLinkAuthorizationClientKernelManaged) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WalletVaultItemSpecLinkAuthorizationClientCustomerManaged struct {
	// Select a provider config by ID or name. Responses return the ID. Renaming a
	// config does not change existing wallet bindings; a wallet cannot switch to a
	// different config after creation.
	ProviderConfig WalletVaultItemSpecLinkAuthorizationClientCustomerManagedProviderConfig `json:"provider_config" api:"required"`
	Type           constant.CustomerManaged                                                `json:"type" default:"customer_managed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ProviderConfig respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletVaultItemSpecLinkAuthorizationClientCustomerManaged) RawJSON() string {
	return r.JSON.raw
}
func (r *WalletVaultItemSpecLinkAuthorizationClientCustomerManaged) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Select a provider config by ID or name. Responses return the ID. Renaming a
// config does not change existing wallet bindings; a wallet cannot switch to a
// different config after creation.
type WalletVaultItemSpecLinkAuthorizationClientCustomerManagedProviderConfig struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletVaultItemSpecLinkAuthorizationClientCustomerManagedProviderConfig) RawJSON() string {
	return r.JSON.raw
}
func (r *WalletVaultItemSpecLinkAuthorizationClientCustomerManagedProviderConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentCard wallet. Omit provider_config to use Kernel-managed credentials, or
// select a customer-owned configuration. Mode (sandbox vs live) is determined by
// the selected credential; there is no per-item test flag. Without user_id,
// creation returns a hosted enrollment action and Kernel polls until the user
// connects. user_id may only reference a user already enrolled by a wallet in this
// organization under the same configuration.
type WalletVaultItemSpecAgentcard struct {
	Provider constant.Agentcard `json:"provider" default:"agentcard"`
	// Select an AgentCard configuration. The wallet's configuration cannot be changed
	// after creation.
	ProviderConfig WalletVaultItemSpecAgentcardProviderConfig `json:"provider_config"`
	UserID         string                                     `json:"user_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Provider       respjson.Field
		ProviderConfig respjson.Field
		UserID         respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletVaultItemSpecAgentcard) RawJSON() string { return r.JSON.raw }
func (r *WalletVaultItemSpecAgentcard) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Select an AgentCard configuration. The wallet's configuration cannot be changed
// after creation.
type WalletVaultItemSpecAgentcardProviderConfig struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletVaultItemSpecAgentcardProviderConfig) RawJSON() string { return r.JSON.raw }
func (r *WalletVaultItemSpecAgentcardProviderConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WalletVaultItemStateUnion contains all possible properties and values from
// [WalletVaultItemStateLink], [WalletVaultItemStateAgentcard].
//
// Use the [WalletVaultItemStateUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type WalletVaultItemStateUnion struct {
	// Any of "link", "agentcard".
	Provider     string `json:"provider"`
	Status       string `json:"status"`
	StatusReason string `json:"status_reason"`
	// This field is from variant [WalletVaultItemStateAgentcard].
	UserID string `json:"user_id"`
	JSON   struct {
		Provider     respjson.Field
		Status       respjson.Field
		StatusReason respjson.Field
		UserID       respjson.Field
		raw          string
	} `json:"-"`
}

// anyWalletVaultItemState is implemented by each variant of
// [WalletVaultItemStateUnion] to add type safety for the return type of
// [WalletVaultItemStateUnion.AsAny]
type anyWalletVaultItemState interface {
	implWalletVaultItemStateUnion()
}

func (WalletVaultItemStateLink) implWalletVaultItemStateUnion()      {}
func (WalletVaultItemStateAgentcard) implWalletVaultItemStateUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := WalletVaultItemStateUnion.AsAny().(type) {
//	case kernel.WalletVaultItemStateLink:
//	case kernel.WalletVaultItemStateAgentcard:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u WalletVaultItemStateUnion) AsAny() anyWalletVaultItemState {
	switch u.Provider {
	case "link":
		return u.AsLink()
	case "agentcard":
		return u.AsAgentcard()
	}
	return nil
}

func (u WalletVaultItemStateUnion) AsLink() (v WalletVaultItemStateLink) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WalletVaultItemStateUnion) AsAgentcard() (v WalletVaultItemStateAgentcard) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u WalletVaultItemStateUnion) RawJSON() string { return u.JSON.raw }

func (r *WalletVaultItemStateUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WalletVaultItemStateLink struct {
	Provider constant.Link `json:"provider" default:"link"`
	// Any of "pending_authorization", "connected", "declined", "reconnect_required",
	// "degraded".
	Status       string `json:"status" api:"required"`
	StatusReason string `json:"status_reason"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Provider     respjson.Field
		Status       respjson.Field
		StatusReason respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletVaultItemStateLink) RawJSON() string { return r.JSON.raw }
func (r *WalletVaultItemStateLink) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WalletVaultItemStateAgentcard struct {
	Provider constant.Agentcard `json:"provider" default:"agentcard"`
	// Any of "pending_authorization", "connected", "degraded".
	Status       string `json:"status" api:"required"`
	StatusReason string `json:"status_reason"`
	// AgentCard user id linked to this wallet. Present once connected.
	UserID string `json:"user_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Provider     respjson.Field
		Status       respjson.Field
		StatusReason respjson.Field
		UserID       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletVaultItemStateAgentcard) RawJSON() string { return r.JSON.raw }
func (r *WalletVaultItemStateAgentcard) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemGetParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`
	// Hold for up to this many seconds while the item is pending authorization,
	// approval, or credential collection. Return the current item when ready or when
	// the wait elapses. This does not wait for edits to an already-ready credential;
	// poll GET without wait and compare version to observe changes after collect.
	Wait param.Opt[int64] `query:"wait,omitzero" json:"-"`
	// Live fields advertised by `available_expansions` to include in `expanded`.
	//
	// Any of "payment_methods".
	Expand []string `query:"expand,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [VaultItemGetParams]'s query parameters as `url.Values`.
func (r VaultItemGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type VaultItemUpdateParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`

	//
	// Request body variants
	//

	// This field is a request body variant, only one variant field can be set.
	OfCardVaultItemUpdateRequest *VaultItemUpdateParamsBodyCardVaultItemUpdateRequest `json:",inline"`
	// This field is a request body variant, only one variant field can be set.
	// Atomically update description and selected values. Omitted properties are
	// preserved. Field names, types, required flags, and sensitivity cannot change.
	// Unknown field names return 400; stale versions or mismatched item types return
	// 409 without changing the item. A successful update increments version and
	// invalidates outstanding Kernel-hosted collection sessions. If required values
	// remain missing, return pending_collection and a fresh collection action.
	// Otherwise return ready without an action; collect can open the form again
	// without clearing values. Customer URLs have no Kernel-managed expiry.
	OfCredentialVaultItemUpdateRequest *CredentialVaultItemUpdateRequestParam `json:",inline"`

	paramObj
}

func (u VaultItemUpdateParams) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfCardVaultItemUpdateRequest, u.OfCredentialVaultItemUpdateRequest)
}
func (r *VaultItemUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Spec is required.
type VaultItemUpdateParamsBodyCardVaultItemUpdateRequest struct {
	// Live payment card. Test-mode card creation is not supported.
	Spec CardVaultItemSpecUnionParam `json:"spec,omitzero" api:"required"`
	// Any of "card".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r VaultItemUpdateParamsBodyCardVaultItemUpdateRequest) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpdateParamsBodyCardVaultItemUpdateRequest
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpdateParamsBodyCardVaultItemUpdateRequest) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[VaultItemUpdateParamsBodyCardVaultItemUpdateRequest](
		"type", "card",
	)
}

type VaultItemDeleteParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`
	paramObj
}

type VaultItemEventsParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`
	// Return events after this event ID.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Long-poll for new events for up to this many seconds.
	Wait param.Opt[int64] `query:"wait,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [VaultItemEventsParams]'s query parameters as `url.Values`.
func (r VaultItemEventsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type VaultItemPerformOperationParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`

	//
	// Request body variants
	//

	// This field is a request body variant, only one variant field can be set.
	// Authorize a Link card using its existing purchase specification. Use only after
	// explicit user approval and when the item advertises authorize. Do not
	// automatically retry provider failures or indeterminate outcomes. Checkout
	// context is not accepted.
	OfAuthorize *AuthorizeVaultItemOperationRequestParam `json:",inline"`
	// This field is a request body variant, only one variant field can be set. Return
	// the credential item with its collection action. Supported for ready and
	// pending_collection credential items. Always render the same form from every
	// form-supported field; totp fields have no form input and are omitted. No
	// caller-selected field subsets or form overrides are accepted. Reuse an active
	// Kernel-hosted session or renew an expired session atomically. Customer-hosted
	// forms use their own backend and ordinary item GET/PATCH. Opening the form does
	// not clear values or change readiness or item version. To observe edits on a
	// ready item, record its version and poll GET without wait until the version
	// changes, then reconcile the returned state. Version changes may also come from
	// PATCH; they do not identify a particular form submission. Customer-hosted apps
	// use their own submission callback, including for unchanged forms. The wait
	// parameter waits for readiness, not edits.
	OfCollect *CollectVaultItemOperationRequestParam `json:",inline"`
	// This field is a request body variant, only one variant field can be set. Prepare
	// an unused AgentCard card for a supported tokenization checkout. Deliver the
	// returned approval URL and keep the approval page open. Poll the item until
	// ready_to_submit, then submit native Pay before preparation.expires_at. Readiness
	// lasts at most 30 seconds. Unused preparations expire automatically. Preparations
	// are single-use even after failure or expiry; do not automatically retry and
	// reconcile uncertain outcomes with the merchant.
	OfPrepareCheckout *PrepareCheckoutVaultItemOperationRequestParam `json:",inline"`
	// This field is a request body variant, only one variant field can be set. Fill
	// selected fields from one ready credential or ready, unexpired Link card into a
	// browser linked to its vault. Only invoke when the item advertises `fill`.
	// Browser and vault must belong to the same project. Kernel checks access and
	// allowed destinations before filling; providing a page URL does not authorize a
	// destination.
	//
	// Find exactly one open page matching `page_url`. Credential items may omit
	// `page_url` to require exactly one open page; cards require an HTTPS page URL.
	// Credentials have no destination allowlist. TOTP fields generate a current code
	// immediately before writing; their seeds never enter the browser. For each
	// selector, search the main frame and all descendant frames for editable inputs or
	// selects matched directly or contained within matching elements. Each selector
	// must resolve to one unique editable element across all frames; zero or multiple
	// candidates fail. Count each element once, even if multiple matching containers
	// contain it. Validate all bindings before filling. Select elements match an
	// option by its value, not its label. If the page navigates or a target disappears
	// during filling, stop rather than selecting a different page or element.
	//
	// Fill in request order and stop on the first failure. This operation is not
	// atomic: previously filled fields are not rolled back. Never submit the form or
	// click buttons, though input/change events may trigger site behavior. Fill is the
	// preferred browser-checkout path. Aliases remain an alternative for explicitly
	// chosen egress-substitution integrations. Do not automatically retry or fall back
	// to aliases after a failed or indeterminate operation.
	//
	// Secret values are never returned or included in operation logs, traces, audit
	// events, or error details. This does not prevent an agent with unrestricted
	// browser access from reading values from the page or other browser observation
	// surfaces.
	OfFill *FillVaultItemOperationRequestParam `json:",inline"`

	paramObj
}

func (u VaultItemPerformOperationParams) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAuthorize, u.OfCollect, u.OfPrepareCheckout, u.OfFill)
}
func (r *VaultItemPerformOperationParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemUpsertParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`

	//
	// Request body variants
	//

	// This field is a request body variant, only one variant field can be set.
	OfWallet *VaultItemUpsertParamsBodyWallet `json:",inline"`
	// This field is a request body variant, only one variant field can be set.
	OfCard *VaultItemUpsertParamsBodyCard `json:",inline"`
	// This field is a request body variant, only one variant field can be set. Create
	// a credential item without a wallet or external provider. Do not use credential
	// items to store, collect, or fill credit card data, including card numbers
	// (PANs), security codes (CVV/CVC), or expiration dates. Use wallet and card item
	// types for credit cards and payment checkout instead. If all required fields have
	// values, return ready without a collection action; collect can still open its
	// form. Otherwise return pending_collection with a time-scoped Kernel-hosted
	// collection action. Missing optional fields alone do not trigger collection.
	// Repeating the original creation request returns the current item without
	// overwriting later edits; a different request at the same key returns 409. Use
	// PATCH for updates. Required totp fields must include a valid seed on creation;
	// otherwise return 400 rather than opening a form that cannot collect it. Optional
	// totp fields may be unset and populated later through PATCH.
	OfCredential *CredentialVaultItemRequestParam `json:",inline"`

	paramObj
}

func (u VaultItemUpsertParams) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfWallet, u.OfCard, u.OfCredential)
}
func (r *VaultItemUpsertParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Spec, Type are required.
type VaultItemUpsertParamsBodyWallet struct {
	// AgentCard wallet. Omit provider_config to use Kernel-managed credentials, or
	// select a customer-owned configuration. Mode (sandbox vs live) is determined by
	// the selected credential; there is no per-item test flag. Without user_id,
	// creation returns a hosted enrollment action and Kernel polls until the user
	// connects. user_id may only reference a user already enrolled by a wallet in this
	// organization under the same configuration.
	Spec VaultItemUpsertParamsBodyWalletSpecUnion `json:"spec,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "wallet".
	Type constant.Wallet `json:"type" default:"wallet"`
	paramObj
}

func (r VaultItemUpsertParamsBodyWallet) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpsertParamsBodyWallet
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpsertParamsBodyWallet) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type VaultItemUpsertParamsBodyWalletSpecUnion struct {
	OfLink      *VaultItemUpsertParamsBodyWalletSpecLink      `json:",omitzero,inline"`
	OfAgentcard *VaultItemUpsertParamsBodyWalletSpecAgentcard `json:",omitzero,inline"`
	paramUnion
}

func (u VaultItemUpsertParamsBodyWalletSpecUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLink, u.OfAgentcard)
}
func (u *VaultItemUpsertParamsBodyWalletSpecUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *VaultItemUpsertParamsBodyWalletSpecUnion) asAny() any {
	if !param.IsOmitted(u.OfLink) {
		return u.OfLink
	} else if !param.IsOmitted(u.OfAgentcard) {
		return u.OfAgentcard
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultItemUpsertParamsBodyWalletSpecUnion) GetAuthorization() *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnion {
	if vt := u.OfLink; vt != nil {
		return &vt.Authorization
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultItemUpsertParamsBodyWalletSpecUnion) GetProviderConfig() *VaultItemUpsertParamsBodyWalletSpecAgentcardProviderConfig {
	if vt := u.OfAgentcard; vt != nil {
		return &vt.ProviderConfig
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultItemUpsertParamsBodyWalletSpecUnion) GetUserID() *string {
	if vt := u.OfAgentcard; vt != nil && vt.UserID.Valid() {
		return &vt.UserID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultItemUpsertParamsBodyWalletSpecUnion) GetProvider() *string {
	if vt := u.OfLink; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfAgentcard; vt != nil {
		return (*string)(&vt.Provider)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[VaultItemUpsertParamsBodyWalletSpecUnion](
		"provider",
		apijson.Discriminator[VaultItemUpsertParamsBodyWalletSpecLink]("link"),
		apijson.Discriminator[VaultItemUpsertParamsBodyWalletSpecAgentcard]("agentcard"),
	)
}

// The properties Authorization, Provider are required.
type VaultItemUpsertParamsBodyWalletSpecLink struct {
	// Kernel starts and completes the user's Link authorization flow.
	Authorization VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnion `json:"authorization,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "link".
	Provider constant.Link `json:"provider" default:"link"`
	paramObj
}

func (r VaultItemUpsertParamsBodyWalletSpecLink) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpsertParamsBodyWalletSpecLink
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpsertParamsBodyWalletSpecLink) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnion struct {
	OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput `json:",omitzero,inline"`
	OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput      *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput      `json:",omitzero,inline"`
	paramUnion
}

func (u VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput, u.OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput)
}
func (u *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnion) asAny() any {
	if !param.IsOmitted(u.OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput) {
		return u.OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput
	} else if !param.IsOmitted(u.OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput) {
		return u.OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnion) GetTokens() *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputTokens {
	if vt := u.OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput; vt != nil {
		return &vt.Tokens
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnion) GetMethod() *string {
	if vt := u.OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput; vt != nil {
		return (*string)(&vt.Method)
	} else if vt := u.OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput; vt != nil {
		return (*string)(&vt.Method)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnion) GetClient() (res vaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnionClient) {
	if vt := u.OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput; vt != nil {
		res.any = &vt.Client
	} else if vt := u.OfVaultItemUpsertsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput; vt != nil {
		res.any = &vt.Client
	}
	return
}

// Can have the runtime types
// [*VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInputClient],
// [*VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClient]
type vaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnionClient struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *kernel.VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInputClient:
//	case *kernel.VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClient:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u vaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnionClient) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u vaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnionClient) GetProviderConfig() *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClientProviderConfig {
	switch vt := u.any.(type) {
	case *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClient:
		return &vt.ProviderConfig
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u vaultItemUpsertParamsBodyWalletSpecLinkAuthorizationUnionClient) GetType() *string {
	switch vt := u.any.(type) {
	case *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInputClient:
		return (*string)(&vt.Type)
	case *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClient:
		return (*string)(&vt.Type)
	}
	return nil
}

// Kernel starts and completes the user's Link authorization flow.
//
// The properties Client, Method are required.
type VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput struct {
	Client VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInputClient `json:"client,omitzero" api:"required"`
	// Any of "oauth".
	Method string `json:"method,omitzero" api:"required"`
	paramObj
}

func (r VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInput](
		"method", "oauth",
	)
}

// The property Type is required.
type VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInputClient struct {
	// Any of "kernel_managed".
	Type string `json:"type,omitzero" api:"required"`
	paramObj
}

func (r VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInputClient) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInputClient
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInputClient) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationKernelManagedLinkAuthorizationInputClient](
		"type", "kernel_managed",
	)
}

// The customer's backend completes Link OAuth and supplies the resulting tokens.
// For a new wallet, Kernel verifies the access token can access Link payment
// methods without consuming or rotating the refresh token. Valid access creates a
// wallet with state.status=connected. An expired, invalid, revoked, or
// insufficiently scoped access token returns 400 and no wallet is created. Refresh
// expired tokens in your backend before importing them. A failed import does not
// modify existing wallets. After successful import, Kernel owns subsequent
// refresh-token rotation; the customer must stop refreshing this grant. Import
// does not verify the refresh token: if it or the configured client credentials
// are rejected during a later refresh, the imported wallet becomes degraded. An
// unknown refresh outcome also leaves it degraded; Kernel does not retry a refresh
// token that may already have been consumed. There is no in-place reauthorization
// operation for an imported wallet. If this imported wallet's credentials become
// unusable, obtain a fresh Link OAuth grant in your backend and create a wallet
// under a NEW wallet key. Use the new wallet for NEW cards and payments, not to
// retry an old payment whose outcome is uncertain. This does not replace the old
// grant, rebind existing cards, or resolve their payment outcomes. Retain the old
// wallet and its cards while reconciling any uncertain payments with the provider
// or support. Do not repeat an uncertain payment on the new wallet, and do not
// treat deletion as evidence that it did not execute. Deletion of the old wallet
// can remain blocked by unresolved child cards. Repeating a create for the same
// item key and non-secret spec returns the existing wallet without replacing
// tokens, even if they have rotated or the wallet needs reconnection. ID and name
// references resolving to the same config are equivalent. A different config or
// non-secret spec returns 409. This create operation does not replace an existing
// grant.
//
// The properties Client, Method, Tokens are required.
type VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput struct {
	Client VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClient `json:"client,omitzero" api:"required"`
	// Any of "oauth".
	Method string `json:"method,omitzero" api:"required"`
	// Send the token pair from your backend. Both tokens must be from the same Link
	// grant under the referenced client. Supply a currently valid access token. Kernel
	// refreshes when needed after import and uses the expiry returned by Link for
	// subsequent tokens. Tokens are never returned in wallet responses, events, or
	// logs.
	Tokens VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputTokens `json:"tokens,omitzero" api:"required"`
	paramObj
}

func (r VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInput](
		"method", "oauth",
	)
}

// The properties ProviderConfig, Type are required.
type VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClient struct {
	// Select a provider config by ID or name. Responses return the ID. Renaming a
	// config does not change existing wallet bindings; a wallet cannot switch to a
	// different config after creation.
	ProviderConfig VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClientProviderConfig `json:"provider_config,omitzero" api:"required"`
	// Any of "customer_managed".
	Type string `json:"type,omitzero" api:"required"`
	paramObj
}

func (r VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClient) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClient
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClient) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClient](
		"type", "customer_managed",
	)
}

// Select a provider config by ID or name. Responses return the ID. Renaming a
// config does not change existing wallet bindings; a wallet cannot switch to a
// different config after creation.
type VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClientProviderConfig struct {
	ID   param.Opt[string] `json:"id,omitzero"`
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClientProviderConfig) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClientProviderConfig
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputClientProviderConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send the token pair from your backend. Both tokens must be from the same Link
// grant under the referenced client. Supply a currently valid access token. Kernel
// refreshes when needed after import and uses the expiry returned by Link for
// subsequent tokens. Tokens are never returned in wallet responses, events, or
// logs.
//
// The properties AccessToken, RefreshToken are required.
type VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputTokens struct {
	AccessToken  string `json:"access_token" api:"required"`
	RefreshToken string `json:"refresh_token" api:"required"`
	paramObj
}

func (r VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputTokens) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputTokens
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpsertParamsBodyWalletSpecLinkAuthorizationImportedLinkAuthorizationInputTokens) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentCard wallet. Omit provider_config to use Kernel-managed credentials, or
// select a customer-owned configuration. Mode (sandbox vs live) is determined by
// the selected credential; there is no per-item test flag. Without user_id,
// creation returns a hosted enrollment action and Kernel polls until the user
// connects. user_id may only reference a user already enrolled by a wallet in this
// organization under the same configuration.
//
// The property Provider is required.
type VaultItemUpsertParamsBodyWalletSpecAgentcard struct {
	UserID param.Opt[string] `json:"user_id,omitzero"`
	// Select an AgentCard configuration. The wallet's configuration cannot be changed
	// after creation.
	ProviderConfig VaultItemUpsertParamsBodyWalletSpecAgentcardProviderConfig `json:"provider_config,omitzero"`
	// This field can be elided, and will marshal its zero value as "agentcard".
	Provider constant.Agentcard `json:"provider" default:"agentcard"`
	paramObj
}

func (r VaultItemUpsertParamsBodyWalletSpecAgentcard) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpsertParamsBodyWalletSpecAgentcard
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpsertParamsBodyWalletSpecAgentcard) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Select an AgentCard configuration. The wallet's configuration cannot be changed
// after creation.
type VaultItemUpsertParamsBodyWalletSpecAgentcardProviderConfig struct {
	ID   param.Opt[string] `json:"id,omitzero"`
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r VaultItemUpsertParamsBodyWalletSpecAgentcardProviderConfig) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpsertParamsBodyWalletSpecAgentcardProviderConfig
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpsertParamsBodyWalletSpecAgentcardProviderConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Spec, Type are required.
type VaultItemUpsertParamsBodyCard struct {
	// Live payment card. Test-mode card creation is not supported.
	Spec CardVaultItemSpecUnionParam `json:"spec,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "card".
	Type constant.Card `json:"type" default:"card"`
	paramObj
}

func (r VaultItemUpsertParamsBodyCard) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpsertParamsBodyCard
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpsertParamsBodyCard) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
