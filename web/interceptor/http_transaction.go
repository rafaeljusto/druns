package interceptor

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/log"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type httpTransactioner interface {
	RemoteAddress() net.IP
	SetLogger(log.Logger)
	SetHTTPId(id string)
}

////////////////////////////////////////////////////////////
/////////////////////// AJAX ///////////////////////////////
////////////////////////////////////////////////////////////

type HTTPTransactionAJAX struct {
	trama.NopAJAXInterceptor
	handler httpTransactioner
}

func NewHTTPTransactionAJAX(h httpTransactioner) *HTTPTransactionAJAX {
	return &HTTPTransactionAJAX{handler: h}
}

func (i HTTPTransactionAJAX) Before(w http.ResponseWriter, r *http.Request) {
	identifier := fmt.Sprintf("%s %05d", i.handler.RemoteAddress(), random.Intn(99999))
	i.handler.SetLogger(log.NewLogger(identifier))
}

////////////////////////////////////////////////////////////
/////////////////////// WEB ////////////////////////////////
////////////////////////////////////////////////////////////

type HTTPTransactionWeb struct {
	trama.NopWebInterceptor
	handler httpTransactioner
}

func NewHTTPTransactionWeb(h httpTransactioner) *HTTPTransactionWeb {
	return &HTTPTransactionWeb{handler: h}
}

func (i HTTPTransactionWeb) Before(response trama.Response, r *http.Request) {
	id := fmt.Sprintf("%05d", random.Intn(99999))
	i.handler.SetHTTPId(id)
	i.handler.SetLogger(log.NewLogger(i.handler.RemoteAddress().String() + " " + id))
}
