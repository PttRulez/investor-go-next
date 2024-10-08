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
	// coupon
	CreateCoupon(ctx context.Context, d domain.Coupon, userID int) error
	DeleteCoupon(ctx context.Context, couponID int, userID int) error

	// deal
	CreateDeal(ctx context.Context, d domain.Deal, userID int) (domain.Deal, error)
	DeleteDealByID(ctx context.Context, id int, userID int) error

	// dividend
	CreateDividend(ctx context.Context, d domain.Dividend, userID int) error
	DeleteDividend(ctx context.Context, dividendID int, userID int) error

	// dividend
	CreateExpense(ctx context.Context, d domain.Expense, userID int) error
	DeleteExpense(ctx context.Context, expenseID int, userID int) error

	// portfolio
	CreatePortfolio(ctx context.Context, p domain.Portfolio) (domain.Portfolio, error)
	DeletePortfolio(ctx context.Context, portfolioID int, userID int) error
	GetFullPortfolioByID(ctx context.Context, portfolioID int, userID int) (domain.Portfolio, error)
	GetPortfolioList(ctx context.Context, userID int) ([]domain.Portfolio, error)
	UpdatePortfolio(ctx context.Context, portfolio domain.Portfolio, userID int) (domain.Portfolio, error)

	// position
	AddInfoToPosition(ctx context.Context, i domain.PositionUpdateInfo) error
	GetPositionList(ctx context.Context, userID int) ([]domain.Position, error)

	// transaction
	DeleteTransaction(ctx context.Context, transactionID int, userID int) error
	CreateTransaction(ctx context.Context, t domain.Transaction) (domain.Transaction, error)
}

type MoexService interface {
	GetBondByTicker(ctx context.Context, isin string) (domain.Bond, error)
	GetShareByTicker(ctx context.Context, ticker string) (domain.Share, error)
}

type UserService interface {
	LoginUser(ctx context.Context, user domain.User) (domain.User, error)
	RegisterUser(ctx context.Context, user domain.User) error
	UpdateUser(ctx context.Context, user domain.User) error
}
