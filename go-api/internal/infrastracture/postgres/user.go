package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
)

func (pg *UserPostgres) Insert(ctx context.Context, u *entity.User) error {
	querySting := "INSERT INTO users (email, hashed_password, name, role) VALUES ($1, $2, $3, $4);"
	_, err := pg.db.ExecContext(ctx, querySting, u.Email, u.HashedPassword, u.Name, u.Role)
	if err != nil {
		return fmt.Errorf("[UserPostgres.Insert]: %w", err)
	}

	return nil
}

func (pg *UserPostgres) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	querySting := `SELECT * FROM users WHERE email = $1 LIMIT 1;`
	row := pg.db.QueryRowContext(ctx, querySting, email)
	if row.Err() != nil {
		return nil, fmt.Errorf("[UserPostgres.GetByEmail]: %w", row.Err())
	}

	var u entity.User
	err := row.Scan(&u.Id, &u.Email, &u.HashedPassword, &u.Name, &u.Role)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		break
	default:
		return nil, fmt.Errorf("[UserPostgres.GetByEmail]: %w", err)
	}

	return &u, nil
}

func (pg *UserPostgres) GetById(ctx context.Context, id int) (*entity.User, error) {
	querySting := `SELECT * FROM users WHERE id = $1 LIMIT 1;`
	row := pg.db.QueryRowContext(ctx, querySting, id)
	if row.Err() != nil {
		return nil, fmt.Errorf("[UserPostgres.GetById]: %w", row.Err())
	}

	var u entity.User
	err := row.Scan(&u.Id, &u.Email, &u.HashedPassword, &u.Name, &u.Role)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		break
	default:
		return nil, err
	}

	return &u, nil
}

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}
