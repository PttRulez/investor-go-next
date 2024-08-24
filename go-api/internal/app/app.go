package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/lib/pq" // postgres db adapter
	"github.com/pttrulez/investor-go/config"
	"github.com/pttrulez/investor-go/internal/controller/httpcontrollers"
	"github.com/pttrulez/investor-go/internal/controller/requestvalidator"
	"github.com/pttrulez/investor-go/internal/infrastracture/issclient"
	"github.com/pttrulez/investor-go/internal/infrastracture/postgres"
	"github.com/pttrulez/investor-go/internal/logger"
	"github.com/pttrulez/investor-go/internal/service/deal"
	"github.com/pttrulez/investor-go/internal/service/expert"
	"github.com/pttrulez/investor-go/internal/service/moexbond"
	"github.com/pttrulez/investor-go/internal/service/moexshare"
	"github.com/pttrulez/investor-go/internal/service/portfolio"
	"github.com/pttrulez/investor-go/internal/service/transaction"
	"github.com/pttrulez/investor-go/internal/service/user"
)

func Run() {
	cfg := config.MustLoad()

	logger := logger.NewLogger()
	validator, err := requestvalidator.NewValidator()
	if err != nil {
		panic("Failed to create validator")
	}

	r := chi.NewRouter()
	logger.Info(fmt.Sprintf("cfg.AllowedCors: %v", cfg.AllowedCors))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedCors,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	// Repositories init
	connStr := fmt.Sprintf(`postgresql://%v:%v@%v:%v/%v?sslmode=%v`,
		cfg.Pg.Username, cfg.Pg.Password, cfg.Pg.Host, cfg.Pg.Port, cfg.Pg.DBName, cfg.Pg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error(fmt.Errorf("failed to start postgres db: %w", err))
	}

	dealRepo := postgres.NewDealPostgres(db)
	expertRepo := postgres.NewExpertPostgres(db)
	moexBondRepo := postgres.NewMoexBondsPostgres(db)
	moexShareRepo := postgres.NewMoexSharesPostgres(db)
	portfolioRepo := postgres.NewPortfolioPostgres(db)
	positionRepo := postgres.NewPositionPostgres(db)
	transactionRepo := postgres.NewTransactionPostgres(db)
	userRepo := postgres.NewUserPostgres(db)

	// Services init
	tokenAuth := jwtauth.New("HS256", []byte(cfg.TokenAuthSecret), nil)
	issClient := issclient.NewISSClient()

	expertSerice := expert.NewExpertService(expertRepo)
	moexBondService := moexbond.NewMoexBondService(moexBondRepo, issClient)
	moexShareService := moexshare.NewMoexShareService(moexShareRepo, issClient)
	portfolioService := portfolio.NewPortfolioService(dealRepo, issClient, positionRepo,
		portfolioRepo, transactionRepo)
	transactionService := transaction.NewTransactionService(transactionRepo, portfolioRepo)
	userService := user.NewUserService(userRepo, tokenAuth)
	dealService := deal.NewDealService(issClient, moexBondService, moexShareService, dealRepo)

	// Controllers init
	authController := httpcontrollers.NewAuthController(logger, userService, validator)
	dealController := httpcontrollers.NewDealController(logger, dealService, validator)
	expertController := httpcontrollers.NewExpertController(logger, expertSerice, validator)
	moexBondController := httpcontrollers.NewMoexBondController(logger, moexBondService)
	moexShareController := httpcontrollers.NewMoexShareController(logger, moexShareService)
	portfolioController := httpcontrollers.NewPortfolioController(logger, portfolioService)
	transactionController := httpcontrollers.NewCashoutController(logger, transactionService, validator)

	// Public Routes
	r.Post("/register", authController.RegisterUser)
	r.Post("/login", authController.LoginUser)

	// Protected Routes
	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		// Transactions
		r.Route("/transaction", func(r chi.Router) {
			r.Post("/", transactionController.CreateNewTransaction)
			r.Delete("/{id}", transactionController.DeleteTransaction)
		})

		//// Deals
		r.Route("/deal", func(r chi.Router) {
			r.Post("/", dealController.CreateDeal)
			r.Post("/delete/{id}", dealController.DeleteDeal)
		})

		// Experts
		r.Route("/expert", func(r chi.Router) {
			r.Post("/", expertController.CreateNewExpert)
			r.Get("/", expertController.GetExpertsList)
			r.Delete("/{id}", expertController.GetExpertsList)
		})

		// Moex-Bonds
		r.Route("/moex-bond", func(r chi.Router) {
			r.Get("/{secid}", moexBondController.GetInfoBySecid)
		})

		// Moex-Shares
		r.Route("/moex-share", func(r chi.Router) {
			r.Get("/{secid}", moexShareController.GetInfoBySecid)
		})

		// Portfolios
		r.Route("/portfolio", func(r chi.Router) {
			r.Delete("/{id}", portfolioController.DeletePortfolio)
			r.Get("/", portfolioController.GetListOfPortfoliosOfCurrentUser)
			r.Get("/{id}", portfolioController.GetPortfolioByID)
			r.Post("/", portfolioController.CreateNewPortfolio)
			r.Put("/", portfolioController.UpdatePortfolio)
		})
	})

	// The HTTP Server
	const headerTimeout = time.Second * 10
	address := fmt.Sprintf("%v:%v", cfg.APIHost, cfg.APIPort)
	srv := &http.Server{
		Addr:              address,
		Handler:           r,
		ReadHeaderTimeout: headerTimeout,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscalls
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quit

		const shutdownTimeout = 30 * time.Second
		shutdownCtx, cancel := context.WithTimeout(serverCtx, shutdownTimeout)
		defer func() { cancel() }()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit")
			}
		}()

		//	Startting Graceful shutdown

		// Shutdown database
		e := db.Close()
		if e != nil {
			log.Fatal(e)
		}

		// Shutdown server
		e = srv.Shutdown(shutdownCtx)
		if e != nil {
			log.Fatal("Server forced to shutdown: ", e)
		}

		serverStopCtx()
	}()

	// Run Server
	logger.Info(fmt.Sprintf("Listening %v", cfg.APIPort))
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
