package http_controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	ierrors "github.com/pttrulez/investor-go/internal/errors"

	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, ierrors.ErrNotYours) {
		writeString(w, http.StatusNotFound, ierrors.ErrNotFound.Error())
		return
	}

	var e ierrors.ErrSendToClient
	if errors.As(err, &e) {
		if e.Status == 0 {
			e.Status = http.StatusInternalServerError
		}
		writeString(w, e.Status, e.Error())
	}
}

func validationErrsToResponse(errs validator.ValidationErrors) map[string]string {
	mappedErrors := map[string]string{}

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			mappedErrors[err.Field()] += fmt.Sprintf("Поле %s обязательно для заполнения", err.Field())
		case "email":
			mappedErrors[err.Field()] += fmt.Sprintf("Поле %s должно быть валидным email'ом", err.Field())
		case "price":
			mappedErrors[err.Field()] += fmt.Sprintf(
				"Поле %s должно быть валидной ценой. Это либо целое число либо, десятичное с 1 или 2 знаками после запятой",
				err.Field(),
			)
		case "is-exchange":
			mappedErrors[err.Field()] += "Неверное имя биржи. на данный момент поддерживаются только следующие: Moex"
		case "securityType":
			mappedErrors[err.Field()] += "Указан неправильный тип бумаги"
		case "dealType":
			mappedErrors[err.Field()] += "Тип сделки может быть либо BUY либо SELL"
		default:
			mappedErrors[err.Field()] += fmt.Sprintf("Неверно заполнено поле %s", err.Field())
		}
	}

	return mappedErrors
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(true)
	err := encoder.Encode(value)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "applications/json; charset=utf-8")
	w.WriteHeader(status)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return
	}
}

func writeString(w http.ResponseWriter, status int, value string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, err := w.Write([]byte(value))
	if err != nil {
		return
	}
}
func writeValidationErrorsJSON(w http.ResponseWriter, errs validator.ValidationErrors) {
	writeJSON(w, http.StatusUnprocessableEntity, validationErrsToResponse(errs))
}

// Сахар
func writeNotFound(w http.ResponseWriter) {
	writeString(w, http.StatusNotFound, "Здесь такого нет")
}
func writeOKJSON(w http.ResponseWriter, value any) {
	writeJSON(w, http.StatusOK, value)
}
func writeInternal(w http.ResponseWriter) {
	writeString(w, http.StatusInternalServerError, ierrors.ErrSmthWrong.Error())
}
