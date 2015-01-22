package data

type InternalServerErrorData struct {
	Ticket string
}

func NewInternalServerErrorData(ticket string) InternalServerErrorData {
	return InternalServerErrorData{
		Ticket: ticket,
	}
}
