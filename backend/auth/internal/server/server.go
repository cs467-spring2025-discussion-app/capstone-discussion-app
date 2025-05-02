package server

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"fmt"
	"godiscauth/internal/handlers"
	"godiscauth/internal/middleware"
	"godiscauth/internal/repository"
	"godiscauth/internal/services"
	"godiscauth/pkg/apperrors"
	"godiscauth/pkg/config"

	"github.com/rs/zerolog/log"
)

// APIServer represents the API server with a gin router.
type APIServer struct {
	DB                 *gorm.DB
	Router             *gin.Engine
	HandlerRegistry    *HandlerRegistry
	MiddlewareProvider *MiddlewareProvider
}

// NewAPIServer initializes a new API server with the gin engine as the router.
func NewAPIServer(db *gorm.DB) (*APIServer, error) {
	if db == nil {
		return nil, apperrors.ErrDatabaseIsNil
	}

	repoProvider, err := NewRepoProvider(db)
	if err != nil {
		return nil, err
	}
	serviceProvider, err := NewServiceProvider(repoProvider)
	if err != nil {
		return nil, err
	}
	HandlerRegistry, err := NewHandlerRegistry(serviceProvider)
	if err != nil {
		return nil, err
	}
	middlewareProvider, err := NewMiddlewares(db)
	if err != nil {
		return nil, err
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetTrustedProxies([]string{"127.0.0.1"})

	server := &APIServer{
		DB:                 db,
		Router:             router,
		HandlerRegistry:    HandlerRegistry,
		MiddlewareProvider: middlewareProvider,
	}
	return server, nil
}

// SetupRoutes sets up the routes for the API server.
func (s *APIServer) SetupRoutes() {
	r := s.Router

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://frontend:5173", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", Ping)
	r.POST("/register", s.HandlerRegistry.User.RegisterUser)
	r.POST("/login", s.HandlerRegistry.User.Login)
	r.POST("/logout", s.HandlerRegistry.User.Logout)

	protected := r.Group("")
	protected.Use(s.MiddlewareProvider.Auth.RequireAuth())
	{
		protected.GET("/profile", s.HandlerRegistry.User.GetUserProfile)
		protected.POST("/logouteverywhere", s.HandlerRegistry.User.LogoutEverywhere)
		protected.POST("/updateuser", s.HandlerRegistry.User.UpdateUser)
		protected.DELETE("/deleteaccount", s.HandlerRegistry.User.PermanentlyDeleteUser)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// Run starts the API server and listens for incoming requests.
func (s *APIServer) Run() {
	s.SetupRoutes()
	var port string
	if val, isSet := os.LookupEnv(config.AuthServerPort); isSet {
		port = val
	} else if val, isSet := os.LookupEnv("PORT"); isSet {
		port = val
	}
	if port == "" {
		log.Fatal().Err(fmt.Errorf("Auth server port not set")).Msg("Exiting due to empty port")
	}
	log.Info().Str("port", port).Msg("Starting auth server")
	s.Router.Run(":" + port)
}

func NewRepoProvider(db *gorm.DB) (*RepoProvider, error) {
	if db == nil {
		return nil, apperrors.ErrDatabaseIsNil
	}
	ur, err := repository.NewUserRepository(db)
	if err != nil {
		return nil, err
	}
	sr, err := repository.NewSessionRepository(db)
	if err != nil {
		return nil, err
	}
	return &RepoProvider{
		User:    ur,
		Session: sr,
	}, nil
}

func NewServiceProvider(repos *RepoProvider) (*ServiceProvider, error) {
	if repos == nil {
		return nil, apperrors.ErrRepoProviderIsNil
	}
	us, err := services.NewUserService(repos.User, repos.Session)
	if err != nil {
		return nil, err
	}
	return &ServiceProvider{
		User: us,
	}, nil
}

func NewHandlerRegistry(services *ServiceProvider) (*HandlerRegistry, error) {
	uh, err := handlers.NewUserHandler(services.User)
	if err != nil {
		return nil, err
	}
	return &HandlerRegistry{
		User: uh,
	}, nil
}

func NewMiddlewares(db *gorm.DB) (*MiddlewareProvider, error) {
	mw, err := middleware.NewAuthMiddleware(db)
	if err != nil {
		return nil, err
	}
	return &MiddlewareProvider{
		Auth: mw,
	}, nil
}

type RepoProvider struct {
	User    *repository.UserRepository
	Session *repository.SessionRepository
}

type ServiceProvider struct {
	User *services.UserService
}

type HandlerRegistry struct {
	User *handlers.UserHandler
}

type MiddlewareProvider struct {
	Auth *middleware.AuthMiddleware
}

// Ping godoc
// @Summary ping the server
// @Schemes
// @Description do ping
// @Accept json
// @Produce json
// @Success 200 {object} models.MessageResponse "response with message field"
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
