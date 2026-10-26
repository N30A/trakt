package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// https://www.postgresql.org/docs/current/errcodes-appendix.html

const (
	PGErrUniqueViolation = "23505"
)

func IsError(err error, errorCode string) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	return ok && pgErr.Code == errorCode
}
