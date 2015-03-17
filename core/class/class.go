package class

import (
	"time"

	"github.com/rafaeljusto/druns/core/group"
)

type Class struct {
	Id       int
	Group    group.Group
	BeginAt  time.Time
	EndAt    time.Time
	Students []Student
	revision uint64
}
