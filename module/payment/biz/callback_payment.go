package biz

import (
	"context"
	"strings"

	"stockflow/module/payment/model"
)

type CallbackPaymentStore interface {
	CallbackPayment(ctx context.Context, data *model.PaymentCallback) (*model.Payment, error)
}

type callbackPaymentBiz struct {
	store CallbackPaymentStore
}

func NewCallbackPaymentBiz(store CallbackPaymentStore) *callbackPaymentBiz {
	return &callbackPaymentBiz{store: store}
}

func (biz *callbackPaymentBiz) CallbackPayment(ctx context.Context, data *model.PaymentCallback) (*model.Payment, error) {
	if data == nil {
		return nil, model.ErrPaymentCallbackDataIsRequired
	}

	data.PaymentID = strings.TrimSpace(data.PaymentID)
	data.PaymentCode = strings.TrimSpace(data.PaymentCode)
	data.Status = strings.TrimSpace(data.Status)

	if data.ExternalTxnID != nil {
		trimmed := strings.TrimSpace(*data.ExternalTxnID)
		data.ExternalTxnID = &trimmed
	}

	if err := data.Validate(); err != nil {
		return nil, err
	}

	updatedPayment, err := biz.store.CallbackPayment(ctx, data)
	if err != nil {
		return nil, err
	}

	if updatedPayment == nil {
		return nil, model.ErrPaymentNotFound
	}

	return updatedPayment, nil
}
