package converter

import (
	"github.com/pttrulez/investor-go/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go/internal/domain"
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
