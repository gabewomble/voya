package server

import (
	"net/http"
	"server/internal/data"
	"server/internal/repository"
	"server/internal/validator"

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
	if ok, err := s.validateTripAccess(c, validateTripAccessParams{
		TripID: tripID,
		UserID: user.ID,
		IsEdit: true,
	}); !ok {
		s.handleInvalidTripAccess(c, handleInvalidTripAccessParams{
			validator: v,
			err:       err,
		})
		return
	}

	// Validate target user
	if ok, err := s.validateUser(c, input.UserID); !ok {
		s.handleInvalidUser(c, handleInvalidUserParams{
			validator: v,
			err:       err,
		})
		return
	}

	// Validate target user is not already a member
	isMember, err := s.db.Queries().CheckUserIsTripMember(c, repository.CheckUserIsTripMemberParams{
		TripID: tripID,
		UserID: input.UserID,
	})

	if err != nil {
		s.log.LogError(c, "addMemberToTripHandler: CheckUserIsTripMember failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
	}

	if isMember {
		v.AddError("user_id", "user is already a member of this trip")
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
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
	if ok, err := s.validateTripAccess(c, validateTripAccessParams{
		TripID: tripID,
		UserID: currentUser.ID,
		IsEdit: false,
	}); !ok {
		s.handleInvalidTripAccess(c, handleInvalidTripAccessParams{
			validator: v,
			err:       err,
		})
		return
	}

	// Validate target user
	if ok, err := s.validateUser(c, input.UserID); !ok {
		s.handleInvalidUser(c, handleInvalidUserParams{
			validator: v,
			err:       err,
		})
		return
	}

	// Get current member info
	s.log.LogInfo(c, "updateTripMemberStatusHandler: GetTripMember", "trip_id", tripID, "user_id", currentUser.ID)
	currentMember, err := s.db.Queries().GetTripMember(c, repository.GetTripMemberParams{
		TripID: tripID,
		UserID: currentUser.ID,
	})
	if err != nil {
		s.log.LogError(c, "updateTripMemberStatusHandler: GetTripmember failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromMessage("failed to retrieve user member information"))
		return
	}

	// Get target trip member info
	s.log.LogInfo(c, "updateTripMemberStatusHandler: GetTripMember", "trip_id", tripID, "user_id", input.UserID)
	targetMember, err := s.db.Queries().GetTripMember(c, repository.GetTripMemberParams{
		TripID: tripID,
		UserID: input.UserID,
	})
	if err != nil {
		s.log.LogError(c, "updateTripMemberStatusHandler: GetTripMember failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromMessage("failed to retrieve trip member information"))
		return
	}

	// Validate member status
	s.log.LogInfo(c, "updateTripMemberStatusHandler: ValidateUpdateMemberStatus")
	data.ValidateUpdateMemberStatus(data.ValidateUpdateMemberStatusParams{
		Validator:     v,
		Fieldname:     "member_status",
		CurrentMember: currentMember,
		TargetMember:  targetMember,
		MemberStatus:  input.MemberStatus,
	})

	if !v.Valid() {
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
		return
	}

	// Get trip owner
	s.log.LogInfo(c, "updateTripMemberStatusHandler: GetTripOwner", "trip_id", tripID)
	owner, err := s.db.Queries().GetTripOwner(c, tripID)
	if err != nil {
		s.log.LogError(c, "updateTripMemberStatusHandler: GetTripOwner failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromMessage("failed to retrieve trip member information"))
		return
	}

	tx, err := s.db.Tx(c)

	if err != nil {
		s.log.LogError(c, "updateTripMemberStatusHandler: Tx failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	queries := s.db.Queries().WithTx(tx)

	updateTripMemberStatusParams := repository.UpdateTripMemberStatusParams{
		MemberStatus: input.MemberStatus,
		TripID:       tripID,
		UserID:       input.UserID,
		UpdatedBy:    currentUser.ID,
	}

	// Update trip member status
	s.log.LogInfo(c, "updateTripMemberStatusHandler: UpdateTripMemberStatus")
	err = queries.UpdateTripMemberStatus(c, updateTripMemberStatusParams)
	if err != nil {
		tx.Rollback(c)
		s.log.LogError(c, "updateTripMemberStatusHandler: UpdateTripMemberStatus failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	// Create notification(s)
	s.log.LogInfo(c, "updateTripMemberStatusHandler: handleNotifyMemberStatusUpdate")
	err = s.handleNotifyMemberStatusUpdate(c, handleNotifyMemberStatusUpdateParams{
		TripID:       tripID,
		TargetUser:   targetMember,
		Owner:        owner,
		MemberStatus: input.MemberStatus,
		Queries:      queries,
	})

	if err != nil {
		tx.Rollback(c)
		s.log.LogError(c, "updateTripMemberStatusHandler: handleNotifyMemberStatusUpdate failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	err = tx.Commit(c)
	if err != nil {
		s.log.LogError(c, "updateTripMemberStatusHandler: Commit failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user status updated"})
}
