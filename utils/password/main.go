package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rafaeljusto/druns/core/password"
)

var (
	input   string
	command string
)

func init() {
	flag.StringVar(&input, "input", "", "input to encrypt or decrypt")
	flag.StringVar(&command, "command", "encrypt", "encrypt or decrypt")
}

func main() {
	flag.Parse()

	if !checkInputs() {
		usage()
		os.Exit(1)
	}

	var output string
	var err error

	if command == "encrypt" {
		output, err = password.Encrypt(input)
	} else {
		output, err = password.Decrypt(input)
	}

	if err != nil {
		fmt.Printf("Error executing the command '%s'. Details: %s", command, err)
		return
	}

	fmt.Printf("Password %s: %s\n", command, output)
}

func usage() {
	fmt.Printf("Usage: %s <-input 'input'> "+
		"[-command 'command']\n", os.Args[0])
	flag.PrintDefaults()
}

func checkInputs() bool {
	ok := true

	if len(input) == 0 {
		fmt.Println("Password not informed!")
		ok = false
	}

	if command != "encrypt" && command != "decrypt" {
		fmt.Println("Command must be 'encrypt' or 'decrypt'")
		ok = false
	}

	return ok
}
