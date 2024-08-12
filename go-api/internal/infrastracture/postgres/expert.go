package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
)

func (pg *ExpertPostgres) Delete(ctx context.Context, id int) error {
	const op = "ExpertPostgres.Delete"

	queryString := "DELETE FROM experts where id = $1;"
	_, err := pg.db.ExecContext(ctx, queryString, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *ExpertPostgres) Insert(ctx context.Context, e *entity.Expert) error {
	const op = "ExpertPostgres.Insert"

	queryString := "INSERT INTO experts (avatar_url, name, user_id) VALUES ($1, $2, $3);"

	_, err := pg.db.ExecContext(ctx, queryString, e.AvatarURL, e.Name, e.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *ExpertPostgres) Update(ctx context.Context, e *entity.Expert) error {
	const op = "ExpertPostgres.Update"

	queryString := "UPDATE experts SET avatar_url = $1, name = $2 WHERE id = $3;"
	_, err := pg.db.ExecContext(ctx, queryString, e.AvatarURL, e.Name, e.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *ExpertPostgres) GetOneByID(ctx context.Context, id int) (*entity.Expert, error) {
	const op = "ExpertPostgres.GetOneByID"

	queryString := "SELECT * FROM experts WHERE id = $1;"

	row := pg.db.QueryRowContext(ctx, queryString, id)
	if row.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, row.Err())
	}

	var e entity.Expert
	err := row.Scan(&e.ID, &e.AvatarURL, &e.Name, &e.UserID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &e, nil
}

func (pg *ExpertPostgres) GetListByUserID(ctx context.Context, userID int) ([]*entity.Expert, error) {
	const op = "ExpertPostgres.GetListByUserID"

	queryString := "SELECT * FROM experts WHERE user_id = $1;"
	rows, err := pg.db.QueryContext(ctx, queryString, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var experts []*entity.Expert

	for rows.Next() {
		var e entity.Expert
		err = rows.Scan(&e.ID, &e.AvatarURL, &e.Name, &e.UserID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		experts = append(experts, &e)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return experts, nil
}

type ExpertPostgres struct {
	db *sql.DB
}

func NewExpertPostgres(db *sql.DB) *ExpertPostgres {
	return &ExpertPostgres{db: db}
}
