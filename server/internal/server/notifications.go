package server

import (
	"server/internal/dbtypes"
	"server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) listNotificationsHandler(c *gin.Context) {
	user := s.ctxGetUser(c)
	limit := parseStringToInt32(c.DefaultQuery("limit", "10"), 10)
	offset := parseStringToInt32(c.DefaultQuery("offset", "0"), 0)

	notifications, err := s.db.Queries().ListNotifications(c, repository.ListNotificationsParams{
		UserID:             user.ID,
		NotificationLimit:  limit,
		NotificationOffset: offset,
	})

	if err != nil {
		s.log.LogError(c, "listNotificationsHandler: ListNotifications failed", err)
		c.JSON(500, gin.H{"error": "Failed to list notifications"})
		return
	}

	c.JSON(200, gin.H{"notifications": notifications})
}

type handleNotifyMemberStatusUpdateParams struct {
	TripID       uuid.UUID
	TargetUserID uuid.UUID
	OwnerID      uuid.UUID
	MemberStatus repository.MemberStatusEnum
	Queries      *repository.Queries
}

func (s *Server) handleNotifyMemberStatusUpdate(c *gin.Context, params handleNotifyMemberStatusUpdateParams) error {
	insertNotificationParams := repository.InsertNotificationParams{
		UserID: params.TargetUserID,
		TripID: params.TripID,
	}

	switch params.MemberStatus {
	case repository.MemberStatusEnumAccepted:
		insertNotificationParams.Type = repository.NotificationTypeTripInviteAccepted
		insertNotificationParams.Message = "A user has accepted your trip invitation"
		insertNotificationParams.UserID = params.OwnerID

	case repository.MemberStatusEnumDeclined:
		insertNotificationParams.Type = repository.NotificationTypeTripInviteDeclined
		insertNotificationParams.Message = "A user has declined your trip invitation"
		insertNotificationParams.UserID = params.OwnerID

	case repository.MemberStatusEnumRemoved:
		insertNotificationParams.Type = repository.NotificationTypeTripMemberRemoved
		insertNotificationParams.Message = "You have been removed from a trip"

	case repository.MemberStatusEnumOwner:
		insertNotificationParams.Type = repository.NotificationTypeTripOwnershipTransfer
		insertNotificationParams.Message = "You are now the owner of a trip"

	default:
		return nil
	}

	currentUser := s.ctxGetUser(c)

	metadata := dbtypes.NotificationMetadata{
		UserID:   currentUser.ID,
		UserName: currentUser.Name,
	}
	insertNotificationParams.Metadata = metadata

	err := params.Queries.InsertNotification(c, insertNotificationParams)

	if err != nil {
		s.log.LogError(c, "handleNotifyMemberStatusUpdate: InsertNotification failed", err)
		return err
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

	err := params.Queries.InsertNotification(c, repository.InsertNotificationParams{
		UserID:  params.TargetUserID,
		Type:    repository.NotificationTypeTripInvitePending,
		TripID:  params.TripID,
		Message: "You have been invited to a trip",
		Metadata: dbtypes.NotificationMetadata{
			UserID:   currentUser.ID,
			UserName: currentUser.Name,
		},
	})

	if err != nil {
		s.log.LogError(c, "handleNotifyTripInvite: InsertNotification failed", err)
		return err
	}

	return nil
}
