package server

import (
	"errors"
	"net/http"
	"server/internal/data"
	"server/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (s *Server) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Vary", "Authorization")

		token, err := data.GetTokenPlainTextFromContext(c)

		if err != nil {
			s.ctxSetUser(c, data.AnonymousUser)
			s.logger.LogError(c, err)
			switch {
			case errors.Is(err, data.TokenErr.NotFound):
				c.Next()
			case errors.Is(err, data.TokenErr.InvalidToken):
				s.invalidAuthTokenResponse(c)
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		tokenHash := data.GetTokenHash(token)

		user, err := s.db.Queries().GetUserForToken(c, repository.GetUserForTokenParams{
			TokenHash:   tokenHash[:],
			TokenScope:  data.TokenScope.Authentication,
			TokenExpiry: time.Now(),
		})

		if err != nil {
			s.ctxSetUser(c, data.AnonymousUser)
			if errors.Is(err, pgx.ErrNoRows) {
				s.invalidAuthTokenResponse(c)
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		s.ctxSetUser(c, &user)
		c.Next()
	}
}

func (s *Server) requireAuthenticatedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := s.ctxGetUser(c)

		if data.UserIsAnonymous(user) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you must be authenticated to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
