package class

import "github.com/rafaeljusto/druns/core/enrollment"

type Student struct {
	Id         int
	Enrollment enrollment.Enrollment
	Attended   bool
	revision   uint64
}
