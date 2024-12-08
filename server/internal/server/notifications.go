package server

import (
	"encoding/json"
	"server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handleNotifyMemberStatusUpdateParams struct {
	TripID       uuid.UUID
	TargetUserID uuid.UUID
	OwnerID      uuid.UUID
	MemberStatus repository.MemberStatusEnum
	Queries      *repository.Queries
}

func (s *Server) handleNotifyMemberStatusUpdate(c *gin.Context, params handleNotifyMemberStatusUpdateParams) error {
	currentUser := s.ctxGetUser(c)

	insertNotificationParams := repository.InsertNotificationParams{
		UserID: params.TargetUserID,
		TripID: params.TripID,
	}

	switch params.MemberStatus {
	case repository.MemberStatusEnumAccepted:
		metadata, err := json.Marshal(map[string]any{
			"user_id":   currentUser.ID,
			"user_name": currentUser.Name,
		})

		if err != nil {
			s.log.LogError(c, "handleNotifyMemberStatusUpdate: json.Marshal failed", err)
			return err
		}

		insertNotificationParams.Type = repository.NotificationTypeTripInviteAccepted
		insertNotificationParams.Message = "You have accepted a trip invite"
		insertNotificationParams.Metadata = metadata

		err = params.Queries.InsertNotification(c, insertNotificationParams)

		if err != nil {
			s.log.LogError(c, "handleNotifyMemberStatusUpdate: InsertNotification failed", err)
			return err
		}
	case repository.MemberStatusEnumDeclined:
		metadata, err := json.Marshal(map[string]any{
			"user_id":   currentUser.ID,
			"user_name": currentUser.Name,
		})

		if err != nil {
			s.log.LogError(c, "handleNotifyMemberStatusUpdate: json.Marshal failed", err)
			return err
		}

		insertNotificationParams.Type = repository.NotificationTypeTripInviteDeclined
		insertNotificationParams.Message = "You have declined a trip invite"
		insertNotificationParams.Metadata = metadata

		err = params.Queries.InsertNotification(c, insertNotificationParams)

		if err != nil {
			s.log.LogError(c, "handleNotifyMemberStatusUpdate: InsertNotification failed", err)
			return err
		}
	}
	return nil
}

type handleNotifyTripInviteParams struct {
	TripID       uuid.UUID
	TargetUserID uuid.UUID
	Queries      *repository.Queries
}

func (s *Server) handleNotifyTripInvite(c *gin.Context, params handleNotifyTripInviteParams) error {
	currentUser := s.ctxGetUser(c)

	metadata, err := json.Marshal(map[string]any{
		"user_id":   currentUser.ID,
		"user_name": currentUser.Name,
	})

	if err != nil {
		s.log.LogError(c, "handleNotifyTripInvite: json.Marshal failed", err)
		return err
	}

	err = params.Queries.InsertNotification(c, repository.InsertNotificationParams{
		UserID:   params.TargetUserID,
		Type:     repository.NotificationTypeTripInvitePending,
		TripID:   params.TripID,
		Message:  "You have been invited to a trip",
		Metadata: metadata,
	})

	if err != nil {
		s.log.LogError(c, "handleNotifyTripInvite: InsertNotification failed", err)
		return err
	}

	return nil
}
