package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pttrulez/investor-go/internal/types"
	tmoex "github.com/pttrulez/investor-go/internal/types/moex"
)

type MoexSharesPostgres struct {
	db *sql.DB
}

func NewMoexSharesPostgres(db *sql.DB) types.MoexShareRepository {
	return &MoexSharesPostgres{db: db}
}

func (pg *MoexSharesPostgres) GetByTicker(ctx context.Context, ticker string) (*tmoex.Share, error) {
	querySting := `SELECT * FROM moex_shares WHERE ticker = $1;`

	row := pg.db.QueryRowContext(ctx, querySting, ticker)
	if row.Err() != nil {
		return nil, fmt.Errorf("[MoexSharesPostgres GetByTicker]: %w", row.Err())
	}

	var share tmoex.Share
	err := row.Scan(&share.Id, &share.Board, &share.Engine, &share.Market, &share.Name, &share.ShortName, &share.Ticker)
	if err != nil {
		return nil, fmt.Errorf("[MoexSharesPostgres GetByTicker]: %w", err)
	}

	return &share, nil
}

func (pg *MoexSharesPostgres) GetListByIds(ctx context.Context, ids []int) ([]*tmoex.Share, error) {
	queryString := "SELECT * FROM moex_shares WHERE id = ANY($1)"
	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))

	if err != nil {
		return nil, fmt.Errorf("[MoexSharesPostgres GetListByIds]: %w", err)
	}

	var shares []*tmoex.Share
	for rows.Next() {
		var share tmoex.Share
		err := rows.Scan(&share.Id, &share.Board, &share.Engine, &share.Market, &share.Name, &share.ShortName, &share.Ticker)
		if err != nil {
			return nil, fmt.Errorf("[MoexSharesPostgres GetListByIds]: %w", err)
		}
		shares = append(shares, &share)
	}

	return shares, nil
}

func (pg *MoexSharesPostgres) Insert(ctx context.Context, share *tmoex.Share) error {
	querySting := `INSERT INTO moex_shares (board, engine, market, name, shortname, ticker)
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`

	_, err := pg.db.ExecContext(ctx, querySting, share.Board, share.Engine, share.Market, share.Name, share.ShortName, share.Ticker)
	if err != nil {
		return fmt.Errorf("[MoexSharesPostgres Insert]: %w", err)
	}

	return nil
}
