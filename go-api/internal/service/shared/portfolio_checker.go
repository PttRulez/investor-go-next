package shared

import "context"

type PortfolioChecker interface {
	BelongsToUser(ctx context.Context, portfolioId int, userId int) (bool, error)
}
