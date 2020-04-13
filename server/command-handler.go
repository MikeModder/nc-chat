package server

import (
	"fmt"
)

func NewCommandHandler() *CommandHandler {
	c := &CommandHandler{
		Commands: make(map[string]*Command),
	}
	c.AddCommand("help", "Get a list of commands and their functions", ModeUser, c.helpCommand)
	return c
}

func (c *CommandHandler) ExecuteCommand(s *Server, i *Client, cmd string, args []string) (bool, string) {
	// TODO: check for aliases too, once everything is working
	if _, ok := c.Commands[cmd]; ok {
		// Check user's mode (permissoins)
		if c.Commands[cmd].Mode > i.Mode {
			return false, "mode is too low"
		}

		c.Commands[cmd].Run(s, i, args)
		return true, ""
	}

	// returns false on error, with a string message
	// TODO: maybe bubble up errors when needed?
	return false, "command not found"

}

func (c *CommandHandler) AddCommand(name, description string, mode int, run CommandRunFunc) {
	c.Commands[name] = &Command{
		Name: name,
		Description: description,
		Mode: mode,
		Run: run,
	}
}

// func (c *CommandHandler) AddAlias(alias, command string) {}

func (c *CommandHandler) helpCommand(s *Server, i *Client, args []string) {
	// Did the user specify a command to get help on?
	if len(args) >= 1 {
		// Yes, show help about that command
		// But first make sure it's an actual command
		if _, ok := c.Commands[args[0]]; !ok {
			i.SendSystemMessage("command not found!")
			return
		}

		i.SendSystemMessage(fmt.Sprintf("Command %s: %s (mode: %d)", args[0], c.Commands[args[0]].Description, c.Commands[args[0]].Mode))
		return
	}

	// No, show them a list of commands
	for name, cmd := range c.Commands {
		i.SendSystemMessage(fmt.Sprintf("Command %s: %s (mode: %d)", name, cmd.Description, cmd.Mode))
	}
	return
}