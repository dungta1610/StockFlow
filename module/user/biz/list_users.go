package biz

import (
	"context"

	"stockflow/module/user/model"
)

type ListUsersStore interface {
	ListUsers(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.User, error)
}

type listUsersBiz struct {
	store ListUsersStore
}

func NewListUsersBiz(store ListUsersStore) *listUsersBiz {
	return &listUsersBiz{store: store}
}

func (biz *listUsersBiz) ListUsers(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.User, error) {
	if filter != nil {
		filter.Normalize()
	}

	if paging == nil {
		paging = model.NewPaging()
	} else {
		paging.Normalize()
	}

	users, err := biz.store.ListUsers(ctx, filter, paging)
	if err != nil {
		return nil, err
	}

	if users == nil {
		return make([]model.User, 0), nil
	}

	return users, nil
}
