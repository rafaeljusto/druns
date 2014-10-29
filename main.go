package main

import (
	"fmt"
	"net/http"

	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/handler"
)

func main() {
	serverConfig := config.DrunsConfig.Server
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", serverConfig.IP, serverConfig.Port),
		Handler: handler.Mux,
	}

	panic(server.ListenAndServe())
}
