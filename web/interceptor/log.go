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

type logger interface {
	RemoteAddress() net.IP
	SetLogger(log.Logger)
}

type Log struct {
	trama.NopAJAXInterceptor
	handler logger
}

func NewLog(h logger) *Log {
	return &Log{handler: h}
}

func (i Log) Before(w http.ResponseWriter, r *http.Request) {
	identifier := fmt.Sprintf("%s %05d", i.handler.RemoteAddress(), random.Intn(99999))
	i.handler.SetLogger(log.NewLogger(identifier))
}
