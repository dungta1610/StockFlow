package biz

import (
	"context"

	"stockflow/module/product/model"
)

type GetProductStore interface {
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
}

type getProductBiz struct {
	store GetProductStore
}

func NewGetProductBiz(store GetProductStore) *getProductBiz {
	return &getProductBiz{store: store}
}

func (biz *getProductBiz) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	if id == "" {
		return nil, model.ErrProductIDIsBlank
	}

	product, err := biz.store.GetProductByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, model.ErrProductNotFound
	}

	return product, nil
}
