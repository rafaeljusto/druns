package model

import (
	"net"
	"time"
)

const (
	LogOperationCreation LogOperation = "Create"
	LogOperationUpdate   LogOperation = "Update"
	LogOperationRemoval  LogOperation = "Delete"
)

type LogOperation string

type Log struct {
	Id        int64
	Agent     int
	IPAddress net.IP
	ChangedAt time.Time
	Operation LogOperation
}

func NewLog(agent int, ipAddress net.IP, opeation LogOperation) Log {
	return Log{
		Agent:     agent,
		IPAddress: ipAddress,
		ChangedAt: time.Now(),
		Operation: opeation,
	}
}
