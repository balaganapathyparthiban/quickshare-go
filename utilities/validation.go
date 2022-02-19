package utilities

import "github.com/go-playground/validator/v10"

func Validation(d interface{}) []map[string]interface{} {
	var errors []map[string]interface{}

	validate := validator.New()
	err := validate.Struct(d)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := make(map[string]interface{})
			element["failedField"] = err.StructNamespace()
			element["tag"] = err.Tag()
			element["value"] = err.Param()
			errors = append(errors, element)
		}
	}

	return errors
}
