package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pttrulez/investor-go/internal/api/contracts"
	"github.com/pttrulez/investor-go/internal/api/converter"
	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/utils"
)

func (c *Handlers) AllUserPositions(w http.ResponseWriter, r *http.Request) {
	const op = "PositionsController.AllUsersPositions"

	ctx := r.Context()

	positions, err := c.portfolioService.GetPositionList(ctx, utils.GetCurrentUserID(ctx))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	res := make([]contracts.PositionResponse, 0, len(positions))
	for _, p := range positions {
		res = append(res, converter.FromPositionToResponse(p))
	}
	writeJSON(w, http.StatusOK, res)
}

func (c *Handlers) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	const op = "PositionsController.UpdatePosition"

	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, fmt.Sprintf("Проблема с конвертацией айди %s: %s",
			chi.URLParam(r, "id"),
			err.Error()))
	}

	var req contracts.UpdatePositionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	err = c.portfolioService.AddInfoToPosition(ctx, domain.PositionUpdateInfo{
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

// type PositionService interface {
// 	GetPositionListByUserID(ctx context.Context, userID int) ([]domain.Position, error)
// 	AddInfo(ctx context.Context, i domain.PositionUpdateInfo) error
// }

// type PositionsController struct {
// 	logger          *logger.Logger
// 	positionService PositionService
// 	validator       *validator.Validate
// }

// func NewPositionsController(logger *logger.Logger, positionService PositionService,
// 	validator *validator.Validate) *PositionsController {
// 	return &PositionsController{
// 		logger:          logger,
// 		positionService: positionService,
// 		validator:       validator,
// 	}
// }
