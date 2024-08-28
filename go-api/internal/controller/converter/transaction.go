package converter

import (
	"context"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"github.com/pttrulez/investor-go/pkg/api"
)

func FromCreateTransactionRequestToTransaction(ctx context.Context,
	req api.CreateTransactionRequest) (entity.Transaction, error) {
	t, err := transactionType(req.Type)
	if err != nil {
		return entity.Transaction{}, err
	}

	userID := utils.GetCurrentUserID(ctx)

	return entity.Transaction{
		Amount:      req.Amount,
		Date:        req.Date,
		PortfolioID: req.PortfolioId,
		Type:        t,
		UserID:      userID,
	}, nil
}

func FromTransactionToResponse(c entity.Transaction) api.TransactionResponse {
	return api.TransactionResponse{
		Amount:      c.Amount,
		Date:        c.Date,
		Id:          c.ID,
		PortfolioId: c.PortfolioID,
		Type:        api.TransactionType(c.Type),
	}
}
