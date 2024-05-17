package http_controllers

import (
	"encoding/json"
	"fmt"
	"github.com/pttrulez/investor-go/internal/utils/http_response"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/controller/model/converter"
	"github.com/pttrulez/investor-go/internal/controller/model/dto"
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
)

func (c *DealController) CreateDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Анмаршалим данные
	var dto dto.CreateDeal
	var err error
	if err = json.NewDecoder(r.Body).Decode(&dto); err != nil {
		httpresponse.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	// Валидация пришедших данных
	if err = c.services.Validator.Struct(dto); err != nil {
		validateErr := err.(validator.ValidationErrors)
		httpresponse.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	err = c.services.Deal.CreateDeal(ctx,
		converter.FromCreateDealDtoToDeal(&dto), getUserIdFromJwt(r))
	if err != nil {
		httpresponse.SendError(w, err)
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *DealController) DeleteDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto dto.DeleteDeal
	var err error
	if err = json.NewDecoder(r.Body).Decode(&dto); err != nil {
		httpresponse.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	// Валидация пришедших данных
	if err = c.services.Validator.Struct(dto); err != nil {
		validateErr := err.(validator.ValidationErrors)
		httpresponse.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	userId := getUserIdFromJwt(r)
	err = c.services.Deal.DeleteDeal(ctx, converter.FromDeleteDealDtoToDeal(&dto), userId)
	if err != nil {
		fmt.Printf("[MoexShareDealController.DeleteDeal] error: %v\n", err)
		httpresponse.SendError(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

type DealController struct {
	repo     *repository.Repository
	services *service.Container
}

func NewDealController(repo *repository.Repository, services *service.Container) *DealController {
	return &DealController{
		repo:     repo,
		services: services,
	}
}
