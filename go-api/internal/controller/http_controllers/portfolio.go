package http_controllers

import (
	"context"
	"encoding/json"
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
	var dto dto.CreatePortfolio
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request fields
	if err := validator.New().Struct(dto); err != nil {
		validateErr := err.(validator.ValidationErrors)
		WriteValidationErrorsJSON(w, validateErr)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	portfolio := converter.FromCreatePortfolioDtoToPortfolio(&dto)
	portfolio.UserId = int(claims["id"].(float64))

	// Create new Portfolio
	err := c.portfolioService.CreatePortfolio(ctx, portfolio)
	if err != nil {
		SendError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func (c *PortfolioController) GetListOfPortfoliosOfCurrentUser(w http.ResponseWriter,
	r *http.Request) {

	ctx := r.Context()
	portfolios, err := c.portfolioService.GetListByUserId(ctx, getUserIdFromJwt(r))

	var res []*response.ShortPortfolio
	for _, portfolio := range portfolios {
		res = append(res, converter.FromPortfolioToShortPortfolio(portfolio))
	}
	if err != nil {
		SendError(w, err)
		return
	}

	WriteOKJSON(w, res)
}
func (c *PortfolioController) GetPortfolioById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	portfolioId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		WriteString(w, http.StatusBadRequest, "Неверный айди портфолио")
	}

	portfolio, err := c.portfolioService.GetFullPortfolioById(ctx, portfolioId, getUserIdFromJwt(r))
	if err != nil {
		SendError(w, err)
	}

	WriteOKJSON(w, portfolio)
}
func (c *PortfolioController) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		WriteString(w, http.StatusBadRequest, "Неверный айди портфолио")
	}

	err = c.portfolioService.DeletePortfolio(ctx, id, getUserIdFromJwt(r))
	if err != nil {
		WriteString(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusOK)
}
func (c *PortfolioController) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto dto.UpdatePortfolio
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update Portfolio
	err = c.portfolioService.UpdatePortfolio(ctx,
		converter.FromUpdatePortfolioDtoToPortfolio(&dto), getUserIdFromJwt(r))
	if err != nil {
		SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type PortfolioService interface {
	BelongsToUser(ctx context.Context, portfolioId int, userId int) (bool, error)
	// CreatePortfolio(ctx context.Context, p *service.Portfolio) error
	// GetFullPortfolioById(ctx context.Context, portfolioId int, userId int) (*service.Portfolio, error)
	// GetListByUserId(ctx context.Context, userId int) ([]*service.Portfolio, error)
	// GetPortfolioById(ctx context.Context, portfolioId int, userId int) (*service.Portfolio, error)
	// DeletePortfolio(ctx context.Context, portfolioId int, userId int) error
	// UpdatePortfolio(ctx context.Context, portfolio *service.Portfolio, userId int) error
}
type PortfolioController struct {
	portfolioService PortfolioService
}

func NewPortfolioController(portfolioService PortfolioService) *PortfolioController {
	return &PortfolioController{
		portfolioService: portfolioService,
	}
}
