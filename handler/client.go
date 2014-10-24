package handler

import (
	"net/http"

	"github.com/gustavo-hms/trama"
)

func init() {
	Mux.RegisterService("/client/{id}", func() trama.AJAXHandler {
		return new(client)
	})
}

type client struct {
	trama.DefaultAJAXHandler

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
