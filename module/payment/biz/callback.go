package biz

import (
	"context"
	"strings"

	"stockflow/module/payment/model"
)

type CallbackStore interface {
	HandleCallback(ctx context.Context, data *model.Callback) (*model.Payment, error)
}

type callbackBiz struct {
	store CallbackStore
}

func NewCallbackBiz(store CallbackStore) *callbackBiz {
	return &callbackBiz{store: store}
}

func (biz *callbackBiz) HandleCallback(ctx context.Context, data *model.Callback) (*model.Payment, error) {
	if data == nil {
		return nil, model.ErrPaymentCallbackDataIsRequired
	}

	data.PaymentID = strings.TrimSpace(data.PaymentID)
	data.ProviderTxnID = strings.TrimSpace(data.ProviderTxnID)
	data.ProviderOrderCode = strings.TrimSpace(data.ProviderOrderCode)
	data.Status = strings.TrimSpace(strings.ToLower(data.Status))
	data.FailureReason = strings.TrimSpace(data.FailureReason)
	data.RawPayload = strings.TrimSpace(data.RawPayload)
	data.UpdatedBy = strings.TrimSpace(data.UpdatedBy)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	payment, err := biz.store.HandleCallback(ctx, data)
	if err != nil {
		return nil, err
	}

	if payment == nil {
		return nil, model.ErrPaymentNotFound
	}

	return payment, nil
}
