package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
)

type PositionPostgres struct {
	db *sql.DB
}

func NewPositionPostgres(db *sql.DB) *PositionPostgres {
	return &PositionPostgres{db: db}
}

func (pg *PositionPostgres) GetPositionForSecurity(ctx context.Context,
	exchange entity.Exchange, portfolioID int,
	securityType entity.SecurityType, ticker string) (entity.Position, error) {
	const op = "PositionPostgres.GetPositionForSecurity"

	queryString := `SELECT
		id,
		amount,
		average_price,
	  board,
		comment,
		exchange,
		portfolio_id,
		security_type,
		shortname,
		ticker,
		target_price
    FROM positions
    WHERE exchange = $1 AND portfolio_id = $2 AND security_type = $3 AND ticker = $4;`

	var p entity.Position

	err := pg.db.QueryRowContext(
		ctx,
		queryString,
		exchange,
		portfolioID,
		securityType,
		ticker,
	).Scan(
		&p.ID,
		&p.Amount,
		&p.AveragePrice,
		&p.Board,
		&p.Comment,
		&p.Exchange,
		&p.PortfolioID,
		&p.SecurityType,
		&p.ShortName,
		&p.Secid,
		&p.TargetPrice,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Position{}, database.ErrNotFound
	}
	if err != nil {
		return entity.Position{}, fmt.Errorf("%s: %w", op, err)
	}

	return p, nil
}

func (pg *PositionPostgres) GetListByPortfolioID(ctx context.Context, id int, userID int) (
	[]entity.Position, error) {
	const op = "PositionPostgres.GetListByPortfolioId"

	queryString := `
		SELECT * FROM positions 
		WHERE portfolio_id = $1 AND user_id = $2;
	`

	rows, err := pg.db.QueryContext(ctx, queryString, id, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var positions []entity.Position

	for rows.Next() {
		var p entity.Position

		err = rows.Scan(
			&p.ID,
			&p.Amount,
			&p.AveragePrice,
			&p.Board,
			&p.Comment,
			&p.Exchange,
			&p.PortfolioID,
			&p.SecurityType,
			&p.ShortName,
			&p.Secid,
			&p.TargetPrice,
			&p.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		positions = append(positions, p)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return positions, nil
}

func (pg *PositionPostgres) Insert(ctx context.Context, p entity.Position) error {
	const op = "PositionPostgres.Insert"

	queryString := `INSERT INTO positions (
		amount,
		average_price,
	  board,
		comment,
		exchange,
		portfolio_id,
		security_type,
		shortname,
		ticker,
		target_price,
		user_id
    ) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ;`

	_, err := pg.db.ExecContext(ctx, queryString,
		p.Amount,
		p.AveragePrice,
		p.Board,
		p.Comment,
		p.Exchange,
		p.PortfolioID,
		p.SecurityType,
		p.ShortName,
		p.Secid,
		p.TargetPrice,
		p.UserID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (pg *PositionPostgres) Update(ctx context.Context, p entity.Position) error {
	const op = "PositionPostgres.Update"

	queryString := `UPDATE positions SET amount = $1, average_price = $2, comment = $3, exchange = $4,
		portfolio_id = $5, security_type = $6, shortname = $7, ticker = $8, target_price = $9 WHERE id = $10;`

	_, err := pg.db.ExecContext(ctx, queryString,
		p.Amount,
		p.AveragePrice,
		p.Board,
		p.Comment,
		p.Exchange,
		p.PortfolioID,
		p.SecurityType,
		p.ShortName,
		p.Secid,
		p.TargetPrice,
		p.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
