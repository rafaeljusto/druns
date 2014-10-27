package handler

import (
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/webserver/interceptor"
)

func init() {
	Mux.RegisterService("/clients", func() trama.AJAXHandler {
		return new(clients)
	})
}

type clients struct {
	trama.DefaultAJAXHandler
	interceptor.DatabaseCompliant
	interceptor.JSONCompliant
}

func (h *clients) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *clients) Interceptors() trama.AJAXInterceptorChain {
	return trama.NewAJAXInterceptorChain(
		interceptor.NewJSON(h),
		interceptor.NewDatabase(h),
	)
}
