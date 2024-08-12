package converter

import (
	"github.com/pttrulez/investor-go/internal/controller/model/dto"
	"github.com/pttrulez/investor-go/internal/controller/model/response"
	"github.com/pttrulez/investor-go/internal/entity"
)

func FromCreateDtoToTransaction(dto *dto.CreateTransaction) *entity.Transaction {
	return &entity.Transaction{
		Amount:      dto.Amount,
		Date:        dto.Date,
		PortfolioID: dto.PortfolioID,
		Type:        dto.Type,
	}
}

func FromTransactionToResponse(c entity.Transaction) response.Transaction {
	return response.Transaction{
		Amount: c.Amount,
		Date:   c.Date,
		ID:     c.ID,
		Type:   c.Type,
	}
}
