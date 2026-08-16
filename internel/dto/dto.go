package dto

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type Message struct {
	Type    string `json:"type" validate:"required,oneof=chat"`
	Author  string `json:"author" validate:"required,notblank,min=1,max=50"`
	Payload string `json:"payload" validate:"required,min=1,max=1000"`
}

type ErrorMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

var validate = validator.New()

func (m Message) Validate() error {
	return validate.Struct(m)
}

func notBlank(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}

func init() {
	validate.RegisterValidation("notblank", notBlank)
}
