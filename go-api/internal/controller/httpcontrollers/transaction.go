package httpcontrollers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"github.com/pttrulez/investor-go/pkg/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/controller/model/converter"
)

func (c *TransactionController) CreateNewTransaction(w http.ResponseWriter, r *http.Request) {
	const op = "ExpertController.CreateNewExpert"

	// Анмаршалим данные
	var tData api.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&tData); err != nil {
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

	// Конвертируем реквест в модель
	t, err := converter.FromCreateTransactionRequestToTransaction(tData)
	if err != nil {
		writeString(w, http.StatusBadRequest, err.Error())
		return
	}

	// Получаем айди юзера и создаем транзакцию
	ctx := r.Context()
	t.UserID = utils.GetCurrentUserID(ctx)
	tr, err := c.transactionService.CreateTransaction(ctx, t)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, tr)
}

func (c *TransactionController) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	const op = "ExpertController.DeleteTransaction"

	cashoutID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeString(w, http.StatusBadRequest, fmt.Sprintf("Проблема с конвертацией айди %s: %s",
			chi.URLParam(r, "id"),
			err.Error()))
	}

	ctx := r.Context()
	err = c.transactionService.DeleteTransaction(ctx, cashoutID, utils.GetCurrentUserID(ctx))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		c.logger.Error(err)
		writeError(w, err)
		return
	}

	writeString(w, http.StatusOK, "Транзакция удалена")
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, transactionData entity.Transaction) (
		entity.Transaction, error)
	DeleteTransaction(ctx context.Context, transactionID int, userID int) error
}

type TransactionController struct {
	logger             Logger
	transactionService TransactionService
	validator          *validator.Validate
}

func NewCashoutController(logger Logger, transactionService TransactionService,
	validator *validator.Validate) *TransactionController {
	return &TransactionController{
		logger:             logger,
		transactionService: transactionService,
		validator:          validator,
	}
}
