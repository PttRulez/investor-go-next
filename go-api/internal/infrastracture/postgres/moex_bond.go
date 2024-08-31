package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"

	"github.com/lib/pq"
)

func (pg *MoexBondsPostgres) GetByTicker(ctx context.Context, ticker string) (entity.Bond, error) {
	const op = "MoexBondsPostgres.GetByTicker"

	queryString := `SELECT id, board, engine, lotsize, market, name, shortname, ticker
	 FROM moex_bonds WHERE ticker = $1;`

	var bond entity.Bond
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
		return entity.Bond{}, database.ErrNotFound
	}
	if err != nil {
		return entity.Bond{}, fmt.Errorf("%s: %w", op, err)
	}

	return bond, nil
}

func (pg *MoexBondsPostgres) GetListByIDs(ctx context.Context, ids []int) (
	[]entity.Bond, error) {
	const op = "MoexBondsPostgres.GetListByIDs"

	queryString := "SELECT * FROM moex_bonds WHERE id = ANY($1)"

	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var bonds []entity.Bond

	for rows.Next() {
		var bond entity.Bond
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

func (pg *MoexBondsPostgres) Insert(ctx context.Context, b entity.Bond) (entity.Bond, error) {
	const op = "MoexBondsPostgres.Insert"

	querySting := `INSERT INTO moex_bonds (board, coupon_percent, coupon_value, coupon_frequency,
	 engine, face_value, issue_date, lotsize, market, mat_date, name, shortname, ticker)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id, board,
		coupon_percent, coupon_value, coupon_frequency, engine, face_value, issue_date,
		lotsize, market, mat_date, name, shortname, ticker;`

	var r entity.Bond
	err := pg.db.QueryRowContext(ctx, querySting, b.Board, b.CouponPercent, b.CouponValue,
		b.CouponFrequency, b.Engine, b.FaceValue, b.IssueDate, b.LotSize, b.Market,
		b.MatDate, b.Name, b.ShortName, b.Ticker).
		Scan(&r.ID, &r.Board, &r.CouponPercent, &r.CouponValue, &r.CouponFrequency, &r.Engine,
			&r.FaceValue, &r.IssueDate, &r.LotSize, &r.Market, &r.MatDate, &r.Name, &r.ShortName,
			&r.Ticker)
	if err != nil {
		return entity.Bond{}, fmt.Errorf("%s: %w", op, err)
	}

	return r, nil
}

type MoexBondsPostgres struct {
	db *sql.DB
}

func NewMoexBondsPostgres(db *sql.DB) *MoexBondsPostgres {
	return &MoexBondsPostgres{db: db}
}
