package handler

import (
	"github.com/gustavo-hms/trama"
)

func init() {
	Mux.RegisterPage("/clients", func() trama.WebHandler {
		return new(clients)
	})
}

type clients struct {
	handy.DefaultHandler
}

func (h *clients) Get(w http.ResponseWriter, r *http.Request) {

}
