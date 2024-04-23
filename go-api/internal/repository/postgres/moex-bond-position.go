package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/types"
)

type MoexBondPositionPostgres struct {
	db *sql.DB
}

func NewMoexBondPositionPostgres(db *sql.DB) *MoexBondPositionPostgres {
	return &MoexBondPositionPostgres{db: db}
}

func (pg *MoexBondPositionPostgres) Get(ctx context.Context, portfolioId int, securityId int) (*types.BondPosition, error) {
	queryString := `SELECT * FROM positions WHERE portfolio_id = $1 AND security_id = $2;`

	var position types.BondPosition
	position.Exchange = types.EXCH_Moex
	position.PortfolioId = portfolioId
	position.SecurityId = securityId

	row := pg.db.QueryRowContext(ctx, queryString, portfolioId, securityId)

	err := row.Scan(&position.Id, &position.Amount, &position.AveragePrice, &position.Comment,
		&position.TargetPrice)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("[MoexBondPositionPostgres Get]: %w", err)
	}

	return &position, nil
}

func (pg *MoexBondPositionPostgres) GetListByPortfolioId(ctx context.Context, id int) ([]*types.BondPosition, error) {
	queryString := `
		SELECT p.id, p.amount, p.average_price, p.comment, p.portfolio_id,
		p.security_id, p.target_price, m.ticker
		FROM moex_bond_positions p 
		LEFT JOIN moex_bonds m ON p.security_id = m.id
		WHERE p.portfolio_id = $1;
	`

	rows, err := pg.db.QueryContext(ctx, queryString, id)
	if err != nil {
		return nil, fmt.Errorf("<-[MoexBondPositionPostgres GetListByPortfolioId]: \n%w", err)
	}
	defer rows.Close()

	var positions []*types.BondPosition

	for rows.Next() {
		var position types.BondPosition
		position.Exchange = types.EXCH_Moex

		err := rows.Scan(&position.Id, &position.Amount, &position.AveragePrice, &position.Comment,
			&position.PortfolioId, &position.SecurityId, &position.TargetPrice)
		if err != nil {
			return nil, fmt.Errorf("<-[MoexBondPositionPostgres GetListByPortfolioId]: \n%w", err)
		}
		positions = append(positions, &position)
	}

	return positions, nil
}

func (pg *MoexBondPositionPostgres) Insert(ctx context.Context, p *types.BondPosition) error {
	queryString := `INSERT INTO positions (amount, average_price, comment, exchange,
    portfolio_id, security_id, security_type,  target_price) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ;`

	_, err := pg.db.ExecContext(ctx, queryString, p.Amount, p.AveragePrice, p.Comment, p.Exchange,
		p.PortfolioId, p.SecurityId, p.TargetPrice)
	if err != nil {
		return fmt.Errorf("[MoexBondPositionPostgres Insert]: %w", err)
	}
	return nil
}

func (pg *MoexBondPositionPostgres) Update(ctx context.Context, p *types.BondPosition) error {
	queryString := `UPDATE positions SET amount = $1, average_price = $2, comment = $3, exchange = $4,
    portfolio_id = $5, security_id = $6, security_type = $7, target_price = $8
    WHERE id = $9;`

	_, err := pg.db.ExecContext(ctx, queryString, p.Amount, p.AveragePrice, p.Comment, p.Exchange,
		p.PortfolioId, p.SecurityId, p.TargetPrice, p.Id)
	if err != nil {
		return fmt.Errorf("[MoexBondPositionPostgres Update]: %w", err)
	}
	return nil
}
