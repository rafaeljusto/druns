package interceptor

import (
	"net"
	"net/http"
	"strings"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/log"
)

type remoteAddresser interface {
	RemoteAddress() net.IP
	SetRemoteAddress(net.IP)
}

type RemoteAddress struct {
	trama.NopAJAXInterceptor
	handler remoteAddresser
}

func NewRemoteAddress(h remoteAddresser) *RemoteAddress {
	return &RemoteAddress{handler: h}
}

func (i *RemoteAddress) Before(w http.ResponseWriter, r *http.Request) {
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
			i.handler.SetRemoteAddress(net.ParseIP(clientAddress))
			return
		}
	}

	clientAddress, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Notice(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i.handler.SetRemoteAddress(net.ParseIP(clientAddress))
}
