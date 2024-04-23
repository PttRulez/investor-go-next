package services

import (
	"github.com/pttrulez/investor-go/internal/types"
)

type MoexService struct {
	Shares MoexShareService
	Bonds  MoexBondService
}

func NewMoexService(repo *types.Repository) *MoexService {
	return &MoexService{
		Shares: *NewShareService(repo),
	}
}
