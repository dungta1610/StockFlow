package storage

import (
	"context"
	"fmt"
	"strings"

	"stockflow/module/user/model"

	"github.com/jackc/pgx/v5"
)

func (s *SQLStore) CreateUser(ctx context.Context, data *model.UserCreate) (*model.User, error) {
	query := `
		INSERT INTO users (
			email,
			password_hash,
			full_name,
			role,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			email,
			password_hash,
			full_name,
			role,
			is_active,
			created_at,
			updated_at;
	`

	var user model.User

	err := s.db.QueryRow(
		ctx,
		query,
		data.Email,
		data.PasswordHash,
		data.FullName,
		data.Role,
		true,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create user: %w", err)
	}

	return &user, nil
}

func (s *SQLStore) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT
			id,
			email,
			password_hash,
			full_name,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
		LIMIT 1;
	`

	var user model.User

	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get user by id: %w", err)
	}

	return &user, nil
}

func (s *SQLStore) ListUsers(
	ctx context.Context,
	filter *model.Filter,
	paging *model.Paging,
) ([]model.User, error) {
	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0)
	argPos := 1

	queryBuilder.WriteString(`
		SELECT
			id,
			email,
			password_hash,
			full_name,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE 1=1
	`)

	if filter != nil {
		if filter.Email != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND email = $%d", argPos))
			args = append(args, filter.Email)
			argPos++
		}

		if filter.FullName != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND full_name ILIKE $%d", argPos))
			args = append(args, "%"+filter.FullName+"%")
			argPos++
		}

		if filter.Role != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND role = $%d", argPos))
			args = append(args, filter.Role)
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
		return nil, fmt.Errorf("cannot list users: %w", err)
	}
	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var user model.User

		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan user: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate user rows: %w", err)
	}

	return users, nil
}

func (s *SQLStore) UpdateUser(ctx context.Context, id string, data *model.UserUpdate) (*model.User, error) {
	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0)
	argPos := 1

	queryBuilder.WriteString(`
		UPDATE users
		SET
	`)

	queryBuilder.WriteString(fmt.Sprintf(" full_name = $%d,", argPos))
	args = append(args, data.FullName)
	argPos++

	queryBuilder.WriteString(fmt.Sprintf(" role = $%d,", argPos))
	args = append(args, data.Role)
	argPos++

	if data.PasswordHash != nil {
		queryBuilder.WriteString(fmt.Sprintf(" password_hash = $%d,", argPos))
		args = append(args, *data.PasswordHash)
		argPos++
	}

	if data.IsActive != nil {
		queryBuilder.WriteString(fmt.Sprintf(" is_active = $%d,", argPos))
		args = append(args, *data.IsActive)
		argPos++
	}

	queryBuilder.WriteString(" updated_at = NOW() ")
	queryBuilder.WriteString(fmt.Sprintf(" WHERE id = $%d", argPos))
	args = append(args, id)
	argPos++

	queryBuilder.WriteString(`
		RETURNING
			id,
			email,
			password_hash,
			full_name,
			role,
			is_active,
			created_at,
			updated_at;
	`)

	var user model.User

	err := s.db.QueryRow(ctx, queryBuilder.String(), args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot update user: %w", err)
	}

	return &user, nil
}
