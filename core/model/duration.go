package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/tr"
)

// Duration is always represented in minutes
type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, int(d.Minutes()))), nil
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf(`%d`, int(d.Minutes()))), nil
}

func (d *Duration) UnmarshalJSON(data []byte) (err error) {
	value := strings.TrimSpace(string(data))
	d.Duration, err = time.ParseDuration(fmt.Sprintf("%sm", value))
	if err != nil {
		err = errors.NewValidation(tr.CodeInvalidDuration, err)
	}
	return
}

func (d *Duration) UnmarshalText(data []byte) (err error) {
	value := strings.TrimSpace(string(data))
	d.Duration, err = time.ParseDuration(fmt.Sprintf("%sm", value))
	if err != nil {
		err = errors.NewValidation(tr.CodeInvalidDuration, err)
	}
	return
}

func (d Duration) String() string {
	return fmt.Sprintf("%d", d.Minutes())
}
