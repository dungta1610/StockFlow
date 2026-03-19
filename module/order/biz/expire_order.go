package biz

import (
	"context"
	"strings"

	"stockflow/module/order/model"
)

type ExpireOrderStore interface {
	ExpireOrder(ctx context.Context, data *model.OrderExpire) (*model.Order, error)
}

type expireOrderBiz struct {
	store ExpireOrderStore
}

func NewExpireOrderBiz(store ExpireOrderStore) *expireOrderBiz {
	return &expireOrderBiz{store: store}
}

func (biz *expireOrderBiz) ExpireOrder(ctx context.Context, data *model.OrderExpire) (*model.Order, error) {
	if data == nil {
		return nil, model.ErrOrderExpireDataIsRequired
	}

	data.OrderID = strings.TrimSpace(data.OrderID)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	updatedOrder, err := biz.store.ExpireOrder(ctx, data)
	if err != nil {
		return nil, err
	}

	if updatedOrder == nil {
		return nil, model.ErrOrderNotFound
	}

	return updatedOrder, nil
}
