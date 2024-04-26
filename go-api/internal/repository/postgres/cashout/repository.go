package cashout

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/model"
)

func (pg *CashoutPostgres) Delete(ctx context.Context, id int) error {
	queryString := `DELETE FROM cashouts WHERE id = $1;`
	_, err := pg.db.ExecContext(ctx, queryString, id)
	if err != nil {
		return fmt.Errorf("[CashoutPostgres.Delete] %w", err)
	}
	return nil
}

func (pg *CashoutPostgres) GetById(ctx context.Context, id int) (*model.Cashout, error) {
	queryString := `SELECT * FROM cashouts WHERE id = $1;`
	row := pg.db.QueryRowContext(ctx, queryString, id)
	if row.Err() != nil {
		return nil, fmt.Errorf("[CashoutPostgres.GetById]: %w", row.Err())
	}

	var cashout Cashout
	err := row.Scan(&cashout.Id, &cashout.Amount, &cashout.Date, &cashout.PortfolioId)
	if err != nil {
		return nil, fmt.Errorf("[CashoutPostgres.GetById]: %w", err)
	}

	return FromDBToModelCashout(&cashout), nil
}

func (pg *CashoutPostgres) GetListByPortfolioId(ctx context.Context, id int) ([]*model.Cashout, error) {
	queryString := `SELECT * FROM cashouts WHERE portfolio_id = $1;`
	rows, err := pg.db.QueryContext(ctx, queryString, id)
	if err != nil {
		return nil, fmt.Errorf("[CashoutPostgres.GetListByPortfolioId]: %w", err)
	}
	defer rows.Close()

	var cashouts []*model.Cashout
	for rows.Next() {
		var cashout Cashout
		err := rows.Scan(&cashout.Id, &cashout.Amount, &cashout.Date, &cashout.PortfolioId)
		if err != nil {
			return nil, fmt.Errorf("[CashoutPostgres.GetListByPortfolioId]: %w", err)
		}
		cashouts = append(cashouts, FromDBToModelCashout(&cashout))
	}

	return cashouts, nil
}

func (pg *CashoutPostgres) Insert(ctx context.Context, c *model.Cashout) error {
	queryString := `INSERT INTO cashouts (amount, date, portfolio_id) VALUES ($1, $2, $3);`
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
