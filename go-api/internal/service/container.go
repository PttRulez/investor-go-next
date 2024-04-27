package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/model"
)

type Container struct {
	Cashout   Cashout
	Deposit   Deposit
	Deal      Deal
	Expert    Expert
	Moex      MoexContainer
	Opinion   Opinion
	Portfolio Portfolio
	User      User
	Validator *validator.Validate
}

type MoexContainer struct {
	Api       IssApi
	Bond      MoexBond
	BondDeal  MoexBondDeal
	Deal      MoexDeal
	Share     MoexShare
	ShareDeal MoexShareDeal
}

type Cashout interface {
	CreateCashout(ctx context.Context, cashoutData *model.Cashout, userId int) error
	DeleteCashout(ctx context.Context, cashoutId int, userId int) error
}

type Deposit interface {
	CreateDeposit(ctx context.Context, depositData *model.Deposit, userId int) error
	DeleteDeposit(ctx context.Context, depositId int, userId int) error
}

type Expert interface {
	CreateNewExpert(ctx context.Context, expert *model.Expert) error
	DeleteExpert(ctx context.Context, id int, userId int) error
	GetExpertsListByUserId(ctx context.Context, userId int) ([]*model.Expert, error)
	GetExpertById(ctx context.Context, id int, userId int) (*model.Expert, error)
	UpdateExpert(ctx context.Context, expert *model.Expert, userId int) error
}
type IssApi interface {
	GetSecurityInfoBySecid(secid string) (*ISSecurityInfo, error)
	GetStocksCurrentPrices(ctx context.Context, market model.ISSMoexMarket,
		tickers []string) (*model.MoexApiResponseCurrentPrices, error)
}
type MoexBond interface {
	GetByISIN(ctx context.Context, isin string) (*model.MoexBond, error)
	UpdatePositionInDB(ctx context.Context, portfolioId int,
		securityId int) error
}
type MoexShare interface {
	GetByTicker(ctx context.Context, ticker string) (*model.MoexShare, error)
	UpdatePositionInDB(ctx context.Context, portfolioId int,
		securityId int) error
}
type Deal interface {
	CreateDeal(ctx context.Context, deal *model.Deal, userId int) error
	DeleteDeal(ctx context.Context, deal *model.Deal, userId int) error
}
type MoexDeal interface {
	CreateDeal(ctx context.Context, deal *model.Deal, userId int) error
	DeleteDeal(ctx context.Context, deal *model.Deal, userId int) error
}
type MoexBondDeal interface {
	CreateDeal(ctx context.Context, deal *model.Deal, userId int) error
	DeleteDeal(ctx context.Context, dealId int, userId int) error
}
type MoexShareDeal interface {
	CreateDeal(ctx context.Context, deal *model.Deal, userId int) error
	DeleteDeal(ctx context.Context, dealId int, userId int) error
}
type Opinion interface {
	CreateNewOpinion(ctx context.Context, opinion *model.Opinion) error
	DeleteOpinion(ctx context.Context, id int, userId int) error
	GetOpinionsListByUserId(ctx context.Context, userId int) ([]*model.Opinion, error)
	GetOpinionById(ctx context.Context, id int, userId int) (*model.Opinion, error)
	UpdateOpinion(ctx context.Context, opinion *model.Opinion, userId int) error
}
type Portfolio interface {
	CreatePortfolio(ctx context.Context, p *model.Portfolio) error
	GetListByUserId(ctx context.Context, userId int) ([]*model.Portfolio, error)
	GetPortfolioById(ctx context.Context, portfolioId int, userId int) (*model.Portfolio, error)
	DeletePortfolio(ctx context.Context, portfolioId int, userId int) error
	UpdatePortfolio(ctx context.Context, portfolio *model.Portfolio, userId int) error
}

type User interface {
	LoginUser(ctx context.Context, user *model.User) (string, error)
	RegisterUser(ctx context.Context, user *model.User) error
}
