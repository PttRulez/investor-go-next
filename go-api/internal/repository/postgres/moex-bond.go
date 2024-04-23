package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pttrulez/investor-go/internal/types"
	tmoex "github.com/pttrulez/investor-go/internal/types/moex"
)

type MoexBondsPostgres struct {
	db *sql.DB
}

func NewMoexBondsPostgres(db *sql.DB) types.MoexBondRepository {
	return &MoexBondsPostgres{db: db}
}

func (pg *MoexBondsPostgres) GetByTicker(ctx context.Context, ticker string) (*tmoex.Bond, error) {
	queryString := `SELECT * FROM moex_bonds WHERE ticker = $1;`

	row := pg.db.QueryRowContext(ctx, queryString, ticker)

	bond := &tmoex.Bond{}
	err := row.Scan(&bond.Board, &bond.Engine, &bond.Market, &bond.Id, &bond.Name, &bond.ShortName, &bond.Ticker)
	if err != nil {

		return nil, fmt.Errorf("[MoexBondsPostgres GetByTicker]: %w", err)
	}

	return bond, nil
}

func (pg *MoexBondsPostgres) GetListByIds(ctx context.Context, ids []int) ([]*tmoex.Bond, error) {
	queryString := "SELECT * FROM moex_bonds WHERE id = ANY($1)"
	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))

	if err != nil {
		return nil, fmt.Errorf("[MoexBondsPostgres GetListByIds]: %w", err)
	}

	var bonds []*tmoex.Bond
	for rows.Next() {
		var bond tmoex.Bond
		err := rows.Scan(&bond.Id, &bond.Board, &bond.Engine, &bond.Market, &bond.Name, &bond.ShortName, &bond.Ticker)
		if err != nil {
			return nil, fmt.Errorf("[MoexBondsPostgres GetListByIds]: %w", err)
		}
		bonds = append(bonds, &bond)
	}

	return bonds, nil
}

func (pg *MoexBondsPostgres) Insert(ctx context.Context, bond *tmoex.Bond) error {
	querySting := `INSERT INTO moex_bonds (board, engine, market, id, name, shortname, ticker)
    VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;`

	_, err := pg.db.ExecContext(ctx, querySting, bond.Board, bond.Engine, bond.Market, bond.Id, bond.Name, bond.ShortName, bond.Ticker)
	if err != nil {
		return fmt.Errorf("[MoexBondsPostgres Insert]: %w", err)
	}

	return nil
}
