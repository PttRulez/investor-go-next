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
	const op = "MoexBondsPostgres.GetBySecid"

	queryString := `SELECT * FROM moex_bonds WHERE secid = $1;`

	row := pg.db.QueryRowContext(ctx, queryString, secid)

	bond := entity.Bond{}
	err := row.Scan(
		&bond.ID,
		&bond.Board,
		&bond.Engine,
		&bond.Market,
		&bond.Name,
		&bond.ShortName,
		&bond.Secid,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &bond, nil
}

func (pg *MoexBondsPostgres) GetListByIDs(ctx context.Context, ids []int) (
	[]*entity.Bond, error) {
	const op = "MoexBondsPostgres.GetListByIDs"

	queryString := "SELECT * FROM moex_bonds WHERE id = ANY($1)"

	rows, err := pg.db.QueryContext(ctx, queryString, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var bonds []*entity.Bond

	for rows.Next() {
		var bond entity.Bond
		err = rows.Scan(&bond.ID, &bond.Board, &bond.Engine, &bond.Market, &bond.Name,
			&bond.ShortName, &bond.Secid)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		bonds = append(bonds, &bond)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return bonds, nil
}

func (pg *MoexBondsPostgres) Insert(ctx context.Context, bond *entity.Bond) error {
	const op = "MoexBondsPostgres.Insert"

	querySting := `INSERT INTO moex_bonds (board, engine, market, name, shortname, secid)
    VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := pg.db.ExecContext(ctx, querySting, bond.Board, bond.Engine, bond.Market,
		bond.Name, bond.ShortName, bond.Secid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

type MoexBondsPostgres struct {
	db *sql.DB
}

func NewMoexBondsPostgres(db *sql.DB) *MoexBondsPostgres {
	return &MoexBondsPostgres{db: db}
}
