package main

import (
	"fmt"
	"errors"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("No username provided, username is required")
	}
	s.conf.SetUser(cmd.args[0])
	fmt.Printf("User '%s' has been configured for login\n", cmd.args[0])
	return nil
}
