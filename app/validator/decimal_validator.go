package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

// DecimalNonNegativeValidator ensures decimal >= 0
func DecimalNonNegativeValidator(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(decimal.Decimal)
	if !ok {
		return false
	}
	return !val.IsNegative()
}

func RegisterCustomValidators(v *validator.Validate) error {
	return v.RegisterValidation("nonnegdecimal", DecimalNonNegativeValidator)
}
