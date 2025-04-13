package repository_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/matryer/is"

	"godiscauth/internal/models"
	"godiscauth/internal/repository"
	"godiscauth/internal/testutils"
	"godiscauth/pkg/apperrors"
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
	testDB := testutils.TestDBSetup()
	is := is.New(t)
	tx := testDB.Begin()
	defer tx.Rollback()

	ur, err := repository.NewUserRepository(tx)
	is.NoErr(err)

	// Validates user before registration
	t.Run("fails on nil user", func(t *testing.T) {
		err := ur.RegisterUser(nil)
		is.Equal(err, apperrors.ErrUserIsNil)
	})
	t.Run("fails on missing email", func(t *testing.T) {
		user := &models.User{
			Password: "password",
		}
		err := ur.RegisterUser(user)
		is.Equal(err, apperrors.ErrEmailIsEmpty)
	})
	t.Run("fails on missing password", func(t *testing.T) {
		user := &models.User{
			Email: "testRegisterUser@test.com",
		}
		err := ur.RegisterUser(user)
		is.Equal(err, apperrors.ErrPasswordIsEmpty)
	})

	// Inserts user into db with complete User value
	t.Run("registers user", func(t *testing.T) {
		user := &models.User{
			Email:    "testRegisterUser@test.com",
			Password: "password",
		}
		err := ur.RegisterUser(user)
		is.NoErr(err)

		// Lookup expected user in the db
		var dbUser models.User
		tx.First(&dbUser, "ID = ?", user.ID)

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
	testDB := testutils.TestDBSetup()
	is := is.New(t)
	tx := testDB.Begin()
	defer tx.Rollback()

	ur, err := repository.NewUserRepository(tx)
	is.NoErr(err)

	// Register a user to look up
	user := &models.User{
		Email:    "testLookupUser@test.com",
		Password: "password",
	}
	err = ur.RegisterUser(user)
	is.NoErr(err)

	// Error on empty email
	t.Run("empty email", func(t *testing.T) {
		dbUser, err := ur.GetUserByEmail("")
		is.Equal(dbUser, nil)
		is.Equal(err, apperrors.ErrEmailIsEmpty)
	})

	// Error on non-existing user lookup
	t.Run("non-existing user", func(t *testing.T) {
		dbUser, err := ur.GetUserByEmail("doesNotExist@test.com")
		is.Equal(dbUser, nil)
		is.Equal(err, apperrors.ErrUserNotFound)
	})

	// Success on existing-user lookup
	t.Run("existing user", func(t *testing.T) {
		dbUser, err := ur.GetUserByEmail(user.Email)
		is.NoErr(err)

		is.True(dbUser.ID != uuid.UUID{})
		is.Equal(dbUser.Email, "testLookupUser@test.com")
		is.Equal(dbUser.Password, "password")
	})
}

// TestUserRepository_LookupUser tests lookup of registered users in the database
func TestUserRepository_GetUserByID(t *testing.T) {
	testDB := testutils.TestDBSetup()
	is := is.New(t)
	tx := testDB.Begin()
	defer tx.Rollback()

	ur, err := repository.NewUserRepository(tx)
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
		is.Equal(err, apperrors.ErrUserNotFound)
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
	testDB := testutils.TestDBSetup()
	is := is.New(t)
	tx := testDB.Begin()
	defer tx.Rollback()

	ur, err := repository.NewUserRepository(tx)
	is.NoErr(err)

	email := "testPermanentlyDeleteUser@test.com"
	password := "password"

	// Register a user to delete
	user := &models.User{
		Email:    email,
		Password: password,
	}
	err = ur.RegisterUser(user)
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
		rowsAffected, err := ur.PermanentlyDeleteUser(user.ID.String())
		is.Equal(rowsAffected, int64(1))
		is.NoErr(err)
	})
}

func TestUserRepository_IncrementFailedLogins(t *testing.T) {
	is := is.New(t)
	testDB := testutils.TestDBSetup()

	// Error on empty user ID
	t.Run("empty user ID", func(t *testing.T) {
		is := is.New(t)
		tx := testDB.Begin()
		defer tx.Rollback()
		ur, err := repository.NewUserRepository(tx)
		is.NoErr(err)

		err = ur.IncrementFailedLogins("")
		is.Equal(err, apperrors.ErrUserIdEmpty)
	})

	// Error on non-existing user lookup
	t.Run("fails on non-existent user", func(t *testing.T) {
		is := is.New(t)
		tx := testDB.Begin()
		defer tx.Rollback()
		ur, err := repository.NewUserRepository(tx)
		is.NoErr(err)

		err = ur.IncrementFailedLogins(uuid.New().String())
		is.Equal(err, apperrors.ErrUserNotFound)
	})

	// Increments failed login attempts on successful lookup
	t.Run("increments FailedLoginAttempts", func(t *testing.T) {
		is := is.New(t)
		tx := testDB.Begin()
		defer tx.Rollback()
		ur, err := repository.NewUserRepository(tx)
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

		for i := range 10 {
			// Does not err on non-empty string
			err = ur.IncrementFailedLogins(user.ID.String())
			is.NoErr(err)

			// Increments on each call
			user, err = ur.GetUserByEmail(user.Email)
			is.NoErr(err)
			is.Equal(user.FailedLoginAttempts, i+1)
		}
	})
}
