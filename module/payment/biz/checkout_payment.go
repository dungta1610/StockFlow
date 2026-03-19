package biz

import (
	"context"
	"strings"

	"stockflow/module/payment/model"
)

type CheckoutPaymentStore interface {
	CheckoutPayment(ctx context.Context, data *model.PaymentCheckout) (*model.Payment, error)
}

type checkoutPaymentBiz struct {
	store CheckoutPaymentStore
}

func NewCheckoutPaymentBiz(store CheckoutPaymentStore) *checkoutPaymentBiz {
	return &checkoutPaymentBiz{store: store}
}

func (biz *checkoutPaymentBiz) CheckoutPayment(ctx context.Context, data *model.PaymentCheckout) (*model.Payment, error) {
	if data == nil {
		return nil, model.ErrPaymentCheckoutDataIsRequired
	}

	data.OrderID = strings.TrimSpace(data.OrderID)
	data.Method = strings.TrimSpace(data.Method)
	data.IdempotencyKey = strings.TrimSpace(data.IdempotencyKey)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	createdPayment, err := biz.store.CheckoutPayment(ctx, data)
	if err != nil {
		return nil, err
	}

	if createdPayment == nil {
		return nil, model.ErrPaymentNotFound
	}

	return createdPayment, nil
}
