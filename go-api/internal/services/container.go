package services

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/types"
)

type ServiceContainer struct {
	db        *sql.DB
	Cashout   *CashoutService
	Deal      *DealService
	Deposit   *DepositService
	IssApi    *IssApiService
	Moex      *MoexService
	Portfolio *PortfolioService
	Validator *validator.Validate
}

type DealService struct {
	MoexShares *MoexShareDealService
	MoexBonds  *MoexBondDealService
}

func NewServiceContainer(repo *types.Repository) *ServiceContainer {
	container := &ServiceContainer{}
	container.Deal = &DealService{
		MoexShares: NewMoexShareDealService(repo, container),
		MoexBonds:  NewMoexBondDealService(repo, container),
	}
	container.Cashout = NewCashoutService(repo, container)
	container.Deposit = NewDepositService(repo, container)
	container.IssApi = NewIssApiService()
	container.Portfolio = NewPortfolioService(repo, container)
	container.Moex = NewMoexService(repo)
	container.Validator = NewValidator()

	return container
}
