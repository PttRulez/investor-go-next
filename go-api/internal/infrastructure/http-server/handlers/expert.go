package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/converter"
	"github.com/pttrulez/investor-go-next/go-api/internal/utils"

	"github.com/go-playground/validator/v10"
)

func (c *Handlers) CreateNewExpert(w http.ResponseWriter, r *http.Request) {
	const op = "ExpertController.CreateNewExpert"

	ctx := r.Context()
	var req contracts.CreateExpertRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, fmt.Errorf("%s: %w", op, err))
		return
	}
	expert := converter.FromCreateExpertRequestToExpert(req)

	// Validate request fields
	if err := c.validator.Struct(expert); err != nil {
		var valErr validator.ValidationErrors
		if errors.As(err, &valErr) {
			writeValidationErrorsJSON(w, valErr)
		}
	}

	// Adding userID to expert
	userID := utils.GetCurrentUserID(r.Context())
	expert.UserID = userID

	// Save new Expert in DB
	expert, err := c.opinionService.CreateNewExpert(ctx, expert)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, expert)
}

func (c *Handlers) DeleteExpert(w http.ResponseWriter, r *http.Request) {
	const op = "ExpertController.Delete"

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, fmt.Sprintf("Проблема с конвертацией айди %s: %s",
			chi.URLParam(r, "id"),
			err.Error()))
	}

	ctx := r.Context()
	err = c.opinionService.DeleteExpert(ctx, id, utils.GetCurrentUserID(ctx))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Handlers) GetExpertsList(w http.ResponseWriter, r *http.Request) {
	const op = "ExpertController.GetExpertsList"

	ctx := r.Context()

	// Получаем список экспертов, обрабатываем и ретерним в случае ошибки
	experts, err := c.opinionService.GetExpertsList(ctx, utils.GetCurrentUserID(r.Context()))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	// Собираем полученный список экспертов в формат http response'a
	expertsResponse := make([]contracts.ExpertResponse, 0, len(experts))
	for _, e := range experts {
		expertsResponse = append(expertsResponse, converter.FromExpertToExpertResponse(e))
	}

	writeJSON(w, http.StatusOK, expertsResponse)
}
