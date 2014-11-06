package model

import (
	"strconv"
	"time"

	"github.com/rafaeljusto/druns/core/protocol"
)

func parseWeekday(weekday string) (time.Weekday, bool) {
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

type Class struct {
	Weekday  time.Weekday
	Time     time.Time
	Duration time.Duration
}

func (c *Class) Apply(request *protocol.ClassRequest) protocol.Translator {
	var messageHolder protocol.MessagesHolder

	if request.Weekday != nil {
		var ok bool
		c.Weekday, ok = parseWeekday(*request.Weekday)
		if !ok {
			messageHolder.Add(protocol.NewMessageWithField(protocol.MsgCodeInvalidClassWeekday,
				"weekday", *request.Weekday))
		}
	}

	if request.Time != nil {
		var err error
		c.Time, err = time.Parse(time.RFC3339, *request.Time)
		if err != nil {
			messageHolder.Add(protocol.NewMessageWithField(protocol.MsgCodeInvalidClassTime,
				"time", *request.Time))
		}
	}

	if request.Duration != nil {
		var err error
		c.Duration, err = time.ParseDuration(*request.Duration + "m")
		if err != nil {
			messageHolder.Add(protocol.NewMessageWithField(protocol.MsgCodeInvalidClassDuration,
				"duration", *request.Duration))
		}
	}

	return messageHolder.Messages()
}

func (c *Class) Protocol() *protocol.ClassResponse {
	return &protocol.ClassResponse{
		Weekday:  c.Weekday.String(),
		Time:     c.Time.UTC().Format(time.RFC3339),
		Duration: strconv.FormatFloat(c.Duration.Minutes(), 'f', -1, 64),
	}
}

/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////

type Classes []Class

func (c Classes) Apply(requests []protocol.ClassRequest) (Classes, protocol.Translator) {
	var messageHolder protocol.MessagesHolder

	for _, request := range requests {
		if request.Weekday == nil || request.Time == nil {
			messageHolder.Add(protocol.NewMessageWithField(protocol.MsgCodeClassDataMissing, "", ""))
		}

		weekday, ok := parseWeekday(*request.Weekday)
		if !ok {
			messageHolder.Add(protocol.NewMessageWithField(protocol.MsgCodeInvalidClassWeekday,
				"weekday", *request.Weekday))
			continue
		}

		classTime, err := time.Parse(time.RFC3339, *request.Time)
		if err != nil {
			messageHolder.Add(protocol.NewMessageWithField(protocol.MsgCodeInvalidClassTime,
				"time", *request.Time))
			continue
		}

		found := false
		for i, class := range c {
			if class.Weekday == weekday && class.Time.Equal(classTime) {
				messageHolder.Add(class.Apply(&request))

				c[i] = class
				found = true
				break
			}
		}

		if !found {
			var class Class
			messageHolder.Add(class.Apply(&request))
			c = append(c, class)
		}
	}

	// TODO: Remove classes of the system!

	return c, messageHolder.Messages()
}

func (c Classes) Protocol() []protocol.ClassResponse {
	var response []protocol.ClassResponse
	for _, class := range c {
		response = append(response, *class.Protocol())
	}
	return response
}
