package http_controllers

import (
	"context"
	"encoding/json"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/controller/model/converter"
	"github.com/pttrulez/investor-go/internal/controller/model/dto"
)

func (c *CashoutController) CreateNewCashout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Анмаршалим данные
	var cashoutData dto.CreateCashout
	if err := json.NewDecoder(r.Body).Decode(&cashoutData); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация пришедших данных
	if err := c.validator.Struct(cashoutData); err != nil {
		validateErr := err.(validator.ValidationErrors)
		writeValidationErrorsJSON(w, validateErr)
		return
	}

	// Создаем
	cashout := converter.FromCreateCashoutDtoToCashout(&cashoutData)
	cashout.UserId = utils.GetCurrentUserId(r.Context())
	err := c.cashoutService.CreateCashout(ctx, cashout)
	if err != nil {
		writeError(w, err)
	}

	writeString(w, http.StatusCreated, "Кэшаут создан")
}

func (c *CashoutController) DeleteCashout(w http.ResponseWriter, r *http.Request) {
	cashoutId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
	}

	err = c.cashoutService.DeleteCashout(r.Context(), cashoutId)
	if err != nil {
		writeError(w, err)
	}
	writeString(w, http.StatusOK, "Кэшаут удалён")

}

type CashoutService interface {
	CreateCashout(ctx context.Context, cashoutData *entity.Cashout) error
	DeleteCashout(ctx context.Context, cashoutId int) error
}

type CashoutController struct {
	cashoutService CashoutService
	validator      *validator.Validate
}

func NewCashoutController(cashoutService CashoutService, validator *validator.Validate) *CashoutController {
	return &CashoutController{
		cashoutService: cashoutService,
		validator:      validator,
	}
}
