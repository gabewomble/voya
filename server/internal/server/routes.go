package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())
	r.Use(s.authenticate())

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)
	r.POST("/users", s.registerUserHandler)
	r.POST("/tokens/authenticate", s.createAuthTokenHandler)

	protected := r.Group("/")
	protected.Use(s.requireAuthenticatedUser())

	// Users
	protected.GET("/users/current", s.getCurrentUserHandler)
	protected.GET("/users/:id", s.getUserByIdHandler)

	// Trips
	protected.GET("/trips", s.listTripsHandler)
	protected.POST("/trips", s.createTripHandler)
	protected.GET("/trips/:id", s.getTripByIdHandler)
	protected.DELETE("/trips/:id", s.deleteTripByIdHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	user := s.ctxGetUser(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, world!",
		"user":    user,
	})
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
