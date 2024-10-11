package converter

import (
	"context"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go-next/go-api/internal/utils"
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
		Date:        req.Date.Time,
		PortfolioID: req.PortfolioId,
		Type:        t,
		UserID:      userID,
	}, nil
}

func FromTransactionToResponse(c domain.Transaction) contracts.TransactionResponse {
	return contracts.TransactionResponse{
		Amount:      c.Amount,
		Date:        openapi_types.Date{Time: c.Date},
		Id:          c.ID,
		PortfolioId: c.PortfolioID,
		Type:        contracts.TransactionType(c.Type),
	}
}
