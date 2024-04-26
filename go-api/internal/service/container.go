package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/model"
)

type Container struct {
	Cashout       Cashout
	Deal          Deal
	Deposit       Deposit
	Expert        Expert
	MoexApi       IssApi
	MoexBond      MoexBond
	MoexShare     MoexShare
	MoexBondDeal  MoexBondDeal
	MoexShareDeal MoexShareDeal
	Position      Position
	User          User
	Validator     *validator.Validate
}

type Cashout interface {
	CreateCashout(ctx context.Context, cashoutData *model.Cashout, userId int) error
	DeleteCashout(ctx context.Context, cashoutId int, userId int) error
}

type Deal interface {
	CreateDeal(ctx context.Context, dealData *model.Deal, userId int) error
	DeleteDeal(ctx context.Context, dealId int, userId int) error
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
	UpdateExpert(ctx context.Context, expert *model.Expert) error
}

type IssApi interface {
	GetSecurityInfoBySecid(secid string) (*ISSecurityInfo, error)
	GetStocksCurrentPrices(ctx context.Context, market model.ISSMoexMarket,
		tickers []string) (*model.MoexApiResponseCurrentPrices, error)
}

type MoexBond interface {
	GetByTicker(ctx context.Context, ticker string) (*model.MoexBond, error)
	UpdatePositionInDB(ctx context.Context, portfolioId int,
		securityId int) error
}

type MoexShare interface {
	GetByTicker(ctx context.Context, ticker string) (*model.MoexShare, error)
	UpdatePositionInDB(ctx context.Context, portfolioId int,
		securityId int) error
}
type MoexDeal interface {
	CreateDeal(ctx context.Context, deal *model.Deal, userId int) error
	DeleteDeal(ctx context.Context, dealId int, userId int) error
}
type MoexBondDeal interface {
	CreateDeal(ctx context.Context, deal *model.Deal, userId int) error
	DeleteDeal(ctx context.Context, dealId int, userId int) error
}

type MoexShareDeal interface {
	CreateDeal(ctx context.Context, deal *model.Deal, userId int) error
	DeleteDeal(ctx context.Context, dealId int, userId int) error
}

type Position interface {
}

type User interface {
	LoginUser(ctx context.Context, user *model.User) (string, error)
	RegisterUser(ctx context.Context, user *model.User) error
}
