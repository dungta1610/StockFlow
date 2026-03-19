package storage

import (
	"context"
	"fmt"

	"stockflow/module/outbox/model"

	"github.com/jackc/pgx/v5"
)

func (s *SQLStore) EnqueueEvent(ctx context.Context, data *model.OutboxEventCreate) (*model.OutboxEvent, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO outbox_events (
			aggregate_type,
			aggregate_id,
			event_type,
			payload,
			status,
			retry_count,
			next_retry_at,
			error_message,
			processed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, NULL, NULL, NULL)
		RETURNING
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
			updated_at;
	`

	event, err := scanOutboxEvent(s.db.QueryRow(
		ctx,
		query,
		data.AggregateType,
		data.AggregateID,
		data.EventType,
		data.Payload,
		outboxStatusPending,
		0,
	))
	if err != nil {
		return nil, fmt.Errorf("cannot enqueue outbox event: %w", err)
	}

	return event, nil
}

func (s *SQLStore) MarkProcessed(ctx context.Context, data *model.OutboxEventMarkProcessed) (*model.OutboxEvent, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot begin mark processed transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	lockQuery := `
		SELECT
			id
		FROM outbox_events
		WHERE id = $1
		LIMIT 1
		FOR UPDATE;
	`

	var eventID string
	err = tx.QueryRow(ctx, lockQuery, data.EventID).Scan(&eventID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot lock outbox event for mark processed: %w", err)
	}

	updateQuery := `
		UPDATE outbox_events
		SET
			status = $1,
			processed_at = NOW(),
			error_message = NULL,
			next_retry_at = NULL,
			updated_at = NOW()
		WHERE id = $2
		RETURNING
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
			updated_at;
	`

	event, err := scanOutboxEvent(tx.QueryRow(ctx, updateQuery, outboxStatusProcessed, data.EventID))
	if err != nil {
		return nil, fmt.Errorf("cannot mark outbox event as processed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit mark processed transaction: %w", err)
	}

	return event, nil
}

func (s *SQLStore) MarkFailed(ctx context.Context, data *model.OutboxEventMarkFailed) (*model.OutboxEvent, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot begin mark failed transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	lockQuery := `
		SELECT
			id
		FROM outbox_events
		WHERE id = $1
		LIMIT 1
		FOR UPDATE;
	`

	var eventID string
	err = tx.QueryRow(ctx, lockQuery, data.EventID).Scan(&eventID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot lock outbox event for mark failed: %w", err)
	}

	updateQuery := `
		UPDATE outbox_events
		SET
			status = $1,
			retry_count = retry_count + 1,
			error_message = $2,
			next_retry_at = $3,
			processed_at = NULL,
			updated_at = NOW()
		WHERE id = $4
		RETURNING
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
			updated_at;
	`

	event, err := scanOutboxEvent(
		tx.QueryRow(ctx, updateQuery, outboxStatusFailed, data.ErrorMessage, data.NextRetryAt, data.EventID),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot mark outbox event as failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit mark failed transaction: %w", err)
	}

	return event, nil
}
