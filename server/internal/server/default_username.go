package server

import (
	"context"
	"fmt"
	"log"
	"server/internal/repository"
	"strings"
)

func (s *Server) generateDefaultUsername(email string) string {
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return email
	}

	return email[:atIndex]
}

func (s *Server) getUniqueUsername(ctx context.Context, baseUserName string) (string, error) {
	username := baseUserName

	for i := 1; ; i++ {
		exists, err := s.db.Queries().CheckUsernameExists(ctx, baseUserName)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
		username = fmt.Sprintf("%s%d", baseUserName, i)
	}

	return username, nil
}

func (s *Server) generateDefaultUsernamesForUsers(ctx context.Context) {
	users, err := s.db.Queries().GetUsersWithoutUsername(ctx)
	if err != nil {
		log.Fatalf("Failed to get users without username: %v", err)
	}

	for _, user := range users {
		defaultUsername := s.generateDefaultUsername(user.Email)
		uniqueUsername, err := s.getUniqueUsername(ctx, defaultUsername)
		if err != nil {
			log.Printf("Failed to get unique username for user %s: %v", user.Email, err)
			continue
		}

		err = s.db.Queries().UpdateUsername(ctx, repository.UpdateUsernameParams{
			Username: uniqueUsername,
			ID:       user.ID,
		})
		if err != nil {
			log.Printf("Failed to update username for user %s: %v", user.Email, err)
		} else {
			log.Printf("Updated username for user %s to %s", user.Email, uniqueUsername)
		}
	}
}
