package enrollment

const (
	TypeReservation Type = iota
	TypeRegular     Type = iota
	TypeReplacement Type = iota
)

type Type int

type Enrollment struct {
	Type   Type
	Client Client
	Group  Group
}
