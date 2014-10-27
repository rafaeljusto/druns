package handler

import (
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/webserver/interceptor"
)

func init() {
	Mux.RegisterService("/scheduler", func() trama.AJAXHandler {
		return new(scheduler)
	})
}

type scheduler struct {
	trama.DefaultAJAXHandler
	interceptor.DatabaseCompliant
	interceptor.JSONCompliant
}

func (h *scheduler) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *scheduler) Interceptors() trama.AJAXInterceptorChain {
	return trama.NewAJAXInterceptorChain(
		interceptor.NewJSON(h),
		interceptor.NewDatabase(h),
	)
}
