package models_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/matryer/is"

	"godiscauth/internal/models"
	"godiscauth/pkg/apperrors"
)

func TestSessionModel_NewSession(t *testing.T) {
	is := is.New(t)

	// Valid uuid and non-empty token should return non-nil Session value, nil
	// error
	t.Run("new valid session", func(t *testing.T) {
		session, err := models.NewSession(uuid.New(), "test-token", time.Now().Add(24*time.Hour))
		is.True(session != nil)
		is.NoErr(err)
	})

	t.Run("fails when user ID is empty", func(t *testing.T) {
		_, err := models.NewSession(uuid.Nil, "test-token", time.Now().Add(24*time.Hour))
		is.Equal(err, apperrors.ErrUserIdEmpty)
	})

	t.Run("fails when token is empty", func(t *testing.T) {
		_, err := models.NewSession(uuid.New(), "", time.Now().Add(24*time.Hour))
		is.Equal(err, apperrors.ErrTokenIsEmpty)
	})
	t.Run("fails when expiration time is empty", func(t *testing.T) {
		_, err := models.NewSession(uuid.New(), "test-token", time.Time{})
		is.Equal(err, apperrors.ErrExpiresAtIsEmpty)
	})
}
