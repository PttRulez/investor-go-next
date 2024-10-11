package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
)

func (pg *Repository) DeletePortfolio(ctx context.Context, id int, userID int) error {
	const op = "Repository.DeletePortfolio"

	queryString := "DELETE FROM portfolios where id = $1 AND user_id = $2;"

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

func (pg *Repository) GetPortfolio(ctx context.Context, id int, userID int) (
	domain.Portfolio, error) {
	const op = "Repository.GetPortfolio"

	queryString := `SELECT * FROM portfolios where id = $1 AND user_id = $2;`

	var p domain.Portfolio
	err := pg.db.QueryRowContext(ctx, queryString, id, userID).Scan(&p.ID, &p.Compound,
		&p.Name, &p.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Portfolio{}, storage.ErrNotFound
		}
		return domain.Portfolio{}, fmt.Errorf("%s: %w", op, err)
	}

	return p, nil
}

func (pg *Repository) GetPortfolioList(ctx context.Context, userID int) (
	[]domain.Portfolio, error) {
	const op = "Repository.GetPortfolioList"

	queryString := `SELECT id, compound, name FROM portfolios where user_id = $1;`
	rows, err := pg.db.QueryContext(ctx, queryString, strconv.Itoa(userID))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var portfolios []domain.Portfolio
	for rows.Next() {
		var p domain.Portfolio
		err = rows.Scan(&p.ID, &p.Compound, &p.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		portfolios = append(portfolios, p)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return portfolios, nil
}

func (pg *Repository) GetPortfolioListByChatID(ctx context.Context, chatId string) (
	[]domain.Portfolio, error) {
	const op = "Repository.GetPortfolioList"

	queryString := `SELECT id, compound, name FROM portfolios where user_id = (
		SELECT id FROM users WHERE invest_bot_tg_chat_id = $1 LIMIT 1
	);`

	rows, err := pg.db.QueryContext(ctx, queryString, chatId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var portfolios []domain.Portfolio
	for rows.Next() {
		var p domain.Portfolio
		err = rows.Scan(&p.ID, &p.Compound, &p.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		portfolios = append(portfolios, p)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return portfolios, nil
}

func (pg *Repository) InsertPortfolio(ctx context.Context, p domain.Portfolio) (domain.Portfolio, error) {
	const op = "Repository.InsertPortfolio"

	queryString := `INSERT INTO portfolios (compound, name, user_id) VALUES ($1, $2, $3)
		RETURNING id, compound, name;`

	var result domain.Portfolio
	err := pg.db.QueryRowContext(ctx, queryString, p.Compound, p.Name, p.UserID).
		Scan(&result.ID, &result.Compound, &result.Name)
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return domain.Portfolio{}, nil
}

func (pg *Repository) UpdatePortfolio(ctx context.Context, p domain.Portfolio, userID int) (domain.Portfolio, error) {
	const op = "Repository.UpdatePortfolio"

	queryString := `UPDATE portfolios SET compound = $1, name = $2 WHERE id = $3 AND user_id = $4
		RETURNING id, compound, name;`

	var up domain.Portfolio
	err := pg.db.QueryRowContext(ctx, queryString, p.Compound, p.Name, p.ID, userID).
		Scan(&up.ID, up.Compound, up.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Portfolio{}, storage.ErrNotFound
	}
	if err != nil {
		return domain.Portfolio{}, fmt.Errorf("%s: %w", op, err)
	}

	return up, nil
}
