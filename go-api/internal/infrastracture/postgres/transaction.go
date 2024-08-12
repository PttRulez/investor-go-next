package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
)

func (pg *Transaction) Delete(ctx context.Context, id int, userID int) error {
	const op = "TransactionPostgres.Delete"

	queryString := `DELETE FROM transactions WHERE id = $1 AND user_id = $2;`
	_, err := pg.db.ExecContext(ctx, queryString, id, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *Transaction) GetByID(ctx context.Context, id int, userID int) (*entity.Transaction, error) {
	const op = "TransactionPostgres.GetByID"

	queryString := `SELECT * FROM transactions WHERE id = $1;`
	row := pg.db.QueryRowContext(ctx, queryString, id, userID)

	var t = new(entity.Transaction)
	err := row.Scan(t.ID, t.Amount, t.Date, t.PortfolioID, t.Type)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func (pg *Transaction) GetListByPortfolioID(ctx context.Context, id int, userID int) ([]entity.Transaction, error) {
	const op = "TransactionPostgres.GetListByPortfolioID"

	queryString := `SELECT * FROM transactions WHERE portfolio_id = $1 AND user_id = $2;`
	rows, err := pg.db.QueryContext(ctx, queryString, id, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var ts []entity.Transaction
	for rows.Next() {
		var t entity.Transaction
		err = rows.Scan(&t.ID, &t.Amount, &t.Date, &t.PortfolioID, &t.Type, &t.UserID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		ts = append(ts, t)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return ts, nil
}

func (pg *Transaction) Insert(ctx context.Context, t *entity.Transaction) error {
	const op = "TransactionPostgres.Insert"

	queryString := `INSERT INTO transactions (amount, date, portfolio_id, type, user_id)
		VALUES ($1, $2, $3, $4, $5);`
	_, err := pg.db.ExecContext(ctx, queryString, t.Amount, t.Date, t.PortfolioID, t.Type, t.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

type Transaction struct {
	db *sql.DB
}

func NewTransactionPostgres(db *sql.DB) *Transaction {
	return &Transaction{db: db}
}
