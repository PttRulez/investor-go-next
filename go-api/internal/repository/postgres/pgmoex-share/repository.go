package pgmoexshare

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pttrulez/investor-go/internal/model"
	"github.com/pttrulez/investor-go/internal/repository"
)

type MoexSharesPostgres struct {
	db *sql.DB
}

func NewMoexSharesPostgres(db *sql.DB) repository.MoexShareRepository {
	return &MoexSharesPostgres{db: db}
}

func (pg *MoexSharesPostgres) GetByTicker(ctx context.Context, ticker string) (*model.MoexShare, error) {
	querySting := `SELECT * FROM moex_shares WHERE ticker = $1;`

	row := pg.db.QueryRowContext(ctx, querySting, ticker)
	if row.Err() != nil {
		return nil, fmt.Errorf("[MoexSharesPostgres GetByTicker]: %w", row.Err())
	}

	var share MoexShare
	err := row.Scan(
		&share.Id,
		&share.Board,
		&share.Engine,
		&share.Market,
		&share.Name,
		&share.ShortName,
		&share.Ticker,
	)
	if err != nil {
		return nil, fmt.Errorf("[MoexSharesPostgres GetByTicker]: %w", err)
	}

	return FromDBToMoexShare(&share), nil
}

func (pg *MoexSharesPostgres) GetListByIds(ctx context.Context, ids []int) ([]*model.MoexShare, error) {
	queryString := "SELECT * FROM moex_shares WHERE id = ANY($1)"
	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))

	if err != nil {
		return nil, fmt.Errorf("[MoexSharesPostgres GetListByIds]: %w", err)
	}

	var shares []*model.MoexShare
	for rows.Next() {
		var share MoexShare
		err := rows.Scan(&share.Id, &share.Board, &share.Engine, &share.Market, &share.Name, &share.ShortName, &share.Ticker)
		if err != nil {
			return nil, fmt.Errorf("[MoexSharesPostgres GetListByIds]: %w", err)
		}
		shares = append(shares, FromDBToMoexShare(&share))
	}

	return shares, nil
}

func (pg *MoexSharesPostgres) Insert(ctx context.Context, share *model.MoexShare) error {
	querySting := `INSERT INTO moex_shares (board, engine, market, name, shortname, ticker)
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`

	_, err := pg.db.ExecContext(ctx, querySting, share.Board, share.Engine, share.Market, share.Name, share.ShortName, share.Ticker)
	if err != nil {
		return fmt.Errorf("[MoexSharesPostgres Insert]: %w", err)
	}

	return nil
}
