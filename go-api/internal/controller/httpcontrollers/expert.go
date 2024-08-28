package httpcontrollers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pttrulez/investor-go/internal/controller/converter"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"github.com/pttrulez/investor-go/pkg/api"

	"github.com/go-playground/validator/v10"
)

func (c *ExpertController) CreateNewExpert(w http.ResponseWriter, r *http.Request) {
	const op = "ExpertController.CreateNewExpert"

	ctx := r.Context()
	var req api.CreateExpertRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}
	expert := converter.FromCreateExpertRequestToExpert(req)

	// Validate request fields
	if err := c.validator.Struct(expert); err != nil {
		var valErr validator.ValidationErrors
		if errors.As(err, &valErr) {
			writeValidationErrorsJSON(w, valErr)
		}
	}

	// Adding userID to expert
	userID := utils.GetCurrentUserID(r.Context())
	expert.UserID = userID

	// Save new Expert in DB
	expert, err := c.expertService.CreateNewExpert(ctx, expert)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, expert)
}

func (c *ExpertController) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "ExpertController.Delete"

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, fmt.Sprintf("Проблема с конвертацией айди %s: %s",
			chi.URLParam(r, "id"),
			err.Error()))
	}

	ctx := r.Context()
	err = c.expertService.DeleteExpert(ctx, id, utils.GetCurrentUserID(ctx))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ExpertController) GetExpertsList(w http.ResponseWriter, r *http.Request) {
	const op = "ExpertController.GetExpertsList"

	ctx := r.Context()

	// Получаем список экспертов, обрабатываем и ретерним в случае ошибки
	experts, err := c.expertService.GetListByUserID(ctx, utils.GetCurrentUserID(r.Context()))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	// Собираем полученный список экспертов в формат http response'a
	expertsResponse := make([]api.ExpertResponse, 0, len(experts))
	for _, e := range experts {
		expertsResponse = append(expertsResponse, converter.FromExpertToExpertResponse(e))
	}

	writeJSON(w, http.StatusOK, expertsResponse)
}

type ExpertService interface {
	CreateNewExpert(ctx context.Context, expert entity.Expert) (entity.Expert, error)
	DeleteExpert(ctx context.Context, id int, userID int) error
	GetListByUserID(ctx context.Context, userID int) ([]entity.Expert, error)
}

type ExpertController struct {
	logger        Logger
	expertService ExpertService
	validator     *validator.Validate
}

func NewExpertController(logger Logger, expertService ExpertService,
	validator *validator.Validate) *ExpertController {
	return &ExpertController{
		logger:        logger,
		expertService: expertService,
		validator:     validator,
	}
}
