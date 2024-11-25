package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"server/internal/data"
	"server/internal/repository"
	"server/internal/validator"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var appUrl = os.Getenv("APP_URL")

func (s *Server) registerUserHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		s.badRequest(c, errorDetailsFromError(err))
		return
	}

	userInput := &data.UserInput{
		Username: input.Username,
		Name:     input.Name,
		Email:    input.Email,
	}

	userInput.Password.Set(input.Password)

	v := validator.New()

	if userInput.Validate(v); !v.Valid() {
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
		return
	}

	// Insert user
	insertUserParams := repository.InsertUserParams{
		Username:     userInput.Username,
		Name:         userInput.Name,
		Email:        userInput.Email,
		PasswordHash: userInput.Password.Hash,
		Activated:    false,
	}

	s.log.LogInfo(c, "Inserting user", "username", userInput.Username, "email", userInput.Email)

	user, err := s.db.Queries().InsertUser(c, insertUserParams)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_email_key" {
				v.AddError("email", "duplicate email")
			}
			if pgErr.ConstraintName == "idx_users_username" {
				v.AddError("username", "duplicate username")
			}

			s.badRequest(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
			return
		}

		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	// Create activation token
	token, err := data.Token.New(data.Token{}, user.ID, 3*24*time.Hour, data.TokenScope.Activation)
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

	// Send activation email
	s.background(func() {
		activationURL := fmt.Sprintf("%s/activate/%s", appUrl, token.Plaintext)
		data := map[string]any{
			"activationURL": activationURL,
		}

		err = s.mailer.Send(userInput.Email, "user_welcome.tmpl", data)
		if err != nil {
			s.log.LogError(c, "Failed to send activation email", err, "email", userInput.Email)
		}
	})

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

type cleanUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Activated bool      `json:"activated"`
	Version   int32     `json:"version"`
	Username  string    `json:"username"`
}

func sanitizeUser(u *repository.User) cleanUser {
	return cleanUser{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		Name:      u.Name,
		Email:     u.Email,
		Activated: u.Activated,
		Version:   u.Version,
		Username:  u.Username,
	}
}

func (s *Server) getCurrentUserHandler(c *gin.Context) {
	ctxUser := s.ctxGetUser(c)

	if data.UserIsAnonymous(ctxUser) {
		c.JSON(http.StatusOK, gin.H{"user": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": sanitizeUser(ctxUser)})
}

func (s *Server) getUserByIdHandler(c *gin.Context) {
	userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		s.badRequest(c, errorDetailsFromMessage("invalid id"))
		return
	}

	user, err := s.db.Queries().GetUserById(c, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.notFoundResponse(c, errorDetailsFromMessage("user not found"))
			return
		}
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": sanitizeUser(&user)})
}

func (s *Server) activateUserHandler(c *gin.Context) {
	// Parse activation token
	var input struct {
		Token string `json:"token"`
	}
	if err := c.BindJSON(&input); err != nil {
		s.badRequest(c, errorDetailsFromError(err))
		return
	}
	// Validate activation token
	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.Token); !v.Valid() {
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
		return
	}
	// Find user from activation token
	tokenHash := data.GetTokenHash(input.Token)
	user, err := s.db.Queries().GetUserForToken(c, repository.GetUserForTokenParams{
		TokenHash:   tokenHash[:],
		TokenScope:  data.TokenScope.Activation,
		TokenExpiry: time.Now(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v.AddError("token", "invalid or expired activation token")
			s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
			return
		}
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}
	// Update user activation status
	user.Activated = true
	_, err = s.db.Queries().UpdateUser(c, repository.UpdateUserParams{
		Activated:    user.Activated,
		Email:        user.Email,
		ID:           user.ID,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		Version:      user.Version,
		Username:     user.Username,
	})
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}
	// Delete activation tokens
	err = s.db.Queries().DeleteAllTokensForUser(c, repository.DeleteAllTokensForUserParams{
		TokenScope: data.TokenScope.Activation,
		UserID:     user.ID,
	})
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}
	// Create authentication token
	token, err := data.Token.New(data.Token{}, user.ID, 24*time.Hour, data.TokenScope.Authentication)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}
	// Insert authentication token
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
	c.JSON(http.StatusOK, gin.H{"token": token.Plaintext})
}
