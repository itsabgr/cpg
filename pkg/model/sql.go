package model

import (
	"context"
	"database/sql"
	"errors"
)

func Transaction[T any](ctx context.Context, db DB, opt *sql.TxOptions,
	fn func(ctx context.Context, tx Tx) (T, error)) (t T, err error) {
	tx, err := db.BeginTx(ctx, opt)
	if err != nil {
		return t, err
	}
	defer tx.Rollback()
	if t, err = fn(ctx, tx); err != nil {
		return t, tx.Rollback()
	}
	return t, tx.Commit()
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, sql.ErrNoRows)
}

type Tx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
type DB interface {
	Tx
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
type ErrorWithAffected struct {
	Err      error
	Affected int64
}

func (ea ErrorWithAffected) ExactAffect(n int64) error {
	if ea.Err != nil {
		return ea.Err
	}
	if ea.Affected != n {
		return errors.New("not affected as desired")
	}
	return nil
}

func Affected(res sql.Result, err error) ErrorWithAffected {
	if err != nil {
		return ErrorWithAffected{Err: err}
	}
	n, err := res.RowsAffected()
	return ErrorWithAffected{Affected: n, Err: err}
}

type RollbackErr struct{}

func IsRollback(err error) bool {
	_, is := err.(interface{ Rollback() })
	return is
}
func (RollbackErr) Error() string { return "rollback" }

func (RollbackErr) Rollback() {}
