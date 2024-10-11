package converter

import (
	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
)

func FromCreateExpenseRequestToExpense(
	req contracts.CreateExpenseRequest) (domain.Expense, error) {

	return domain.Expense{
		Amount:      req.Amount,
		Date:        req.Date.Time,
		Description: req.Description,
		PortfolioID: req.PortfolioId,
	}, nil
}
