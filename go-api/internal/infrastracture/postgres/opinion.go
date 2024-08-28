package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
)

type OpinionPostgres struct {
	db *sql.DB
}

func NewOpinionPostgres(db *sql.DB) *OpinionPostgres {
	return &OpinionPostgres{db: db}
}

func (pg *OpinionPostgres) Delete(ctx context.Context, id int, userID int) error {
	const op = "OpinionPostgres.Delete"

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
		return database.ErrNotFound
	}

	return nil
}

func (pg *OpinionPostgres) Insert(ctx context.Context, o entity.Opinion) (entity.Opinion, error) {
	const op = "OpinionPostgres.Insert"

	queryString := `INSERT INTO opinions (date, exchange, expert_id, text, security_id,  
    security_type, source_link, target_price, type, user_id) VALUES ($1, $2, $3, $4, $5,
		$6, $7, $8, $9, $10) RETURNING id, date, exchange, expert_id, text, security_id,
    security_type, source_link, target_price, type;`

	var r entity.Opinion
	err := pg.db.QueryRowContext(ctx, queryString, o.Date, o.Exchange, o.ExpertID, o.Text, o.SecurityID,
		o.SecurityType, o.SourceLink, o.TargetPrice, o.Type, o.UserID).
		Scan(&r.ID, &r.Date, &r.Exchange, &r.ExpertID, &r.Text, &r.SecurityID, &r.SecurityType,
			&r.SourceLink, &r.TargetPrice, &r.Type)
	if err != nil {
		return entity.Opinion{}, fmt.Errorf("%s: %w", op, err)
	}

	return r, nil
}

func (pg *OpinionPostgres) Update(ctx context.Context, o entity.Opinion) error {
	const op = "OpinionPostgres.Update"

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

func (pg *OpinionPostgres) GetOpinionsList(ctx context.Context, f entity.OpinionFilters,
	userID int) ([]entity.Opinion, error) {
	const op = "OpinionPostgres.GetListByUserId"

	queryString := `SELECT id, date, exchange, expert_id, text, security_id, security_type,
	 source_link, target_price, type FROM opinions WHERE user_id = $1`
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

	var opinions = make([]entity.Opinion, 0)

	for rows.Next() {
		var o entity.Opinion
		err = rows.Scan(&o.ID, &o.Date.Time, &o.Exchange, &o.ExpertID, &o.Text, &o.SecurityID,
			&o.SecurityType, &o.SourceLink, &o.TargetPrice, &o.Type)
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
