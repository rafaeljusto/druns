package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core/errors"
)

type Address struct {
	value string
}

func NewAddress(value string) Address {
	var Address Address
	Address.Set(value)
	return Address
}

func (n *Address) Set(value string) {
	n.value = strings.TrimSpace(value)
}

func (n Address) MarshalText() ([]byte, error) {
	return []byte(n.value), nil
}

func (n *Address) UnmarshalText(data []byte) (err error) {
	n.Set(string(data))
	return
}

func (n Address) String() string {
	return n.value
}

func (n Address) Value() (driver.Value, error) {
	return n.value, nil
}

func (n *Address) Scan(src interface{}) error {
	if src == nil {
		n.value = ""
		return nil
	}

	switch t := src.(type) {
	case int64, float64, bool, time.Time:
		return errors.New(fmt.Errorf("Unsupported type to convert into an Address"))

	case []byte:
		n.Set(string(t))
	case string:
		n.Set(t)
	}

	return nil
}
