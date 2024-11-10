package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) invalidCredentialsResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication credentials"})
}

func (s *Server) invalidAuthTokenResponse(c *gin.Context) {
	c.Header("WWW-Authenticate", "Bearer")
	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing authentication token"})
}

func (s *Server) notFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"error": message})
}
