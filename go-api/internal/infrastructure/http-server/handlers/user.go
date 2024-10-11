package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/converter"
	"github.com/pttrulez/investor-go-next/go-api/internal/utils"
)

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UpdateUser"

	ctx := r.Context()
	var req contracts.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, fmt.Errorf("%s: %w", op, err))
		return
	}

	user := converter.FromUpdateUserRequestToUser(req)
	user.ID = utils.GetCurrentUserID(ctx)
	err := h.userService.UpdateUser(ctx, user)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		h.logger.Error(err)
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
