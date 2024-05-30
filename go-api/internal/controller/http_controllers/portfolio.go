package http_controllers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/controller/model/converter"
	"github.com/pttrulez/investor-go/internal/controller/model/dto"
	"github.com/pttrulez/investor-go/internal/controller/model/response"
)

func (c *PortfolioController) CreateNewPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var pDto dto.CreatePortfolio
	if err := json.NewDecoder(r.Body).Decode(&pDto); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request fields
	if err := validator.New().Struct(pDto); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		writeValidationErrorsJSON(w, validateErr)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	portfolio := converter.FromCreatePortfolioDtoToPortfolio(&pDto)
	portfolio.UserId = int(claims["id"].(float64))

	// Create new Portfolio
	err := c.portfolioService.CreatePortfolio(ctx, portfolio)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func (c *PortfolioController) GetListOfPortfoliosOfCurrentUser(w http.ResponseWriter,
	r *http.Request) {

	ctx := r.Context()
	portfolios, err := c.portfolioService.GetListByUserId(ctx, utils.GetCurrentUserId(r.Context()))

	var res []*response.ShortPortfolio
	for _, portfolio := range portfolios {
		res = append(res, converter.FromPortfolioToShortPortfolio(portfolio))
	}
	if err != nil {
		writeError(w, err)
		return
	}

	writeOKJSON(w, res)
}
func (c *PortfolioController) GetPortfolioById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	portfolioId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, "Неверный айди портфолио")
	}

	portfolio, err := c.portfolioService.GetFullPortfolioById(ctx, portfolioId, utils.GetCurrentUserId(r.Context()))
	if err != nil {
		writeError(w, err)
	}

	writeOKJSON(w, portfolio)
}
func (c *PortfolioController) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, "Неверный айди портфолио")
	}

	err = c.portfolioService.DeletePortfolio(ctx, id, utils.GetCurrentUserId(r.Context()))
	if err != nil {
		writeString(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusOK)
}
func (c *PortfolioController) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var pDto dto.UpdatePortfolio
	err := json.NewDecoder(r.Body).Decode(&pDto)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update Portfolio
	err = c.portfolioService.UpdatePortfolio(ctx,
		converter.FromUpdatePortfolioDtoToPortfolio(&pDto), utils.GetCurrentUserId(r.Context()))
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type PortfolioService interface {
	DeletePortfolio(ctx context.Context, portfolioId int, userId int) error
	CreatePortfolio(ctx context.Context, p *entity.Portfolio) error
	GetFullPortfolioById(ctx context.Context, portfolioId int,
		userId int) (*response.FullPortfolio, error)
	GetListByUserId(ctx context.Context, userId int) ([]*entity.Portfolio, error)
	UpdatePortfolio(ctx context.Context, portfolio *entity.Portfolio, userId int) error
}
type PortfolioController struct {
	portfolioService PortfolioService
}

func NewPortfolioController(portfolioService PortfolioService) *PortfolioController {
	return &PortfolioController{
		portfolioService: portfolioService,
	}
}
