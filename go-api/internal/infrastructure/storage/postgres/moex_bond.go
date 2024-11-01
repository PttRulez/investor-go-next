package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"

	"github.com/lib/pq"
)

func (pg *Repository) GetMoexBond(ctx context.Context, ticker string) (domain.Bond, error) {
	const op = "Repository.GetMoexBond"

	queryString := `SELECT id, board, engine, lotsize, market, name, shortname, ticker
	 FROM moex_bonds WHERE ticker = $1;`

	var bond domain.Bond
	err := pg.db.QueryRowContext(ctx, queryString, ticker).Scan(
		&bond.ID,
		&bond.Board,
		&bond.Engine,
		&bond.LotSize,
		&bond.Market,
		&bond.Name,
		&bond.ShortName,
		&bond.Ticker,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Bond{}, storage.ErrNotFound
	}
	if err != nil {
		return domain.Bond{}, fmt.Errorf("%s: %w", op, err)
	}

	return bond, nil
}

func (pg *Repository) GetMoexBondList(ctx context.Context, ids []int) (
	[]domain.Bond, error) {
	const op = "Repository.GetMoexBondList"

	queryString := "SELECT * FROM moex_bonds WHERE id = ANY($1)"

	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var bonds []domain.Bond

	for rows.Next() {
		var bond domain.Bond
		err = rows.Scan(&bond.ID, &bond.Board, &bond.Engine, &bond.Market, &bond.Name,
			&bond.ShortName, &bond.Ticker)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		bonds = append(bonds, bond)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return bonds, nil
}

func (pg *Repository) InsertMoexBond(ctx context.Context, b domain.Bond) (domain.Bond, error) {
	const op = "Repository.InsertMoexBond"

	querySting := `INSERT INTO moex_bonds (board, coupon_percent, coupon_value, coupon_frequency,
	 currency, engine, face_value, issue_date, lotsize, market, mat_date, name, shortname, ticker)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, board,
		coupon_percent, coupon_value, coupon_frequency, currency, engine, face_value, issue_date,
		lotsize, market, mat_date, name, shortname, ticker;`

	var r domain.Bond
	err := pg.db.QueryRowContext(ctx, querySting, b.Board, b.CouponPercent, b.CouponValue,
		b.CouponFrequency, b.Currency, b.Engine, b.FaceValue, b.IssueDate, b.LotSize, b.Market,
		b.MatDate, b.Name, b.ShortName, b.Ticker).
		Scan(&r.ID, &r.Board, &r.CouponPercent, &r.CouponValue, &r.CouponFrequency, &r.Currency,
			&r.Engine, &r.FaceValue, &r.IssueDate, &r.LotSize, &r.Market, &r.MatDate, &r.Name,
			&r.ShortName, &r.Ticker)
	if err != nil {
		return domain.Bond{}, fmt.Errorf("%s: %w", op, err)
	}

	return r, nil
}
