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
}
