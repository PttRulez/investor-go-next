package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
)

type PositionPostgres struct {
	db *sql.DB
}

func NewPositionPostgres(db *sql.DB) *PositionPostgres {
	return &PositionPostgres{db: db}
}

func (pg *PositionPostgres) GetForSecurity(ctx context.Context, exchange entity.Exchange, portfolioId int,
	securityType entity.SecurityType, ticker string) (*entity.Position, error) {
	queryString := `SELECT * FROM positions
    WHERE exchange = $1 AND portfolio_id = $2 AND security_type = $3 AND ticker = $4;`

	var p = new(entity.Position)

	row := pg.db.QueryRowContext(
		ctx,
		queryString,
		exchange,
		portfolioId,
		securityType,
		ticker,
	)

	err := row.Scan(
		p.Id,
		p.Amount,
		p.AveragePrice,
		p.Board,
		p.Comment,
		p.Exchange,
		p.PortfolioId,
		p.SecurityType,
		p.Ticker,
		p.TargetPrice,
		p.ShortName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("[PositionPostgres Get]: %w", err)
	}

	return p, nil
}

func (pg *PositionPostgres) GetListByPortfolioId(ctx context.Context, id int, userId int) (
	[]*entity.Position, error) {
	queryString := `
		SELECT * FROM positions 
		WHERE portfolio_id = $1 AND user_id = $2;
	`

	rows, err := pg.db.QueryContext(ctx, queryString, id, userId)
	if err != nil {
		return nil, fmt.Errorf("[PositionPostgres GetListByPortfolioId]: \n%w", err)
	}
	defer rows.Close()

	var positions []*entity.Position

	for rows.Next() {
		var p *entity.Position

		err := rows.Scan(
			p.Id,
			p.Amount,
			p.AveragePrice,
			p.Board,
			p.Comment,
			p.Exchange,
			p.PortfolioId,
			p.SecurityType,
			p.ShortName,
			p.Ticker,
			p.TargetPrice,
			p.UserId,
		)
		if err != nil {
			return nil, fmt.Errorf("[MoexSharePositionPostgres GetListByPortfolioId]: %w", err)
		}
		positions = append(positions, p)
	}

	return positions, nil
}

func (pg *PositionPostgres) Insert(ctx context.Context, p *entity.Position) error {
	queryString := `INSERT INTO positions (
		amount,
		average_price,
	   	board,
		comment,
		exchange,
		portfolio_id,
		security_type,
		short_name,
		ticker,
		target_price,
		user_id
    ) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ;`

	_, err := pg.db.ExecContext(ctx, queryString,
		p.Amount,
		p.AveragePrice,
		p.Board,
		p.Comment,
		p.Exchange,
		p.PortfolioId,
		p.SecurityType,
		p.ShortName,
		p.Ticker,
		p.TargetPrice,
		p.UserId,
	)
	if err != nil {
		return fmt.Errorf("[PositionPostgres Insert]: %w", err)
	}
	return nil
}

func (pg *PositionPostgres) Update(ctx context.Context, p *entity.Position) error {
	queryString := `UPDATE positions SET amount = $1, average_price = $2, comment = $3, exchange = $4,
		portfolio_id = $5, security_type = $6, short_name = $7, ticker = $8, target_price = $9 WHERE id = $10;`

	_, err := pg.db.ExecContext(ctx, queryString,
		p.Amount,
		p.AveragePrice,
		p.Board,
		p.Comment,
		p.Exchange,
		p.PortfolioId,
		p.SecurityType,
		p.ShortName,
		p.Ticker,
		p.TargetPrice,
		p.Id,
	)
	if err != nil {
		return fmt.Errorf("[PositionPostgres Update]: %w", err)
	}
	return nil
}
