package server

import (
	"errors"
	"net/http"
	"server/internal/data"
	"server/internal/repository"
	"server/internal/validator"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (s *Server) createAuthTokenHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v := validator.New()

	data.ValidateUserEmail(v, input.Email)
	data.ValidateUserPasswordPlaintext(v, input.Password)

	if !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
		return
	}

	user, err := s.db.Queries().GetUserByEmail(c, input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.invalidCredentialsResponse(c)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	match, err := data.CompareUserHashAndPassword(user.PasswordHash, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !match {
		s.invalidCredentialsResponse(c)
		return
	}

	token, err := data.Token.New(data.Token{}, user.ID, 24*time.Hour, data.TokenScope.Authentication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.db.Queries().InsertToken(c, repository.InsertTokenParams(token.Model))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token.Plaintext, "scope": data.TokenScope.Authentication})
}

func (s *Server) deleteAuthTokenHandler(c *gin.Context) {
	token, err := data.GetTokenPlainTextFromContext(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	tokenHash := data.GetTokenHash(token)

	err = s.db.Queries().DeleteToken(c, tokenHash[:])

	if err != nil {
		s.logger.LogError(c, err)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "operation not permitted"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete token"})
		return
	}

	c.Status(http.StatusNoContent)
}
