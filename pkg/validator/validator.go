package validator

import "github.com/go-playground/validator/v10"

var Validator *validator.Validate

func NewValidator() {
	Validator = validator.New(validator.WithRequiredStructEnabled())
}

func Validate(data interface{}) error {
	err := Validator.Struct(data)
	if err != nil {
		return err
	}
	return nil
}
