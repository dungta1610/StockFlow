package model

import (
	"strings"
	"time"
)

type Inventory struct {
	ID           string    `json:"id" db:"id"`
	ProductID    string    `json:"product_id" db:"product_id"`
	WarehouseID  string    `json:"warehouse_id" db:"warehouse_id"`
	AvailableQty int       `json:"available_qty" db:"available_qty"`
	ReservedQty  int       `json:"reserved_qty" db:"reserved_qty"`
	Version      int       `json:"version" db:"version"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type InventoryCreate struct {
	ProductID    string `json:"product_id"`
	WarehouseID  string `json:"warehouse_id"`
	AvailableQty int    `json:"available_qty"`
	ReservedQty  int    `json:"reserved_qty"`
}

func (i *InventoryCreate) Validate() error {
	if i == nil {
		return ErrInventoryDataIsNil
	}

	i.ProductID = strings.TrimSpace(i.ProductID)
	i.WarehouseID = strings.TrimSpace(i.WarehouseID)

	if i.ProductID == "" {
		return ErrInventoryProductIDIsBlank
	}

	if i.WarehouseID == "" {
		return ErrInventoryWarehouseIDIsBlank
	}

	if i.AvailableQty < 0 {
		return ErrInventoryAvailableQtyInvalid
	}

	if i.ReservedQty < 0 {
		return ErrInventoryReservedQtyInvalid
	}

	return nil
}

type InventoryAdjust struct {
	ProductID   string `json:"product_id"`
	WarehouseID string `json:"warehouse_id"`
	Quantity    int    `json:"quantity"`
	Reason      string `json:"reason"`
	CreatedBy   string `json:"created_by"`
}

func (i *InventoryAdjust) Validate() error {
	if i == nil {
		return ErrInventoryAdjustDataRequired
	}

	i.ProductID = strings.TrimSpace(i.ProductID)
	i.WarehouseID = strings.TrimSpace(i.WarehouseID)
	i.Reason = strings.TrimSpace(i.Reason)
	i.CreatedBy = strings.TrimSpace(i.CreatedBy)

	if i.ProductID == "" {
		return ErrInventoryProductIDIsBlank
	}

	if i.WarehouseID == "" {
		return ErrInventoryWarehouseIDIsBlank
	}

	if i.Quantity == 0 {
		return ErrInventoryAdjustQtyInvalid
	}

	return nil
}

type InventoryReserve struct {
	OrderID     string `json:"order_id"`
	OrderItemID string `json:"order_item_id"`
	ProductID   string `json:"product_id"`
	WarehouseID string `json:"warehouse_id"`
	Quantity    int    `json:"quantity"`
	CreatedBy   string `json:"created_by"`
}

func (i *InventoryReserve) Validate() error {
	if i == nil {
		return ErrInventoryReserveDataRequired
	}

	i.OrderID = strings.TrimSpace(i.OrderID)
	i.OrderItemID = strings.TrimSpace(i.OrderItemID)
	i.ProductID = strings.TrimSpace(i.ProductID)
	i.WarehouseID = strings.TrimSpace(i.WarehouseID)
	i.CreatedBy = strings.TrimSpace(i.CreatedBy)

	if i.OrderID == "" {
		return ErrInventoryOrderIDIsBlank
	}

	if i.OrderItemID == "" {
		return ErrInventoryOrderItemIDIsBlank
	}

	if i.ProductID == "" {
		return ErrInventoryProductIDIsBlank
	}

	if i.WarehouseID == "" {
		return ErrInventoryWarehouseIDIsBlank
	}

	if i.Quantity <= 0 {
		return ErrInventoryReserveQtyInvalid
	}

	return nil
}

type Filter struct {
	ProductID   string `json:"product_id" form:"product_id"`
	WarehouseID string `json:"warehouse_id" form:"warehouse_id"`
}

func (f *Filter) Normalize() {
	if f == nil {
		return
	}

	f.ProductID = strings.TrimSpace(f.ProductID)
	f.WarehouseID = strings.TrimSpace(f.WarehouseID)
}
