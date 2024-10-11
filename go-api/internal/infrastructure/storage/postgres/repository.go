package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type ctxKey string

var txKey = ctxKey("tx")

type withDB struct {
	db *sql.DB
}

type Repository struct {
	withDB
}

func (r *Repository) t(ctx context.Context) DXBT {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	if !ok {
		return r.db
	}
	return tx
}

type Transactioner struct {
	withDB
}

func (t *Transactioner) Transac(ctx context.Context, opts *sql.TxOptions,
	fn func(ctx context.Context) error) (err error) {
	if opts == nil {
		opts = &sql.TxOptions{}
	}

	// Начинаем транзакцию
	tx, err := t.db.BeginTx(ctx, opts)
	if err == nil {
		return fmt.Errorf("failed to begin transaction %w", err)
	}
	ctx = context.WithValue(ctx, txKey, tx)

	// Роллбэк в случае ошибки
	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			err = fmt.Errorf("failed to begin transaction %w", err)
		}
	}()

	err = fn(ctx)
	if err == nil {
		return tx.Commit()
	}

	return nil
}

func (t *Repository) Transac(ctx context.Context, opts *sql.TxOptions,
	fn func(ctx context.Context) error) (err error) {
	if opts == nil {
		opts = &sql.TxOptions{}
	}

	// Начинаем транзакцию
	tx, err := t.db.BeginTx(ctx, opts)
	if err == nil {
		return fmt.Errorf("failed to begin transaction %w", err)
	}
	ctx = context.WithValue(ctx, txKey, tx)

	// Роллбэк в случае ошибки
	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			err = fmt.Errorf("failed to begin transaction %w", err)
		}
	}()

	err = fn(ctx)
	if err == nil {
		return tx.Commit()
	}

	return nil
}

func NewRepository(db *sql.DB) (*Repository, *Transactioner) {
	return &Repository{
			withDB{db},
		}, &Transactioner{
			withDB{db},
		}
}

type DXBT interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
