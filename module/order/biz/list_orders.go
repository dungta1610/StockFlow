package biz

import (
	"context"

	"stockflow/module/order/model"
)

type ListOrdersStore interface {
	ListOrders(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.Order, error)
}

type listOrdersBiz struct {
	store ListOrdersStore
}

func NewListOrdersBiz(store ListOrdersStore) *listOrdersBiz {
	return &listOrdersBiz{store: store}
}

func (biz *listOrdersBiz) ListOrders(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.Order, error) {
	if filter != nil {
		filter.Normalize()
	}

	if paging == nil {
		paging = model.NewPaging()
	} else {
		paging.Normalize()
	}

	orders, err := biz.store.ListOrders(ctx, filter, paging)
	if err != nil {
		return nil, err
	}

	if orders == nil {
		return make([]model.Order, 0), nil
	}

	return orders, nil
}
