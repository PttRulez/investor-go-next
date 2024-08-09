package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/lib/pq"
	"github.com/pttrulez/investor-go/config"
	"github.com/pttrulez/investor-go/internal/controller/http_controllers"
	"github.com/pttrulez/investor-go/internal/controller/request_validator"
	"github.com/pttrulez/investor-go/internal/infrastracture/iss_client"
	"github.com/pttrulez/investor-go/internal/infrastracture/postgres"
	"github.com/pttrulez/investor-go/internal/service/deal"
	"github.com/pttrulez/investor-go/internal/service/expert"
	"github.com/pttrulez/investor-go/internal/service/moex_bond"
	"github.com/pttrulez/investor-go/internal/service/moex_share"
	"github.com/pttrulez/investor-go/internal/service/portfolio"
	"github.com/pttrulez/investor-go/internal/service/transaction"
	"github.com/pttrulez/investor-go/internal/service/user"
)

func Run() {
	cfg := config.MustLoad()

	logger := slog.Default()
	validator := request_validator.NewValidator()
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedCors,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Repositories init
	connStr := fmt.Sprintf(`postgresql://%v:%v@%v:%v/%v?sslmode=%v`,
		cfg.Pg.Username, cfg.Pg.Password, cfg.Pg.Host, cfg.Pg.Port, cfg.Pg.DBName, cfg.Pg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error(fmt.Sprintf("Error starting postgres db: %v", err))
	}

	dealRepo := postgres.NewDealPostgres(db)
	expertRepo := postgres.NewExpertPostgres(db)
	moexBondRepo := postgres.NewMoexBondsPostgres(db)
	moexShareRepo := postgres.NewMoexSharesPostgres(db)
	portfolioRepo := postgres.NewPortfolioPostgres(db)
	positionRepo := postgres.NewPositionPostgres(db)
	transactionRepo := postgres.NewTransactionPostgres(db)
	userRepo := postgres.NewUserPostgres(db)

	//Services init
	tokenAuth := jwtauth.New("HS256", []byte(cfg.TokenAuthSecret), nil)
	issClient := iss_client.NewISSClient()

	expertSerice := expert.NewExpertService(expertRepo)
	moexBondService := moex_bond.NewMoexBondService(moexBondRepo, issClient)
	moexShareService := moex_share.NewMoexShareService(moexShareRepo, issClient)
	portfolioService := portfolio.NewPortfolioService(dealRepo, issClient, positionRepo, portfolioRepo, transactionRepo)
	transactionService := transaction.NewTransactionService(transactionRepo, portfolioRepo)
	userService := user.NewUserService(userRepo, tokenAuth)
	dealService := deal.NewDealService(issClient, moexBondService, moexShareService, dealRepo)

	// Controllers init
	authController := http_controllers.NewAuthController(userService, validator)
	dealController := http_controllers.NewDealController(dealService, validator)
	expertController := http_controllers.NewExpertController(expertSerice, validator)
	moexBondController := http_controllers.NewMoexBondController(moexBondService)
	moexShareController := http_controllers.NewMoexShareController(moexShareService)
	portfolioController := http_controllers.NewPortfolioController(portfolioService)
	transactionController := http_controllers.NewCashoutController(transactionService, validator)

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
			r.Get("/{id}", portfolioController.GetPortfolioById)
			r.Post("/", portfolioController.CreateNewPortfolio)
			r.Put("/", portfolioController.UpdatePortfolio)
		})
	})

	// The HTTP Server
	address := fmt.Sprintf("%v:%v", cfg.ApiHost, cfg.ApiPort)
	srv := &http.Server{
		Addr:    address,
		Handler: r,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscalls
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quit

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit")
			}
		}()

		//	Startting Graceful shutdown

		// Shutdown database
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}

		// Shutdown server
		err = srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal("Server forced to shutdown: ", err)
		}

		serverStopCtx()
	}()

	// Run Server
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
