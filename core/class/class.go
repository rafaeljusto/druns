package class

import (
	"time"

	"github.com/rafaeljusto/druns/core/group"
)

type Class struct {
	Id       int
	Group    group.Group
	Date     time.Time
	revision uint64
}
