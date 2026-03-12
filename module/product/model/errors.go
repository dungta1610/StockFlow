package model

import "errors"

var (
	ErrProductDataIsNil          = errors.New("product data is required")
	ErrProductSKUIsBlank         = errors.New("product sku is required")
	ErrProductNameIsBlank        = errors.New("product name is required")
	ErrProductPriceInvalid       = errors.New("product price must be greater than or equal to 0")
	ErrProductSKUAlreadyExists   = errors.New("product sku already exists")
	ErrProductNotFound           = errors.New("product not found")
	ErrProductUpdateDataRequired = errors.New("product update data is required")
	ErrProductIDIsBlank          = errors.New("product id is required")
)
