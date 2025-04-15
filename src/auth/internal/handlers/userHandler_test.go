package handlers_test

import (
	"testing"

	"github.com/matryer/is"

	"godiscauth/internal/handlers"
	"godiscauth/internal/repository"
	"godiscauth/internal/services"
	"godiscauth/internal/testutils"
	"godiscauth/pkg/apperrors"
)

// TestHandlers_NewUserHandler checks the NewUserHandler constructor returns a valid UsrHandler
func TestHandlers_NewUserHandler(t *testing.T) {
	is := is.New(t)

	t.Run("err on nil user service", func(t *testing.T) {
		uh, err := handlers.NewUserHandler(nil)
		is.Equal(uh, nil)
		is.Equal(err, apperrors.ErrUserServiceIsNil)
	})

	t.Run("creates new user handler", func(t *testing.T) {
		uh := setupUserHandler(t)
		is.True(uh != nil)
	})
}

func setupUserHandler(t *testing.T) *handlers.UserHandler {
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
	uh, err := handlers.NewUserHandler(us)
	if err != nil {
		t.Fatalf("failed to create user handler: %v", err)
	}
	return uh
}
