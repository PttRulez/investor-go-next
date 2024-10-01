package interfaces

import (
	"context"

	"github.com/pttrulez/investor-go/internal/domain"
)

type OpinionService interface {
	AttachToPosition(ctx context.Context, opinionID, positionID int) error
	CreateNewExpert(ctx context.Context, expert domain.Expert) (domain.Expert, error)
	CreateOpinion(ctx context.Context, opinion domain.Opinion) (domain.Opinion, error)
	DeleteExpert(ctx context.Context, id int, userID int) error
	DeleteOpinionByID(ctx context.Context, id int, userID int) error
	GetExpertsList(ctx context.Context, userID int) ([]domain.Expert, error)
	GetOpinionsList(ctx context.Context, f domain.OpinionFilters, userID int) ([]domain.Opinion,
		error)
}

type PortfolioService interface {
	AddInfoToPosition(ctx context.Context, i domain.PositionUpdateInfo) error
	DeletePortfolio(ctx context.Context, portfolioID int, userID int) error
	DeleteDealByID(ctx context.Context, id int, userID int) error
	DeleteTransaction(ctx context.Context, transactionID int, userID int) error
	CreateDeal(ctx context.Context, d domain.Deal) (domain.Deal, error)
	CreatePortfolio(ctx context.Context, p domain.Portfolio) (domain.Portfolio, error)
	CreateTransaction(ctx context.Context, t domain.Transaction) (domain.Transaction, error)
	GetFullPortfolioByID(ctx context.Context, portfolioID int,
		userID int) (domain.Portfolio, error)
	GetListByUserID(ctx context.Context, userID int) ([]domain.Portfolio, error)
	GetPositionList(ctx context.Context, userID int) ([]domain.Position, error)
	UpdatePortfolio(ctx context.Context, portfolio domain.Portfolio, userID int) (
		domain.Portfolio, error)
}

type MoexService interface {
	GetBondByTicker(ctx context.Context, isin string) (domain.Bond, error)
	GetShareByTicker(ctx context.Context, ticker string) (domain.Share, error)
}

type UserService interface {
	LoginUser(ctx context.Context, user domain.User) (domain.User, error)
	RegisterUser(ctx context.Context, user domain.User) error
}
