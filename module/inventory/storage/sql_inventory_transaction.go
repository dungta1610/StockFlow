package storage

import (
	"context"
	"fmt"
	"strings"

	"stockflow/module/inventory/model"
)

func (s *SQLStore) CreateInventoryTransaction(ctx context.Context, data *model.InventoryTransactionCreate) error {
	query := `
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at;
	`

	if err := data.Validate(); err != nil {
		return err
	}

	err := s.db.QueryRow(
		ctx,
		query,
		data.InventoryID,
		data.ProductID,
		data.WarehouseID,
		data.OrderID,
		data.ReservationID,
		data.TxnType,
		data.Quantity,
		data.BeforeAvailableQty,
		data.AfterAvailableQty,
		data.BeforeReservedQty,
		data.AfterReservedQty,
		data.Reason,
		data.CreatedBy,
	).Scan(&data.ID, &data.CreatedAt)
	if err != nil {
		return fmt.Errorf("cannot create inventory transaction: %w", err)
	}

	return nil
}

func (s *SQLStore) ListInventoryTransactions(ctx context.Context, filter *model.TransactionFilter, paging *model.Paging) ([]model.InventoryTransaction, error) {
	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0)
	argPos := 1

	queryBuilder.WriteString(`
		SELECT
			id,
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
			created_by,
			created_at
		FROM inventory_transactions
		WHERE 1=1
	`)

	if filter != nil {
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

		if filter.OrderID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND order_id = $%d", argPos))
			args = append(args, filter.OrderID)
			argPos++
		}

		if filter.ReservationID != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND reservation_id = $%d", argPos))
			args = append(args, filter.ReservationID)
			argPos++
		}

		if filter.TxnType != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND txn_type = $%d", argPos))
			args = append(args, filter.TxnType)
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
		return nil, fmt.Errorf("cannot list inventory transactions: %w", err)
	}
	defer rows.Close()

	transactions := make([]model.InventoryTransaction, 0)

	for rows.Next() {
		var item model.InventoryTransaction

		if err := rows.Scan(
			&item.ID,
			&item.InventoryID,
			&item.ProductID,
			&item.WarehouseID,
			&item.OrderID,
			&item.ReservationID,
			&item.TxnType,
			&item.Quantity,
			&item.BeforeAvailableQty,
			&item.AfterAvailableQty,
			&item.BeforeReservedQty,
			&item.AfterReservedQty,
			&item.Reason,
			&item.CreatedBy,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan inventory transaction: %w", err)
		}

		transactions = append(transactions, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate inventory transaction rows: %w", err)
	}

	return transactions, nil
}
