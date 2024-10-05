package validators

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func PasswordPatternValidator(fl validator.FieldLevel) bool {
	specialChars := "!@#$%^&*()_+=[{}];':\"\\|,.<>?/-"

	hasLetter := false
	hasNumber := false
	hasSpecial := false

	for _, char := range fl.Field().String() {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
		case unicode.IsDigit(char):
			hasNumber = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	return hasLetter && hasNumber && hasSpecial
}
