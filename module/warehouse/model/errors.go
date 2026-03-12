package model

import "errors"

var (
	ErrWarehouseDataIsNil          = errors.New("warehouse data is required")
	ErrWarehouseIDIsBlank          = errors.New("warehouse id is required")
	ErrWarehouseCodeIsBlank        = errors.New("warehouse code is required")
	ErrWarehouseNameIsBlank        = errors.New("warehouse name is required")
	ErrWarehouseCodeAlreadyExists  = errors.New("warehouse code already exists")
	ErrWarehouseNotFound           = errors.New("warehouse not found")
	ErrWarehouseUpdateDataRequired = errors.New("warehouse update data is required")
)
