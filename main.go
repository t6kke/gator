package main

import (
	"os"
	"log"
	"github.com/t6kke/gator/internal/config"
)

func main() {
	cnf, err := config.ReadConfig() //in Go, to take the address of a value returned from a function, you'll need to store the return value in a variable first.
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	var s state
	s.conf = &cnf

	commands := newCommandsStruct()
	commands.register("login", handlerLogin)

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
