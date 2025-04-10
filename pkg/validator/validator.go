package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate
var PrettyError string

func AtLeastOneFieldRequired(fl validator.FieldLevel) bool {
	// Check the parent struct value
	parent := fl.Parent()

	// Loop through all the fields in the struct
	for i := 0; i < parent.NumField(); i++ {
		field := parent.Field(i)
		// If the field is not zero, return true
		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			return true
		}
	}

	// If no non-zero field is found, return false
	return false
}

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
