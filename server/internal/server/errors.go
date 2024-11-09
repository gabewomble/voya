package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) invalidCredentialsResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication credentials"})
}
