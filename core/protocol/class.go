package protocol

import (
	"time"
)

type ClassRequest struct {
	Weekday time.Weekday `json:"weekday"`
	Time    time.Time    `json:"time"`
}

type ClassResponse struct {
	Weekday time.Weekday `json:"weekday"`
	Time    time.Time    `json:"time"`
}
