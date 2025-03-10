//go:build go1.18
// +build go1.18

package sqlext

import (
	"context"
	"database/sql"

	resultext "github.com/khulnasoft-lab/utils/values/result"
)

// DoTransaction is a helper function that abstracts some complexities of dealing with a transaction and rolling it back.
func DoTransaction[T any](ctx context.Context, opts *sql.TxOptions, conn *sql.DB, fn func(context.Context, *sql.Tx) resultext.Result[T, error]) resultext.Result[T, error] {
	tx, err := conn.BeginTx(ctx, opts)
	if err != nil {
		return resultext.Err[T, error](err)
	}
	result := fn(ctx, tx)
	if result.IsErr() {
		_ = tx.Rollback()
		return result
	}
	err = tx.Commit()
	if err != nil {
		return resultext.Err[T, error](err)
	}
	return result
}
