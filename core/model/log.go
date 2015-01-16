package model

import (
	"net"
	"time"
)

const (
	LogOperationUpdate   LogOperation = iota
	LogOperationRemoval  LogOperation = iota
	LogOperationCreation LogOperation = iota
)

type LogOperation int

type Log struct {
	Id        int64
	Handle    string
	IPAddress net.IP
	ChangedAt time.Time
	Operation LogOperation
}

func NewLog() Log {
	return Log{
		ChangedAt: time.Now(),
		Operation: LogOperationUpdate,
	}
}
