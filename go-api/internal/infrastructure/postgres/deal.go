package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/database"
)

func (pg *DealPostgres) Delete(ctx context.Context, id int, userID int) (domain.Deal, error) {
	const op = "DealPostgres.Delete"

	// Удаляем через QueryRowContext т.к нам в сервисе нужна полная инфа по сделке, чтобы
	// пересчитать позицию, куда входила сделка
	queryString := `DELETE FROM deals WHERE id = $1 AND user_id = $2 RETURNING *;`
	row := pg.db.QueryRowContext(ctx, queryString, id, userID)

	var d domain.Deal

	err := row.Scan(
		d.Amount,
		d.Commission,
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
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Deal{}, database.ErrNotFound
	}
	if err != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, err)
	}
	if row.Err() != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, row.Err())
	}

	return d, nil
}

func (pg *DealPostgres) GetDealListByPortoflioID(ctx context.Context,
	portfolioID int, userID int) ([]domain.Deal, error) {
	const op = "DealPostgres.GetDealListByPortoflioID"

	queryString := `SELECT * FROM deals
		WHERE portfolio_id = $1 AND user_id = $2
		ORDER BY date DESC, id DESC;`

	rows, err := pg.db.QueryContext(ctx, queryString, portfolioID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var deals []domain.Deal
	for rows.Next() {
		var deal = domain.Deal{}
		e := rows.Scan(
			&deal.Amount,
			&deal.Commission,
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
		deals = append(deals, deal)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return deals, nil
}

func (pg *DealPostgres) GetDealListForSecurity(ctx context.Context, exchange domain.Exchange, portfolioID int,
	securityType domain.SecurityType, ticker string) ([]domain.Deal, error) {
	const op = "DealPostgres.GetDealListForSecurity"

	queryString := `SELECT amount, commission, date, exchange, id, portfolio_id, price,
		security_type, ticker, type
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
		return nil, fmt.Errorf("%s (QueryContext): %w", op, err)
	}
	defer rows.Close()

	var deals []domain.Deal
	for rows.Next() {
		var deal domain.Deal
		e := rows.Scan(
			&deal.Amount,
			&deal.Commission,
			&deal.Date,
			&deal.Exchange,
			&deal.ID,
			&deal.PortfolioID,
			&deal.Price,
			&deal.SecurityType,
			&deal.Ticker,
			&deal.Type,
		)
		if e != nil {
			return nil, fmt.Errorf("%s (rows.Scan): %w", op, e)
		}
		deals = append(deals, deal)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s (rows.Err()): %w", op, rows.Err())
	}

	return deals, nil
}

func (pg *DealPostgres) Insert(ctx context.Context, d domain.Deal) (domain.Deal, error) {
	const op = "DealPostgres.Insert"

	queryString := `INSERT INTO deals (amount, commission, date, exchange, portfolio_id, price,
		security_type, ticker, type, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING amount, commission, date, exchange, id, portfolio_id, price, security_type,
		ticker, type;`

	var deal domain.Deal
	err := pg.db.QueryRowContext(
		ctx,
		queryString,
		d.Amount,
		d.Commission,
		d.Date,
		d.Exchange,
		d.PortfolioID,
		d.Price,
		d.SecurityType,
		d.Ticker,
		d.Type,
		d.UserID,
	).Scan(
		&deal.Amount,
		&deal.Commission,
		&deal.Date,
		&deal.Exchange,
		&deal.ID,
		&deal.PortfolioID,
		&deal.Price,
		&deal.SecurityType,
		&deal.Ticker,
		&deal.Type,
	)
	if err != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	return deal, nil
}

func (pg *DealPostgres) Update(ctx context.Context, d domain.Deal) (domain.Deal, error) {
	const op = "DealPostgres.Update"

	queryString := `UPDATE deals SET amount = $1, date = $2, exchange = $4, portfolio_id = $5,
		price = $6, security_type = $7, ticker = $8, type = $9 WHERE id = $10
		RETURNING amount, date, exchange, portfolio_id, price, security_type, ticker, type;`

	var deal domain.Deal
	err := pg.db.QueryRowContext(
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
	).Scan(
		&deal.Amount,       // $1
		&deal.Date,         // $2
		&deal.Exchange,     // $3
		&deal.PortfolioID,  // $4
		&deal.Price,        // $5
		&deal.SecurityType, // $6
		&deal.Ticker,       // $7
		&deal.Type,         // $8
	)
	if err != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	return deal, nil
}

type DealPostgres struct {
	db *sql.DB
}

func NewDealPostgres(db *sql.DB) *DealPostgres {
	return &DealPostgres{db: db}
}
