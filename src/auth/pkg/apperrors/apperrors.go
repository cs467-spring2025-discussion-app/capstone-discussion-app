package apperrors

import (
	"errors"
)

var New = errors.New

var (
	ErrEmailMaxLength = New("Email exceeds max length of 254 characters")

	ErrDatabaseIsNil = New("Database is nil")
	ErrUserIsNil     = New("User is nil")

	ErrExpiresAtIsEmpty = New("Expiration time is empty")
	ErrTokenIsEmpty     = New("Token is empty")
	ErrUserIdEmpty      = New("User ID is empty")
)
