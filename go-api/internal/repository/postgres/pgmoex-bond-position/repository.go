package pgmoexbondposition

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/model"
)

type MoexBondPositionPostgres struct {
	db *sql.DB
}

func NewMoexBondPositionPostgres(db *sql.DB) *MoexBondPositionPostgres {
	return &MoexBondPositionPostgres{db: db}
}

func (pg *MoexBondPositionPostgres) Get(ctx context.Context, portfolioId int, securityId int) (*model.Position, error) {
	queryString := `SELECT p.*, b.isin, b.shortname
    FROM moex_bond_positions p
    LEFT JOIN moex_bonds b ON p.security_id = b.id
		WHERE p.portfolio_id = $1 AND p.security_id = $2;`

	var p MoexBondPosition
	var isin string

	row := pg.db.QueryRowContext(ctx, queryString, portfolioId, securityId)
	err := row.Scan(
		&p.Id,
		&p.Amount,
		&p.AveragePrice,
		&p.Comment,
		&p.PortfolioId,
		&p.SecurityId,
		&p.TargetPrice,
		isin,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("[postgres/moex-bond-position Get]: \n%w", err)
	}

	position := FromDBtoPosition(&p)
	position.Secid = isin
	return position, nil
}

func (pg *MoexBondPositionPostgres) GetListByPortfolioId(ctx context.Context, id int) ([]*model.Position, error) {
	queryString := `
		SELECT p.*, b.isin, b.shortname
		FROM moex_bond_positions p 
		LEFT JOIN moex_bonds b ON p.security_id = b.id
		WHERE p.portfolio_id = $1;
	`

	rows, err := pg.db.QueryContext(ctx, queryString, id)
	if err != nil {
		return nil, fmt.Errorf("<-[postgres/moex-bond-position GetListByPortfolioId]: \n%w", err)
	}
	defer rows.Close()

	var positions []*model.Position

	for rows.Next() {
		var position MoexBondPosition

		err := rows.Scan(
			&position.Id,
			&position.Amount,
			&position.AveragePrice,
			&position.Comment,
			&position.PortfolioId,
			&position.SecurityId,
			&position.TargetPrice,
			&position.Isin,
			&position.ShortName,
		)
		if err != nil {
			return nil, fmt.Errorf("<-[postgres/moex-bond-position GetListByPortfolioId]: \n%w", err)
		}
		positions = append(positions, FromDBtoPosition(&position))
	}

	return positions, nil
}

func (pg *MoexBondPositionPostgres) Insert(ctx context.Context, position *model.Position) error {
	p := FromPositionToDB(position)
	queryString := `INSERT INTO moex_bond_positions (amount, average_price, comment, portfolio_id,
		 security_id, target_price) VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := pg.db.ExecContext(ctx, queryString,
		p.Amount,
		p.AveragePrice,
		p.Comment,
		p.PortfolioId,
		p.SecurityId,
		p.TargetPrice,
	)
	if err != nil {
		return fmt.Errorf("[postgres/moex-bond-position Insert]: \n%w", err)
	}
	return nil
}

func (pg *MoexBondPositionPostgres) Update(ctx context.Context, position *model.Position) error {
	p := FromPositionToDB(position)
	queryString := `UPDATE moex_bond_positions SET amount = $1, average_price = $2,
		comment = $3, portfolio_id = $4, security_id = $5, target_price = $6
    WHERE id = $7;`

	_, err := pg.db.ExecContext(ctx, queryString, p.Amount, p.AveragePrice, p.Comment,
		p.PortfolioId, p.SecurityId, p.TargetPrice, p.Id)
	if err != nil {
		return fmt.Errorf("[postgres/moex-bond-position Update]: %w", err)
	}
	return nil
}
