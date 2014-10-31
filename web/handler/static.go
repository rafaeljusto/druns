package handler

import (
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/web/interceptor"
)

func init() {
	Mux.RegisterPage("/", func() trama.WebHandler {
		return new(static)
	})
}

type static struct {
	trama.DefaultWebHandler
	interceptor.LanguageCompliant
}

func (h *static) Get(response trama.Response, r *http.Request) {
	response.ExecuteTemplate("index.html", nil)
}

func (h *static) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewAcceptLanguagePage(h),
	)
}
