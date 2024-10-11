package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pttrulez/investor-go-next/go-api/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/converter"
)

func (c *Handlers) CreateNewTransaction(w http.ResponseWriter, r *http.Request) {
	const op = "CreateNewTransaction.CreateNewTransaction"
	ctx := r.Context()

	// Анмаршалим данные
	var tData contracts.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&tData); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация пришедших данных
	if err := c.validator.Struct(tData); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		writeValidationErrorsJSON(w, validateErr)
		return
	}

	// Конвертируем реквест в модель
	t, err := converter.FromCreateTransactionRequestToTransaction(ctx, tData)
	if err != nil {
		writeString(w, http.StatusBadRequest, err.Error())
		return
	}

	// Получаем айди юзера и создаем транзакцию
	t.UserID = utils.GetCurrentUserID(ctx)
	tr, err := c.portfolioService.CreateTransaction(ctx, t)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, converter.FromTransactionToResponse(tr))
}

func (c *Handlers) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.DeleteTransaction"

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeString(w, http.StatusBadRequest, fmt.Sprintf("Проблема с конвертацией айди %s: %s",
			chi.URLParam(r, "id"),
			err.Error()))
		return
	}

	ctx := r.Context()
	err = c.portfolioService.DeleteTransaction(ctx, id, utils.GetCurrentUserID(ctx))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeString(w, http.StatusOK, "Транзакция удалена")
}
