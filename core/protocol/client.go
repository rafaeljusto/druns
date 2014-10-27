package protocol

type ClientRequest struct {
	Name    string         `json:"name"`
	Classes []ClassRequest `json:"classes"`
}

type ClientResponse struct {
	Id      int             `json:"id"`
	Name    string          `json:"name"`
	Classes []ClassResponse `json:"classes"`
}
