package storage

import (
	"context"
	"fmt"
	"time"

	"stockflow/module/outbox/model"

	"github.com/jackc/pgx/v5"
)

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
		FOR UPDATE;
	`

	var event model.OutboxEvent

	err = tx.QueryRow(ctx, lockQuery, data.ID).Scan(
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
		return nil, fmt.Errorf("cannot lock outbox event for mark processed: %w", err)
	}

	if event.Status == model.OutboxStatusProcessed {
		return nil, model.ErrOutboxEventAlreadyProcessed
	}

	if event.Status != model.OutboxStatusPending && event.Status != model.OutboxStatusProcessing {
		return nil, model.ErrOutboxEventCannotBeProcessed
	}

	now := time.Now()

	updateQuery := `
		UPDATE outbox_events
		SET
			status = $1,
			processed_at = $2,
			last_error = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at;
	`

	err = tx.QueryRow(
		ctx,
		updateQuery,
		model.OutboxStatusProcessed,
		now,
		"",
		event.ID,
	).Scan(&event.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot mark outbox event as processed: %w", err)
	}

	event.Status = model.OutboxStatusProcessed
	event.ProcessedAt = &now
	event.LastError = ""

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit mark processed transaction: %w", err)
	}

	return &event, nil
}
