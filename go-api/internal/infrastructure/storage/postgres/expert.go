package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/storage"
)

func (pg *Repository) DeleteExpert(ctx context.Context, id int, userID int) error {
	const op = "Repository.DeleteExpert"

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
		return storage.ErrNotFound
	}

	return nil
}

func (pg *Repository) GetExpert(ctx context.Context, id int, userID int) (domain.Expert, error) {
	const op = "Repository.GetExpert"

	queryString := "SELECT * FROM experts WHERE id = $1 AND user_id = $2;"

	var e domain.Expert

	err := pg.db.QueryRowContext(ctx, queryString, id, userID).
		Scan(&e.ID, &e.AvatarURL, &e.Name, &e.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Expert{}, storage.ErrNotFound
	}
	if err != nil {
		return domain.Expert{}, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}

func (pg *Repository) GetExpertList(ctx context.Context, userID int) ([]domain.Expert, error) {
	const op = "Repository.GetExpertList"

	queryString := "SELECT * FROM experts WHERE user_id = $1;"
	rows, err := pg.db.QueryContext(ctx, queryString, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var experts []domain.Expert

	for rows.Next() {
		var e domain.Expert
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

func (pg *Repository) InsertExpert(ctx context.Context, e domain.Expert) (domain.Expert, error) {
	const op = "Repository.InsertExpert"

	queryString := `INSERT INTO experts (avatar_url, name, user_id) VALUES ($1, $2, $3)
		RETURNING id, avatar_url, name;`

	var expert domain.Expert
	err := pg.db.QueryRowContext(ctx, queryString, e.AvatarURL, e.Name, e.UserID).
		Scan(&expert.ID, &expert.AvatarURL, &expert.Name)
	if err != nil {
		return expert, fmt.Errorf("%s: %w", op, err)
	}

	return expert, nil
}

func (pg *Repository) UpdateExpert(ctx context.Context, e domain.Expert, userID int) (domain.Expert, error) {
	const op = "Repository.UpdateExpert"

	queryString := `UPDATE experts SET avatar_url = $1, name = $2 WHERE id = $3 AND user_id = $4
		RETURNING id, avatar_url, name;`

	var ue domain.Expert
	err := pg.db.QueryRowContext(ctx, queryString, e.AvatarURL, e.Name, e.ID, userID).
		Scan(&ue.ID, &ue.AvatarURL, &ue.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Expert{}, storage.ErrNotFound
	}
	if err != nil {
		return domain.Expert{}, fmt.Errorf("%s: %w", op, err)
	}

	return ue, nil
}
