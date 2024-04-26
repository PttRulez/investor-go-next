package moexbond

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pttrulez/investor-go/internal/model"
)

type MoexBondsPostgres struct {
	db *sql.DB
}

func NewMoexBondsPostgres(db *sql.DB) *MoexBondsPostgres {
	return &MoexBondsPostgres{db: db}
}

func (pg *MoexBondsPostgres) GetByISIN(ctx context.Context, isin string) (*model.MoexBond, error) {
	queryString := `SELECT * FROM moex_bonds WHERE isin = $1;`

	row := pg.db.QueryRowContext(ctx, queryString, isin)

	bond := MoexBond{}
	err := row.Scan(
		&bond.Id,
		&bond.Board,
		&bond.Engine,
		&bond.Market,
		&bond.Name,
		&bond.ShortName,
		&bond.Isin,
	)
	if err != nil {
		return nil, fmt.Errorf("\n<-[MoexBondsPostgres GetByISIN]:\n %w", err)
	}

	return FromDBToMoexBond(&bond), nil
}

func (pg *MoexBondsPostgres) GetListByIds(ctx context.Context, ids []int) ([]*model.MoexBond, error) {
	queryString := "SELECT * FROM moex_bonds WHERE id = ANY($1)"
	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))

	if err != nil {
		return nil, fmt.Errorf("[MoexBondsPostgres GetListByIds]: %w", err)
	}

	var bonds []*model.MoexBond
	for rows.Next() {
		var bond MoexBond
		err := rows.Scan(&bond.Id, &bond.Board, &bond.Engine, &bond.Market, &bond.Name, &bond.ShortName, &bond.Isin)
		if err != nil {
			return nil, fmt.Errorf("[MoexBondsPostgres GetListByIds]: %w", err)
		}
		bonds = append(bonds, FromDBToMoexBond(&bond))
	}

	return bonds, nil
}

func (pg *MoexBondsPostgres) Insert(ctx context.Context, bond *model.MoexBond) error {
	querySting := `INSERT INTO moex_bonds (board, engine, market, name, shortname, isin)
    VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := pg.db.ExecContext(ctx, querySting, bond.Board, bond.Engine, bond.Market, bond.Name, bond.ShortName, bond.Isin)
	if err != nil {
		return fmt.Errorf("[MoexBondsPostgres Insert]: %w", err)
	}

	return nil
}
