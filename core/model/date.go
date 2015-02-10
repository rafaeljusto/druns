package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/tr"
)

type Date struct {
	time.Time
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

func (d *Date) UnmarshalJSON(data []byte) (err error) {
	value := strings.TrimSpace(string(data))
	d.Time, err = time.Parse(`"2006-01-02"`, value)
	if err != nil {
		err = errors.NewValidation(tr.CodeInvalidDate, err)
	}
	return
}

func (d *Date) UnmarshalText(data []byte) (err error) {
	value := strings.TrimSpace(string(data))
	d.Time, err = time.Parse(`2006-01-02`, value)
	if err != nil {
		err = errors.NewValidation(tr.CodeInvalidDate, err)
	}
	return
}

func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.Format("2006-01-02")
}
