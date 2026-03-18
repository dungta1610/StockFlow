package model

import (
	"strings"
	"time"
)

type OrderItem struct {
	ID          string    `json:"id" db:"id"`
	OrderID     string    `json:"order_id" db:"order_id"`
	ProductID   string    `json:"product_id" db:"product_id"`
	ProductSKU  string    `json:"product_sku" db:"product_sku"`
	ProductName string    `json:"product_name" db:"product_name"`
	Quantity    int       `json:"quantity" db:"quantity"`
	UnitPrice   float64   `json:"unit_price" db:"unit_price"`
	LinePrice   float64   `json:"line_price" db:"line_price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type OrderItemCreate struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

func (i *OrderItemCreate) Validate() error {
	if i == nil {
		return ErrOrderItemDataIsRequired
	}

	i.ProductID = strings.TrimSpace(i.ProductID)

	if i.ProductID == "" {
		return ErrOrderItemProductIDIsBlank
	}

	if i.Quantity <= 0 {
		return ErrOrderItemQuantityInvalid
	}

	if i.UnitPrice < 0 {
		return ErrOrderItemUnitPriceInvalid
	}

	return nil
}
