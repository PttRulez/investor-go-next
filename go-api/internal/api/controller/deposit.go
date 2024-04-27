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

func (c *DepositController) CreateNewDeposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Анмаршалим данные
	var deposit dto.CreateDeposit
	if err := json.NewDecoder(r.Body).Decode(&deposit); err != nil {
		httpresponse.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация пришедших данных
	if err := c.services.Validator.Struct(deposit); err != nil {
		validateErr := err.(validator.ValidationErrors)
		httpresponse.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	// Создаем
	err := c.services.Deposit.CreateDeposit(ctx, converter.FromCreateDepositDtoToDeposit(&deposit), getUserIdFromJwt(r))
	if err != nil {
		httpresponse.SendError(w, err)
	}

	httpresponse.WriteString(w, http.StatusCreated, "Депозит создан")
}

func (c *DepositController) DeleteDeposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	depositId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpresponse.SendError(w, err)
	}

	err = c.services.Deposit.DeleteDeposit(ctx, depositId, getUserIdFromJwt(r))
	if err != nil {
		httpresponse.SendError(w, err)
	}
	httpresponse.WriteString(w, http.StatusOK, "Депозит удалён")
}

type DepositController struct {
	repo     *repository.Repository
	services *service.Container
}

func NewDepositController(repo *repository.Repository, services *service.Container) *DepositController {
	return &DepositController{
		repo:     repo,
		services: services,
	}
}
