package biz

import (
	"context"

	"stockflow/module/outbox/model"
)

type ListPendingEventsStore interface {
	ListPendingEvents(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.OutboxEvent, error)
}

type listPendingEventsBiz struct {
	store ListPendingEventsStore
}

func NewListPendingEventsBiz(store ListPendingEventsStore) *listPendingEventsBiz {
	return &listPendingEventsBiz{store: store}
}

func (biz *listPendingEventsBiz) ListPendingEvents(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.OutboxEvent, error) {
	if filter == nil {
		filter = &model.Filter{}
	}
	filter.Normalize()

	if filter.Status == "" {
		filter.Status = model.OutboxStatusPending
	}

	if paging == nil {
		paging = model.NewPaging()
	} else {
		paging.Normalize()
	}

	events, err := biz.store.ListPendingEvents(ctx, filter, paging)
	if err != nil {
		return nil, err
	}

	if events == nil {
		return make([]model.OutboxEvent, 0), nil
	}

	return events, nil
}
