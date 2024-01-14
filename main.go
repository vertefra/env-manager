package main

import (
	"fmt"

	"github.com/vertefra/env/cli"
	"github.com/vertefra/env/manager"
)

type ISecret interface {
	GetSecret() string
}

// main is the entry point of the program.
func main() {
	// Take input
	c := cli.ParseArgs()
	fmt.Print("Command: ")
	fmt.Println(c.Command)
	fmt.Print("File: ")
	fmt.Println(c.FromFile)

	s := manager.InitSecret()

	if c.Command == "add" {
		if c.FromFile == "" {
			panic("No file path provided")
		}
		init_(c.FromFile, &s)
	}

	if c.Command == "get" {
		if c.Identifier == "" {
			panic("No identifier provided")
		}

		get(c.Identifier, &s)
	}

}

// get retrieves the environment file identified by the given identifier,
// restores it using the provided secret, and prints "Get" to the console.
func get(identifier string, s ISecret) {
	fmt.Println("Get")
	e := manager.GetEnvFile(identifier)
	secret := s.GetSecret()
	manager.RestoreEnvFile(e, secret)
}

// init_ initializes the environment by reading the environment file from the given file path,
// saving the environment variables along with the secret provided by ISecret interface.
// It panics if no file path is provided.
func init_(filePath string, s ISecret) {
	fmt.Print("\nAdd... ")
	fmt.Print(filePath)

	if filePath != "" {
		manager.CreateInitFolderIfNotExist()
		e := manager.ReadEnvFile(filePath)
		secret := s.GetSecret()
		manager.SaveEnvFile(e, secret)
	} else {
		panic("No file path provided")
	}
}
