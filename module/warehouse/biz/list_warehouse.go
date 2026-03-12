package biz

import (
	"context"
	"stockflow/module/warehouse/model"
)

type ListWarehousesStore interface {
	ListWarehouses(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.Warehouse, error)
}

type listWarehousesBiz struct {
	store ListWarehousesStore
}

func NewListWarehousesBiz(store ListWarehousesStore) *listWarehousesBiz {
	return &listWarehousesBiz{store: store}
}

func (biz *listWarehousesBiz) ListWarehouses(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.Warehouse, error) {
	if filter != nil {
		filter.Normalize()
	}

	if paging != nil {
		paging = model.NewPaging()
	} else {
		paging.Normalize()
	}

	warehouses, err := biz.store.ListWarehouses(ctx, filter, paging)

	if err != nil {
		return nil, err
	}

	if warehouses == nil {
		return []model.Warehouse{}, nil
	}

	return warehouses, nil
}
