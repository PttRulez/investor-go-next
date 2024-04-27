package pgopinion

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/model"
	"github.com/pttrulez/investor-go/internal/repository"
)

type OpinionPostgres struct {
	db *sql.DB
}

func NewOpinionPostgres(db *sql.DB) repository.OpinionRepository {
	return &OpinionPostgres{db: db}
}

func (pg *OpinionPostgres) Delete(ctx context.Context, id int) error {
	queryString := "DELETE FROM opinions where id = $1;"
	_, err := pg.db.ExecContext(ctx, queryString, id)
	if err != nil {
		return fmt.Errorf("[OpinionPostgres Delete]: %w", err)
	}
	return nil
}

func (pg *OpinionPostgres) Insert(ctx context.Context, o *model.Opinion) error {
	queryString := `INSERT INTO opinions (date, exchange, expert_id, security_id, 
    security_type, source_link, target_price, type, user_id) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

	_, err := pg.db.ExecContext(ctx, queryString, o.Date, o.Exchange, o.ExpertId, o.SecurityId,
		o.SecurityType, o.SourceLink, o.TargetPrice, o.Type, o.UserId)
	if err != nil {
		return fmt.Errorf("[OpinionPostgres Insert]: %w", err)
	}
	return nil
}

func (pg *OpinionPostgres) Update(ctx context.Context, model *model.Opinion) error {
	queryString := `UPDATE opinions SET date = $1, exchange = $2, expert_id = $3,
    security_id = $4, security_type = $5, source_link = $6, target_price = $7,
    type = $8, user_id = $9 WHERE id = $10;`

	o := FromOpinionToDBOpinion(model)
	_, err := pg.db.ExecContext(ctx, queryString, o.Date, o.Exchange, o.ExpertId, o.SecurityId,
		o.SecurityType, o.SourceLink, o.TargetPrice, o.Type, o.UserId, o.Id)
	if err != nil {
		return fmt.Errorf("[OpinionPostgres Update]: %w", err)
	}
	return nil
}

func (pg *OpinionPostgres) GetListByUserId(ctx context.Context, userId int) ([]*model.Opinion, error) {
	queryString := `SELECT * FROM opinions WHERE user_id = $1;`
	rows, err := pg.db.QueryContext(ctx, queryString, userId)
	if err != nil {
		return nil, fmt.Errorf("[OpinionPostgres GetListByUserId]: %w", err)
	}
	defer rows.Close()

	opinions := []*model.Opinion{}

	for rows.Next() {
		var o Opinion
		err = rows.Scan(&o.Id, &o.Date, &o.Exchange, &o.ExpertId, &o.SecurityId,
			&o.SecurityType, &o.SourceLink, &o.TargetPrice, &o.Type, &o.UserId)
		if err != nil {
			return nil, fmt.Errorf("[OpinionPostgres GetListByUserId]: %w", err)
		}
		opinions = append(opinions, FromDBToModelOpinion(&o))
	}

	return opinions, nil
}
