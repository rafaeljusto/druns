package handler

import (
	"fmt"
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/core/model"
	"github.com/rafaeljusto/druns/core/protocol"
	"github.com/rafaeljusto/druns/web/interceptor"
)

func init() {
	Mux.RegisterService("/client", func() trama.AJAXHandler {
		return new(newClient)
	})

	Mux.RegisterService("/client/{id}", func() trama.AJAXHandler {
		return new(client)
	})
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

type newClient struct {
	trama.DefaultAJAXHandler
	interceptor.DatabaseCompliant
	interceptor.JSONCompliant
	interceptor.LanguageCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LogCompliant

	handle  string
	Request *protocol.ClientRequest `request:"post"`
}

func (h *newClient) SetHandle(handle string) {
	h.handle = handle
}

func (h *newClient) AuthSecret(secretId string) (string, error) {
	// TODO!
	return "abc123", nil
}

func (h *newClient) Post(w http.ResponseWriter, r *http.Request) {
	var client model.Client
	if err := client.Apply(h.Request); err != nil {
		h.SetMessage(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	clientDAO := dao.NewClient(h.Tx(), h.RemoteAddress(), h.handle)
	if err := clientDAO.Save(&client); err != nil {
		h.Logger().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/client/%d", client.Id))
	w.WriteHeader(http.StatusNoContent)
}

func (h *newClient) Interceptors() trama.AJAXInterceptorChain {
	return trama.NewAJAXInterceptorChain(
		interceptor.NewAcceptLanguage(h),
		interceptor.NewRemoteAddress(h),
		interceptor.NewLog(h),
		interceptor.NewJSON(h),
		interceptor.NewValidate(h),
		interceptor.NewAuth(h),
		interceptor.NewDatabase(h),
	)
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

type client struct {
	trama.DefaultAJAXHandler
	interceptor.DatabaseCompliant
	interceptor.JSONCompliant
	interceptor.LanguageCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LogCompliant

	handle   string
	Id       int                      `param:"id"`
	Request  *protocol.ClientRequest  `request:"put"`
	Response *protocol.ClientResponse `response:"get"`
}

func (h *client) SetHandle(handle string) {
	h.handle = handle
}

func (h *client) AuthSecret(secretId string) (string, error) {
	// TODO!
	return "abc123", nil
}

func (h *client) Get(w http.ResponseWriter, r *http.Request) {
	clientDAO := dao.NewClient(h.Tx(), h.RemoteAddress(), h.handle)

	client, err := clientDAO.FindById(h.Id)
	if err == core.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		h.Logger().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Response = client.Protocol()
	w.WriteHeader(http.StatusOK)
}

func (h *client) Put(w http.ResponseWriter, r *http.Request) {
	clientDAO := dao.NewClient(h.Tx(), h.RemoteAddress(), h.handle)

	client, err := clientDAO.FindById(h.Id)
	if err == core.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		h.Logger().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := client.Apply(h.Request); err != nil {
		h.SetMessage(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := clientDAO.Save(&client); err != nil {
		h.Logger().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *client) Interceptors() trama.AJAXInterceptorChain {
	return trama.NewAJAXInterceptorChain(
		interceptor.NewAcceptLanguage(h),
		interceptor.NewRemoteAddress(h),
		interceptor.NewLog(h),
		interceptor.NewJSON(h),
		interceptor.NewValidate(h),
		interceptor.NewAuth(h),
		interceptor.NewDatabase(h),
	)
}
