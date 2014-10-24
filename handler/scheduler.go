package handler

import (
	"net/http"

	"github.com/gustavo-hms/trama"
)

func init() {
	Mux.RegisterService("/scheduler", func() trama.AJAXHandler {
		return new(scheduler)
	})
}

type scheduler struct {
	trama.DefaultAJAXHandler
}

func (h *scheduler) Get(w http.ResponseWriter, r *http.Request) {

}
