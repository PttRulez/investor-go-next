package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/handlers"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/interfaces"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/metrics"
	"github.com/pttrulez/investor-go-next/go-api/pkg/logger"
)

func StartApiServer(cfg Config, s Services, log *logger.Logger) (*http.Server, error) {

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
	mtrcs := metrics.Register()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedCors,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))
	// r.Use(mwLogger.New(logger.SetupPrettySlog()))

	// Controllers init
	h := handlers.NewHandlers(log, s.OpinionService, s.PortfolioService, tokenAuth,
		s.MoexService, s.UserService, validator)

	hw := NewWrapper(log, mtrcs)

	// Public Routes
	r.Post("/register", h.RegisterUser)
	r.Post("/login", h.LoginUser)
	r.Get("/metrics", mtrcs.Default().ServeHTTP)

	// Protected Routes
	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		// portfolio
		r.Route("/portfolio", func(r chi.Router) {
			r.Delete("/{id}", h.DeletePortfolio)
			r.Get("/", hw.makeHttpHandler(h.GetListOfPortfoliosOfCurrentUser))
			r.Get("/{id}", hw.makeHttpHandler(h.GetPortfolioByID))
			r.Post("/", h.CreateNewPortfolio)
			r.Put("/", h.UpdatePortfolio)
		})

		// coupons
		r.Route("/coupon", func(r chi.Router) {
			r.Post("/", h.CreateCoupon)
			r.Delete("/{id}", h.DeleteCoupon)
		})

		// deals
		r.Route("/deal", func(r chi.Router) {
			r.Post("/", h.CreateDeal)
			r.Delete("/{id}", hw.makeHttpHandler(h.DeleteDeal))
		})

		// dividends
		r.Route("/dividend", func(r chi.Router) {
			r.Post("/", h.CreateDividend)
			r.Delete("/{id}", h.DeleteDividend)
		})

		// expenses
		r.Route("/expense", func(r chi.Router) {
			r.Post("/", h.CreateExpense)
			r.Delete("/{id}", h.DeleteExpense)
		})

		// deposits and cashouts
		r.Route("/transaction", func(r chi.Router) {
			r.Post("/", h.CreateNewTransaction)
			r.Delete("/{id}", h.DeleteTransaction)
		})

		// positions
		r.Route("/position", func(r chi.Router) {
			r.Get("/", h.AllUserPositions)
			r.Patch("/{id}", h.UpdatePosition)
		})

		// opinion
		r.Route("/opinion", func(r chi.Router) {
			r.Delete("/{id}", h.DeleteOpinion)
			r.Get("/{opinionID}/attach-position/{positionID}", h.AttachToPosition)
			r.Post("/", h.CreateOpinion)
			r.Get("/list", h.GetOpinionsList)
		})

		// experts
		r.Route("/expert", func(r chi.Router) {
			r.Post("/", h.CreateNewExpert)
			r.Get("/", h.GetExpertsList)
			r.Delete("/{id}", h.DeleteExpert)
		})

		// moex
		r.Route("/moex", func(r chi.Router) {
			r.Get("/bond/{ticker}", h.GetBondInfoByTicker)
			r.Get("/share/{ticker}", h.GetShareInfoByTicker)
		})

		// users
		r.Route("/user", func(r chi.Router) {
			r.Put("/", h.UpdateUser)
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

	log.Info(fmt.Sprintf("starting http server on port %v", cfg.APIPort))
	// Starting server
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return nil, err
	}

	return srv, nil
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
