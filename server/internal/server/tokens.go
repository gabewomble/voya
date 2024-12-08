package server

import (
	"errors"
	"net/http"
	"server/internal/data"
	"server/internal/repository"
	"server/internal/validator"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Server) createAuthTokenHandler(c *gin.Context) {
	var input struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		s.badRequest(c, errorDetailsFromError(err))
		return
	}

	v := validator.New()

	// Validate Identifier
	data.ValidateIdentifier(v, input.Identifier)
	// Validate Password
	data.ValidateUserPasswordPlaintext(v, input.Password)

	if !v.Valid() {
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
		return
	}

	isEmail := validator.Matches(input.Identifier, validator.EmailRX)

	var user repository.User
	var err error

	if isEmail {
		user, err = s.db.Queries().GetUserByEmail(c, input.Identifier)
	} else {
		user, err = s.db.Queries().GetUserByUsername(c, input.Identifier)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.invalidCredentialsResponse(c)
			return
		}
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	if !user.Activated {
		s.errorResponse(c, http.StatusForbidden, errorDetailsFromMessage("user not activated"))
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

	s.generateAndRespondWithToken(c, user.ID)
}

func (s *Server) refreshAuthTokenHandler(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BindJSON(&input); err != nil {
		s.badRequest(c, errorDetailsFromError(err))
		return
	}

	v := validator.New()
	data.ValidateTokenPlaintext(v, input.RefreshToken, "refresh_token")
	if !v.Valid() {
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
		return
	}

	refreshHash := data.GetTokenHash(input.RefreshToken)
	expiry := time.Now()
	user, err := s.db.Queries().GetUserForRefreshToken(c, repository.GetUserForRefreshTokenParams{
		RefreshToken: refreshHash[:],
		TokenScope:   data.TokenScope.Authentication,
		TokenExpiry:  &expiry,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.invalidCredentialsResponse(c)
			return
		}
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	s.generateAndRespondWithToken(c, user.ID)
}

func (s *Server) generateAndRespondWithToken(c *gin.Context, userID uuid.UUID) {
	// Delete the expired tokens for the current user
	err := s.db.Queries().DeleteExpiredTokensForUser(c, userID)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	token, err := data.Token.New(data.Token{}, userID, 24*time.Hour, data.TokenScope.Authentication)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	err = s.db.Queries().InsertToken(c, repository.InsertTokenParams{
		TokenHash:    token.Model.Hash,
		UserID:       token.Model.UserID,
		TokenExpiry:  token.Model.Expiry,
		TokenScope:   token.Model.Scope,
		RefreshToken: token.Model.RefreshToken,
	})
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token.Plaintext, "refresh_token": token.RefreshToken, "scope": data.TokenScope.Authentication})
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
