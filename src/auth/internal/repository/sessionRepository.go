package repository

import (
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
