package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"godiscauth/internal/handlers"
	"godiscauth/internal/repository"
	"godiscauth/internal/services"
	"godiscauth/pkg/apperrors"
)

// APIServer represents the API server with a gin router.
type APIServer struct {
	DB *gorm.DB
	Router   *gin.Engine
	Handlers *HandlerRegistry
	// TODO: middleware
}

// NewAPIServer initializes a new API server with the gin engine as the router.
func NewAPIServer(db *gorm.DB) (*APIServer, error) {
	if db == nil {
		return nil, apperrors.ErrDatabaseIsNil
	}

	repos, err := NewRepoProvider(db)
	if err != nil {
		return nil, err
	}
	services, err := NewServiceProvider(repos)
	if err != nil {
		return nil, err
	}
	handlers, err := NewHandlerRegistry(services)
	if err != nil {
		return nil, err
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetTrustedProxies([]string{"127.0.0.1"})

	server := &APIServer{
		DB: db,
		Router:   router,
		Handlers: handlers,
	}
	return server, nil
}

// SetupRoutes sets up the routes for the API server.
func (s *APIServer) SetupRoutes() {
	r := s.Router
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	r.POST("/register", s.Handlers.User.RegisterUser)

	// TODO: add protected (requireAuth)routes

	// TODO: add admin routes
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
