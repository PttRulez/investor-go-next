package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/types"
)

type PortfolioPostgres struct {
	db *sql.DB
}

func NewPortfolioPostgres(db *sql.DB) types.PortfolioRepository {
	return &PortfolioPostgres{db: db}
}

func (pg *PortfolioPostgres) Delete(ctx context.Context, id int) error {
	queryString := "DELETE FROM portfolios where id = $1;"
	_, err := pg.db.ExecContext(ctx, queryString, fmt.Sprint(id))
	if err != nil {
		return fmt.Errorf("[PortfolioPostgres Delete]: %w", err)
	}
	return nil
}

func (pg *PortfolioPostgres) GetById(ctx context.Context, id int) (*types.Portfolio, error) {
	queryString := `SELECT * FROM portfolios where id = $1;`
	// JOIN deals WHERE portfolios.id = deals.portfolio_id
	row := pg.db.QueryRowContext(ctx, queryString, fmt.Sprint(id))
	if row.Err() != nil {
		return nil, fmt.Errorf("[PortfolioPostgres GetById]: %w", row.Err())
	}

	var p types.Portfolio
	err := row.Scan(&p.Id, &p.Compound, &p.Name, &p.UserId)
	if err != nil {
		return nil, fmt.Errorf("[PortfolioPostgres GetById]: %w", err)
	}
	return &p, nil
}

func (pg *PortfolioPostgres) GetByIdAndScan(ctx context.Context, id int, p *types.FullPortfolioData) error {
	queryString := `SELECT * FROM portfolios where id = $1;`
	// JOIN deals WHERE portfolios.id = deals.portfolio_id
	row := pg.db.QueryRowContext(ctx, queryString, fmt.Sprint(id))
	if row.Err() != nil {
		return fmt.Errorf("[PortfolioPostgres GetByIdAndScan]: %w", row.Err())
	}

	err := row.Scan(&p.Id, &p.Compound, &p.Name, &p.UserId)
	if err != nil {
		return fmt.Errorf("[PortfolioPostgres GetByIdAndScan]: %w", err)
	}
	return nil
}

func (pg *PortfolioPostgres) GetListByUserId(ctx context.Context, id int) ([]*types.Portfolio, error) {
	queryString := `SELECT * FROM portfolios where user_id = $1;`
	// JOIN deals WHERE portfolios.id = deals.portfolio_id
	rows, err := pg.db.QueryContext(ctx, queryString, fmt.Sprint(id))
	if err != nil {
		return nil, fmt.Errorf("[PortfolioPostgres GetListByUserId]: %w", err)
	}
	defer rows.Close()

	portfolios := []*types.Portfolio{}
	for rows.Next() {
		p := &types.Portfolio{}
		err := rows.Scan(&p.Id, &p.Compound, &p.Name, &p.UserId)
		if err != nil {
			return nil, fmt.Errorf("[PortfolioPostgres GetListByUserId] %w", err)
		}
		portfolios = append(portfolios, p)
	}

	return portfolios, nil
}

func (pg *PortfolioPostgres) Insert(ctx context.Context, u *types.Portfolio) error {
	queryString := "INSERT INTO portfolios (compound, name, user_id) VALUES ($1, $2, $3);"
	_, err := pg.db.ExecContext(ctx, queryString, u.Compound, u.Name, u.UserId)
	if err != nil {
		return fmt.Errorf("[PortfolioPostgres Insert]: %w", err)
	}
	return nil
}

func (pg *PortfolioPostgres) Update(ctx context.Context, data *types.PortfolioUpdate) error {
	queryString := "UPDATE portfolios SET compound = $1, name = $2 WHERE id = $3;"
	_, err := pg.db.ExecContext(ctx, queryString, data.Compound, data.Name, data.Id)
	if err != nil {
		return fmt.Errorf("[PortfolioPostgres Update]: %w", err)
	}
	return nil
}
