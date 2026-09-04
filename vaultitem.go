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
// instead of a partial item.
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

// Update a card specification before or between authorizations
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

// List vault items without secret values
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

// Delete a vault item and invalidate its secret value
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
// `available_operations`, following its natural-language description. Operations
// may call an external provider and can return the item's updated state. If the
// provider rate limits spend-request creation, returns HTTP 429 with code
// `spend_request_rate_limited`; stop and back off before retrying.
func (r *VaultItemService) PerformOperation(ctx context.Context, key string, params VaultItemPerformOperationParams, opts ...option.RequestOption) (res *VaultItemUnion, err error) {
	opts = slices.Concat(r.Options, opts)
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

// Create or retrieve an identical vault item by immutable key
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
	JSON          struct {
		Provider      respjson.Field
		Status        respjson.Field
		Aliases       respjson.Field
		Domains       respjson.Field
		Masks         respjson.Field
		StatusReason  respjson.Field
		Authorization respjson.Field
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
	// Any of "requested", "pending_authorization", "ready", "consumed", "expired",
	// "declined".
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
	// Any of "requested", "ready", "pending_approval", "degraded".
	Status  string           `json:"status" api:"required"`
	Aliases VaultCardAliases `json:"aliases"`
	// The in-flight or most recent checkout authorization. Present while a checkout is
	// pending approval and after it settles.
	Authorization AgentcardCheckoutAuthorization   `json:"authorization"`
	Masks         CardVaultItemStateAgentcardMasks `json:"masks"`
	StatusReason  string                           `json:"status_reason"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Provider      respjson.Field
		Status        respjson.Field
		Aliases       respjson.Field
		Authorization respjson.Field
		Masks         respjson.Field
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

// VaultItemUnion contains all possible properties and values from
// [VaultItemWallet], [VaultItemCard].
//
// Use the [VaultItemUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type VaultItemUnion struct {
	ID string `json:"id"`
	// This field is a union of [[]VaultItemWalletAvailableExpansion],
	// [[]VaultItemCardAvailableExpansion]
	AvailableExpansions VaultItemUnionAvailableExpansions `json:"available_expansions"`
	// This field is a union of [[]VaultItemWalletAvailableOperation],
	// [[]VaultItemCardAvailableOperation]
	AvailableOperations VaultItemUnionAvailableOperations `json:"available_operations"`
	CreatedAt           time.Time                         `json:"created_at"`
	Key                 string                            `json:"key"`
	// This field is a union of [WalletVaultItemSpecUnion], [CardVaultItemSpecUnion]
	Spec VaultItemUnionSpec `json:"spec"`
	// This field is a union of [WalletVaultItemStateUnion], [CardVaultItemStateUnion]
	State VaultItemUnionState `json:"state"`
	// Any of "wallet", "card".
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	// This field is from variant [VaultItemWallet].
	Action VaultItemActionUnion `json:"action"`
	// This field is from variant [VaultItemWallet].
	Expanded  VaultItemWalletExpanded `json:"expanded"`
	ExpiresAt time.Time               `json:"expires_at"`
	JSON      struct {
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
		raw                 string
	} `json:"-"`
}

// anyVaultItem is implemented by each variant of [VaultItemUnion] to add type
// safety for the return type of [VaultItemUnion.AsAny]
type anyVaultItem interface {
	implVaultItemUnion()
}

func (VaultItemWallet) implVaultItemUnion() {}
func (VaultItemCard) implVaultItemUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := VaultItemUnion.AsAny().(type) {
//	case kernel.VaultItemWallet:
//	case kernel.VaultItemCard:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u VaultItemUnion) AsAny() anyVaultItem {
	switch u.Type {
	case "wallet":
		return u.AsWallet()
	case "card":
		return u.AsCard()
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
// OfVaultItemCardAvailableExpansions]
type VaultItemUnionAvailableExpansions struct {
	// This field will be present if the value is a
	// [[]VaultItemWalletAvailableExpansion] instead of an object.
	OfVaultItemWalletAvailableExpansions []VaultItemWalletAvailableExpansion `json:",inline"`
	// This field will be present if the value is a [[]VaultItemCardAvailableExpansion]
	// instead of an object.
	OfVaultItemCardAvailableExpansions []VaultItemCardAvailableExpansion `json:",inline"`
	JSON                               struct {
		OfVaultItemWalletAvailableExpansions respjson.Field
		OfVaultItemCardAvailableExpansions   respjson.Field
		raw                                  string
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
// OfVaultItemCardAvailableOperations]
type VaultItemUnionAvailableOperations struct {
	// This field will be present if the value is a
	// [[]VaultItemWalletAvailableOperation] instead of an object.
	OfVaultItemWalletAvailableOperations []VaultItemWalletAvailableOperation `json:",inline"`
	// This field will be present if the value is a [[]VaultItemCardAvailableOperation]
	// instead of an object.
	OfVaultItemCardAvailableOperations []VaultItemCardAvailableOperation `json:",inline"`
	JSON                               struct {
		OfVaultItemWalletAvailableOperations respjson.Field
		OfVaultItemCardAvailableOperations   respjson.Field
		raw                                  string
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
	JSON   struct {
		Authorization   respjson.Field
		Provider        respjson.Field
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
	JSON          struct {
		Provider      respjson.Field
		Status        respjson.Field
		StatusReason  respjson.Field
		UserID        respjson.Field
		Aliases       respjson.Field
		Domains       respjson.Field
		Masks         respjson.Field
		Authorization respjson.Field
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

type VaultItemWallet struct {
	ID                  string                              `json:"id" api:"required"`
	AvailableExpansions []VaultItemWalletAvailableExpansion `json:"available_expansions" api:"required"`
	AvailableOperations []VaultItemWalletAvailableOperation `json:"available_operations" api:"required"`
	CreatedAt           time.Time                           `json:"created_at" api:"required" format:"date-time"`
	// Immutable item key assigned when the item is created.
	Key string `json:"key" api:"required"`
	// AgentCard wallet. Mode (sandbox vs live) is fixed by the deployment's AgentCard
	// credential; there is no per-item test flag. user_id may only reference a user
	// already enrolled by a wallet in this organization.
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
	// Any of "authorize".
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
	// Any of "authorize".
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
	UserID string `json:"user_id"`
	JSON   struct {
		Authorization respjson.Field
		Provider      respjson.Field
		UserID        respjson.Field
		raw           string
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

// ToParam converts this WalletVaultItemSpecUnion to a
// WalletVaultItemSpecUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// WalletVaultItemSpecUnionParam.Overrides()
func (r WalletVaultItemSpecUnion) ToParam() WalletVaultItemSpecUnionParam {
	return param.Override[WalletVaultItemSpecUnionParam](json.RawMessage(r.RawJSON()))
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
	Client WalletVaultItemSpecLinkAuthorizationClient `json:"client" api:"required"`
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

type WalletVaultItemSpecLinkAuthorizationClient struct {
	// Any of "kernel_managed".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletVaultItemSpecLinkAuthorizationClient) RawJSON() string { return r.JSON.raw }
func (r *WalletVaultItemSpecLinkAuthorizationClient) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentCard wallet. Mode (sandbox vs live) is fixed by the deployment's AgentCard
// credential; there is no per-item test flag. user_id may only reference a user
// already enrolled by a wallet in this organization.
type WalletVaultItemSpecAgentcard struct {
	Provider constant.Agentcard `json:"provider" default:"agentcard"`
	UserID   string             `json:"user_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Provider    respjson.Field
		UserID      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletVaultItemSpecAgentcard) RawJSON() string { return r.JSON.raw }
func (r *WalletVaultItemSpecAgentcard) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func WalletVaultItemSpecParamOfLink(authorization WalletVaultItemSpecLinkAuthorizationParam) WalletVaultItemSpecUnionParam {
	var link WalletVaultItemSpecLinkParam
	link.Authorization = authorization
	return WalletVaultItemSpecUnionParam{OfLink: &link}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type WalletVaultItemSpecUnionParam struct {
	OfLink      *WalletVaultItemSpecLinkParam      `json:",omitzero,inline"`
	OfAgentcard *WalletVaultItemSpecAgentcardParam `json:",omitzero,inline"`
	paramUnion
}

func (u WalletVaultItemSpecUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLink, u.OfAgentcard)
}
func (u *WalletVaultItemSpecUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *WalletVaultItemSpecUnionParam) asAny() any {
	if !param.IsOmitted(u.OfLink) {
		return u.OfLink
	} else if !param.IsOmitted(u.OfAgentcard) {
		return u.OfAgentcard
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WalletVaultItemSpecUnionParam) GetAuthorization() *WalletVaultItemSpecLinkAuthorizationParam {
	if vt := u.OfLink; vt != nil {
		return &vt.Authorization
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WalletVaultItemSpecUnionParam) GetUserID() *string {
	if vt := u.OfAgentcard; vt != nil && vt.UserID.Valid() {
		return &vt.UserID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WalletVaultItemSpecUnionParam) GetProvider() *string {
	if vt := u.OfLink; vt != nil {
		return (*string)(&vt.Provider)
	} else if vt := u.OfAgentcard; vt != nil {
		return (*string)(&vt.Provider)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[WalletVaultItemSpecUnionParam](
		"provider",
		apijson.Discriminator[WalletVaultItemSpecLinkParam]("link"),
		apijson.Discriminator[WalletVaultItemSpecAgentcardParam]("agentcard"),
	)
}

// The properties Authorization, Provider are required.
type WalletVaultItemSpecLinkParam struct {
	Authorization WalletVaultItemSpecLinkAuthorizationParam `json:"authorization,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "link".
	Provider constant.Link `json:"provider" default:"link"`
	paramObj
}

func (r WalletVaultItemSpecLinkParam) MarshalJSON() (data []byte, err error) {
	type shadow WalletVaultItemSpecLinkParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WalletVaultItemSpecLinkParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Client, Method are required.
type WalletVaultItemSpecLinkAuthorizationParam struct {
	Client WalletVaultItemSpecLinkAuthorizationClientParam `json:"client,omitzero" api:"required"`
	// Any of "oauth".
	Method string `json:"method,omitzero" api:"required"`
	paramObj
}

func (r WalletVaultItemSpecLinkAuthorizationParam) MarshalJSON() (data []byte, err error) {
	type shadow WalletVaultItemSpecLinkAuthorizationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WalletVaultItemSpecLinkAuthorizationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[WalletVaultItemSpecLinkAuthorizationParam](
		"method", "oauth",
	)
}

// The property Type is required.
type WalletVaultItemSpecLinkAuthorizationClientParam struct {
	// Any of "kernel_managed".
	Type string `json:"type,omitzero" api:"required"`
	paramObj
}

func (r WalletVaultItemSpecLinkAuthorizationClientParam) MarshalJSON() (data []byte, err error) {
	type shadow WalletVaultItemSpecLinkAuthorizationClientParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WalletVaultItemSpecLinkAuthorizationClientParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[WalletVaultItemSpecLinkAuthorizationClientParam](
		"type", "kernel_managed",
	)
}

// AgentCard wallet. Mode (sandbox vs live) is fixed by the deployment's AgentCard
// credential; there is no per-item test flag. user_id may only reference a user
// already enrolled by a wallet in this organization.
//
// The property Provider is required.
type WalletVaultItemSpecAgentcardParam struct {
	UserID param.Opt[string] `json:"user_id,omitzero"`
	// This field can be elided, and will marshal its zero value as "agentcard".
	Provider constant.Agentcard `json:"provider" default:"agentcard"`
	paramObj
}

func (r WalletVaultItemSpecAgentcardParam) MarshalJSON() (data []byte, err error) {
	type shadow WalletVaultItemSpecAgentcardParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WalletVaultItemSpecAgentcardParam) UnmarshalJSON(data []byte) error {
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
	// Hold for up to this many seconds while the item is pending authorization or
	// approval.
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
	// Live payment card. Test-mode card creation is not supported.
	Spec CardVaultItemSpecUnionParam `json:"spec,omitzero" api:"required"`
	paramObj
}

func (r VaultItemUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
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
	// Any of "authorize".
	Type VaultItemPerformOperationParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r VaultItemPerformOperationParams) MarshalJSON() (data []byte, err error) {
	type shadow VaultItemPerformOperationParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultItemPerformOperationParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultItemPerformOperationParamsType string

const (
	VaultItemPerformOperationParamsTypeAuthorize VaultItemPerformOperationParamsType = "authorize"
)

type VaultItemUpsertParams struct {
	IDOrName string `path:"id_or_name" api:"required" json:"-"`

	//
	// Request body variants
	//

	// This field is a request body variant, only one variant field can be set.
	OfWallet *VaultItemUpsertParamsBodyWallet `json:",inline"`
	// This field is a request body variant, only one variant field can be set.
	OfCard *VaultItemUpsertParamsBodyCard `json:",inline"`

	paramObj
}

func (u VaultItemUpsertParams) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfWallet, u.OfCard)
}
func (r *VaultItemUpsertParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Spec, Type are required.
type VaultItemUpsertParamsBodyWallet struct {
	// AgentCard wallet. Mode (sandbox vs live) is fixed by the deployment's AgentCard
	// credential; there is no per-item test flag. user_id may only reference a user
	// already enrolled by a wallet in this organization.
	Spec WalletVaultItemSpecUnionParam `json:"spec,omitzero" api:"required"`
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
