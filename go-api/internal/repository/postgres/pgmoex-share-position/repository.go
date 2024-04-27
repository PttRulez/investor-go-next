package pgmoexshareposition

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/model"
)

type MoexSharePositionPostgres struct {
	db *sql.DB
}

func NewMoexSharePositionPostgres(db *sql.DB) *MoexSharePositionPostgres {
	return &MoexSharePositionPostgres{db: db}
}

func (pg *MoexSharePositionPostgres) Get(ctx context.Context, portfolioId int, securityId int) (*model.Position, error) {
	queryString := `SELECT p.*, s.isin, s.shortname
    FROM moex_shared_positions p
    LEFT JOIN moex_shares s ON p.security_id = s.id
		WHERE p.portfolio_id = $1 AND p.security_id = $2;`

	var p MoexSharePosition

	row := pg.db.QueryRowContext(ctx, queryString, portfolioId, securityId)

	err := row.Scan(
		&p.Id,
		&p.Amount,
		&p.AveragePrice,
		&p.Comment,
		&p.PortfolioId,
		&p.SecurityId,
		&p.TargetPrice,
		&p.Ticker,
		&p.ShortName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("[MoexSharePositionPostgres Get]: %w", err)
	}

	return FromDBtoPosition(&p), nil
}

func (pg *MoexSharePositionPostgres) GetListByPortfolioId(ctx context.Context, id int) ([]*model.Position, error) {
	queryString := `
		SELECT p.*, m.ticker, m.shortname
		FROM moex_share_positions p 
		LEFT JOIN moex_shares m ON p.security_id = m.id
		WHERE p.portfolio_id = $1;
	`

	rows, err := pg.db.QueryContext(ctx, queryString, id)
	if err != nil {
		return nil, fmt.Errorf("<-[MoexSharePositionPostgres GetListByPortfolioId]: \n%w", err)
	}
	defer rows.Close()

	var positions []*model.Position
	for rows.Next() {
		var p MoexSharePosition

		err := rows.Scan(
			&p.Id,
			&p.Amount,
			&p.AveragePrice,
			&p.Comment,
			&p.PortfolioId,
			&p.SecurityId,
			&p.TargetPrice,
			&p.Ticker,
			&p.ShortName,
		)
		if err != nil {
			return nil, fmt.Errorf("<-[MoexSharePositionPostgres GetListByPortfolioId]: \n%w", err)
		}
		positions = append(positions, FromDBtoPosition(&p))
	}

	return positions, nil
}

func (pg *MoexSharePositionPostgres) Insert(ctx context.Context, position *model.Position) error {
	p := FromPositionToDB(position)
	queryString := `INSERT INTO moex_share_positions (amount, average_price, comment,
    portfolio_id, security_id, target_price) 
    VALUES ($1, $2, $3, $4, $5, $6) ;`

	_, err := pg.db.ExecContext(ctx, queryString,
		p.Amount,
		p.AveragePrice,
		p.Comment,
		p.PortfolioId,
		p.SecurityId,
		p.TargetPrice,
	)
	if err != nil {
		return fmt.Errorf("[MoexSharePositionPostgres Insert]: %w", err)
	}
	return nil
}

func (pg *MoexSharePositionPostgres) Update(ctx context.Context, position *model.Position) error {
	p := FromPositionToDB(position)
	queryString := `UPDATE moex_share_positions SET amount = $1, average_price = $2, comment = $3,
    portfolio_id = $4, security_id = $5, target_price = $6 WHERE id = $7;`

	_, err := pg.db.ExecContext(ctx, queryString,
		p.Amount,
		p.AveragePrice,
		p.Comment,
		p.PortfolioId,
		p.SecurityId,
		p.TargetPrice,
		p.Id,
	)
	if err != nil {
		return fmt.Errorf("[MoexSharePositionPostgres Update]: %w", err)
	}
	return nil
}
