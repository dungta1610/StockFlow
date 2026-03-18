package model

import "errors"

var (
	ErrOutboxEventDataIsRequired         = errors.New("outbox event data is required")
	ErrOutboxEventIDIsBlank              = errors.New("outbox event id is required")
	ErrOutboxAggregateTypeIsBlank        = errors.New("outbox aggregate type is required")
	ErrOutboxAggregateIDIsBlank          = errors.New("outbox aggregate id is required")
	ErrOutboxEventTypeIsBlank            = errors.New("outbox event type is required")
	ErrOutboxPayloadIsRequired           = errors.New("outbox payload is required")
	ErrOutboxPayloadInvalid              = errors.New("outbox payload is invalid")
	ErrOutboxStatusIsBlank               = errors.New("outbox status is required")
	ErrOutboxEventNotFound               = errors.New("outbox event not found")
	ErrOutboxEventAlreadyProcessed       = errors.New("outbox event already processed")
	ErrOutboxEventCannotBeProcessed      = errors.New("outbox event cannot be processed")
	ErrOutboxMarkProcessedDataIsRequired = errors.New("outbox mark processed data is required")
)
