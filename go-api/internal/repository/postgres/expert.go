package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/types"
)

type ExpertPostgres struct {
	db *sql.DB
}

func NewExpertPostgres(db *sql.DB) types.ExpertRepository {
	return &ExpertPostgres{db: db}
}

func (pg *ExpertPostgres) Delete(ctx context.Context, id int) error {
	queryString := "DELETE FROM experts where id = $1;"
	row := pg.db.QueryRow(queryString, id)
	if row.Err() != nil {
		return fmt.Errorf("[ExpertPostgres.Delete] %w", row.Err())
	}
	return nil
}

func (pg *ExpertPostgres) Insert(ctx context.Context, e *types.Expert) error {
	queryString := "INSERT INTO experts (avatar_url, name, user_id) VALUES ($1, $2, $3);"
	row := pg.db.QueryRow(queryString, e.AvatarUrl, e.Name, e.UserId)
	if row.Err() != nil {
		return fmt.Errorf("[ExpertPostgres.Insert] %w", row.Err())
	}
	return nil
}

func (pg *ExpertPostgres) Update(ctx context.Context, e *types.Expert) error {
	queryString := "UPDATE experts SET avatar_url = $1, name = $2 WHERE id = $3;"
	row := pg.db.QueryRowContext(ctx, queryString, e.AvatarUrl, e.Name, e.Id)
	if row.Err() != nil {
		return fmt.Errorf("[ExpertPostgres.Update] %w", row.Err())
	}
	return nil
}

func (pg *ExpertPostgres) GetListByUserId(ctx context.Context, userId int) ([]*types.Expert, error) {
	queryString := "SELECT * FROM experts WHERE user_id = $1;"
	rows, err := pg.db.QueryContext(ctx, queryString, userId)
	if err != nil {
		return nil, fmt.Errorf("[ExpertPostgres.GetListByUserId] %w", err)
	}

	var experts []*types.Expert

	for rows.Next() {
		var e types.Expert
		err := rows.Scan(&e.Id, &e.AvatarUrl, &e.Name, &e.UserId)
		if err != nil {
			return nil, fmt.Errorf("[ExpertPostgres.GetListByUserId] %w", err)
		}
		experts = append(experts, &e)
	}

	return experts, nil
}
