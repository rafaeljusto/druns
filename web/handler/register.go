package handler

import (
	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/log"
)

var (
	Logger = log.NewLogger("system")
	Mux    = trama.New(Logger.Error)
)
