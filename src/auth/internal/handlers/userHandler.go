package handlers

import (
	"godiscauth/internal/services"
	"godiscauth/pkg/apperrors"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) (*UserHandler, error) {
	if userService == nil {
		return nil, apperrors.ErrUserServiceIsNil
	}
	return &UserHandler{UserService: userService}, nil
}
