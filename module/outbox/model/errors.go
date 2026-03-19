package model

import "errors"

var (
	ErrOutboxEventDataIsRequired         = errors.New("outbox event data is required")
	ErrOutboxEventCreateDataIsRequired   = errors.New("outbox event create data is required")
	ErrOutboxMarkProcessedDataIsRequired = errors.New("outbox mark processed data is required")
	ErrOutboxMarkFailedDataIsRequired    = errors.New("outbox mark failed data is required")

	ErrOutboxEventIDIsBlank       = errors.New("outbox event id is required")
	ErrOutboxAggregateTypeIsBlank = errors.New("outbox aggregate type is required")
	ErrOutboxAggregateIDIsBlank   = errors.New("outbox aggregate id is required")
	ErrOutboxEventTypeIsBlank     = errors.New("outbox event type is required")
	ErrOutboxPayloadIsBlank       = errors.New("outbox payload is required")
	ErrOutboxPayloadInvalid       = errors.New("outbox payload must be valid json")
	ErrOutboxStatusIsBlank        = errors.New("outbox status is required")
	ErrOutboxRetryCountInvalid    = errors.New("outbox retry count must be greater than or equal to 0")
	ErrOutboxErrorMessageIsBlank  = errors.New("outbox error message is required")

	ErrOutboxEventNotFound = errors.New("outbox event not found")
)
