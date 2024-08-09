package http_controllers

import (
	"context"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (c *MoexShareController) GetInfoBySecid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	moexShare, err := c.moexSharService.GetBySecid(ctx, chi.URLParam(r, "secid"))
	if err != nil {
		writeError(w, err)
		fmt.Println(err)
		return
	}
	writeOKJSON(w, moexShare)
}

type MoexShareService interface {
	GetBySecid(ctx context.Context, ticker string) (*entity.Share, error)
}
type MoexShareController struct {
	moexSharService MoexShareService
}

func NewMoexShareController(moexService MoexShareService) *MoexShareController {
	return &MoexShareController{
		moexSharService: moexService,
	}
}
