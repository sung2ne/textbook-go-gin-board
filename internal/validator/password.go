package validator

import (
    "regexp"
    "unicode"

    "github.com/go-playground/validator/v10"
)

// RegisterPasswordValidators 비밀번호 유효성 검사기 등록
func RegisterPasswordValidators(v *validator.Validate) {
    v.RegisterValidation("password", validatePassword)
}

// validatePassword 비밀번호 정책 검증
// - 최소 8자
// - 대문자 1개 이상
// - 소문자 1개 이상
// - 숫자 1개 이상
// - 특수문자 1개 이상
func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    if len(password) < 8 {
        return false
    }

    var hasUpper, hasLower, hasDigit, hasSpecial bool

    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsDigit(char):
            hasDigit = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }

    return hasUpper && hasLower && hasDigit && hasSpecial
}
