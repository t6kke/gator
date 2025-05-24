package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No username provided, username is required --- Usage: %s <name>", cmd.name)
	}

	err := s.conf.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User '%s' has been successfully configured for session\n", cmd.args[0])
	return nil
}
