package postgres

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
)

func (pg *Repository) DeleteDividend(ctx context.Context, id int, userID int) error {
	const op = "Repository.DeleteDividend"

	queryString := `DELETE FROM dividends WHERE id = $1 AND user_id = $2;`
	result, err := pg.db.ExecContext(ctx, queryString, id, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (pg *Repository) GetDividendList(ctx context.Context, portfolioID int) (
	[]domain.Dividend, error) {
	const op = "Repository.GetDividendList"

	queryString := `SELECT date, exchange, id, payment_per_share, payment_period,
		portfolio_id, shares_count, ticker FROM dividends WHERE portfolio_id = $1
		ORDER BY date DESC, id DESC;`

	rows, err := pg.db.QueryContext(ctx, queryString, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var dividends []domain.Dividend
	for rows.Next() {
		var d domain.Dividend
		e := rows.Scan(
			&d.Date,
			&d.Exchange,
			&d.ID,
			&d.PaymentPerShare,
			&d.PaymentPeriod,
			&d.PortfolioID,
			&d.SharesCount,
			&d.Ticker,
		)
		if e != nil {
			return nil, fmt.Errorf("%s: %w", op, e)
		}
		dividends = append(dividends, d)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return dividends, nil
}

func (pg *Repository) InsertDividend(ctx context.Context, d domain.Dividend,
	userID int) error {
	const op = "Repository.InsertDividend"

	queryString := `INSERT INTO dividends (date, exchange, payment_per_share,
    payment_period, portfolio_id, shares_count, ticker, user_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := pg.db.ExecContext(ctx, queryString,
		d.Date,
		d.Exchange,
		d.PaymentPerShare,
		d.PaymentPeriod,
		d.PortfolioID,
		d.SharesCount,
		d.Ticker,
		userID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
