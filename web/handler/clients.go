package handler

import (
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/core/protocol"
	"github.com/rafaeljusto/druns/web/interceptor"
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

	Response []protocol.ClientResponse `response:"get"`
}

func (h *clients) Get(w http.ResponseWriter, r *http.Request) {
	clientDAO := dao.NewClient(h.DB())

	clients, err := clientDAO.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Response = clients.Protocol()
	w.WriteHeader(http.StatusOK)
}

func (h *clients) Interceptors() trama.AJAXInterceptorChain {
	return trama.NewAJAXInterceptorChain(
		interceptor.NewJSON(h),
		interceptor.NewDatabase(h),
	)
}
