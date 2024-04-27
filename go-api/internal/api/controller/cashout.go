package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/api/model/dto"
	"github.com/pttrulez/investor-go/internal/api/model/converter"
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
	httpresponse "github.com/pttrulez/investor-go/pkg/http-response"
)

func (c *CashoutController) CreateNewCashout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Анмаршалим данные
	var cashout dto.CreateCashout
	if err := json.NewDecoder(r.Body).Decode(&cashout); err != nil {
		httpresponse.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация пришедших данных
	if err := c.services.Validator.Struct(cashout); err != nil {
		validateErr := err.(validator.ValidationErrors)
		httpresponse.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	// Создаем
	err := c.services.Cashout.CreateCashout(ctx, converter.FromCreateCashoutDtoToCashout(&cashout), getUserIdFromJwt(r))
	if err != nil {
		httpresponse.SendError(w, err)
	}

	httpresponse.WriteString(w, http.StatusCreated, "Кэшаут создан")
}

func (c *CashoutController) DeleteCashout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cashoutId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpresponse.SendError(w, err)
	}

	err = c.services.Cashout.DeleteCashout(ctx, cashoutId, getUserIdFromJwt(r))
	if err != nil {
		httpresponse.SendError(w, err)
	}
	httpresponse.WriteString(w, http.StatusOK, "Кэшаут удалён")
}

type CashoutController struct {
	repo     *repository.Repository
	services *service.Container
}

func NewCashoutController(repo *repository.Repository, services *service.Container) *CashoutController {
	return &CashoutController{
		repo:     repo,
		services: services,
	}
}
