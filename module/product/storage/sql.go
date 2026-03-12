package storage

import "github.com/jackc/pgx/v5/pgxpool"

type SQLStore struct {
	db *pgxpool.Pool
}

func NewSQLStore(db *pgxpool.Pool) *SQLStore {
	return &SQLStore{db: db}
}
