package handler

import (
	"net/http"

	"github.com/gustavo-hms/trama"
	//"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/webserver/interceptor"
)

func init() {
	Mux.RegisterService("/client/{id}", func() trama.AJAXHandler {
		return new(client)
	})
}

type client struct {
	trama.DefaultAJAXHandler
	interceptor.DatabaseCompliant
	interceptor.JSONCompliant

	Id string `param:"id"`
}

func (h *client) Get(w http.ResponseWriter, r *http.Request) {
	//clientDAO := dao.NewClient(h.DB())
	//client, err := clientDAO.FindById(h.Id)
}

func (h *client) Post(w http.ResponseWriter, r *http.Request) {

}

func (h *client) Put(w http.ResponseWriter, r *http.Request) {

}

func (h *client) Delete(w http.ResponseWriter, r *http.Request) {

}

func (h *client) Interceptors() trama.AJAXInterceptorChain {
	return trama.NewAJAXInterceptorChain(
		interceptor.NewJSON(h),
		interceptor.NewDatabase(h),
	)
}
