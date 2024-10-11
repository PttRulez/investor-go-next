package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ctxKey string

var txKey = ctxKey("tx")

type Repository struct {
	db *sql.DB
}

func (r *Repository) t(ctx context.Context) DXBT {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	if !ok {
		return r.db
	}
	return tx
}

func (t *Repository) Transac(ctx context.Context, opts *sql.TxOptions,
	fn func(ctx context.Context) error) (err error) {
	if opts == nil {
		opts = &sql.TxOptions{}
	}

	// Начинаем транзакцию
	tx, err := t.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction %w", err)
	}
	ctx = context.WithValue(ctx, txKey, tx)

	// Роллбэк в случае ошибки
	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = fmt.Errorf("failed to rollback transaction %w", err)
		}
	}()

	err = fn(ctx)
	if err == nil {
		return tx.Commit()
	}

	return err
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

type DXBT interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
