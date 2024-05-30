package app

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/lib/pq"
	"github.com/pttrulez/investor-go/config"
	"github.com/pttrulez/investor-go/internal/controller/http_controllers"
	"github.com/pttrulez/investor-go/internal/controller/request_validator"
	"github.com/pttrulez/investor-go/internal/infrastracture/iss_client"
	"github.com/pttrulez/investor-go/internal/infrastracture/postgres"
	"github.com/pttrulez/investor-go/internal/service/cashout"
	"github.com/pttrulez/investor-go/internal/service/deal"
	"github.com/pttrulez/investor-go/internal/service/deposit"
	"github.com/pttrulez/investor-go/internal/service/expert"
	"github.com/pttrulez/investor-go/internal/service/moex_bond"
	"github.com/pttrulez/investor-go/internal/service/moex_share"
	"github.com/pttrulez/investor-go/internal/service/user"
	"log"
	"log/slog"
	"net/http"
)

func Run() {
	cfg := config.MustLoad()

	logger := slog.Default()
	validator := request_validator.NewValidator()
	r := chi.NewRouter()

	// Repositories init
	connStr := fmt.Sprintf(`postgresql://%v:%v@%v:%v/%v?sslmode=%v`,
		cfg.Pg.Username, cfg.Pg.Password, cfg.Pg.Host, cfg.Pg.Port, cfg.Pg.DBName, cfg.Pg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error(fmt.Sprintf("Error starting postgres db: %v", err))
	}

	cashoutRepo := postgres.NewCashoutPostgres(db)
	depositRepo := postgres.NewDepositPostgres(db)
	dealRepo := postgres.NewDealPostgres(db)
	expertRepo := postgres.NewExpertPostgres(db)
	moexBondRepo := postgres.NewMoexBondsPostgres(db)
	moexShareRepo := postgres.NewMoexSharesPostgres(db)
	portfolioRepo := postgres.NewPortfolioPostgres(db)
	userRepo := postgres.NewUserPostgres(db)

	//Services init
	tokenAuth := jwtauth.New("HS256", []byte(cfg.TokenAuthSecret), nil)
	issClient := iss_client.NewISSClient()

	cashoutService := cashout.NewCashoutService(cashoutRepo, portfolioRepo)
	depositService := deposit.NewDepositService(depositRepo, portfolioRepo)
	expertSerice := expert.NewExpertService(expertRepo)
	moexBondService := moex_bond.NewMoexBondService(moexBondRepo, issClient)
	moexShareService := moex_share.NewMoexShareService(moexShareRepo, issClient)
	userService := user.NewUserService(userRepo, tokenAuth)
	dealService := deal.NewDealService(issClient, moexBondService, moexShareService, dealRepo)

	// Controllers init
	authController := http_controllers.NewAuthController(userService, validator)
	cashoutController := http_controllers.NewCashoutController(cashoutService, validator)
	depositController := http_controllers.NewDepositController(depositService, validator)
	dealController := http_controllers.NewDealController(dealService, validator)
	expertController := http_controllers.NewExpertController(expertSerice, validator)
	moexBondController := http_controllers.NewMoexBondController(moexBondService)
	moexShareController := http_controllers.NewMoexShareController(moexShareService)

	// Public Routes
	r.Post("/register", authController.RegisterUser)
	r.Post("/login", authController.LoginUser)

	// Protected Routes
	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		// Cashouts
		r.Route("/entity", func(r chi.Router) {
			r.Post("/", cashoutController.CreateNewCashout)
			r.Delete("/{id}", cashoutController.DeleteCashout)
		})

		//// Deals
		r.Route("/deal", func(r chi.Router) {
			r.Post("/", dealController.CreateDeal)
			r.Post("/delete/{id}", dealController.DeleteDeal)
		})

		// Deposits
		r.Route("/deposit", func(r chi.Router) {
			r.Post("/", depositController.CreateNewDeposit)
			r.Delete("/{id}", depositController.DeleteDeposit)
		})

		// Experts
		r.Route("/expert", func(r chi.Router) {
			r.Post("/", expertController.CreateNewExpert)
			r.Get("/", expertController.GetExpertsList)
		})

		// Moex-Bonds
		r.Route("/moex-bond", func(r chi.Router) {
			r.Get("/{isin}", moexBondController.GetInfoByISIN)
		})

		// Moex-Shares
		r.Route("/moex-share", func(r chi.Router) {
			r.Get("/{ticker}", moexShareController.GetInfoByTicker)
		})
		//
		//// Portfolios
		//r.Route("/portfolio", func(r chi.Router) {
		//	r.Delete("/{id}", c.Portfolio.DeletePortfolio)
		//	r.Get("/", c.Portfolio.GetListOfPortfoliosOfCurrentUser)
		//	r.Get("/{id}", c.Portfolio.GetPortfolioById)
		//	r.Post("/", c.Portfolio.CreateNewPortfolio)
		//	r.Put("/", c.Portfolio.UpdatePortfolio)
		//})
	})

	address := fmt.Sprintf("%v:%v", cfg.ApiHost, cfg.ApiPort)
	srv := &http.Server{
		Addr:    address,
		Handler: r,
	}
	logger.Info(fmt.Sprintf("Listening on  %v", address))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}

}
