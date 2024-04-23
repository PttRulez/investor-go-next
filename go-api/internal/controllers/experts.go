package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/lib/response"
	"github.com/pttrulez/investor-go/internal/services"
	"github.com/pttrulez/investor-go/internal/types"
)

func (c *ExpertController) CreateNewExpert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var expert types.Expert
	if err := json.NewDecoder(r.Body).Decode(&expert); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request fields
	if err := c.services.Validator.Struct(expert); err != nil {
		validateErr := err.(validator.ValidationErrors)
		response.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	// Expert must be created by user
	_, claims, _ := jwtauth.FromContext(r.Context())
	expert.UserId = int(claims["id"].(float64))

	// Save new Expert in DB
	err := c.repo.Expert.Insert(ctx, &expert)
	if err != nil {
		response.WriteString(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *ExpertController) GetExpertsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	experts, err := c.repo.Expert.GetListByUserId(ctx, getUserIdFromJwt(r))
	if err != nil {
		response.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteOKJSON(w, experts)
}

type ExpertController struct {
	repo     *types.Repository
	services *services.ServiceContainer
}

func NewExpertController(repo *types.Repository, services *services.ServiceContainer) types.ExpertController {
	return &ExpertController{
		repo:     repo,
		services: services,
	}
}
