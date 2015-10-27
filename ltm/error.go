package ltm

import "fmt"

// ErrorLTM provides fields for error
// handling which come from F5 devices
type ErrorLTM struct {
	Code       int           `json:"code"`
	Message    string        `json:"message"`
	ErrorStack []interface{} `json:"errorStack"`
}

// ErrorMessage returns only error message
// code status might be changed at the client-side
func (err ErrorLTM) ErrorMessage() string {
	return fmt.Sprintf("%s", err.Message)
}
