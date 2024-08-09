package http_controllers

import (
	"context"
	"github.com/pttrulez/investor-go/internal/entity"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (c *MoexBondController) GetInfoBySecid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	moexBond, err := c.moexBondService.GetBySecid(ctx, chi.URLParam(r, "secid"))
	if err != nil {
		logger.Error(err)
		writeError(w, err)
		return
	}
	writeOKJSON(w, moexBond)
}

type MoexService interface {
	GetBySecid(ctx context.Context, isin string) (*entity.Bond, error)
}
type MoexBondController struct {
	moexBondService MoexService
}

func NewMoexBondController(moexService MoexService) *MoexBondController {
	return &MoexBondController{
		moexBondService: moexService,
	}
}
