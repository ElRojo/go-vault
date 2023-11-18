package utility

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type validationErrors struct {
	required []string
	equality []string
}

func ValidateRequestFields(s interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(jsonTagName)

	if err := validate.Struct(s); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			valErrors := extractValidationErrors(ve)

			var messages []string
			if len(valErrors.required) > 0 {
				messages = append(messages, fmt.Sprintf("request missing required field(s): %s", strings.Join(valErrors.required, ", ")))
			}
			if len(valErrors.equality) > 0 {
				messages = append(messages, fmt.Sprintf("incompatible field values: %s", strings.Join(valErrors.equality, ", ")))
			}

			return fmt.Errorf(strings.Join(messages, " "))
		}
		return err
	}
	return nil
}

func jsonTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 1)[0]
	if name == "-" {
		return ""
	}
	return name
}

func extractValidationErrors(ve validator.ValidationErrors) validationErrors {
	var valErrors validationErrors
	for _, fe := range ve {
		switch tag := fe.Tag(); tag {
		case "nefield":
			valErrors.equality = append(valErrors.equality, fmt.Sprintf("%s: %s", fe.Field(), fe.Param()))
		case "required":
			valErrors.required = append(valErrors.required, fe.Field())
		}
	}

	return valErrors
}
