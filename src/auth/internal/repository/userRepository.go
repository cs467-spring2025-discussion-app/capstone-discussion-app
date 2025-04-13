package repository

import (
	"gorm.io/gorm"

	"godiscauth/internal/models"
	"godiscauth/pkg/apperrors"
)

// UserRepository represents the entry point into the database for managing
// the `users` table
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository returns a value for the UserRepository struct
func NewUserRepository(db *gorm.DB) (*UserRepository, error) {
	if db == nil {
		return nil, apperrors.ErrDatabaseIsNil
	}
	return &UserRepository{DB: db}, nil
}

// RegisterUser inserts a new user into the `users` table
func (ur *UserRepository) RegisterUser(u *models.User) error {
	if u == nil {
		return apperrors.ErrUserIsNil
	}
	return nil
}
