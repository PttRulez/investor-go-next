package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/lib/response"
	"github.com/pttrulez/investor-go/internal/services"
	"github.com/pttrulez/investor-go/internal/types"
)

func (c *DepositController) CreateNewDeposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Анмаршалим данные
	var deposit types.Deposit
	if err := json.NewDecoder(r.Body).Decode(&deposit); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация пришедших данных
	if err := c.services.Validator.Struct(deposit); err != nil {
		validateErr := err.(validator.ValidationErrors)
		response.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	// Создаем
	err := c.services.Deposit.CreateDeposit(ctx, &deposit, getUserIdFromJwt(r))
	if err != nil {
		response.WriteString(w, http.StatusInternalServerError, "Что-то пошло не так")
	}

	response.WriteString(w, http.StatusCreated, "Депозит создан")
}

func (c *DepositController) DeleteDeposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	depositId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteString(w, http.StatusBadRequest, types.ErrWrongId.Error())
	}

	err = c.services.Deposit.DeleteDeposit(ctx, depositId, getUserIdFromJwt(r))
	if err != nil {
		response.SendError(w, err)
	}
	response.WriteString(w, http.StatusOK, "Депозит удалён")
}

type DepositController struct {
	repo     *types.Repository
	services *services.ServiceContainer
}

func NewDepositController(repo *types.Repository, services *services.ServiceContainer) types.DepositController {
	return &DepositController{
		repo:     repo,
		services: services,
	}
}
