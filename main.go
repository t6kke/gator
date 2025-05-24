package main
import _ "github.com/lib/pq"
import (
	"os"
	"log"
	"database/sql"

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

	dbQueries := database.New(db)
	s.dbq = dbQueries

	commands := newCommandsStruct()
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)

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
