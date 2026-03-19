package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"stockflow/module/outbox/model"
)

const (
	outboxStatusPending    = "pending"
	outboxStatusProcessing = "processing"
	outboxStatusProcessed  = "processed"
	outboxStatusFailed     = "failed"
)

type rowScanner interface {
	Scan(dest ...any) error
}

func scanOutboxEvent(scanner rowScanner) (*model.OutboxEvent, error) {
	var (
		event        model.OutboxEvent
		payloadBytes []byte
	)

	err := scanner.Scan(
		&event.ID,
		&event.AggregateType,
		&event.AggregateID,
		&event.EventType,
		&payloadBytes,
		&event.Status,
		&event.RetryCount,
		&event.NextRetryAt,
		&event.ErrorMessage,
		&event.ProcessedAt,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	event.Payload = json.RawMessage(payloadBytes)

	return &event, nil
}

func (s *SQLStore) ListPendingEvents(
	ctx context.Context,
	filter *model.Filter,
	paging *model.Paging,
) ([]model.OutboxEvent, error) {
	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0)
	argPos := 1

	queryBuilder.WriteString(`
		SELECT
			id,
			aggregate_type,
			aggregate_id,
			event_type,
			payload,
			status,
			retry_count,
			next_retry_at,
			error_message,
			processed_at,
			created_at,
			updated_at
		FROM outbox_events
		WHERE status IN ('pending', 'failed')
		  AND processed_at IS NULL
		  AND (next_retry_at IS NULL OR next_retry_at <= NOW())
	`)

	if filter != nil {
		if filter.AggregateType != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND aggregate_type = $%d", argPos))
			args = append(args, filter.AggregateType)
			argPos++
		}

		if filter.AggregateID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND aggregate_id = $%d", argPos))
			args = append(args, filter.AggregateID)
			argPos++
		}

		if filter.EventType != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND event_type = $%d", argPos))
			args = append(args, filter.EventType)
			argPos++
		}

		if filter.Status != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", argPos))
			args = append(args, filter.Status)
			argPos++
		}
	}

	queryBuilder.WriteString(" ORDER BY created_at ASC")

	if paging != nil {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1))
		args = append(args, paging.Limit, paging.Offset())
	}

	rows, err := s.db.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("cannot list pending outbox events: %w", err)
	}
	defer rows.Close()

	events := make([]model.OutboxEvent, 0)

	for rows.Next() {
		event, err := scanOutboxEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("cannot scan outbox event: %w", err)
		}

		events = append(events, *event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate outbox event rows: %w", err)
	}

	return events, nil
}
