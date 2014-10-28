package model

import (
	"time"

	"github.com/rafaeljusto/druns/core/protocol"
)

type Class struct {
	Weekday time.Weekday
	Time    time.Time
}

func (c *Class) Apply(request *protocol.ClassRequest) protocol.Translator {
	return nil
}

func (c *Class) Protocol() *protocol.ClassResponse {
	return &protocol.ClassResponse{
		Weekday: c.Weekday.String(),
		Time:    c.Time.Format(time.RFC1123),
	}
}

/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////

type Classes []Class

func (c Classes) Protocol() []protocol.ClassResponse {
	var response []protocol.ClassResponse
	for _, class := range c {
		response = append(response, *class.Protocol())
	}
	return response
}
