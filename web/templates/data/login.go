package data

type Login struct {
	Email   string
	Message string
}

func NewLogin(email string, message string) Login {
	return Login{
		Email:   email,
		Message: message,
	}
}
