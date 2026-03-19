package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"stockflow/module/payment/model"

	"github.com/jackc/pgx/v5"
)

const (
	paymentStatusPending   = "pending"
	paymentStatusSuccess   = "success"
	paymentStatusFailed    = "failed"
	paymentStatusCancelled = "cancelled"
	paymentStatusRefunded  = "refunded"
)

func (s *SQLStore) CheckoutPayment(ctx context.Context, data *model.PaymentCheckout) (*model.Payment, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot begin checkout payment transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	existingQuery := `
		SELECT
			id,
			order_id,
			payment_code,
			method,
			status,
			amount,
			idempotency_key,
			external_txn_id,
			paid_at,
			failed_at,
			created_at,
			updated_at
		FROM payments
		WHERE idempotency_key = $1
		LIMIT 1
		FOR UPDATE;
	`

	var existing model.Payment
	err = tx.QueryRow(ctx, existingQuery, data.IdempotencyKey).Scan(
		&existing.ID,
		&existing.OrderID,
		&existing.PaymentCode,
		&existing.Method,
		&existing.Status,
		&existing.Amount,
		&existing.IdempotencyKey,
		&existing.ExternalTxnID,
		&existing.PaidAt,
		&existing.FailedAt,
		&existing.CreatedAt,
		&existing.UpdatedAt,
	)
	if err == nil {
		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("cannot commit existing checkout payment transaction: %w", err)
		}
		return &existing, nil
	}
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("cannot check existing payment by idempotency key: %w", err)
	}

	paymentCode, err := generatePaymentCode(ctx, tx)
	if err != nil {
		return nil, err
	}

	insertQuery := `
		INSERT INTO payments (
			order_id,
			payment_code,
			method,
			status,
			amount,
			idempotency_key,
			external_txn_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id,
			order_id,
			payment_code,
			method,
			status,
			amount,
			idempotency_key,
			external_txn_id,
			paid_at,
			failed_at,
			created_at,
			updated_at;
	`

	var payment model.Payment

	err = tx.QueryRow(
		ctx,
		insertQuery,
		data.OrderID,
		paymentCode,
		data.Method,
		paymentStatusPending,
		data.Amount,
		data.IdempotencyKey,
		nil,
	).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.PaymentCode,
		&payment.Method,
		&payment.Status,
		&payment.Amount,
		&payment.IdempotencyKey,
		&payment.ExternalTxnID,
		&payment.PaidAt,
		&payment.FailedAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot checkout payment: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit checkout payment transaction: %w", err)
	}

	return &payment, nil
}

func (s *SQLStore) CallbackPayment(ctx context.Context, data *model.PaymentCallback) (*model.Payment, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot begin callback payment transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var (
		lockQuery string
		lockArg   interface{}
		payment   model.Payment
	)

	if data.PaymentID != "" {
		lockQuery = `
			SELECT
				id,
				order_id,
				payment_code,
				method,
				status,
				amount,
				idempotency_key,
				external_txn_id,
				paid_at,
				failed_at,
				created_at,
				updated_at
			FROM payments
			WHERE id = $1
			LIMIT 1
			FOR UPDATE;
		`
		lockArg = data.PaymentID
	} else {
		lockQuery = `
			SELECT
				id,
				order_id,
				payment_code,
				method,
				status,
				amount,
				idempotency_key,
				external_txn_id,
				paid_at,
				failed_at,
				created_at,
				updated_at
			FROM payments
			WHERE payment_code = $1
			LIMIT 1
			FOR UPDATE;
		`
		lockArg = data.PaymentCode
	}

	err = tx.QueryRow(ctx, lockQuery, lockArg).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.PaymentCode,
		&payment.Method,
		&payment.Status,
		&payment.Amount,
		&payment.IdempotencyKey,
		&payment.ExternalTxnID,
		&payment.PaidAt,
		&payment.FailedAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot lock payment for callback: %w", err)
	}

	currentStatus := strings.TrimSpace(strings.ToLower(payment.Status))
	newStatus := strings.TrimSpace(strings.ToLower(data.Status))

	if currentStatus == paymentStatusSuccess {
		return nil, model.ErrPaymentAlreadyPaid
	}

	if currentStatus == paymentStatusFailed {
		return nil, model.ErrPaymentAlreadyFailed
	}

	if newStatus == "" {
		return nil, model.ErrPaymentStatusIsBlank
	}

	var paidAt *time.Time
	var failedAt *time.Time

	now := time.Now()

	switch newStatus {
	case paymentStatusSuccess:
		paidAt = &now
		failedAt = nil
	case paymentStatusFailed:
		failedAt = &now
		paidAt = nil
	default:
		paidAt = payment.PaidAt
		failedAt = payment.FailedAt
	}

	updateQuery := `
		UPDATE payments
		SET
			status = $1,
			external_txn_id = $2,
			paid_at = $3,
			failed_at = $4,
			updated_at = NOW()
		WHERE id = $5
		RETURNING
			id,
			order_id,
			payment_code,
			method,
			status,
			amount,
			idempotency_key,
			external_txn_id,
			paid_at,
			failed_at,
			created_at,
			updated_at;
	`

	err = tx.QueryRow(
		ctx,
		updateQuery,
		newStatus,
		data.ExternalTxnID,
		paidAt,
		failedAt,
		payment.ID,
	).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.PaymentCode,
		&payment.Method,
		&payment.Status,
		&payment.Amount,
		&payment.IdempotencyKey,
		&payment.ExternalTxnID,
		&payment.PaidAt,
		&payment.FailedAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot callback payment: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit callback payment transaction: %w", err)
	}

	return &payment, nil
}

func generatePaymentCode(ctx context.Context, tx pgx.Tx) (string, error) {
	query := `
		SELECT CONCAT('PAY-', TO_CHAR(NOW(), 'YYYYMMDD'), '-', LPAD((FLOOR(RANDOM() * 1000000))::text, 6, '0'));
	`

	var code string
	if err := tx.QueryRow(ctx, query).Scan(&code); err != nil {
		return "", fmt.Errorf("cannot generate payment code: %w", err)
	}

	return code, nil
}
