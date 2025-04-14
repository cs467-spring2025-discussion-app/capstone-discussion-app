package services_test

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"golang.org/x/crypto/bcrypt"

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

		userService := setupUserService(t)
		is.True(userService != nil)
	})
}

// TestUserService_RegisterUser tests that the user service can register a user
func TestUserService_RegisterUser(t *testing.T) {
	is := is.New(t)

	us := setupUserService(t)

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
		result := us.UserRepo.DB.First(&user, "email = ?", email)
		is.NoErr(result.Error)
	})
}

// TestUserService_UpdateUser test user updates into the database
func TestUserService_UpdateUser(t *testing.T) {
	is := is.New(t)

	// Error on empty user ID
	t.Run("empty user ID", func(t *testing.T) {
		us := setupUserService(t)
		err := us.UpdateUser("", map[string]any{})
		is.Equal(err, apperrors.ErrUserIdEmpty)
	})

	// Error on non-existent user
	t.Run("non-existent user", func(t *testing.T) {
		us := setupUserService(t)
		err := us.UpdateUser("doesNotExist", map[string]any{})
		is.Equal(err, apperrors.ErrUserNotFound)
	})

	// Can update user
	t.Run("can update user", func(t *testing.T) {
		// Register a user to update
		email := "testUpdateUser@test.com"
		password := testutils.TestingPassword
		us := setupUserService(t)

		err := us.RegisterUser(email, password)
		is.NoErr(err)

		// Update user
		referenceTime := time.Now().Add(time.Hour * 24).UTC().Truncate(time.Second)
		var user models.User
		result := us.UserRepo.DB.First(&user, "email = ?", email)
		is.NoErr(result.Error)
		err = us.UpdateUser(user.ID.String(), map[string]any{
			"email":                 "newUserName@test.com",
			"password":              "new" + password,
			"last_login":            referenceTime,
			"failed_login_attempts": 99,
			"account_locked":        true,
			"account_locked_until":  referenceTime,
		})
		is.NoErr(err)

		// Get updated user and check for updated fields
		updatedUser, err := us.UserRepo.GetUserByID(user.ID.String())
		is.NoErr(err)
		t.Run("updates email", func(t *testing.T) {
			is.Equal(updatedUser.Email, "newUserName@test.com")
		})
		t.Run("updates password", func(t *testing.T) {
			err := bcrypt.CompareHashAndPassword(
				[]byte(updatedUser.Password),
				[]byte("new"+password),
			)
			is.NoErr(err)
		})
		t.Run("updates last_login", func(t *testing.T) {
			is.Equal(updatedUser.LastLogin, &referenceTime)
		})
		t.Run("updates failed_login_attempts", func(t *testing.T) {
			is.Equal(updatedUser.FailedLoginAttempts, 99)
		})
		t.Run("updates account_locked", func(t *testing.T) {
			is.True(updatedUser.AccountLocked)
		})
		t.Run("updates account_locked_until", func(t *testing.T) {
			is.Equal(updatedUser.LastLogin, &referenceTime)
		})
	})
}

func setupUserService(t *testing.T) *services.UserService {
	t.Helper()

	testDB := testutils.TestDBSetup()
	tx := testDB.Begin()
	t.Cleanup(func() { tx.Rollback() })

	ur, err := repository.NewUserRepository(tx)
	if err != nil {
		t.Fatalf("failed to create user repository: %v", err)
	}
	sr, err := repository.NewSessionRepository(tx)
	if err != nil {
		t.Fatalf("failed to create session repository: %v", err)
	}
	us, err := services.NewUserService(ur, sr)
	if err != nil {
		t.Fatalf("failed to create user service: %v", err)
	}
	return us
}
