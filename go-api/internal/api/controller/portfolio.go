package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/api/model/converter"
	"github.com/pttrulez/investor-go/internal/api/model/dto"
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
	httpresponse "github.com/pttrulez/investor-go/pkg/http-response"
)

func (c *PortfolioController) CreateNewPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var dto dto.CreatePortfolio
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		httpresponse.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request fields
	if err := validator.New().Struct(dto); err != nil {
		validateErr := err.(validator.ValidationErrors)
		httpresponse.WriteValidationErrorsJSON(w, validateErr)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	portfolio := converter.FromCreatePortfolioDtoToPortfolio(&dto)
	portfolio.UserId = int(claims["id"].(float64))

	// Create new Portfolio
	err := c.services.Portfolio.CreatePortfolio(ctx, portfolio)
	if err != nil {
		httpresponse.SendError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *PortfolioController) GetListOfPortfolios(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	portfolios, err := c.services.Portfolio.GetListByUserId(ctx, getUserIdFromJwt(r))
	if err != nil {
		httpresponse.SendError(w, err)
		return
	}
	httpresponse.WriteOKJSON(w, portfolios)
}

func (c *PortfolioController) GetPortfolioById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	portfolioId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpresponse.WriteString(w, http.StatusBadRequest, "Неверный айди портфолио")
	}

	portfolio, err := c.services.Portfolio.GetPortfolioById(ctx, portfolioId, getUserIdFromJwt(r))
	if err != nil {
		httpresponse.SendError(w, err)
	}
	httpresponse.WriteOKJSON(w, portfolio)
}

func (c *PortfolioController) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpresponse.WriteString(w, http.StatusBadRequest, "Неверный айди портфолио")
	}

	err = c.services.Portfolio.DeletePortfolio(ctx, id, getUserIdFromJwt(r))
	if err != nil {
		httpresponse.WriteString(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

func (c *PortfolioController) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var dto dto.UpdatePortfolio
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		httpresponse.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update Portfolio
	err = c.services.Portfolio.UpdatePortfolio(ctx,
		converter.FromUpdatePortfolioDtoToPortfolio(&dto), getUserIdFromJwt(r))
	if err != nil {
		httpresponse.SendError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type PortfolioController struct {
	repo     *repository.Repository
	services *service.Container
}

func NewPortfolioController(r *repository.Repository, s *service.Container) *PortfolioController {
	return &PortfolioController{
		repo:     r,
		services: s,
	}
}
