package http_controllers

import (
	"github.com/pttrulez/investor-go/internal/utils/http_response"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
)

func (c *MoexBondController) GetInfoByISIN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	moexBond, err := c.services.Moex.Bond.GetByISIN(ctx, chi.URLParam(r, "isin"))
	if err != nil {
		slog.Error(err.Error())
		httpresponse.SendError(w, err)
		return
	}
	httpresponse.WriteOKJSON(w, moexBond)
}

type MoexBondController struct {
	repo     *repository.Repository
	services *service.Container
}

func NewMoexBondController(repo *repository.Repository, services *service.Container) *MoexBondController {
	return &MoexBondController{
		repo:     repo,
		services: services,
	}
}
