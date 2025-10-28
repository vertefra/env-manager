package cli

import (
	"github.com/alexflint/go-arg"
)

type CommandType struct {
	Command    string `arg:"positional,required" help:"Command to execute: add, get, list, create, remove"`
	FromFile   string `arg:"-f" help:"Path to environment file"`
	Identifier string `arg:"-i" help:"Unique identifier for the environment configuration"`
	RestoreAs  string `arg:"-r" help:"Filename to restore the environment file as (default: .env)"`
}

func (CommandType) Description() string {
	return `Environment Manager - Securely store and manage environment configurations

Commands:
  add      Add an environment file with headers (requires -f)
  get      Retrieve and restore an environment configuration (requires -i)
  list     List all saved environment configurations
  create   Create environment configuration from a file without headers (requires -f, -i)
  remove   Remove an environment configuration (requires -i)

Examples:
  env-manager add -f .env.local
  env-manager create -f secrets.txt -i production -r .env.prod
  env-manager get -i production
  env-manager list
  env-manager remove -i production`
}

func (c *CommandType) validateCommand() {
	for _, cmd := range validCommands {
		if c.Command == cmd {
			return
		}
	}

	panic("Invalid command. Valid commands are add, get, list, remove, and create")
}

var validCommands [5]string = [5]string{"add", "get", "list", "remove", "create"}

func ParseArgs() CommandType {
	var cmd CommandType

	arg.MustParse(&cmd)
	cmd.validateCommand()
	// Check if command is valid

	return cmd
}
