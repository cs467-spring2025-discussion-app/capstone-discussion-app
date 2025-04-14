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

func TestSessionRepository_CreateSession(t *testing.T) {
	testDB := testutils.TestDBSetup()
	is := is.New(t)

	t.Run("fails on nil session", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)

		err = sr.CreateSession(nil)
		is.Equal(err, apperrors.ErrSessionIsNil)
	})

	t.Run("creates session", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)

		// Register test user
		user := &models.User{
			Email:    "testCreatesSession@test.com",
			Password: "password",
		}
		err = tx.Create(user).Error
		is.NoErr(err)

		session, err := models.NewSession(user.ID, uuid.New().String(), time.Now().Add(1*time.Hour))
		is.NoErr(err)

		err = sr.CreateSession(session)
		is.NoErr(err)
	})

	t.Run("fails on duplicate session", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)

		// Register test user
		user := &models.User{
			Email:    "testCreatesSession@test.com",
			Password: "password",
		}
		err = tx.Create(user).Error
		is.NoErr(err)

		tokenStr := uuid.New().String()

		// Insert first session
		sessionOne, err := models.NewSession(user.ID, tokenStr, time.Now().Add(1*time.Hour))
		is.NoErr(err)

		err = sr.CreateSession(sessionOne)
		is.NoErr(err)

		// Insert second session with same token (expect error)
		sessionTwo, err := models.NewSession(
			user.ID,
			tokenStr,
			time.Now().Add(1*time.Hour),
		)
		is.NoErr(err)

		err = sr.CreateSession(sessionTwo)
		is.Equal(err, apperrors.ErrSessionAlreadyExists)
	})
}

func TestSessionRepository_GetSessionByToken(t *testing.T) {
	testDB := testutils.TestDBSetup()
	is := is.New(t)

	t.Run("fails on empty token", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)

		session, err := sr.GetUnexpiredSessionByToken("")
		is.Equal(session, nil)
		is.Equal(err, apperrors.ErrTokenIsEmpty)
	})

	t.Run("fails on non-existing token", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)

		session, err := sr.GetUnexpiredSessionByToken(uuid.New().String())
		is.Equal(session, nil)
		is.Equal(err, gorm.ErrRecordNotFound)
	})

	// Valid session in db can be retrieved
	t.Run("retrieves session by token", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)

		// Register test user
		user := &models.User{
			Email:    "testGetUnexpiredSessionByToken@test.com",
			Password: "password",
		}
		err = tx.Create(user).Error
		is.NoErr(err)

		// Insert associated session
		session, err := models.NewSession(user.ID, uuid.New().String(), time.Now().Add(1*time.Hour))
		is.NoErr(err)
		err = sr.CreateSession(session)
		is.NoErr(err)

		// Retrieved session has same values we passed in
		retrievedSession, err := sr.GetUnexpiredSessionByToken(session.Token)
		is.NoErr(err)
		is.Equal(retrievedSession.UserID, session.UserID)
		is.Equal(retrievedSession.Token, session.Token)
	})
}

func TestSessionRepository_DeleteSessionByToken(t *testing.T) {
	testDB := testutils.TestDBSetup()
	is := is.New(t)

	t.Run("fails on empty token", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)

		err = sr.DeleteSessionByToken("")
		is.Equal(err, apperrors.ErrTokenIsEmpty)
	})

	t.Run("fails on non-existing token", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)

		err = sr.DeleteSessionByToken(uuid.New().String())
		is.Equal(err, gorm.ErrRecordNotFound)
	})

	t.Run("deletes session by token", func(t *testing.T) {
		tx := testDB.Begin()
		defer tx.Rollback()
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)

		// Register test user
		user := &models.User{
			Email:    "testDeleteSessionByToken@test.com",
			Password: "password",
		}
		err = tx.Create(user).Error
		is.NoErr(err)

		// Insert session to delete
		session, err := models.NewSession(user.ID, uuid.New().String(), time.Now().Add(1*time.Hour))
		is.NoErr(err)
		err = sr.CreateSession(session)
		is.NoErr(err)

		// Delete session
		err = sr.DeleteSessionByToken(session.Token)
		is.NoErr(err)

		// Attempt to retrieve deleted session
		retrievedSession, err := sr.GetUnexpiredSessionByToken(session.Token)
		is.Equal(retrievedSession, nil)
		is.Equal(err, gorm.ErrRecordNotFound)
	})
}

func TestSessionRepository_DeleteSessionsByUserID(t *testing.T) {
	testDB := testutils.TestDBSetup()
	is := is.New(t)
	tx := testDB.Begin()
	defer tx.Rollback()

	t.Run("fails on non-existing user ID", func(t *testing.T) {
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)
		err = sr.DeleteSessionsByUserID(uuid.New().String())
		is.Equal(err, gorm.ErrRecordNotFound)
	})

	t.Run("fails on empty user ID", func(t *testing.T) {
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)
		err = sr.DeleteSessionsByUserID("")
		is.Equal(err, apperrors.ErrUserIdEmpty)
	})

	t.Run("deletes session by user ID", func(t *testing.T) {
		sr, err := repository.NewSessionRepository(tx)
		is.NoErr(err)

		// Register test user
		user := &models.User{
			Email:    "testDeleteSessionsByUserID@test.com",
			Password: "password",
		}
		err = tx.Create(user).Error
		is.NoErr(err)
		is.NoErr(err)

		// Insert first session associated with user
		sessionOne, err := models.NewSession(user.ID, uuid.New().String(), time.Now().Add(1*time.Hour))
		is.NoErr(err)
		err = sr.CreateSession(sessionOne)
		is.NoErr(err)

		// Insert second session associated with user
		sessionTwo, err := models.NewSession(user.ID, uuid.New().String(), time.Now().Add(1*time.Hour))
		is.NoErr(err)
		err = sr.CreateSession(sessionTwo)
		is.NoErr(err)

		// Delete both sessions
		err = sr.DeleteSessionsByUserID(user.ID.String())
		is.NoErr(err)

		// Attempt to retrieve both sessions
		for _, session := range []*models.Session{sessionOne, sessionTwo} {
			retrievedSession, err := sr.GetUnexpiredSessionByToken(session.Token)
			if session.UserID == user.ID {
				is.Equal(retrievedSession, nil)
				is.Equal(err, gorm.ErrRecordNotFound)
			} else {
				is.NoErr(err)
				is.Equal(retrievedSession.UserID, user.ID)
			}
		}
	})
}
