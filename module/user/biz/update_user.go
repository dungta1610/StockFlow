package biz

import (
	"context"
	"strings"

	"stockflow/module/user/model"
)

type UpdateUserStore interface {
	UpdateUser(ctx context.Context, id string, data *model.UserUpdate) (*model.User, error)
}

type updateUserBiz struct {
	store UpdateUserStore
}

func NewUpdateUserBiz(store UpdateUserStore) *updateUserBiz {
	return &updateUserBiz{store: store}
}

func (biz *updateUserBiz) UpdateUser(ctx context.Context, id string, data *model.UserUpdate) (*model.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, model.ErrUserIDIsBlank
	}

	if data == nil {
		return nil, model.ErrUserUpdateDataIsRequired
	}

	data.FullName = strings.TrimSpace(data.FullName)
	data.Role = strings.TrimSpace(data.Role)

	if data.PasswordHash != nil {
		trimmed := strings.TrimSpace(*data.PasswordHash)
		data.PasswordHash = &trimmed
	}

	if err := data.Validate(); err != nil {
		return nil, err
	}

	updatedUser, err := biz.store.UpdateUser(ctx, id, data)
	if err != nil {
		return nil, err
	}

	if updatedUser == nil {
		return nil, model.ErrUserNotFound
	}

	return updatedUser, nil
}
