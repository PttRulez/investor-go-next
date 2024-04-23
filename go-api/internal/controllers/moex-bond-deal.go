package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/lib/response"
	"github.com/pttrulez/investor-go/internal/services"
	"github.com/pttrulez/investor-go/internal/types"
)

func (c *MoexBondDealController) CreateNewDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Анмаршалим данные
	var deal types.Deal
	if err := json.NewDecoder(r.Body).Decode(&deal); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация пришедших данных
	if err := c.services.Validator.Struct(deal); err != nil {
		validateErr := err.(validator.ValidationErrors)
		response.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	// Создаем сделку
	err := c.services.Deal.MoexBond.CreateDeal(ctx, &deal, getUserIdFromJwt(r))
	if err != nil {
		fmt.Println(err)
		response.SendError(w, err)
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *MoexBondDealController) DeleteDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	err := c.repo.Deal.MoexBond.Delete(ctx, id)
	if err != nil {
		response.WriteString(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

type MoexBondDealController struct {
	repo     *types.Repository
	services *services.ServiceContainer
}

func NewMoexBondDealController(repo *types.Repository, services *services.ServiceContainer) types.MoexBondDealController {
	return &MoexBondDealController{
		repo:     repo,
		services: services,
	}
}
