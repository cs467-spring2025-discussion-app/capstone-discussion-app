package services

import (
	"godiscauth/internal/models"
	"godiscauth/internal/repository"
	"godiscauth/pkg/apperrors"
)

// UserService is a struct that contains the repositories needed for user-related operations
type UserService struct {
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
}

// NewUserService returns a value of type UserService
func NewUserService(ur *repository.UserRepository, sr *repository.SessionRepository) (*UserService, error) {
	if ur == nil {
		return nil, apperrors.ErrUserRepoIsNil
	}
	if sr == nil {
		return nil, apperrors.ErrSessionRepoIsNil
	}
	return &UserService{
		UserRepo:    ur,
		SessionRepo: sr,
	}, nil
}

// RegisterUser mediates the new User value creation and the insertion of a user into the database
func (us *UserService) RegisterUser(email, password string) error {
	// Check for empty fields
	if email == "" {
		return apperrors.ErrEmailIsEmpty
	}
	if password == "" {
		return apperrors.ErrPasswordIsEmpty
	}

	user, err := models.NewUser(email, password)
	if err != nil {
		return err
	}
	return us.UserRepo.RegisterUser(user)
}
