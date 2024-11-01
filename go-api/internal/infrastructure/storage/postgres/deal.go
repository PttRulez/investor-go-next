package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
)

func (pg *Repository) DeleteDeal(ctx context.Context, id int, userID int) (domain.Deal, error) {
	const op = "Repository.DeleteDeal"

	// Удаляем через QueryRowContext т.к нам в сервисе нужна полная инфа по сделке, чтобы
	// пересчитать позицию, куда входила сделка
	queryString := `DELETE FROM deals WHERE id = $1 AND user_id = $2 RETURNING amount,
		commission, date, exchange, id, portfolio_id, price, security_type,
		ticker, type, nkd, shortname;`
	row := pg.t(ctx).QueryRowContext(ctx, queryString, id, userID)

	var d domain.Deal

	err := row.Scan(
		&d.Amount,
		&d.Commission,
		&d.Date,
		&d.Exchange,
		&d.ID,
		&d.PortfolioID,
		&d.Price,
		&d.SecurityType,
		&d.Ticker,
		&d.Type,
		&d.Nkd,
		&d.ShortName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Deal{}, storage.ErrNotFound
	}
	if err != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, err)
	}
	if row.Err() != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, row.Err())
	}

	return d, nil
}

func (pg *Repository) GetDealList(ctx context.Context,
	portfolioID int, userID int) ([]domain.Deal, error) {
	const op = "Repository.GetDealList"

	queryString := `
	SELECT amount, commission, date, exchange, id, portfolio_id, price,
		security_type, ticker, type, nkd, shortname
		FROM deals
		WHERE portfolio_id = $1 AND user_id = $2
		ORDER BY date DESC, id DESC;`

	rows, err := pg.t(ctx).QueryContext(ctx, queryString, portfolioID, userID)
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
			&deal.Nkd,
			&deal.ShortName,
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

func (pg *Repository) GetDealListForSecurity(ctx context.Context, exchange domain.Exchange, portfolioID int,
	securityType domain.SecurityType, ticker string) ([]domain.Deal, error) {
	const op = "Repository.GetDealListForSecurity"

	queryString := `SELECT amount, commission, date, exchange, id, portfolio_id, price,
		security_type, ticker, type, nkd, shortname
		FROM deals d 
		WHERE d.exchange = $1 AND d.security_type = $2 AND d.ticker = $3 AND d.portfolio_id = $4
		ORDER BY d.date DESC, d.id DESC;`

	rows, err := pg.t(ctx).QueryContext(
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
			&deal.Nkd,
			&deal.ShortName,
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

func (pg *Repository) InsertDeal(ctx context.Context, d domain.Deal) (domain.Deal, error) {
	const op = "Repository.InsertDeal"
	var err error

	queryString := `INSERT INTO deals (amount, commission, date, exchange, nkd, portfolio_id, price,
		security_type, ticker, type, user_id, shortname) VALUES ($1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12)
		RETURNING amount, commission, date, exchange, id, nkd, portfolio_id, price, security_type,
			shortname, ticker, type;`

	var deal domain.Deal
	err = pg.t(ctx).QueryRowContext(
		ctx,
		queryString,
		d.Amount,
		d.Commission,
		d.Date,
		d.Exchange,
		d.Nkd,
		d.PortfolioID,
		d.Price,
		d.SecurityType,
		d.Ticker,
		d.Type,
		d.UserID,
		d.ShortName,
	).Scan(
		&deal.Amount,
		&deal.Commission,
		&deal.Date,
		&deal.Exchange,
		&deal.ID,
		&deal.Nkd,
		&deal.PortfolioID,
		&deal.Price,
		&deal.SecurityType,
		&deal.Ticker,
		&deal.Type,
		&deal.ShortName,
	)
	if err != nil {
		fmt.Println("LALA")
		return domain.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	return deal, nil
}

func (pg *Repository) UpdateDeal(ctx context.Context, d domain.Deal) (domain.Deal, error) {
	const op = "Repository.UpdateDeal"

	queryString := `UPDATE deals SET amount = $1, date = $2, exchange = $4, portfolio_id = $5,
		price = $6, security_type = $7, ticker = $8, type = $9, nkd = $10, shortname = $11
		WHERE id = $12
		RETURNING amount, date, exchange, portfolio_id, price, security_type, ticker, type,
			nkd, shortname;`

	var deal domain.Deal
	err := pg.t(ctx).QueryRowContext(
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
		d.Nkd,
		d.ShortName,
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
		&deal.Nkd,          // $8
		&deal.ShortName,    // $8
	)
	if err != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	return deal, nil
}
