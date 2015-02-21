package core

import (
	"fmt"
)

const (
	ValidationErrorCodeAuthenticationError ValidationErrorCode = "authentication-error"
	ValidationErrorCodeInvalidDate         ValidationErrorCode = "invalid-date"
	ValidationErrorCodeInvalidDuration     ValidationErrorCode = "invalid-duration"
	ValidationErrorCodeInvalidEmail        ValidationErrorCode = "invalid-email"
	ValidationErrorCodeInvalidGroupType    ValidationErrorCode = "invalid-group-type"
	ValidationErrorCodeInvalidTime         ValidationErrorCode = "invalid-time"
	ValidationErrorCodeInvalidWeekday      ValidationErrorCode = "invalid-weekday"
	ValidationErrorCodeSessionExpired      ValidationErrorCode = "session-expired" // TODO: Remove this from the core?
)

type ValidationErrorCode string

type ValidationError struct {
	Err  error
	Code ValidationErrorCode
}

func NewValidationError(code ValidationErrorCode, err error) error {
	return ValidationError{
		Err:  err,
		Code: code,
	}
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s", v.Code)
}
