package portfolio

import (
	"context"
	"database/sql"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/iss-client"
	"github.com/pttrulez/investor-go/internal/service/moex"
)

type Repository interface {
	ExecAsTransaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error

	// coupons
	DeleteCoupon(ctx context.Context, id int, userID int) error
	GetCouponList(ctx context.Context, portfolioID int) ([]domain.Coupon, error)
	InsertCoupon(ctx context.Context, d domain.Coupon, userID int) error

	// deals
	InsertDeal(ctx context.Context, tx *sql.Tx, d domain.Deal) (domain.Deal, error)
	DeleteDeal(ctx context.Context, id int, userID int) (domain.Deal, error)
	GetDealList(ctx context.Context, portfolioID int, userID int) ([]domain.Deal, error)
	GetDealListForSecurity(ctx context.Context, exchange domain.Exchange, portfolioID int,
		securityType domain.SecurityType, ticker string) ([]domain.Deal, error)

	// dividends
	DeleteDividend(ctx context.Context, id int, userID int) error
	GetDividendList(ctx context.Context, portfolioID int) ([]domain.Dividend, error)
	InsertDividend(ctx context.Context, d domain.Dividend, userID int) error

	// expenses
	DeleteExpense(ctx context.Context, id int, userID int) error
	GetExpenseList(ctx context.Context, portfolioID int) ([]domain.Expense, error)
	InsertExpense(ctx context.Context, d domain.Expense, userID int) error

	// positions
	AddPositionInfo(ctx context.Context, i domain.PositionUpdateInfo) error
	GetPosition(ctx context.Context, exchange domain.Exchange, portfolioID int,
		securityType domain.SecurityType, ticker string) (domain.Position, error)
	GetPortfolioPositionList(ctx context.Context, portfolioID int, userID int) ([]domain.Position, error)
	GetUserPositionList(ctx context.Context, userID int) ([]domain.Position, error)
	InsertPosition(ctx context.Context, p domain.Position) error
	UpdatePosition(ctx context.Context, tx *sql.Tx, p domain.Position) error

	// moex-share
	GetMoexShare(ctx context.Context, ticker string) (domain.Share, error)

	// portfolio
	DeletePortfolio(ctx context.Context, id int, userID int) error
	GetPortfolio(ctx context.Context, id int, userID int) (domain.Portfolio, error)
	GetPortfolioList(ctx context.Context, userID int) ([]domain.Portfolio, error)
	InsertPortfolio(ctx context.Context, p domain.Portfolio) (domain.Portfolio, error)
	UpdatePortfolio(ctx context.Context, p domain.Portfolio, userID int) (domain.Portfolio, error)

	// transaction
	DeleteTransaction(ctx context.Context, id int, userID int) error
	GetTransaction(ctx context.Context, id int, userID int) (domain.Transaction, error)
	GetTransactionList(ctx context.Context, portfolioID int, userID int) ([]domain.Transaction, error)
	InsertTransaction(ctx context.Context, t domain.Transaction) (domain.Transaction, error)
}

type Service struct {
	repo        Repository
	issClient   *issclient.IssClient
	moexService *moex.Service
}

func NewPortfolioService(
	issClient *issclient.IssClient,
	moexService *moex.Service,
	repo Repository,
) *Service {
	return &Service{
		repo:        repo,
		issClient:   issClient,
		moexService: moexService,
	}
}
