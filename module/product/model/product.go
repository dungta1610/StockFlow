package model

import (
	"strings"
	"time"
)

type Product struct {
	ID          string    `json:"id" db:"id"`
	SKU         string    `json:"sku" db:"sku"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ProductCreate struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (p *ProductCreate) Validate() error {
	if p == nil {
		return ErrProductDataIsNil
	}

	p.SKU = strings.TrimSpace(strings.ToUpper(p.SKU))
	p.Name = strings.TrimSpace(p.Name)
	p.Description = strings.TrimSpace(p.Description)

	if p.SKU == "" {
		return ErrProductSKUIsBlank
	}

	if p.Name == "" {
		return ErrProductNameIsBlank
	}

	if p.Price < 0 {
		return ErrProductPriceInvalid
	}

	return nil
}

type ProductUpdate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsActive    *bool   `json:"is_active"`
}

func (p *ProductUpdate) Validate() error {
	if p == nil {
		return ErrProductUpdateDataRequired
	}

	p.Name = strings.TrimSpace(p.Name)
	p.Description = strings.TrimSpace(p.Description)

	if p.Name == "" {
		return ErrProductNameIsBlank
	}

	if p.Price < 0 {
		return ErrProductPriceInvalid
	}

	return nil
}

type Filter struct {
	SKU      string `json:"sku" form:"sku"`
	Name     string `json:"name" form:"name"`
	IsActive *bool  `json:"is_active" form:"is_active"`
}

func (f *Filter) Normalize() {
	if f == nil {
		return
	}

	f.SKU = strings.TrimSpace(strings.ToUpper(f.SKU))
	f.Name = strings.TrimSpace(f.Name)
}
