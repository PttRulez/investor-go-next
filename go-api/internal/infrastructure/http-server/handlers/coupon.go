package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/converter"
	"github.com/pttrulez/investor-go-next/go-api/internal/utils"
)

func (c *Handlers) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.CreateCoupon"

	ctx := r.Context()
	var req contracts.CreateCouponRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	// Validate request fields
	if err := c.validator.Struct(req); err != nil {
		var valErr validator.ValidationErrors
		if errors.As(err, &valErr) {
			err = fmt.Errorf("%s: %w", op, err)
			c.logger.Error(err)
			writeValidationErrorsJSON(w, valErr)
			return
		}
	}

	d, err := converter.FromCreateCouponRequestToCoupon(req)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeString(w, http.StatusBadRequest, err.Error())
	}

	// Save new Coupon in DB
	err = c.portfolioService.CreateCoupon(ctx, d,
		utils.GetCurrentUserID(r.Context()))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeString(w, http.StatusCreated, "Coupon created successfully")
}

func (c *Handlers) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	const op = "ExpertController.Delete"

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, fmt.Sprintf("Проблема с конвертацией айди %s: %s",
			chi.URLParam(r, "id"),
			err.Error()))
	}

	ctx := r.Context()
	err = c.opinionService.DeleteExpert(ctx, id, utils.GetCurrentUserID(ctx))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
