package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// APIServer represents the API server with a gin router.
type APIServer struct {
	Router *gin.Engine
	// TODO: handlers
	// TODO: middleware
}

// NewAPIServer initializes a new API server with the gin engine as the router.
func NewAPIServer(db *gorm.DB) *APIServer {
	server := &APIServer{
		Router: gin.Default(),
	}
	return server
}

// SetupRoutes sets up the routes for the API server.
func (s *APIServer) SetupRoutes() {
	r := s.Router
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	// TODO: add protected (requireAuth)routes

	// TODO: add admin routes
}

// Run starts the API server and listens for incoming requests.
func (s *APIServer) Run() {
	s.SetupRoutes()
	s.Router.Run()
}
