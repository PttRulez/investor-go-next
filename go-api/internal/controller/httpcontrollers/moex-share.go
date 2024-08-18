package httpcontrollers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pttrulez/investor-go/internal/controller/model/converter"
	"github.com/pttrulez/investor-go/internal/entity"

	"github.com/go-chi/chi/v5"
)

func (c *MoexShareController) GetInfoBySecid(w http.ResponseWriter, r *http.Request) {
	const op = "MoexShareController.GetInfoBySecid"

	ctx := r.Context()

	moexShare, err := c.moexSharService.GetBySecid(ctx, chi.URLParam(r, "secid"))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.FromMoexShareToMoexShareResponse(*moexShare))
}

type MoexShareService interface {
	GetBySecid(ctx context.Context, ticker string) (*entity.Share, error)
}
type MoexShareController struct {
	logger          Logger
	moexSharService MoexShareService
}

func NewMoexShareController(logger Logger, moexService MoexShareService) *MoexShareController {
	return &MoexShareController{
		logger:          logger,
		moexSharService: moexService,
	}
}
