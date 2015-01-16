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
	interceptor.LanguageCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LogCompliant

	handle   string
	Response []protocol.ClientResponse `response:"get"`
}

func (h *clients) SetHandle(handle string) {
	h.handle = handle
}

func (h *clients) AuthSecret(secretId string) (string, error) {
	// TODO!
	return "abc123", nil
}

func (h *clients) Get(w http.ResponseWriter, r *http.Request) {
	clientDAO := dao.NewClient(h.Tx(), h.RemoteAddress(), h.handle)

	clients, err := clientDAO.FindAll()
	if err != nil {
		h.Logger().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Response = clients.Protocol()
	w.WriteHeader(http.StatusOK)
}

func (h *clients) Interceptors() trama.AJAXInterceptorChain {
	return trama.NewAJAXInterceptorChain(
		interceptor.NewAcceptLanguage(h),
		interceptor.NewRemoteAddress(h),
		interceptor.NewLog(h),
		interceptor.NewJSON(h),
		interceptor.NewAuth(h),
		interceptor.NewDatabase(h),
	)
}
