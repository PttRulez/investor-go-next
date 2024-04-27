package converter

import (
	"github.com/pttrulez/investor-go/internal/api/model/dto"
	"github.com/pttrulez/investor-go/internal/model"
)

func FromCreateDepositDtoToDeposit(dto *dto.CreateDeposit) *model.Deposit {
	return &model.Deposit{
		Amount:      dto.Amount,
		Date:        dto.Date,
		PortfolioId: dto.PortfolioId,
	}
}
