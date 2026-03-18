package storage

import (
	"context"
	"fmt"

	"stockflow/module/payment/model"

	"github.com/jackc/pgx/v5"
)

func (s *SQLStore) GetPaymentByID(ctx context.Context, id string) (*model.Payment, error) {
	query := `
		SELECT
			id,
			order_id,
			provider,
			method,
			status,
			amount,
			currency,
			provider_txn_id,
			provider_order_code,
			checkout_url,
			callback_payload,
			failure_reason,
			paid_at,
			expired_at,
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
		&payment.Provider,
		&payment.Method,
		&payment.Status,
		&payment.Amount,
		&payment.Currency,
		&payment.ProviderTxnID,
		&payment.ProviderOrderCode,
		&payment.CheckoutURL,
		&payment.CallbackPayload,
		&payment.FailureReason,
		&payment.PaidAt,
		&payment.ExpiredAt,
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
