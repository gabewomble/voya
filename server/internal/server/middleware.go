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
			switch {
			case errors.Is(err, data.TokenErr.NotFound):
				c.Next()
			case errors.Is(err, data.TokenErr.InvalidToken):
				s.invalidAuthTokenResponse(c)
			default:
				s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
			}
			return
		}

		tokenHash := data.GetTokenHash(token)

		expiry := time.Now()
		user, err := s.db.Queries().GetUserForToken(c, repository.GetUserForTokenParams{
			TokenHash:   tokenHash[:],
			TokenScope:  data.TokenScope.Authentication,
			TokenExpiry: &expiry,
		})

		if err != nil {
			s.ctxSetUser(c, data.AnonymousUser)
			if errors.Is(err, pgx.ErrNoRows) {
				s.invalidAuthTokenResponse(c)
				c.Abort()
				return
			}
			s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
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
			s.errorResponse(c, http.StatusUnauthorized, errorDetailsFromMessage("you must be authenticated to access this resource"))
			c.Abort()
			return
		}

		c.Next()
	}
}
