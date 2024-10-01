package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/api/contracts"
	"github.com/pttrulez/investor-go/internal/api/converter"
	"github.com/pttrulez/investor-go/internal/utils"
)

func (c *Handlers) CreateNewPortfolio(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.CreateNewPortfolio"

	ctx := r.Context()

	var req contracts.CreatePortfolioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.logger.Error(err)
		writeString(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request fields
	if err := validator.New().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		writeValidationErrorsJSON(w, validateErr)
		return
	}
	portfolio := converter.FromCreatePortfolioRequestToPortfolio(req)
	portfolio.UserID = utils.GetCurrentUserID(r.Context())

	// Create new Portfolio
	portfolio, err := c.portfolioService.CreatePortfolio(ctx, portfolio)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, portfolio)
}

func (c *Handlers) GetListOfPortfoliosOfCurrentUser(w http.ResponseWriter,
	r *http.Request) {
	const op = "Handlers.GetListOfPortfoliosOfCurrentUser"

	ctx := r.Context()
	portfolios, err := c.portfolioService.GetListByUserID(ctx, utils.GetCurrentUserID(r.Context()))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	var res []contracts.PortfolioResponse
	for _, portfolio := range portfolios {
		res = append(res, converter.FromPortfolioToPortfolioResponse(portfolio))
	}

	writeJSON(w, http.StatusOK, res)
}

func (c *Handlers) GetPortfolioByID(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.GetPortfolioByID"
	ctx := r.Context()
	portfolioID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, fmt.Sprintf("%s: Проблема с конвертацией айди %s: %s",
			op,
			chi.URLParam(r, "id"),
			err.Error()))
	}

	portfolio, err := c.portfolioService.GetFullPortfolioByID(ctx, portfolioID,
		utils.GetCurrentUserID(r.Context()))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.FromPortfolioToFullPortfolioResponse(portfolio))
}

func (c *Handlers) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.DeletePortfolio"

	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, fmt.Sprintf("Проблема с конвертацией айди %s: %s",
			chi.URLParam(r, "id"),
			err.Error()))
	}

	err = c.portfolioService.DeletePortfolio(ctx, id, utils.GetCurrentUserID(r.Context()))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Handlers) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UpdatePortfolio"

	ctx := r.Context()

	var req contracts.UpdatePortfolioRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeString(w, http.StatusBadRequest,
			fmt.Sprintf("Проблемы с парсингом JSON'а: %s", err.Error()))
		return
	}

	portfolio := converter.FromUpdatePortfolioRequestToPortfolio(req)

	// Update Portfolio
	up, err := c.portfolioService.UpdatePortfolio(ctx, portfolio,
		utils.GetCurrentUserID(r.Context()))
	if err != nil {
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.FromPortfolioToPortfolioResponse(up))
}

// type PortfolioController struct {
// 	logger           *logger.Logger
// 	portfolioService interfaces.PortfolioService
// }

// func NewPortfolioController(logger *logger.Logger,
// 	portfolioService interfaces.PortfolioService) *PortfolioController {
// 	return &PortfolioController{
// 		logger:           logger,
// 		portfolioService: portfolioService,
// 	}
// }
