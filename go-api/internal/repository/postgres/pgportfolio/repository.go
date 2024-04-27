package pgportfolio

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/model"
)

func (pg *PortfolioPostgres) Delete(ctx context.Context, id int) error {
	queryString := "DELETE FROM portfolios where id = $1;"
	_, err := pg.db.ExecContext(ctx, queryString, fmt.Sprint(id))
	if err != nil {
		return fmt.Errorf("[PortfolioPostgres Delete]: %w", err)
	}
	return nil
}

func (pg *PortfolioPostgres) GetById(ctx context.Context, id int) (*model.Portfolio, error) {
	queryString := `SELECT * FROM portfolios where id = $1;`
	row := pg.db.QueryRowContext(ctx, queryString, fmt.Sprint(id))
	if row.Err() != nil {
		return nil, fmt.Errorf("[PortfolioPostgres GetById]: %w", row.Err())
	}

	var p Portfolio
	err := row.Scan(&p.Id, &p.Compound, &p.Name, &p.UserId)
	if err != nil {
		return nil, fmt.Errorf("[PortfolioPostgres GetById]: %w", err)
	}
	return FromDBToPortfolio(&p), nil
}

func (pg *PortfolioPostgres) GetListByUserId(ctx context.Context, id int) ([]*model.Portfolio, error) {
	queryString := `SELECT * FROM portfolios where user_id = $1;`
	rows, err := pg.db.QueryContext(ctx, queryString, fmt.Sprint(id))
	if err != nil {
		return nil, fmt.Errorf("[PortfolioPostgres GetListByUserId]: %w", err)
	}
	defer rows.Close()

	portfolios := []*model.Portfolio{}
	for rows.Next() {
		p := Portfolio{}
		err := rows.Scan(&p.Id, &p.Compound, &p.Name, &p.UserId)
		if err != nil {
			return nil, fmt.Errorf("[PortfolioPostgres GetListByUserId] %w", err)
		}
		portfolios = append(portfolios, FromDBToPortfolio(&p))
	}

	return portfolios, nil
}

func (pg *PortfolioPostgres) Insert(ctx context.Context, portfolio *model.Portfolio) error {
	p := FromPortfolioToDB(portfolio)
	queryString := "INSERT INTO portfolios (compound, name, user_id) VALUES ($1, $2, $3);"
	_, err := pg.db.ExecContext(ctx, queryString, p.Compound, p.Name, p.UserId)
	if err != nil {
		return fmt.Errorf("[PortfolioPostgres Insert]: %w", err)
	}
	return nil
}

func (pg *PortfolioPostgres) Update(ctx context.Context, portfolio *model.Portfolio) error {
	p := FromPortfolioToDB(portfolio)
	queryString := "UPDATE portfolios SET compound = $1, name = $2 WHERE id = $3;"
	_, err := pg.db.ExecContext(ctx, queryString, p.Compound, p.Name, p.Id)
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
