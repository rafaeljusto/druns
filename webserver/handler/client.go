package handler

import (
	"fmt"
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/core/model"
	"github.com/rafaeljusto/druns/core/protocol"
	"github.com/rafaeljusto/druns/webserver/interceptor"
	"gopkg.in/mgo.v2"
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

	Request *protocol.ClientRequest `request:"post"`
}

func (h *newClient) Post(w http.ResponseWriter, r *http.Request) {
	var client model.Client
	if err := client.Apply(h.Request); err != nil {
		// TODO
	}

	clientDAO := dao.NewClient(h.DB())
	if err := clientDAO.Save(&client); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/client/%s", client.Id.String()))
	w.WriteHeader(http.StatusNoContent)
}

func (h *newClient) Interceptors() trama.AJAXInterceptorChain {
	return trama.NewAJAXInterceptorChain(
		interceptor.NewJSON(h),
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

	Id       string                   `param:"id"`
	Request  *protocol.ClientRequest  `request:"post"`
	Response *protocol.ClientResponse `response:"get"`
}

func (h *client) Get(w http.ResponseWriter, r *http.Request) {
	clientDAO := dao.NewClient(h.DB())

	client, err := clientDAO.FindById(h.Id)
	if err == mgo.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Response = client.Protocol()
	w.WriteHeader(http.StatusOK)
}

func (h *client) Put(w http.ResponseWriter, r *http.Request) {
	clientDAO := dao.NewClient(h.DB())

	client, err := clientDAO.FindById(h.Id)
	if err != mgo.ErrNotFound && err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := client.Apply(h.Request); err != nil {
		// TODO
	}

	if err := clientDAO.Save(&client); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *client) Delete(w http.ResponseWriter, r *http.Request) {
	clientDAO := dao.NewClient(h.DB())

	client, err := clientDAO.FindById(h.Id)
	if err == mgo.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := clientDAO.Delete(&client); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *client) Interceptors() trama.AJAXInterceptorChain {
	return trama.NewAJAXInterceptorChain(
		interceptor.NewJSON(h),
		interceptor.NewDatabase(h),
	)
}
