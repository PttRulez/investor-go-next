package httpcontrollers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pttrulez/investor-go/internal/controller/converter"
	"github.com/pttrulez/investor-go/internal/entity"

	"github.com/go-chi/chi/v5"
)

func (c *MoexBondController) GetInfoByTicker(w http.ResponseWriter, r *http.Request) {
	const op = "MoexBondController.GetInfoByTicker"

	ctx := r.Context()

	moexBond, err := c.moexBondService.GetByTicker(ctx, chi.URLParam(r, "ticker"))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.FromMoexBondToMoexBondResponse(moexBond))
}

type MoexBondService interface {
	GetByTicker(ctx context.Context, isin string) (entity.Bond, error)
}
type MoexBondController struct {
	logger          Logger
	moexBondService MoexBondService
}

func NewMoexBondController(logger Logger, moexService MoexBondService) *MoexBondController {
	return &MoexBondController{
		logger:          logger,
		moexBondService: moexService,
	}
}
