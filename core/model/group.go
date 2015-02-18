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

const (
	GroupTypeWeekley string = "Weekley"
	GroupTypeOnce    string = "Once"
)

type GroupType struct {
	value string
}

func (g *GroupType) Set(value string) (err error) {
	value = strings.TrimSpace(value)
	switch value {
	case GroupTypeWeekley:
		g.value = GroupTypeWeekley
	case GroupTypeOnce:
		g.value = GroupTypeOnce
	default:
		err = errors.NewValidation(tr.CodeInvalidGroupType, err)
	}
	return
}

func (g GroupType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, g.value)), nil
}

func (g GroupType) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf(`%s`, g.value)), nil
}

func (g *GroupType) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), ` "`)
	return g.Set(value)
}

func (g *GroupType) UnmarshalText(data []byte) error {
	return g.Set(string(data))
}

func (g GroupType) String() string {
	return g.value
}

func (g GroupType) Value() (driver.Value, error) {
	return g.value, nil
}

func (g *GroupType) Scan(src interface{}) error {
	if src == nil {
		g.value = ""
		return core.NewError(fmt.Errorf("Unsupported type to convert into a GroupType"))
	}

	switch t := src.(type) {
	case int64, float64, bool, time.Time:
		return core.NewError(fmt.Errorf("Unsupported type to convert into a GroupType"))

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
	Name     Name
	Weekday  Weekday
	Time     Time
	Duration Duration
	Type     GroupType
	Capacity int
}

func (g *Group) Equal(other Group) bool {
	return *g == other
}

/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////

type Groups []Group
