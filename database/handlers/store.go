package database

import "github.com/jackc/pgx/v5"

type Store interface {
	Querier
}

type PostgresStore struct {
	*Queries
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) Store {
	return &PostgresStore{
		db:      db,
		Queries: New(db),
	}
}
