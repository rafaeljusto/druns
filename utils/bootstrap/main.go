package main

import (
	"flag"
	"fmt"
	"net"
	"net/mail"
	"os"

	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/model"
	"github.com/rafaeljusto/druns/web/config"
)

var (
	configPath string
	name       string
	email      string
	password   string
)

func init() {
	flag.StringVar(&configPath, "config", "", "configuration file path of the web server")
	flag.StringVar(&name, "name", "", "administrator's name")
	flag.StringVar(&email, "email", "", "administrator's e-mail")
	flag.StringVar(&password, "password", "", "administrator's password")
}

func main() {
	flag.Parse()

	if !checkInputs() {
		usage()
		os.Exit(1)
	}

	if err := config.LoadConfig(configPath); err != nil {
		fmt.Printf("Error loading configuration file. Details: %s\n", err)
		os.Exit(2)
	}

	if err := initializeDatabase(); err != nil {
		fmt.Printf("Error initializing database. Details: %s\n", err)
		os.Exit(3)
	}
	defer db.DB.Close()

	e, err := mail.ParseAddress(email)
	if err != nil {
		fmt.Printf("Invalid e-mail. Details: %s\n", err)
		os.Exit(4)
	}

	user := model.User{
		Name:     name,
		Email:    e,
		Password: password,
	}

	addr, err := localAddress()
	if err != nil {
		fmt.Printf("Error retrieving the local address. Details: %s\n", err)
		os.Exit(5)

	} else if addr == nil {
		fmt.Printf("Couldn't retrieve the local address")
		os.Exit(6)
	}

	tx, err := db.DB.Begin()
	if err != nil {
		fmt.Printf("Error starting database transaction. Details: %s\n", err)
		os.Exit(7)
	}

	userDAO := dao.NewUser(db.DB, addr, "BOOTSTRAP")
	if users, err := userDAO.FindAll(); err != nil {
		fmt.Printf("Error retrieving users. Details: %s\n", err)
		os.Exit(8)

	} else if len(users) > 0 {
		fmt.Println("Database already initialized")
		return
	}

	if err := userDAO.Save(&user); err != nil {
		fmt.Printf("Error saving the user. Details: %s\n", err)
		os.Exit(9)
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("Error commiting database transaction. Details: %s\n", err)
		os.Exit(10)
	}

	fmt.Println("Administrator created successfully")
}

func usage() {
	fmt.Printf("Usage: %s <-config 'path'> <-email 'email'> "+
		"<-name 'name'> <-password 'password'>\n", os.Args[0])
	flag.PrintDefaults()
}

func checkInputs() bool {
	ok := true

	if len(configPath) == 0 {
		fmt.Println("Configuration path not informed!")
		ok = false
	}

	if len(name) == 0 {
		fmt.Println("Name not informed!")
		ok = false
	}

	if len(email) == 0 {
		fmt.Println("E-mail not informed!")
		ok = false
	}

	if len(password) == 0 {
		fmt.Println("Password not informed!")
		ok = false
	}

	return ok
}

func initializeDatabase() error {
	return db.Start(
		config.DrunsConfig.Database.Host,
		config.DrunsConfig.Database.Port,
		config.DrunsConfig.Database.User,
		config.DrunsConfig.Database.Password,
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
