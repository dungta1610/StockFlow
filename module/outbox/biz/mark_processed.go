package biz

import (
	"context"
	"strings"

	"stockflow/module/outbox/model"
)

type MarkProcessedStore interface {
	MarkProcessed(ctx context.Context, data *model.OutboxEventMarkProcessed) (*model.OutboxEvent, error)
}

type markProcessedBiz struct {
	store MarkProcessedStore
}

func NewMarkProcessedBiz(store MarkProcessedStore) *markProcessedBiz {
	return &markProcessedBiz{store: store}
}

func (biz *markProcessedBiz) MarkProcessed(ctx context.Context, data *model.OutboxEventMarkProcessed) (*model.OutboxEvent, error) {
	if data == nil {
		return nil, model.ErrOutboxMarkProcessedDataIsRequired
	}

	data.EventID = strings.TrimSpace(data.EventID)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	event, err := biz.store.MarkProcessed(ctx, data)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, model.ErrOutboxEventNotFound
	}

	return event, nil
}
