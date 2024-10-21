package server

import (
	"errors"
	"fmt"
	"net/http"
	"server/internal/repository"
	"server/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Server) listTripsHandler (c *gin.Context) {
	trips, err := s.db.Queries().ListTrips(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if trips == nil {
		trips = make([]repository.Trip, 0)
	}

	c.JSON(http.StatusOK, gin.H{ "trips": trips })
}

func (s *Server) createTripHandler(c *gin.Context) {
	input := repository.InsertTripParams{}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	v := validator.New()

	if validator.ValidateTrip(v, &input); !v.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors, "error": "form is not valid" })
		return
	}

	trip, err := s.db.Queries().InsertTrip(c, input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	s.logger.LogInfo(c, fmt.Sprintf("trip created: %s", trip.ID.String()))

	c.Redirect(http.StatusCreated, fmt.Sprintf("/trips/%d", trip.ID))
}

func (s *Server) getTripByIdHandler(c *gin.Context) {
	tripID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	trip, err := s.db.Queries().GetTripById(c, tripID)
	if err != nil {
		switch {
			case errors.Is(err, pgx.ErrNoRows):
				c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusFound, gin.H{"trip": trip})
}
