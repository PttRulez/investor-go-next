package services

import (
	"github.com/pttrulez/investor-go/internal/types"
)

type PositionService struct {
	repo *types.Repository
}

func NewPositionService(repo *types.Repository) *PositionService {
	return &PositionService{repo: repo}
}
