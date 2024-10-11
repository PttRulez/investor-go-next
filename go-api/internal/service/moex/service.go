package moex

import (
	"context"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/iss-client"
)

func NewMoexService(
	issClient *issclient.IssClient,
	repo Repository,
) *Service {
	return &Service{
		issClient: issClient,
		repo:      repo,
	}
}

type Service struct {
	issClient *issclient.IssClient
	repo      Repository
}

type Repository interface {
	GetMoexBond(ctx context.Context, ticker string) (domain.Bond, error)
	GetMoexShare(ctx context.Context, ticker string) (domain.Share, error)
	InsertMoexBond(ctx context.Context, bond domain.Bond) (domain.Bond, error)
	InsertMoexShare(ctx context.Context, share domain.Share) (domain.Share, error)
}
