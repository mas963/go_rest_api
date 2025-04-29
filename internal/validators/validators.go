package validators

import "github.com/go-playground/validator/v10"

func StrongPasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return len(password) >= 8 && containsUppercase(password) && containsNumber(password)
}

func containsUppercase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func containsNumber(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}