package opinion

import (
	"context"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
)

type Repository interface {
	// experts
	DeleteExpert(ctx context.Context, id int, userID int) error
	InsertExpert(ctx context.Context, expert domain.Expert) (domain.Expert, error)
	GetExpert(ctx context.Context, id int, userID int) (domain.Expert, error)
	GetExpertList(ctx context.Context, userID int) ([]domain.Expert, error)
	UpdateExpert(ctx context.Context, expert domain.Expert, userID int) (domain.Expert, error)

	// opininons
	AttachOpinionToPosition(ctx context.Context, opinionID, positionID int) error
	DeleteOpinion(ctx context.Context, id int, userID int) error
	GetOpinionsList(ctx context.Context, f domain.OpinionFilters, userID int) ([]domain.Opinion,
		error)
	InsertOpinion(ctx context.Context, o domain.Opinion) (domain.Opinion, error)
}

type OpinionRepository interface {
}

type Service struct {
	repo Repository
}

func NewOpinionService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
