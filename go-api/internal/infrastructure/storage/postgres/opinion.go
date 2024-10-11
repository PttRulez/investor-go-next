package postgres

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
)

func (pg *Repository) AttachOpinionToPosition(ctx context.Context, opinionID int, positionID int) error {
	const op = "Repository.AttachToPosition"

	// queryString := `INSERT INTO opinions_on_positions (position_id, opinion_id) VALUES ($1, $2)
	// 	ON CONFLICT (position_id, opinion_id) DO NOTHING;`
	queryString := `WITH deleted AS (
	    DELETE FROM opinions_on_positions
	    WHERE position_id = $1 AND opinion_id = $2
	    RETURNING *
		)
		INSERT INTO opinions_on_positions (position_id, opinion_id)
		SELECT $1, $2
		WHERE NOT EXISTS (SELECT 1 FROM deleted);`

	_, err := pg.db.ExecContext(ctx, queryString, positionID, opinionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *Repository) DeleteOpinion(ctx context.Context, id int, userID int) error {
	const op = "Repository.DeleteOpinion"

	queryString := "DELETE FROM opinions where id = $1 and user_id = $2;"
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

func (pg *Repository) GetOpinionsList(ctx context.Context, f domain.OpinionFilters,
	userID int) ([]domain.Opinion, error) {
	const op = "Repository.GetOpinionsList"

	queryString := `SELECT o.id, date, exchange, expert_id, text, security_id, security_type,
		source_link, target_price, type, name AS expertname
	FROM opinions o
	LEFT JOIN experts e ON e.id = o.expert_id
	WHERE o.user_id = $1`
	args := []interface{}{userID}
	count := 2

	if f.ExpertID != nil {
		queryString += fmt.Sprintf(" AND expert_id = $%d", count)
		args = append(args, *f.ExpertID)
		count++
	}
	if f.Exchange != nil {
		queryString += fmt.Sprintf(" AND exchange = $%d", count)
		args = append(args, *f.Exchange)
		count++
	}
	if f.SecurityID != nil {
		queryString += fmt.Sprintf(" AND security_id = $%d", count)
		args = append(args, *f.SecurityID)
		count++
	}
	if f.SecurityType != nil {
		queryString += fmt.Sprintf(" AND security_type = $%d", count)
		args = append(args, *f.SecurityType)
	}
	queryString += ";"

	rows, err := pg.db.QueryContext(ctx, queryString, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var opinions = make([]domain.Opinion, 0)

	for rows.Next() {
		var o domain.Opinion
		err = rows.Scan(&o.ID, &o.Date.Time, &o.Exchange, &o.ExpertID, &o.Text, &o.SecurityID,
			&o.SecurityType, &o.SourceLink, &o.TargetPrice, &o.Type, &o.Expert.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		opinions = append(opinions, o)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return opinions, nil
}

func (pg *Repository) InsertOpinion(ctx context.Context, o domain.Opinion) (domain.Opinion, error) {
	const op = "Repository.InsertOpinion"

	queryString := `INSERT INTO opinions (date, exchange, expert_id, text, security_id,  
    security_type, source_link, target_price, ticker, type, user_id) VALUES ($1, $2, $3, $4, $5,
		$6, $7, $8, $9, $10, $11) RETURNING id, date, exchange, expert_id, text, security_id,
    security_type, source_link, target_price, ticker, type;`

	var r domain.Opinion
	err := pg.db.QueryRowContext(ctx, queryString, o.Date.Time, o.Exchange, o.ExpertID, o.Text, o.SecurityID,
		o.SecurityType, o.SourceLink, o.TargetPrice, o.Ticker, o.Type, o.UserID).
		Scan(&r.ID, &r.Date.Time, &r.Exchange, &r.ExpertID, &r.Text, &r.SecurityID, &r.SecurityType,
			&r.SourceLink, &r.TargetPrice, &r.Ticker, &r.Type)
	if err != nil {
		return domain.Opinion{}, fmt.Errorf("%s: %w", op, err)
	}

	return r, nil
}

func (pg *Repository) Update(ctx context.Context, o domain.Opinion) error {
	const op = "Repository.UpdateOpinion"

	queryString := `UPDATE opinions SET date = $1, exchange = $2, expert_id = $3,
    security_id = $4, security_type = $5, source_link = $6, target_price = $7,
    type = $8, user_id = $9 WHERE id = $10;`

	_, err := pg.db.ExecContext(ctx, queryString, o.Date, o.Exchange, o.ExpertID, o.SecurityID,
		o.SecurityType, o.SourceLink, o.TargetPrice, o.Type, o.UserID, o.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
