package user

import (
	"github.com/rafaeljusto/druns/core/types"
)

type User struct {
	Id       int
	Name     types.Name
	Email    types.Email
	Password string
	revision uint64
}
