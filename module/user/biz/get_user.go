package biz

import (
	"context"
	"strings"

	"stockflow/module/user/model"
)

type GetUserStore interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

type getUserBiz struct {
	store GetUserStore
}

func NewGetUserBiz(store GetUserStore) *getUserBiz {
	return &getUserBiz{store: store}
}

func (biz *getUserBiz) GetUser(ctx context.Context, id string) (*model.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, model.ErrUserIDIsBlank
	}

	user, err := biz.store.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, model.ErrUserNotFound
	}

	return user, nil
}
