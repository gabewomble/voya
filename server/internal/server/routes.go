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
	r.GET("/users/current", s.getCurrentUserHandler)
	r.PUT("/users/activated", s.activateUserHandler)
	r.POST("/users/resend-activation", s.resendActivationHandler)
	r.POST("/tokens/authenticate", s.createAuthTokenHandler)
	r.POST("/tokens/refresh", s.refreshAuthTokenHandler)

	protected := r.Group("/")
	protected.Use(s.requireAuthenticatedUser())
	{
		// Tokens
		protected.DELETE("/tokens/current", s.deleteAuthTokenHandler)
		// Users
		protected.GET("/users/u/:username", s.getUserByUsernameHandler)
		protected.PATCH("/users/u/:username", s.updateUserProfileHandler)
		protected.POST("/users/search", s.searchUsersHandler)

		// Trips
		protected.GET("/trips", s.listTripsHandler)
		protected.POST("/trips", s.createTripHandler)
		protected.GET("/trips/t/:id", s.getTripByIdHandler)
		protected.DELETE("/trips/t/:id", s.deleteTripByIdHandler)
	}

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
