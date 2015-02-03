package data

type Form struct {
	FieldMessage map[string]string
	Message      string
}

func NewForm() Form {
	return Form{
		FieldMessage: make(map[string]string),
	}
}
