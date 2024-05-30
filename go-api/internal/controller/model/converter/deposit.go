package converter

import (
	"github.com/pttrulez/investor-go/internal/controller/model/dto"
	"github.com/pttrulez/investor-go/internal/controller/model/response"
	"github.com/pttrulez/investor-go/internal/entity"
)

func FromCreateDepositDtoToDeposit(dto *dto.CreateDeposit) *entity.Deposit {
	return &entity.Deposit{
		Amount:      dto.Amount,
		Date:        dto.Date,
		PortfolioId: dto.PortfolioId,
	}
}

func FromDepositToResponse(c entity.Deposit) response.Deposit {
	return response.Deposit{
		Amount: c.Amount,
		Date:   c.Date,
	}
}
