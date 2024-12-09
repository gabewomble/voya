package server

import (
	"net/http"
	"server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	if notifications == nil {
		notifications = make([]repository.ListNotificationsRow, 0)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(200, gin.H{"notifications": notifications, "total": 0})
			return
		}
		s.log.LogError(c, "listNotificationsHandler: ListNotifications failed", err)
		c.JSON(500, gin.H{"error": "Failed to list notifications"})
		return
	}

	c.JSON(200, gin.H{"notifications": notifications, "total": len(notifications)})
}

func (s *Server) listUnreadNotificationsHandler(c *gin.Context) {
	user := s.ctxGetUser(c)

	notifications, err := s.db.Queries().GetUnreadNotifications(c, user.ID)

	if notifications == nil {
		notifications = make([]repository.GetUnreadNotificationsRow, 0)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(200, gin.H{"notifications": notifications, "total": 0})
			return
		}
		s.log.LogError(c, "listUnreadNotificationsHandler: ListUnreadNotifications failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromMessage("Failed to list unread notifications"))
		return
	}

	c.JSON(200, gin.H{"notifications": notifications, "total": len(notifications)})
}

func (s *Server) countUnreadNotificationsHandler(c *gin.Context) {
	user := s.ctxGetUser(c)

	count, err := s.db.Queries().CountUnreadNotifications(c, user.ID)

	if err != nil {
		s.log.LogError(c, "countUnreadNotificationsHandler: CountUnreadNotifications failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromMessage("Failed to count unread notifications"))
		return
	}

	c.JSON(200, gin.H{"count": count})
}

func (s *Server) markNotificationsAsReadHandler(c *gin.Context) {
	user := s.ctxGetUser(c)

	err := s.db.Queries().MarkNotificationsAsRead(c, user.ID)

	if err != nil {
		s.log.LogError(c, "markNotificationsAsReadHandler: MarkNotificationsAsRead failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromMessage("Failed to mark notifications as read"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifications marked as read"})
}

func (s *Server) getNotificationByIdHandler(c *gin.Context) {
	user := s.ctxGetUser(c)
	notificationID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, errorDetailsFromMessage("Invalid notification ID"))
		return
	}

	notification, err := s.db.Queries().GetNotificationById(c, repository.GetNotificationByIdParams{
		UserID: user.ID,
		ID:     notificationID,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			s.errorResponse(c, http.StatusNotFound, errorDetailsFromMessage("Notification not found"))
			return
		}
		s.log.LogError(c, "getNotificationByIdHandler: GetNotificationByID failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromMessage("Failed to get notification"))
		return
	}

	c.JSON(http.StatusOK, notification)
}

func (s *Server) markNotificationAsReadHandler(c *gin.Context) {
	user := s.ctxGetUser(c)
	notificationID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, errorDetailsFromMessage("Invalid notification ID"))
		return
	}

	err = s.db.Queries().MarkNotificationAsRead(c, repository.MarkNotificationAsReadParams{
		UserID: user.ID,
		ID:     notificationID,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			s.errorResponse(c, http.StatusNotFound, errorDetailsFromMessage("Notification not found"))
			return
		}
		s.log.LogError(c, "markNotificationAsReadHandler: MarkNotificationAsRead failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromMessage("Failed to mark notification as read"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (s *Server) deleteNotificationHandler(c *gin.Context) {
	user := s.ctxGetUser(c)
	notificationID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, errorDetailsFromMessage("Invalid notification ID"))
		return
	}

	err = s.db.Queries().DeleteNotification(c, repository.DeleteNotificationParams{
		UserID: user.ID,
		ID:     notificationID,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			s.errorResponse(c, http.StatusNotFound, errorDetailsFromMessage("Notification not found"))
			return
		}
		s.log.LogError(c, "deleteNotificationHandler: DeleteNotification failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromMessage("Failed to delete notification"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted"})
}

type handleNotifyMemberStatusUpdateParams struct {
	TripID       uuid.UUID
	Owner        repository.GetTripOwnerRow
	TargetUser   repository.GetTripMemberRow
	MemberStatus repository.MemberStatusEnum
	Queries      *repository.Queries
}

func (s *Server) handleNotifyMemberStatusUpdate(c *gin.Context, params handleNotifyMemberStatusUpdateParams) error {
	switch params.MemberStatus {
	case repository.MemberStatusEnumAccepted:
		return s.handleNotifyInviteAccepted(c, params)
	case repository.MemberStatusEnumDeclined:
		return s.handleNotifyInviteDeclined(c, params)
	case repository.MemberStatusEnumRemoved:
		return s.handleNotifyUserRemoved(c, params)
	case repository.MemberStatusEnumOwner:
		return s.handleNotifyOwnershipTransfer(c, params)
	case repository.MemberStatusEnumCancelled:
		return s.handleNotifyInviteCancelled(c, params)

	default:
		return nil
	}
}

func (s *Server) handleNotifyInviteAccepted(c *gin.Context, params handleNotifyMemberStatusUpdateParams) error {
	// Remove pending notification for user
	err := s.removePendingInvite(c, params)

	if err != nil {
		return err
	}

	// Notify trip members
	err = params.Queries.NotifyOtherTripMembers(c, repository.NotifyOtherTripMembersParams{
		TripID:       params.TripID,
		Message:      "A user has accepted the trip invitation",
		Type:         repository.NotificationTypeTripInviteAccepted,
		CreatedBy:    params.TargetUser.ID,
		TargetUserID: params.TargetUser.ID,
	})
	if err != nil {
		s.log.LogError(c, "handleNotifyInviteAccepted: NotifyTripMembers failed", err)
		return err
	}

	return nil
}

func (s *Server) handleNotifyInviteDeclined(c *gin.Context, params handleNotifyMemberStatusUpdateParams) error {
	return s.removePendingInvite(c, params)
}

func (s *Server) handleNotifyInviteCancelled(c *gin.Context, params handleNotifyMemberStatusUpdateParams) error {
	return s.removePendingInvite(c, params)
}

func (s *Server) handleNotifyUserRemoved(c *gin.Context, params handleNotifyMemberStatusUpdateParams) error {
	currentUser := s.ctxGetUser(c)
	err := params.Queries.NotifyOtherTripMembers(c, repository.NotifyOtherTripMembersParams{
		TripID:       params.TripID,
		Message:      "A user has been removed from the trip",
		Type:         repository.NotificationTypeTripMemberRemoved,
		CreatedBy:    currentUser.ID,
		TargetUserID: params.TargetUser.ID,
	})

	if err != nil {
		s.log.LogError(c, "handleNotifyUserRemoved: NotifyTripMembers failed", err)
		return err
	}

	err = params.Queries.DeleteNotificationsByType(c, repository.DeleteNotificationsByTypeParams{
		UserID: params.TargetUser.ID,
		TripID: params.TripID,
		Type:   repository.NotificationTypeTripInviteAccepted,
	})

	if err != nil {
		s.log.LogError(c, "handleNotifyUserRemoved: DeleteNotificationsByType failed", err)
		return err
	}

	return nil
}

func (s *Server) handleNotifyOwnershipTransfer(c *gin.Context, params handleNotifyMemberStatusUpdateParams) error {
	currentUser := s.ctxGetUser(c)
	err := params.Queries.NotifyOtherTripMembers(c, repository.NotifyOtherTripMembersParams{
		TripID:       params.TripID,
		Message:      "The trip ownership has been transferred",
		Type:         repository.NotificationTypeTripOwnershipTransfer,
		CreatedBy:    currentUser.ID,
		TargetUserID: params.TargetUser.ID,
	})

	if err != nil {
		s.log.LogError(c, "handleNotifyOwnershipTransfer: NotifyTripMembers failed", err)
		return err
	}

	return nil
}

func (s *Server) removePendingInvite(c *gin.Context, params handleNotifyMemberStatusUpdateParams) error {
	// Remove pending notification for user
	err := params.Queries.DeleteNotificationsByType(c, repository.DeleteNotificationsByTypeParams{
		UserID: params.TargetUser.ID,
		TripID: params.TripID,
		Type:   repository.NotificationTypeTripInvitePending,
	})

	if err != nil {
		s.log.LogError(c, "handleNotifyInviteDeclined: DeleteNotificationsByType failed", err)
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

	err := params.Queries.DeleteNotificationsByType(c, repository.DeleteNotificationsByTypeParams{
		UserID: params.TargetUserID,
		TripID: params.TripID,
		Type:   repository.NotificationTypeTripMemberRemoved,
	})

	if err != nil {
		s.log.LogError(c, "handleNotifyTripInvite: DeleteNotificationsByType failed", err)
		return err
	}

	err = params.Queries.InsertNotification(c, repository.InsertNotificationParams{
		UserID:       params.TargetUserID,
		Type:         repository.NotificationTypeTripInvitePending,
		TripID:       params.TripID,
		Message:      "You have been invited to a trip",
		CreatedBy:    currentUser.ID,
		TargetUserID: params.TargetUserID,
	})

	if err != nil {
		s.log.LogError(c, "handleNotifyTripInvite: InsertNotification failed", err)
		return err
	}

	return nil
}
