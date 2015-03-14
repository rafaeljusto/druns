package class

import (
	"time"

	"github.com/rafaeljusto/druns/core/group"
)

type Class struct {
	Id       int
	Group    group.Group
	Date     time.Time
	Students []Student
	revision uint64
}
