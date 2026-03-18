package model

import (
	"strings"
	"time"
)

const (
	OrderStatusPending         = "pending"
	OrderStatusReserved        = "reserved"
	OrderStatusAwaitingPayment = "awaiting_payment"
	OrderStatusPaid            = "paid"
	OrderStatusCompleted       = "completed"
	OrderStatusCanceled        = "canceled"
	OrderStatusExpired         = "expired"
)

const (
	PaymentStatusPending  = "pending"
	PaymentStatusPaid     = "paid"
	PaymentStatusFailed   = "failed"
	PaymentStatusRefunded = "refunded"
)

type Order struct {
	ID            string     `json:"id" db:"id"`
	Code          string     `json:"code" db:"code"`
	UserID        string     `json:"user_id" db:"user_id"`
	WarehouseID   string     `json:"warehouse_id" db:"warehouse_id"`
	Status        string     `json:"status" db:"status"`
	PaymentStatus string     `json:"payment_status" db:"payment_status"`
	SubtotalPrice float64    `json:"subtotal_price" db:"subtotal_price"`
	DiscountPrice float64    `json:"discount_price" db:"discount_price"`
	TotalPrice    float64    `json:"total_price" db:"total_price"`
	Note          string     `json:"note" db:"note"`
	ExpiredAt     *time.Time `json:"expired_at,omitempty" db:"expired_at"`
	CanceledAt    *time.Time `json:"canceled_at,omitempty" db:"canceled_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`

	Items []OrderItem `json:"items,omitempty"`
}

type OrderCreate struct {
	UserID      string            `json:"user_id"`
	WarehouseID string            `json:"warehouse_id"`
	Note        string            `json:"note"`
	ExpiredAt   *time.Time        `json:"expired_at"`
	CreatedBy   string            `json:"created_by"`
	Items       []OrderItemCreate `json:"items"`
}

func (o *OrderCreate) Validate() error {
	if o == nil {
		return ErrOrderDataIsRequired
	}

	o.UserID = strings.TrimSpace(o.UserID)
	o.WarehouseID = strings.TrimSpace(o.WarehouseID)
	o.Note = strings.TrimSpace(o.Note)
	o.CreatedBy = strings.TrimSpace(o.CreatedBy)

	if o.UserID == "" {
		return ErrOrderUserIDIsBlank
	}

	if o.WarehouseID == "" {
		return ErrOrderWarehouseIDIsBlank
	}

	if len(o.Items) == 0 {
		return ErrOrderItemsIsEmpty
	}

	for idx := range o.Items {
		if err := o.Items[idx].Validate(); err != nil {
			return err
		}
	}

	return nil
}

type OrderCancel struct {
	OrderID string `json:"order_id"`
	Reason  string `json:"reason"`
	By      string `json:"by"`
}

func (o *OrderCancel) Validate() error {
	if o == nil {
		return ErrOrderCancelDataIsRequired
	}

	o.OrderID = strings.TrimSpace(o.OrderID)
	o.Reason = strings.TrimSpace(o.Reason)
	o.By = strings.TrimSpace(o.By)

	if o.OrderID == "" {
		return ErrOrderIDIsBlank
	}

	if o.Reason == "" {
		return ErrOrderCancelReasonIsBlank
	}

	return nil
}

type OrderExpire struct {
	OrderID string `json:"order_id"`
	Reason  string `json:"reason"`
	By      string `json:"by"`
}

func (o *OrderExpire) Validate() error {
	if o == nil {
		return ErrOrderExpireDataIsRequired
	}

	o.OrderID = strings.TrimSpace(o.OrderID)
	o.Reason = strings.TrimSpace(o.Reason)
	o.By = strings.TrimSpace(o.By)

	if o.OrderID == "" {
		return ErrOrderIDIsBlank
	}

	return nil
}

type Filter struct {
	Code          string `json:"code" form:"code"`
	UserID        string `json:"user_id" form:"user_id"`
	WarehouseID   string `json:"warehouse_id" form:"warehouse_id"`
	Status        string `json:"status" form:"status"`
	PaymentStatus string `json:"payment_status" form:"payment_status"`
}

func (f *Filter) Normalize() {
	if f == nil {
		return
	}

	f.Code = strings.TrimSpace(strings.ToUpper(f.Code))
	f.UserID = strings.TrimSpace(f.UserID)
	f.WarehouseID = strings.TrimSpace(f.WarehouseID)
	f.Status = strings.TrimSpace(f.Status)
	f.PaymentStatus = strings.TrimSpace(f.PaymentStatus)
}
