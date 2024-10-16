package handlers

import (
	"bytes"
	"encoding/json"
	"errors"

	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go-next/go-api/internal/service"
)

func writeErr(w http.ResponseWriter, err error) error {
	var code int
	if errors.Is(err, service.ErrDomainNotFound) {
		code = http.StatusNotFound
		writeString(w, code, "Не найдено")
	}

	code = http.StatusInternalServerError
	writeString(w, code, err.Error())

	return APIError{
		Code: code,
		Err:  err,
	}
}

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrDomainNotFound) {
		writeString(w, http.StatusNotFound, "Не найдено")
		return
	}

	writeString(w, http.StatusInternalServerError, err.Error())
}

func writeJS(w http.ResponseWriter, status int, value any) error {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(true)
	err := encoder.Encode(value)
	if err != nil {
		return APIError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	w.Header().Set("Content-Type", "applications/json; charset=utf-8")
	w.WriteHeader(status)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return APIError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
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
			mappedErrors[err.Field()] += "Неверное имя биржи. на данный момент поддерживаются только следующие: MOEX"
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

type APIError struct {
	Code int
	Err  error
}

func (e APIError) Error() string {
	return e.Err.Error()
}
