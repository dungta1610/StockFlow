package model

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	OutboxStatusPending    = "pending"
	OutboxStatusProcessing = "processing"
	OutboxStatusProcessed  = "processed"
	OutboxStatusFailed     = "failed"
)

type OutboxEvent struct {
	ID            string     `json:"id" db:"id"`
	AggregateType string     `json:"aggregate_type" db:"aggregate_type"`
	AggregateID   string     `json:"aggregate_id" db:"aggregate_id"`
	EventType     string     `json:"event_type" db:"event_type"`
	Payload       string     `json:"payload" db:"payload"`
	Status        string     `json:"status" db:"status"`
	RetryCount    int        `json:"retry_count" db:"retry_count"`
	LastError     string     `json:"last_error" db:"last_error"`
	AvailableAt   *time.Time `json:"available_at,omitempty" db:"available_at"`
	ProcessedAt   *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type OutboxEventCreate struct {
	AggregateType string      `json:"aggregate_type"`
	AggregateID   string      `json:"aggregate_id"`
	EventType     string      `json:"event_type"`
	Payload       interface{} `json:"payload"`
	AvailableAt   *time.Time  `json:"available_at"`
}

func (d *OutboxEventCreate) Validate() error {
	if d == nil {
		return ErrOutboxEventDataIsRequired
	}

	d.AggregateType = strings.TrimSpace(strings.ToLower(d.AggregateType))
	d.AggregateID = strings.TrimSpace(d.AggregateID)
	d.EventType = strings.TrimSpace(strings.ToLower(d.EventType))

	if d.AggregateType == "" {
		return ErrOutboxAggregateTypeIsBlank
	}

	if d.AggregateID == "" {
		return ErrOutboxAggregateIDIsBlank
	}

	if d.EventType == "" {
		return ErrOutboxEventTypeIsBlank
	}

	if d.Payload == nil {
		return ErrOutboxPayloadIsRequired
	}

	return nil
}

func (d *OutboxEventCreate) PayloadJSON() (string, error) {
	if err := d.Validate(); err != nil {
		return "", err
	}

	raw, err := json.Marshal(d.Payload)
	if err != nil {
		return "", ErrOutboxPayloadInvalid
	}

	return string(raw), nil
}

type OutboxEventMarkProcessed struct {
	ID string `json:"id"`
}

func (d *OutboxEventMarkProcessed) Validate() error {
	if d == nil {
		return ErrOutboxMarkProcessedDataIsRequired
	}

	d.ID = strings.TrimSpace(d.ID)
	if d.ID == "" {
		return ErrOutboxEventIDIsBlank
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

	f.AggregateType = strings.TrimSpace(strings.ToLower(f.AggregateType))
	f.AggregateID = strings.TrimSpace(f.AggregateID)
	f.EventType = strings.TrimSpace(strings.ToLower(f.EventType))
	f.Status = strings.TrimSpace(strings.ToLower(f.Status))
}
