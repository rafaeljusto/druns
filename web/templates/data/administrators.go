package data

import "github.com/rafaeljusto/druns/core/model"

type Administrators struct {
	Logged
	Users []model.User
}

func NewAdministrators(username string, menu Menu, users []model.User) Administrators {
	return Administrators{
		Logged: NewLogged(username, menu),
		Users:  users,
	}
}