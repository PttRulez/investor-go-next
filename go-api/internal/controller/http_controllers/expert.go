package http_controllers

import (
	"context"
	"encoding/json"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func (c *ExpertController) CreateNewExpert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var expert entity.Expert
	if err := json.NewDecoder(r.Body).Decode(&expert); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request fields
	if err := c.validator.Struct(expert); err != nil {
		validateErr := err.(validator.ValidationErrors)
		writeValidationErrorsJSON(w, validateErr)
		return
	}

	// Expert must be created by user
	userId := utils.GetCurrentUserId(r.Context())
	expert.UserId = userId

	// Save new Expert in DB
	err := c.expertService.CreateNewExpert(ctx, &expert)
	if err != nil {
		writeString(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *ExpertController) GetExpertsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	experts, err := c.expertService.GetListByUserId(ctx, utils.GetCurrentUserId(r.Context()))
	if err != nil {
		writeString(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOKJSON(w, experts)
}

type ExpertService interface {
	CreateNewExpert(ctx context.Context, expert *entity.Expert) error
	GetListByUserId(ctx context.Context, userId int) ([]*entity.Expert, error)
}

type ExpertController struct {
	expertService ExpertService
	validator     *validator.Validate
}

func NewExpertController(expertService ExpertService, validator *validator.Validate) *ExpertController {
	return &ExpertController{
		expertService: expertService,
		validator:     validator,
	}
}
