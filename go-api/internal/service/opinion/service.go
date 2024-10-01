package opinion

import (
	"context"

	"github.com/pttrulez/investor-go/internal/domain"
)

type ExpertRepository interface {
	Insert(ctx context.Context, expert domain.Expert) (domain.Expert, error)
	Delete(ctx context.Context, id int, userID int) error
	Update(ctx context.Context, expert domain.Expert, userID int) (domain.Expert, error)
	GetOneByID(ctx context.Context, id int, userID int) (domain.Expert, error)
	GetListByUserID(ctx context.Context, userID int) ([]domain.Expert, error)
}

type OpinionRepository interface {
	AttachToPosition(ctx context.Context, opinionID, positionID int) error
	Delete(ctx context.Context, id int, userID int) error
	GetOpinionsList(ctx context.Context, f domain.OpinionFilters, userID int) ([]domain.Opinion,
		error)
	Insert(ctx context.Context, o domain.Opinion) (domain.Opinion, error)
}

type Service struct {
	expertRepo  ExpertRepository
	opinionRepo OpinionRepository
}

func NewOpinionService(expertRepo ExpertRepository, opinionRepo OpinionRepository) *Service {
	return &Service{
		expertRepo:  expertRepo,
		opinionRepo: opinionRepo,
	}
}
