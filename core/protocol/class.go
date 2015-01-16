package protocol

type ClassRequest struct {
}

func (c *ClassRequest) Normalize() {

}

func (c *ClassRequest) Validate(checkMandatoryFields bool) Translator {
	var messageHolder MessagesHolder

	return messageHolder.Messages()
}

type ClassResponse struct {
}
