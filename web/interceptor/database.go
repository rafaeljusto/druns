package interceptor

import (
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/log"
)

type sqler interface {
	Logger() *log.Logger
	Tx() db.Transaction
	SetTx(tx db.Transaction)
}

type Database struct {
	handler sqler
}

func NewDatabase(h sqler) *Database {
	return &Database{handler: h}
}

func (i *Database) Before(w http.ResponseWriter, r *http.Request) {
	tx, err := db.DB.Begin()
	if err != nil {
		i.handler.Logger().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i.handler.SetTx(tx)
}

func (i *Database) After(w http.ResponseWriter, r *http.Request) {
	if i.handler.Tx() == nil {
		i.handler.Logger().Warning("Nil transaction found")
		return
	}

	if responseWriter, ok := w.(*trama.BufferedResponseWriter); ok {
		if responseWriter.Status() >= 200 && responseWriter.Status() < 400 {
			if err := i.handler.Tx().Commit(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				i.handler.Logger().Error(err)
			}
		} else {
			if err := i.handler.Tx().Rollback(); err != nil {
				i.handler.Logger().Error(err)
			}
		}

	} else {
		i.handler.Logger().Warning("Unknown ResponseWriter")

		if err := i.handler.Tx().Rollback(); err != nil {
			i.handler.Logger().Error(err)
		}
	}
}
