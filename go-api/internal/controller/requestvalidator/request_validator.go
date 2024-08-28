package requestvalidator

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/pttrulez/investor-go/pkg/api"
)

func NewValidator() (*validator.Validate, error) {
	validate := validator.New()
	err := validate.RegisterValidation("dealType", ValidateDealType)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("is-exchange", ValidateExchange)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("moex-market", ValidateMoexMarket)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("opinionType", ValidateOpinion)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("price", ValidatePrice)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("securityType", ValidateSecurityType)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("transactionType", ValidateTransactionType)
	if err != nil {
		return nil, err
	}
	return validate, nil
}

func ValidateExchange(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(api.Exchange).Validate()
}
func ValidateDealType(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(api.DealType).Validate()
}

func ValidateMoexMarket(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(api.ISSMoexMarket).Validate()
}

func ValidateOpinion(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(api.OpinionType).Validate()
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
	return fl.Field().Interface().(api.Exchange).Validate()
}

func ValidateSecurityType(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(api.SecurityType).Validate()
}
func ValidateTransactionType(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(api.TransactionType).Validate()
}
