package http_controllers

import (
	"context"
	"github.com/pttrulez/investor-go/internal/entity"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (c *MoexBondController) GetInfoByISIN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	moexBond, err := c.moexBondService.GetByISIN(ctx, chi.URLParam(r, "isin"))
	if err != nil {
		slog.Error(err.Error())
		writeError(w, err)
		return
	}
	writeOKJSON(w, moexBond)
}

type MoexService interface {
	GetByISIN(ctx context.Context, isin string) (*entity.Bond, error)
}
type MoexBondController struct {
	moexBondService MoexService
}

func NewMoexBondController(moexService MoexService) *MoexBondController {
	return &MoexBondController{
		moexBondService: moexService,
	}
}
