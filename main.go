package main

import (
	"os"
	"fmt"
	"log"
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/t6kke/gator/internal/config"
	"github.com/t6kke/gator/internal/database"
)

type state struct {
	dbq   *database.Queries
	conf  *config.Config
}

func main() {
	cnf, err := config.ReadConfig() //in Go, to take the address of a value returned from a function, you'll need to store the return value in a variable first.
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	var s state
	s.conf = &cnf

	db, err := sql.Open("postgres", s.conf.DB_url)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	s.dbq = dbQueries

	commands := newCommandsStruct()
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	raw_args := os.Args
	args := raw_args[1:]
	if len(args) == 0 {
		log.Fatalf("no arguments provided --- Usage: cli <command> [args...]")
		os.Exit(1)
	}

	err = commands.run(&s, command{args[0], args[1:]})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}



func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		new_ctx := context.Background()
		current_user := s.conf.Current_user_name
		user, err := s.dbq.GetUser(new_ctx, current_user)
		if err != nil && err.Error() != "sql: no rows in result set" {
			return fmt.Errorf("%w", err)
		}
		if err != nil {
			return fmt.Errorf("Current user '%s' not found in database", current_user)
		}
		return handler(s, cmd, user)
	}
}
