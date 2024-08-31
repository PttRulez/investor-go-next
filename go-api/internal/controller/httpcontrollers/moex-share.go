package httpcontrollers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pttrulez/investor-go/internal/controller/converter"
	"github.com/pttrulez/investor-go/internal/entity"

	"github.com/go-chi/chi/v5"
)

func (c *MoexShareController) GetInfoByTicker(w http.ResponseWriter, r *http.Request) {
	const op = "MoexShareController.GetInfoByTicker"

	ctx := r.Context()

	moexShare, err := c.moexSharService.GetByTicker(ctx, chi.URLParam(r, "ticker"))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.FromMoexShareToMoexShareResponse(moexShare))
}

type MoexShareService interface {
	GetByTicker(ctx context.Context, ticker string) (entity.Share, error)
}
type MoexShareController struct {
	logger          Logger
	moexSharService MoexShareService
}

func NewMoexShareController(logger Logger, moexService MoexShareService) *MoexShareController {
	return &MoexShareController{
		logger:          logger,
		moexSharService: moexService,
	}
}
