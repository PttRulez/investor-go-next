package service

import (
	"context"
	"github.com/pttrulez/investor-go/internal/entity"
)

//	type Repository struct {
//		Cashout   CashoutRepository
//		Deal      DealRepository
//		Deposit   DepositRepository
//		Expert    ExpertRepository
//		Moex      MoexRepository
//		Portfolio PortfolioRepository
//		Position  PositionRepository
//		User      UserRepository
//	}

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

type DepositRepository interface {
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*entity.Deposit, error)
	GetListByPortfolioId(ctx context.Context, id int) ([]*entity.Deposit, error)
	Insert(ctx context.Context, c *entity.Deposit) error
}
type MoexBondDealRepository interface {
	Delete(ctx context.Context, id int) error
	GetDealListByBondId(ctx context.Context, portfolioId int,
		securityId int) ([]*entity.Deal, error)
	GetDealListByPortoflioId(ctx context.Context, portfolioId int) ([]*entity.Deal, error)
	Insert(ctx context.Context, d *entity.Deal) error
	Update(ctx context.Context, d *entity.Deal) error
}
type MoexShareDealRepository interface {
	Delete(ctx context.Context, id int) error
	GetDealListByShareId(ctx context.Context, portfolioId int,
		securityId int) ([]*entity.Deal, error)
	GetDealListByPortoflioId(ctx context.Context, portfolioId int) ([]*entity.Deal, error)
	Insert(ctx context.Context, d *entity.Deal) error
	Update(ctx context.Context, d *entity.Deal) error
}
type ExpertRepository interface {
	Delete(ctx context.Context, id int) error
	GetListByUserId(ctx context.Context, userId int) ([]*entity.Expert, error)
	GetOneById(ctx context.Context, id int) (*entity.Expert, error)
	Insert(ctx context.Context, e *entity.Expert) error
	Update(ctx context.Context, e *entity.Expert) error
}
type MoexBondRepository interface {
	Insert(ctx context.Context, bond *entity.Bond) error
	GetByISIN(ctx context.Context, isin string) (*entity.Bond, error)
	GetListByIds(ctx context.Context, ids []int) ([]*entity.Bond, error)
}
type MoexShareRepository interface {
	Insert(ctx context.Context, share *entity.Share) error
	GetByTicker(ctx context.Context, ticker string) (*entity.Share, error)
	GetListByIds(ctx context.Context, ids []int) ([]*entity.Share, error)
}
type MoexBondPositionRepository interface {
	GetListByPortfolioId(ctx context.Context, id int) ([]*entity.Position, error)
	Get(ctx context.Context, portfolioId int, securityId int) (*entity.Position, error)
	Insert(ctx context.Context, p *entity.Position) error
	Update(ctx context.Context, p *entity.Position) error
}
type MoexSharePositionRepository interface {
	GetListByPortfolioId(ctx context.Context, id int) ([]*entity.Position, error)
	Get(ctx context.Context, portfolioId int, securityId int) (*entity.Position, error)
	Insert(ctx context.Context, p *entity.Position) error
	Update(ctx context.Context, p *entity.Position) error
}
type OpinionRepository interface {
	Insert(ctx context.Context, o *entity.Opinion) error
	GetListByUserId(ctx context.Context, userId int) ([]*entity.Opinion, error)
	Update(ctx context.Context, o *entity.Opinion) error
	Delete(ctx context.Context, id int) error
}
type PortfolioRepository interface {
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*entity.Portfolio, error)
	GetListByUserId(ctx context.Context, userId int) ([]*entity.Portfolio, error)
	Insert(ctx context.Context, u *entity.Portfolio) error
	Update(ctx context.Context, u *entity.Portfolio) error
}
type UserRepository interface {
	Insert(ctx context.Context, u *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetById(ctx context.Context, id int) (*entity.User, error)
}
