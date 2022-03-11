package errors

import (
	"net/http"
)

type ErrorCode string

const (
	BadRequest    = "ERR_API_BAD_REQUEST"
	InvalidTicket = "ERR_API_INVALID_TICKET"
	UnableToTrack = "ERR_API_UNABLE_TO_TRACK"
)

var ApiErrors = map[ErrorCode]string{
	BadRequest:    "Invalid request body",
	InvalidTicket: "Invalid ticket",
	UnableToTrack: "Unable to track source and destination for the given tickets",
}

type ErrorResponse struct {
	HttpStatusCode int       `json:"status"`
	ErrorCode      ErrorCode `json:"error_code,omitempty"`
	ErrorMessage   string    `json:"error_message,omitempty"`
}


//Create new error responses
func NewErrorResponse(httpStatusCode int, errorCode ErrorCode, errorMessage string) *ErrorResponse {
	return &ErrorResponse{
		HttpStatusCode: httpStatusCode,
		ErrorCode:      errorCode,
		ErrorMessage:   errorMessage,
	}
}

func (e ErrorResponse) Error() string {
	return e.ErrorMessage
}

var ErrBadRequest = NewErrorResponse(http.StatusBadRequest, BadRequest, ApiErrors[BadRequest])
var ErrInvalidTicket = NewErrorResponse(http.StatusBadRequest, InvalidTicket, ApiErrors[InvalidTicket])
var ErrUnableToTrack = NewErrorResponse(http.StatusUnprocessableEntity, UnableToTrack, ApiErrors[UnableToTrack])
