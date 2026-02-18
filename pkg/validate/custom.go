package validate

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("safe_string", validateSafeString)
	v.RegisterValidation("alphanum_unicode", validateAlphanumUnicode)
	v.RegisterValidation("safe_url", validateSafeURL)
}

func validateSafeString(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	scriptPattern := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	if scriptPattern.MatchString(s) {
		return false
	}

	eventPattern := regexp.MustCompile(`(?i)on\w+\s*=`)
	if eventPattern.MatchString(s) {
		return false
	}

	return true
}

func validateAlphanumUnicode(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' {
			return false
		}
	}
	return true
}

func validateSafeURL(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	if regexp.MustCompile(`(?i)^javascript:`).MatchString(s) {
		return false
	}

	if regexp.MustCompile(`(?i)^data:`).MatchString(s) {
		return false
	}

	return true
}
