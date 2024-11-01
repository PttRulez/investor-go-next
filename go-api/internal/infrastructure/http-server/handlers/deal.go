package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/converter"
)

func (c *Handlers) CreateDeal(w http.ResponseWriter, r *http.Request) {
	const op = "DealController.CreateDeal"

	ctx := r.Context()

	// Анмаршалим данные
	var dealReq contracts.CreateDealRequest
	var err error
	if err = json.NewDecoder(r.Body).Decode(&dealReq); err != nil {
		c.logger.Error(err)
		writeString(w, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация пришедших данных
	if err = c.validator.Struct(dealReq); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		writeValidationErrorsJSON(w, validateErr)
		return
	}

	deal, err := converter.FromCreateDealRequestToDeal(ctx, dealReq)
	if err != nil {
		c.logger.Error(err)
		writeString(w, http.StatusBadRequest, err.Error())
	}

	if deal.SecurityType == domain.STBond && deal.Nkd == nil {
		writeJS(w, http.StatusBadRequest, map[string]string{"nkd": "Необходимо указать НКД"})
		return
	}

	result, err := c.portfolioService.CreateDeal(ctx, deal, utils.GetCurrentUserID(ctx))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (c *Handlers) DeleteDeal(w http.ResponseWriter, r *http.Request) error {
	const op = "DealController.DeleteDeal"

	ctx := r.Context()

	i := chi.URLParam(r, "id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return writeS(w, http.StatusBadRequest, fmt.Sprintf(
			"Проблема с конвертацией айди %s: %s",
			chi.URLParam(r, "id"),
			err.Error()))
	}

	userID := utils.GetCurrentUserID(ctx)

	err = c.portfolioService.DeleteDealByID(ctx, id, userID)
	if err != nil {
		return writeErr(w, fmt.Errorf("%s, %w", op, err))
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
