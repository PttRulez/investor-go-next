package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
)

func (pg *DealPostgres) Delete(ctx context.Context, id int, userID int) (*entity.Deal, error) {
	const op = "DealPostgres.Delete"

	queryString := `DELETE FROM deals WHERE id = $1 AND user_id = $2 RETURNING *;`
	row := pg.db.QueryRowContext(ctx, queryString, id, userID)

	var d = new(entity.Deal)
	err := row.Scan(
		d.Amount,
		d.Date,
		d.Exchange,
		d.ID,
		d.PortfolioID,
		d.Price,
		d.SecurityType,
		d.Ticker,
		d.Type,
		d.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return d, nil
}

func (pg *DealPostgres) GetDealListByPortoflioID(ctx context.Context,
	portfolioID int, userID int) ([]*entity.Deal, error) {
	const op = "DealPostgres.GetDealListByPortoflioID"

	queryString := `SELECT * FROM deals
		WHERE portfolio_id = $1 AND user_id = $2
		ORDER BY date DESC, id DESC;`

	rows, err := pg.db.QueryContext(ctx, queryString, portfolioID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var deals []*entity.Deal
	for rows.Next() {
		var deal = entity.Deal{}
		e := rows.Scan(
			&deal.Amount,
			&deal.Date,
			&deal.Exchange,
			&deal.ID,
			&deal.PortfolioID,
			&deal.Price,
			&deal.SecurityType,
			&deal.Ticker,
			&deal.Type,
			&deal.UserID,
		)
		if e != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		deals = append(deals, &deal)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return deals, nil
}

func (pg *DealPostgres) GetDealListForSecurity(ctx context.Context, exchange entity.Exchange, portfolioID int,
	securityType entity.SecurityType, ticker string) ([]*entity.Deal, error) {
	const op = "DealPostgres.GetDealListForSecurity"

	queryString := `SELECT d.*
		FROM deals d 
		WHERE d.exchange = $1 AND d.security_type = $2 AND d.ticker = $3 AND d.portfolio_id = $4
		ORDER BY d.date DESC, d.id DESC;`

	rows, err := pg.db.QueryContext(
		ctx,
		queryString,
		exchange,
		securityType,
		ticker,
		portfolioID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var deals []*entity.Deal
	for rows.Next() {
		var deal = new(entity.Deal)
		e := rows.Scan(
			deal.Amount,
			deal.Date,
			deal.Exchange,
			deal.ID,
			deal.PortfolioID,
			deal.Price,
			deal.SecurityType,
			deal.Ticker,
			deal.Type,
		)
		if e != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		deals = append(deals, deal)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return deals, nil
}

func (pg *DealPostgres) Insert(ctx context.Context, d *entity.Deal) error {
	const op = "DealPostgres.Insert"

	queryString := `INSERT INTO deals (amount, date, exchange, portfolio_id, price, security_type, ticker, type,
    	user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	_, err := pg.db.ExecContext(
		ctx,
		queryString,
		d.Amount,       // $1
		d.Date,         // $2
		d.Exchange,     // $3
		d.PortfolioID,  // $4
		d.Price,        // $5
		d.SecurityType, // $6
		d.Ticker,       // $7
		d.Type,         // $8
		d.UserID,       // $9
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (pg *DealPostgres) Update(ctx context.Context, d *entity.Deal) error {
	const op = "DealPostgres.Update"

	queryString := `UPDATE deals SET amount = $1, date = $2, exchange = $4, portfolio_id = $5,
    	price = $6, security_type = $7, ticker = $8, type = $9 WHERE id = $10;`
	_, err := pg.db.ExecContext(
		ctx,
		queryString,
		d.Amount,
		d.Date,
		d.Exchange,
		d.PortfolioID,
		d.Price,
		d.SecurityType,
		d.Ticker,
		d.Type,
		d.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

type DealPostgres struct {
	db *sql.DB
}

func NewDealPostgres(db *sql.DB) *DealPostgres {
	return &DealPostgres{db: db}
}
