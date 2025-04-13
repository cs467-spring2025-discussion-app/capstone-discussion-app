package models_test

import (
	"testing"

	"github.com/matryer/is"

	"godiscauth/internal/models"
	"godiscauth/internal/testutils"
)

func TestNewUser(t *testing.T) {
	is := is.New(t)

	t.Run("valid email and password", func(t *testing.T) {
		validEmail := "test@newuser.com"
		validPassword := testutils.TestingPassword

		_, err := models.NewUser(validEmail, validPassword)
		is.NoErr(err)
	})

	t.Run("email exceeds 254 chars", func(t *testing.T) {
		exceeds254Chars := "a" + string(make([]byte, 253)) + "@example.com"
		validPassword := testutils.TestingPassword
		_, err := models.NewUser(exceeds254Chars, validPassword)
		is.True(err != nil)
	})
}
