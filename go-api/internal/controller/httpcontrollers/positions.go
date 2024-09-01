package httpcontrollers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/controller/converter"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"github.com/pttrulez/investor-go/pkg/api"
)

func (c *PositionsController) AllUserPositions(w http.ResponseWriter, r *http.Request) {
	const op = "PositionsController.AllUsersPositions"

	ctx := r.Context()

	positions, err := c.positionService.GetListByUserID(ctx, utils.GetCurrentUserID(ctx))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	res := make([]api.PositionResponse, 0, len(positions))
	for _, p := range positions {
		res = append(res, converter.FromPositionToResponse(p))
	}
	writeJSON(w, http.StatusOK, res)
}

func (c *PositionsController) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	const op = "PositionsController.UpdatePosition"

	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, fmt.Sprintf("Проблема с конвертацией айди %s: %s",
			chi.URLParam(r, "id"),
			err.Error()))
	}

	var req api.UpdatePositionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	err = c.positionService.AddInfo(ctx, entity.PositionUpdateInfo{
		ID:          id,
		Comment:     req.Comment,
		TargetPrice: req.TargetPrice,
		UserID:      utils.GetCurrentUserID(ctx),
	})
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

type PostionService interface {
	GetListByUserID(ctx context.Context, userID int) ([]entity.Position, error)
	AddInfo(ctx context.Context, i entity.PositionUpdateInfo) error
}

type PositionsController struct {
	logger          Logger
	positionService PostionService
	validator       *validator.Validate
}

func NewPositionsController(logger Logger, positionService PostionService,
	validator *validator.Validate) *PositionsController {
	return &PositionsController{
		logger:          logger,
		positionService: positionService,
		validator:       validator,
	}
}
