package group

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/place"
	"github.com/rafaeljusto/druns/core/types"
)

const (
	TypeWeekley string = "Weekley"
	TypeOnce    string = "Once"
)

type Type struct {
	value string
}

func (g *Type) Set(value string) (err error) {
	value = strings.TrimSpace(value)
	switch value {
	case TypeWeekley:
		g.value = TypeWeekley
	case TypeOnce:
		g.value = TypeOnce
	default:
		err = errors.NewValidation(errors.ValidationCodeInvalidGroupType, err)
	}
	return
}

func (g Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, g.value)), nil
}

func (g Type) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf(`%s`, g.value)), nil
}

func (g *Type) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), ` "`)
	return g.Set(value)
}

func (g *Type) UnmarshalText(data []byte) error {
	return g.Set(string(data))
}

func (g Type) String() string {
	return g.value
}

func (g Type) Value() (driver.Value, error) {
	return g.value, nil
}

func (g *Type) Scan(src interface{}) error {
	if src == nil {
		g.value = ""
		return errors.New(fmt.Errorf("Unsupported type to convert into a Type"))
	}

	switch t := src.(type) {
	case int64, float64, bool, time.Time:
		return errors.New(fmt.Errorf("Unsupported type to convert into a Type"))

	case []byte:
		return g.Set(string(t))
	case string:
		return g.Set(t)
	}

	return nil
}

/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////

type Group struct {
	Id       int
	Name     types.Name
	Place    place.Place
	Weekday  types.Weekday
	Time     types.Time
	Duration types.Duration
	Type     Type
	Capacity int
}

func (g Group) Equal(other Group) bool {
	return g == other
}

/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////

type Groups []Group
