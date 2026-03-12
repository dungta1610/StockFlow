package biz

import (
	"context"
	"stockflow/module/warehouse/model"
)

type GetWarehouseStore interface {
	GetWarehouseByID(ctx context.Context, id string) (*model.Warehouse, error)
}

type getWarehouseBiz struct {
	store GetWarehouseStore
}

func NewGetWarehouseBiz(store GetWarehouseStore) *getWarehouseBiz {
	return &getWarehouseBiz{store: store}
}

func (biz *getWarehouseBiz) GetWarehouse(ctx context.Context, id string) (*model.Warehouse, error) {
	if id == "" {
		return nil, model.ErrWarehouseIDIsBlank
	}

	warehouse, err := biz.store.GetWarehouseByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if warehouse == nil {
		return nil, model.ErrWarehouseNotFound
	}

	return warehouse, nil
}
