package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	*validator.Validate
}

var Val *validator.Validate

func New() *Validator {
	Val = validator.New(validator.WithRequiredStructEnabled())

	Val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})
	return &Validator{
		Validate: Val,
	}

}

func (v *Validator) ValidateStruct(data any) error {
	err := Val.Struct(data)
	if err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			firstError := validationErr[0]
			errName := firstError.Tag()
			fieldName := firstError.Field()

			return fmt.Errorf("field '%s' is %s", fieldName, errName)
		}
		return err
	}
	return nil
}
