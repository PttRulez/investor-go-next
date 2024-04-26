package controller

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pttrulez/investor-go/internal/lib/response"
	"github.com/pttrulez/investor-go/internal/services"
	"github.com/pttrulez/investor-go/internal/types"
)

func (c *MoexBondController) GetInfoByISIN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	moexBond, err := c.services.MoexBond.GetByISIN(ctx, chi.URLParam(r, "isin"))
	if err != nil {
		slog.Error(err.Error())
		response.SendError(w, err)
		return
	}
	response.WriteOKJSON(w, moexBond)
}

type MoexBondController struct {
	repo     *types.Repository
	services *services.ServiceContainer
}

func NewMoexBondController(repo *types.Repository, services *services.ServiceContainer) *MoexBondController {
	return &MoexBondController{
		repo:     repo,
		services: services,
	}
}
