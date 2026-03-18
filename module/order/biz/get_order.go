package biz

import (
	"context"
	"strings"

	"stockflow/module/order/model"
)

type GetOrderStore interface {
	GetOrderByID(ctx context.Context, id string) (*model.Order, error)
}

type getOrderBiz struct {
	store GetOrderStore
}

func NewGetOrderBiz(store GetOrderStore) *getOrderBiz {
	return &getOrderBiz{store: store}
}

func (biz *getOrderBiz) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, model.ErrOrderIDIsBlank
	}

	order, err := biz.store.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, model.ErrOrderNotFound
	}

	return order, nil
}
