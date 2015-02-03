package data

type Login struct {
	Form
	Email string
}

func NewLogin(email string) Login {
	return Login{
		Form:  NewForm(),
		Email: email,
	}
}
