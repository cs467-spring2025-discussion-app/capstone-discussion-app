package services_test

import (
	"testing"

	"github.com/matryer/is"

	"godiscauth/internal/models"
	"godiscauth/internal/repository"
	"godiscauth/internal/services"
	"godiscauth/internal/testutils"
	"godiscauth/pkg/apperrors"
)

// TestUserService_NewUserService tests the creation of a new UserService
func TestUserService_NewUserService(t *testing.T) {
	is := is.New(t)
	testDB := testutils.TestDBSetup()

	t.Run("returns err with nil user repo", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()

		sr, err := repository.NewSessionRepository(tx)
		is.Equal(err, nil)

		userService, err := services.NewUserService(nil, sr)
		is.Equal(userService, nil)
		is.Equal(err, apperrors.ErrUserRepoIsNil)
	})
	t.Run("returns err with nil session repo", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()

		ur, err := repository.NewUserRepository(tx)
		is.Equal(err, nil)

		userService, err := services.NewUserService(ur, nil)
		is.Equal(userService, nil)
		is.Equal(err, apperrors.ErrSessionRepoIsNil)
	})

	t.Run("creates user service", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()

		ur, err := repository.NewUserRepository(tx)
		is.Equal(err, nil)
		sr, err := repository.NewSessionRepository(tx)
		is.Equal(err, nil)

		userService, err := services.NewUserService(ur, sr)
		is.True(userService != nil)
		is.NoErr(err)
	})
}

// TestUserService_RegisterUser tests that the user service can register a user
func TestUserService_RegisterUser(t *testing.T) {
	is := is.New(t)

	testDB := testutils.TestDBSetup()
	tx := testDB.Begin()
	defer tx.Rollback()

	// Set up user service
	userRepo, err := repository.NewUserRepository(tx)
	is.NoErr(err)
	sessionRepo, err := repository.NewSessionRepository(tx)
	is.NoErr(err)
	us, err := services.NewUserService(userRepo, sessionRepo)
	is.NoErr(err)

	// Error on empty email
	t.Run("empty email", func(t *testing.T) {
		err := us.RegisterUser("", "password")
		is.Equal(err, apperrors.ErrEmailIsEmpty)
	})

	// Error on empty password
	t.Run("empty password", func(t *testing.T) {
		err := us.RegisterUser("some@test.com", "")
		is.Equal(err, apperrors.ErrPasswordIsEmpty)
	})

	// Can register with a valid User value
	t.Run("valid user", func(t *testing.T) {
		email := "testUserServiceRegisterUser@test.com"
		password := testutils.TestingPassword
		err := us.RegisterUser(email, password)
		is.NoErr(err)

		// Check user actually exists in db
		var user models.User
		result := tx.First(&user, "email = ?", email)
		is.NoErr(result.Error)
	})
}
