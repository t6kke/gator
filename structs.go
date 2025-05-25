package main

import (
	"fmt"
)

type command struct {
	name  string
	args  []string
}

type commands struct {
	commands  map[string]func(*state, command) error
}

func newCommandsStruct() commands {
	return commands{
		commands: make(map[string]func(*state, command) error, 1),
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.commands[cmd.name]
	if !ok {
		//return errors.New("unknown command: %s", cmd.args[0])
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	err := f(s, cmd)
	return err
}





