package handlers

import (
	"fmt"
	"net/http"

	"github.com/pttrulez/investor-go/internal/api/converter"

	"github.com/go-chi/chi/v5"
)

func (c *Handlers) GetBondInfoByTicker(w http.ResponseWriter, r *http.Request) {
	const op = "MoexBondController.GetBondInfoByTicker"

	ctx := r.Context()

	moexBond, err := c.moexService.GetBondByTicker(ctx, chi.URLParam(r, "ticker"))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.FromMoexBondToMoexBondResponse(moexBond))
}

// type MoexBondController struct {
// 	logger      *logger.Logger
// 	moexService interfaces.MoexService
// }

// func NewMoexBondController(logger *logger.Logger, moexService interfaces.MoexService) *MoexBondController {
// 	return &MoexBondController{
// 		logger:      logger,
// 		moexService: moexService,
// 	}
// }
