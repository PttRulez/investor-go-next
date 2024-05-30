package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
)

func (pg *PortfolioPostgres) Delete(ctx context.Context, id int, userId int) error {
	queryString := "DELETE FROM portfolios where id = $1 AND user_id = $2;"
	_, err := pg.db.ExecContext(ctx, queryString, id, userId)
	if err != nil {
		return fmt.Errorf("[PortfolioPostgres Delete]: %w", err)
	}
	return nil
}

func (pg *PortfolioPostgres) GetById(ctx context.Context, id int, userId int) (*entity.Portfolio, error) {
	queryString := `SELECT * FROM portfolios where id = $1 AND user_id = $2;`
	row := pg.db.QueryRowContext(ctx, queryString, id, userId)
	if row.Err() != nil {
		return nil, fmt.Errorf("[PortfolioPostgres GetById]: %w", row.Err())
	}

	var p entity.Portfolio
	err := row.Scan(&p.Id, &p.Compound, &p.Name, &p.UserId)
	if err != nil {
		return nil, fmt.Errorf("[PortfolioPostgres GetById]: %w", err)
	}
	return &p, nil
}

func (pg *PortfolioPostgres) GetListByUserId(ctx context.Context, id int) ([]*entity.Portfolio, error) {
	queryString := `SELECT * FROM portfolios where user_id = $1;`
	rows, err := pg.db.QueryContext(ctx, queryString, fmt.Sprint(id))
	if err != nil {
		return nil, fmt.Errorf("[PortfolioPostgres GetListByUserId]: %w", err)
	}
	defer rows.Close()

	var portfolios []*entity.Portfolio
	for rows.Next() {
		p := entity.Portfolio{}
		err := rows.Scan(&p.Id, &p.Compound, &p.Name, &p.UserId)
		if err != nil {
			return nil, fmt.Errorf("[PortfolioPostgres GetListByUserId] %w", err)
		}
		portfolios = append(portfolios, &p)
	}

	return portfolios, nil
}

func (pg *PortfolioPostgres) Insert(ctx context.Context, p *entity.Portfolio) error {
	queryString := "INSERT INTO portfolios (compound, name, user_id) VALUES ($1, $2, $3);"
	_, err := pg.db.ExecContext(ctx, queryString, p.Compound, p.Name, p.UserId)
	if err != nil {
		return fmt.Errorf("[PortfolioPostgres Insert]: %w", err)
	}
	return nil
}

func (pg *PortfolioPostgres) Update(ctx context.Context, p *entity.Portfolio, userId int) error {
	queryString := "UPDATE portfolios SET compound = $1, name = $2 WHERE id = $3 AND user_id = $4;"
	_, err := pg.db.ExecContext(ctx, queryString, p.Compound, p.Name, p.Id, userId)
	if err != nil {
		return fmt.Errorf("[PortfolioPostgres Update]: %w", err)
	}
	return nil
}

type PortfolioPostgres struct {
	db *sql.DB
}

func NewPortfolioPostgres(db *sql.DB) *PortfolioPostgres {
	return &PortfolioPostgres{db: db}
}
