package enrollment

import (
	"github.com/rafaeljusto/druns/core/client"
	"github.com/rafaeljusto/druns/core/group"
)

const (
	TypeReservation Type = "Reservation"
	TypeRegular     Type = "Regular"
	TypeReplacement Type = "Replacement"
)

type Type string

type Enrollment struct {
	Id     int
	Type   Type
	Client client.Client
	Group  group.Group
}

func (e Enrollment) Equal(other Enrollment) bool {
	return e == other
}

/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////

type Enrollments []Enrollment
