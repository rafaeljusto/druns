package enrollment

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/client"
	"github.com/rafaeljusto/druns/core/group"
)

const (
	TypeReservation string = "Reservation"
	TypeRegular     string = "Regular"
	TypeReplacement string = "Replacement"
)

type Type struct {
	value string
}

func (t *Type) Set(value string) (err error) {
	value = strings.TrimSpace(value)
	switch value {
	case TypeReservation:
		t.value = TypeReservation
	case TypeRegular:
		t.value = TypeRegular
	case TypeReplacement:
		t.value = TypeReplacement
	default:
		err = core.NewValidationError(core.ValidationErrorCodeInvalidEnrollmentType, err)
	}
	return
}

func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.value)), nil
}

func (t Type) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf(`%s`, t.value)), nil
}

func (t *Type) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), ` "`)
	return t.Set(value)
}

func (t *Type) UnmarshalText(data []byte) error {
	return t.Set(string(data))
}

func (t Type) String() string {
	return t.value
}

func (t Type) Value() (driver.Value, error) {
	return t.value, nil
}

func (t *Type) Scan(src interface{}) error {
	if src == nil {
		t.value = ""
		return core.NewError(fmt.Errorf("Unsupported type to convert into a Type"))
	}

	switch value := src.(type) {
	case int64, float64, bool, time.Time:
		return core.NewError(fmt.Errorf("Unsupported type to convert into a Type"))

	case []byte:
		return t.Set(string(value))
	case string:
		return t.Set(value)
	}

	return nil
}

/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////

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
