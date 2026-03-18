package biz

import (
	"context"
	"strings"

	"stockflow/module/payment/model"
)

type CheckoutStore interface {
	Checkout(ctx context.Context, data *model.Checkout) (*model.Payment, error)
}

type checkoutBiz struct {
	store CheckoutStore
}

func NewCheckoutBiz(store CheckoutStore) *checkoutBiz {
	return &checkoutBiz{store: store}
}

func (biz *checkoutBiz) Checkout(ctx context.Context, data *model.Checkout) (*model.Payment, error) {
	if data == nil {
		return nil, model.ErrCheckoutDataIsRequired
	}

	data.OrderID = strings.TrimSpace(data.OrderID)
	data.Method = strings.TrimSpace(strings.ToLower(data.Method))
	data.Provider = strings.TrimSpace(strings.ToLower(data.Provider))
	data.CreatedBy = strings.TrimSpace(data.CreatedBy)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	payment, err := biz.store.Checkout(ctx, data)
	if err != nil {
		return nil, err
	}

	if payment == nil {
		return nil, model.ErrPaymentNotFound
	}

	return payment, nil
}
