package protocol

import (
	"strings"
	"time"
)

type ClassRequest struct {
	Weekday  *string `json:"weekday"`
	Time     *string `json:"time"`
	Duration *string `json:"duration"`
}

func (c *ClassRequest) Normalize() {
	if c.Weekday != nil {
		weekday := strings.ToLower(*c.Weekday)
		weekday = strings.TrimSpace(weekday)
		c.Weekday = &weekday
	}
}

func (c *ClassRequest) Validate(checkMandatoryFields bool) Translator {
	var messageHolder MessagesHolder

	if checkMandatoryFields {
		if c.Weekday == nil || c.Time == nil {
			messageHolder.Add(NewMessageWithField(MsgCodeClassDataMissing, "", ""))
		}
	}

	if c.Weekday != nil {
		if _, ok := ParseWeekday(*c.Weekday); !ok {
			messageHolder.Add(NewMessageWithField(MsgCodeInvalidClassWeekday,
				"weekday", *c.Weekday))
		}
	}

	if c.Time != nil {
		if _, err := time.Parse(time.RFC3339, *c.Time); err != nil {
			messageHolder.Add(NewMessageWithField(MsgCodeInvalidClassTime,
				"time", *c.Time))
		}
	}

	if c.Duration != nil {
		duration, err := time.ParseDuration(*c.Duration + "m")
		if err != nil || duration < (30*time.Minute) || duration > (120*time.Minute) {
			messageHolder.Add(NewMessageWithField(MsgCodeInvalidClassDuration,
				"duration", *c.Duration))
		}
	}

	return messageHolder.Messages()
}

type ClassResponse struct {
	Weekday  string `json:"weekday"`
	Time     string `json:"time"`
	Duration string `json:"duration"`
}

func ParseWeekday(weekday string) (time.Weekday, bool) {
	switch weekday {
	case "sunday":
		return time.Sunday, true
	case "monday":
		return time.Monday, true
	case "tuesday":
		return time.Tuesday, true
	case "wednesday":
		return time.Wednesday, true
	case "thursday":
		return time.Thursday, true
	case "friday":
		return time.Friday, true
	case "saturday":
		return time.Saturday, true
	default:
		return time.Sunday, false
	}
}
