package errors

import (
	"fmt"
)

const (
	ValidationCodeAuthenticationError   ValidationCode = "authentication-error"
	ValidationCodeInvalidDate           ValidationCode = "invalid-date"
	ValidationCodeInvalidDuration       ValidationCode = "invalid-duration"
	ValidationCodeInvalidEmail          ValidationCode = "invalid-email"
	ValidationCodeInvalidEnrollmentType ValidationCode = "invalid-enrollment-type"
	ValidationCodeInvalidGroupType      ValidationCode = "invalid-group-type"
	ValidationCodeInvalidTime           ValidationCode = "invalid-time"
	ValidationCodeInvalidWeekday        ValidationCode = "invalid-weekday"
	ValidationCodeSessionExpired        ValidationCode = "session-expired" // TODO: Remove this from the core?
)

type ValidationCode string

type Validation struct {
	Err  error
	Code ValidationCode
}

func NewValidation(code ValidationCode, err error) error {
	return Validation{
		Err:  err,
		Code: code,
	}
}

func (v Validation) Error() string {
	return fmt.Sprintf("%s", v.Code)
}
