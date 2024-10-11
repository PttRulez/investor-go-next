package portfolio

import (
	"context"
	"database/sql"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	issclient "github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/iss-client"
	"github.com/pttrulez/investor-go-next/go-api/internal/service/moex"
	"github.com/pttrulez/investor-go-next/go-api/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Transac(ctx context.Context, opts *sql.TxOptions, fn func(ctx context.Context) error) (err error)

	// coupons
	DeleteCoupon(ctx context.Context, id int, userID int) error
	GetCouponList(ctx context.Context, portfolioID int) ([]domain.Coupon, error)
	InsertCoupon(ctx context.Context, d domain.Coupon, userID int) error

	// deals
	InsertDeal(ctx context.Context, d domain.Deal) (domain.Deal, error)
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
	UpdatePosition(ctx context.Context, p domain.Position) error

	// moex-share
	GetMoexShare(ctx context.Context, ticker string) (domain.Share, error)

	// portfolio
	DeletePortfolio(ctx context.Context, id int, userID int) error
	GetPortfolio(ctx context.Context, id int, userID int) (domain.Portfolio, error)
	GetPortfolioList(ctx context.Context, userID int) ([]domain.Portfolio, error)
	GetPortfolioListByChatID(ctx context.Context, chatId string) ([]domain.Portfolio, error)
	InsertPortfolio(ctx context.Context, p domain.Portfolio) (domain.Portfolio, error)
	UpdatePortfolio(ctx context.Context, p domain.Portfolio, userID int) (domain.Portfolio, error)

	// transaction
	DeleteTransaction(ctx context.Context, id int, userID int) error
	GetTransaction(ctx context.Context, id int, userID int) (domain.Transaction, error)
	GetTransactionList(ctx context.Context, portfolioID int, userID int) ([]domain.Transaction, error)
	InsertTransaction(ctx context.Context, t domain.Transaction) (domain.Transaction, error)
}

type Telegram interface {
	SendMsg(ctx context.Context, text string)
}

type Service struct {
	issClient   *issclient.IssClient
	logger      *logger.Logger
	moexService *moex.Service
	redisClient *redis.Client
	repo        Repository
	tg          Telegram
}

func NewPortfolioService(
	issClient *issclient.IssClient,
	logger *logger.Logger,
	moexService *moex.Service,
	repo Repository,
	redisClient *redis.Client,
	tg Telegram,
) *Service {
	return &Service{
		repo:        repo,
		issClient:   issClient,
		logger:      logger,
		moexService: moexService,
		redisClient: redisClient,
		tg:          tg,
	}
}
