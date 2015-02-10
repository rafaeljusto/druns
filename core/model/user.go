package model

type User struct {
	Id       int
	Name     Name
	Email    Email
	Password string
}

func (u *User) Equal(other User) bool {
	if u.Id != other.Id ||
		u.Name != other.Name ||
		u.Email != other.Email ||
		u.Password != other.Password {

		return false
	}

	return true
}
