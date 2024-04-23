package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/pttrulez/investor-go/internal/services"
	"github.com/pttrulez/investor-go/internal/types"
)

func Init(r *chi.Mux, repo *types.Repository, services *services.ServiceContainer) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	c := NewControllers(repo, services, tokenAuth)

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
			r.Route("/moex-bond", func(r chi.Router) {
				r.Post("/", c.Deal.MoexBond.CreateNewDeal)
			})
			r.Route("/moex-share", func(r chi.Router) {
				r.Post("/", c.Deal.MoexShare.CreateNewDeal)
			})
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

		// Portfolios
		r.Route("/portfolio", func(r chi.Router) {
			r.Delete("/{id}", c.Portfolio.DeletePortfolio)
			r.Get("/", c.Portfolio.GetListOfPortfolios)
			r.Get("/{id}", c.Portfolio.GetPortfolioById)
			r.Post("/", c.Portfolio.CreateNewPortfolio)
			r.Put("/", c.Portfolio.UpdatePortfolio)
		})
	})
}

type Controllers struct {
	Auth types.AuthController
	Deal struct {
		MoexBond  types.MoexBondDealController
		MoexShare types.MoexShareDealController
	}
	Deposit   types.DepositController
	Cashout   types.CashoutController
	Expert    types.ExpertController
	Portfolio types.PortfolioController
}

type DealController struct {
	MoexBond  types.MoexBondDealController
	MoexShare types.MoexShareDealController
}

func NewControllers(repo *types.Repository, services *services.ServiceContainer, tokenAuth *jwtauth.JWTAuth) *Controllers {
	return &Controllers{
		Auth:    NewAuthController(repo, tokenAuth),
		Cashout: NewCashoutController(repo, services),
		Deal: DealController{
			MoexBond:  NewMoexBondDealController(repo, services),
			MoexShare: NewMoexShareDealController(repo, services),
		},
		Deposit:   NewDepositController(repo, services),
		Expert:    NewExpertController(repo, services),
		Portfolio: NewPortfolioController(repo, services),
	}
}
