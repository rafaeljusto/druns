package protocol

type ClientRequest struct {
	Name    *string        `json:"name"`
	Classes []ClassRequest `json:"classes"`
}

func (c *ClientRequest) Normalize() {
	for i, class := range c.Classes {
		class.Normalize()
		c.Classes[i] = class
	}
}

func (c *ClientRequest) Validate(checkMandatoryFields bool) Translator {
	var messagesHolder MessagesHolder
	for _, class := range c.Classes {
		messagesHolder.Add(class.Validate(checkMandatoryFields))
	}
	return messagesHolder.Messages()
}

type ClientResponse struct {
	Id      string          `json:"id"`
	Name    string          `json:"name"`
	Classes []ClassResponse `json:"classes"`
}
