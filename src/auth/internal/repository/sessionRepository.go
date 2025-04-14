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
	return nil
}
