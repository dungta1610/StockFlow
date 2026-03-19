package biz

import (
	"context"

	"stockflow/module/payment/model"
)

type ListPaymentsStore interface {
	ListPayments(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.Payment, error)
}

type listPaymentsBiz struct {
	store ListPaymentsStore
}

func NewListPaymentsBiz(store ListPaymentsStore) *listPaymentsBiz {
	return &listPaymentsBiz{store: store}
}

func (biz *listPaymentsBiz) ListPayments(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.Payment, error) {
	if filter != nil {
		filter.Normalize()
	}

	if paging == nil {
		paging = model.NewPaging()
	} else {
		paging.Normalize()
	}

	payments, err := biz.store.ListPayments(ctx, filter, paging)
	if err != nil {
		return nil, err
	}

	if payments == nil {
		return make([]model.Payment, 0), nil
	}

	return payments, nil
}
