// Package validate provides validation utilities for the application.
package validate

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Struct(s interface{}) error {
	return validate.Struct(s)
}
