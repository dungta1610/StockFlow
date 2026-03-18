package model

import (
	"strings"
	"time"
)

const (
	PaymentMethodCOD          = "cod"
	PaymentMethodBankTransfer = "bank_transfer"
	PaymentMethodMomo         = "momo"
	PaymentMethodVNPay        = "vnpay"
	PaymentMethodStripe       = "stripe"
)

const (
	PaymentStatusPending    = "pending"
	PaymentStatusProcessing = "processing"
	PaymentStatusSucceeded  = "succeeded"
	PaymentStatusFailed     = "failed"
	PaymentStatusCanceled   = "canceled"
	PaymentStatusExpired    = "expired"
)

type Payment struct {
	ID                string     `json:"id" db:"id"`
	OrderID           string     `json:"order_id" db:"order_id"`
	Provider          string     `json:"provider" db:"provider"`
	Method            string     `json:"method" db:"method"`
	Status            string     `json:"status" db:"status"`
	Amount            float64    `json:"amount" db:"amount"`
	Currency          string     `json:"currency" db:"currency"`
	ProviderTxnID     string     `json:"provider_txn_id" db:"provider_txn_id"`
	ProviderOrderCode string     `json:"provider_order_code" db:"provider_order_code"`
	CheckoutURL       string     `json:"checkout_url" db:"checkout_url"`
	CallbackPayload   string     `json:"callback_payload" db:"callback_payload"`
	FailureReason     string     `json:"failure_reason" db:"failure_reason"`
	PaidAt            *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	ExpiredAt         *time.Time `json:"expired_at,omitempty" db:"expired_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

type Checkout struct {
	OrderID   string `json:"order_id"`
	Method    string `json:"method"`
	Provider  string `json:"provider"`
	CreatedBy string `json:"created_by"`
}

func (c *Checkout) Validate() error {
	if c == nil {
		return ErrCheckoutDataIsRequired
	}

	c.OrderID = strings.TrimSpace(c.OrderID)
	c.Method = strings.TrimSpace(strings.ToLower(c.Method))
	c.Provider = strings.TrimSpace(strings.ToLower(c.Provider))
	c.CreatedBy = strings.TrimSpace(c.CreatedBy)

	if c.OrderID == "" {
		return ErrPaymentOrderIDIsBlank
	}

	if c.Method == "" {
		return ErrPaymentMethodIsBlank
	}

	if c.Provider == "" {
		return ErrPaymentProviderIsBlank
	}

	return nil
}

type Callback struct {
	PaymentID         string `json:"payment_id"`
	ProviderTxnID     string `json:"provider_txn_id"`
	ProviderOrderCode string `json:"provider_order_code"`
	Status            string `json:"status"`
	FailureReason     string `json:"failure_reason"`
	RawPayload        string `json:"raw_payload"`
	UpdatedBy         string `json:"updated_by"`
}

func (c *Callback) Validate() error {
	if c == nil {
		return ErrPaymentCallbackDataIsRequired
	}

	c.PaymentID = strings.TrimSpace(c.PaymentID)
	c.ProviderTxnID = strings.TrimSpace(c.ProviderTxnID)
	c.ProviderOrderCode = strings.TrimSpace(c.ProviderOrderCode)
	c.Status = strings.TrimSpace(strings.ToLower(c.Status))
	c.FailureReason = strings.TrimSpace(c.FailureReason)
	c.RawPayload = strings.TrimSpace(c.RawPayload)
	c.UpdatedBy = strings.TrimSpace(c.UpdatedBy)

	if c.PaymentID == "" {
		return ErrPaymentIDIsBlank
	}

	if c.Status == "" {
		return ErrPaymentStatusIsBlank
	}

	return nil
}

type Filter struct {
	OrderID  string `json:"order_id" form:"order_id"`
	Method   string `json:"method" form:"method"`
	Provider string `json:"provider" form:"provider"`
	Status   string `json:"status" form:"status"`
}

func (f *Filter) Normalize() {
	if f == nil {
		return
	}

	f.OrderID = strings.TrimSpace(f.OrderID)
	f.Method = strings.TrimSpace(strings.ToLower(f.Method))
	f.Provider = strings.TrimSpace(strings.ToLower(f.Provider))
	f.Status = strings.TrimSpace(strings.ToLower(f.Status))
}
