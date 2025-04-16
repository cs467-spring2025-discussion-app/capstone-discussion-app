package repository_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/matryer/is"
	"gorm.io/gorm"

	"godiscauth/internal/models"
	"godiscauth/internal/repository"
	"godiscauth/internal/testutils"
	"godiscauth/pkg/apperrors"
	"godiscauth/pkg/config"
)

// TestUserRepository_NewUserRepository tests creation of UserRepository
// structs in the `repository` package
func TestUserRepository_NewUserRepository(t *testing.T) {
	is := is.New(t)

	testDB := testutils.TestDBSetup()

	t.Run("creates new user repo", func(t *testing.T) {
		ur, err := repository.NewUserRepository(testDB)
		is.True(ur != nil)
		is.NoErr(err)
	})

	t.Run("returns err with nil db", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()

		ur, err := repository.NewUserRepository(nil)
		is.Equal(ur, nil)
		is.Equal(err, apperrors.ErrDatabaseIsNil)
	})
}

// TestUserRepository_RegisterUser tests insertion of new users into the
// `users` table of the database
func TestUserRepository_RegisterUser(t *testing.T) {
	is := is.New(t)

	// Validates user before registration
	t.Run("fails on nil user", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		err = ur.RegisterUser(nil)
		is.Equal(err, apperrors.ErrUserIsNil)
	})
	t.Run("fails on missing email", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		user := &models.User{
			Password: "password",
		}
		err = ur.RegisterUser(user)
		is.Equal(err, apperrors.ErrEmailIsEmpty)
	})
	t.Run("fails on missing password", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		user := &models.User{
			Email: "testRegisterUser@test.com",
		}
		err = ur.RegisterUser(user)
		is.Equal(err, apperrors.ErrPasswordIsEmpty)
	})

	// Inserts user into db with complete User value
	t.Run("registers user", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		user := &models.User{
			Email:    "testRegisterUser@test.com",
			Password: "password",
		}
		err = ur.RegisterUser(user)
		is.NoErr(err)

		// Lookup expected user in the db
		var dbUser models.User
		ur.DB.First(&dbUser, "ID = ?", user.ID)

		// Sets given user values
		is.True(dbUser.ID != uuid.UUID{})
		is.Equal(dbUser.Email, "testRegisterUser@test.com")
		is.Equal(dbUser.Password, "password")
		// Sets default values
		is.Equal(dbUser.LastLogin, nil)
		is.Equal(dbUser.FailedLoginAttempts, 0)
		is.True(!dbUser.AccountLocked)
		is.Equal(dbUser.AccountLockedUntil, nil)
	})
}

// TestUserRepository_GetUserByEmail tests lookup of registered users in the database
func TestUserRepository_GetUserByEmail(t *testing.T) {
	is := is.New(t)

	// Error on empty email
	t.Run("empty email", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		dbUser, err := ur.GetUserByEmail("")
		is.Equal(dbUser, nil)
		is.Equal(err, apperrors.ErrEmailIsEmpty)
	})

	// Error on non-existing user lookup
	t.Run("non-existing user", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		dbUser, err := ur.GetUserByEmail("doesNotExist@test.com")
		is.Equal(dbUser, nil)
		is.Equal(err, gorm.ErrRecordNotFound)
	})

	// Success on existing-user lookup
	t.Run("existing user", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		// Register a user to look up
		user := &models.User{
			Email:    "testGetUserByEmail@test.com",
			Password: "password",
		}
		err = ur.RegisterUser(user)
		is.NoErr(err)

		dbUser, err := ur.GetUserByEmail(user.Email)
		is.NoErr(err)

		is.True(dbUser.ID != uuid.UUID{})
		is.Equal(dbUser.Email, "testGetUserByEmail@test.com")
		is.Equal(dbUser.Password, "password")
	})
}

// TestUserRepository_LookupUser tests lookup of registered users in the database
func TestUserRepository_GetUserByID(t *testing.T) {
	is := is.New(t)

	ur, err := setupUserRepository(t)
	is.NoErr(err)

	email := "testGetUserByID@test.com"
	password := "password"

	// Register a user to look up
	user := &models.User{
		Email:    email,
		Password: password,
	}
	err = ur.RegisterUser(user)
	is.NoErr(err)

	// Error on empty user ID
	t.Run("empty user ID", func(t *testing.T) {
		dbUser, err := ur.GetUserByID("")
		is.Equal(dbUser, nil)
		is.Equal(err, apperrors.ErrUserIdEmpty)
	})

	// Error on non-existing user lookup
	t.Run("non-existing user", func(t *testing.T) {
		randUUID := uuid.New()

		dbUser, err := ur.GetUserByID(randUUID.String())
		is.Equal(dbUser, nil)
		is.Equal(err, gorm.ErrRecordNotFound)
	})

	// Success on existing-user lookup
	t.Run("existing user", func(t *testing.T) {
		dbUser, err := ur.GetUserByID(user.ID.String())
		is.NoErr(err)

		is.True(dbUser.ID != uuid.UUID{})
		is.Equal(dbUser.Email, email)
		is.Equal(dbUser.Password, password)
		is.Equal(dbUser.ID, user.ID)
	})
}

// TestUserRepository_PermanentlyDeleteUser tests deletion of existing users in database
func TestUserRepository_PermanentlyDeleteUser(t *testing.T) {
	is := is.New(t)

	ur, err := setupUserRepository(t)
	is.NoErr(err)

	// Error on empty user ID
	t.Run("empty user ID", func(t *testing.T) {
		rowsAffected, err := ur.PermanentlyDeleteUser("")
		is.Equal(rowsAffected, int64(0))
		is.Equal(err, apperrors.ErrUserIdEmpty)
	})

	// Error on non-existing user lookup
	t.Run("non-existing user", func(t *testing.T) {
		randUUID := uuid.New()
		rowsAffected, err := ur.PermanentlyDeleteUser(randUUID.String())
		is.Equal(rowsAffected, int64(0))
		is.NoErr(err)
	})

	// Success on existing-user lookup
	t.Run("existing user", func(t *testing.T) {
		email := "testPermanentlyDeleteUser@test.com"
		password := "password"

		// Register a user to delete
		user := &models.User{
			Email:    email,
			Password: password,
		}
		err = ur.RegisterUser(user)
		is.NoErr(err)
		rowsAffected, err := ur.PermanentlyDeleteUser(user.ID.String())
		is.Equal(rowsAffected, int64(1))
		is.NoErr(err)
		// Confirm user is no longer in the database
		user, err = ur.GetUserByEmail(email)
		is.Equal(user, nil)
		is.Equal(err, gorm.ErrRecordNotFound)
	})
}

// TestUserRepository_UpdateUser test user updates into the database
func TestUserRepository_UpdateUser(t *testing.T) {
	is := is.New(t)

	// Error on empty user ID
	t.Run("empty user ID", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		err = ur.UpdateUser("", map[string]any{})
		is.Equal(err, apperrors.ErrUserIdEmpty)
	})

	// Error on non-existent user
	t.Run("non-existent user", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		err = ur.UpdateUser("doesNotExist", map[string]any{})
		is.Equal(err, apperrors.ErrUserNotFound)
	})

	// Can update user
	t.Run("can update user", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		// Register a user to update
		email := "testUpdateUser@test.com"
		password := "password"
		user := &models.User{
			Email:    email,
			Password: password,
		}
		err = ur.RegisterUser(user)
		is.NoErr(err)

		// Update user
		referenceTime := time.Now().Add(time.Hour * 24).UTC().Truncate(time.Second)
		err = ur.UpdateUser(user.ID.String(), map[string]any{
			"email":                 "newUserName@test.com",
			"password":              "newpassword",
			"last_login":            referenceTime,
			"failed_login_attempts": 99,
			"account_locked":        true,
			"account_locked_until":  referenceTime,
		})
		is.NoErr(err)

		// Get updated user and check for updated fields
		user, err = ur.GetUserByID(user.ID.String())
		is.NoErr(err)
		t.Run("updates email", func(t *testing.T) {
			is.Equal(user.Email, "newUserName@test.com")
		})
		t.Run("updates password", func(t *testing.T) {
			is.Equal(user.Password, "newpassword")
		})
		t.Run("updates last_login", func(t *testing.T) {
			is.Equal(user.LastLogin, &referenceTime)
		})
		t.Run("updates failed_login_attempts", func(t *testing.T) {
			is.Equal(user.FailedLoginAttempts, 99)
		})
		t.Run("updates account_locked", func(t *testing.T) {
			is.True(user.AccountLocked)
		})
		t.Run("updates account_locked_until", func(t *testing.T) {
			is.Equal(user.LastLogin, &referenceTime)
		})
	})
}

func TestUserRepository_IncrementFailedLogins(t *testing.T) {
	is := is.New(t)

	// Error on empty user ID
	t.Run("empty user ID", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		err = ur.IncrementFailedLogins("")
		is.Equal(err, apperrors.ErrUserIdEmpty)
	})

	// Error on non-existing user lookup
	t.Run("fails on non-existent user", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		err = ur.IncrementFailedLogins(uuid.New().String())
		is.Equal(err, apperrors.ErrUserNotFound)
	})

	// Increments failed login attempts on successful lookup
	t.Run("increments FailedLoginAttempts", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		// Register a user to fail logins with
		email := "testHandleFailedLogin@test.com"
		password := "password"
		user := &models.User{
			Email:    email,
			Password: password,
		}
		err = ur.RegisterUser(user)
		is.NoErr(err)

		// Increment login attempts
		for i := range 10 {
			err = ur.IncrementFailedLogins(user.ID.String())
			is.NoErr(err)

			user, err = ur.GetUserByEmail(user.Email)
			is.NoErr(err)
			is.Equal(user.FailedLoginAttempts, i+1)
		}
	})
}

func TestUserRepository_LockAccount(t *testing.T) {
	is := is.New(t)

	t.Run("locks on existing user", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		// Register test user
		email := "testLockAccount@test.com"
		password := "password"
		user := &models.User{
			Email:    email,
			Password: password,
		}
		err = ur.RegisterUser(user)
		is.NoErr(err)

		// Lock account
		err = ur.LockAccount(user.ID.String())
		user, err = ur.GetUserByEmail(user.Email)
		is.NoErr(err)

		is.True(user.AccountLocked)
		// HACK: assumes test won't be blocked for 2 seconds, unlock time is > 2 seconds from now
		approxLockMax := time.Now().Add(config.AccountLockoutLength*time.Second - 2)
		is.True(user.AccountLockedUntil.After(approxLockMax))
	})

	t.Run("fails on non-existent user", func(t *testing.T) {
		ur, err := setupUserRepository(t)
		is.NoErr(err)

		err = ur.LockAccount(uuid.New().String())
		is.Equal(err, apperrors.ErrUserNotFound)
	})
}

func setupUserRepository(t *testing.T) (*repository.UserRepository, error) {
	testDB := testutils.TestDBSetup()
	ur, err := repository.NewUserRepository(testDB)
	tx := testDB.Begin()
	t.Cleanup(func() { tx.Rollback() })
	return ur, err
}
