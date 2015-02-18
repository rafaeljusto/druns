package model

import (
	"fmt"
	"strings"

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

func (g GroupType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, g.value)), nil
}

func (g GroupType) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf(`%s`, g.value)), nil
}

func (g *GroupType) UnmarshalJSON(data []byte) (err error) {
	value := strings.TrimSpace(string(data))
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

func (g *GroupType) UnmarshalText(data []byte) (err error) {
	value := strings.TrimSpace(string(data))
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

func (g GroupType) String() string {
	return g.value
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
