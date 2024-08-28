package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
)

func (pg *ExpertPostgres) Delete(ctx context.Context, id int, userID int) error {
	const op = "ExpertPostgres.Delete"

	queryString := "DELETE FROM experts where id = $1 AND user_id = $2;"

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

func (pg *ExpertPostgres) GetOneByID(ctx context.Context, id int, userID int) (entity.Expert, error) {
	const op = "ExpertPostgres.GetOneByID"

	queryString := "SELECT * FROM experts WHERE id = $1 AND user_id = $2;"

	var e entity.Expert

	err := pg.db.QueryRowContext(ctx, queryString, id, userID).
		Scan(&e.ID, &e.AvatarURL, &e.Name, &e.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Expert{}, database.ErrNotFound
	}
	if err != nil {
		return entity.Expert{}, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}

func (pg *ExpertPostgres) GetListByUserID(ctx context.Context, userID int) ([]entity.Expert, error) {
	const op = "ExpertPostgres.GetListByUserID"

	queryString := "SELECT * FROM experts WHERE user_id = $1;"
	rows, err := pg.db.QueryContext(ctx, queryString, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var experts []entity.Expert

	for rows.Next() {
		var e entity.Expert
		err = rows.Scan(&e.ID, &e.AvatarURL, &e.Name, &e.UserID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		experts = append(experts, e)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return experts, nil
}

func (pg *ExpertPostgres) Insert(ctx context.Context, e entity.Expert) (entity.Expert, error) {
	const op = "ExpertPostgres.Insert"

	queryString := `INSERT INTO experts (avatar_url, name, user_id) VALUES ($1, $2, $3)
		RETURNING id, avatar_url, name;`

	var expert entity.Expert
	err := pg.db.QueryRowContext(ctx, queryString, e.AvatarURL, e.Name, e.UserID).
		Scan(&expert.ID, &expert.AvatarURL, &expert.Name)
	if err != nil {
		return expert, fmt.Errorf("%s: %w", op, err)
	}

	return expert, nil
}

func (pg *ExpertPostgres) Update(ctx context.Context, e entity.Expert, userID int) (entity.Expert, error) {
	const op = "ExpertPostgres.Update"

	queryString := `UPDATE experts SET avatar_url = $1, name = $2 WHERE id = $3 AND user_id = $4
		RETURNING id, avatar_url, name;`

	var ue entity.Expert
	err := pg.db.QueryRowContext(ctx, queryString, e.AvatarURL, e.Name, e.ID, userID).
		Scan(&ue.ID, &ue.AvatarURL, &ue.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Expert{}, database.ErrNotFound
	}
	if err != nil {
		return entity.Expert{}, fmt.Errorf("%s: %w", op, err)
	}

	return ue, nil
}

type ExpertPostgres struct {
	db *sql.DB
}

func NewExpertPostgres(db *sql.DB) *ExpertPostgres {
	return &ExpertPostgres{db: db}
}
