package model

import "time"

const (
	GroupTypeWeekley GroupType = iota
	GroupTypeOnce    GroupType = iota
)

type GroupType int

type Group struct {
	Weekday  time.Weekday
	Time     time.Time
	Duration time.Duration
	Type     GroupType
	Capacity int
}
