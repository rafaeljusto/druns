package protocol

import (
	"strings"
)

type ClassRequest struct {
	Weekday  *string `json:"weekday"`
	Time     *string `json:"time"`
	Duration *string `json:"duration"`
}

func (c *ClassRequest) Normalize() {
	if c.Weekday != nil {
		weekday := strings.ToLower(*c.Weekday)
		c.Weekday = &weekday
	}
}

func (c *ClassRequest) Validate() {
	// 1. Time is in the correct format (hh:mm:ss)
	// 2. Duration should be between 30 and 120
}

type ClassResponse struct {
	Weekday  string `json:"weekday"`
	Time     string `json:"time"`
	Duration string `json:"duration"`
}
