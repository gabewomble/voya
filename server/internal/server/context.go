package server

import (
	"server/internal/repository"

	"github.com/gin-gonic/gin"
)

type ctxKey string

const userCtxKey = ctxKey("user")

func (s *Server) ctxSetUser(c *gin.Context, user *repository.User) {
	c.Set(string(userCtxKey), user)
}

func (s *Server) ctxGetUser(c *gin.Context) *repository.User {
	user, ok := c.Get(string(userCtxKey))
	if !ok {
		panic("missing user value in request context")
	}

	return user.(*repository.User)
}
