package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

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

func (uh *UserHandler) RegisterUser(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	clientIP := c.ClientIP()

	// Expect both email and password
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Info().
			Str("email", body.Email).
			Str("clientIP", clientIP).
			Str("error", err.Error()).
			Msg("Bad user registration request")

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Attempt registration
	if err := uh.UserService.RegisterUser(body.Email, body.Password); err != nil {
		log.Info().
			Str("email", body.Email).
			Str("clientIP", clientIP).
			Str("error", err.Error()).
			Msg("User registration failed")


		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Registration success
	log.Info().
		Str("email", body.Email).
		Str("clientIP", clientIP).
		Msg("User registration success")

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s created", body.Email)})
}
