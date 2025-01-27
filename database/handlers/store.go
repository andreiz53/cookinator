package database

import "github.com/jackc/pgx/v5"

type Store struct {
	*Queries
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}
