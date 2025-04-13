package models_test

import (
	"testing"

	"github.com/matryer/is"

	"godiscauth/internal/models"
	"godiscauth/internal/testutils"
	"godiscauth/pkg/apperrors"
)

func TestNewUser(t *testing.T) {
	is := is.New(t)

	validEmail := "test@newuser.com"
	validPassword := testutils.TestingPassword
	t.Run("valid email and password", func(t *testing.T) {
		_, err := models.NewUser(validEmail, validPassword)
		is.NoErr(err)
	})

	t.Run("email exceeds 254 chars", func(t *testing.T) {
		exceeds254Chars := "example@"
		for range 256 {
			exceeds254Chars += "a"
		}
		exceeds254Chars += ".com"

		_, err := models.NewUser(exceeds254Chars, validPassword)
		is.True(err == apperrors.ErrEmailMaxLength)
	})

	invalidFormats := map[string]string{
		"emptyString":     "",
		"missingAt":       "example.com",
		"multipleAts":     "at@at@at.com",
		"missingLocal":    "@at.com",
		"missingDomain":   "example@",
		"consecutiveDots": "example@at..com",
		"leadingDot":      ".example@at.com",
		"trailingDot":     "exmaple@at.com.",
		"space":           "example @at.com",
	}

	for name, email := range invalidFormats {
		t.Run(name, func(t *testing.T) {
			_, err := models.NewUser(email, validPassword)
			is.True(err != nil)
		})
	}

	// TODO: test password complexity
}
