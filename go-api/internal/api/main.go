package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/pttrulez/investor-go/internal/api/contracts"
	"github.com/pttrulez/investor-go/internal/api/handlers"
	"github.com/pttrulez/investor-go/internal/api/interfaces"
	mwLogger "github.com/pttrulez/investor-go/internal/api/middleware/logger"
	"github.com/pttrulez/investor-go/internal/logger"
	"github.com/pttrulez/investor-go/internal/logger/slogpretty"
)

func StartApiServer(cfg Config, s Services) *http.Server {
	// logger init
	log := setupPrettySlog()
	logger := logger.NewLogger(log)

	// Validator init
	validator, err := contracts.NewValidator()
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		//nolint:mnd // Функция взята из примера от разработчика пакета validate
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	if err != nil {
		panic("Failed to create validator")
	}

	tokenAuth := jwtauth.New("HS256", []byte(cfg.TokenAuthSecret), nil)

	// Router init
	r := chi.NewRouter()
	
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedCors,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))
	r.Use(mwLogger.New(log))

	// Controllers init
	h := handlers.NewHandlers(logger, s.OpinionService, s.PortfolioService, tokenAuth,
		s.MoexService, s.UserService, validator)

	// Public Routes
	r.Post("/register", h.RegisterUser)
	r.Post("/login", h.LoginUser)

	// Protected Routes
	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		// Transactions
		r.Route("/transaction", func(r chi.Router) {
			r.Post("/", h.CreateNewTransaction)
			r.Delete("/{id}", h.DeleteTransaction)
		})

		//// Deals
		r.Route("/deal", func(r chi.Router) {
			r.Post("/", h.CreateDeal)
			r.Post("/delete/{id}", h.DeleteDeal)
		})

		// Experts
		r.Route("/expert", func(r chi.Router) {
			r.Post("/", h.CreateNewExpert)
			r.Get("/", h.GetExpertsList)
			r.Delete("/{id}", h.GetExpertsList)
		})

		// Moex-Bonds
		r.Route("/moex-bond", func(r chi.Router) {
			r.Get("/{ticker}", h.GetBondInfoByTicker)
		})

		// Moex-Shares
		r.Route("/moex-share", func(r chi.Router) {
			r.Get("/{ticker}", h.GetShareInfoByTicker)
		})

		// Moex-Shares
		r.Route("/opinion", func(r chi.Router) {
			r.Delete("/{id}", h.DeleteOpinion)
			r.Get("/{opinionID}/attach-position/{positionID}", h.AttachToPosition)
			r.Post("/", h.CreateOpinion)
			r.Get("/list", h.GetOpinionsList)
		})

		// Portfolios
		r.Route("/portfolio", func(r chi.Router) {
			r.Delete("/{id}", h.DeletePortfolio)
			r.Get("/", h.GetListOfPortfoliosOfCurrentUser)
			r.Get("/{id}", h.GetPortfolioByID)
			r.Post("/", h.CreateNewPortfolio)
			r.Put("/", h.UpdatePortfolio)
		})

		// Positions
		r.Route("/position", func(r chi.Router) {
			r.Get("/", h.AllUserPositions)
			r.Patch("/{id}", h.UpdatePosition)
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

	// Starting server
	logger.Info(fmt.Sprintf("Listening %v", cfg.APIPort))
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err)
	}

	return srv
}

type Config struct {
	APIHost         string
	APIPort         int
	AllowedCors     []string
	BaseCtx         context.Context
	TokenAuthSecret string
}

type Services struct {
	OpinionService   interfaces.OpinionService
	PortfolioService interfaces.PortfolioService
	MoexService      interfaces.MoexService
	UserService      interfaces.UserService
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
