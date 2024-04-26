package moexbonddeal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pttrulez/investor-go/internal/model"
)

func (pg *MoexBondDealPostgres) Delete(ctx context.Context, id int) error {
	queryString := `DELETE FROM moex_bond_deals WHERE id = $1;`
	_, err := pg.db.ExecContext(ctx, queryString, id)
	if err != nil {
		return fmt.Errorf("\n<-[MoexBondDealPostgres.Delete]: %w", err)
	}

	return nil
}

func (pg *MoexBondDealPostgres) GetDealListByPortoflioId(ctx context.Context,
	portfolioId int) ([]*model.Deal, error) {
	queryString := `SELECT d.*, s.isin
    FROM moex_bond_deals d 
    LEFT JOIN moex_bonds s ON d.security_id = s.id
		WHERE d.portfolio_id = $1
    ORDER BY d.date DESC, d.id DESC;`
	rows, err := pg.db.QueryContext(ctx, queryString, portfolioId)
	if err != nil {
		return nil, fmt.Errorf("\n<-[MoexBondDealPostgres.GetDealsListByPortoflioId]: %w", err)
	}

	deals := []*model.Deal{}
	for rows.Next() {
		var deal = MoexBondDeal{}
		err := rows.Scan(
			&deal.Id,
			&deal.Amount,
			&deal.Date,
			&deal.PortfolioId,
			&deal.Price,
			&deal.SecurityId,
			&deal.Type,
			&deal.Isin,
		)
		if err != nil {
			return nil, fmt.Errorf("\n<-[MoexBondDealPostgres.GetDealsListByPortoflioId]: %w", err)
		}
		deals = append(deals, FromDBToDeal(&deal))
	}
	return deals, nil
}

func (pg *MoexBondDealPostgres) GetDealListByBondId(ctx context.Context, portfolioId int,
	securityId int) ([]*model.Deal, error) {

	queryString := `SELECT d.*, s.isin
    FROM moex_bond_deals d 
    LEFT JOIN moex_bonds s ON d.security_id = s.id
		WHERE d.security_id = $1 AND d.portfolio_id = $2
    ORDER BY d.date DESC, d.id DESC;`
	rows, err := pg.db.QueryContext(ctx, queryString, securityId, portfolioId)
	if err != nil {
		return nil, fmt.Errorf("\n<-[MoexBondDealPostgres.GetDealListByBondId]: %w", err)
	}

	deals := []*model.Deal{}
	for rows.Next() {
		var deal = MoexBondDeal{}
		err := rows.Scan(
			&deal.Id,
			&deal.Amount,
			&deal.Date,
			&deal.PortfolioId,
			&deal.Price,
			&deal.SecurityId,
			&deal.Type,
			&deal.Isin,
		)
		if err != nil {
			return nil, fmt.Errorf("\n<-[MoexBondDealPostgres.GetDealsBySecurityId]: %w", err)
		}
		deals = append(deals, FromDBToDeal(&deal))
	}
	return deals, nil
}

func (pg *MoexBondDealPostgres) Insert(ctx context.Context, deal *model.Deal) error {
	d := FromDealToDB(deal)
	queryString := `INSERT INTO moex_bond_deals (amount, date, portfolio_id, price,
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
		return fmt.Errorf("\n<-[MoexBondDealPostgres.Insert]: %w", err)
	}
	return nil
}

func (pg *MoexBondDealPostgres) Update(ctx context.Context, deal *model.Deal) error {
	d := FromDealToDB(deal)
	queryString := `UPDATE moex_bond_deals SET amount = $1, date = $2,	portfolio_id = $4,
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
		return fmt.Errorf("\n<-[MoexBondDealPostgres.Update]: %w", err)
	}
	return nil
}

type MoexBondDealPostgres struct {
	db *sql.DB
}

func NewMoexBondDealPostgres(db *sql.DB) *MoexBondDealPostgres {
	return &MoexBondDealPostgres{db: db}
}
