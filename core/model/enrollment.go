package model

const (
	EnrollmentTypeReservation EnrollmentType = iota
	EnrollmentTypeRegular     EnrollmentType = iota
	EnrollmentTypeReplacement EnrollmentType = iota
)

type EnrollmentType int

type Enrollment struct {
	Type   EnrollmentType
	Client Client
	Group  Group
}
