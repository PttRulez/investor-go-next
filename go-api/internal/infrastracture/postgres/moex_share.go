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

func (pg *MoexSharesPostgres) GetBySecid(ctx context.Context, secid string) (
	entity.Share, error) {
	const op = "MoexSharesPostgres.GetBySecid"

	querySting := `SELECT id, board, engine, lotsize, market, name, price_decimals,
		shortname, secid FROM moex_shares WHERE secid = $1;`

	row := pg.db.QueryRowContext(ctx, querySting, secid)

	var share entity.Share
	err := row.Scan(
		&share.ID,
		&share.Board,
		&share.Engine,
		&share.LotSize,
		&share.Market,
		&share.Name,
		&share.PriceDecimals,
		&share.ShortName,
		&share.Secid,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Share{}, database.ErrNotFound
	}
	if err != nil {
		return entity.Share{}, fmt.Errorf("%s: %w", op, err)
	}

	return share, nil
}

func (pg *MoexSharesPostgres) GetListByIDs(ctx context.Context, ids []int) (
	[]entity.Share, error) {
	const op = "MoexSharesPostgres.GetListByIDs"

	queryString := `SELECT id, board, engine, lotsize, market, name, price_decimals,
		shortname, secid FROM moex_shares WHERE id = ANY($1)`

	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var shares []entity.Share
	for rows.Next() {
		var share entity.Share
		err = rows.Scan(&share.ID, &share.Board, &share.Engine, &share.Market, &share.Name,
			&share.PriceDecimals, &share.ShortName, &share.Secid)
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

func (pg *MoexSharesPostgres) Insert(ctx context.Context, s entity.Share) (entity.Share, error) {
	const op = "MoexSharesPostgres.Insert"

	querySting := `INSERT INTO moex_shares (board, engine, lotsize, market, name, price_decimals,
		shortname, secid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING
		id, board, engine, lotsize, market, name, price_decimals, shortname, secid;`

	var r entity.Share
	err := pg.db.QueryRowContext(ctx, querySting, s.Board, s.Engine, s.LotSize,
		s.Market, s.Name, &s.PriceDecimals, s.ShortName, s.Secid).
		Scan(&r.ID, &r.Board, &r.Engine, &r.LotSize, &r.Market, &r.Name, &r.PriceDecimals,
			&r.ShortName, &r.Secid)
	if err != nil {
		return entity.Share{}, fmt.Errorf("%s: %w", op, err)
	}

	return r, nil
}

type MoexSharesPostgres struct {
	db *sql.DB
}

func NewMoexSharesPostgres(db *sql.DB) *MoexSharesPostgres {
	return &MoexSharesPostgres{db: db}
}
