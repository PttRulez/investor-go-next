package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/pttrulez/investor-go/internal/entity"
)

func (pg *PortfolioPostgres) Delete(ctx context.Context, id int, userID int) error {
	const op = "PortfolioPostgres.Delete"

	queryString := "DELETE FROM portfolios where id = $1 AND user_id = $2;"

	_, err := pg.db.ExecContext(ctx, queryString, id, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *PortfolioPostgres) GetByID(ctx context.Context, id int, userID int) (*entity.Portfolio, error) {
	const op = "PortfolioPostgres.GetByID"

	queryString := `SELECT * FROM portfolios where id = $1 AND user_id = $2;`

	row := pg.db.QueryRowContext(ctx, queryString, id, userID)
	if row.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, row.Err())
	}

	var p entity.Portfolio

	err := row.Scan(&p.ID, &p.Compound, &p.Name, &p.UserID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &p, nil
}

func (pg *PortfolioPostgres) GetListByUserID(ctx context.Context, id int) ([]*entity.Portfolio, error) {
	const op = "PortfolioPostgres.GetListByUserId"

	queryString := `SELECT * FROM portfolios where user_id = $1;`
	rows, err := pg.db.QueryContext(ctx, queryString, strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var portfolios []*entity.Portfolio
	for rows.Next() {
		p := entity.Portfolio{}
		err = rows.Scan(&p.ID, &p.Compound, &p.Name, &p.UserID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		portfolios = append(portfolios, &p)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return portfolios, nil
}

func (pg *PortfolioPostgres) Insert(ctx context.Context, p *entity.Portfolio) error {
	const op = "PortfolioPostgres.Insert"

	queryString := "INSERT INTO portfolios (compound, name, user_id) VALUES ($1, $2, $3);"
	_, err := pg.db.ExecContext(ctx, queryString, p.Compound, p.Name, p.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (pg *PortfolioPostgres) Update(ctx context.Context, p *entity.Portfolio, userID int) error {
	const op = "PortfolioPostgres.Update"

	queryString := "UPDATE portfolios SET compound = $1, name = $2 WHERE id = $3 AND user_id = $4;"
	_, err := pg.db.ExecContext(ctx, queryString, p.Compound, p.Name, p.ID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

type PortfolioPostgres struct {
	db *sql.DB
}

func NewPortfolioPostgres(db *sql.DB) *PortfolioPostgres {
	return &PortfolioPostgres{db: db}
}
