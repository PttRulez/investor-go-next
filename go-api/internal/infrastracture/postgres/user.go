package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
)

func (pg *UserPostgres) Insert(ctx context.Context, u *entity.User) error {
	const op = "UserPostgres.Insert"

	querySting := "INSERT INTO users (email, hashed_password, name, role) VALUES ($1, $2, $3, $4);"
	_, err := pg.db.ExecContext(ctx, querySting, u.Email, u.HashedPassword, u.Name, u.Role)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *UserPostgres) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	const op = "UserPostgres.GetByEmail"

	querySting := `SELECT * FROM users WHERE email = $1 LIMIT 1;`
	row := pg.db.QueryRowContext(ctx, querySting, email)

	var u entity.User
	err := row.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Name, &u.Role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &u, nil
}

func (pg *UserPostgres) GetByID(ctx context.Context, id int) (*entity.User, error) {
	const op = "UserPostgres.GetByID"

	querySting := `SELECT * FROM users WHERE id = $1 LIMIT 1;`

	row := pg.db.QueryRowContext(ctx, querySting, id)

	var u entity.User
	err := row.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Name, &u.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, database.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &u, nil
}

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}
