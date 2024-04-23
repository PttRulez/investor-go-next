package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/types"
)

func (pg *CashoutPostgres) Delete(ctx context.Context, id int) error {
	queryString := `DELETE FROM cashouts WHERE id = $1;`
	_, err := pg.db.ExecContext(ctx, queryString, id)
	if err != nil {
		return fmt.Errorf("[CashoutPostgres.Delete] %w", err)
	}
	return nil
}

func (pg *CashoutPostgres) GetById(ctx context.Context, id int) (*types.Cashout, error) {
	queryString := `SELECT * FROM cashouts WHERE id = $1;`
	row := pg.db.QueryRowContext(ctx, queryString, id)
	if row.Err() != nil {
		return nil, fmt.Errorf("[CashoutPostgres.GetById]: %w", row.Err())
	}

	var cashout types.Cashout
	err := row.Scan(&cashout.Id, &cashout.Amount, &cashout.Date, &cashout.PortfolioId)
	if err != nil {
		return nil, fmt.Errorf("[CashoutPostgres.GetById]: %w", err)
	}

	return &cashout, nil
}

func (pg *CashoutPostgres) GetListByPortfolioId(ctx context.Context, id int) ([]*types.Cashout, error) {
	queryString := `SELECT * FROM cashouts WHERE portfolio_id = $1;`
	rows, err := pg.db.QueryContext(ctx, queryString, id)
	if err != nil {
		return nil, fmt.Errorf("[CashoutPostgres.GetListByPortfolioId]: %w", err)
	}
	defer rows.Close()

	var cashouts []*types.Cashout
	for rows.Next() {
		var cashout types.Cashout
		err := rows.Scan(&cashout.Id, &cashout.Amount, &cashout.Date, &cashout.PortfolioId)
		if err != nil {
			return nil, fmt.Errorf("[CashoutPostgres.GetListByPortfolioId]: %w", err)
		}
		cashouts = append(cashouts, &cashout)
	}

	return cashouts, nil
}

func (pg *CashoutPostgres) Insert(ctx context.Context, c *types.Cashout) error {
	queryString := `INSERT INTO cashouts (amount, date, portfolio_id) VALUES ($1, $2, $3) RETURNING *;`
	_, err := pg.db.ExecContext(ctx, queryString, c.Amount, c.Date, c.PortfolioId)
	if err != nil {
		return fmt.Errorf("[CashoutPostgres.Insert]: %w", err)
	}
	return nil
}

type CashoutPostgres struct {
	db *sql.DB
}

func NewCashoutPostgres(db *sql.DB) *CashoutPostgres {
	return &CashoutPostgres{db: db}
}
