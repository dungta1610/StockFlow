package model

import "errors"

var (
	ErrPaymentDataIsRequired         = errors.New("payment data is required")
	ErrPaymentCheckoutDataIsRequired = errors.New("payment checkout data is required")
	ErrPaymentCallbackDataIsRequired = errors.New("payment callback data is required")

	ErrPaymentIDIsBlank             = errors.New("payment id is required")
	ErrPaymentCodeIsBlank           = errors.New("payment code is required")
	ErrPaymentOrderIDIsBlank        = errors.New("payment order id is required")
	ErrPaymentMethodIsBlank         = errors.New("payment method is required")
	ErrPaymentStatusIsBlank         = errors.New("payment status is required")
	ErrPaymentAmountInvalid         = errors.New("payment amount must be greater than or equal to 0")
	ErrPaymentIdempotencyKeyIsBlank = errors.New("payment idempotency key is required")

	ErrPaymentNotFound          = errors.New("payment not found")
	ErrPaymentAlreadyPaid       = errors.New("payment already paid")
	ErrPaymentAlreadyFailed     = errors.New("payment already failed")
	ErrPaymentCannotBeProcessed = errors.New("payment cannot be processed")
)
