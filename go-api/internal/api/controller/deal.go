package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/dtos"
	"github.com/pttrulez/investor-go/internal/lib/response"
	"github.com/pttrulez/investor-go/internal/services"
	"github.com/pttrulez/investor-go/internal/types"
	tmoex "github.com/pttrulez/investor-go/internal/types/moex"
)

func (c *MoexDealController) CreateNewDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Анмаршалим данные
	var dto dtos.CreateMoexDeal
	var err error
	if err = json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	// Валидация пришедших данных
	if err = c.services.Validator.Struct(dto); err != nil {
		validateErr := err.(validator.ValidationErrors)
		response.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	c.services.Mo

	// Создаем сделку
	if dto.Market == tmoex.Market_Bonds {
		shareDeal := dto.ToMoexShareDeal()
		err = c.services.Deal.MoexShare.CreateDeal(ctx, &deal, getUserIdFromJwt(r))
	} else {
		err = c.services.Deal.MoexShare.CreateDeal(ctx, &deal, getUserIdFromJwt(r))
	}
	if err != nil {
		fmt.Println(err)
		response.SendError(w, err)
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *MoexDealController) DeleteDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	err := c.services.Deal.MoexShare.DeleteDeal(ctx, id, getUserIdFromJwt(r))
	if err != nil {
		fmt.Printf("[MoexShareDealController.DeleteDeal] error: %v\n", err)
		response.SendError(w, err)
	}
	w.WriteHeader(http.StatusOK)
}

type MoexDealController struct {
	repo     *types.Repository
	services *services.ServiceContainer
}

func NewMoexDealController(repo *types.Repository,
	services *services.ServiceContainer) *MoexDealController {
	return &MoexDealController{
		repo:     repo,
		services: services,
	}
}
