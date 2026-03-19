package storage

import (
	"context"
	"fmt"
	"strings"

	"stockflow/module/payment/model"

	"github.com/jackc/pgx/v5"
)

func (s *SQLStore) GetPaymentByID(ctx context.Context, id string) (*model.Payment, error) {
	query := `
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
		LIMIT 1;
	`

	var payment model.Payment

	err := s.db.QueryRow(ctx, query, id).Scan(
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
		return nil, fmt.Errorf("cannot get payment by id: %w", err)
	}

	return &payment, nil
}

func (s *SQLStore) ListPayments(
	ctx context.Context,
	filter *model.Filter,
	paging *model.Paging,
) ([]model.Payment, error) {
	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0)
	argPos := 1

	queryBuilder.WriteString(`
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
		WHERE 1=1
	`)

	if filter != nil {
		if filter.OrderID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND order_id = $%d", argPos))
			args = append(args, filter.OrderID)
			argPos++
		}

		if filter.PaymentCode != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND payment_code = $%d", argPos))
			args = append(args, filter.PaymentCode)
			argPos++
		}

		if filter.Method != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND method = $%d", argPos))
			args = append(args, filter.Method)
			argPos++
		}

		if filter.Status != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", argPos))
			args = append(args, filter.Status)
			argPos++
		}
	}

	queryBuilder.WriteString(" ORDER BY created_at DESC")

	if paging != nil {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1))
		args = append(args, paging.Limit, paging.Offset())
	}

	rows, err := s.db.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("cannot list payments: %w", err)
	}
	defer rows.Close()

	payments := make([]model.Payment, 0)

	for rows.Next() {
		var payment model.Payment

		if err := rows.Scan(
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
		); err != nil {
			return nil, fmt.Errorf("cannot scan payment: %w", err)
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate payment rows: %w", err)
	}

	return payments, nil
}
