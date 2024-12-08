package server

import (
	"net/http"
	"server/internal/data"
	"server/internal/repository"
	"server/internal/validator"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) addMemberToTripHandler(c *gin.Context) {
	var input struct {
		UserID uuid.UUID `json:"user_id"`
	}

	if err := c.BindJSON(&input); err != nil {
		s.badRequest(c, errorDetailsFromError(err))
		return
	}

	v := validator.New()
	v.Check(input.UserID != uuid.Nil, "user_id", "must be provided")

	if !v.Valid() {
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
		return
	}

	tripID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		s.badRequest(c, errorDetailsFromMessage("invalid trip id"))
		return
	}

	user := s.ctxGetUser(c)

	// Validate trip access
	if err = s.validateTripAccess(c, tripID, user.ID); err != nil {
		if err == ErrTripNotFound {
			s.badRequest(c, errorDetailsFromMessage("unable to add member to this trip"))
			return
		}
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	// Validate target user
	if err = s.validateUser(c, input.UserID); err != nil {
		if err == ErrUserNotFound {
			v.AddError("user_id", "unable to find user for user_id")
			s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
			return
		}
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	tx, err := s.db.Tx(c)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	queries := s.db.Queries().WithTx(tx)

	// Add target user to trip
	err = queries.AddUserToTrip(c, repository.AddUserToTripParams{
		TripID:    tripID,
		UserID:    input.UserID,
		InvitedBy: user.ID,
	})

	if err != nil {
		tx.Rollback(c)
		s.log.LogError(c, "addMemberToTripHandler: AddUserToTrip failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	// Add notification for target user
	err = s.handleNotifyTripInvite(c, handleNotifyTripInviteParams{
		TripID:       tripID,
		TargetUserID: input.UserID,
		Queries:      queries,
	})

	if err != nil {
		tx.Rollback(c)
		s.log.LogError(c, "addMemberToTripHandler: handleNotifyTripInvite failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
	}

	tx.Commit(c)

	c.JSON(http.StatusOK, gin.H{"message": "user invited to trip"})
}

func (s *Server) updateTripMemberStatusHandler(c *gin.Context) {
	var input struct {
		UserID       uuid.UUID                   `json:"user_id"`
		MemberStatus repository.MemberStatusEnum `json:"member_status"`
	}

	if err := c.BindJSON(&input); err != nil {
		s.badRequest(c, errorDetailsFromError(err))
		return
	}

	v := validator.New()
	v.Check(input.UserID != uuid.Nil, "user_id", "must be provided")
	v.Check(input.MemberStatus != "", "member_status", "must be provided")

	if !v.Valid() {
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
		return
	}

	tripID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		s.badRequest(c, errorDetailsFromMessage("invalid trip id"))
		return
	}

	currentUser := s.ctxGetUser(c)

	// Validate trip access
	if err = s.validateTripAccess(c, tripID, currentUser.ID); err != nil {
		if err == ErrTripNotFound {
			v.AddError("trip_id", "unable to find or access trip for trip_id")
			s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
			return
		}
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	// Validate target user
	if err = s.validateUser(c, input.UserID); err != nil {
		if err == ErrUserNotFound {
			v.AddError("user_id", "unable to find user for user_id")
			s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
			return
		}
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	// Get trip owner id
	owner, err := s.db.Queries().GetTripOwner(c, tripID)
	if err != nil {
		s.log.LogError(c, "updateTripMemberStatusHandler: GetTripOwner failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromMessage("failed to retrieve trip information"))
		return
	}

	// Validate member status
	data.ValidateUpdateMemberStatus(data.ValidateUpdateMemberStatusParams{
		Validator:    v,
		Fieldname:    "member_status",
		UserID:       currentUser.ID,
		TargetUserID: input.UserID,
		MemberStatus: input.MemberStatus,
		OwnerID:      owner.ID,
	})

	if !v.Valid() {
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
		return
	}

	tx, err := s.db.Tx(c)

	if err != nil {
		s.log.LogError(c, "updateTripMemberStatusHandler: Tx failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	queries := s.db.Queries().WithTx(tx)

	updateTripParams := repository.UpdateTripMemberStatusParams{
		MemberStatus: input.MemberStatus,
		TripID:       tripID,
		UserID:       input.UserID,
	}

	if input.MemberStatus == repository.MemberStatusEnumRemoved {
		user := s.ctxGetUser(c)
		updateTripParams.RemovedBy = user.ID
		updateTripParams.RemovedAt = time.Now()
	}

	// Update trip member status
	err = queries.UpdateTripMemberStatus(c, updateTripParams)
	if err != nil {
		tx.Rollback(c)
		s.log.LogError(c, "updateTripMemberStatusHandler: UpdateTripMemberStatus failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	// Create notification(s)
	err = s.handleNotifyMemberStatusUpdate(c, handleNotifyMemberStatusUpdateParams{
		TripID:       tripID,
		TargetUserID: input.UserID,
		OwnerID:      owner.ID,
		MemberStatus: input.MemberStatus,
		Queries:      queries,
	})

	if err != nil {
		tx.Rollback(c)
		s.log.LogError(c, "updateTripMemberStatusHandler: handleNotifyMemberStatusUpdate failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user status updated"})
}
