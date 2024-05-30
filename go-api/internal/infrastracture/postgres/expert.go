package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
)

func (pg *ExpertPostgres) Delete(ctx context.Context, id int) error {
	queryString := "DELETE FROM experts where id = $1;"
	row := pg.db.QueryRow(queryString, id)
	if row.Err() != nil {
		return fmt.Errorf("[ExpertPostgres.Delete] %w", row.Err())
	}
	return nil
}
func (pg *ExpertPostgres) Insert(ctx context.Context, e *entity.Expert) error {
	queryString := "INSERT INTO experts (avatar_url, name, user_id) VALUES ($1, $2, $3);"
	row := pg.db.QueryRow(queryString, e.AvatarUrl, e.Name, e.UserId)
	if row.Err() != nil {
		return fmt.Errorf("[ExpertPostgres.Insert] %w", row.Err())
	}
	return nil
}
func (pg *ExpertPostgres) Update(ctx context.Context, e *entity.Expert) error {
	queryString := "UPDATE experts SET avatar_url = $1, name = $2 WHERE id = $3;"
	row := pg.db.QueryRowContext(ctx, queryString, e.AvatarUrl, e.Name, e.Id)
	if row.Err() != nil {
		return fmt.Errorf("[ExpertPostgres.Update] %w", row.Err())
	}
	return nil
}
func (pg *ExpertPostgres) GetOneById(ctx context.Context, id int) (*entity.Expert, error) {
	queryString := "SELECT * FROM experts WHERE id = $1;"

	row := pg.db.QueryRowContext(ctx, queryString, id)
	if row.Err() != nil {
		return nil, fmt.Errorf("[ExpertPostgres.GetOneById] %w", row.Err())
	}

	var e entity.Expert
	err := row.Scan(&e.Id, &e.AvatarUrl, &e.Name, &e.UserId)
	if err != nil {
		return nil, fmt.Errorf("[ExpertPostgres.GetOneById] %w", err)
	}

	return &e, nil
}
func (pg *ExpertPostgres) GetListByUserId(ctx context.Context, userId int) ([]*entity.Expert, error) {
	queryString := "SELECT * FROM experts WHERE user_id = $1;"
	rows, err := pg.db.QueryContext(ctx, queryString, userId)
	if err != nil {
		return nil, fmt.Errorf("[ExpertPostgres.GetListByUserId] %w", err)
	}

	var experts []*entity.Expert

	for rows.Next() {
		var e entity.Expert
		err := rows.Scan(&e.Id, &e.AvatarUrl, &e.Name, &e.UserId)
		if err != nil {
			return nil, fmt.Errorf("[ExpertPostgres.GetListByUserId] %w", err)
		}
		experts = append(experts, &e)
	}

	return experts, nil
}

type ExpertPostgres struct {
	db *sql.DB
}

func NewExpertPostgres(db *sql.DB) *ExpertPostgres {
	return &ExpertPostgres{db: db}
}
