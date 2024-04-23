package types

import (
	"context"

	tmoex "github.com/pttrulez/investor-go/internal/types/moex"
)

type Repository struct {
	Cashout           CashoutRepository
	Deal              DealRepository
	Deposit           DepositRepository
	Expert            ExpertRepository
	Moex              MoexRepository
	MoexBondPosition  MoexBondPositionRepository
	MoexSharePosition MoexSharePositionRepository
	Portfolio         PortfolioRepository
	User              UserRepository
}

type CashoutRepository interface {
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*Cashout, error)
	GetListByPortfolioId(ctx context.Context, id int) ([]*Cashout, error)
	Insert(ctx context.Context, c *Cashout) error
}
type DepositRepository interface {
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*Deposit, error)
	GetListByPortfolioId(ctx context.Context, id int) ([]*Deposit, error)
	Insert(ctx context.Context, c *Deposit) error
}
type MoexBondDealRepository interface {
	Delete(ctx context.Context, id int) error
	GetDealsListForSecurity(ctx context.Context, portfolioId int,
		securityId int) ([]*Deal, error)
	GetDealsListByPortoflioId(ctx context.Context, portfolioId int) ([]*Deal, error)
	Insert(ctx context.Context, d *Deal) error
	Update(ctx context.Context, d *Deal) error
}
type MoexShareDealRepository interface {
	Delete(ctx context.Context, id int) error
	GetDealsListForSecurity(ctx context.Context, portfolioId int,
		securityId int) ([]*Deal, error)
	GetDealsListByPortoflioId(ctx context.Context, portfolioId int) ([]*Deal, error)
	Insert(ctx context.Context, d *Deal) error
	Update(ctx context.Context, d *Deal) error
}
type ExpertRepository interface {
	Delete(ctx context.Context, id int) error
	GetListByUserId(ctx context.Context, userId int) ([]*Expert, error)
	Insert(ctx context.Context, e *Expert) error
	Update(ctx context.Context, e *Expert) error
}
type MoexRepository struct {
	Bonds  MoexBondRepository
	Shares MoexShareRepository
}
type DealRepository struct {
	MoexBonds  MoexBondDealRepository
	MoexShares MoexShareDealRepository
}
type MoexBondRepository interface {
	Insert(ctx context.Context, bond *tmoex.Bond) error
	GetByTicker(ctx context.Context, ticker string) (*tmoex.Bond, error)
	GetListByIds(ctx context.Context, ids []int) ([]*tmoex.Bond, error)
}
type MoexShareRepository interface {
	Insert(ctx context.Context, share *tmoex.Share) error
	GetByTicker(ctx context.Context, ticker string) (*tmoex.Share, error)
	GetListByIds(ctx context.Context, ids []int) ([]*tmoex.Share, error)
}
type MoexBondPositionRepository interface {
	GetListByPortfolioId(ctx context.Context, id int) ([]*BondPosition, error)
	Get(ctx context.Context, portfolioId int, securityId int) (*BondPosition, error)
	Insert(ctx context.Context, p *BondPosition) error
	Update(ctx context.Context, p *BondPosition) error
}
type MoexSharePositionRepository interface {
	GetListByPortfolioId(ctx context.Context, id int) ([]*SharePosition, error)
	Get(ctx context.Context, portfolioId int, securityId int) (*SharePosition, error)
	Insert(ctx context.Context, p *SharePosition) error
	Update(ctx context.Context, p *SharePosition) error
}
type OpinionRepository interface {
	Insert(ctx context.Context, o *Opinion) error
	GetListByUserId(ctx context.Context, userId int) ([]*Opinion, error)
	Update(ctx context.Context, o *Opinion) error
	Delete(ctx context.Context, id int) error
}
type PortfolioRepository interface {
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*Portfolio, error)
	GetByIdAndScan(ctx context.Context, id int, p *FullPortfolioData) error
	// GetByIdTx(id int, tx *sql.Tx, p *FullPortfolioData) error
	GetListByUserId(ctx context.Context, userId int) ([]*Portfolio, error)
	Insert(ctx context.Context, u *Portfolio) error
	Update(ctx context.Context, u *PortfolioUpdate) error
}
type UserRepository interface {
	Insert(ctx context.Context, u User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetById(ctx context.Context, id int) (*User, error)
}
