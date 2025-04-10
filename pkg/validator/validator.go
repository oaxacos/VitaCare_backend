package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate
var PrettyError string

func NewValidator() {
	Validator = validator.New(validator.WithRequiredStructEnabled())

	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

}

func Validate(data interface{}) error {

	err := Validator.Struct(data)
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
