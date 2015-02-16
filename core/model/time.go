package model

import (
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/tr"
)

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(`"15:04"`)), nil
}

func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.Format("15:04")), nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	value := strings.TrimSpace(string(data))
	t.Time, err = time.Parse(`"15:04"`, value)
	if err != nil {
		err = errors.NewValidation(tr.CodeInvalidTime, err)
	}
	return
}

func (t *Time) UnmarshalText(data []byte) (err error) {
	value := strings.TrimSpace(string(data))
	t.Time, err = time.Parse(`15:04`, value)
	if err != nil {
		err = errors.NewValidation(tr.CodeInvalidTime, err)
	}
	return
}

func (t Time) String() string {
	if t.IsZero() {
		return ""
	}

	return t.Format("15:04")
}
