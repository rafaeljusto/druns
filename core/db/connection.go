package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"

	_ "github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/lib/pq"
	"github.com/rafaeljusto/druns/core/errors"
)

// We are using an interface to allow overwriting DB with a custom type. Useful for the integration
// tests where we don't want to open many transactions
type Database interface {
	Begin() (Transaction, error)
	Close() error
	Driver() driver.Driver
	Exec(query string, args ...interface{}) (sql.Result, error)
	Ping() error
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
}

// We are using this transaction to avoid commiting twice when reusing a transaction for the
// integration tests
type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Rollback() error
	Commit() error
}

type database struct {
	*sql.DB
}

func (db *database) Begin() (Transaction, error) {
	return db.DB.Begin()
}

var (
	DB     Database
	dbLock sync.Mutex
	Driver = "postgres"
)

func Start(host string, port int, user, password, name string) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	if DB != nil {
		// If we already have a connection, use it!
		if err := DB.Ping(); err == nil {
			return nil
		}
	}

	connParams := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		user,
		password,
		host,
		port,
		name,
	)

	sqlDB, err := sql.Open(Driver, connParams)
	if err != nil {
		return errors.New(err)
	}

	if err := sqlDB.Ping(); err != nil {
		return errors.New(err)
	}

	// sql.Open returns a concrete type and not an interface, so we must set the global variable only
	// after to make sure that it is initialized
	DB = &database{sqlDB}
	DB.SetMaxIdleConns(16)
	return nil
}
