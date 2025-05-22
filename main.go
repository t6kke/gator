package main

import (
	"os"
	"fmt"
	"github.com/t6kke/gator/internal/config"
)

func main() {
	var s state
	cnf := config.ReadConfig() //in Go, to take the address of a value returned from a function, you'll need to store the return value in a variable first.
	s.conf = &cnf

	commands := newCommandsStruct()
	commands.register("login", handlerLogin)

	raw_args := os.Args
	args := raw_args[1:]
	if len(args) == 0 {
		fmt.Println("no arguments provided")
		os.Exit(1)
	}
	fmt.Println(args[0], args[1:])

	err := commands.run(&s, command{args[0], args[1:]})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
