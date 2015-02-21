package core

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Duration is always represented in minutes
type Duration struct {
	time.Duration
}

func (d *Duration) Set(value string) (err error) {
	value = strings.TrimSpace(value)

	if d.Duration, err = time.ParseDuration(fmt.Sprintf("%sm", value)); err != nil {
		err = NewValidationError(ValidationErrorCodeInvalidDuration, err)
	}

	return
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, int64(d.Minutes()))), nil
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf(`%d`, int64(d.Minutes()))), nil
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	return d.Set(string(data))
}

func (d *Duration) UnmarshalText(data []byte) error {
	return d.Set(string(data))
}

func (d Duration) String() string {
	return fmt.Sprintf("%d", int64(d.Minutes()))
}

func (d Duration) Value() (driver.Value, error) {
	return int64(d.Minutes()), nil
}

func (d *Duration) Scan(src interface{}) (err error) {
	if src == nil {
		d.Duration = 0
		return
	}

	switch t := src.(type) {
	case bool, time.Time:
		return NewError(fmt.Errorf("Unsupported type to convert into a Duration"))

	case int64:
		d.Duration, err = time.ParseDuration(fmt.Sprintf("%dm", int64(t)))
	case float64:
		d.Duration, err = time.ParseDuration(fmt.Sprintf("%dm", int64(t)))
	case []byte:
		err = d.Set(string(t))
	case string:
		err = d.Set(t)
	}

	return
}
