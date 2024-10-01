package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/database"

	"github.com/lib/pq"
)

func (pg *MoexSharesPostgres) GetByTicker(ctx context.Context, ticker string) (
	domain.Share, error) {
	const op = "MoexSharesPostgres.GetByTicker"

	querySting := `SELECT id, board, engine, lotsize, market, name, price_decimals,
		shortname, ticker FROM moex_shares WHERE ticker = $1;`

	row := pg.db.QueryRowContext(ctx, querySting, ticker)

	var share domain.Share
	err := row.Scan(
		&share.ID,
		&share.Board,
		&share.Engine,
		&share.LotSize,
		&share.Market,
		&share.Name,
		&share.PriceDecimals,
		&share.ShortName,
		&share.Ticker,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Share{}, database.ErrNotFound
	}
	if err != nil {
		return domain.Share{}, fmt.Errorf("%s: %w", op, err)
	}

	return share, nil
}

func (pg *MoexSharesPostgres) GetListByIDs(ctx context.Context, ids []int) (
	[]domain.Share, error) {
	const op = "MoexSharesPostgres.GetListByIDs"

	queryString := `SELECT id, board, engine, lotsize, market, name, price_decimals,
		shortname, ticker FROM moex_shares WHERE id = ANY($1)`

	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var shares []domain.Share
	for rows.Next() {
		var share domain.Share
		err = rows.Scan(&share.ID, &share.Board, &share.Engine, &share.Market, &share.Name,
			&share.PriceDecimals, &share.ShortName, &share.Ticker)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		shares = append(shares, share)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return shares, nil
}

func (pg *MoexSharesPostgres) Insert(ctx context.Context, s domain.Share) (domain.Share, error) {
	const op = "MoexSharesPostgres.Insert"

	querySting := `INSERT INTO moex_shares (board, engine, lotsize, market, name, price_decimals,
		shortname, ticker) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING
		id, board, engine, lotsize, market, name, price_decimals, shortname, ticker;`

	var r domain.Share
	err := pg.db.QueryRowContext(ctx, querySting, s.Board, s.Engine, s.LotSize,
		s.Market, s.Name, &s.PriceDecimals, s.ShortName, s.Ticker).
		Scan(&r.ID, &r.Board, &r.Engine, &r.LotSize, &r.Market, &r.Name, &r.PriceDecimals,
			&r.ShortName, &r.Ticker)
	if err != nil {
		return domain.Share{}, fmt.Errorf("%s: %w", op, err)
	}

	return r, nil
}

type MoexSharesPostgres struct {
	db *sql.DB
}

func NewMoexSharesPostgres(db *sql.DB) *MoexSharesPostgres {
	return &MoexSharesPostgres{db: db}
}
