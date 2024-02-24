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

	if c.Command == "list" {
		list()
	}

}

// list retrieves all the environment files from the default environment folder
// and prints "List" to the console.
func list() {
	fmt.Println(">> Listing environment configurations...")
	envFiles := manager.GetEnvFiles(&manager.DEFAULT_ENV_FOLDER)
	fmt.Printf("\n>> Found %d environment configurations\n", len(envFiles))
	for _, e := range envFiles {
		fmt.Printf("\t> %s\n", e.Identifier())
	}
}

// get retrieves the environment file identified by the given identifier,
// restores it using the provided secret, and prints "Get" to the console.
func get(identifier string, s ISecret) {
	fmt.Printf("\n>> Getting environment configuration for %s...\n", identifier)
	e := manager.GetEnvFile(identifier, &manager.DEFAULT_ENV_FOLDER)
	secret := s.GetSecret()
	manager.RestoreEnvFile(e, secret)
	fmt.Printf("\t> Environment configuration restored as %s", e.RestoreAs())
}

// init_ initializes the environment by reading the environment file from the given file path,
// saving the environment variables along with the secret provided by ISecret interface.
// It panics if no file path is provided.
func init_(filePath string, s ISecret) {
	fmt.Printf("\n>> Initializing environment configuration from %s...\n", filePath)
	if filePath != "" {
		manager.CreateInitFolderIfNotExist(&manager.DEFAULT_ENV_FOLDER)
		e := manager.ReadEnvFile(filePath, false)
		fmt.Println("\t> Saving environment configuration...")
		secret := s.GetSecret()
		manager.SaveEnvFile(e, secret, &manager.DEFAULT_ENV_FOLDER)
		fmt.Println("\t> Environment configuration saved")
	} else {
		panic("No file path provided")
	}
}
