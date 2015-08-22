package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/core/password"
	"github.com/rafaeljusto/druns/core/user"
	"github.com/rafaeljusto/druns/web/config"
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
		return
	}

	if err := initializeLogger(); err != nil {
		fmt.Printf("Error initializing logger. Details: %s\n", err)
		return
	}

	if err := initializeDatabase(); err != nil {
		Logger.Critf("Error initializing database. Details: %s", err)
		return
	}
	defer db.DB.Close()

	addr, err := localAddress()
	if err != nil {
		Logger.Errorf("Error retrieving the local address. Details: %s\n", err)
		return

	} else if addr == nil {
		Logger.Errorf("Couldn't retrieve the local address")
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		Logger.Errorf("Error creating a database transaction. Details: %s", err)
		return
	}

	systemUser, err := user.NewService(tx).SystemUser()
	if err == errors.NotFound {
		Logger.Errorf("System user not found!")
		return

	} else if err != nil {
		Logger.Errorf("Error retrieving system user. Details: %s", err)
		return
	}

	// TODO!
}

func initializeLogger() error {
	logAddr := net.JoinHostPort(config.DrunsConfig.Log.Host, strconv.Itoa(config.DrunsConfig.Log.Port))
	if err := log.Connect("druns", logAddr); err != nil {
		return errors.New(err)
	}
	return nil
}

func initializeDatabase() error {
	dbPassword, err := password.Decrypt(config.DrunsConfig.Database.Password)
	if err != nil {
		return err
	}

	return db.Start(
		config.DrunsConfig.Database.Host,
		config.DrunsConfig.Database.Port,
		config.DrunsConfig.Database.User,
		dbPassword,
		config.DrunsConfig.Database.Name,
	)
}

func localAddress() (net.IP, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return nil, err
	}

	if len(addrs) > 0 {
		return net.ParseIP(addrs[0]), nil
	}

	return nil, nil
}
