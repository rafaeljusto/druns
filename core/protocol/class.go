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
		weekday = strings.TrimSpace(weekday)
		c.Weekday = &weekday
	}
}

func (c *ClassRequest) Validate(checkMandatoryFields bool) Translator {
	// TODO
	// 1. Time is in the correct format (RFC3339)
	// 2. Duration should be between 30 and 120
	return nil
}

type ClassResponse struct {
	Weekday  string `json:"weekday"`
	Time     string `json:"time"`
	Duration string `json:"duration"`
}
