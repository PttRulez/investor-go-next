package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"

	"github.com/lib/pq"
)

func (pg *MoexBondsPostgres) GetBySecid(ctx context.Context, secid string) (
	*entity.Bond, error) {

	queryString := `SELECT * FROM moex_bonds WHERE secid = $1;`

	row := pg.db.QueryRowContext(ctx, queryString, secid)

	bond := entity.Bond{}
	err := row.Scan(
		&bond.Id,
		&bond.Board,
		&bond.Engine,
		&bond.Market,
		&bond.Name,
		&bond.ShortName,
		&bond.Secid,
	)
	if err != nil {
		return nil, fmt.Errorf("\n<-[MoexBondsPostgres GetBySecid]:\n %w", err)
	}

	return &bond, nil
}

func (pg *MoexBondsPostgres) GetListByIds(ctx context.Context, ids []int) (
	[]*entity.Bond, error) {

	queryString := "SELECT * FROM moex_bonds WHERE id = ANY($1)"
	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))

	if err != nil {
		return nil, fmt.Errorf("[MoexBondsPostgres GetListByIds]: %w", err)
	}

	var bonds []*entity.Bond
	for rows.Next() {
		var bond entity.Bond
		err := rows.Scan(&bond.Id, &bond.Board, &bond.Engine, &bond.Market, &bond.Name,
			&bond.ShortName, &bond.Secid)
		if err != nil {
			return nil, fmt.Errorf("[MoexBondsPostgres GetListByIds]: %w", err)
		}
		bonds = append(bonds, &bond)
	}

	return bonds, nil
}

func (pg *MoexBondsPostgres) Insert(ctx context.Context, bond *entity.Bond) error {
	querySting := `INSERT INTO moex_bonds (board, engine, market, name, shortname, secid)
    VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := pg.db.ExecContext(ctx, querySting, bond.Board, bond.Engine, bond.Market,
		bond.Name, bond.ShortName, bond.Secid)
	if err != nil {
		return fmt.Errorf("[MoexBondsPostgres Insert]: %w", err)
	}

	return nil
}

type MoexBondsPostgres struct {
	db *sql.DB
}

func NewMoexBondsPostgres(db *sql.DB) *MoexBondsPostgres {
	return &MoexBondsPostgres{db: db}
}
