package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"stockflow/module/outbox/model"

	"github.com/jackc/pgx/v5"
)

func (s *SQLStore) EnqueueEvent(ctx context.Context, data *model.OutboxEventCreate) (*model.OutboxEvent, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	payloadJSON, err := data.PayloadJSON()
	if err != nil {
		return nil, err
	}

	status := model.OutboxStatusPending
	now := time.Now()

	availableAt := data.AvailableAt
	if availableAt == nil {
		availableAt = &now
	}

	query := `
		INSERT INTO outbox_events (
			aggregate_type,
			aggregate_id,
			event_type,
			payload,
			status,
			retry_count,
			last_error,
			available_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			id,
			aggregate_type,
			aggregate_id,
			event_type,
			payload,
			status,
			retry_count,
			last_error,
			available_at,
			processed_at,
			created_at,
			updated_at;
	`

	var event model.OutboxEvent

	err = s.db.QueryRow(
		ctx,
		query,
		data.AggregateType,
		data.AggregateID,
		data.EventType,
		payloadJSON,
		status,
		0,
		"",
		availableAt,
	).Scan(
		&event.ID,
		&event.AggregateType,
		&event.AggregateID,
		&event.EventType,
		&event.Payload,
		&event.Status,
		&event.RetryCount,
		&event.LastError,
		&event.AvailableAt,
		&event.ProcessedAt,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot enqueue outbox event: %w", err)
	}

	return &event, nil
}

func (s *SQLStore) ListPendingEvents(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.OutboxEvent, error) {
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
			last_error,
			available_at,
			processed_at,
			created_at,
			updated_at
		FROM outbox_events
		WHERE 1=1
	`)

	if filter != nil {
		if filter.AggregateType != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND aggregate_type = $%d", argPos))
			args = append(args, strings.TrimSpace(strings.ToLower(filter.AggregateType)))
			argPos++
		}

		if filter.AggregateID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND aggregate_id = $%d", argPos))
			args = append(args, strings.TrimSpace(filter.AggregateID))
			argPos++
		}

		if filter.EventType != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND event_type = $%d", argPos))
			args = append(args, strings.TrimSpace(strings.ToLower(filter.EventType)))
			argPos++
		}

		if filter.Status != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", argPos))
			args = append(args, strings.TrimSpace(strings.ToLower(filter.Status)))
			argPos++
		}
	}

	queryBuilder.WriteString(fmt.Sprintf(" AND (available_at IS NULL OR available_at <= $%d)", argPos))
	args = append(args, time.Now())
	argPos++

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
		var event model.OutboxEvent

		if err := rows.Scan(
			&event.ID,
			&event.AggregateType,
			&event.AggregateID,
			&event.EventType,
			&event.Payload,
			&event.Status,
			&event.RetryCount,
			&event.LastError,
			&event.AvailableAt,
			&event.ProcessedAt,
			&event.CreatedAt,
			&event.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan outbox event: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate outbox event rows: %w", err)
	}

	return events, nil
}

func (s *SQLStore) GetOutboxEventByID(ctx context.Context, id string) (*model.OutboxEvent, error) {
	query := `
		SELECT
			id,
			aggregate_type,
			aggregate_id,
			event_type,
			payload,
			status,
			retry_count,
			last_error,
			available_at,
			processed_at,
			created_at,
			updated_at
		FROM outbox_events
		WHERE id = $1
		LIMIT 1;
	`

	var event model.OutboxEvent

	err := s.db.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.AggregateType,
		&event.AggregateID,
		&event.EventType,
		&event.Payload,
		&event.Status,
		&event.RetryCount,
		&event.LastError,
		&event.AvailableAt,
		&event.ProcessedAt,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get outbox event by id: %w", err)
	}

	return &event, nil
}
