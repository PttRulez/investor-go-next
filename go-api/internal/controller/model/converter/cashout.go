package converter

import (
	"github.com/pttrulez/investor-go/internal/controller/model/dto"
	"github.com/pttrulez/investor-go/internal/entity"
)

func FromCreateCashoutDtoToCashout(dto *dto.CreateCashout) *entity.Cashout {
	return &entity.Cashout{
		Amount:      dto.Amount,
		Date:        dto.Date,
		PortfolioId: dto.PortfolioId,
	}
}
