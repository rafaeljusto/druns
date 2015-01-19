package interceptor

import (
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/core/protocol"
	"github.com/rafaeljusto/druns/web/security"
)

type authSecreter interface {
	Logger() *log.Logger
	AuthSecret(secretId string) (string, error)
	SetHandle(handle string)
	SetMessage(protocol.Translator)
}

type Auth struct {
	trama.NopAJAXInterceptor
	handler authSecreter
}

func NewAuth(h authSecreter) *Auth {
	return &Auth{
		handler: h,
	}
}

func (i *Auth) Before(w http.ResponseWriter, r *http.Request) {
	authorized, handle, msg, err := security.CheckAuthorization(r, i.handler.AuthSecret)
	i.handler.SetHandle(handle)

	if err != nil {
		i.handler.Logger().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if msg != nil {
		w.WriteHeader(http.StatusBadRequest)
		i.handler.SetMessage(msg)
	} else if !authorized {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
