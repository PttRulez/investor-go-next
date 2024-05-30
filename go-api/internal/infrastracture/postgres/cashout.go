package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
)

func (pg *Cashout) Delete(ctx context.Context, id int, userId int) error {
	queryString := `DELETE FROM cashouts WHERE id = $1 AND user_id = $2;`
	_, err := pg.db.ExecContext(ctx, queryString, id, userId)
	if err != nil {
		return fmt.Errorf("[postgres.Cashout.Delete] %w", err)
	}
	return nil
}

func (pg *Cashout) GetById(ctx context.Context, id int, userId int) (*entity.Cashout, error) {
	queryString := `SELECT * FROM cashouts WHERE id = $1;`
	row := pg.db.QueryRowContext(ctx, queryString, id, userId)
	if row.Err() != nil {
		return nil, fmt.Errorf("[CashoutPostgres.GetById]: %w", row.Err())
	}

	var c entity.Cashout
	err := row.Scan(&c.Id, &c.Amount, &c.Date, &c.PortfolioId)
	if err != nil {
		return nil, fmt.Errorf("[CashoutPostgres.GetById]: %w", err)
	}

	return &c, nil
}

func (pg *Cashout) GetListByPortfolioId(ctx context.Context, id int) ([]entity.Cashout, error) {
	queryString := `SELECT * FROM cashouts WHERE portfolio_id = $1 AND user_id = $2;`
	rows, err := pg.db.QueryContext(ctx, queryString, id)
	if err != nil {
		return nil, fmt.Errorf("[CashoutPostgres.GetListByPortfolioId]: %w", err)
	}
	defer rows.Close()

	var cashouts []entity.Cashout
	for rows.Next() {
		var c entity.Cashout
		err := rows.Scan(&c.Id, &c.Amount, &c.Date, &c.PortfolioId)
		if err != nil {
			return nil, fmt.Errorf("[CashoutPostgres.GetListByPortfolioId]: %w", err)
		}
		cashouts = append(cashouts, c)
	}

	return cashouts, nil
}

func (pg *Cashout) Insert(ctx context.Context, c *entity.Cashout) error {
	queryString := `INSERT INTO cashouts (amount, date, portfolio_id) VALUES ($1, $2, $3);`
	_, err := pg.db.ExecContext(ctx, queryString, c.Amount, c.Date, c.PortfolioId)
	if err != nil {
		return fmt.Errorf("[CashoutPostgres.Insert]: %w", err)
	}
	return nil
}

type Cashout struct {
	db *sql.DB
}

func NewCashoutPostgres(db *sql.DB) *Cashout {
	return &Cashout{db: db}
}
