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
	ValidationCodeSunday                ValidationCode = "sunday"          // TODO: Remove this from the core?
	ValidationCodeMonday                ValidationCode = "monday"          // TODO: Remove this from the core?
	ValidationCodeTuesday               ValidationCode = "tuesday"         // TODO: Remove this from the core?
	ValidationCodeWednesday             ValidationCode = "wednesday"       // TODO: Remove this from the core?
	ValidationCodeThursday              ValidationCode = "thursday"        // TODO: Remove this from the core?
	ValidationCodeFriday                ValidationCode = "friday"          // TODO: Remove this from the core?
	ValidationCodeSaturday              ValidationCode = "saturday"        // TODO: Remove this from the core?
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
