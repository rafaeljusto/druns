package data

type Former interface {
	AddFieldMessage(field, message string)
	AddMessage(message string)
}

type Form struct {
	FieldMessage map[string]string
	Message      string
}

func NewForm() Form {
	return Form{
		FieldMessage: make(map[string]string),
	}
}

func (f *Form) AddFieldMessage(field, message string) {
	f.FieldMessage[field] = message
}

func (f *Form) AddMessage(message string) {
	f.Message = message
}
