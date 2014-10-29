package interceptor

import (
	"net/http"
	"sync"
	"time"

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
}

type Database struct {
	databaseHandler DatabaseHandler
}

func NewDatabase(h DatabaseHandler) *Database {
	return &Database{databaseHandler: h}
}

func (i *Database) Before(w http.ResponseWriter, r *http.Request) {
	err := initializeSession(config.DrunsConfig.Database.URI, config.DrunsConfig.Database.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	localSession := session.Copy()
	database := localSession.DB(config.DrunsConfig.Database.Name)

	i.databaseHandler.SetDBSession(localSession)
	i.databaseHandler.SetDB(database)
}

func (i *Database) After(w http.ResponseWriter, r *http.Request) {
	i.databaseHandler.DBSession().Close()
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
	session, err = mgo.DialWithInfo(&dialInfo)
	return err
}
