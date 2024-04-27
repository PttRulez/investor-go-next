package converter

import (
	"github.com/pttrulez/investor-go/internal/api/model/dto"
	"github.com/pttrulez/investor-go/internal/model"
)

func FromCreateCashoutDtoToCashout(dto *dto.CreateCashout) *model.Cashout {
	return &model.Cashout{
		Amount:      dto.Amount,
		Date:        dto.Date,
		PortfolioId: dto.PortfolioId,
	}
}
