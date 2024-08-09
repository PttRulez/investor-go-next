package http_controllers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/controller/model/converter"
	"github.com/pttrulez/investor-go/internal/controller/model/dto"
)

func (c *TransactionController) CreateNewTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Анмаршалим данные
	var tData = new(dto.CreateTransaction)
	if err := json.NewDecoder(r.Body).Decode(tData); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация пришедших данных
	if err := c.validator.Struct(tData); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		writeValidationErrorsJSON(w, validateErr)
		return
	}

	// Создаем
	t := converter.FromCreateDtoToTransaction(tData)
	t.UserId = utils.GetCurrentUserId(r.Context())
	err := c.transactionService.CreateTransaction(ctx, t)
	if err != nil {
		writeError(w, err)
	}

	writeString(w, http.StatusCreated, "Транзакция создана")
}

func (c *TransactionController) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	cashoutId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
	}

	err = c.transactionService.DeleteTransaction(r.Context(), cashoutId, utils.GetCurrentUserId(r.Context()))
	if err != nil {
		writeError(w, err)
	}
	writeString(w, http.StatusOK, "Транзакция удалена")

}

type TransactionService interface {
	CreateTransaction(ctx context.Context, transactionData *entity.Transaction) error
	DeleteTransaction(ctx context.Context, transactionId int, userId int) error
}

type TransactionController struct {
	transactionService TransactionService
	validator          *validator.Validate
}

func NewCashoutController(transactionService TransactionService, validator *validator.Validate) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
		validator:          validator,
	}
}
