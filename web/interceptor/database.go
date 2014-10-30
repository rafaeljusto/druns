package interceptor

import (
	"net/http"
	"sync"
	"time"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/web/config"
	"gopkg.in/mgo.v2"
)

var (
	session     *mgo.Session
	sessionLock sync.Mutex

	Timeout = time.Duration(5) * time.Second
)

type DatabaseHandler interface {
	SetDBSession(*mgo.Session)
	DBSession() *mgo.Session
	SetDB(*mgo.Database)
	Logger() *log.Logger
}

type Database struct {
	handler DatabaseHandler
}

func NewDatabase(h DatabaseHandler) *Database {
	return &Database{handler: h}
}

func (i *Database) Before(w http.ResponseWriter, r *http.Request) {
	err := initializeSession(config.DrunsConfig.Database.URI, config.DrunsConfig.Database.Name)
	if err != nil {
		i.handler.Logger().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	localSession := session.Copy()
	database := localSession.DB(config.DrunsConfig.Database.Name)

	i.handler.SetDBSession(localSession)
	i.handler.SetDB(database)
}

func (i *Database) After(w http.ResponseWriter, r *http.Request) {
	i.handler.DBSession().Close()
}

func initializeSession(uri string, databaseName string) error {
	sessionLock.Lock()
	defer sessionLock.Unlock()

	if session != nil {
		return nil
	}

	dialInfo := mgo.DialInfo{
		Addrs:    []string{uri},
		Timeout:  Timeout,
		Database: databaseName,
	}

	var err error
	if session, err = mgo.DialWithInfo(&dialInfo); err != nil {
		return core.NewError(err)
	}
	return nil
}
