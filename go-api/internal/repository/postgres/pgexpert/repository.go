package pgexpert

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/model"
)

func (pg *ExpertPostgres) Delete(ctx context.Context, id int) error {
	queryString := "DELETE FROM experts where id = $1;"
	row := pg.db.QueryRow(queryString, id)
	if row.Err() != nil {
		return fmt.Errorf("[ExpertPostgres.Delete] %w", row.Err())
	}
	return nil
}

func (pg *ExpertPostgres) Insert(ctx context.Context, e *model.Expert) error {
	queryString := "INSERT INTO experts (avatar_url, name, user_id) VALUES ($1, $2, $3);"
	row := pg.db.QueryRow(queryString, e.AvatarUrl, e.Name, e.UserId)
	if row.Err() != nil {
		return fmt.Errorf("[ExpertPostgres.Insert] %w", row.Err())
	}
	return nil
}

func (pg *ExpertPostgres) Update(ctx context.Context, e *model.Expert) error {
	queryString := "UPDATE experts SET avatar_url = $1, name = $2 WHERE id = $3;"
	row := pg.db.QueryRowContext(ctx, queryString, e.AvatarUrl, e.Name, e.Id)
	if row.Err() != nil {
		return fmt.Errorf("[ExpertPostgres.Update] %w", row.Err())
	}
	return nil
}
func (pg *ExpertPostgres) GetOneById(ctx context.Context, id int) (*model.Expert, error) {
	queryString := "SELECT * FROM experts WHERE id = $1;"
	row := pg.db.QueryRowContext(ctx, queryString, id)
	if row.Err() != nil {
		return nil, fmt.Errorf("[ExpertPostgres.GetOneById] %w", row.Err())
	}
	var e Expert
	err := row.Scan(&e.Id, &e.AvatarUrl, &e.Name, &e.UserId)
	if err != nil {
		return nil, fmt.Errorf("[ExpertPostgres.GetOneById] %w", err)
	}
	return FromDBToModelExpert(&e), nil
}
func (pg *ExpertPostgres) GetListByUserId(ctx context.Context, userId int) ([]*model.Expert, error) {
	queryString := "SELECT * FROM experts WHERE user_id = $1;"
	rows, err := pg.db.QueryContext(ctx, queryString, userId)
	if err != nil {
		return nil, fmt.Errorf("[ExpertPostgres.GetListByUserId] %w", err)
	}

	var experts []*model.Expert

	for rows.Next() {
		var e Expert
		err := rows.Scan(&e.Id, &e.AvatarUrl, &e.Name, &e.UserId)
		if err != nil {
			return nil, fmt.Errorf("[ExpertPostgres.GetListByUserId] %w", err)
		}
		experts = append(experts, FromDBToModelExpert(&e))
	}

	return experts, nil
}

type ExpertPostgres struct {
	db *sql.DB
}

func NewExpertPostgres(db *sql.DB) *ExpertPostgres {
	return &ExpertPostgres{db: db}
}
