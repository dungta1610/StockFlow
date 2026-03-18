package storage

import (
	"context"
	"fmt"
	"strings"

	"stockflow/module/order/model"

	"github.com/jackc/pgx/v5"
)

func (s *SQLStore) GetOrderByID(ctx context.Context, id string) (*model.Order, error) {
	orderQuery := `
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
		LIMIT 1;
	`

	var order model.Order

	err := s.db.QueryRow(ctx, orderQuery, id).Scan(
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
		return nil, fmt.Errorf("cannot get order by id: %w", err)
	}

	itemsQuery := `
		SELECT
			id,
			order_id,
			product_id,
			product_sku,
			product_name,
			quantity,
			unit_price,
			line_price,
			created_at,
			updated_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at ASC;
	`

	rows, err := s.db.Query(ctx, itemsQuery, order.ID)
	if err != nil {
		return nil, fmt.Errorf("cannot get order items: %w", err)
	}
	defer rows.Close()

	items := make([]model.OrderItem, 0)
	for rows.Next() {
		var item model.OrderItem

		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.ProductSKU,
			&item.ProductName,
			&item.Quantity,
			&item.UnitPrice,
			&item.LinePrice,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan order item: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate order item rows: %w", err)
	}

	order.Items = items

	return &order, nil
}

func (s *SQLStore) ListOrders(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.Order, error) {
	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0)
	argPos := 1

	queryBuilder.WriteString(`
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
		WHERE 1=1
	`)

	if filter != nil {
		if filter.Code != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND code = $%d", argPos))
			args = append(args, strings.ToUpper(strings.TrimSpace(filter.Code)))
			argPos++
		}

		if filter.UserID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND user_id = $%d", argPos))
			args = append(args, strings.TrimSpace(filter.UserID))
			argPos++
		}

		if filter.WarehouseID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND warehouse_id = $%d", argPos))
			args = append(args, strings.TrimSpace(filter.WarehouseID))
			argPos++
		}

		if filter.Status != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", argPos))
			args = append(args, strings.TrimSpace(filter.Status))
			argPos++
		}

		if filter.PaymentStatus != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND payment_status = $%d", argPos))
			args = append(args, strings.TrimSpace(filter.PaymentStatus))
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
		return nil, fmt.Errorf("cannot list orders: %w", err)
	}
	defer rows.Close()

	orders := make([]model.Order, 0)

	for rows.Next() {
		var order model.Order

		if err := rows.Scan(
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
		); err != nil {
			return nil, fmt.Errorf("cannot scan order: %w", err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate order rows: %w", err)
	}

	return orders, nil
}
