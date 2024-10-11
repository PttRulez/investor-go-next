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

	"github.com/pttrulez/investor-go-next/go-api/config"
	grpcServer "github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/grpc/server"
	api "github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server"
	issclient "github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/iss-client"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage/postgres"
	"github.com/pttrulez/investor-go-next/go-api/internal/service/moex"
	"github.com/pttrulez/investor-go-next/go-api/internal/service/opinion"
	"github.com/pttrulez/investor-go-next/go-api/internal/service/portfolio"
	"github.com/pttrulez/investor-go-next/go-api/internal/service/telegram"
	"github.com/pttrulez/investor-go-next/go-api/internal/service/user"
	"github.com/pttrulez/investor-go-next/go-api/pkg/logger"
	"github.com/redis/go-redis/v9"
)

func Run() {
	// logger init
	l := logger.SetupPrettySlog()
	logger := logger.NewLogger(l)

	// config init
	cfg := config.MustLoad()

	// init repository
	connStr := fmt.Sprintf(`postgresql://%v:%v@%v:%v/%v?sslmode=%v`,
		cfg.Pg.Username, cfg.Pg.Password, cfg.Pg.Host, cfg.Pg.Port, cfg.Pg.DBName, cfg.Pg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to start postgres db: %w", err))
	}
	repo := postgres.NewRepository(db)

	// init external apis
	issClient := issclient.NewISSClient()

	// Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// init tgClient
	telega := telegram.New(cfg.TgClientEndpoint, logger)
	fmt.Println("Telegram client is ready to send to ednpoint: ", cfg.TgClientEndpoint)

	// init services
	moex := moex.NewMoexService(issClient, repo)
	opinion := opinion.NewOpinionService(repo)
	portfolio := portfolio.NewPortfolioService(issClient, logger, moex, repo, redisClient,
		telega)
	user := user.NewUserService(repo)

	// grpc Server
	grpcServ := grpcServer.NewGRPCServer(portfolio)
	go func() {
		grpcServ.MustStart(cfg.GrpcServerPort)
	}()

	// start API Server
	apiServer, err := api.StartApiServer(cfg.API, api.Services{
		OpinionService:   opinion,
		PortfolioService: portfolio,
		MoexService:      moex,
		UserService:      user,
	}, logger)
	if err != nil {
		logger.Error(fmt.Errorf("failed to start API server: %w", err))
		return
	}

	// Graceful shutdown logic:
	appCtx, stopAppCtx := context.WithCancel(context.Background())
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
