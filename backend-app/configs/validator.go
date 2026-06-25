package configs

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func NewValidator() *validator.Validate {
	v := validator.New()

	v.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		date := fl.Field().String()
		regex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		return regex.MatchString(date)
	})

	return v
}
