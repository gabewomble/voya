package server

import (
	"net/http"
	"server/internal/data"
	"server/internal/repository"
	"server/internal/validator"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

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

	// TODO: Activation email, token / cookie

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

type UserWithoutPassword struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Activated bool      `json:"activated"`
	Version   int32     `json:"version"`
}

func (s *Server) getCurrentUserHandler(c *gin.Context) {
	ctxUser := s.ctxGetUser(c)

	user := &UserWithoutPassword{
		ID:        ctxUser.ID,
		CreatedAt: ctxUser.CreatedAt,
		Name:      ctxUser.Name,
		Email:     ctxUser.Email,
		Activated: ctxUser.Activated,
		Version:   ctxUser.Version,
	}

	c.JSON(http.StatusFound, gin.H{"user": user})
}
