package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var V = validator.New()

func init() {
	V.RegisterValidation("notblank", notBlank)
}

func notBlank(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}
