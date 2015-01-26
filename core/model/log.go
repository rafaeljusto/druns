package model

import (
	"net"
	"time"
)

const (
	LogOperationCreation LogOperation = "CREATE"
	LogOperationUpdate   LogOperation = "UPDATE"
	LogOperationRemoval  LogOperation = "DELETE"
)

type LogOperation string

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
