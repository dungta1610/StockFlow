package model

import "errors"

var (
	ErrOrderDataIsRequired          = errors.New("order data is required")
	ErrOrderIDIsBlank               = errors.New("order id is required")
	ErrOrderCodeIsBlank             = errors.New("order code is required")
	ErrOrderUserIDIsBlank           = errors.New("order user id is required")
	ErrOrderWarehouseIDIsBlank      = errors.New("order warehouse id is required")
	ErrOrderStatusIsBlank           = errors.New("order status is required")
	ErrOrderItemsIsEmpty            = errors.New("order items are required")
	ErrOrderTotalAmountInvalid      = errors.New("order total amount must be greater than or equal to 0")
	ErrOrderNotFound                = errors.New("order not found")
	ErrOrderCannotBeCanceled        = errors.New("order cannot be canceled")
	ErrOrderCannotBeCancelled       = ErrOrderCannotBeCanceled
	ErrOrderCannotBeExpired         = errors.New("order cannot be expired")
	ErrOrderCannotBePaid            = errors.New("order cannot be paid")
	ErrOrderAlreadyCanceled         = errors.New("order already canceled")
	ErrOrderAlreadyCancelled        = ErrOrderAlreadyCanceled
	ErrOrderAlreadyExpired          = errors.New("order already expired")
	ErrOrderAlreadyPaid             = errors.New("order already paid")
	ErrOrderInvalidStatusTransition = errors.New("order invalid status transition")

	ErrOrderCancelDataIsRequired = errors.New("order cancel data is required")
	ErrOrderExpireDataIsRequired = errors.New("order expire data is required")

	ErrOrderItemDataIsRequired   = errors.New("order item data is required")
	ErrOrderItemIDIsBlank        = errors.New("order item id is required")
	ErrOrderItemOrderIDIsBlank   = errors.New("order item order id is required")
	ErrOrderItemProductIDIsBlank = errors.New("order item product id is required")
	ErrOrderItemQuantityInvalid  = errors.New("order item quantity must be greater than 0")
	ErrOrderItemUnitPriceInvalid = errors.New("order item unit price must be greater than or equal to 0")
	ErrOrderItemLineTotalInvalid = errors.New("order item line total must be greater than or equal to 0")
	ErrOrderItemLinePriceInvalid = ErrOrderItemLineTotalInvalid
)
