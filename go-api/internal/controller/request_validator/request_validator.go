package request_validator

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("dealType", ValidateDealType)
	validate.RegisterValidation("is-exchange", ValidateExchange)
	validate.RegisterValidation("moex-market", ValidateMoexMarket)
	validate.RegisterValidation("opinionType", ValidatePrice)
	validate.RegisterValidation("price", ValidatePrice)
	validate.RegisterValidation("securityType", ValidateSecurityType)
	validate.RegisterValidation("transactionType", ValidateTransactionType)
	return validate
}

func ValidateExchange(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(entity.Exchange).Validate()
}
func ValidateDealType(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(entity.DealType).Validate()
}

func ValidateMoexMarket(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(entity.ISSMoexMarket).Validate()
}

func ValidateOpinion(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(entity.OpinionType).Validate()
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
	return fl.Field().Interface().(entity.Exchange).Validate()
}

func ValidateSecurityType(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(entity.SecurityType).Validate()
}
func ValidateTransactionType(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(entity.TransactionType).Validate()
}
