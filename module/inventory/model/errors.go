package model

import "errors"

var (
	ErrInventoryDataIsNil                = errors.New("inventory data is required")
	ErrInventoryIDIsBlank                = errors.New("inventory id is required")
	ErrInventoryProductIDIsBlank         = errors.New("inventory product id is required")
	ErrInventoryWarehouseIDIsBlank       = errors.New("inventory warehouse id is required")
	ErrInventoryOrderIDIsBlank           = errors.New("inventory order id is required")
	ErrInventoryOrderItemIDIsBlank       = errors.New("inventory order item id is required")
	ErrInventoryAvailableQtyInvalid      = errors.New("inventory available quantity must be greater than or equal to 0")
	ErrInventoryReservedQtyInvalid       = errors.New("inventory reserved quantity must be greater than or equal to 0")
	ErrInventoryAdjustDataRequired       = errors.New("inventory adjust data is required")
	ErrInventoryAdjustQtyInvalid         = errors.New("inventory adjust quantity must not be 0")
	ErrInventoryReserveDataRequired      = errors.New("inventory reserve data is required")
	ErrInventoryReserveQtyInvalid        = errors.New("inventory reserve quantity must be greater than 0")
	ErrInventoryNotFound                 = errors.New("inventory not found")
	ErrInventoryAlreadyExists            = errors.New("inventory already exists")
	ErrInventoryNotEnoughStock           = errors.New("inventory not enough stock")
	ErrInventoryTransactionNotFound      = errors.New("inventory transaction not found")
	ErrInventoryReservationNotFound      = errors.New("inventory reservation not found")
	ErrInventoryTransactionDataRequired  = errors.New("inventory transaction data is required")
	ErrInventoryTxnTypeIsBlank           = errors.New("inventory transaction type is required")
	ErrInventoryTransactionQtyInvalid    = errors.New("inventory transaction quantity must be greater than 0")
	ErrInventoryReservationDataRequired  = errors.New("inventory reservation data is required")
	ErrInventoryReservationQtyInvalid    = errors.New("inventory reservation quantity must be greater than 0")
	ErrInventoryReservationStatusIsBlank = errors.New("inventory reservation status is required")
)
