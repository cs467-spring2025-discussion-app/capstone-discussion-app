package apperrors

import (
	"errors"
)

var New = errors.New

var ErrEmailMaxLength = New("Email exceeds max length of 254 characters")
