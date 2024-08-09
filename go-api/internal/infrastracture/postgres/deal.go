package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
)

func (pg *DealPostgres) Delete(ctx context.Context, id int, userId int) (*entity.Deal, error) {
	queryString := `DELETE FROM deals WHERE id = $1 AND user_id = $2 RETURNING *;`
	row, err := pg.db.QueryContext(ctx, queryString, id, userId)
	if err != nil {
		return nil, fmt.Errorf("[DealPostgres Delete]: %w", err)
	}

	var d = new(entity.Deal)
	err = row.Scan(
		d.Amount,
		d.Date,
		d.Exchange,
		d.Id,
		d.PortfolioId,
		d.Price,
		d.SecurityType,
		d.Ticker,
		d.Type,
		d.UserId,
	)
	if err != nil {
		return nil, fmt.Errorf("[DealPostgres.Delete]: %w", err)
	}

	return d, nil
}

func (pg *DealPostgres) GetDealListByPortoflioId(ctx context.Context,
	portfolioId int, userId int) ([]*entity.Deal, error) {
	queryString := `SELECT * FROM deals
		WHERE portfolio_id = $1 AND user_id = $2
		ORDER BY date DESC, id DESC;`

	rows, err := pg.db.QueryContext(ctx, queryString, portfolioId, userId)
	if err != nil {
		return nil, fmt.Errorf("[DealPostgres.GetDealsListByPortoflioId]: %w", err)
	}

	var deals []*entity.Deal
	for rows.Next() {
		var deal = entity.Deal{}
		err := rows.Scan(
			&deal.Amount,
			&deal.Date,
			&deal.Exchange,
			&deal.Id,
			&deal.PortfolioId,
			&deal.Price,
			&deal.SecurityType,
			&deal.Ticker,
			&deal.Type,
			&deal.UserId,
		)
		if err != nil {
			return nil, fmt.Errorf("[MoexBondDealPostgres.GetDealsListByPortoflioId]: %w", err)
		}
		deals = append(deals, &deal)
	}
	return deals, nil
}

func (pg *DealPostgres) GetDealListForSecurity(ctx context.Context, exchange entity.Exchange, portfolioId int,
	securityType entity.SecurityType, ticker string) ([]*entity.Deal, error) {
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
		portfolioId,
	)
	if err != nil {
		return nil, fmt.Errorf("[DealPostgres.GetDealsListForSecurity]: %w", err)
	}
	var deals []*entity.Deal
	for rows.Next() {
		var deal = new(entity.Deal)
		err := rows.Scan(
			deal.Amount,
			deal.Date,
			deal.Exchange,
			deal.Id,
			deal.PortfolioId,
			deal.Price,
			deal.SecurityType,
			deal.Ticker,
			deal.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("[DealPostgres.GetDealsListForSecurity]: %w", err)
		}
		deals = append(deals, deal)
	}

	return deals, nil
}

func (pg *DealPostgres) Insert(ctx context.Context, d *entity.Deal) error {
	queryString := `INSERT INTO deals (amount, date, exchange, portfolio_id, price, security_type, ticker, type,
    	user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	_, err := pg.db.ExecContext(
		ctx,
		queryString,
		d.Amount,       // $1
		d.Date,         // $2
		d.Exchange,     // $3
		d.PortfolioId,  // $4
		d.Price,        // $5
		d.SecurityType, // $6
		d.Ticker,       // $7
		d.Type,         // $8
		d.UserId,       // $9
	)
	if err != nil {
		return fmt.Errorf("\n<-[DealPostgres.Insert]: %w", err)
	}
	return nil
}

func (pg *DealPostgres) Update(ctx context.Context, d *entity.Deal) error {
	queryString := `UPDATE deals SET amount = $1, date = $2, exchange = $4, portfolio_id = $5,
    	price = $6, security_type = $7, ticker = $8, type = $9 WHERE id = $10;`
	_, err := pg.db.ExecContext(
		ctx,
		queryString,
		d.Amount,
		d.Date,
		d.Exchange,
		d.PortfolioId,
		d.Price,
		d.SecurityType,
		d.Ticker,
		d.Type,
		d.Id,
	)
	if err != nil {
		return fmt.Errorf("\n<-[DealPostgres.Update]: %w", err)
	}
	return nil
}

type DealPostgres struct {
	db *sql.DB
}

func NewDealPostgres(db *sql.DB) *DealPostgres {
	return &DealPostgres{db: db}
}
