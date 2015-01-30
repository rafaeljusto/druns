package data

import "github.com/rafaeljusto/druns/web/config"

type Login struct {
	Action  string
	Email   string
	Message string
}

func NewLogin(email string, message string) Login {
	return Login{
		Action:  config.DrunsConfig.URLs.GetHTTPS("login"),
		Email:   email,
		Message: message,
	}
}
