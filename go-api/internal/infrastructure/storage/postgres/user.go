package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
)

func (pg *Repository) InsertUser(ctx context.Context, u domain.User) error {
	const op = "Repository.InsertUser"

	querySting := "INSERT INTO users (email, hashed_password, name, role) VALUES ($1, $2, $3, $4);"
	_, err := pg.db.ExecContext(ctx, querySting, u.Email, u.HashedPassword, u.Name, u.Role)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *Repository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	const op = "Repository.GetUserByEmail"

	querySting := `SELECT * FROM users WHERE email = $1 LIMIT 1;`
	row := pg.db.QueryRowContext(ctx, querySting, email)

	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Name, &u.Role, &u.TgChatID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, storage.ErrNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (pg *Repository) GetUser(ctx context.Context, id int) (domain.User, error) {
	const op = "Repository.GetUser"

	querySting := `SELECT * FROM users WHERE id = $1 LIMIT 1;`

	row := pg.db.QueryRowContext(ctx, querySting, id)

	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Name, &u.Role, &u.TgChatID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, storage.ErrNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (pg *Repository) UpdateUser(ctx context.Context, u domain.User) error {
	const op = "Repository.UpdateUser"
	var args []any
	var sets []string

	q := "UPDATE users SET "
	count := 0
	if u.Name != "" {
		count++
		sets = append(sets, fmt.Sprintf("name = $%d", count))
		args = append(args, u.Name)
	}
	if u.TgChatID != nil {
		count++
		sets = append(sets, fmt.Sprintf("invest_bot_tg_chat_id = $%d", count))
		args = append(args, u.TgChatID)
	}

	q = q + strings.Join(sets, ", ") + fmt.Sprintf(" WHERE id = $%d", count+1)

	args = append(args, u.ID)

	_, err := pg.db.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
