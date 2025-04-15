package validatorLib

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	var errorMessages []string
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "email":
				errorMessages = append(errorMessages, "Invalid Email Format")
			case "required":
				errorMessages = append(errorMessages, "Field "+err.Field()+" wajib diisi")
			case "min":
				if err.Field() == "Password" {
					errorMessages = append(errorMessages, "Password harus 8 karakter.")
				}
			case "eqfield":
				errorMessages = append(errorMessages, err.Field()+" harus sama dengan "+err.Param()+".")
			default:
				errorMessages = append(errorMessages, "Field "+err.Field()+" tidak valid")
			}
		}

		return errors.New("Validasi gagal: " + joinMessage(errorMessages))
	}

	return nil
}

func joinMessage(messages []string) string {
	result := ""
	for i, message := range messages {
		if i > 0 {
			result += ", "
		}
		result += message
	}

	return result
}
