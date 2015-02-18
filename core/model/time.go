package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/tr"
)

type Time struct {
	time.Time
}

func (t *Time) Set(value string) (err error) {
	value = strings.TrimSpace(value)

	if t.Time, err = time.Parse(`15:04`, value); err != nil {
		err = errors.NewValidation(tr.CodeInvalidTime, err)
	}

	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(`"15:04"`)), nil
}

func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.Format("15:04")), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), ` "`)
	return t.Set(value)
}

func (t *Time) UnmarshalText(data []byte) error {
	return t.Set(string(data))
}

func (t Time) String() string {
	if t.IsZero() {
		return ""
	}

	return t.Format("15:04")
}

func (t Time) Value() (driver.Value, error) {
	return t.Format("15:04"), nil
}

func (t *Time) Scan(src interface{}) (err error) {
	if src == nil {
		t.Time = time.Time{}
		return
	}

	switch v := src.(type) {
	case int64, float64, bool:
		return core.NewError(fmt.Errorf("Unsupported type to convert into a Time"))

	case time.Time:
		t.Time = v
	case []byte:
		err = t.Set(string(v))
	case string:
		err = t.Set(v)
	}

	return
}
