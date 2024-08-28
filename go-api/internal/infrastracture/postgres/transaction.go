package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
)

func (pg *Transaction) Delete(ctx context.Context, id int, userID int) error {
	const op = "TransactionPostgres.Delete"

	queryString := `DELETE FROM transactions WHERE id = $1 AND user_id = $2;`
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

func (pg *Transaction) GetByID(ctx context.Context, id int, userID int) (entity.Transaction, error) {
	const op = "TransactionPostgres.GetByID"

	queryString := `SELECT * FROM transactions WHERE id = $1 AND user_id = $2;`
	row := pg.db.QueryRowContext(ctx, queryString, id, userID)

	var t entity.Transaction
	err := row.Scan(&t.ID, &t.Amount, &t.Date, &t.PortfolioID, &t.Type)
	if err != nil {
		return entity.Transaction{}, fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func (pg *Transaction) GetListByPortfolioID(ctx context.Context, id int, userID int) (
	[]entity.Transaction, error) {
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

func (pg *Transaction) Insert(ctx context.Context, t entity.Transaction) (entity.Transaction, error) {
	const op = "TransactionPostgres.Insert"

	queryString := `INSERT INTO transactions (amount, date, portfolio_id, type, user_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, amount, date, portfolio_id, type;`

	var res entity.Transaction
	err := pg.db.QueryRowContext(ctx, queryString, t.Amount, t.Date, t.PortfolioID, t.Type,
		t.UserID).
		Scan(&res.ID, &res.Amount, &res.Date, &res.PortfolioID, &res.Type)
	if err != nil {
		return entity.Transaction{}, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

type Transaction struct {
	db *sql.DB
}

func NewTransactionPostgres(db *sql.DB) *Transaction {
	return &Transaction{db: db}
}
