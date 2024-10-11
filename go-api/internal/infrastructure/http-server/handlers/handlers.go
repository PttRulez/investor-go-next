package handlers

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/interfaces"
	"github.com/pttrulez/investor-go-next/go-api/pkg/logger"
)

func NewHandlers(
	logger *logger.Logger,
	opinionService interfaces.OpinionService,
	portfolioService interfaces.PortfolioService,
	tokenAuth *jwtauth.JWTAuth,
	moexService interfaces.MoexService,
	userService interfaces.UserService,
	validator *validator.Validate) *Handlers {
	return &Handlers{
		logger:           logger,
		opinionService:   opinionService,
		portfolioService: portfolioService,
		tokenAuth:        tokenAuth,
		moexService:      moexService,
		userService:      userService,
		validator:        validator,
	}
}

type Handlers struct {
	logger           *logger.Logger
	opinionService   interfaces.OpinionService
	portfolioService interfaces.PortfolioService
	tokenAuth        *jwtauth.JWTAuth
	moexService      interfaces.MoexService
	userService      interfaces.UserService
	validator        *validator.Validate
}
