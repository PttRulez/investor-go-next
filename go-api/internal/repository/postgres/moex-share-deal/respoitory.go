package moexsharedeal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/model"
)

func (pg *MoexShareDealPostgres) Delete(ctx context.Context, id int) error {
	queryString := `DELETE FROM moex_share_deals WHERE id = $1;`
	_, err := pg.db.ExecContext(ctx, queryString, id)
	if err != nil {
		return fmt.Errorf("\n<-[MoexShareDealPostgres.Delete]: %w", err)
	}

	return nil
}

func (pg *MoexShareDealPostgres) GetDealListByPortoflioId(ctx context.Context,
	portfolioId int) ([]*model.Deal, error) {
	queryString := `SELECT d.*, s.ticker
    FROM moex_share_deals d 
    LEFT JOIN moex_shares s ON d.security_id = s.id
		WHERE d.portfolio_id = $1
    ORDER BY d.date DESC, d.id DESC;`
	rows, err := pg.db.QueryContext(ctx, queryString, portfolioId)
	if err != nil {
		return nil, fmt.Errorf("<-[MoexShareDealPostgres.GetDealsListByPortoflioId]: \n%w", err)
	}

	deals := []*model.Deal{}
	for rows.Next() {
		var deal = MoexShareDeal{}
		err := rows.Scan(
			&deal.Id,
			&deal.Amount,
			&deal.Date,
			&deal.PortfolioId,
			&deal.Price,
			&deal.SecurityId,
			&deal.Type,
			&deal.Ticker,
		)
		if err != nil {
			return nil, fmt.Errorf("\n<-[MoexShareDealPostgres.GetDealsListByPortoflioId]: \n%w", err)
		}
		deals = append(deals, FromDBToMoexShareDeal(&deal))
	}
	return deals, nil
}

func (pg *MoexShareDealPostgres) GetDealListByShareId(ctx context.Context,
	portfolioId int, securityId int) ([]*model.Deal, error) {
	queryString := `SELECT d.id, d.amount, d.date, d.portfolio_id, d.price, 
		d.security_id, d.type, s.ticker
    FROM moex_share_deals d 
    LEFT JOIN moex_shares s ON d.security_id = s.id
		WHERE d.security_id = $1 AND d.portfolio_id = $2
    ORDER BY id DESC;`
	rows, err := pg.db.QueryContext(ctx, queryString, securityId, portfolioId)
	if err != nil {
		return nil, fmt.Errorf("\n<-[Postgres.MoexShareDeal.GetDealsListForSecurity]: %w", err)
	}

	deals := make([]*model.Deal, 0)
	for rows.Next() {
		var deal = MoexShareDeal{}
		err := rows.Scan(
			&deal.Id,
			&deal.Amount,
			&deal.Date,
			&deal.PortfolioId,
			&deal.Price,
			&deal.SecurityId,
			&deal.Type,
			&deal.Ticker,
		)
		if err != nil {
			return nil, fmt.Errorf("\n<-[Postgres.MoexShareDeal.GetDealsListForSecurity]: %w", err)
		}
		deals = append(deals, FromDBToMoexShareDeal(&deal))
	}
	return deals, nil
}

func (pg *MoexShareDealPostgres) Insert(ctx context.Context, deal *model.Deal) error {
	d := FromMoexShareDealToDB(deal)
	queryString := `INSERT INTO moex_share_deals (amount, date, portfolio_id, price,
    security_id, type) VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := pg.db.ExecContext(
		ctx,
		queryString,
		d.Amount,
		d.Date,
		d.PortfolioId,
		d.Price,
		d.SecurityId,
		d.Type,
	)
	if err != nil {
		return fmt.Errorf("\n<-[Postgres.MoexShareDeal.Insert]: %w", err)
	}
	return nil
}

func (pg *MoexShareDealPostgres) Update(ctx context.Context, deal *model.Deal) error {
	d := FromMoexShareDealToDB(deal)
	queryString := `UPDATE moex_share_deals SET amount = $1, date = $2,	portfolio_id = $4,
   price = $5, security_id = $6, type = $7 WHERE id = $8;`
	_, err := pg.db.ExecContext(
		ctx,
		queryString,
		d.Amount,
		d.Date,
		d.PortfolioId,
		d.Price,
		d.SecurityId,
		d.Type,
		d.Id,
	)
	if err != nil {
		return fmt.Errorf("\n<-[Postgres.MoexShareDeal.Update]: %w", err)
	}
	return nil
}

type MoexShareDealPostgres struct {
	db *sql.DB
}

func NewMoexShareDealPostgres(db *sql.DB) *MoexShareDealPostgres {
	return &MoexShareDealPostgres{db: db}
}
