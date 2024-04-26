package service

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/internal/model"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("dealType", ValidateDealType)
	validate.RegisterValidation("is-exchange", ValidateExchange)
	validate.RegisterValidation("opinionType", ValidatePrice)
	validate.RegisterValidation("price", ValidatePrice)
	validate.RegisterValidation("securityType", ValidateSecurityType)
	return validate
}

func ValidateExchange(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(model.Exchange).Validate()
}

func ValidateSecurityType(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(model.SecurityType).Validate()
}

func ValidateDealType(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(model.DealType).Validate()
}

func ValidateOpinion(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(model.OpinionType).Validate()
}

func ValidatePrice(fl validator.FieldLevel) bool {
	priceFloat, ok := fl.Field().Interface().(float64)
	if !ok {
		return false
	}
	match, _ := regexp.MatchString(`^\d+$|^\d+\.\d{1,2}$`, strconv.FormatFloat(priceFloat, 'f', -1, 64))
	return match
}

func ValidateRole(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(model.Exchange).Validate()
}
