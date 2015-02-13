package model

import "time"

const (
	GroupTypeWeekley GroupType = "Weekley"
	GroupTypeOnce    GroupType = "Once"
)

type GroupType string

type Group struct {
	Id       int
	Weekday  time.Weekday
	Time     time.Time
	Duration time.Duration
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
