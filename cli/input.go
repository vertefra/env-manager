package cli

import (
	"github.com/alexflint/go-arg"
)

type CommandType struct {
	Command    string `arg:"positional,required"`
	FromFile   string `arg:"-f"`
	Identifier string `arg:"-i"`
}

func (c *CommandType) validateCommand() {
	for _, cmd := range validCommands {
		if c.Command == cmd {
			return
		}
	}

	panic("Invalid command. Valid commands are init and get")
}

var validCommands [4]string = [4]string{"add", "get", "list", "remove"}

func ParseArgs() CommandType {
	var cmd CommandType

	arg.MustParse(&cmd)
	cmd.validateCommand()
	// Check if command is valid

	return cmd
}
