package model

const (
	GroupTypeWeekley GroupType = iota
	GroupTypeOnce    GroupType
)

type GroupType int

type Group struct {
	Weekday  time.Weekday
	Time     time.Time
	Duration time.Duration
	Type     GroupType
	Capacity int
}
