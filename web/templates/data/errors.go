package data

type InternalServerError struct {
	Ticket string
}

func NewInternalServerError(ticket string) InternalServerError {
	return InternalServerError{
		Ticket: ticket,
	}
}
