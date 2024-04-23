package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/types"
)

type MoexSharePositionPostgres struct {
	db *sql.DB
}

func NewMoexSharePositionPostgres(db *sql.DB) *MoexSharePositionPostgres {
	return &MoexSharePositionPostgres{db: db}
}

func (pg *MoexSharePositionPostgres) Get(ctx context.Context, portfolioId int, securityId int) (*types.SharePosition, error) {
	queryString := `SELECT * FROM moex_share_positions WHERE portfolio_id = $1 AND security_id = $2;`

	var position types.SharePosition
	position.Exchange = types.EXCH_Moex
	position.PortfolioId = portfolioId
	position.SecurityId = securityId

	row := pg.db.QueryRowContext(ctx, queryString, portfolioId, securityId)

	err := row.Scan(
		&position.Id,
		&position.Amount,
		&position.AveragePrice,
		&position.Comment,
		&position.PortfolioId,
		&position.SecurityId,
		&position.TargetPrice,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("[MoexSharePositionPostgres Get]: %w", err)
	}

	return &position, nil
}

func (pg *MoexSharePositionPostgres) GetListByPortfolioId(ctx context.Context, id int) ([]*types.SharePosition, error) {
	queryString := `
		SELECT p.id, p.amount, p.average_price, p.comment, p.portfolio_id,
		p.security_id, p.target_price, m.ticker
		FROM moex_share_positions p 
		LEFT JOIN moex_shares m ON p.security_id = m.id
		WHERE p.portfolio_id = $1;
	`

	rows, err := pg.db.QueryContext(ctx, queryString, id)
	if err != nil {
		return nil, fmt.Errorf("<-[MoexSharePositionPostgres GetListByPortfolioId]: \n%w", err)
	}
	defer rows.Close()

	var positions []*types.SharePosition
	for rows.Next() {
		var position types.SharePosition
		position.Exchange = types.EXCH_Moex

		err := rows.Scan(&position.Id, &position.Amount, &position.AveragePrice, &position.Comment,
			&position.PortfolioId, &position.SecurityId, &position.TargetPrice, &position.Ticker)
		if err != nil {
			return nil, fmt.Errorf("<-[MoexSharePositionPostgres GetListByPortfolioId]: \n%w", err)
		}
		positions = append(positions, &position)
	}

	return positions, nil
}

func (pg *MoexSharePositionPostgres) Insert(ctx context.Context, p *types.SharePosition) error {
	queryString := `INSERT INTO moex_share_positions (amount, average_price, comment,
    portfolio_id, security_id, target_price) 
    VALUES ($1, $2, $3, $4, $5, $6) ;`

	_, err := pg.db.ExecContext(ctx, queryString, p.Amount, p.AveragePrice, p.Comment,
		p.PortfolioId, p.SecurityId, p.TargetPrice)
	if err != nil {
		return fmt.Errorf("[MoexSharePositionPostgres Insert]: %w", err)
	}
	return nil
}

func (pg *MoexSharePositionPostgres) Update(ctx context.Context, p *types.SharePosition) error {
	queryString := `UPDATE moex_share_positions SET amount = $1, average_price = $2, comment = $3,
    portfolio_id = $4, security_id = $5, target_price = $6 WHERE id = $7;`

	_, err := pg.db.ExecContext(ctx, queryString, p.Amount, p.AveragePrice, p.Comment,
		p.PortfolioId, p.SecurityId, p.TargetPrice, p.Id)
	if err != nil {
		return fmt.Errorf("[MoexSharePositionPostgres Update]: %w", err)
	}
	return nil
}
