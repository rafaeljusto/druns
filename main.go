package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/core/protocol"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/handler"
	"github.com/rafaeljusto/druns/web/tr"
)

var (
	Logger = log.NewLogger("system")
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config-file>\n", os.Args[0])
		return
	}

	if err := config.LoadConfig(os.Args[1]); err != nil {
		Logger.Critf("Error loading configuration file. Details: %s", err)
		return
	}

	if err := initializeLogger(); err != nil {
		Logger.Critf("Error initializing logger. Details: %s", err)
		return
	}

	if err := initializeProtocolTranslations(); err != nil {
		Logger.Critf("Error initializing protocol translations. Details: %s", err)
		return
	}

	if err := initializeWebTranslations(); err != nil {
		Logger.Critf("Error initializing web translations. Details: %s", err)
		return
	}

	if err := initializeTrama(); err != nil {
		Logger.Critf("Error initializing trama framework. Details: %s", err)
		return
	}

	serverConfig := config.DrunsConfig.Server
	server := http.Server{
		Addr:        fmt.Sprintf("%s:%d", serverConfig.IP, serverConfig.Port),
		Handler:     handler.Mux,
		ReadTimeout: 5 * time.Second,
	}

	panic(server.ListenAndServe())
}

func initializeLogger() error {
	logAddr := net.JoinHostPort(config.DrunsConfig.Log.Host, strconv.Itoa(config.DrunsConfig.Log.Port))
	if err := log.Connect("druns", logAddr); err != nil {
		return core.NewError(err)
	}
	return nil
}

func initializeProtocolTranslations() error {
	translationsPath := path.Join(config.DrunsConfig.Paths.Home,
		config.DrunsConfig.Paths.ProtocolTranslations)

	return protocol.LoadTranslations(translationsPath)
}

func initializeWebTranslations() error {
	translationsPath := path.Join(config.DrunsConfig.Paths.Home,
		config.DrunsConfig.Paths.WebTranslations)

	return tr.LoadTranslations(translationsPath)
}

func initializeTrama() error {
	handler.Mux.Recover = func(r interface{}) {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		Logger.Critf("Panic detected. Details: %v\n%s", r, buf)
	}
	handler.Mux.SetTemplateDelims("[[", "]]")

	groupSet := trama.NewTemplateGroupSet(nil)
	for _, language := range config.DrunsConfig.Languages {
		templates := config.DrunsConfig.HTMLTemplates(language, "global")

		groupSet.Insert(trama.TemplateGroup{
			Name:  language,
			Files: templates,
		})
	}
	handler.Mux.GlobalTemplates = groupSet

	handler.Mux.RegisterStatic("/assets", http.Dir(config.DrunsConfig.Paths.WebAssets))

	if err := handler.Mux.ParseTemplates(); err != nil {
		return core.NewError(err)
	}

	return nil
}
