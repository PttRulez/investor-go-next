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

func (c *CashoutController) CreateNewCashout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Анмаршалим данные
	var cashout types.Cashout
	if err := json.NewDecoder(r.Body).Decode(&cashout); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация пришедших данных
	if err := c.services.Validator.Struct(cashout); err != nil {
		validateErr := err.(validator.ValidationErrors)
		response.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	// Создаем
	err := c.services.Cashout.CreateCashout(ctx, &cashout, getUserIdFromJwt(r))
	if err != nil {
		response.WriteString(w, http.StatusInternalServerError, "Что-то пошло не так")
	}

	response.WriteString(w, http.StatusCreated, "Кэшаут создан")
}

func (c *CashoutController) DeleteCashout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cashoutId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteString(w, http.StatusBadRequest, types.ErrWrongId.Error())
	}

	err = c.services.Cashout.DeleteCashout(ctx, cashoutId, getUserIdFromJwt(r))
	if err != nil {
		response.SendError(w, err)
	}
	response.WriteString(w, http.StatusOK, "Кэшаут удалён")
}

type CashoutController struct {
	repo     *types.Repository
	services *services.ServiceContainer
}

func NewCashoutController(repo *types.Repository, services *services.ServiceContainer) types.CashoutController {
	return &CashoutController{
		repo:     repo,
		services: services,
	}
}
