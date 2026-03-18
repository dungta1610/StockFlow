package storage

import (
	"context"
	"fmt"
	"time"

	"stockflow/module/payment/model"

	"github.com/jackc/pgx/v5"
)

type orderPaymentSnapshot struct {
	ID            string
	Status        string
	PaymentStatus string
	TotalPrice    float64
}

func (s *SQLStore) Checkout(ctx context.Context, data *model.Checkout) (*model.Payment, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot begin checkout transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	order, err := getOrderForPayment(ctx, tx, data.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, model.ErrPaymentNotFound
	}

	if order.PaymentStatus == "paid" {
		return nil, model.ErrPaymentAlreadySucceeded
	}

	existingPendingPayment, err := getLatestPendingPaymentByOrderID(ctx, tx, data.OrderID)
	if err != nil {
		return nil, err
	}
	if existingPendingPayment != nil {
		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("cannot commit checkout transaction: %w", err)
		}
		return existingPendingPayment, nil
	}

	expiredAt := time.Now().Add(15 * time.Minute)
	providerOrderCode := buildProviderOrderCode(data.OrderID)

	payment := &model.Payment{
		OrderID:           data.OrderID,
		Provider:          data.Provider,
		Method:            data.Method,
		Status:            model.PaymentStatusPending,
		Amount:            order.TotalPrice,
		Currency:          "VND",
		ProviderOrderCode: providerOrderCode,
		CheckoutURL:       buildCheckoutURL(data.Provider, providerOrderCode),
		ExpiredAt:         &expiredAt,
	}

	insertQuery := `
		INSERT INTO payments (
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
			expired_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING
			id,
			created_at,
			updated_at;
	`

	err = tx.QueryRow(
		ctx,
		insertQuery,
		payment.OrderID,
		payment.Provider,
		payment.Method,
		payment.Status,
		payment.Amount,
		payment.Currency,
		payment.ProviderTxnID,
		payment.ProviderOrderCode,
		payment.CheckoutURL,
		payment.CallbackPayload,
		payment.FailureReason,
		payment.PaidAt,
		payment.ExpiredAt,
	).Scan(
		&payment.ID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create payment: %w", err)
	}

	updateOrderQuery := `
		UPDATE orders
		SET
			payment_status = $1,
			updated_at = NOW()
		WHERE id = $2;
	`

	if _, err := tx.Exec(ctx, updateOrderQuery, "pending", order.ID); err != nil {
		return nil, fmt.Errorf("cannot update order payment status after checkout: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit checkout transaction: %w", err)
	}

	return payment, nil
}

func (s *SQLStore) HandleCallback(ctx context.Context, data *model.Callback) (*model.Payment, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot begin callback transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	payment, err := getPaymentForUpdate(ctx, tx, data.PaymentID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, nil
	}

	if payment.Status == model.PaymentStatusSucceeded {
		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("cannot commit callback transaction: %w", err)
		}
		return payment, nil
	}

	now := time.Now()

	nextStatus := data.Status
	var paidAt *time.Time
	if nextStatus == model.PaymentStatusSucceeded {
		paidAt = &now
	}

	updatePaymentQuery := `
		UPDATE payments
		SET
			status = $1,
			provider_txn_id = $2,
			provider_order_code = $3,
			callback_payload = $4,
			failure_reason = $5,
			paid_at = $6,
			updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at;
	`

	err = tx.QueryRow(
		ctx,
		updatePaymentQuery,
		nextStatus,
		data.ProviderTxnID,
		data.ProviderOrderCode,
		data.RawPayload,
		data.FailureReason,
		paidAt,
		payment.ID,
	).Scan(&payment.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot update payment from callback: %w", err)
	}

	payment.Status = nextStatus
	payment.ProviderTxnID = data.ProviderTxnID
	payment.ProviderOrderCode = data.ProviderOrderCode
	payment.CallbackPayload = data.RawPayload
	payment.FailureReason = data.FailureReason
	payment.PaidAt = paidAt

	orderPaymentStatus := "pending"
	if nextStatus == model.PaymentStatusSucceeded {
		orderPaymentStatus = "paid"
	}

	updateOrderQuery := `
		UPDATE orders
		SET
			payment_status = $1,
			updated_at = NOW()
		WHERE id = $2;
	`

	if _, err := tx.Exec(ctx, updateOrderQuery, orderPaymentStatus, payment.OrderID); err != nil {
		return nil, fmt.Errorf("cannot update order payment status from callback: %w", err)
	}

	if nextStatus == model.PaymentStatusSucceeded {
		updateOrderStatusQuery := `
			UPDATE orders
			SET
				status = $1,
				updated_at = NOW()
			WHERE id = $2 AND status IN ($3, $4);
		`

		if _, err := tx.Exec(
			ctx,
			updateOrderStatusQuery,
			"paid",
			payment.OrderID,
			"pending",
			"awaiting_payment",
		); err != nil {
			return nil, fmt.Errorf("cannot update order status after payment success: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit callback transaction: %w", err)
	}

	return payment, nil
}

func getOrderForPayment(ctx context.Context, tx pgx.Tx, orderID string) (*orderPaymentSnapshot, error) {
	query := `
		SELECT
			id,
			status,
			payment_status,
			total_price
		FROM orders
		WHERE id = $1
		FOR UPDATE;
	`

	var order orderPaymentSnapshot

	err := tx.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.Status,
		&order.PaymentStatus,
		&order.TotalPrice,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get order for payment: %w", err)
	}

	return &order, nil
}

func getLatestPendingPaymentByOrderID(ctx context.Context, tx pgx.Tx, orderID string) (*model.Payment, error) {
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
		WHERE order_id = $1
			AND status IN ($2, $3)
		ORDER BY created_at DESC
		LIMIT 1;
	`

	var payment model.Payment

	err := tx.QueryRow(
		ctx,
		query,
		orderID,
		model.PaymentStatusPending,
		model.PaymentStatusProcessing,
	).Scan(
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
		return nil, fmt.Errorf("cannot get latest pending payment by order id: %w", err)
	}

	return &payment, nil
}

func getPaymentForUpdate(ctx context.Context, tx pgx.Tx, paymentID string) (*model.Payment, error) {
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
		FOR UPDATE;
	`

	var payment model.Payment

	err := tx.QueryRow(ctx, query, paymentID).Scan(
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
		return nil, fmt.Errorf("cannot get payment for update: %w", err)
	}

	return &payment, nil
}

func buildProviderOrderCode(orderID string) string {
	return fmt.Sprintf("PAY-%s", orderID)
}

func buildCheckoutURL(provider, providerOrderCode string) string {
	return fmt.Sprintf("https://pay.local/%s/%s", provider, providerOrderCode)
}
