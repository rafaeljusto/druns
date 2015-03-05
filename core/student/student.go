package model

import (
	"github.com/rafaeljusto/druns/core/class"
)

type Student struct {
	Class    class.Class
	Attended bool
}
