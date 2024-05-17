package app

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/pttrulez/investor-go/internal/controller/request_validator"
	"github.com/pttrulez/investor-go/internal/infrastracture/iss_client"
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
	"github.com/pttrulez/investor-go/internal/service/cashout"
	"github.com/pttrulez/investor-go/internal/service/deal"
	"github.com/pttrulez/investor-go/internal/service/deposit"
	"github.com/pttrulez/investor-go/internal/service/expert"
	"github.com/pttrulez/investor-go/internal/service/moex_bond"
	"github.com/pttrulez/investor-go/internal/service/moex_bond_deal"
	"github.com/pttrulez/investor-go/internal/service/moex_deal"
	"github.com/pttrulez/investor-go/internal/service/moex_share"
	"github.com/pttrulez/investor-go/internal/service/moex_share_deal"
	"github.com/pttrulez/investor-go/internal/service/portfolio"
	"github.com/pttrulez/investor-go/internal/service/user"
)

func NewServiceContainer(repo *repository.Repository, tokenAuth *jwtauth.JWTAuth) *service.Container {
	container := &service.Container{}

	container.Cashout = cashout.NewCashoutService(repo, container)
	container.Deposit = deposit.NewDepositService(repo, container)
	container.Deal = deal.NewDealService(container)
	container.Expert = expert.NewExpertService(repo, container)
	container.Moex = service.MoexContainer{
		Api:       iss_client.NewIssApiService(),
		Bond:      moex_bond.NewMoexBondService(repo, container),
		BondDeal:  moex_bond_deal.NewMoexBondDealService(repo, container),
		Deal:      moex_deal.NewMoexDealService(repo, container),
		Share:     moex_share.NewMoexShareService(repo, container),
		ShareDeal: moex_share_deal.NewMoexShareDealService(repo, container),
	}
	container.Portfolio = portfolio.NewPortfolioService(repo, container)
	container.User = user.NewUserService(repo, tokenAuth)
	container.Validator = request_validator.NewValidator()

	return container
}
