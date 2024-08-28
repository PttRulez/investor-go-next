package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
)

func (pg *PortfolioPostgres) Delete(ctx context.Context, id int, userID int) error {
	const op = "PortfolioPostgres.Delete"

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
		return database.ErrNotFound
	}

	return nil
}

func (pg *PortfolioPostgres) GetByID(ctx context.Context, id int, userID int) (entity.Portfolio, error) {
	const op = "PortfolioPostgres.GetByID"

	queryString := `SELECT * FROM portfolios where id = $1 AND user_id = $2;`

	var p entity.Portfolio
	err := pg.db.QueryRowContext(ctx, queryString, id, userID).Scan(&p.ID, &p.Compound, &p.Name, &p.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Portfolio{}, database.ErrNotFound
		}
		return entity.Portfolio{}, fmt.Errorf("%s: %w", op, err)
	}

	return p, nil
}

func (pg *PortfolioPostgres) GetListByUserID(ctx context.Context, userID int) ([]entity.Portfolio, error) {
	const op = "PortfolioPostgres.GetListByUserId"

	queryString := `SELECT id, compound, name FROM portfolios where user_id = $1;`
	rows, err := pg.db.QueryContext(ctx, queryString, strconv.Itoa(userID))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var portfolios []entity.Portfolio
	for rows.Next() {
		var p entity.Portfolio
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

func (pg *PortfolioPostgres) Insert(ctx context.Context, p entity.Portfolio) (entity.Portfolio, error) {
	const op = "PortfolioPostgres.Insert"

	queryString := `INSERT INTO portfolios (compound, name, user_id) VALUES ($1, $2, $3)
		RETURNING id, compound, name;`

	var result entity.Portfolio
	err := pg.db.QueryRowContext(ctx, queryString, p.Compound, p.Name, p.UserID).
		Scan(&result.ID, &result.Compound, &result.Name)
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return entity.Portfolio{}, nil
}

func (pg *PortfolioPostgres) Update(ctx context.Context, p entity.Portfolio, userID int) (entity.Portfolio, error) {
	const op = "PortfolioPostgres.Update"

	queryString := `UPDATE portfolios SET compound = $1, name = $2 WHERE id = $3 AND user_id = $4
		RETURNING id, compound, name;`

	var up entity.Portfolio
	err := pg.db.QueryRowContext(ctx, queryString, p.Compound, p.Name, p.ID, userID).
		Scan(&up.ID, up.Compound, up.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Portfolio{}, database.ErrNotFound
	}
	if err != nil {
		return entity.Portfolio{}, fmt.Errorf("%s: %w", op, err)
	}

	return up, nil
}

type PortfolioPostgres struct {
	db *sql.DB
}

func NewPortfolioPostgres(db *sql.DB) *PortfolioPostgres {
	return &PortfolioPostgres{db: db}
}
