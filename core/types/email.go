package types

import (
	"database/sql/driver"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core/errors"
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	var email Email
	err := email.Set(value)
	return email, err
}

func (e *Email) Set(value string) (err error) {
	e.value = strings.TrimSpace(value)
	e.value = strings.ToLower(e.value)

	if _, err = mail.ParseAddress(e.value); err != nil {
		err = errors.NewValidation(errors.ValidationCodeInvalidEmail, err)
	}

	return
}

func (e Email) MarshalText() ([]byte, error) {
	return []byte(e.value), nil
}

func (e *Email) UnmarshalText(data []byte) (err error) {
	return e.Set(string(data))
}

func (e Email) String() string {
	return e.value
}

func (e Email) Value() (driver.Value, error) {
	return e.value, nil
}

func (e *Email) Scan(src interface{}) (err error) {
	if src == nil {
		e.value = ""
		return
	}

	switch t := src.(type) {
	case bool, time.Time, int64, float64:
		return fmt.Errorf("Unsupported type to convert into an Email")

	case []byte:
		err = e.Set(string(t))
	case string:
		err = e.Set(t)
	}

	return
}
