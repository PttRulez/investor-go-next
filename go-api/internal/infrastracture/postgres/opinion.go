package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
)

type OpinionPostgres struct {
	db *sql.DB
}

func NewOpinionPostgres(db *sql.DB) *OpinionPostgres {
	return &OpinionPostgres{db: db}
}

func (pg *OpinionPostgres) Delete(ctx context.Context, id int) error {
	const op = "OpinionPostgres.Delete"

	queryString := "DELETE FROM opinions where id = $1;"
	_, err := pg.db.ExecContext(ctx, queryString, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (pg *OpinionPostgres) Insert(ctx context.Context, o *entity.Opinion) error {
	const op = "OpinionPostgres.Insert"

	queryString := `INSERT INTO opinions (date, exchange, expert_id, text, security_id,  
    security_type, source_link, target_price, type, user_id) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	_, err := pg.db.ExecContext(ctx, queryString, o.Date, o.Exchange, o.ExpertID, o.Text, o.SecurityID,
		o.SecurityType, o.SourceLink, o.TargetPrice, o.Type, o.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (pg *OpinionPostgres) Update(ctx context.Context, o *entity.Opinion) error {
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

func (pg *OpinionPostgres) GetListByUserID(ctx context.Context, userID int) ([]*entity.Opinion, error) {
	const op = "OpinionPostgres.GetListByUserId"

	queryString := `SELECT * FROM opinions WHERE user_id = $1;`
	rows, err := pg.db.QueryContext(ctx, queryString, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var opinions []*entity.Opinion

	for rows.Next() {
		o := new(entity.Opinion)
		err = rows.Scan(&o.ID, &o.Date, &o.Exchange, &o.ExpertID, &o.SecurityID,
			&o.SecurityType, &o.SourceLink, &o.TargetPrice, &o.Type, &o.UserID)
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
