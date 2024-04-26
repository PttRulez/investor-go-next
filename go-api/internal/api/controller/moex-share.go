package controller

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pttrulez/investor-go/internal/lib/response"
	"github.com/pttrulez/investor-go/internal/services"
	"github.com/pttrulez/investor-go/internal/types"
)

func (c *MoexShareController) GetInfoByTicker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	moexShare, err := c.services.MoexShare.GetByTicker(ctx, chi.URLParam(r, "ticker"))
	if err != nil {
		response.SendError(w, err)
		fmt.Println(err)
		return
	}
	response.WriteOKJSON(w, moexShare)
}

type MoexShareController struct {
	repo     *types.Repository
	services *services.ServiceContainer
}

func NewMoexShareController(repo *types.Repository, services *services.ServiceContainer) *MoexShareController {
	return &MoexShareController{
		repo:     repo,
		services: services,
	}
}
