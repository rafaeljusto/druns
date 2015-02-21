package user

import (
	"github.com/rafaeljusto/druns/core"
)

type User struct {
	Id       int
	Name     core.Name
	Email    core.Email
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
