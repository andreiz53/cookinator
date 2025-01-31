package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CodeDuplicateKey = "23505"
)

var ErrDuplicateKey = &pgconn.PgError{
	Code: CodeDuplicateKey,
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
