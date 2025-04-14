package repository_test

import (
	"testing"

	"github.com/matryer/is"

	"godiscauth/internal/repository"
	"godiscauth/internal/testutils"
	"godiscauth/pkg/apperrors"
)

// TestSessionRepository_NewSessionRepository tests creation of SessionRepository
// structs in the `repository` package
func TestSessionRepository_NewSessionRepository(t *testing.T) {
	is := is.New(t)

	testDB := testutils.TestDBSetup()

	t.Run("creates new session repo", func(t *testing.T) {
		sr, err := repository.NewSessionRepository(testDB)
		is.True(sr != nil)
		is.NoErr(err)
	})

	t.Run("returns err with nil db", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()

		sr, err := repository.NewSessionRepository(nil)
		is.Equal(sr, nil)
		is.Equal(err, apperrors.ErrDatabaseIsNil)
	})
}
