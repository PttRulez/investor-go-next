package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/lib/response"
	"github.com/pttrulez/investor-go/internal/services"
	"github.com/pttrulez/investor-go/internal/types"
)

func (c *PortfolioController) CreateNewPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var portfolioData types.Portfolio
	if err := json.NewDecoder(r.Body).Decode(&portfolioData); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request fields
	if err := validator.New().Struct(portfolioData); err != nil {
		validateErr := err.(validator.ValidationErrors)
		response.WriteValidationErrorsJSON(w, validateErr)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	portfolioData.UserId = int(claims["id"].(float64))

	// Create new Portfolio
	err := c.repo.Portfolio.Insert(ctx, &portfolioData)
	if err != nil {
		response.WriteString(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *PortfolioController) GetListOfPortfolios(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	portfolios, err := c.repo.Portfolio.GetListByUserId(ctx, getUserIdFromJwt(r))
	if err != nil {
		response.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteOKJSON(w, portfolios)
}

func (c *PortfolioController) GetPortfolioById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	portfolioId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteString(w, http.StatusBadRequest, "Неверный айди портфолио")
	}

	portfolio, err := c.services.Portfolio.GetPortfolio(ctx, portfolioId, getUserIdFromJwt(r))
	if err != nil {
		log.Println("GetPortfolioById:\n", err)
		response.SendError(w, err)
	}
	response.WriteOKJSON(w, portfolio)
}

func (c *PortfolioController) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	err := c.repo.Portfolio.Delete(ctx, id)
	if err != nil {
		response.WriteString(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

func (c *PortfolioController) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var portfolioData types.PortfolioUpdate
	err := json.NewDecoder(r.Body).Decode(&portfolioData)
	if err != nil {
		fmt.Println("Patch portfolioData err", err)
	}

	// Update Portfolio
	err = c.repo.Portfolio.Update(ctx, &portfolioData)
	if err != nil {
		response.WriteString(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

type PortfolioController struct {
	repo     *types.Repository
	services *services.ServiceContainer
}

func NewPortfolioController(r *types.Repository, s *services.ServiceContainer) types.PortfolioController {
	return &PortfolioController{
		repo:     r,
		services: s,
	}
}
