package interceptor

import (
	"net"
	"net/http"
	"strings"

	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/web/templates/data"
)

type remoteAddresser interface {
	RemoteAddress() net.IP
	SetRemoteAddress(net.IP)
}

////////////////////////////////////////////////////////////
/////////////////////// AJAX ///////////////////////////////
////////////////////////////////////////////////////////////

type RemoteAddressAJAX struct {
	trama.NopAJAXInterceptor
	handler remoteAddresser
}

func NewRemoteAddressAJAX(h remoteAddresser) *RemoteAddressAJAX {
	return &RemoteAddressAJAX{handler: h}
}

func (i *RemoteAddressAJAX) Before(w http.ResponseWriter, r *http.Request) {
	clientAddress, err := remoteAddress(r)
	if err != nil {
		log.Notice(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i.handler.SetRemoteAddress(clientAddress)
}

////////////////////////////////////////////////////////////
/////////////////////// WEB ////////////////////////////////
////////////////////////////////////////////////////////////

type RemoteAddressWeb struct {
	trama.NopWebInterceptor
	handler remoteAddresser
}

func NewRemoteAddressWeb(h remoteAddresser) *RemoteAddressWeb {
	return &RemoteAddressWeb{handler: h}
}

func (i *RemoteAddressWeb) Before(response trama.Response, r *http.Request) {
	clientAddress, err := remoteAddress(r)
	if err != nil {
		log.Notice(err.Error())
		response.ExecuteTemplate("500.html", data.NewInternalServerError("N/A"))
		return
	}

	i.handler.SetRemoteAddress(clientAddress)
}

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////

func remoteAddress(r *http.Request) (net.IP, error) {
	var clientAddress string

	xff := r.Header.Get("X-Forwarded-For")
	xff = strings.TrimSpace(xff)

	if len(xff) > 0 {
		xffParts := strings.Split(xff, ",")
		if len(xffParts) == 1 {
			clientAddress = strings.TrimSpace(xffParts[0])
		} else if len(xffParts) > 1 {
			clientAddress = strings.TrimSpace(xffParts[len(xffParts)-2])
		}

	} else {
		clientAddress = strings.TrimSpace(r.Header.Get("X-Real-IP"))
	}

	if len(clientAddress) > 0 {
		if address := net.ParseIP(clientAddress); address != nil {
			return net.ParseIP(clientAddress), nil
		}
	}

	clientAddress, _, err := net.SplitHostPort(r.RemoteAddr)
	return net.ParseIP(clientAddress), err
}
