package converter

import (
	"github.com/pttrulez/investor-go/internal/controller/model/response"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/pkg/api"
)

func FromCreateTransactionRequestToTransaction(req api.CreateTransactionRequest) (entity.Transaction, error) {
	transactionType, err := transactionType(req.Type)
	if err != nil {
		return entity.Transaction{}, err
	}

	return entity.Transaction{
		Amount:      req.Amount,
		Date:        req.Date.Time,
		PortfolioID: req.PortfolioId,
		Type:        transactionType,
	}, nil
}

func FromTransactionToResponse(c entity.Transaction) response.Transaction {
	return response.Transaction{
		Amount: c.Amount,
		Date:   c.Date,
		ID:     c.ID,
		Type:   c.Type,
	}
}
