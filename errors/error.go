package errors

import (
	"fmt"
	"net/http"
)

var (
	// ErrorCodeBadRequestParse 400 bad request parse message
	ErrorCodeBadRequestParse = Error{
		Status:  http.StatusBadRequest,
		Message: "Not able to parse url",
	}
	// ErrorCodeNotFoundPattern 404 pattern error message
	ErrorCodeNotFoundPattern = Error{
		Status:  http.StatusNotFound,
		Message: "Pattern not found, use itm (for internal) or gtm (for external)",
	}
	// ErrorCodeBadRequestMarshal 400 marshal message
	ErrorCodeBadRequestMarshal = Error{
		Status:  http.StatusBadRequest,
		Message: "Not able to marshal http payload",
	}
)

// Error provides fields for error
// handling which come from F5 devices
// can also be used by common errors
type Error struct {
	Status     int           `json:"status"`
	Message    string        `json:"message"`
	ErrorStack []interface{} `json:"errorStack"`
}

// ErrorMessage returns only error message
// code status might be changed at the client-side
func (err Error) ErrorMessage() string {
	return fmt.Sprintf("%s", err.Message)
}
