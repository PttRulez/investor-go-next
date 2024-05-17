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

func (c *DepositController) CreateNewDeposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Анмаршалим данные
	var depositData dto.CreateDeposit
	if err := json.NewDecoder(r.Body).Decode(&depositData); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация пришедших данных
	if err := c.validator.Struct(depositData); err != nil {
		validateErr := err.(validator.ValidationErrors)
		writeValidationErrorsJSON(w, validateErr)
		return
	}

	// Создаем
	deposit := converter.FromCreateDepositDtoToDeposit(&depositData)
	deposit.UserId = utils.GetCurrentUserId(r.Context())
	err := c.depositService.CreateDeposit(ctx, deposit)
	if err != nil {
		writeError(w, err)
	}

	writeString(w, http.StatusCreated, "Депозит создан")
}

func (c *DepositController) DeleteDeposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	depositId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
	}

	err = c.depositService.DeleteDeposit(ctx, depositId)
	if err != nil {
		writeError(w, err)
	}
	writeString(w, http.StatusOK, "Депозит удалён")
}

type DepositService interface {
	CreateDeposit(ctx context.Context, depositData *entity.Deposit) error
	DeleteDeposit(ctx context.Context, depositId int) error
}

type DepositController struct {
	depositService DepositService
	validator      *validator.Validate
}

func NewDepositController(depositService DepositService, validator *validator.Validate) *DepositController {
	return &DepositController{
		depositService: depositService,
		validator:      validator,
	}
}
