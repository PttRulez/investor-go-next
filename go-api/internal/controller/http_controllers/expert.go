package http_controllers

import (
	"encoding/json"
	"github.com/pttrulez/investor-go/internal/entity"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
)

func (c *ExpertController) CreateNewExpert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var expert entity.Expert
	if err := json.NewDecoder(r.Body).Decode(&expert); err != nil {
		httpresponse.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request fields
	if err := c.services.Validator.Struct(expert); err != nil {
		validateErr := err.(validator.ValidationErrors)
		httpresponse.WriteValidationErrorsJSON(w, validateErr)
		return
	}

	// Expert must be created by user
	userId := getUserIdFromJwt(r)
	expert.UserId = userId

	// Save new Expert in DB
	err := c.repo.Expert.Insert(ctx, &expert)
	if err != nil {
		httpresponse.WriteString(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *ExpertController) GetExpertsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	experts, err := c.repo.Expert.GetListByUserId(ctx, getUserIdFromJwt(r))
	if err != nil {
		httpresponse.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteOKJSON(w, experts)
}

type ExpertController struct {
	repo     *repository.Repository
	services *service.Container
}

func NewExpertController(repo *repository.Repository, services *service.Container) *ExpertController {
	return &ExpertController{
		repo:     repo,
		services: services,
	}
}
