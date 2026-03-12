package biz

import (
	"context"
	"stockflow/module/warehouse/model"
	"strings"
)

type CreateWarehouseStore interface {
	CreateWarehouse(ctx context.Context, data *model.Warehouse) error
	FindWarehouseByCode(ctx context.Context, code string) (*model.Warehouse, error)
}

type createWarehouseBiz struct {
	store CreateWarehouseStore
}

func NewCreateWarehouseBiz(store CreateWarehouseStore) *createWarehouseBiz {
	return &createWarehouseBiz{store: store}
}

func (biz *createWarehouseBiz) CreateWarehouse(ctx context.Context, data *model.WarehouseCreate) (*model.Warehouse, error) {
	if data == nil {
		return nil, model.ErrWarehouseDataIsNil
	}

	data.Code = strings.TrimSpace(strings.ToUpper(data.Code))
	data.Name = strings.TrimSpace(data.Name)
	data.Address = strings.TrimSpace(data.Address)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	exitedWarehouse, err := biz.store.FindWarehouseByCode(ctx, data.Code)

	if err != nil {
		return nil, err
	}

	if exitedWarehouse != nil {
		return nil, model.ErrWarehouseCodeAlreadyExists
	}

	warehouse := &model.Warehouse{
		Code:     data.Code,
		Name:     data.Name,
		Address:  data.Address,
		IsActive: true,
	}

	if err := biz.store.CreateWarehouse(ctx, warehouse); err != nil {
		return nil, err
	}

	return warehouse, nil
}
