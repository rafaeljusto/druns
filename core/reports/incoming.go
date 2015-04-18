package reports

import (
	"time"

	"github.com/rafaeljusto/druns/core/group"
)

type Incoming struct {
	Group    group.Group
	Month    time.Time
	Foreseen float64
	Value    float64
}
