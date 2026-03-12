package biz

import (
	"context"

	"stockflow/module/inventory/model"
)

type ListInventoryTransactionsStore interface {
	ListInventoryTransactions(ctx context.Context, filter *model.TransactionFilter, paging *model.Paging) ([]model.InventoryTransaction, error)
}

type listInventoryTransactionsBiz struct {
	store ListInventoryTransactionsStore
}

func NewListInventoryTransactionsBiz(store ListInventoryTransactionsStore) *listInventoryTransactionsBiz {
	return &listInventoryTransactionsBiz{store: store}
}

func (biz *listInventoryTransactionsBiz) ListInventoryTransactions(ctx context.Context, filter *model.TransactionFilter, paging *model.Paging) ([]model.InventoryTransaction, error) {
	if filter != nil {
		filter.Normalize()
	}

	if paging == nil {
		paging = model.NewPaging()
	} else {
		paging.Normalize()
	}

	transactions, err := biz.store.ListInventoryTransactions(ctx, filter, paging)

	if err != nil {
		return nil, err
	}

	if transactions == nil {
		return []model.InventoryTransaction{}, nil
	}

	return transactions, nil
}
