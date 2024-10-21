package server

import (
	"errors"
	"fmt"
	"net/http"
	"server/internal/repository"
	"server/internal/validator"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"errors": v.Errors, "error": errors.New("form is not valid") })
		return
	}

	trip, err := s.db.Queries().InsertTrip(c, input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.Redirect(http.StatusCreated, fmt.Sprintf("/trips/%d", trip.ID))
}
