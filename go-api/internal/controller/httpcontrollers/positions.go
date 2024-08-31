package httpcontrollers

import (
	"context"
	"fmt"
	"net/http"

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

type PostionService interface {
	GetListByUserID(ctx context.Context, userID int) ([]entity.Position, error)
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
