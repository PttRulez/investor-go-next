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
	const op = "MoexSharesPostgres.GetBySecid"

	querySting := `SELECT * FROM moex_shares WHERE secid = $1;`

	row := pg.db.QueryRowContext(ctx, querySting, secid)

	var share entity.Share
	err := row.Scan(
		&share.ID,
		&share.Board,
		&share.Engine,
		&share.Market,
		&share.Name,
		&share.ShortName,
		&share.Secid,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &share, nil
}

func (pg *MoexSharesPostgres) GetListByIDs(ctx context.Context, ids []int) (
	[]*entity.Share, error) {
	const op = "MoexSharesPostgres.GetListByIDs"

	queryString := "SELECT * FROM moex_shares WHERE id = ANY($1)"

	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var shares []*entity.Share
	for rows.Next() {
		var share entity.Share
		err = rows.Scan(&share.ID, &share.Board, &share.Engine, &share.Market, &share.Name,
			&share.ShortName, &share.Secid)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		shares = append(shares, &share)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return shares, nil
}

func (pg *MoexSharesPostgres) Insert(ctx context.Context, share *entity.Share) error {
	const op = "MoexSharesPostgres.Insert"

	querySting := `INSERT INTO moex_shares (board, engine, market, name, shortname, secid)
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`

	_, err := pg.db.ExecContext(ctx, querySting, share.Board, share.Engine, share.Market,
		share.Name, share.ShortName, share.Secid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

type MoexSharesPostgres struct {
	db *sql.DB
}

func NewMoexSharesPostgres(db *sql.DB) *MoexSharesPostgres {
	return &MoexSharesPostgres{db: db}
}