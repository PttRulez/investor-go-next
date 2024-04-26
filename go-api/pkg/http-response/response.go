package httpresponse

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"

	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func SendError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotYours) || errors.Is(err, sql.ErrNoRows) {
		WriteNotFound(w)
	} else if e, ok := err.(ErrSendToClient); ok {
		if e.Status == 0 {
			e.Status = http.StatusInternalServerError
		}
		WriteString(w, e.Status, e.Error())
	} else {
		SmthWentWrong(w)
	}
}
func SmthWentWrong(w http.ResponseWriter) {
	WriteString(w, http.StatusInternalServerError, "Что-то пошло не так")
}
func ValidationErrsToResponse(errs validator.ValidationErrors) map[string]string {
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
func WriteCreatedJSON(w http.ResponseWriter, value any) (int, error) {
	return WriteJSON(w, http.StatusCreated, value)
}
func WriteJSON(w http.ResponseWriter, status int, value any) (int, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(true)
	encoder.Encode(value)

	w.Header().Set("Content-Type", "applications/json; charset=utf-8")
	w.WriteHeader(status)

	return w.Write(buf.Bytes())
}
func WriteNotFound(w http.ResponseWriter) {
	WriteString(w, http.StatusNotFound, "Здесь такого нет")
}
func WriteOKJSON(w http.ResponseWriter, value any) (int, error) {
	return WriteJSON(w, http.StatusOK, value)
}
func WriteString(w http.ResponseWriter, status int, value string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(value))
}
func WriteValidationErrorsJSON(w http.ResponseWriter, errs validator.ValidationErrors) (int, error) {
	return WriteJSON(w, http.StatusUnprocessableEntity, ValidationErrsToResponse(errs))
}
