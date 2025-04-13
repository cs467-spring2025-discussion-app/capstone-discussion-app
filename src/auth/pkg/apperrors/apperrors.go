package apperrors

import (
	"errors"
)

var New = errors.New

var (
	ErrDuplicateEmail = New("Email already exists in database")
	ErrEmailMaxLength = New("Email exceeds max length of 254 characters")

	ErrDatabaseIsNil = New("Database is nil")
	ErrUserIsNil     = New("User is nil")

	ErrEmailIsEmpty     = New("Email is empty")
	ErrExpiresAtIsEmpty = New("Expiration time is empty")
	ErrPasswordIsEmpty  = New("Password is empty")
	ErrTokenIsEmpty     = New("Token is empty")
	ErrUserIdEmpty      = New("User ID is empty")

	ErrUserNotFound = New("User not found")

	ErrCouldNotIncrementFailedLogins = New("Could not increment users.failed_login_attempts")
)
