package storage

import (
	"context"
	"fmt"
	"strings"

	"stockflow/module/inventory/model"

	"github.com/jackc/pgx/v5"
)

func (s *SQLStore) AdjustStock(ctx context.Context, data *model.InventoryAdjust) (*model.Inventory, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT
			id,
			product_id,
			warehouse_id,
			available_qty,
			reserved_qty,
			version,
			created_at,
			updated_at
		FROM inventory
		WHERE product_id = $1 AND warehouse_id = $2
		FOR UPDATE;
	`

	var inventory model.Inventory

	err = tx.QueryRow(ctx, query, data.ProductID, data.WarehouseID).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.WarehouseID,
		&inventory.AvailableQty,
		&inventory.ReservedQty,
		&inventory.Version,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get inventory for adjust stock: %w", err)
	}

	beforeAvailable := inventory.AvailableQty
	beforeReserved := inventory.ReservedQty

	afterAvailable := inventory.AvailableQty + data.Quantity
	if afterAvailable < 0 {
		return nil, model.ErrInventoryNotEnoughStock
	}

	updateQuery := `
		UPDATE inventory
		SET
			available_qty = $1,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $2
		RETURNING version, updated_at;
	`

	err = tx.QueryRow(ctx, updateQuery, afterAvailable, inventory.ID).Scan(
		&inventory.Version,
		&inventory.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot update inventory: %w", err)
	}

	inventory.AvailableQty = afterAvailable

	var createdBy *string
	if data.CreatedBy != "" {
		createdBy = &data.CreatedBy
	}

	reason := data.Reason
	if reason == "" {
		reason = "manual_adjustment"
	}

	insertTxnQuery := `
		INSERT INTO inventory_transactions (
			inventory_id,
			product_id,
			warehouse_id,
			order_id,
			reservation_id,
			txn_type,
			quantity,
			before_available_qty,
			after_available_qty,
			before_reserved_qty,
			after_reserved_qty,
			reason,
			created_by
		)
		VALUES ($1, $2, $3, NULL, NULL, $4, $5, $6, $7, $8, $9, $10, $11);
	`

	_, err = tx.Exec(
		ctx,
		insertTxnQuery,
		inventory.ID,
		inventory.ProductID,
		inventory.WarehouseID,
		"manual_adjustment",
		abs(data.Quantity),
		beforeAvailable,
		afterAvailable,
		beforeReserved,
		beforeReserved,
		reason,
		createdBy,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create inventory transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("cannot commit adjust stock transaction: %w", err)
	}

	return &inventory, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (s *SQLStore) CreateInventory(ctx context.Context, data *model.Inventory) error {
	query := `
		INSERT INTO inventory (
			product_id,
			warehouse_id,
			available_qty,
			reserved_qty,
			version
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at;
	`

	err := s.db.QueryRow(
		ctx,
		query,
		data.ProductID,
		data.WarehouseID,
		data.AvailableQty,
		data.ReservedQty,
		data.Version,
	).Scan(&data.ID, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot create inventory: %w", err)
	}

	return nil
}

func (s *SQLStore) GetInventoryByID(ctx context.Context, id string) (*model.Inventory, error) {
	query := `
		SELECT
			id,
			product_id,
			warehouse_id,
			available_qty,
			reserved_qty,
			version,
			created_at,
			updated_at
		FROM inventory
		WHERE id = $1
		LIMIT 1;
	`

	var inventory model.Inventory

	err := s.db.QueryRow(ctx, query, id).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.WarehouseID,
		&inventory.AvailableQty,
		&inventory.ReservedQty,
		&inventory.Version,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get inventory by id: %w", err)
	}

	return &inventory, nil
}

func (s *SQLStore) GetInventoryByProductAndWarehouse(ctx context.Context, productID, warehouseID string) (*model.Inventory, error) {
	query := `
		SELECT
			id,
			product_id,
			warehouse_id,
			available_qty,
			reserved_qty,
			version,
			created_at,
			updated_at
		FROM inventory
		WHERE product_id = $1 AND warehouse_id = $2
		LIMIT 1;
	`

	var inventory model.Inventory

	err := s.db.QueryRow(ctx, query, productID, warehouseID).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.WarehouseID,
		&inventory.AvailableQty,
		&inventory.ReservedQty,
		&inventory.Version,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get inventory by product and warehouse: %w", err)
	}

	return &inventory, nil
}

func (s *SQLStore) ListInventory(ctx context.Context, filter *model.Filter, paging *model.Paging) ([]model.Inventory, error) {
	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0)
	argPos := 1

	queryBuilder.WriteString(`
		SELECT
			id,
			product_id,
			warehouse_id,
			available_qty,
			reserved_qty,
			version,
			created_at,
			updated_at
		FROM inventory
		WHERE 1=1
	`)

	if filter != nil {
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
	}

	queryBuilder.WriteString(" ORDER BY created_at DESC")

	if paging != nil {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1))
		args = append(args, paging.Limit, paging.Offset())
	}

	rows, err := s.db.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("cannot list inventory: %w", err)
	}
	defer rows.Close()

	inventories := make([]model.Inventory, 0)

	for rows.Next() {
		var inventory model.Inventory

		if err := rows.Scan(
			&inventory.ID,
			&inventory.ProductID,
			&inventory.WarehouseID,
			&inventory.AvailableQty,
			&inventory.ReservedQty,
			&inventory.Version,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan inventory: %w", err)
		}

		inventories = append(inventories, inventory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate inventory rows: %w", err)
	}

	return inventories, nil
}
