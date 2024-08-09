package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"

	"github.com/lib/pq"
)

func (pg *MoexSharesPostgres) GetBySecid(ctx context.Context, secid string) (
	*entity.Share, error) {
	querySting := `SELECT * FROM moex_shares WHERE secid = $1;`

	row := pg.db.QueryRowContext(ctx, querySting, secid)
	if row.Err() != nil {
		return nil, fmt.Errorf("[MoexSharesPostgres GetBySecid]: %w", row.Err())
	}

	var share entity.Share
	err := row.Scan(
		&share.Id,
		&share.Board,
		&share.Engine,
		&share.Market,
		&share.Name,
		&share.ShortName,
		&share.Secid,
	)
	if err != nil {
		return nil, fmt.Errorf("[MoexSharesPostgres GetBySecid]: %w", err)
	}

	return &share, nil
}

func (pg *MoexSharesPostgres) GetListByIds(ctx context.Context, ids []int) (
	[]*entity.Share, error) {
	queryString := "SELECT * FROM moex_shares WHERE id = ANY($1)"
	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))

	if err != nil {
		return nil, fmt.Errorf("[MoexSharesPostgres GetListByIds]: %w", err)
	}

	var shares []*entity.Share
	for rows.Next() {
		var share entity.Share
		err := rows.Scan(&share.Id, &share.Board, &share.Engine, &share.Market, &share.Name,
			&share.ShortName, &share.Secid)
		if err != nil {
			return nil, fmt.Errorf("[MoexSharesPostgres GetListByIds]: %w", err)
		}
		shares = append(shares, &share)
	}

	return shares, nil
}

func (pg *MoexSharesPostgres) Insert(ctx context.Context, share *entity.Share) error {
	querySting := `INSERT INTO moex_shares (board, engine, market, name, shortname, secid)
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`

	_, err := pg.db.ExecContext(ctx, querySting, share.Board, share.Engine, share.Market,
		share.Name, share.ShortName, share.Secid)
	if err != nil {
		return fmt.Errorf("[MoexSharesPostgres Insert]: %w", err)
	}

	return nil
}

type MoexSharesPostgres struct {
	db *sql.DB
}

func NewMoexSharesPostgres(db *sql.DB) *MoexSharesPostgres {
	return &MoexSharesPostgres{db: db}
}
