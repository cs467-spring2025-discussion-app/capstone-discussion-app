package repository

import (
	"gorm.io/gorm"

	"godiscauth/internal/models"
	"godiscauth/pkg/apperrors"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) (*UserRepository, error) {
	if db == nil {
		return nil, apperrors.ErrDatabaseIsNil
	}
	return &UserRepository{DB: db}, nil
}

func (ur *UserRepository) RegisterUser(u *models.User) error {
	if u == nil {
		return apperrors.ErrUserIsNil
	}
	return nil
}
