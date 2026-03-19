package storage

import (
	"context"
	"fmt"
	"time"

	"stockflow/module/order/model"

	"github.com/jackc/pgx/v5"
)

func (s *SQLStore) CreateOrder(ctx context.Context, data *model.OrderCreate) (*model.Order, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot begin create order transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	orderCode, err := generateOrderCode(ctx, tx)
	if err != nil {
		return nil, err
	}

	totalAmount := 0.0
	items := make([]model.OrderItem, 0, len(data.Items))

	for _, reqItem := range data.Items {
		lineTotal := reqItem.UnitPrice * float64(reqItem.Quantity)
		totalAmount += lineTotal

		items = append(items, model.OrderItem{
			ProductID: reqItem.ProductID,
			Quantity:  reqItem.Quantity,
			UnitPrice: reqItem.UnitPrice,
			LineTotal: lineTotal,
		})
	}

	status := model.OrderStatusPending
	if data.ReservationExpiresAt != nil {
		status = model.OrderStatusReserved
	}

	insertOrderQuery := `
		INSERT INTO orders (
			order_code,
			user_id,
			warehouse_id,
			status,
			total_amount,
			reservation_expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			created_at,
			updated_at;
	`

	var order model.Order
	order.OrderCode = orderCode
	order.UserID = data.UserID
	order.WarehouseID = data.WarehouseID
	order.Status = status
	order.TotalAmount = totalAmount
	order.ReservationExpiresAt = data.ReservationExpiresAt

	err = tx.QueryRow(
		ctx,
		insertOrderQuery,
		order.OrderCode,
		order.UserID,
		order.WarehouseID,
		order.Status,
		order.TotalAmount,
		order.ReservationExpiresAt,
	).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create order: %w", err)
	}

	insertItemQuery := `
		INSERT INTO order_items (
			order_id,
			product_id,
			quantity,
			unit_price,
			line_total
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			created_at;
	`

	for i := range items {
		items[i].OrderID = order.ID

		err := tx.QueryRow(
			ctx,
			insertItemQuery,
			items[i].OrderID,
			items[i].ProductID,
			items[i].Quantity,
			items[i].UnitPrice,
			items[i].LineTotal,
		).Scan(
			&items[i].ID,
			&items[i].CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot create order item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit create order transaction: %w", err)
	}

	return s.GetOrderByID(ctx, order.ID)
}

func (s *SQLStore) CancelOrder(ctx context.Context, data *model.OrderCancel) (*model.Order, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot begin cancel order transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	lockQuery := `
		SELECT
			id,
			order_code,
			user_id,
			warehouse_id,
			status,
			total_amount,
			reservation_expires_at,
			paid_at,
			cancelled_at,
			fulfilled_at,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		FOR UPDATE;
	`

	var order model.Order

	err = tx.QueryRow(ctx, lockQuery, data.OrderID).Scan(
		&order.ID,
		&order.OrderCode,
		&order.UserID,
		&order.WarehouseID,
		&order.Status,
		&order.TotalAmount,
		&order.ReservationExpiresAt,
		&order.PaidAt,
		&order.CancelledAt,
		&order.FulfilledAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot lock order for cancel: %w", err)
	}

	switch order.Status {
	case model.OrderStatusCancelled:
		return nil, model.ErrOrderAlreadyCancelled
	case model.OrderStatusPaid, model.OrderStatusFulfilled, model.OrderStatusCompleted, model.OrderStatusExpired:
		return nil, model.ErrOrderCannotBeCancelled
	}

	now := time.Now()

	updateQuery := `
		UPDATE orders
		SET
			status = $1,
			cancelled_at = $2,
			updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at;
	`

	err = tx.QueryRow(ctx, updateQuery, model.OrderStatusCancelled, now, order.ID).Scan(&order.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot cancel order: %w", err)
	}

	order.Status = model.OrderStatusCancelled
	order.CancelledAt = &now

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit cancel order transaction: %w", err)
	}

	return s.GetOrderByID(ctx, order.ID)
}

func (s *SQLStore) ExpireOrder(ctx context.Context, data *model.OrderExpire) (*model.Order, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot begin expire order transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	lockQuery := `
		SELECT
			id,
			order_code,
			user_id,
			warehouse_id,
			status,
			total_amount,
			reservation_expires_at,
			paid_at,
			cancelled_at,
			fulfilled_at,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		FOR UPDATE;
	`

	var order model.Order

	err = tx.QueryRow(ctx, lockQuery, data.OrderID).Scan(
		&order.ID,
		&order.OrderCode,
		&order.UserID,
		&order.WarehouseID,
		&order.Status,
		&order.TotalAmount,
		&order.ReservationExpiresAt,
		&order.PaidAt,
		&order.CancelledAt,
		&order.FulfilledAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot lock order for expire: %w", err)
	}

	switch order.Status {
	case model.OrderStatusExpired:
		return nil, model.ErrOrderAlreadyExpired
	case model.OrderStatusPaid, model.OrderStatusCancelled, model.OrderStatusFulfilled, model.OrderStatusCompleted:
		return nil, model.ErrOrderCannotBeExpired
	}

	updateQuery := `
		UPDATE orders
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at;
	`

	err = tx.QueryRow(ctx, updateQuery, model.OrderStatusExpired, order.ID).Scan(&order.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot expire order: %w", err)
	}

	order.Status = model.OrderStatusExpired

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit expire order transaction: %w", err)
	}

	return s.GetOrderByID(ctx, order.ID)
}

func generateOrderCode(ctx context.Context, tx pgx.Tx) (string, error) {
	query := `
		SELECT CONCAT('ORD-', TO_CHAR(NOW(), 'YYYYMMDD'), '-', LPAD((FLOOR(RANDOM() * 1000000))::text, 6, '0'));
	`

	var code string
	if err := tx.QueryRow(ctx, query).Scan(&code); err != nil {
		return "", fmt.Errorf("cannot generate order code: %w", err)
	}

	return code, nil
}
