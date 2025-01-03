package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
)

func (pg *Repository) AddPositionInfo(ctx context.Context, i domain.PositionUpdateInfo) error {
	const op = "Repository.AddPositionInfo"
	var (
		queryString = "UPDATE positions SET "
		count       = 1
		s           = make([]string, 0)
		args        = make([]interface{}, 0)
	)

	if i.Comment != nil {
		s = append(s, fmt.Sprintf("comment = $%d", count))
		args = append(args, *i.Comment)
		count++
	}
	if i.TargetPrice != nil {
		s = append(s, fmt.Sprintf("target_price = $%d", count))
		args = append(args, *i.TargetPrice)
		count++
	}
	queryString += strings.Join(s, ", ")
	queryString += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d;", count, count+1)
	args = append(args, i.ID, i.UserID)

	fmt.Println(queryString, args)
	result, err := pg.t(ctx).ExecContext(ctx, queryString, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (pg *Repository) DeletePosition(ctx context.Context, id int, userID int) error {
	const op = "Repository.DeletePosition"

	queryString := "DELETE FROM positions where id = $1 AND user_id = $2;"

	result, err := pg.db.ExecContext(ctx, queryString, id, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (pg *Repository) GetPosition(ctx context.Context,
	exchange domain.Exchange, portfolioID int,
	securityType domain.SecurityType, ticker string) (domain.Position, error) {
	const op = "Repository.GetPosition"

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

	var p domain.Position

	err := pg.t(ctx).QueryRowContext(
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
		&p.Ticker,
		&p.TargetPrice,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Position{}, storage.ErrNotFound
	}
	if err != nil {
		return domain.Position{}, fmt.Errorf("%s: %w", op, err)
	}

	return p, nil
}

func (pg *Repository) GetPortfolioPositionList(ctx context.Context, id int, userID int) (
	[]domain.Position, error) {
	const op = "Repository.GetPortfolioPositionList"

	queryString := `
		SELECT
	    p.id AS position_id,
	    p.amount,
	    p.average_price,
	    p.board,
	    p.comment,
	    p.exchange,
	    p.portfolio_id,
	    p.security_type,
	    p.shortname,
	    p.ticker,
	    p.target_price,
			COALESCE(
		    json_agg(
		        JSON_BUILD_OBJECT(
	            'id', o.id,
	            'date', o.date,
	            'sourceLink', o.source_link,
	            'targetPrice', o.target_price,
	            'text', o.text,
							'ticker', o.ticker,
							'type', o.type,
	            'expert', JSON_BUILD_OBJECT(
	                'name', e.name,
	                'avatarUrl', e.avatar_url
	            )
	        )
		    ) FILTER (WHERE o.id IS NOT NULL),
	        '[]'::json
			)AS opinions,
		COALESCE(MAX(bond_data.currency), MAX(share_data.currency))
		FROM
		    positions p
		LEFT JOIN
		    opinions_on_positions oop ON p.id = oop.position_id
		LEFT JOIN
		    opinions o ON oop.opinion_id = o.id
		LEFT JOIN
		    experts e ON o.expert_id = e.id
		LEFT JOIN moex_bonds bond_data 
       ON p.security_type = 'BOND' AND p.ticker = bond_data.ticker
		LEFT JOIN moex_shares share_data 
		   ON p.security_type = 'SHARE' AND p.ticker = share_data.ticker
		WHERE
		    p.portfolio_id = $1 AND p.user_id = $2
		GROUP BY
		    p.id, p.amount, p.average_price, p.board, p.comment, p.exchange, p.portfolio_id,
		    p.security_type, p.shortname, p.ticker, p.target_price
		ORDER BY
		    p.id;
	`

	rows, err := pg.t(ctx).QueryContext(ctx, queryString, id, userID)
	if err != nil {
		return nil, fmt.Errorf("%s QueryContext: %w", op, err)
	}
	defer rows.Close()

	var positions []domain.Position

	for rows.Next() {
		var p domain.Position
		var opinionsData []byte

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
			&p.Ticker,
			&p.TargetPrice,
			&opinionsData,
			&p.Currency,
		)
		if err != nil {
			return nil, fmt.Errorf("%s rows.Scan: %w", op, err)
		}
		err = json.Unmarshal(opinionsData, &p.Opinions)
		if err != nil {
			return nil, fmt.Errorf("%s Unmarshal: %w", op, err)
		}
		positions = append(positions, p)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return positions, nil
}

func (pg *Repository) GetUserPositionList(ctx context.Context, userID int) (
	[]domain.Position, error) {
	const op = "Repository.GetUserPositionList"
	var opinionIDs pq.Int64Array

	queryString := `
		SELECT p.id, amount, average_price, board, comment, exchange, portfolio_id,
		security_type, shortname, ticker, target_price, pr.name,
		(SELECT ARRAY_AGG(oop.opinion_id)
			FROM opinions_on_positions oop
			WHERE oop.position_id = p.id
		)AS opinion_ids
		FROM positions p
		LEFT JOIN portfolios pr
		ON pr.id = p.portfolio_id
		WHERE p.user_id = $1;
	`

	rows, err := pg.t(ctx).QueryContext(ctx, queryString, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var positions = make([]domain.Position, 0)

	for rows.Next() {
		var p domain.Position
		err = rows.Scan(&p.ID, &p.Amount, &p.AveragePrice, &p.Board, &p.Comment, &p.Exchange,
			&p.PortfolioID, &p.SecurityType, &p.ShortName, &p.Ticker, &p.TargetPrice,
			&p.PortfolioName, &opinionIDs)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		p.OpinionIDs = convertPqInt64ArrayToIntSlice(opinionIDs)
		positions = append(positions, p)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return positions, nil
}

func (pg *Repository) InsertPosition(ctx context.Context, p domain.Position) error {
	const op = "Repository.InsertPosition"

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

	_, err := pg.t(ctx).ExecContext(ctx, queryString,
		p.Amount,
		p.AveragePrice,
		p.Board,
		p.Comment,
		p.Exchange,
		p.PortfolioID,
		p.SecurityType,
		p.ShortName,
		p.Ticker,
		p.TargetPrice,
		p.UserID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (pg *Repository) UpdatePosition(ctx context.Context, p domain.Position) error {
	const op = "Repository.UpdatePosition"
	var err error

	queryString := `UPDATE positions SET amount = $1, average_price = $2, board = $3,
		comment = $4, exchange = $5, portfolio_id = $6, security_type = $7, shortname = $8,
		ticker = $9, target_price = $10 WHERE id = $11;`

	_, err = pg.t(ctx).ExecContext(ctx, queryString,
		p.Amount,
		p.AveragePrice,
		p.Board,
		p.Comment,
		p.Exchange,
		p.PortfolioID,
		p.SecurityType,
		p.ShortName,
		p.Ticker,
		p.TargetPrice,
		p.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
