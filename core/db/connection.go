package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
	"github.com/rafaeljusto/druns/core"
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
	Driver = "mysql"
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

	// Collation specifies the charset used in our system. The mysql server is responsable for
	// converting between UTF-8 and ISO-8859-1
	connParams := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?collation=utf8_general_ci&parseTime=true&loc=UTC",
		user,
		password,
		host,
		port,
		name,
	)

	sqlDB, err := sql.Open(Driver, connParams)
	if err != nil {
		return core.NewError(err)
	}

	// sql.Open returns a concrete type and not an interface, so we must set the global variable only
	// after to make sure that it is initialized
	DB = &database{sqlDB}
	DB.SetMaxIdleConns(16)
	return nil
}
