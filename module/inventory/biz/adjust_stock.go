package biz

import (
	"context"
	"strings"

	"stockflow/module/inventory/model"
)

type AdjustStockStore interface {
	GetInventoryByProductAndWarehouse(ctx context.Context, productID, warehouseID string) (*model.Inventory, error)
	CreateInventory(ctx context.Context, data *model.Inventory) error
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

	inventory, err := biz.store.GetInventoryByProductAndWarehouse(
		ctx,
		data.ProductID,
		data.WarehouseID,
	)

	if err != nil {
		return nil, err
	}

	if inventory == nil {
		if data.Quantity < 0 {
			return nil, model.ErrInventoryNotFound
		}

		newInventory := &model.Inventory{
			ProductID:    data.ProductID,
			WarehouseID:  data.WarehouseID,
			AvailableQty: data.Quantity,
			ReservedQty:  0,
			Version:      1,
		}

		if err := biz.store.CreateInventory(ctx, newInventory); err != nil {
			return nil, err
		}

		return newInventory, nil
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
