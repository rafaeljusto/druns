package payment

import (
	"time"

	"github.com/rafaeljusto/druns/core/client"
)

const (
	StatusPending Status = iota
	StatusPayed
)

type Status int

type Payment struct {
	Id        int
	Client    client.Client
	Status    Status
	ExpiresAt time.Time
	Value     int
	revision  uint64
}
