package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
)

func (pg *Transaction) Delete(ctx context.Context, id int, userId int) error {
	queryString := `DELETE FROM transactions WHERE id = $1 AND user_id = $2;`
	_, err := pg.db.ExecContext(ctx, queryString, id, userId)
	if err != nil {
		return fmt.Errorf("[TransactionPostgres.Delete] %w", err)
	}
	return nil
}

func (pg *Transaction) GetById(ctx context.Context, id int, userId int) (*entity.Transaction, error) {
	queryString := `SELECT * FROM transactions WHERE id = $1;`
	row := pg.db.QueryRowContext(ctx, queryString, id, userId)
	if row.Err() != nil {
		return nil, fmt.Errorf("[TransactionPostgres.GetById]: %w", row.Err())
	}

	var t = new(entity.Transaction)
	err := row.Scan(t.Id, t.Amount, t.Date, t.PortfolioId, t.Type)
	if err != nil {
		return nil, fmt.Errorf("[TransactionPostgres.GetById]: %w", err)
	}

	return t, nil
}

func (pg *Transaction) GetListByPortfolioId(ctx context.Context, id int, userId int) ([]entity.Transaction, error) {
	queryString := `SELECT * FROM transactions WHERE portfolio_id = $1 AND user_id = $2;`
	rows, err := pg.db.QueryContext(ctx, queryString, id, userId)
	if err != nil {
		return nil, fmt.Errorf("TransactionPostgres.GetListByPortfolioId -> %w", err)
	}
	defer rows.Close()

	var ts []entity.Transaction
	for rows.Next() {
		var t entity.Transaction
		err := rows.Scan(&t.Id, &t.Amount, &t.Date, &t.PortfolioId, &t.Type, &t.UserId)
		if err != nil {
			return nil, fmt.Errorf("TransactionPostgres.GetListByPortfolioId -> %w", err)
		}
		ts = append(ts, t)
	}

	return ts, nil
}

func (pg *Transaction) Insert(ctx context.Context, t *entity.Transaction) error {
	queryString := `INSERT INTO transactions (amount, date, portfolio_id, type, user_id)
		VALUES ($1, $2, $3, $4, $5);`
	_, err := pg.db.ExecContext(ctx, queryString, t.Amount, t.Date, t.PortfolioId, t.Type, t.UserId)
	if err != nil {
		return fmt.Errorf("[TransactionPostgres.Insert]: %w", err)
	}
	return nil
}

type Transaction struct {
	db *sql.DB
}

func NewTransactionPostgres(db *sql.DB) *Transaction {
	return &Transaction{db: db}
}
