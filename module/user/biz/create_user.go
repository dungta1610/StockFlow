package biz

import (
	"context"
	"strings"

	"stockflow/module/user/model"
)

type CreateUserStore interface {
	CreateUser(ctx context.Context, data *model.UserCreate) (*model.User, error)
}

type createUserBiz struct {
	store CreateUserStore
}

func NewCreateUserBiz(store CreateUserStore) *createUserBiz {
	return &createUserBiz{store: store}
}

func (biz *createUserBiz) CreateUser(ctx context.Context, data *model.UserCreate) (*model.User, error) {
	if data == nil {
		return nil, model.ErrUserDataIsRequired
	}

	data.Email = strings.TrimSpace(strings.ToLower(data.Email))
	data.PasswordHash = strings.TrimSpace(data.PasswordHash)
	data.FullName = strings.TrimSpace(data.FullName)
	data.Role = strings.TrimSpace(data.Role)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	createdUser, err := biz.store.CreateUser(ctx, data)
	if err != nil {
		return nil, err
	}

	if createdUser == nil {
		return nil, model.ErrUserNotFound
	}

	return createdUser, nil
}
