package biz

import (
	"context"

	"stockflow/module/product/model"
)

type ListProductsStore interface {
	ListProducts(
		ctx context.Context,
		filter *model.Filter,
		paging *model.Paging,
	) ([]model.Product, error)
}

type listProductsBiz struct {
	store ListProductsStore
}

func NewListProductsBiz(store ListProductsStore) *listProductsBiz {
	return &listProductsBiz{store: store}
}

func (biz *listProductsBiz) ListProducts(
	ctx context.Context,
	filter *model.Filter,
	paging *model.Paging,
) ([]model.Product, error) {
	if filter != nil {
		filter.Normalize()
	}

	if paging == nil {
		paging = model.NewPaging()
	} else {
		paging.Normalize()
	}

	products, err := biz.store.ListProducts(ctx, filter, paging)
	if err != nil {
		return nil, err
	}

	if products == nil {
		return []model.Product{}, nil
	}

	return products, nil
}
