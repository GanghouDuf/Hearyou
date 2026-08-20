package dto

import "project_chat/internel/validation"

type Message struct {
	Type    string `json:"type" validate:"required,oneof=chat"`
	Author  string `json:"author" validate:"required,notblank,min=1,max=50"`
	Payload string `json:"payload" validate:"required,min=1,max=1000"`
}

type ErrorMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (m Message) Validate() error {
	return validation.V.Struct(m)
}
