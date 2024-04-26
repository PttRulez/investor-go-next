package repository

import (
	"context"

	"github.com/pttrulez/investor-go/internal/model"
)

type Repository struct {
	Cashout   CashoutRepository
	Deal      DealRepository
	Deposit   DepositRepository
	Expert    ExpertRepository
	Moex      MoexRepository
	Portfolio PortfolioRepository
	Position  PositionRepository
	User      UserRepository
}
type DealRepository struct {
	MoexBond  MoexBondDealRepository
	MoexShare MoexShareDealRepository
}
type MoexRepository struct {
	Bond  MoexBondRepository
	Share MoexShareRepository
}
type PositionRepository struct {
	MoexBond  MoexBondPositionRepository
	MoexShare MoexSharePositionRepository
}

type CashoutRepository interface {
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*model.Cashout, error)
	GetListByPortfolioId(ctx context.Context, id int) ([]*model.Cashout, error)
	Insert(ctx context.Context, c *model.Cashout) error
}
type DepositRepository interface {
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*model.Deposit, error)
	GetListByPortfolioId(ctx context.Context, id int) ([]*model.Deposit, error)
	Insert(ctx context.Context, c *model.Deposit) error
}
type MoexBondDealRepository interface {
	Delete(ctx context.Context, id int) error
	GetDealListByBondId(ctx context.Context, portfolioId int,
		securityId int) ([]*model.Deal, error)
	GetDealListByPortoflioId(ctx context.Context, portfolioId int) ([]*model.Deal, error)
	Insert(ctx context.Context, d *model.Deal) error
	Update(ctx context.Context, d *model.Deal) error
}
type MoexShareDealRepository interface {
	Delete(ctx context.Context, id int) error
	GetDealListByShareId(ctx context.Context, portfolioId int,
		securityId int) ([]*model.Deal, error)
	GetDealListByPortoflioId(ctx context.Context, portfolioId int) ([]*model.Deal, error)
	Insert(ctx context.Context, d *model.Deal) error
	Update(ctx context.Context, d *model.Deal) error
}
type ExpertRepository interface {
	Delete(ctx context.Context, id int) error
	GetListByUserId(ctx context.Context, userId int) ([]*model.Expert, error)
	GetOneById(ctx context.Context, id int) (*model.Expert, error)
	Insert(ctx context.Context, e *model.Expert) error
	Update(ctx context.Context, e *model.Expert) error
}
type MoexBondRepository interface {
	Insert(ctx context.Context, bond *model.MoexBond) error
	GetByISIN(ctx context.Context, isin string) (*model.MoexBond, error)
	GetListByIds(ctx context.Context, ids []int) ([]*model.MoexBond, error)
}
type MoexShareRepository interface {
	Insert(ctx context.Context, share *model.MoexShare) error
	GetByTicker(ctx context.Context, ticker string) (*model.MoexShare, error)
	GetListByIds(ctx context.Context, ids []int) ([]*model.MoexShare, error)
}
type MoexBondPositionRepository interface {
	GetListByPortfolioId(ctx context.Context, id int) ([]*model.Position, error)
	Get(ctx context.Context, portfolioId int, securityId int) (*model.Position, error)
	Insert(ctx context.Context, p *model.Position) error
	Update(ctx context.Context, p *model.Position) error
}
type MoexSharePositionRepository interface {
	GetListByPortfolioId(ctx context.Context, id int) ([]*model.Position, error)
	Get(ctx context.Context, portfolioId int, securityId int) (*model.Position, error)
	Insert(ctx context.Context, p *model.Position) error
	Update(ctx context.Context, p *model.Position) error
}
type OpinionRepository interface {
	Insert(ctx context.Context, o *model.Opinion) error
	GetListByUserId(ctx context.Context, userId int) ([]*model.Opinion, error)
	Update(ctx context.Context, o *model.Opinion) error
	Delete(ctx context.Context, id int) error
}
type PortfolioRepository interface {
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*model.Portfolio, error)
	GetByIdAndScan(ctx context.Context, id int, p *model.Portfolio) (*model.Portfolio, error)
	GetListByUserId(ctx context.Context, userId int) ([]*model.Portfolio, error)
	Insert(ctx context.Context, u *model.Portfolio) error
	Update(ctx context.Context, u *model.Portfolio) error
}
type UserRepository interface {
	Insert(ctx context.Context, u *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetById(ctx context.Context, id int) (*model.User, error)
}
