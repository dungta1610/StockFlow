package model

import (
	"strings"
	"time"
)

type InventoryTransaction struct {
	ID                 string    `json:"id" db:"id"`
	InventoryID        string    `json:"inventory_id" db:"inventory_id"`
	ProductID          string    `json:"product_id" db:"product_id"`
	WarehouseID        string    `json:"warehouse_id" db:"warehouse_id"`
	OrderID            *string   `json:"order_id,omitempty" db:"order_id"`
	ReservationID      *string   `json:"reservation_id,omitempty" db:"reservation_id"`
	TxnType            string    `json:"txn_type" db:"txn_type"`
	Quantity           int       `json:"quantity" db:"quantity"`
	BeforeAvailableQty int       `json:"before_available_qty" db:"before_available_qty"`
	AfterAvailableQty  int       `json:"after_available_qty" db:"after_available_qty"`
	BeforeReservedQty  int       `json:"before_reserved_qty" db:"before_reserved_qty"`
	AfterReservedQty   int       `json:"after_reserved_qty" db:"after_reserved_qty"`
	Reason             string    `json:"reason" db:"reason"`
	CreatedBy          *string   `json:"created_by,omitempty" db:"created_by"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

type InventoryTransactionCreate struct {
	InventoryID        string  `json:"inventory_id"`
	ProductID          string  `json:"product_id"`
	WarehouseID        string  `json:"warehouse_id"`
	OrderID            *string `json:"order_id"`
	ReservationID      *string `json:"reservation_id"`
	TxnType            string  `json:"txn_type"`
	Quantity           int     `json:"quantity"`
	BeforeAvailableQty int     `json:"before_available_qty"`
	AfterAvailableQty  int     `json:"after_available_qty"`
	BeforeReservedQty  int     `json:"before_reserved_qty"`
	AfterReservedQty   int     `json:"after_reserved_qty"`
	Reason             string  `json:"reason"`
	CreatedBy          *string `json:"created_by"`
}

type TransactionFilter struct {
	InventoryID   string `json:"inventory_id" form:"inventory_id"`
	ProductID     string `json:"product_id" form:"product_id"`
	WarehouseID   string `json:"warehouse_id" form:"warehouse_id"`
	OrderID       string `json:"order_id" form:"order_id"`
	ReservationID string `json:"reservation_id" form:"reservation_id"`
	TxnType       string `json:"txn_type" form:"txn_type"`
}

func (f *TransactionFilter) Normalize() {
	if f == nil {
		return
	}

	f.InventoryID = strings.TrimSpace(f.InventoryID)
	f.ProductID = strings.TrimSpace(f.ProductID)
	f.WarehouseID = strings.TrimSpace(f.WarehouseID)
	f.OrderID = strings.TrimSpace(f.OrderID)
	f.ReservationID = strings.TrimSpace(f.ReservationID)
	f.TxnType = strings.TrimSpace(f.TxnType)
}

func (i *InventoryTransactionCreate) Validate() error {
	if i == nil {
		return ErrInventoryTransactionDataRequired
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

	if i.TxnType == "" {
		return ErrInventoryTxnTypeIsBlank
	}

	if i.Quantity <= 0 {
		return ErrInventoryTransactionQtyInvalid
	}

	if i.BeforeAvailableQty < 0 || i.AfterAvailableQty < 0 {
		return ErrInventoryAvailableQtyInvalid
	}

	if i.BeforeReservedQty < 0 || i.AfterReservedQty < 0 {
		return ErrInventoryReservedQtyInvalid
	}

	return nil
}
