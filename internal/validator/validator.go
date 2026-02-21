package validator

import (
    "errors"
    "regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
    if email == "" {
        return errors.New("email is required")
    }
    if len(email) > 254 {
        return errors.New("email too long")
    }
    if !emailRegex.MatchString(email) {
        return errors.New("invalid email format")
    }
    return nil
}
