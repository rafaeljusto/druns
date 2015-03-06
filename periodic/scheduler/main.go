package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/group"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/core/password"
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
		Logger.Critf("Error loading configuration file. Details: %s", err)
		return
	}

	if err := initializeLogger(); err != nil {
		fmt.Printf("Error initializing logger. Details: %s\n", err)
		Logger.Critf("Error initializing logger. Details: %s", err)
		return
	}

	if err := initializeDatabase(); err != nil {
		Logger.Critf("Error initializing database. Details: %s", err)
		return
	}
	defer db.DB.Close()

	tx, err := db.DB.Begin()
	if err != nil {
		Logger.Errorf("Error creating a database transaction. Details: %s", err)
		return
	}

	_, err = group.NewService(tx).FindAll()
	if err != nil {
		Logger.Errorf("Error retrieving groups. Details: %s", err)
		return
	}
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
