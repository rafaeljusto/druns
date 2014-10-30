package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/handler"
)

var (
	Logger = log.NewLogger("system")
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config-file>\n", os.Args[0])
		return
	}

	if err := initializeConfig(os.Args[1]); err != nil {
		Logger.Error(err)
		return
	}

	if err := initializeLogger(); err != nil {
		Logger.Error(err)
		return
	}

	handler.Mux.Recover = func(r interface{}) {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		Logger.Critf("Panic detected. Details: %v\n%s", r, buf)
	}

	serverConfig := config.DrunsConfig.Server
	server := http.Server{
		Addr:        fmt.Sprintf("%s:%d", serverConfig.IP, serverConfig.Port),
		Handler:     handler.Mux,
		ReadTimeout: 5 * time.Second,
	}

	panic(server.ListenAndServe())
}

func initializeConfig(configFile string) error {
	if err := config.LoadConfig(configFile); err != nil {
		return core.NewError(err)
	}
	return nil
}

func initializeLogger() error {
	logAddr := net.JoinHostPort(config.DrunsConfig.Log.Host, strconv.Itoa(config.DrunsConfig.Log.Port))
	if err := log.Connect("druns", logAddr); err != nil {
		return core.NewError(err)
	}
	return nil
}
