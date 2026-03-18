package biz

import (
	"context"
	"strings"

	"stockflow/module/order/model"
)

type CreateOrderStore interface {
	CreateOrder(ctx context.Context, data *model.OrderCreate) (*model.Order, error)
}

type createOrderBiz struct {
	store CreateOrderStore
}

func NewCreateOrderBiz(store CreateOrderStore) *createOrderBiz {
	return &createOrderBiz{store: store}
}

func (biz *createOrderBiz) CreateOrder(ctx context.Context, data *model.OrderCreate) (*model.Order, error) {
	if data == nil {
		return nil, model.ErrOrderDataIsRequired
	}

	data.UserID = strings.TrimSpace(data.UserID)
	data.WarehouseID = strings.TrimSpace(data.WarehouseID)
	data.Note = strings.TrimSpace(data.Note)
	data.CreatedBy = strings.TrimSpace(data.CreatedBy)

	for i := range data.Items {
		data.Items[i].ProductID = strings.TrimSpace(data.Items[i].ProductID)
	}

	if err := data.Validate(); err != nil {
		return nil, err
	}

	createdOrder, err := biz.store.CreateOrder(ctx, data)
	if err != nil {
		return nil, err
	}

	if createdOrder == nil {
		return nil, model.ErrOrderNotFound
	}

	return createdOrder, nil
}
