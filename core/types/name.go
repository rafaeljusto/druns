package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Name struct {
	value string
}

func NewName(value string) Name {
	var name Name
	name.Set(value)
	return name
}

func (n *Name) Set(value string) {
	n.value = strings.TrimSpace(value)
	n.value = strings.Title(n.value)
}

func (n Name) MarshalText() ([]byte, error) {
	return []byte(n.value), nil
}

func (n *Name) UnmarshalText(data []byte) (err error) {
	n.Set(string(data))
	return
}

func (n Name) String() string {
	return n.value
}

func (n Name) Value() (driver.Value, error) {
	return n.value, nil
}

func (n *Name) Scan(src interface{}) error {
	if src == nil {
		n.value = ""
		return nil
	}

	switch t := src.(type) {
	case int64, float64, bool, time.Time:
		return fmt.Errorf("Unsupported type to convert into a Name")

	case []byte:
		n.Set(string(t))
	case string:
		n.Set(t)
	}

	return nil
}
