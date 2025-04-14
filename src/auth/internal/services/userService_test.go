package services_test

import (
	"testing"

	"github.com/matryer/is"

	"godiscauth/internal/repository"
	"godiscauth/internal/services"
	"godiscauth/internal/testutils"
	"godiscauth/pkg/apperrors"
)

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
