package repository

import (
	"strings"

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
	// Validate user
	if u == nil {
		return apperrors.ErrUserIsNil
	}
	if u.Email == "" {
		return apperrors.ErrEmailIsEmpty
	}
	if u.Password == "" {
		return apperrors.ErrPasswordIsEmpty
	}

	err := ur.DB.Create(u).Error
	if err != nil && strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_pkey"`) {
		return apperrors.ErrDuplicateEmail
	}
	return err
}

func (r *UserRepository) LookupUser(email string) (*models.User, error) {
	return nil, nil
}
