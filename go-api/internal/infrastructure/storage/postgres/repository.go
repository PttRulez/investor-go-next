package postgres

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) ExecAsTransaction(ctx context.Context, fn func(ctx context.Context,
	tx *sql.Tx) error) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	if err := fn(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
