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
	Database *gorm.DB
	Router   *gin.Engine
	Handlers *HandlerRegistry
	// TODO: middleware
}

// NewAPIServer initializes a new API server with the gin engine as the router.
func NewAPIServer(db *gorm.DB) (*APIServer, error) {
	if db == nil {
		return nil, apperrors.ErrDatabaseIsNil
	}

	repos, err := initRepoProvider(db)
	if err != nil {
		return nil, apperrors.ErrCouldNotInitRepoProvider
	}
	services, err := initServiceProvider(repos)
	if err != nil {
		return nil, apperrors.ErrCouldNotInitServiceProvider
	}
	handlers, err := initHandlerRegistry(services)
	if err != nil {
		return nil, apperrors.ErrCouldNotInitHandlerRegistry
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetTrustedProxies([]string{"127.0.0.1"})

	server := &APIServer{
		Database: db,
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

func initRepoProvider(db *gorm.DB) (*RepoProvider, error) {
	if db == nil {
		return nil, apperrors.ErrDatabaseIsNil
	}
	ur, _ := repository.NewUserRepository(db)
	if ur == nil {
		return nil, apperrors.ErrUserRepoIsNil
	}
	sr, _ := repository.NewSessionRepository(db)
	if sr == nil {
		return nil, apperrors.ErrSessionRepoIsNil
	}
	return &RepoProvider{
		User:    ur,
		Session: sr,
	}, nil
}

func initServiceProvider(repos *RepoProvider) (*Services, error) {
	if repos == nil {
		return nil, apperrors.ErrRepoProviderIsNil
	}
	us, _ := services.NewUserService(repos.User, repos.Session)
	if us == nil {
		return nil, apperrors.ErrUserServiceIsNil
	}
	return &Services{
		User: us,
	}, nil
}

func initHandlerRegistry(services *Services) (*HandlerRegistry, error) {
	uh, _ := handlers.NewUserHandler(services.User)
	if uh == nil {
		return nil, apperrors.ErrUserHandlerIsNil
	}
	return &HandlerRegistry{
		User: uh,
	}, nil
}

type RepoProvider struct {
	User    *repository.UserRepository
	Session *repository.SessionRepository
}

type Services struct {
	User *services.UserService
}

type HandlerRegistry struct {
	User *handlers.UserHandler
}
