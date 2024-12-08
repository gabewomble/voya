package server

import (
	"net/http"
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

	// TODO: Make trip & user validation reusable, copy paste easier for now
	tripID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		s.badRequest(c, errorDetailsFromMessage("invalid trip id"))
		return
	}

	user := s.ctxGetUser(c)

	// Check tripID is valid
	_, err = s.db.Queries().GetTripById(c, repository.GetTripByIdParams{
		ID:     tripID,
		UserID: user.ID,
	})
	if err != nil {
		s.log.LogError(c, "addMemberToTripHandler: GetTripById failed", err)
		s.badRequest(c, errorDetailsFromMessage("unable to add member to this trip"))
		return
	}

	// Check userID is valid
	_, err = s.db.Queries().GetUserById(c, input.UserID)
	if err != nil {
		s.log.LogError(c, "addMemberToTripHandler: GetUserById failed", err)
		v.AddError("user_id", "unable to find user for user_id")
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
		return
	}

	err = s.db.Queries().AddUserToTrip(c, repository.AddUserToTripParams{
		TripID:    tripID,
		UserID:    input.UserID,
		InvitedBy: user.ID,
	})

	if err != nil {
		s.log.LogError(c, "addMemberToTripHandler: AddUserToTrip failed", err)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user invited to trip"})
}

func (s *Server) updateTripMemberStatusHandler(c *gin.Context) {
	var input struct {
		UserID       uuid.UUID `json:"user_id"`
		MemberStatus repository.MemberStatusEnum    `json:"member_status"`
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

	user := s.ctxGetUser(c)

	// Check tripID is valid
	_, err = s.db.Queries().GetTripById(c, repository.GetTripByIdParams{
		ID:     tripID,
		UserID: user.ID,
	})
	if err != nil {
		s.log.LogError(c, "updateTripMemberStatusHandler: GetTripById failed", err)
		s.badRequest(c, errorDetailsFromMessage("unable to add member to this trip"))
		return
	}

	// Check userID is valid
	_, err = s.db.Queries().GetUserById(c, input.UserID)
	if err != nil {
		s.log.LogError(c, "updateTripMemberStatusHandler: GetUserById failed", err)
		v.AddError("user_id", "unable to find user for user_id")
		s.unprocessableEntity(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v}))
		return
	}

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

	err = s.db.Queries().UpdateTripMemberStatus(c, updateTripParams)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user status updated"})
}
