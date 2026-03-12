package model

import (
	"strings"
	"time"
)

type Warehouse struct {
	ID        string    `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	Address   string    `json:"address" db:"address"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type WarehouseCreate struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (w *WarehouseCreate) Validate() error {
	if w == nil {
		return ErrWarehouseDataIsNil
	}

	w.Code = strings.TrimSpace(strings.ToUpper(w.Code))
	w.Name = strings.TrimSpace(w.Name)
	w.Address = strings.TrimSpace(w.Address)

	if w.Code == "" {
		return ErrWarehouseCodeIsBlank
	}

	if w.Name == "" {
		return ErrWarehouseNameIsBlank
	}

	return nil
}

type WarehouseUpdate struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	IsActive *bool  `json:"is_active"`
}

func (w *WarehouseUpdate) Validate() error {
	if w == nil {
		return ErrWarehouseUpdateDataRequired
	}

	w.Name = strings.TrimSpace(w.Name)
	w.Address = strings.TrimSpace(w.Address)

	if w.Name == "" {
		return ErrWarehouseNameIsBlank
	}

	return nil
}

type Filter struct {
	Code     string `json:"code" form:"code"`
	Name     string `json:"name" form:"name"`
	IsActive *bool  `json:"is_active" form:"is_active"`
}

func (f *Filter) Normalize() {
	if f == nil {
		return
	}

	f.Code = strings.TrimSpace(strings.ToUpper(f.Code))
	f.Name = strings.TrimSpace(f.Name)
}
