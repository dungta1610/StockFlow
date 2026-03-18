package model

import "errors"

var (
	ErrPaymentDataIsRequired          = errors.New("payment data is required")
	ErrPaymentIDIsBlank               = errors.New("payment id is required")
	ErrPaymentOrderIDIsBlank          = errors.New("payment order id is required")
	ErrPaymentProviderIsBlank         = errors.New("payment provider is required")
	ErrPaymentMethodIsBlank           = errors.New("payment method is required")
	ErrPaymentStatusIsBlank           = errors.New("payment status is required")
	ErrPaymentAmountInvalid           = errors.New("payment amount must be greater than 0")
	ErrPaymentCurrencyIsBlank         = errors.New("payment currency is required")
	ErrPaymentNotFound                = errors.New("payment not found")
	ErrPaymentAlreadySucceeded        = errors.New("payment already succeeded")
	ErrPaymentAlreadyFailed           = errors.New("payment already failed")
	ErrPaymentAlreadyCanceled         = errors.New("payment already canceled")
	ErrPaymentInvalidStatusTransition = errors.New("payment invalid status transition")

	ErrCheckoutDataIsRequired        = errors.New("checkout data is required")
	ErrPaymentCallbackDataIsRequired = errors.New("payment callback data is required")
)
