package handlers

import (
	"fmt"
	"net/http"

	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/converter"

	"github.com/go-chi/chi/v5"
)

func (c *Handlers) GetShareInfoByTicker(w http.ResponseWriter, r *http.Request) {
	const op = "MoexShareController.GetShareInfoByTicker"

	ctx := r.Context()

	moexShare, err := c.moexService.GetShareByTicker(ctx, chi.URLParam(r, "ticker"))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.FromMoexShareToMoexShareResponse(moexShare))
}
