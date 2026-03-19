package model

import (
	"encoding/json"
	"strings"
	"time"
)

type OutboxEvent struct {
	ID            string          `json:"id" db:"id"`
	AggregateType string          `json:"aggregate_type" db:"aggregate_type"`
	AggregateID   string          `json:"aggregate_id" db:"aggregate_id"`
	EventType     string          `json:"event_type" db:"event_type"`
	Payload       json.RawMessage `json:"payload" db:"payload"`
	Status        string          `json:"status" db:"status"`
	RetryCount    int             `json:"retry_count" db:"retry_count"`
	NextRetryAt   *time.Time      `json:"next_retry_at,omitempty" db:"next_retry_at"`
	ErrorMessage  *string         `json:"error_message,omitempty" db:"error_message"`
	ProcessedAt   *time.Time      `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}

type OutboxEventCreate struct {
	AggregateType string          `json:"aggregate_type"`
	AggregateID   string          `json:"aggregate_id"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
}

func (o *OutboxEventCreate) Validate() error {
	if o == nil {
		return ErrOutboxEventCreateDataIsRequired
	}

	o.AggregateType = strings.TrimSpace(o.AggregateType)
	o.AggregateID = strings.TrimSpace(o.AggregateID)
	o.EventType = strings.TrimSpace(o.EventType)

	if o.AggregateType == "" {
		return ErrOutboxAggregateTypeIsBlank
	}

	if o.AggregateID == "" {
		return ErrOutboxAggregateIDIsBlank
	}

	if o.EventType == "" {
		return ErrOutboxEventTypeIsBlank
	}

	if len(o.Payload) == 0 {
		return ErrOutboxPayloadIsBlank
	}

	if !json.Valid(o.Payload) {
		return ErrOutboxPayloadInvalid
	}

	return nil
}

type OutboxEventMarkProcessed struct {
	EventID string `json:"event_id"`
}

func (o *OutboxEventMarkProcessed) Validate() error {
	if o == nil {
		return ErrOutboxMarkProcessedDataIsRequired
	}

	o.EventID = strings.TrimSpace(o.EventID)

	if o.EventID == "" {
		return ErrOutboxEventIDIsBlank
	}

	return nil
}

type OutboxEventMarkFailed struct {
	EventID      string     `json:"event_id"`
	ErrorMessage string     `json:"error_message"`
	NextRetryAt  *time.Time `json:"next_retry_at"`
}

func (o *OutboxEventMarkFailed) Validate() error {
	if o == nil {
		return ErrOutboxMarkFailedDataIsRequired
	}

	o.EventID = strings.TrimSpace(o.EventID)
	o.ErrorMessage = strings.TrimSpace(o.ErrorMessage)

	if o.EventID == "" {
		return ErrOutboxEventIDIsBlank
	}

	if o.ErrorMessage == "" {
		return ErrOutboxErrorMessageIsBlank
	}

	return nil
}

type Filter struct {
	AggregateType string `json:"aggregate_type" form:"aggregate_type"`
	AggregateID   string `json:"aggregate_id" form:"aggregate_id"`
	EventType     string `json:"event_type" form:"event_type"`
	Status        string `json:"status" form:"status"`
}

func (f *Filter) Normalize() {
	if f == nil {
		return
	}

	f.AggregateType = strings.TrimSpace(f.AggregateType)
	f.AggregateID = strings.TrimSpace(f.AggregateID)
	f.EventType = strings.TrimSpace(f.EventType)
	f.Status = strings.TrimSpace(f.Status)
}
