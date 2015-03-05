package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core/errors"
)

type Weekday struct {
	time.Weekday
}

func NewWeekday(value time.Weekday) Weekday {
	return Weekday{
		Weekday: value,
	}
}

func (w *Weekday) Set(value string) (err error) {
	value = strings.TrimSpace(value)
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
		err = errors.NewValidation(errors.ValidationCodeInvalidWeekday, err)
	}
	return
}

func (w Weekday) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, w.String())), nil
}

func (w Weekday) MarshalText() ([]byte, error) {
	return []byte(w.String()), nil
}

func (w *Weekday) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), ` "`)
	return w.Set(value)
}

func (w *Weekday) UnmarshalText(data []byte) error {
	return w.Set(string(data))
}

func (w Weekday) Value() (driver.Value, error) {
	return w.String(), nil
}

func (w *Weekday) Scan(src interface{}) error {
	if src == nil {
		return errors.New(fmt.Errorf("Unsupported type to convert into a Weekday"))
	}

	switch t := src.(type) {
	case int64, float64, bool, time.Time:
		return errors.New(fmt.Errorf("Unsupported type to convert into a Weekday"))

	case []byte:
		w.Set(string(t))
	case string:
		w.Set(t)
	}

	return nil
}
