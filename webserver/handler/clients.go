package handler

import (
	"net/http"

	"github.com/gustavo-hms/trama"
)

func init() {
	Mux.RegisterService("/clients", func() trama.AJAXHandler {
		return new(clients)
	})
}

type clients struct {
	trama.DefaultAJAXHandler
}

func (h *clients) Get(w http.ResponseWriter, r *http.Request) {

}
