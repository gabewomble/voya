package server

import (
	"errors"
	"fmt"
	"net/http"
	"server/internal/data"
	"server/internal/repository"
	"server/internal/validator"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (s *Server) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Vary", "Authorization")

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			s.ctxSetUser(c, data.AnonymousUser)
			c.Next()
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			s.logger.LogInfo(c, "aborted auth at header parsing")
			s.invalidAuthTokenResponse(c)
			c.Abort()
			return
		}

		token := headerParts[1]
		s.logger.LogInfo(c, fmt.Sprintf("token: %s", token))

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			s.invalidAuthTokenResponse(c)
			c.Abort()
			return
		}

		tokenHash := data.GetTokenHash(token)
		s.logger.LogInfo(c, fmt.Sprintf("hash: %s", tokenHash))

		user, err := s.db.Queries().GetUserForToken(c, repository.GetUserForTokenParams{
			TokenHash:   tokenHash[:],
			TokenScope:  data.ScopeAuthentication,
			TokenExpiry: time.Now(),
		})

		if err != nil {
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