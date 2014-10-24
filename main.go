package main

import (
	"net/http"

	"github.com/rafaeljusto/druns/handler"
)

func main() {
	server := http.Server{
		Handler: handler.Mux,
	}

	panic(server.ListenAndServe())
}
