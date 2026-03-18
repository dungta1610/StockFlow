package biz

import (
	"context"
	"strings"

	"stockflow/module/inventory/model"
)

type AdjustStockStore interface {
	AdjustStock(ctx context.Context, data *model.InventoryAdjust) (*model.Inventory, error)
}

type adjustStockBiz struct {
	store AdjustStockStore
}

func NewAdjustStockBiz(store AdjustStockStore) *adjustStockBiz {
	return &adjustStockBiz{store: store}
}

func (biz *adjustStockBiz) AdjustStock(ctx context.Context, data *model.InventoryAdjust) (*model.Inventory, error) {
	if data == nil {
		return nil, model.ErrInventoryAdjustDataRequired
	}

	data.ProductID = strings.TrimSpace(data.ProductID)
	data.WarehouseID = strings.TrimSpace(data.WarehouseID)
	data.Reason = strings.TrimSpace(data.Reason)
	data.CreatedBy = strings.TrimSpace(data.CreatedBy)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	updatedInventory, err := biz.store.AdjustStock(ctx, data)
	if err != nil {
		return nil, err
	}

	if updatedInventory == nil {
		return nil, model.ErrInventoryNotFound
	}

	return updatedInventory, nil
}
