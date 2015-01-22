package interceptor

import (
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/web/templates/data"
)

type sqler interface {
	Logger() *log.Logger
	Tx() db.Transaction
	SetTx(tx db.Transaction)
	HTTPId() string
}

////////////////////////////////////////////////////////////
/////////////////////// AJAX ///////////////////////////////
////////////////////////////////////////////////////////////

type DatabaseAJAX struct {
	handler sqler
}

func NewDatabaseAJAX(h sqler) *DatabaseAJAX {
	return &DatabaseAJAX{handler: h}
}

func (i *DatabaseAJAX) Before(w http.ResponseWriter, r *http.Request) {
	tx, err := db.DB.Begin()
	if err != nil {
		i.handler.Logger().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i.handler.SetTx(tx)
}

func (i *DatabaseAJAX) After(w http.ResponseWriter, r *http.Request) {
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

////////////////////////////////////////////////////////////
/////////////////////// WEB ////////////////////////////////
////////////////////////////////////////////////////////////

type DatabaseWeb struct {
	handler sqler
}

func NewDatabaseWeb(h sqler) *DatabaseWeb {
	return &DatabaseWeb{handler: h}
}

func (i *DatabaseWeb) Before(response trama.Response, r *http.Request) {
	tx, err := db.DB.Begin()
	if err != nil {
		i.handler.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerErrorData(i.handler.HTTPId()))
		return
	}

	i.handler.SetTx(tx)
}

func (i *DatabaseWeb) After(response trama.Response, r *http.Request) {
	if i.handler.Tx() == nil {
		i.handler.Logger().Warning("Nil transaction found")
		return
	}

	if response.TemplateName() == "500.html" {
		if err := i.handler.Tx().Rollback(); err != nil {
			i.handler.Logger().Error(err)
		}

	} else {
		if err := i.handler.Tx().Commit(); err != nil {
			i.handler.Logger().Error(err)
			response.ExecuteTemplate("500.html", data.NewInternalServerErrorData(i.handler.HTTPId()))
		}
	}
}
