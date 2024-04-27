package controller

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pttrulez/investor-go/internal/repository"
	"github.com/pttrulez/investor-go/internal/service"
	httpresponse "github.com/pttrulez/investor-go/pkg/http-response"
)

func (c *MoexShareController) GetInfoByTicker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	moexShare, err := c.services.Moex.Share.GetByTicker(ctx, chi.URLParam(r, "ticker"))
	if err != nil {
		httpresponse.SendError(w, err)
		fmt.Println(err)
		return
	}
	httpresponse.WriteOKJSON(w, moexShare)
}

type MoexShareController struct {
	repo     *repository.Repository
	services *service.Container
}

func NewMoexShareController(repo *repository.Repository, services *service.Container) *MoexShareController {
	return &MoexShareController{
		repo:     repo,
		services: services,
	}
}
