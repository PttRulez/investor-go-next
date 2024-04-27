package app

import (
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
)

func NewServiceContainer(repo *repository.Repository) *service.Container {
	container := &service.Container{}

	container.Cashout = service.NewCashoutService(repo, container)
	container.Deposit = service.NewDepositService(repo, container)
	container.Deal = service.NewDealService(container)
	container.Expert = service.NewExpertService(repo, container)
	container.Moex = service.MoexContainer{
		Api:       service.NewIssApiService(),
		Bond:      service.NewMoexBondService(repo, container),
		BondDeal:  service.NewMoexBondDealService(repo, container),
		Deal:      service.NewMoexDealService(repo, container),
		Share:     service.NewMoexShareService(repo, container),
		ShareDeal: service.NewMoexShareDealService(repo, container),
	}
	container.Portfolio = service.NewPortfolioService(repo, container)
	container.User = service.NewUserService(repo)
	container.Validator = service.NewValidator()

	return container
}
