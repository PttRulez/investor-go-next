package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
)

func Init(r *chi.Mux, repo *repository.Repository, serviceContainer *service.Container, tokenAuth *jwtauth.JWTAuth) {
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
			r.Route("/moex", func(r chi.Router) {
				r.Post("/", c.Deal.Moex.CreateNewDeal)
				r.Delete("/", c.Deal.Moex.DeleteDeal)
			})
			// r.Route("/moex-bond", func(r chi.Router) {
			// 	r.Post("/", c.Deal.MoexBond.CreateNewDeal)
			// 	r.Delete("/{id}", c.Deal.MoexBond.DeleteDeal)
			// })
			// r.Route("/moex-share", func(r chi.Router) {
			// 	r.Post("/", c.Deal.MoexShare.CreateNewDeal)
			// 	r.Delete("/{id}", c.Deal.MoexShare.DeleteDeal)
			// })
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
			r.Get("/", c.Portfolio.GetListOfPortfolios)
			r.Get("/{id}", c.Portfolio.GetPortfolioById)
			r.Post("/", c.Portfolio.CreateNewPortfolio)
			r.Put("/", c.Portfolio.UpdatePortfolio)
		})
	})
}

type Controllers struct {
	Auth      *AuthController
	Deal      *DealController
	Deposit   *DepositController
	Cashout   *CashoutController
	Expert    *ExpertController
	MoexBond  *MoexBondController
	MoexShare *MoexShareController
	Portfolio *PortfolioController
}

type DealController struct {
	Moex MoexDealController
	// MoexBond  types.MoexBondDealController
	// MoexShare types.MoexShareDealController
}

func NewControllers(repo *repository.Repository, services *service.Container, tokenAuth *jwtauth.JWTAuth) *Controllers {
	return &Controllers{
		Auth:      NewAuthController(repo, tokenAuth),
		Cashout:   NewCashoutController(repo, services),
		Deposit:   NewDepositController(repo, services),
		Expert:    NewExpertController(repo, services),
		MoexBond:  NewMoexBondController(repo, services),
		MoexShare: NewMoexShareController(repo, services),
		Portfolio: NewPortfolioController(repo, services),
	}
}
