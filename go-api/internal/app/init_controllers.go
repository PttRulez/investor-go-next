package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/pttrulez/investor-go/internal/api"
	"github.com/pttrulez/investor-go/internal/api/controller"
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
)

func InitControllers(r *chi.Mux, repo *repository.Repository, serviceContainer *service.Container,
	tokenAuth *jwtauth.JWTAuth) {
	c := NewControllers(repo, serviceContainer, tokenAuth)

	//                                   Public routes
	r.Post("/register", c.Auth.RegisterUser)
	r.Post("/login", c.Auth.LoginUser)

	//                               		Protected routes
	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		// Cashouts
		r.Route("/cashout", func(r chi.Router) {
			r.Post("/", c.Cashout.CreateNewCashout)
			r.Delete("/{id}", c.Cashout.DeleteCashout)
		})

		// Deals
		r.Route("/deal", func(r chi.Router) {
			r.Post("/", c.Deal.CreateDeal)
			r.Post("/delete", c.Deal.DeleteDeal)
		})

		// Deposits
		r.Route("/deposit", func(r chi.Router) {
			r.Post("/", c.Deposit.CreateNewDeposit)
			r.Delete("/{id}", c.Deposit.DeleteDeposit)
		})

		// Experts
		r.Route("/expert", func(r chi.Router) {
			r.Post("/", c.Expert.CreateNewExpert)
			r.Get("/", c.Expert.GetExpertsList)
		})

		// Moex-Bonds
		r.Route("/moex-share", func(r chi.Router) {
			r.Get("/{ticker}", c.MoexShare.GetInfoByTicker)
		})

		// Moex-Shares
		r.Route("/moex-bond", func(r chi.Router) {
			r.Get("/{isin}", c.MoexBond.GetInfoByISIN)
		})

		// Portfolios
		r.Route("/portfolio", func(r chi.Router) {
			r.Delete("/{id}", c.Portfolio.DeletePortfolio)
			r.Get("/", c.Portfolio.GetListOfPortfoliosOfCurrentUser)
			r.Get("/{id}", c.Portfolio.GetPortfolioById)
			r.Post("/", c.Portfolio.CreateNewPortfolio)
			r.Put("/", c.Portfolio.UpdatePortfolio)
		})
	})
}

type Controllers struct {
	Auth      api.AuthController
	Deposit   api.DepositController
	Cashout   api.CashoutController
	Deal      api.DealController
	Expert    api.ExpertController
	MoexBond  api.MoexBondController
	MoexShare api.MoexShareController
	Portfolio api.PortfolioController
}

func NewControllers(repo *repository.Repository, services *service.Container, tokenAuth *jwtauth.JWTAuth) *Controllers {
	return &Controllers{
		Auth:      controller.NewAuthController(repo, services),
		Cashout:   controller.NewCashoutController(repo, services),
		Deal:      controller.NewDealController(repo, services),
		Deposit:   controller.NewDepositController(repo, services),
		Expert:    controller.NewExpertController(repo, services),
		MoexBond:  controller.NewMoexBondController(repo, services),
		MoexShare: controller.NewMoexShareController(repo, services),
		Portfolio: controller.NewPortfolioController(repo, services),
	}
}
