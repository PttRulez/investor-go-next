package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pttrulez/investor-go/config"
	"github.com/pttrulez/investor-go/internal/api"
	"github.com/pttrulez/investor-go/internal/infrastructure/issclient"
	"github.com/pttrulez/investor-go/internal/infrastructure/postgres"
	"github.com/pttrulez/investor-go/internal/service/moex"
	"github.com/pttrulez/investor-go/internal/service/opinion"
	"github.com/pttrulez/investor-go/internal/service/portfolio"
	"github.com/pttrulez/investor-go/internal/service/user"
)

func Run() {
	// config init
	cfg := config.MustLoad()

	// repositories init
	connStr := fmt.Sprintf(`postgresql://%v:%v@%v:%v/%v?sslmode=%v`,
		cfg.Pg.Username, cfg.Pg.Password, cfg.Pg.Host, cfg.Pg.Port, cfg.Pg.DBName, cfg.Pg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to start postgres db: %w", err))
	}

	dealRepo := postgres.NewDealPostgres(db)
	expertRepo := postgres.NewExpertPostgres(db)
	moexBondRepo := postgres.NewMoexBondsPostgres(db)
	moexShareRepo := postgres.NewMoexSharesPostgres(db)
	opinionRepo := postgres.NewOpinionPostgres(db)
	portfolioRepo := postgres.NewPortfolioPostgres(db)
	positionRepo := postgres.NewPositionPostgres(db)
	transactionRepo := postgres.NewTransactionPostgres(db)
	userRepo := postgres.NewUserPostgres(db)

	// external apis
	issClient := issclient.NewISSClient()

	// services init
	moex := moex.NewMoexService(moexBondRepo, moexShareRepo, issClient)
	opinion := opinion.NewOpinionService(expertRepo, opinionRepo)
	portfolio := portfolio.NewPortfolioService(dealRepo, issClient, moexShareRepo, moex,
		portfolioRepo, positionRepo, transactionRepo)
	user := user.NewUserService(userRepo)

	// start API Server
	apiServer := api.StartApiServer(cfg.API, api.Services{
		OpinionService:   opinion,
		PortfolioService: portfolio,
		MoexService:      moex,
		UserService:      user,
	})

	// Server run context
	appCtx, stopAppCtx := context.WithCancel(context.Background())

	// Listen for syscalls
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quit

		const apiShutdownTimeout = 30 * time.Second
		apiShutdownCtx, cancel := context.WithTimeout(appCtx, apiShutdownTimeout)
		defer func() { cancel() }()

		// Shutdown database
		e := db.Close()
		if e != nil {
			log.Fatal(e)
		}

		// Shutdown server
		e = apiServer.Shutdown(apiShutdownCtx)
		if e != nil {
			log.Fatal("Server forced to shutdown: ", e)
		}

		stopAppCtx()
	}()

	// Wait for server context to be stopped
	<-appCtx.Done()
}
