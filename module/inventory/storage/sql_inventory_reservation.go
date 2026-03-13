package storage

import (
	"context"
	"fmt"
	"strings"

	"stockflow/module/inventory/model"

	"github.com/jackc/pgx/v5"
)

func (s *SQLStore) CreateInventoryReservation(ctx context.Context, data *model.InventoryReservationCreate) error {
	if err := data.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO inventory_reservations (
			order_id,
			order_item_id,
			inventory_id,
			product_id,
			warehouse_id,
			quantity,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, reserved_at, created_at, updated_at;
	`

	err := s.db.QueryRow(
		ctx,
		query,
		data.OrderID,
		data.OrderItemID,
		data.InventoryID,
		data.ProductID,
		data.WarehouseID,
		data.Quantity,
		data.Status,
	).Scan(&data.ID, &data.ReservedAt, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot create inventory reservation: %w", err)
	}

	return nil
}

func (s *SQLStore) GetInventoryReservationByID(ctx context.Context, id string) (*model.InventoryReservation, error) {
	query := `
		SELECT
			id,
			order_id,
			order_item_id,
			inventory_id,
			product_id,
			warehouse_id,
			quantity,
			status,
			reserved_at,
			released_at,
			consumed_at,
			created_at,
			updated_at
		FROM inventory_reservations
		WHERE id = $1
		LIMIT 1;
	`

	var item model.InventoryReservation

	err := s.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.OrderID,
		&item.OrderItemID,
		&item.InventoryID,
		&item.ProductID,
		&item.WarehouseID,
		&item.Quantity,
		&item.Status,
		&item.ReservedAt,
		&item.ReleasedAt,
		&item.ConsumedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get inventory reservation by id: %w", err)
	}

	return &item, nil
}

func (s *SQLStore) ListInventoryReservations(ctx context.Context, filter *model.ReservationFilter, paging *model.Paging) ([]model.InventoryReservation, error) {
	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0)
	argPos := 1

	queryBuilder.WriteString(`
		SELECT
			id,
			order_id,
			order_item_id,
			inventory_id,
			product_id,
			warehouse_id,
			quantity,
			status,
			reserved_at,
			released_at,
			consumed_at,
			created_at,
			updated_at
		FROM inventory_reservations
		WHERE 1=1
	`)

	if filter != nil {
		if filter.OrderID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND order_id = $%d", argPos))
			args = append(args, filter.OrderID)
			argPos++
		}

		if filter.OrderItemID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND order_item_id = $%d", argPos))
			args = append(args, filter.OrderItemID)
			argPos++
		}

		if filter.InventoryID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND inventory_id = $%d", argPos))
			args = append(args, filter.InventoryID)
			argPos++
		}

		if filter.ProductID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND product_id = $%d", argPos))
			args = append(args, filter.ProductID)
			argPos++
		}

		if filter.WarehouseID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND warehouse_id = $%d", argPos))
			args = append(args, filter.WarehouseID)
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
		return nil, fmt.Errorf("cannot list inventory reservations: %w", err)
	}
	defer rows.Close()

	items := make([]model.InventoryReservation, 0)

	for rows.Next() {
		var item model.InventoryReservation

		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.OrderItemID,
			&item.InventoryID,
			&item.ProductID,
			&item.WarehouseID,
			&item.Quantity,
			&item.Status,
			&item.ReservedAt,
			&item.ReleasedAt,
			&item.ConsumedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan inventory reservation: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate inventory reservation rows: %w", err)
	}

	return items, nil
}
