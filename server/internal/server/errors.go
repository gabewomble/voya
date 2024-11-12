package server

import (
	"net/http"
	"server/internal/validator"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Errors []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func errorDetailsFromMessage(message string) []ErrorDetail {
	return []ErrorDetail{{Message: message}}
}

func errorDetailsFromError(err error) []ErrorDetail {
	return []ErrorDetail{{Message: err.Error()}}
}

type ErrorDetailFromValidatorInput struct {
	v       *validator.Validator
	message string
}

func errorDetailsFromValidator(arg ErrorDetailFromValidatorInput) []ErrorDetail {
	message := "form is not valid"

	if arg.message != "" {
		message = arg.message
	}

	response := []ErrorDetail{{Message: message}}

	for key, value := range arg.v.Errors {
		response = append(response, ErrorDetail{Field: key, Message: value})
	}

	return response
}

func (s *Server) unprocessableEntity(c *gin.Context, errors []ErrorDetail) {
	s.errorResponse(c, http.StatusUnprocessableEntity, errors)
}

func (s *Server) badRequest(c *gin.Context, errors []ErrorDetail) {
	s.errorResponse(c, http.StatusBadRequest, errors)
}

func (s *Server) invalidCredentialsResponse(c *gin.Context) {
	s.errorResponse(c, http.StatusUnauthorized, []ErrorDetail{{Message: "invalid authentication credentials"}})
}

func (s *Server) invalidAuthTokenResponse(c *gin.Context) {
	c.Header("WWW-Authenticate", "Bearer")
	s.errorResponse(c, http.StatusUnauthorized, []ErrorDetail{{Message: "invalid or missing authentication token"}})
}

func (s *Server) notFoundResponse(c *gin.Context, errors []ErrorDetail) {
	s.errorResponse(c, http.StatusNotFound, errors)
}

func (s *Server) errorResponse(c *gin.Context, status int, errors []ErrorDetail) {
	s.log.LogInfo(c, "Error response", "status", status, "errors", errors)
	c.JSON(status, ErrorResponse{Errors: errors})
}
