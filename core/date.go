package core

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

func NewDate(value time.Time) Date {
	return Date{
		Time: value,
	}
}

func (d *Date) Set(value string) (err error) {
	value = strings.TrimSpace(value)

	if d.Time, err = time.Parse(`2006-01-02`, value); err != nil {
		err = NewValidationError(ValidationErrorCodeInvalidDate, err)
	}

	return
}

// http://golang.org/src/time/time.go # MarshalJSON()
func (d Date) MarshalJSON() ([]byte, error) {
	if y := d.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, fmt.Errorf("Time.MarshalJSON: year outside of range [0,9999]")
	}

	return []byte(d.Format(`"2006-01-02"`)), nil
}

func (d Date) MarshalText() ([]byte, error) {
	if y := d.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, fmt.Errorf("Time.MarshalTEXT: year outside of range [0,9999]")
	}

	return []byte(d.Format(`2006-01-02`)), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), ` "`)
	return d.Set(value)
}

func (d *Date) UnmarshalText(data []byte) error {
	return d.Set(string(data))
}

func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.Format("2006-01-02")
}

func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}

func (d *Date) Scan(src interface{}) (err error) {
	if src == nil {
		d.Time = time.Time{}
		return
	}

	switch t := src.(type) {
	case int64, float64, bool:
		return NewError(fmt.Errorf("Unsupported type to convert into a Date"))

	case time.Time:
		d.Time = t
	case []byte:
		err = d.Set(string(t))
	case string:
		err = d.Set(t)
	}

	return
}
