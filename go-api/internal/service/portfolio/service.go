package portfolio

import (
	"context"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/issclient"
	"github.com/pttrulez/investor-go/internal/service/moex"
)

type DealRepository interface {
	Insert(ctx context.Context, d domain.Deal) (domain.Deal, error)
	GetDealListForSecurity(ctx context.Context, exchange domain.Exchange, portfolioID int,
		securityType domain.SecurityType, ticker string) ([]domain.Deal, error)
	GetDealListByPortoflioID(ctx context.Context, portfolioID int, userID int) ([]domain.Deal, error)
	Delete(ctx context.Context, id int, userID int) (domain.Deal, error)
}

type MoexShareRepo interface {
	GetByTicker(ctx context.Context, ticker string) (domain.Share, error)
}

type PositionRepository interface {
	AddInfo(ctx context.Context, i domain.PositionUpdateInfo) error
	GetListByPortfolioID(ctx context.Context, portfolioID int, userID int) ([]domain.Position, error)
	GetListByUserID(ctx context.Context, userID int) ([]domain.Position, error)
	GetPositionForSecurity(ctx context.Context, exchange domain.Exchange, portfolioID int,
		securityType domain.SecurityType, ticker string) (domain.Position, error)
	Insert(ctx context.Context, p domain.Position) error
	Update(ctx context.Context, p domain.Position) error
}

type PortfolioRepository interface {
	Delete(ctx context.Context, id int, userID int) error
	GetByID(ctx context.Context, id int, userID int) (domain.Portfolio, error)
	GetListByUserID(ctx context.Context, id int) ([]domain.Portfolio, error)
	Insert(ctx context.Context, p domain.Portfolio) (domain.Portfolio, error)
	Update(ctx context.Context, p domain.Portfolio, userID int) (domain.Portfolio, error)
}

type TransactionRepository interface {
	Delete(ctx context.Context, id int, userID int) error
	GetByID(ctx context.Context, id int, userID int) (domain.Transaction, error)
	GetListByPortfolioID(ctx context.Context, portfolioID int, userID int) ([]domain.Transaction, error)
	Insert(ctx context.Context, t domain.Transaction) (domain.Transaction, error)
}

type Service struct {
	dealRepo        DealRepository
	issClient       *issclient.IssClient
	moexShareRepo   MoexShareRepo
	moexService     *moex.Service
	portfolioRepo   PortfolioRepository
	positionRepo    PositionRepository
	transactionRepo TransactionRepository
}

func NewPortfolioService(
	dealRepo DealRepository,
	issClient *issclient.IssClient,
	moexShareRepo MoexShareRepo,
	moexService *moex.Service,
	portfolioRepo PortfolioRepository,
	positionRepo PositionRepository,
	transactionRepo TransactionRepository,
) *Service {
	return &Service{
		dealRepo:        dealRepo,
		issClient:       issClient,
		moexShareRepo:   moexShareRepo,
		moexService:     moexService,
		portfolioRepo:   portfolioRepo,
		positionRepo:    positionRepo,
		transactionRepo: transactionRepo,
	}
}
