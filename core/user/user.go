package user

import (
	"github.com/rafaeljusto/druns/core/types"
)

type User struct {
	Id       int
	Name     types.Name
	Email    types.Email
	Password string
}

func (u User) Equal(other User) bool {
	if u.Id != other.Id ||
		u.Name != other.Name ||
		u.Email != other.Email ||
		u.Password != other.Password {

		return false
	}

	return true
}
