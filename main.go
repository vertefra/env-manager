package main

import (
	"fmt"
	"os"

	"github.com/thinktwiceco/env-manager/cli"
	"github.com/thinktwiceco/env-manager/manager"
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

	if c.Command == "create" {
		if c.Identifier == "" {
			panic("No identifier provided")
		}
		if c.FromFile == "" {
			panic("No file path provided")
		}
		create(c.FromFile, c.Identifier, c.RestoreAs, &s)
	}

	if c.Command == "remove" {
		if c.Identifier == "" {
			panic("No identifier provided")
		}
		remove(c.Identifier)
	}

}

// list retrieves all the environment files from the default environment folder
// and prints "List" to the console.
func list() {
	fmt.Println(">> Listing environment configurations...")
	envFiles, err := manager.GetEnvFiles(&manager.DEFAULT_ENV_FOLDER)
	if err != nil {
		fmt.Printf("Error getting environment files: %v\n", err)
		return
	}
	fmt.Printf("\n>> Found %d environment configurations\n", len(envFiles))
	for _, e := range envFiles {
		fmt.Printf("\t> %s\n", e.Identifier())
	}
}

// get retrieves the environment file identified by the given identifier,
// restores it using the provided secret, and prints "Get" to the console.
func get(identifier string, s ISecret) {
	fmt.Printf("\n>> Getting environment configuration for %s...\n", identifier)
	e, err := manager.GetEnvFile(identifier, &manager.DEFAULT_ENV_FOLDER)
	if err != nil {
		fmt.Printf("Error getting environment file: %v\n", err)
		return
	}
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
		f, err := manager.GetOrCreateFolder(&manager.DEFAULT_ENV_FOLDER)
		if err != nil {
			panic(err)
		}
		e := manager.ReadEnvFile(filePath)
		fmt.Println("\t> Saving environment configuration...")
		secret := s.GetSecret()
		f.AddFileIdentifier(manager.EnvFilePath(filePath), manager.EnvFileIdentifier(e.Identifier()))
		manager.SaveEnvFile(e, secret, &manager.DEFAULT_ENV_FOLDER)
		fmt.Println("\t> Environment configuration saved")
	} else {
		panic("No file path provided")
	}
}

// create creates a new environment file from a source file without headers.
// It uses InitEnvFile to create the env file with the given identifier and restoreAs.
func create(filePath string, identifier string, restoreAs string, s ISecret) {
	fmt.Printf("\n>> Creating environment configuration '%s' from %s...\n", identifier, filePath)

	f, err := manager.GetOrCreateFolder(&manager.DEFAULT_ENV_FOLDER)
	if err != nil {
		panic(err)
	}

	// Set default restoreAs if not provided
	if restoreAs == "" {
		restoreAs = manager.DEFAULT_RESTORE_AS
	}

	// Create new env file using InitEnvFile
	e := manager.InitEnvFile(identifier, restoreAs)

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error reading file: %v", err))
	}

	// Set the content (this will add headers automatically)
	e.SetContent(string(content))

	fmt.Println("\t> Saving environment configuration...")
	secret := s.GetSecret()
	f.AddFileIdentifier(manager.EnvFilePath(filePath), manager.EnvFileIdentifier(identifier))
	manager.SaveEnvFile(e, secret, &manager.DEFAULT_ENV_FOLDER)
	fmt.Printf("\t> Environment configuration '%s' saved\n", identifier)
}

// remove deletes an environment configuration from the env-manager folder.
func remove(identifier string) {
	fmt.Printf("\n>> Removing environment configuration '%s'...\n", identifier)

	f, err := manager.GetOrCreateFolder(&manager.DEFAULT_ENV_FOLDER)
	if err != nil {
		panic(err)
	}

	// Find and remove the file from manifest
	var filePathToRemove manager.EnvFilePath
	for filePath, id := range f.GetIdentifiers() {
		if string(id) == identifier {
			filePathToRemove = filePath
			break
		}
	}

	if filePathToRemove == "" {
		fmt.Printf("Error: identifier '%s' not found\n", identifier)
		return
	}

	// Remove from manifest
	err = f.EvictFileIdentifier(filePathToRemove)
	if err != nil {
		panic(fmt.Sprintf("Error removing from manifest: %v", err))
	}

	// Remove the actual encrypted file
	filePath := fmt.Sprintf("%s/%s%s", manager.DEFAULT_ENV_FOLDER, manager.SAVED_PREFIX, identifier)
	err = os.Remove(filePath)
	if err != nil {
		fmt.Printf("Warning: could not remove file %s: %v\n", filePath, err)
	}

	fmt.Printf("\t> Environment configuration '%s' removed\n", identifier)
}
