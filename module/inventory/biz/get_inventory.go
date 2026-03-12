package biz

import (
	"context"
	"strings"

	"stockflow/module/inventory/model"
)

type GetInventoryStore interface {
	GetInventoryByProductAndWarehouse(
		ctx context.Context,
		productID, warehouseID string,
	) (*model.Inventory, error)

	GetInventoryByID(ctx context.Context, id string) (*model.Inventory, error)
}

type getInventoryBiz struct {
	store GetInventoryStore
}

func NewGetInventoryBiz(store GetInventoryStore) *getInventoryBiz {
	return &getInventoryBiz{store: store}
}

func (biz *getInventoryBiz) GetInventory(ctx context.Context, id, productID, warehouseID string) (*model.Inventory, error) {
	id = strings.TrimSpace(id)
	productID = strings.TrimSpace(productID)
	warehouseID = strings.TrimSpace(warehouseID)

	if id != "" {
		inventory, err := biz.store.GetInventoryByID(ctx, id)

		if err != nil {
			return nil, err
		}

		if inventory == nil {
			return nil, model.ErrInventoryNotFound
		}

		return inventory, nil
	}

	if productID == "" {
		return nil, model.ErrInventoryProductIDIsBlank
	}

	if warehouseID == "" {
		return nil, model.ErrInventoryWarehouseIDIsBlank
	}

	inventory, err := biz.store.GetInventoryByProductAndWarehouse(ctx, productID, warehouseID)

	if err != nil {
		return nil, err
	}

	if inventory == nil {
		return nil, model.ErrInventoryNotFound
	}

	return inventory, nil
}
