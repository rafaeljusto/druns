package handler

import (
	"github.com/gustavo-hms/trama"
)

func init() {
	Mux.RegisterPage("/scheduler", func() trama.WebHandler {
		return new(scheduler)
	})
}

type scheduler struct {
	handy.DefaultHandler
}

func (h *scheduler) Get(w http.ResponseWriter, r *http.Request) {

}
