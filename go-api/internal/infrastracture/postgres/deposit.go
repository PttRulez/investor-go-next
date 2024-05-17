package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
)

func (pg *DepositPostgres) Delete(ctx context.Context, id int) error {
	queryString := `DELETE FROM deposits WHERE id = $1;`
	_, err := pg.db.ExecContext(ctx, queryString, id)
	if err != nil {
		return fmt.Errorf("[DepositPostgres.Delete] %w", err)
	}
	return nil
}

func (pg *DepositPostgres) GetById(ctx context.Context, id int) (*entity.Deposit, error) {
	queryString := `SELECT * FROM deposits WHERE id = $1;`
	row := pg.db.QueryRowContext(ctx, queryString, id)
	if row.Err() != nil {
		return nil, fmt.Errorf("[DepositPostgres.GetById]: %w", row.Err())
	}

	var deposit entity.Deposit
	err := row.Scan(&deposit.Id, &deposit.Amount, &deposit.Date, &deposit.PortfolioId)
	if err != nil {
		return nil, fmt.Errorf("[DepositPostgres.GetById]: %w", err)
	}

	return &deposit, nil
}

func (pg *DepositPostgres) GetListByPortfolioId(ctx context.Context, id int) ([]*entity.Deposit, error) {
	queryString := `SELECT * FROM deposits WHERE portfolio_id = $1;`
	rows, err := pg.db.QueryContext(ctx, queryString, id)
	if err != nil {
		return nil, fmt.Errorf("[DepositPostgres.GetListByPortfolioId]: %w", err)
	}
	defer rows.Close()

	var deposits []*entity.Deposit
	for rows.Next() {
		var deposit entity.Deposit
		err := rows.Scan(&deposit.Id, &deposit.Amount, &deposit.Date, &deposit.PortfolioId)
		if err != nil {
			return nil, fmt.Errorf("[DepositPostgres.GetListByPortfolioId]: %w", err)
		}
		deposits = append(deposits, &deposit)
	}

	return deposits, nil
}

func (pg *DepositPostgres) Insert(ctx context.Context, c *entity.Deposit) error {
	queryString := `INSERT INTO deposits (amount, date, portfolio_id) VALUES ($1, $2, $3) RETURNING *;`
	_, err := pg.db.ExecContext(ctx, queryString, c.Amount, c.Date, c.PortfolioId)
	if err != nil {
		return fmt.Errorf("[DepositPostgres.Insert]: %w", err)
	}
	return nil
}

type DepositPostgres struct {
	db *sql.DB
}

func NewDepositPostgres(db *sql.DB) *DepositPostgres {
	return &DepositPostgres{db: db}
}
