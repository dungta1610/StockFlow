package biz

import (
	"context"
	"strings"

	"stockflow/module/outbox/model"
)

type EnqueueEventStore interface {
	EnqueueEvent(ctx context.Context, data *model.OutboxEventCreate) (*model.OutboxEvent, error)
}

type enqueueEventBiz struct {
	store EnqueueEventStore
}

func NewEnqueueEventBiz(store EnqueueEventStore) *enqueueEventBiz {
	return &enqueueEventBiz{store: store}
}

func (biz *enqueueEventBiz) EnqueueEvent(ctx context.Context, data *model.OutboxEventCreate) (*model.OutboxEvent, error) {
	if data == nil {
		return nil, model.ErrOutboxEventCreateDataIsRequired
	}

	data.AggregateType = strings.TrimSpace(data.AggregateType)
	data.AggregateID = strings.TrimSpace(data.AggregateID)
	data.EventType = strings.TrimSpace(data.EventType)

	if err := data.Validate(); err != nil {
		return nil, err
	}

	event, err := biz.store.EnqueueEvent(ctx, data)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, model.ErrOutboxEventNotFound
	}

	return event, nil
}
