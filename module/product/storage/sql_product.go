package storage

import (
	"context"
	"fmt"
	"strings"

	"stockflow/module/product/model"

	"github.com/jackc/pgx/v5"
)

func (s *SQLStore) CreateProduct(ctx context.Context, data *model.Product) error {
	query := `
		INSERT INTO products (
			sku,
			name,
			description,
			price,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at;
	`

	err := s.db.QueryRow(
		ctx,
		query,
		data.SKU,
		data.Name,
		data.Description,
		data.Price,
		data.IsActive,
	).Scan(&data.ID, &data.CreatedAt, &data.UpdatedAt)

	if err != nil {
		return fmt.Errorf("cannot create product: %w", err)
	}

	return nil
}

func (s *SQLStore) FindProductBySKU(ctx context.Context, sku string) (*model.Product, error) {
	query := `
		SELECT
			id,
			sku,
			name,
			description,
			price,
			is_active,
			created_at,
			updated_at
		FROM products
		WHERE sku = $1
		LIMIT 1;
	`

	var product model.Product

	err := s.db.QueryRow(ctx, query, sku).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("cannot find product by sku: %w", err)
	}

	return &product, nil
}

func (s *SQLStore) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	query := `
		SELECT
			id,
			sku,
			name,
			description,
			price,
			is_active,
			created_at,
			updated_at
		FROM products
		WHERE id = $1
		LIMIT 1;
	`

	var product model.Product

	err := s.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get product by id: %w", err)
	}

	return &product, nil
}

func (s *SQLStore) ListProducts(
	ctx context.Context,
	filter *model.Filter,
	paging *model.Paging,
) ([]model.Product, error) {
	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0)
	argPos := 1

	queryBuilder.WriteString(`
		SELECT
			id,
			sku,
			name,
			description,
			price,
			is_active,
			created_at,
			updated_at
		FROM products
		WHERE 1=1
	`)

	if filter != nil {
		if filter.SKU != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND sku = $%d", argPos))
			args = append(args, filter.SKU)
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
		return nil, fmt.Errorf("cannot list products: %w", err)
	}
	defer rows.Close()

	products := make([]model.Product, 0)

	for rows.Next() {
		var product model.Product

		if err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan product: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate product rows: %w", err)
	}

	return products, nil
}
