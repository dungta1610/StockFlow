package biz

import (
	"context"
	"strings"

	"stockflow/module/order/model"
)

type CancelOrderStore interface {
	CancelOrder(ctx context.Context, data *model.OrderCancel) (*model.Order, error)
}

type cancelOrderBiz struct {
	store CancelOrderStore
}

func NewCancelOrderBiz(store CancelOrderStore) *cancelOrderBiz {
	return &cancelOrderBiz{store: store}
}

func (biz *cancelOrderBiz) CancelOrder(ctx context.Context, data *model.OrderCancel) (*model.Order, error) {
	if data == nil {
		return nil, model.ErrOrderCancelDataIsRequired
	}

	data.OrderID = strings.TrimSpace(data.OrderID)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	updatedOrder, err := biz.store.CancelOrder(ctx, data)
	if err != nil {
		return nil, err
	}

	if updatedOrder == nil {
		return nil, model.ErrOrderNotFound
	}

	return updatedOrder, nil
}
