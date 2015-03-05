package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"runtime"

	"github.com/rafaeljusto/druns/core/log"
)

var (
	NotFound = Error{
		Err: errors.New("Object not found"),
	}
)

type Error struct {
	Err  error
	File string
	Line int
}

func New(err error) error {
	if err == nil {
		return nil

	} else if err == sql.ErrNoRows {
		return NotFound
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	if _, ok := err.(Error); ok {
		log.Warningf("Trying to create a errors.Error from another errors.Error on '%s', line '%d'!",
			file, line)
		return err
	}

	return Error{
		Err:  err,
		File: file,
		Line: line,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Err.Error())
}
