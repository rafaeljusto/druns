package handler

import (
	"github.com/gustavo-hms/trama"
)

func init() {
	Mux.RegisterPage("/client/{id}", func() trama.WebHandler {
		return new(client)
	})
}

type client struct {
	handy.DefaultHandler

	Id string `param:"id"`
}

func (h *client) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *client) Post(w http.ResponseWriter, r *http.Request) {

}

func (h *client) Put(w http.ResponseWriter, r *http.Request) {

}

func (h *client) Delete(w http.ResponseWriter, r *http.Request) {

}
