package requestvalidator

import (
	"regexp"
	"strconv"

	"github.com/pttrulez/investor-go/internal/entity"

	"github.com/go-playground/validator/v10"
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
	err = validate.RegisterValidation("opinionType", ValidatePrice)
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
