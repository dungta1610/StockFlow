package model

import (
	"strings"
	"time"
)

type Payment struct {
	ID             string     `json:"id" db:"id"`
	OrderID        string     `json:"order_id" db:"order_id"`
	PaymentCode    string     `json:"payment_code" db:"payment_code"`
	Method         string     `json:"method" db:"method"`
	Status         string     `json:"status" db:"status"`
	Amount         float64    `json:"amount" db:"amount"`
	IdempotencyKey string     `json:"idempotency_key" db:"idempotency_key"`
	ExternalTxnID  *string    `json:"external_txn_id,omitempty" db:"external_txn_id"`
	PaidAt         *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	FailedAt       *time.Time `json:"failed_at,omitempty" db:"failed_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type PaymentCheckout struct {
	OrderID        string  `json:"order_id"`
	Method         string  `json:"method"`
	Amount         float64 `json:"amount"`
	IdempotencyKey string  `json:"idempotency_key"`
}

func (p *PaymentCheckout) Validate() error {
	if p == nil {
		return ErrPaymentCheckoutDataIsRequired
	}

	p.OrderID = strings.TrimSpace(p.OrderID)
	p.Method = strings.TrimSpace(p.Method)
	p.IdempotencyKey = strings.TrimSpace(p.IdempotencyKey)

	if p.OrderID == "" {
		return ErrPaymentOrderIDIsBlank
	}

	if p.Method == "" {
		return ErrPaymentMethodIsBlank
	}

	if p.Amount < 0 {
		return ErrPaymentAmountInvalid
	}

	if p.IdempotencyKey == "" {
		return ErrPaymentIdempotencyKeyIsBlank
	}

	return nil
}

type PaymentCallback struct {
	PaymentID     string  `json:"payment_id"`
	PaymentCode   string  `json:"payment_code"`
	Status        string  `json:"status"`
	ExternalTxnID *string `json:"external_txn_id"`
}

func (p *PaymentCallback) Validate() error {
	if p == nil {
		return ErrPaymentCallbackDataIsRequired
	}

	p.PaymentID = strings.TrimSpace(p.PaymentID)
	p.PaymentCode = strings.TrimSpace(p.PaymentCode)
	p.Status = strings.TrimSpace(p.Status)

	if p.PaymentID == "" && p.PaymentCode == "" {
		return ErrPaymentIDIsBlank
	}

	if p.Status == "" {
		return ErrPaymentStatusIsBlank
	}

	if p.ExternalTxnID != nil {
		trimmed := strings.TrimSpace(*p.ExternalTxnID)
		p.ExternalTxnID = &trimmed
	}

	return nil
}

type Filter struct {
	OrderID     string `json:"order_id" form:"order_id"`
	PaymentCode string `json:"payment_code" form:"payment_code"`
	Method      string `json:"method" form:"method"`
	Status      string `json:"status" form:"status"`
}

func (f *Filter) Normalize() {
	if f == nil {
		return
	}

	f.OrderID = strings.TrimSpace(f.OrderID)
	f.PaymentCode = strings.TrimSpace(strings.ToUpper(f.PaymentCode))
	f.Method = strings.TrimSpace(f.Method)
	f.Status = strings.TrimSpace(f.Status)
}
