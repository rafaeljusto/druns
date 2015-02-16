package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/tr"
)

type Weekday struct {
	time.Weekday
}

func (w Weekday) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, w.String())), nil
}

func (w Weekday) MarshalText() ([]byte, error) {
	return []byte(w.String()), nil
}

func (w *Weekday) UnmarshalJSON(data []byte) (err error) {
	value := strings.TrimSpace(string(data))
	switch value {
	case time.Sunday.String():
		w.Weekday = time.Sunday
	case time.Monday.String():
		w.Weekday = time.Monday
	case time.Tuesday.String():
		w.Weekday = time.Tuesday
	case time.Wednesday.String():
		w.Weekday = time.Wednesday
	case time.Thursday.String():
		w.Weekday = time.Thursday
	case time.Friday.String():
		w.Weekday = time.Friday
	case time.Saturday.String():
		w.Weekday = time.Saturday
	default:
		err = errors.NewValidation(tr.CodeInvalidWeekday, err)
	}
	return
}

func (w *Weekday) UnmarshalText(data []byte) (err error) {
	value := strings.TrimSpace(string(data))
	switch value {
	case time.Sunday.String():
		w.Weekday = time.Sunday
	case time.Monday.String():
		w.Weekday = time.Monday
	case time.Tuesday.String():
		w.Weekday = time.Tuesday
	case time.Wednesday.String():
		w.Weekday = time.Wednesday
	case time.Thursday.String():
		w.Weekday = time.Thursday
	case time.Friday.String():
		w.Weekday = time.Friday
	case time.Saturday.String():
		w.Weekday = time.Saturday
	default:
		err = errors.NewValidation(tr.CodeInvalidWeekday, err)
	}
	return
}
