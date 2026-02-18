package service

import "errors"

var (
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrUsernameExists      = errors.New("username already exists")
	ErrWrongPassword       = errors.New("wrong password")
	ErrSamePassword        = errors.New("new password same as current")
	ErrPasswordMismatch    = errors.New("password confirmation mismatch")
	ErrCannotWithdrawAdmin = errors.New("admin cannot withdraw")
)
