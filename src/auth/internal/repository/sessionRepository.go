package repository

import (
	"time"

	"gorm.io/gorm"

	"godiscauth/internal/models"
	"godiscauth/pkg/apperrors"
)

// SessionRepository represents the entry point into the database for managing
// the `sessions` table
type SessionRepository struct {
	DB *gorm.DB
}

// NewSessionRepository returns a value for the SessionRepository struct
func NewSessionRepository(db *gorm.DB) (*SessionRepository, error) {
	if db == nil {
		return nil, apperrors.ErrDatabaseIsNil
	}
	return &SessionRepository{DB: db}, nil
}

// CreateSession inserts a new session into the `sessions` table
func (sr *SessionRepository) CreateSession(session *models.Session) error {
	if session == nil {
		return apperrors.ErrSessionIsNil
	}
	// Lookup existing session by token
	var existingSession models.Session
	result := sr.DB.Where("token = ?", session.Token).First(&existingSession)
	if result.Error == nil {
		// Session already exists
		return apperrors.ErrSessionAlreadyExists
	} else if result.Error != gorm.ErrRecordNotFound { // RecordNotFound is what we want
		// Some other error
		return result.Error
	}
	return sr.DB.Create(session).Error
}

// GetUnexpiredSessionByToken retrieves a session from the database by token, but ignores any
// expired tokens
func (sr *SessionRepository) GetUnexpiredSessionByToken(token string) (*models.Session, error) {
	if token == "" {
		return nil, apperrors.ErrTokenIsEmpty
	}
	var session models.Session
	result := sr.DB.Where("token = ? AND expires_at > ?", token, time.Now()).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}


// DeleteSessionByToken deletes a single session from the database by token
func (sr *SessionRepository) DeleteSessionByToken(token string) error {
	if token == "" {
		return apperrors.ErrTokenIsEmpty
	}
	result := sr.DB.Where("token = ?", token).Delete(&models.Session{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
