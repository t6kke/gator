package main

import (
	"fmt"
	"github.com/t6kke/gator/internal/config"
)

type state struct {
	conf  *config.Config
}

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

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.commands[cmd.name]
	if ok {
		err := f(s, cmd)
		return err
	} else {
		fmt.Println("unknown command")
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}



