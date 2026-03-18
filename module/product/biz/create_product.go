package biz

import (
	"context"
	"strings"

	"stockflow/module/product/model"
)

type CreateProductStore interface {
	CreateProduct(ctx context.Context, data *model.Product) error
	FindProductBySKU(ctx context.Context, sku string) (*model.Product, error)
}

type createProductBiz struct {
	store CreateProductStore
}

func NewCreateProductBiz(store CreateProductStore) *createProductBiz {
	return &createProductBiz{store: store}
}

func (biz *createProductBiz) CreateProduct(ctx context.Context, data *model.ProductCreate) (*model.Product, error) {
	if data == nil {
		return nil, model.ErrProductDataIsNil
	}

	data.Name = strings.TrimSpace(data.Name)
	data.SKU = strings.TrimSpace(strings.ToUpper(data.SKU))
	data.Description = strings.TrimSpace(data.Description)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	existedProduct, err := biz.store.FindProductBySKU(ctx, data.SKU)

	if err != nil {
		return nil, err
	}

	if existedProduct != nil {
		return nil, model.ErrProductSKUAlreadyExists
	}

	product := &model.Product{
		SKU:         data.SKU,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		IsActive:    true,
	}

	if err := biz.store.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}
