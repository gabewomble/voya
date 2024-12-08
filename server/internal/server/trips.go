package server

import (
	"errors"
	"fmt"
	"net/http"
	"server/internal/data"
	"server/internal/repository"
	"server/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Server) listTripsHandler(c *gin.Context) {
	trips, err := s.db.Queries().ListTrips(c, s.ctxGetUser(c).ID)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, []ErrorDetail{{Message: err.Error()}})
		return
	}

	if trips == nil {
		trips = make([]repository.Trip, 0)
	}

	c.JSON(http.StatusOK, gin.H{"trips": trips})
}

func (s *Server) createTripHandler(c *gin.Context) {
	input := repository.InsertTripParams{}

	if err := c.BindJSON(&input); err != nil {
		s.badRequest(c, []ErrorDetail{{Message: err.Error()}})
		return
	}

	v := validator.New()

	if data.ValidateTrip(v, &input); !v.Valid() {
		s.badRequest(c, errorDetailsFromValidator(ErrorDetailFromValidatorInput{v: v, message: "form is not valid"}))
		return
	}

	userID := s.ctxGetUser(c).ID

	tx, err := s.db.Tx(c)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	queries := s.db.Queries().WithTx(tx)

	trip, err := queries.InsertTrip(c, input)

	if err != nil {
		tx.Rollback(c)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	err = queries.InsertTripOwner(c, repository.InsertTripOwnerParams{
		TripID:  trip.ID,
		OwnerID: userID,
	})

	if err != nil {
		tx.Rollback(c)
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	err = tx.Commit(c)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	s.log.LogInfo(c, "createTripHandler: trip created", "trip", trip)

	c.Header("Location", fmt.Sprintf("/trips/%s", trip.ID))
	c.JSON(http.StatusCreated, gin.H{"trip": trip})
}

func (s *Server) getTripByIdHandler(c *gin.Context) {
	tripID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		s.badRequest(c, errorDetailsFromMessage("invalid id"))
		return
	}

	user := s.ctxGetUser(c)

	trip, err := s.db.Queries().GetTripById(c, repository.GetTripByIdParams{
		ID:     tripID,
		UserID: user.ID,
	})
	if err != nil {
		s.log.LogError(c, "getTripByIdHandler: GetTripById failed", err)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.notFoundResponse(c, errorDetailsFromMessage("trip not found"))
		default:
			s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		}
		return
	}

	members, err := s.db.Queries().GetTripMembers(c, trip.ID)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	if members == nil {
		members = make([]repository.GetTripMembersRow, 0)
	}

	c.JSON(http.StatusOK, gin.H{"trip": trip, "members": members})
}

func (s *Server) deleteTripByIdHandler(c *gin.Context) {
	tripID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		s.badRequest(c, errorDetailsFromMessage("invalid id"))
		return
	}

	_, err = s.db.Queries().GetTripById(c, repository.GetTripByIdParams{
		ID:     tripID,
		UserID: s.ctxGetUser(c).ID,
	})
	if err != nil {
		s.badRequest(c, errorDetailsFromError(err))
		return
	}

	err = s.db.Queries().DeleteTripById(c, repository.DeleteTripByIdParams{
		ID:     tripID,
		UserID: s.ctxGetUser(c).ID,
	})
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, errorDetailsFromError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
