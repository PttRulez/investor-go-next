package app

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/pttrulez/investor-go/config"
	"github.com/pttrulez/investor-go/internal/controller/http_controllers"
	"github.com/pttrulez/investor-go/internal/controller/request_validator"
	"github.com/pttrulez/investor-go/internal/infrastracture/postgres"
	"github.com/pttrulez/investor-go/internal/service/cashout"
	"github.com/pttrulez/investor-go/internal/service/deposit"
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
	portfolioRepo := postgres.NewPortfolioPostgres(db)
	userRepo := postgres.NewUserPostgres(db)

	//Services init
	tokenAuth := jwtauth.New("HS256", []byte(cfg.TokenAuthSecret), nil)

	cashoutService := cashout.NewCashoutService(cashoutRepo, portfolioRepo)
	depositService := deposit.NewDepositService(depositRepo, portfolioRepo)
	userService := user.NewUserService(userRepo, tokenAuth)

	// Controllers init
	authController := http_controllers.NewAuthController(userService, validator)
	cashoutController := http_controllers.NewCashoutController(cashoutService, validator)
	depositController := http_controllers.NewDepositController(depositService, validator)

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
		//r.Route("/deal", func(r chi.Router) {
		//	r.Post("/", c.Deal.CreateDeal)
		//	r.Post("/delete", c.Deal.DeleteDeal)
		//})

		// Deposits
		r.Route("/deposit", func(r chi.Router) {
			r.Post("/", depositController.CreateNewDeposit)
			r.Delete("/{id}", depositController.DeleteDeposit)
		})

		//// Experts
		//r.Route("/expert", func(r chi.Router) {
		//	r.Post("/", c.Expert.CreateNewExpert)
		//	r.Get("/", c.Expert.GetExpertsList)
		//})
		//
		//// Moex-Bonds
		//r.Route("/moex-share", func(r chi.Router) {
		//	r.Get("/{ticker}", c.MoexShare.GetInfoByTicker)
		//})
		//
		//// Moex-Shares
		//r.Route("/moex_bond", func(r chi.Router) {
		//	r.Get("/{isin}", c.MoexBond.GetInfoByISIN)
		//})
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
