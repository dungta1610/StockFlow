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
		LIMIT 1;
	`

	var order model.Order

	err := s.db.QueryRow(ctx, orderQuery, id).Scan(
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
		return nil, fmt.Errorf("cannot get order by id: %w", err)
	}

	itemsQuery := `
		SELECT
			id,
			order_id,
			product_id,
			quantity,
			unit_price,
			line_total,
			created_at
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
			&item.Quantity,
			&item.UnitPrice,
			&item.LineTotal,
			&item.CreatedAt,
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
		WHERE 1=1
	`)

	if filter != nil {
		if filter.OrderCode != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND order_code = $%d", argPos))
			args = append(args, strings.ToUpper(strings.TrimSpace(filter.OrderCode)))
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
