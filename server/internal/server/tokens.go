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
		s.badRequest(c, errorDetailsFromError(err))
		return
	}

	v := validator.New()

	data.ValidateUserEmail(v, input.Email)
	data.ValidateUserPasswordPlaintext(v, input.Password)

	if !v.Valid() {
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
		return
	}

	user, err := s.db.Queries().GetUserByEmail(c, input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.invalidCredentialsResponse(c)
			return
		}
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	match, err := data.CompareUserHashAndPassword(user.PasswordHash, input.Password)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	if !match {
		s.invalidCredentialsResponse(c)
		return
	}

	token, err := data.Token.New(data.Token{}, user.ID, 24*time.Hour, data.TokenScope.Authentication)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	err = s.db.Queries().InsertToken(c, repository.InsertTokenParams(token.Model))
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token.Plaintext, "scope": data.TokenScope.Authentication})
}

func (s *Server) deleteAuthTokenHandler(c *gin.Context) {
	token, err := data.GetTokenPlainTextFromContext(c)

	if err != nil {
		s.errorResponse(c, http.StatusUnauthorized, errorDetailsFromError(err))
		return
	}

	tokenHash := data.GetTokenHash(token)

	err = s.db.Queries().DeleteToken(c, tokenHash[:])

	if err != nil {
		s.log.LogError(c, "deleteAuthTokenHandler: Failed to delete token", err)
		if errors.Is(err, pgx.ErrNoRows) {
			s.errorResponse(c, http.StatusUnauthorized, errorDetailsFromMessage("operation not allowed"))
			return
		}
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromMessage("failed to delete token"))
		return
	}

	c.Status(http.StatusNoContent)
}
