package model

import (
	"strconv"
	"time"

	"github.com/rafaeljusto/druns/core/protocol"
)

type Class struct {
	Weekday  time.Weekday
	Time     time.Time
	Duration time.Duration
}

func (c *Class) Apply(request *protocol.ClassRequest) protocol.Translator {
	var messageHolder protocol.MessagesHolder

	if request.Weekday != nil {
		var ok bool
		c.Weekday, ok = protocol.ParseWeekday(*request.Weekday)
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
	var classes Classes
	var messageHolder protocol.MessagesHolder

	for _, request := range requests {
		if request.Weekday == nil || request.Time == nil {
			messageHolder.Add(protocol.NewMessageWithField(protocol.MsgCodeClassDataMissing, "", ""))
		}

		weekday, ok := protocol.ParseWeekday(*request.Weekday)
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
		for _, class := range c {
			if class.Weekday == weekday && class.Time.Equal(classTime) {
				messageHolder.Add(class.Apply(&request))

				classes = append(classes, class)
				found = true
				break
			}
		}

		if !found {
			var class Class
			messageHolder.Add(class.Apply(&request))
			classes = append(classes, class)
		}
	}

	return classes, messageHolder.Messages()
}

func (c Classes) Protocol() []protocol.ClassResponse {
	var response []protocol.ClassResponse
	for _, class := range c {
		response = append(response, *class.Protocol())
	}
	return response
}
