package utils

import (
	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info.
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validate struct fields.
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
