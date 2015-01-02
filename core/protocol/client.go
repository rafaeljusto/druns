package protocol

type ClientRequest struct {
}

func (c *ClientRequest) Normalize() {
}

func (c *ClientRequest) Validate(checkMandatoryFields bool) Translator {
	var messagesHolder MessagesHolder
	return messagesHolder.Messages()
}

type ClientResponse struct {
}
