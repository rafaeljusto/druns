package dblog

import (
	"net"
	"time"
)

const (
	OperationCreation Operation = "Create"
	OperationUpdate   Operation = "Update"
	OperationRemoval  Operation = "Delete"
)

type Operation string

type DBLog struct {
	Id        int64
	Agent     int
	IPAddress net.IP
	ChangedAt time.Time
	Operation Operation
}

func NewDBLog(agent int, ipAddress net.IP, operation Operation) DBLog {
	return DBLog{
		Agent:     agent,
		IPAddress: ipAddress,
		ChangedAt: time.Now(),
		Operation: operation,
	}
}
