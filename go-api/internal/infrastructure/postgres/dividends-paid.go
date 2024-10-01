package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/infrastructure/database"
)

func (pg *DividendsPaidPostgres) Delete(ctx context.Context, id int, userID int) error {
	const op = "DividendsPaidPostgres.Delete"

	queryString := `DELETE FROM dividends_paid WHERE id = $1 AND user_id = $2;`
	result, err := pg.db.ExecContext(ctx, queryString, id, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}

type DividendsPaidPostgres struct {
	db *sql.DB
}

func NewDividendsPaidPostgres(db *sql.DB) *DividendsPaidPostgres {
	return &DividendsPaidPostgres{db: db}
}
