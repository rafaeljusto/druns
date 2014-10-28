package protocol

type ClassRequest struct {
	Weekday *string `json:"weekday"`
	Time    *string `json:"time"`
}

type ClassResponse struct {
	Weekday string `json:"weekday"`
	Time    string `json:"time"`
}
