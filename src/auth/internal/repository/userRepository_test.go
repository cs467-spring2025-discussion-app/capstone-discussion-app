package repository_test

import (
	"testing"

	"github.com/matryer/is"

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

	t.Run("fails on nil user", func(t *testing.T) {
		err := ur.RegisterUser(nil)
		is.Equal(err, apperrors.ErrUserIsNil)
	})
}
