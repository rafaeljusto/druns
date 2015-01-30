package data

type Administrators struct {
	Logged
}

func NewAdministrators(user string, menu Menu) Administrators {
	return Administrators{
		Logged: NewLogged(user, menu),
	}
}
