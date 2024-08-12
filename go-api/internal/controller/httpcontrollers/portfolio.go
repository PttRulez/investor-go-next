package httpcontrollers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"

	"github.com/go-chi/chi/v5"
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
	portfolio := converter.FromCreatePortfolioDtoToPortfolio(&pDto)
	portfolio.UserID = utils.GetCurrentUserID(r.Context())

	// Create new Portfolio
	err := c.portfolioService.CreatePortfolio(ctx, portfolio)
	if err != nil {
		c.logger.Error(err)
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *PortfolioController) GetListOfPortfoliosOfCurrentUser(w http.ResponseWriter,
	r *http.Request) {
	ctx := r.Context()
	portfolios, err := c.portfolioService.GetListByUserID(ctx, utils.GetCurrentUserID(r.Context()))
	if err != nil {
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	var res []*response.ShortPortfolio
	for _, portfolio := range portfolios {
		res = append(res, converter.FromPortfolioToShortPortfolio(portfolio))
	}

	writeJSON(w, http.StatusOK, res)
}

func (c *PortfolioController) GetPortfolioByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	portfolioID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, "Неверный айди портфолио")
	}

	portfolio, err := c.portfolioService.GetFullPortfolioByID(ctx, portfolioID, utils.GetCurrentUserID(r.Context()))
	if err != nil {
		c.logger.Error(err)

		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, portfolio)
}

func (c *PortfolioController) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, "Неверный айди портфолио")
	}

	err = c.portfolioService.DeletePortfolio(ctx, id, utils.GetCurrentUserID(r.Context()))
	if err != nil {
		c.logger.Error(err)
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
		converter.FromUpdatePortfolioDtoToPortfolio(&pDto), utils.GetCurrentUserID(r.Context()))
	if err != nil {
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type PortfolioService interface {
	DeletePortfolio(ctx context.Context, portfolioID int, userID int) error
	CreatePortfolio(ctx context.Context, p *entity.Portfolio) error
	GetFullPortfolioByID(ctx context.Context, portfolioID int,
		userID int) (*response.FullPortfolio, error)
	GetListByUserID(ctx context.Context, userID int) ([]*entity.Portfolio, error)
	UpdatePortfolio(ctx context.Context, portfolio *entity.Portfolio, userID int) error
}
type PortfolioController struct {
	logger           Logger
	portfolioService PortfolioService
}

func NewPortfolioController(portfolioService PortfolioService) *PortfolioController {
	return &PortfolioController{
		portfolioService: portfolioService,
	}
}
