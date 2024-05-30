package http_controllers

import (
	"context"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (c *MoexShareController) GetInfoByTicker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	moexShare, err := c.moexSharService.GetByTicker(ctx, chi.URLParam(r, "ticker"))
	if err != nil {
		writeError(w, err)
		fmt.Println(err)
		return
	}
	writeOKJSON(w, moexShare)
}

type MoexShareService interface {
	GetByTicker(ctx context.Context, ticker string) (*entity.Share, error)
}
type MoexShareController struct {
	moexSharService MoexShareService
}

func NewMoexShareController(moexService MoexShareService) *MoexShareController {
	return &MoexShareController{
		moexSharService: moexService,
	}
}
