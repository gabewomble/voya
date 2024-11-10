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

	// Trips
	r.GET("/trips", s.listTripsHandler)
	r.POST("/trips", s.createTripHandler)
	r.GET("/trips/:id", s.getTripByIdHandler)
	r.DELETE("/trips/:id", s.deleteTripByIdHandler)

	r.POST("/users", s.registerUserHandler)

	r.POST("/tokens/authenticate", s.createAuthTokenHandler)

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
