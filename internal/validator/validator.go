package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var usernameRegex = regexp.MustCompile(`^[가-힣a-zA-Z0-9_]+$`)

// RegisterCustomValidators 커스텀 유효성 검사기 등록
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("username", validateUsername)
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return usernameRegex.MatchString(username)
}
