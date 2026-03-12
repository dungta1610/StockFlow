package storage

import (
	"context"
	"fmt"
	"strings"

	"stockflow/module/warehouse/model"

	"github.com/jackc/pgx/v5"
)

func (s *SQLStore) CreateWarehouse(ctx context.Context, data *model.Warehouse) error {
	query := `
		INSERT INTO warehouses (
			code,
			name,
			address,
			is_active
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at;
	`

	err := s.db.QueryRow(
		ctx,
		query,
		data.Code,
		data.Name,
		data.Address,
		data.IsActive,
	).Scan(&data.ID, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot create warehouse: %w", err)
	}

	return nil
}

func (s *SQLStore) FindWarehouseByCode(ctx context.Context, code string) (*model.Warehouse, error) {
	query := `
		SELECT
			id,
			code,
			name,
			address,
			is_active,
			created_at,
			updated_at
		FROM warehouses
		WHERE code = $1
		LIMIT 1;
	`

	var warehouse model.Warehouse

	err := s.db.QueryRow(ctx, query, code).Scan(
		&warehouse.ID,
		&warehouse.Code,
		&warehouse.Name,
		&warehouse.Address,
		&warehouse.IsActive,
		&warehouse.CreatedAt,
		&warehouse.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot find warehouse by code: %w", err)
	}

	return &warehouse, nil
}

func (s *SQLStore) GetWarehouseByID(ctx context.Context, id string) (*model.Warehouse, error) {
	query := `
		SELECT
			id,
			code,
			name,
			address,
			is_active,
			created_at,
			updated_at
		FROM warehouses
		WHERE id = $1
		LIMIT 1;
	`

	var warehouse model.Warehouse

	err := s.db.QueryRow(ctx, query, id).Scan(
		&warehouse.ID,
		&warehouse.Code,
		&warehouse.Name,
		&warehouse.Address,
		&warehouse.IsActive,
		&warehouse.CreatedAt,
		&warehouse.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get warehouse by id: %w", err)
	}

	return &warehouse, nil
}

func (s *SQLStore) ListWarehouses(
	ctx context.Context,
	filter *model.Filter,
	paging *model.Paging,
) ([]model.Warehouse, error) {
	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0)
	argPos := 1

	queryBuilder.WriteString(`
		SELECT
			id,
			code,
			name,
			address,
			is_active,
			created_at,
			updated_at
		FROM warehouses
		WHERE 1=1
	`)

	if filter != nil {
		if filter.Code != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND code = $%d", argPos))
			args = append(args, filter.Code)
			argPos++
		}

		if filter.Name != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND name ILIKE $%d", argPos))
			args = append(args, "%"+filter.Name+"%")
			argPos++
		}

		if filter.IsActive != nil {
			queryBuilder.WriteString(fmt.Sprintf(" AND is_active = $%d", argPos))
			args = append(args, *filter.IsActive)
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
		return nil, fmt.Errorf("cannot list warehouses: %w", err)
	}
	defer rows.Close()

	warehouses := make([]model.Warehouse, 0)

	for rows.Next() {
		var warehouse model.Warehouse

		if err := rows.Scan(
			&warehouse.ID,
			&warehouse.Code,
			&warehouse.Name,
			&warehouse.Address,
			&warehouse.IsActive,
			&warehouse.CreatedAt,
			&warehouse.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan warehouse: %w", err)
		}

		warehouses = append(warehouses, warehouse)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate warehouse rows: %w", err)
	}

	return warehouses, nil
}
