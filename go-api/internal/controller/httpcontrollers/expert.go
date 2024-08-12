package httpcontrollers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"

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
		var valErr validator.ValidationErrors
		if errors.As(err, &valErr) {
			writeValidationErrorsJSON(w, valErr)
		}
	}

	// Expert must be created by user
	userID := utils.GetCurrentUserID(r.Context())
	expert.UserID = userID

	// Save new Expert in DB
	err := c.expertService.CreateNewExpert(ctx, &expert)
	if err != nil {
		writeString(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *ExpertController) GetExpertsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	experts, err := c.expertService.GetListByUserID(ctx, utils.GetCurrentUserID(r.Context()))
	if err != nil {
		writeString(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, experts)
}

type ExpertService interface {
	CreateNewExpert(ctx context.Context, expert *entity.Expert) error
	GetListByUserID(ctx context.Context, userID int) ([]*entity.Expert, error)
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
