package main

import (
	"fmt"
	"os"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/user"
	"github.com/rafaeljusto/shelter/config"
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
