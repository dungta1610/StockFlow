package storage

import "github.com/jackc/pgx/v5/pgxpool"

type sqlStore struct {
	db *pgxpool.Pool
}

func NewSQLStore(db *pgxpool.Pool) *sqlStore {
	return &sqlStore{db: db}
}
