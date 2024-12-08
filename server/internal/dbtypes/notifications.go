package dbtypes

import "github.com/google/uuid"

type NotificationMetadata struct {
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
}
