package main

import (
	"crypto/tls"
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
	"github.com/rafaeljusto/druns/core/db"
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
		fmt.Printf("Error loading configuration file. Details: %s\n", err)
		Logger.Critf("Error loading configuration file. Details: %s", err)
		return
	}

	if err := initializeLogger(); err != nil {
		fmt.Printf("Error initializing logger. Details: %s\n", err)
		Logger.Critf("Error initializing logger. Details: %s", err)
		return
	}

	if err := initializeProtocolTranslations(); err != nil {
		fmt.Printf("Error initializing protocol translations. Details: %s\n", err)
		Logger.Critf("Error initializing protocol translations. Details: %s", err)
		return
	}

	if err := initializeWebTranslations(); err != nil {
		fmt.Printf("Error initializing web translations. Details: %s\n", err)
		Logger.Critf("Error initializing web translations. Details: %s", err)
		return
	}

	if err := initializeTrama(); err != nil {
		fmt.Printf("Error initializing trama framework. Details: %s\n", err)
		Logger.Critf("Error initializing trama framework. Details: %s", err)
		return
	}

	if err := initializeDatabase(); err != nil {
		fmt.Printf("Error initializing database. Details: %s\n", err)
		Logger.Critf("Error initializing database. Details: %s", err)
		return
	}
	defer db.DB.Close()

	certPath, privKeyPath := config.DrunsConfig.TLS()
	cert, err := tls.LoadX509KeyPair(certPath, privKeyPath)
	if err != nil {
		fmt.Printf("Error on TLS setup. Details: %s\n", err)
		Logger.Critf("Error on TLS setup. Details: %s", err)
		return
	}

	serverConfig := config.DrunsConfig.Server
	ipAndPort := net.JoinHostPort(serverConfig.IP, strconv.Itoa(serverConfig.Port))
	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}}

	ln, err := tls.Listen("tcp", ipAndPort, &tlsConfig)
	if err != nil {
		fmt.Printf("Error listening to interface. Details: %s\n", err)
		Logger.Critf("Error listening to interface. Details: %s", err)
		return
	}

	server := http.Server{
		Handler:     handler.Mux,
		ReadTimeout: 5 * time.Second,
	}

	panic(server.Serve(ln))
}

func initializeLogger() error {
	logAddr := net.JoinHostPort(config.DrunsConfig.Log.Host, strconv.Itoa(config.DrunsConfig.Log.Port))
	if err := log.Connect("druns", logAddr); err != nil {
		return core.NewError(err)
	}
	return nil
}

func initializeProtocolTranslations() error {
	translationsPath := path.Join(config.DrunsConfig.Home,
		config.DrunsConfig.Paths.ProtocolTranslations)

	return protocol.LoadTranslations(translationsPath)
}

func initializeWebTranslations() error {
	translationsPath := path.Join(config.DrunsConfig.Home,
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

	handler.Mux.RegisterStatic("/assets", http.Dir(path.Join(config.DrunsConfig.Home,
		config.DrunsConfig.Paths.WebAssets)))

	if err := handler.Mux.ParseTemplates(); err != nil {
		return core.NewError(err)
	}

	return nil
}

func initializeDatabase() error {
	return db.Start(
		config.DrunsConfig.Database.Host,
		config.DrunsConfig.Database.Port,
		config.DrunsConfig.Database.User,
		config.DecryptPassword(config.DrunsConfig.Database.Password),
		config.DrunsConfig.Database.Name,
	)
}
