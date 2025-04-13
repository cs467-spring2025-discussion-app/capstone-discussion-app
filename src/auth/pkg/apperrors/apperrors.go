package apperrors

import (
	"errors"
)

var New = errors.New

var (
	ErrEmailMaxLength = New("Email exceeds max length of 254 characters")

	ErrExpiresAtIsEmpty = New("Expiration time is empty")
	ErrTokenIsEmpty     = New("Token is empty")
	ErrUserIdEmpty      = New("User ID is empty")
)
