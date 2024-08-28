package httpcontrollers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/controller/converter"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"github.com/pttrulez/investor-go/pkg/api"
)

func (c *OpinionController) CreateOpinion(w http.ResponseWriter, r *http.Request) {
	const op = "OpinionController.CreateOpinion"

	ctx := r.Context()

	// Анмаршалим данные
	var opinionReq api.CreateOpinionRequest
	var err error
	if err = json.NewDecoder(r.Body).Decode(&opinionReq); err != nil {
		c.logger.Error(fmt.Errorf("%s json decoding: %w", op, err))
		writeString(w, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация пришедших данных
	if err = c.validator.Struct(opinionReq); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		writeValidationErrorsJSON(w, validateErr)
		return
	}

	opinion, err := converter.FromCreateOpinionRequestToOpinion(ctx, opinionReq)
	if err != nil {
		c.logger.Error(err)
		writeString(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := c.opinionService.CreateOpinion(ctx, opinion)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (c *OpinionController) DeleteOpinion(w http.ResponseWriter, r *http.Request) {
	const op = "OpinionController.DeleteOpinion"

	ctx := r.Context()

	i := chi.URLParam(r, "id")
	id, err := strconv.Atoi(i)
	if err != nil {
		writeString(w, http.StatusBadRequest, fmt.Sprintf(
			"Проблема с конвертацией айди %s: %s",
			chi.URLParam(r, "id"),
			err.Error()))
		return
	}

	userID := utils.GetCurrentUserID(ctx)

	err = c.opinionService.DeleteOpinionByID(ctx, id, userID)
	if err != nil {
		writeError(w, fmt.Errorf("%s, %w", op, err))
	}

	w.WriteHeader(http.StatusOK)
}

func (c *OpinionController) GetOpinionsList(w http.ResponseWriter, r *http.Request) {
	const op = "OpinionController.GetOpinionsList"

	ctx := r.Context()
	var filters entity.OpinionFilters

	if e := r.URL.Query().Get("exchange"); e != "" {
		exchange := entity.Exchange(e)
		if !exchange.Validate() {
			writeString(w, http.StatusBadRequest, "Неверный формат exchange в запросе")
			return
		}
		filters.Exchange = &exchange
	}

	if id := r.URL.Query().Get("expertId"); id != "" {
		expertID, err := strconv.Atoi(id)
		if err != nil {
			writeString(w, http.StatusBadRequest, "Неверный формат expertId в запросе")
			return
		}
		filters.ExpertID = &expertID
	}

	if id := r.URL.Query().Get("securityId"); id != "" {
		securityID, err := strconv.Atoi(id)
		if err != nil {
			writeString(w, http.StatusBadRequest, "Неверный формат securityId в запросе")
			return
		}
		filters.SecurityID = &securityID
	}

	if t := r.URL.Query().Get("securityType"); t != "" {
		securityType := entity.SecurityType(t)
		if !securityType.Validate() {
			writeString(w, http.StatusBadRequest, "Неверный формат securityType в запросе")
			return
		}
		filters.SecurityType = &securityType
	}

	opinions, err := c.opinionService.GetOpinions(ctx, filters, utils.GetCurrentUserID(ctx))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, opinions)
}

type OpinionService interface {
	CreateOpinion(ctx context.Context, opinion entity.Opinion) (entity.Opinion, error)
	GetOpinions(ctx context.Context, f entity.OpinionFilters, userID int) ([]entity.Opinion,
		error)
	DeleteOpinionByID(ctx context.Context, id int, userID int) error
}

type OpinionController struct {
	logger         Logger
	opinionService OpinionService
	validator      *validator.Validate
}

func NewOpinionController(logger Logger, opinionService OpinionService,
	validator *validator.Validate) *OpinionController {
	return &OpinionController{
		logger:         logger,
		opinionService: opinionService,
		validator:      validator,
	}
}
