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
	OrderStatusFulfilled       = "fulfilled"
	OrderStatusCompleted       = "completed"
	OrderStatusCancelled       = "cancelled"
	OrderStatusExpired         = "expired"
)

type Order struct {
	ID                   string     `json:"id" db:"id"`
	OrderCode            string     `json:"order_code" db:"order_code"`
	UserID               string     `json:"user_id" db:"user_id"`
	WarehouseID          string     `json:"warehouse_id" db:"warehouse_id"`
	Status               string     `json:"status" db:"status"`
	TotalAmount          float64    `json:"total_amount" db:"total_amount"`
	ReservationExpiresAt *time.Time `json:"reservation_expires_at,omitempty" db:"reservation_expires_at"`
	PaidAt               *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	CancelledAt          *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
	FulfilledAt          *time.Time `json:"fulfilled_at,omitempty" db:"fulfilled_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`

	Items []OrderItem `json:"items,omitempty"`
}

type OrderCreate struct {
	UserID               string            `json:"user_id"`
	WarehouseID          string            `json:"warehouse_id"`
	ReservationExpiresAt *time.Time        `json:"reservation_expires_at"`
	ExpiredAt            *time.Time        `json:"expired_at"`
	Items                []OrderItemCreate `json:"items"`
}

func (o *OrderCreate) Validate() error {
	if o == nil {
		return ErrOrderDataIsRequired
	}

	o.UserID = strings.TrimSpace(o.UserID)
	o.WarehouseID = strings.TrimSpace(o.WarehouseID)

	if o.ReservationExpiresAt == nil && o.ExpiredAt != nil {
		o.ReservationExpiresAt = o.ExpiredAt
	}

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
}

func (o *OrderCancel) Validate() error {
	if o == nil {
		return ErrOrderCancelDataIsRequired
	}

	o.OrderID = strings.TrimSpace(o.OrderID)

	if o.OrderID == "" {
		return ErrOrderIDIsBlank
	}

	return nil
}

type OrderExpire struct {
	OrderID string `json:"order_id"`
}

func (o *OrderExpire) Validate() error {
	if o == nil {
		return ErrOrderExpireDataIsRequired
	}

	o.OrderID = strings.TrimSpace(o.OrderID)

	if o.OrderID == "" {
		return ErrOrderIDIsBlank
	}

	return nil
}

type Filter struct {
	OrderCode   string `json:"order_code" form:"order_code"`
	UserID      string `json:"user_id" form:"user_id"`
	WarehouseID string `json:"warehouse_id" form:"warehouse_id"`
	Status      string `json:"status" form:"status"`
}

func (f *Filter) Normalize() {
	if f == nil {
		return
	}

	f.OrderCode = strings.TrimSpace(strings.ToUpper(f.OrderCode))
	f.UserID = strings.TrimSpace(f.UserID)
	f.WarehouseID = strings.TrimSpace(f.WarehouseID)
	f.Status = strings.TrimSpace(f.Status)
}
