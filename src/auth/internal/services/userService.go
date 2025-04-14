package services

import (
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
