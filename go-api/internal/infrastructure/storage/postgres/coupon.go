package postgres

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
)

func (pg *Repository) DeleteCoupon(ctx context.Context, id int, userID int) error {
	const op = "Repository.DeleteCoupon"

	queryString := `DELETE FROM coupons WHERE id = $1 AND user_id = $2;`
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

func (pg *Repository) GetCouponList(ctx context.Context, portfolioID int) (
	[]domain.Coupon, error) {
	const op = "Repository.GetCouponList"

	queryString := `SELECT bonds_count, coupon_amount, date, exchange, id, payment_period,
		portfolio_id, ticker FROM coupons WHERE portfolio_id = $1 ORDER BY date DESC, id DESC;`

	rows, err := pg.db.QueryContext(ctx, queryString, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var coupons []domain.Coupon
	for rows.Next() {
		var c domain.Coupon
		e := rows.Scan(
			&c.BondsCount,
			&c.CouponAmount,
			&c.Date,
			&c.Exchange,
			&c.ID,
			&c.PaymentPeriod,
			&c.PortfolioID,
			&c.Ticker,
		)
		if e != nil {
			return nil, fmt.Errorf("%s: %w", op, e)
		}
		coupons = append(coupons, c)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return coupons, nil
}

func (pg *Repository) InsertCoupon(ctx context.Context, c domain.Coupon,
	userID int) error {
	const op = "Repository.InsertCoupon"

	queryString := `INSERT INTO coupons (bonds_count, coupon_amount, date, exchange,
    payment_period, portfolio_id, ticker, user_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := pg.db.ExecContext(ctx, queryString,
		c.BondsCount,
		c.CouponAmount,
		c.Date,
		c.Exchange,
		c.PaymentPeriod,
		c.PortfolioID,
		c.Ticker,
		userID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
