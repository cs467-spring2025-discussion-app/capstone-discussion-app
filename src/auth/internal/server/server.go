package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"godiscauth/internal/handlers"
	"godiscauth/internal/middleware"
	"godiscauth/internal/repository"
	"godiscauth/internal/services"
	"godiscauth/pkg/apperrors"
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
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	r.POST("/register", s.HandlerRegistry.User.RegisterUser)
	r.POST("/login", s.HandlerRegistry.User.Login)
	r.POST("/logout", s.HandlerRegistry.User.Logout)

	protected := r.Group("")
	protected.Use(s.MiddlewareProvider.Auth.RequireAuth())
	{
		protected.POST("/logouteverywhere", s.HandlerRegistry.User.LogoutEverywhere)
		protected.POST("/updateuser", s.HandlerRegistry.User.UpdateUser)
		protected.DELETE("/deleteaccount", s.HandlerRegistry.User.PermanentlyDeleteUser)
	}
}

// Run starts the API server and listens for incoming requests.
func (s *APIServer) Run() {
	s.SetupRoutes()
	s.Router.Run()
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
