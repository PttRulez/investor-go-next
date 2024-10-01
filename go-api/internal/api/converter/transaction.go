package converter

import (
	"context"

	"github.com/pttrulez/investor-go/internal/api/contracts"
	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/utils"
)

func FromCreateTransactionRequestToTransaction(ctx context.Context,
	req contracts.CreateTransactionRequest) (domain.Transaction, error) {
	t, err := transactionType(req.Type)
	if err != nil {
		return domain.Transaction{}, err
	}

	userID := utils.GetCurrentUserID(ctx)

	return domain.Transaction{
		Amount:      req.Amount,
		Date:        req.Date,
		PortfolioID: req.PortfolioId,
		Type:        t,
		UserID:      userID,
	}, nil
}

func FromTransactionToResponse(c domain.Transaction) contracts.TransactionResponse {
	return contracts.TransactionResponse{
		Amount:      c.Amount,
		Date:        c.Date,
		Id:          c.ID,
		PortfolioId: c.PortfolioID,
		Type:        contracts.TransactionType(c.Type),
	}
}
