package moex

import (
	"context"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/issclient"
)

func NewMoexService(
	bondRepo MoexBondRepository,
	shareRepo MoexShareRepository,
	issClient *issclient.IssClient,
) *Service {
	return &Service{
		issClient: issClient,
		bondRepo:  bondRepo,
		shareRepo: shareRepo,
	}
}

type Service struct {
	issClient *issclient.IssClient
	bondRepo  MoexBondRepository
	shareRepo MoexShareRepository
}

type MoexBondRepository interface {
	GetByTicker(ctx context.Context, ticker string) (domain.Bond, error)
	Insert(ctx context.Context, bond domain.Bond) (domain.Bond, error)
}

type MoexShareRepository interface {
	GetByTicker(ctx context.Context, ticker string) (domain.Share, error)
	Insert(ctx context.Context, share domain.Share) (domain.Share, error)
}
