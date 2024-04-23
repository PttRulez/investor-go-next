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
	MoexBond  *MoexBondService
	MoexShare *MoexShareService
	Portfolio *PortfolioService
	Validator *validator.Validate
}

type DealService struct {
	MoexShare *MoexShareDealService
	MoexBond  *MoexBondDealService
}

func NewServiceContainer(repo *types.Repository) *ServiceContainer {
	container := &ServiceContainer{}
	container.Deal = &DealService{
		MoexShare: NewMoexShareDealService(repo, container),
		MoexBond:  NewMoexBondDealService(repo, container),
	}
	container.Cashout = NewCashoutService(repo, container)
	container.Deposit = NewDepositService(repo, container)
	container.IssApi = NewIssApiService()
	container.MoexBond = NewMoexBondService(repo)
	container.MoexShare = NewMoexShareService(repo)
	container.Portfolio = NewPortfolioService(repo, container)
	container.Validator = NewValidator()

	return container
}
