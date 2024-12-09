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
		// Hate that I can't have both /users/search and /users/:username
		// -1 for Gin
		protected.GET("/user/:username", s.getUserByUsernameHandler)
		protected.PATCH("/user/:username", s.updateUserProfileHandler)
		protected.POST("/users/search", s.searchUsersHandler)
		protected.POST("/users/batch", s.getUsersByIdsHandler)

		// Trips
		protected.GET("/trips", s.listTripsHandler)
		protected.POST("/trips", s.createTripHandler)
		protected.GET("/trip/:id", s.getTripByIdHandler)
		protected.DELETE("/trip/:id", s.deleteTripByIdHandler)

		// Gin doesn't correctly parse parameters, so it considers /memberes as part of the :id
		// -1 for Gin
		// Trip Members
		protected.POST("/trip/:id/members", s.addMemberToTripHandler)
		protected.PATCH("/trip/:id/members", s.updateTripMemberStatusHandler)

		// Notifications
		protected.GET("/notification/:id", s.getNotificationByIdHandler)
		protected.POST("/notification/:id/read", s.markNotificationAsReadHandler)
		protected.DELETE("/notification/:id", s.deleteNotificationHandler)
		protected.GET("/notifications", s.listNotificationsHandler)
		protected.GET("/notifications/unread", s.listUnreadNotificationsHandler)
		protected.GET("/notifications/unread/count", s.countUnreadNotificationsHandler)
		protected.POST("/notifications/read", s.markNotificationsAsReadHandler)
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
