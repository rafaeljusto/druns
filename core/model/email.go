package model

import (
	"net/mail"
	"strings"

	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/tr"
)

type Email struct {
	value string
}

func (e *Email) Set(value string) error {
	e.value = strings.TrimSpace(value)
	e.value = strings.ToLower(e.value)

	if _, err := mail.ParseAddress(e.value); err != nil {
		return errors.NewValidation(tr.CodeInvalidEmail, err)
	}

	return nil
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
