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
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInput := &data.UserInput{
		Name:  input.Name,
		Email: input.Email,
	}

	userInput.Password.Set(input.Password)

	v := validator.New()

	if userInput.Validate(v); !v.Valid() {
		s.logger.LogInfo(c, "invalid user input")
		for key, value := range v.Errors {
			s.logger.LogInfo(c, fmt.Sprintf("%s: %s", key, value))
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
		return
	}

	// Insert user
	insertUserParams := repository.InsertUserParams{
		Name:         userInput.Name,
		Email:        userInput.Email,
		PasswordHash: userInput.Password.Hash,
		Activated:    false,
	}

	user, err := s.db.Queries().InsertUser(c, insertUserParams)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			v.AddError("email", "duplicate email")
			c.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create activation token
	token, err := data.Token.New(data.Token{}, user.ID, 3*24*time.Hour, data.TokenScope.Activation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = s.db.Queries().InsertToken(c, repository.InsertTokenParams(token.Model))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			s.logger.LogError(c, err)
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
}

func sanitizeUser(u *repository.User) cleanUser {
	return cleanUser{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		Name:      u.Name,
		Email:     u.Email,
		Activated: u.Activated,
		Version:   u.Version}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := s.db.Queries().GetUserById(c, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.notFoundResponse(c, "user not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate activation token
	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.Token); !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
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
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Delete activation tokens
	err = s.db.Queries().DeleteAllTokensForUser(c, repository.DeleteAllTokensForUserParams{
		TokenScope: data.TokenScope.Activation,
		UserID:     user.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Create authentication token
	token, err := data.Token.New(data.Token{}, user.ID, 24*time.Hour, data.TokenScope.Authentication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Insert authentication token
	err = s.db.Queries().InsertToken(c, repository.InsertTokenParams(token.Model))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token.Plaintext})
}
