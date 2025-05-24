package main

import (
	"fmt"
	"time"
	"context"
	"github.com/google/uuid"
	"github.com/t6kke/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No username provided, username is required --- Usage: %s <name>", cmd.name)
	}

	new_ctx := context.Background()
	_, err := s.dbq.GetUser(new_ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = s.conf.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User '%s' has been successfully configured for session\n", cmd.args[0])
	return nil
}



func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No username provided, username is required --- Usage: %s <name>", cmd.name)
	}

	new_uuid := uuid.New()
	current_time := time.Now()
	user_name := cmd.args[0]

	new_user := database.CreateUserParams{
		ID:        new_uuid,
		CreatedAt: current_time,
		UpdatedAt: current_time,
		Name:      user_name,
	}

	new_ctx := context.Background()
	user, err := s.dbq.GetUser(new_ctx, user_name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return fmt.Errorf("%w", err)
	}

	if user.Name != "" {
		fmt.Errorf("User with name already exists")
	}

	user, err = s.dbq.CreateUser(new_ctx, new_user)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = s.conf.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User '%s' has been successfully added to the database\n", user_name)
	fmt.Printf("DEBUG --- uuid: '%v' --- timestamp: '%v' --- user: '%s'\n", new_user.ID, new_user.CreatedAt, new_user.Name)

	return nil
}
