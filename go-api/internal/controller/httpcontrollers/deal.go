package httpcontrollers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"github.com/pttrulez/investor-go/pkg/api"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/controller/model/converter"
)

func (c *DealController) CreateDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Анмаршалим данные
	var dealReq api.CreateDealRequest
	var err error
	if err = json.NewDecoder(r.Body).Decode(&dealReq); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	// Валидация пришедших данных
	if err = c.validator.Struct(dealReq); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		writeValidationErrorsJSON(w, validateErr)
		return
	}

	deal, err := converter.FromCreateDealRequestToDeal(dealReq)
	if err != nil {
		writeString(w, http.StatusBadRequest, err.Error())
	}

	err = c.dealService.CreateDeal(ctx, &deal)
	if err != nil {
		writeError(w, err)
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *DealController) DeleteDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	i := chi.URLParam(r, "id")
	id, err := strconv.Atoi(i)
	if err != nil {
		writeString(w, http.StatusBadRequest, fmt.Sprintf("Проблема с конвертацией айди %s: %s",
			chi.URLParam(r, "id"),
			err.Error()))
	}

	userID := utils.GetCurrentUserID(ctx)

	err = c.dealService.DeleteDealByID(ctx, id, userID)
	if err != nil {
		writeError(w, fmt.Errorf("[DealController.DeleteDeal] %w", err))
	}

	w.WriteHeader(http.StatusOK)
}

type DealService interface {
	CreateDeal(ctx context.Context, d *entity.Deal) error
	DeleteDealByID(ctx context.Context, id int, userID int) error
}

type DealController struct {
	dealService DealService
	validator   *validator.Validate
}

func NewDealController(dealService DealService, validator *validator.Validate) *DealController {
	return &DealController{
		dealService: dealService,
		validator:   validator,
	}
}
