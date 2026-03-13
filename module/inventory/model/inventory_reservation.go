package model

import (
	"strings"
	"time"
)

type InventoryReservation struct {
	ID          string     `json:"id" db:"id"`
	OrderID     string     `json:"order_id" db:"order_id"`
	OrderItemID string     `json:"order_item_id" db:"order_item_id"`
	InventoryID string     `json:"inventory_id" db:"inventory_id"`
	ProductID   string     `json:"product_id" db:"product_id"`
	WarehouseID string     `json:"warehouse_id" db:"warehouse_id"`
	Quantity    int        `json:"quantity" db:"quantity"`
	Status      string     `json:"status" db:"status"`
	ReservedAt  time.Time  `json:"reserved_at" db:"reserved_at"`
	ReleasedAt  *time.Time `json:"released_at,omitempty" db:"released_at"`
	ConsumedAt  *time.Time `json:"consumed_at,omitempty" db:"consumed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type InventoryReservationCreate struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"order_id"`
	OrderItemID string    `json:"order_item_id"`
	InventoryID string    `json:"inventory_id"`
	ProductID   string    `json:"product_id"`
	WarehouseID string    `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"`
	ReservedAt  time.Time `json:"reserved_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (i *InventoryReservationCreate) Validate() error {
	if i == nil {
		return ErrInventoryReservationDataRequired
	}

	i.OrderID = strings.TrimSpace(i.OrderID)
	i.OrderItemID = strings.TrimSpace(i.OrderItemID)
	i.InventoryID = strings.TrimSpace(i.InventoryID)
	i.ProductID = strings.TrimSpace(i.ProductID)
	i.WarehouseID = strings.TrimSpace(i.WarehouseID)
	i.Status = strings.TrimSpace(i.Status)

	if i.OrderID == "" {
		return ErrInventoryOrderIDIsBlank
	}

	if i.OrderItemID == "" {
		return ErrInventoryOrderItemIDIsBlank
	}

	if i.InventoryID == "" {
		return ErrInventoryIDIsBlank
	}

	if i.ProductID == "" {
		return ErrInventoryProductIDIsBlank
	}

	if i.WarehouseID == "" {
		return ErrInventoryWarehouseIDIsBlank
	}

	if i.Quantity <= 0 {
		return ErrInventoryReservationQtyInvalid
	}

	if i.Status == "" {
		return ErrInventoryReservationStatusIsBlank
	}

	return nil
}

type ReservationFilter struct {
	OrderID     string `json:"order_id" form:"order_id"`
	OrderItemID string `json:"order_item_id" form:"order_item_id"`
	InventoryID string `json:"inventory_id" form:"inventory_id"`
	ProductID   string `json:"product_id" form:"product_id"`
	WarehouseID string `json:"warehouse_id" form:"warehouse_id"`
	Status      string `json:"status" form:"status"`
}

func (f *ReservationFilter) Normalize() {
	if f == nil {
		return
	}

	f.OrderID = strings.TrimSpace(f.OrderID)
	f.OrderItemID = strings.TrimSpace(f.OrderItemID)
	f.InventoryID = strings.TrimSpace(f.InventoryID)
	f.ProductID = strings.TrimSpace(f.ProductID)
	f.WarehouseID = strings.TrimSpace(f.WarehouseID)
	f.Status = strings.TrimSpace(f.Status)
}
