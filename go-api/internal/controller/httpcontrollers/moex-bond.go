package httpcontrollers

import (
	"context"
	"net/http"

	"github.com/pttrulez/investor-go/internal/entity"

	"github.com/go-chi/chi/v5"
)

func (c *MoexBondController) GetInfoBySecid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	moexBond, err := c.moexBondService.GetBySecid(ctx, chi.URLParam(r, "secid"))
	if err != nil {
		c.logger.Error(err)
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, moexBond)
}

type MoexBondService interface {
	GetBySecid(ctx context.Context, isin string) (*entity.Bond, error)
}
type MoexBondController struct {
	logger          Logger
	moexBondService MoexBondService
}

func NewMoexBondController(moexService MoexBondService, logger Logger) *MoexBondController {
	return &MoexBondController{
		logger:          logger,
		moexBondService: moexService,
	}
}
