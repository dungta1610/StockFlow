package biz

import (
	"context"
	"strings"

	"stockflow/module/outbox/model"
)

type MarkFailedStore interface {
	MarkFailed(ctx context.Context, data *model.OutboxEventMarkFailed) (*model.OutboxEvent, error)
}

type markFailedBiz struct {
	store MarkFailedStore
}

func NewMarkFailedBiz(store MarkFailedStore) *markFailedBiz {
	return &markFailedBiz{store: store}
}

func (biz *markFailedBiz) MarkFailed(ctx context.Context, data *model.OutboxEventMarkFailed) (*model.OutboxEvent, error) {
	if data == nil {
		return nil, model.ErrOutboxMarkFailedDataIsRequired
	}

	data.EventID = strings.TrimSpace(data.EventID)
	data.ErrorMessage = strings.TrimSpace(data.ErrorMessage)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	event, err := biz.store.MarkFailed(ctx, data)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, model.ErrOutboxEventNotFound
	}

	return event, nil
}
