package biz

import (
	"context"
	"strings"

	"stockflow/module/payment/model"
)

type GetPaymentStore interface {
	GetPaymentByID(ctx context.Context, id string) (*model.Payment, error)
}

type getPaymentBiz struct {
	store GetPaymentStore
}

func NewGetPaymentBiz(store GetPaymentStore) *getPaymentBiz {
	return &getPaymentBiz{store: store}
}

func (biz *getPaymentBiz) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, model.ErrPaymentIDIsBlank
	}

	payment, err := biz.store.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if payment == nil {
		return nil, model.ErrPaymentNotFound
	}

	return payment, nil
}
