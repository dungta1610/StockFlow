package storage

import (
	"context"
	"fmt"
	"time"

	"stockflow/module/order/model"

	"github.com/jackc/pgx/v5"
)

type productSnapshot struct {
	ID    string
	SKU   string
	Name  string
	Price float64
}

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

	var subtotal float64
	items := make([]model.OrderItem, 0, len(data.Items))

	for _, reqItem := range data.Items {
		product, err := getProductSnapshotForOrder(ctx, tx, reqItem.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, fmt.Errorf("product not found: %s", reqItem.ProductID)
		}

		unitPrice := reqItem.UnitPrice
		if unitPrice <= 0 {
			unitPrice = product.Price
		}

		linePrice := unitPrice * float64(reqItem.Quantity)
		subtotal += linePrice

		items = append(items, model.OrderItem{
			ProductID:   product.ID,
			ProductSKU:  product.SKU,
			ProductName: product.Name,
			Quantity:    reqItem.Quantity,
			UnitPrice:   unitPrice,
			LinePrice:   linePrice,
		})
	}

	orderStatus := model.OrderStatusPending
	if data.ExpiredAt != nil {
		orderStatus = model.OrderStatusReserved
	}

	paymentStatus := model.PaymentStatusPending
	discountPrice := 0.0
	totalPrice := subtotal - discountPrice

	orderInsertQuery := `
		INSERT INTO orders (
			code,
			user_id,
			warehouse_id,
			status,
			payment_status,
			subtotal_price,
			discount_price,
			total_price,
			note,
			expired_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING
			id,
			created_at,
			updated_at;
	`

	var order model.Order
	order.Code = orderCode
	order.UserID = data.UserID
	order.WarehouseID = data.WarehouseID
	order.Status = orderStatus
	order.PaymentStatus = paymentStatus
	order.SubtotalPrice = subtotal
	order.DiscountPrice = discountPrice
	order.TotalPrice = totalPrice
	order.Note = data.Note
	order.ExpiredAt = data.ExpiredAt

	err = tx.QueryRow(
		ctx,
		orderInsertQuery,
		order.Code,
		order.UserID,
		order.WarehouseID,
		order.Status,
		order.PaymentStatus,
		order.SubtotalPrice,
		order.DiscountPrice,
		order.TotalPrice,
		order.Note,
		order.ExpiredAt,
	).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create order: %w", err)
	}

	itemInsertQuery := `
		INSERT INTO order_items (
			order_id,
			product_id,
			product_sku,
			product_name,
			quantity,
			unit_price,
			line_price
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id,
			created_at,
			updated_at;
	`

	for idx := range items {
		items[idx].OrderID = order.ID

		err := tx.QueryRow(
			ctx,
			itemInsertQuery,
			items[idx].OrderID,
			items[idx].ProductID,
			items[idx].ProductSKU,
			items[idx].ProductName,
			items[idx].Quantity,
			items[idx].UnitPrice,
			items[idx].LinePrice,
		).Scan(
			&items[idx].ID,
			&items[idx].CreatedAt,
			&items[idx].UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot create order item: %w", err)
		}
	}

	order.Items = items

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit create order transaction: %w", err)
	}

	return &order, nil
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
			code,
			user_id,
			warehouse_id,
			status,
			payment_status,
			subtotal_price,
			discount_price,
			total_price,
			note,
			expired_at,
			canceled_at,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		FOR UPDATE;
	`

	var order model.Order

	err = tx.QueryRow(ctx, lockQuery, data.OrderID).Scan(
		&order.ID,
		&order.Code,
		&order.UserID,
		&order.WarehouseID,
		&order.Status,
		&order.PaymentStatus,
		&order.SubtotalPrice,
		&order.DiscountPrice,
		&order.TotalPrice,
		&order.Note,
		&order.ExpiredAt,
		&order.CanceledAt,
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
	case model.OrderStatusCanceled:
		return nil, model.ErrOrderAlreadyCanceled
	case model.OrderStatusExpired:
		return nil, model.ErrOrderCannotBeCanceled
	case model.OrderStatusCompleted:
		return nil, model.ErrOrderCannotBeCanceled
	}

	if order.PaymentStatus == model.PaymentStatusPaid {
		return nil, model.ErrOrderCannotBeCanceled
	}

	now := time.Now()

	updateQuery := `
		UPDATE orders
		SET
			status = $1,
			canceled_at = $2,
			updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at;
	`

	err = tx.QueryRow(ctx, updateQuery, model.OrderStatusCanceled, now, order.ID).Scan(&order.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot cancel order: %w", err)
	}

	order.Status = model.OrderStatusCanceled
	order.CanceledAt = &now

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
			code,
			user_id,
			warehouse_id,
			status,
			payment_status,
			subtotal_price,
			discount_price,
			total_price,
			note,
			expired_at,
			canceled_at,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		FOR UPDATE;
	`

	var order model.Order

	err = tx.QueryRow(ctx, lockQuery, data.OrderID).Scan(
		&order.ID,
		&order.Code,
		&order.UserID,
		&order.WarehouseID,
		&order.Status,
		&order.PaymentStatus,
		&order.SubtotalPrice,
		&order.DiscountPrice,
		&order.TotalPrice,
		&order.Note,
		&order.ExpiredAt,
		&order.CanceledAt,
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
	case model.OrderStatusCanceled:
		return nil, model.ErrOrderCannotBeExpired
	case model.OrderStatusCompleted:
		return nil, model.ErrOrderCannotBeExpired
	}

	if order.PaymentStatus == model.PaymentStatusPaid {
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

func getProductSnapshotForOrder(ctx context.Context, tx pgx.Tx, productID string) (*productSnapshot, error) {
	query := `
		SELECT
			id,
			sku,
			name,
			price
		FROM products
		WHERE id = $1
		LIMIT 1;
	`

	var item productSnapshot

	err := tx.QueryRow(ctx, query, productID).Scan(
		&item.ID,
		&item.SKU,
		&item.Name,
		&item.Price,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get product snapshot: %w", err)
	}

	return &item, nil
}
