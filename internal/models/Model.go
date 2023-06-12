package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt *time.Time     `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		var fieldErr validator.FieldError
		for _, fieldErr = range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = fieldErr.StructNamespace()
			element.Tag = fieldErr.Tag()
			element.Value = fieldErr.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
